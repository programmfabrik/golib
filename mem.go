package lib

import (
	"fmt"
	"runtime"
	"strings"
)

// PrintMem dumps some mem info to the console
func PrintMem(info ...string) {
	ms := runtime.MemStats{}
	runtime.ReadMemStats(&ms)
	println(fmt.Sprintf(strings.Join(info, " ")+" MEM: Heap: %s Sys: %s", HumanByteSize(ms.HeapAlloc), HumanByteSize(ms.Sys)))
}

// 257262: 40.9 MB
// 1564397
