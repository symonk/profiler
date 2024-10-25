## Profiler

`profiler` is a utility library to easily enable various level of profiling for go programs.
The various modes available are:

-----

##### CPU Profiling

Profiling CPU can help identify areas of your code where most CPU cycle execution time is spent.  For
programs with a lot of IO wait etc, the graph may not be extremely useful, trace may be of
better benefit there.

##### Heap Profiling

Heap profiling reports memory allocation samples, useful for monitoring current and historical memory
usage and to check for potential memory leaks in your program.  Heap profiling tracks the allocation
sites for all live objects in the application memory and for all objects allocated since the program
has started.  By default heap profiling will display live objects, scaled by size.

##### Alloc Profiling

Alloc profiling is essentially the same as heap profiling except rather than the default of live objects
scaled by size, it reports the `-alloc_space` data, that is the total number of bytes allocated since the
program has began (including garbage collected bytes).

##### Block Profiling


## 

```go
package main

import (
    "github.com/symonk/profiler
)

/*
By default, calling profiler.Start() will enable CPU profiling, but for verbosity
we will include it in this example.  Similarly all profile files are written to
the executing directory by default unless options for file location are provided.

All of the examples are outlined in a single main() call below, however you should
typically avoid trying to profile more than one thing at once, use them as a 
reference.
*/

func main() {
    // CPU Profiling writing to a custom location for the cpu.pprof file
    defer profiler.Start(profiler.WithCPUProfiler(), profiler.WithProfileFileLocation("/tmp/profiles")).Stop()

    // Heap Profiling writing to the current directory
    defer profiler.Start(profiler.WithHeapMemoryProfiling()).Stop()

    // Allocation profiling writing to the current directory
    defer profiler.Start(profiler.WithAllocMemoryProfiling()).Stop()

    // Block profiling writing to the current directory
    
}
```

----

## Available Options


