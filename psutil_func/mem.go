package psutil_func

import (
	"fmt"
	"github.com/shirou/gopsutil/mem"
)

func PrintMem() {
	v, _ := mem.VirtualMemory()

	// almost every return value is a struct
	fmt.Printf("Total: %v, Free:%v, UsedPercent:%f%%\n", v.Total, v.Free, v.UsedPercent)
	fmt.Printf("total: %v, available:%v, used:%v, buffer:%v, cache:%v \n", v.Total, v.Available, v.Used, v.Buffers, v.Cached)

	// convert to JSON. String() is also implemented
	fmt.Println(v)
}