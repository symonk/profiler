package profiler

// ProfileOption is a functional option to configure
// profiler instances.
type ProfileOption func(*Profiler)

// WithProfileFileLocation allows a custom output path for the profile
// file that is written to disk.
func WithProfileFileLocation(path string) ProfileOption {
	return func(p *Profiler) {
		p.profileFolder = path
	}
}

// WithCPUProfiler enables the CPU Profiler.
// CPU Profiling is useful for determining where a program is
// spending CPU cycles (as opposed) to sleeping or waiting for
// IO.
func WithCPUProfiler() ProfileOption {
	return func(p *Profiler) {
		p.profileMode = CPUMode
	}
}

// WithHeapMemoryProfiling enables the Heap Profiler.
// Heap Profiling is useful for determining where memory is
// being allocated and where it is being retained.
func WithHeapMemoryProfiling() ProfileOption {
	return func(p *Profiler) {
		p.profileMode = MemoryHeapMode
	}
}

// WithAllocMemoryProfiling enables the Alloc Profiler.
// Alloc Profiling is useful for determining where memory is
// being allocated and where it is being retained.
// This is different to Heap Profiling as it will show you
// where memory is being allocated, but not necessarily
// where it is being retained.
// This is useful for finding memory leaks.
// This is only available in Go 1.12 and later.
// The rate at which the profiler samples memory allocations
// can be set with the WithMemoryProfilingRate option.
func WithAllocMemoryProfiling() ProfileOption {
	return func(p *Profiler) {
		p.profileMode = MemoryAllocMode
	}
}

// WithMemoryProfilingRate sets the rate at which the
// memory profiler samples memory allocations for both
// Heap and Alloc profiling.  By default this is set to
// the runtime.MemProfileRate value which is 512 * 1024.
// This can be set to a higher value to increase the
// resolution of the memory profile.
func WithMemoryProfilingRate(rate int) ProfileOption {
	return func(p *Profiler) {
		p.memoryProfileRate = rate
	}
}

// WithoutSignalHandling disables the signal handling
// for the profiler.  This is useful for cases where
// you want to handle the signal yourself.
// Be sure to invoke profiler.Stop() yourself in your
// code and handle the os.Exit() yourself etc.
func WithoutSignalHandling() ProfileOption {
	return func(p *Profiler) {
		p.signalHandling = false
	}
}

// WithCallback executes a user defined function when
// clean up occurs.  This function is also fired on
// sigterm handling when the option is enabled.
// Callbacks have access to the underlying *Profiler
// instance, this is typically useful if you wanted to
// do some logic with the profile files that are written
// as the callback is only fired when the profile is
// complete, such as persisting a profile file to a central
// store etc.
func WithCallback(callback CallbackFunc) ProfileOption {
	return func(p *Profiler) {
		p.callback = callback
	}
}

// WithQuietOutput prevents the profiling from writing
// logger events.
func WithQuietOutput() ProfileOption {
	return func(p *Profiler) {
		p.quiet = true
	}
}

// WithTracing enables the tracing profiler.
// Tracing is useful for determining the flow of a program
// and where it is spending time.
// Utilising the trace api within your code can add some
// extra context to the trace output (logs, tasks etc).
// but is not the responsibility of this package.
func WithTracing() ProfileOption {
	return func(p *Profiler) {
		p.profileMode = TraceMode
	}
}

// WithLiveTracing enables live tracing of the program
// as it runs for cases which allow it.  This exposes
// trace data via the runtime/pprof http server.
func WithRealTimeData() ProfileOption {
	return func(p *Profiler) {
		p.live = true
	}
}

// WithMutexFraction sets the rate at which the mutex profiler
// samples mutex contention.  By default this is set to 1.
func WithMutexFraction(rate int) ProfileOption {
	return func(p *Profiler) {
		p.profileMode = MutexMode
	}
}

// WithClockProfiling utilises wall clock profiling powered by
// https://github.com/felixge/fgprof.  This allows you to profile
// both CPU ON and OFF wait in tandem, painting a nice picture.
// Go runtimes built in CPU profiler only displays cpu ON time.
func WithClockProfiling() ProfileOption {
	return func(p *Profiler) {
		p.profileMode = ClockMode
	}
}
