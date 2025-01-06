package profiler

import (
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"sync/atomic"
	"syscall"
)

const (
	allocProfileName = "allocs"
	heapProfileName  = "heap"
)

const (
	CPUFileName          = "cpu.pprof"
	MemoryFileName       = "memory.pprof" // Covers heap and alloc
	BlockFileName        = "block.pprof"
	GoroutineFileName    = "goroutine.pprof"
	MutexFileName        = "mutex.pprof"
	ThreadCreateFileName = "threadcreate.pprof"
	TraceFileName        = "trace.out"
	ClockFileName        = "clock.pprof"
)

// FinalizerFunc is a function that is invokved during the teardown period
// of the profiling instance.
type FinalizerFunc func() error

// CallbackFunc is a function that can be supplied with the
// WithCallback option to be executed when the profiling instance
// is performing teardown.  It has access to the *Profiler instance.
type CallbackFunc func(p *Profiler)

type Mode int

const (
	// List of available runtime profiles
	CPUMode Mode = iota
	MemoryHeapMode
	MemoryAllocMode
	BlockMode
	GoroutineMode
	MutexMode
	ThreadCreateMode
	TraceMode
	ClockMode
)

// profileActive is used as a flag to determine if a profiling
// session has begun to manage cases of Start/Stop calls out of
// order, prevent any human error.
var profilingActive uint32

// Profiler encapsulates a profiling instance.
type Profiler struct {
	profileFolder     string
	profileFile       *os.File
	signalHandling    bool
	profileMode       Mode
	memoryProfileRate int
	quiet             bool
	callback          CallbackFunc
	finalizer         FinalizerFunc
	live              bool
	interrupted       bool
}

// New returns a new instance of the Profiler.
func New(options ...ProfileOption) *Profiler {
	p := &Profiler{
		profileFolder:     ".",
		signalHandling:    true,
		memoryProfileRate: runtime.MemProfileRate,
	}
	for _, opt := range options {
		opt(p)
	}
	return p
}

// Stop stops the profiling instance.
// If no profiling instance is active, this function
// will cause an exit.
func (p *Profiler) Stop() {
	if !atomic.CompareAndSwapUint32(&profilingActive, 1, 0) {
		die("profiler instance was not started")
	}
	if err := p.finalizer(); err != nil {
		die(err.Error())
	}
	if p.callback != nil {
		p.callback(p)
	}

	absPath, err := filepath.Abs(p.profileFile.Name())
	if err != nil {
		die(err.Error())
	}
	// Handle reporting data for improved user experience when not running
	// in a suppressed mode.
	extension := filepath.Ext(absPath)
	wasTrace := strings.HasSuffix(absPath, ".out")
	cmd := "go tool pprof -http :8080"
	if wasTrace {
		cmd = "go tool trace"
	}
	p.report("profiling completed.  You can find the %s file at %s", extension, absPath)
	p.report("to view the profile, run `%s %s`", cmd, absPath)
	if p.interrupted {
		p.report("[warning] profiling was interrupted, data may be incomplete")
	}
	if !wasTrace {
		p.report("port can be any ephemeral port you wish to use.")
		p.report("Graph interpretation is outlined here: https://github.com/google/pprof/blob/main/doc/README.md#graphical-reports")
	}
}

// SetProfileFile sets the profile file for the profiler instance.
// not to be confused with the folder location provided by the functional
// options.
func (p *Profiler) SetProfileFile(name string) {
	profileFile, err := CreateProfileFile(p.profileFolder, name)
	if err != nil {
		die(err.Error())
	}
	p.profileFile = profileFile
}

// report writes a formatted log statement to stderr.
// If the WithSuppressedOutput option is provided, this
// will be a no-op.
func (p *Profiler) report(format string, args ...any) {
	if !p.quiet {
		log.Printf(format, args...)
	}
}

// Start starts a new profiling instance.
// If no mode option is provided, the default behavious
// is to perform CPU profiling.
// Start returns the underlying profile instance
// typically deferred in simple scenarios. In more complex
// scenarios keeping a handle to the stop function and calling
// it yourself in some of your own signal handling code for
// example is wise, this should be used with the option:
// WithNoSignalShutdownHandling.
func Start(options ...ProfileOption) *Profiler {

	// Ensure that StartProfiling is not invoked multiple times
	if !atomic.CompareAndSwapUint32(&profilingActive, 0, 1) {
		die("profiler instance has already been started")
	}

	p := New(options...)
	profileFunc, ok := StrategyMap[p.profileMode]
	if !ok {
		die("profiler mode not implemented, this should never happen")
	}
	finalizer, err := profileFunc(p)
	if err != nil {
		die(err.Error())
	}
	p.finalizer = finalizer

	// Register an asynchronous sig term handler if the user
	// has not opted to take full control of exit handling
	// themselves.
	if p.signalHandling {
		go func() {
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
			<-ch
			p.report("sigterm received, performing tear down")
			p.interrupted = true
			p.Stop()
			os.Exit(0)
		}()
	}
	return p
}

// die causes the profiler instance to die with a message.
// This is useful for cases where you want to exit the program
// immediately with a message.
func die(because string) {
	log.Fatalf("profiler instance exited: %s", because)
}
