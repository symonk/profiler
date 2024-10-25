### Profiler 

`profiler` is an abstraction layer around the golang profiling stdlib, namely `runtime/pprof`
and `runtime/trace` to simplify and make profiling configuration easy and fast.  It offers an
assortment of options and configurations as well as various non stdlib options such as the `fgprof`
wall clock option, which pieces together a single `.pprof` file which includes both cpu `ON` and `OFF`
data.  A basic example for each of the profilers is outlined below.

By default `profiler` will listen for `SIGTERM` & `SIGINT` in order to close out the pprof/trace files
correctly, however in some cases this is undesirable.  If you would like full control of cleaning up
then use the `WithoutSignalHandling()` functional option to any invocation of `.Start()`.

-----

### :one: CPU Profiling

Profiling CPU can help identify areas of your code where most CPU cycle execution time is spent.  For
programs with a lot of IO wait etc, the graph may not be extremely useful, trace may be of
better benefit there.

```go
import (
    "github.com/symonk/profiler"
)

func main() {
    defer profiler.Start(profiler.WithCPUProfiler()).Stop()
    /* your code here */
}
```

-----

### :two: Heap Profiling 

Heap profiling reports memory allocation samples, useful for monitoring current and historical memory
usage and to check for potential memory leaks in your program.  Heap profiling tracks the allocation
sites for all live objects in the application memory and for all objects allocated since the program
has started.  By default heap profiling will display live objects, scaled by size.

```go
package main

import (
    "github.com/symonk/profiler"
)

func main() {
    // additional options available here are: `profiler.WithMemoryProfileRate(...)`
    defer profiler.Start(profiler.WithHeapProfiler()).Stop()
    /* your code here */
}
```

-----


### :three: Alloc Profiling

Alloc profiling is essentially the same as heap profiling except rather than the default of live objects
scaled by size, it reports the `-alloc_space` data, that is the total number of bytes allocated since the
program has began (including garbage collected bytes).

```go
package main

import (
    "github.com/symonk/profiler"
)

func main() {
    // additional options available here are: `profiler.WithMemoryProfileRate(...)`
    defer profiler.Start(profiler.WithAllocProfiler()).Stop()
    /* your code here */
}
```

------

### :four: Block Profiling

Block profiling captures how long a program spends off CPU blocked by either a mutex or a channel
operation.  The following events are recorded:

 * `select` operations
 * `channel send` operations
 * `channel receive` operations
 * `semacquire` operations (`sync => Mutex.Lock()`, `sync => RWMutex.RLock()`, `sync => RWMutex.Lock()`, `sync => WaitGroup.Wait()`)
 * `notify list` operations (`cond => Wait()`)

> [!CAUTION]
> Block profiles do not include time in sleep, IO wait etc and block events are only recorded 
> when the block has cleared, as such they are not appropriate to see why a program is currently hanging.

```go
package main

import (
    "github.com/symonk/profiler"    
)

func main() {
    defer profiler.Start(profiler.WithBlockProfiler()).Stop()
    /* your code here */
}
```


-----

### :five: Thread Profiling

Thread profiling shows stack traces of code that caused new OS level threads to be created by the
go scheduler.  This is implemented in this library but it has been broken since 2013.

> [!CAUTION]
> This has been broken since 2013, do not use it!

```go
package main

import (
    "github.com/symonk/profiler"
)

func main() {
    defer profiler.Start(profiler.WithThreadProfiler()).Stop()
    /* your code here */
}
```

-----

### :six: Goroutine Profiling

...

-----

### :seven: Mutex Profiling

...

-----

### :eight: Clock Profiling

...

-----


## Available Options

* `WithAllocProfiler` => Enables allocation (memory) profiling.
* `WithBlockProfiler` => Enables block profiling.
* `WithCPUProfiler` => Enables CPU profiling (default).
* `WithCallback` => User defined callback that has the profiler in scope, invoked after teardown.
* `WithClockProfiling` => Enables CPU on & off profiling (non stdlib).
* `WithHeapProfiler` =>  Enables heap (memory) profiling.
* `WithMemoryProfilingRate` => Sets the profiling rate for memory related profiling samples.
* `WithMutexFraction` => Sets the fraction rate used in conjunction with mutex profiling.
* `WithProfileFileLocation` => Sets the custom folder location for the pprof / trace files. 
* `WithQuietOutput` => Suppresses writing to stdout/printing.
* `WithRealTimeData` => Spins a http server for the lifetime of the profiling for real curl/fetching if desired.
* `WithThreadProfiler` => Enables the os thread creation profiling.
* `WithTracing` => Enables the tracing.
* `WithoutSignalHandling` => Prevents the profiler tool signal handling, allow more fine grained user control.

