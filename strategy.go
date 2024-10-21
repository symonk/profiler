package profiler

import (
	"runtime"
	"runtime/pprof"
	"runtime/trace"
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
		defer p.profileFile.Close()
	}, nil
}

func heapStrategyFn(p *Profiler) (FinalizerFunc, error) {
	rate := runtime.MemProfileRate
	p.SetProfileFile(MemoryFileName)
	pprof.Lookup(heapProfileName).WriteTo(p.profileFile, 0)
	return func() {
		defer p.profileFile.Close()
		runtime.MemProfileRate = rate
	}, nil
}

func allocStrategyFn(p *Profiler) (FinalizerFunc, error) {
	rate := runtime.MemProfileRate
	p.SetProfileFile(MemoryFileName)
	pprof.Lookup(allocProfileName).WriteTo(p.profileFile, 0)
	return func() {
		defer p.profileFile.Close()
		runtime.MemProfileRate = rate
	}, nil
}

func mutexStrategyFn(p *Profiler) (FinalizerFunc, error) {
	p.SetProfileFile(MutexFileName)
	pprof.Lookup("mutex").WriteTo(p.profileFile, 0)
	return func() {
		defer p.profileFile.Close()
	}, nil
}

func blockStrategyFn(p *Profiler) (FinalizerFunc, error) {
	p.SetProfileFile(BlockFileName)
	pprof.Lookup("block").WriteTo(p.profileFile, 0)
	return func() {
		defer p.profileFile.Close()
	}, nil
}

func goroutineStrategyFn(p *Profiler) (FinalizerFunc, error) {
	p.SetProfileFile(GoroutineFileName)
	pprof.Lookup("goroutine").WriteTo(p.profileFile, 0)
	return func() {
		defer p.profileFile.Close()
	}, nil
}

func threadCreateStrategyFn(p *Profiler) (FinalizerFunc, error) {
	p.SetProfileFile(ThreadCreateFileName)
	pprof.Lookup("threadcreate").WriteTo(p.profileFile, 0)
	return func() {
		defer p.profileFile.Close()
	}, nil
}

func traceStrategyFn(p *Profiler) (FinalizerFunc, error) {
	p.SetProfileFile(TraceFileName)
	trace.Start(p.profileFile)
	return func() {
		defer p.profileFile.Close()
		trace.Stop()
	}, nil
}
