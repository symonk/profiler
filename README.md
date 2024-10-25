### Profiler 

`profiler` is an abstraction layer around the golang profiling stdlib, namely `runtime/pprof`
and `runtime/trace` to simplify and make profiling configuration easy and fast.  It offers an
assortment of options and configurations as well as various non stdlib options such as the `fgprof`
wall clock option, which pieces together a single `.pprof` file which includes both cpu `ON` and `OFF`
data.  A basic example for each of the profilers is outlined below.

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
operation.  The following things cases are measured:

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
    defer profiler.Start(profiler.WithBlockProfiler(), profiler.WithBlockProfileRate(100_000_000)).Stop()
}
```


-----

### :five: Thread Profiling

...

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


