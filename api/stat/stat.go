package stat

import (
	"runtime"
	
	"admincheckapi/api/resource"
)

//
// bToMb connverts bytes to Mbytes
//
func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

//
// Info provides mem stat info
//
func Info() resource.Stat {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	s := resource.Stat{
		m.Alloc,
		m.TotalAlloc,
		m.Sys,
		m.NumGC,
	}
	
	return s
}
