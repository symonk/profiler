## Profiler :star:

`profiler` is a utility library to easily enable various level of profiling for go programs.
The various modes available are outlined below.

> [!NOTE] Some functional options change the behaviour of various profiling setups.

> [!NOTE] By default `cpu profiling` is enabled if no profile is provided.

> [!NOTE] By default the profile files are written to the directory executing your program.

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
}
```

### :two: Heap Profiling 

Heap profiling reports memory allocation samples, useful for monitoring current and historical memory
usage and to check for potential memory leaks in your program.  Heap profiling tracks the allocation
sites for all live objects in the application memory and for all objects allocated since the program
has started.  By default heap profiling will display live objects, scaled by size.

### :three: Alloc Profiling

Alloc profiling is essentially the same as heap profiling except rather than the default of live objects
scaled by size, it reports the `-alloc_space` data, that is the total number of bytes allocated since the
program has began (including garbage collected bytes).

### :four: Block Profiling

...

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


