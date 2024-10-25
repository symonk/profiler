package profiler

import (
	"runtime"
	"runtime/pprof"
	"runtime/trace"

	"github.com/felixge/fgprof"
)

// StrategyFunc is the custom type for an implementation
// that controls pre/post profiling setup and teardown.
type StrategyFunc func(p *Profiler) (FinalizerFunc, error)

var StrategyMap = map[Mode]StrategyFunc{
	CPUMode:          cpuStrategyFn,
	MemoryHeapMode:   heapStrategyFn,
	MemoryAllocMode:  allocStrategyFn,
	MutexMode:        mutexStrategyFn,
	BlockMode:        blockStrategyFn,
	GoroutineMode:    goroutineStrategyFn,
	ThreadCreateMode: threadCreateStrategyFn,
	TraceMode:        traceStrategyFn,
	ClockMode:        clockStrategyFn,
}

// cpuStrategyFn handles configuring the cpu profiler and
// deferring it's teardown.
// the output of using this strategy is a `cpu.pprof`
// file written to disk.
func cpuStrategyFn(p *Profiler) (FinalizerFunc, error) {
	p.SetProfileFile(CPUFileName)
	if err := pprof.StartCPUProfile(p.profileFile); err != nil {
		return nil, err
	}
	return func() {
		pprof.StopCPUProfile()
		p.profileFile.Close()
	}, nil
}

func heapStrategyFn(p *Profiler) (FinalizerFunc, error) {
	rate := runtime.MemProfileRate
	p.SetProfileFile(MemoryFileName)
	runtime.MemProfileRate = p.memoryProfileRate
	return func() {
		pprof.Lookup(heapProfileName).WriteTo(p.profileFile, 0)
		runtime.GC()
		p.profileFile.Close()
		runtime.MemProfileRate = rate
	}, nil
}

func allocStrategyFn(p *Profiler) (FinalizerFunc, error) {
	rate := runtime.MemProfileRate
	p.SetProfileFile(MemoryFileName)
	runtime.MemProfileRate = p.memoryProfileRate
	return func() {
		pprof.Lookup(allocProfileName).WriteTo(p.profileFile, 0)
		runtime.GC()
		p.profileFile.Close()
		runtime.MemProfileRate = rate
	}, nil
}

func mutexStrategyFn(p *Profiler) (FinalizerFunc, error) {
	p.SetProfileFile(MutexFileName)
	pprof.Lookup("mutex").WriteTo(p.profileFile, 0)
	return func() {
		p.profileFile.Close()
	}, nil
}

func blockStrategyFn(p *Profiler) (FinalizerFunc, error) {
	p.SetProfileFile(BlockFileName)
	// for now, we do not allow customising the runtime.SetBlockProfileRate
	// if it is useful in future, change is welcome here.
	return func() {
		pprof.Lookup("block").WriteTo(p.profileFile, 0)
		p.profileFile.Close()
		runtime.SetBlockProfileRate(0)
	}, nil
}

func goroutineStrategyFn(p *Profiler) (FinalizerFunc, error) {
	p.SetProfileFile(GoroutineFileName)
	pprof.Lookup("goroutine").WriteTo(p.profileFile, 0)
	return func() {
		p.profileFile.Close()
	}, nil
}

func threadCreateStrategyFn(p *Profiler) (FinalizerFunc, error) {
	p.SetProfileFile(ThreadCreateFileName)
	pprof.Lookup("threadcreate").WriteTo(p.profileFile, 0)
	return func() {
		p.profileFile.Close()
	}, nil
}

func traceStrategyFn(p *Profiler) (FinalizerFunc, error) {
	p.SetProfileFile(TraceFileName)
	trace.Start(p.profileFile)
	return func() {
		p.profileFile.Close()
		trace.Stop()
	}, nil
}

func clockStrategyFn(p *Profiler) (FinalizerFunc, error) {
	p.SetProfileFile(ClockFileName)
	teardown := fgprof.Start(p.profileFile, fgprof.FormatPprof)
	return func() {
		p.profileFile.Close()
		teardown()
	}, nil
}
