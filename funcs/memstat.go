package funcs

import (
	"github.com/StackExchange/wmi"
	"github.com/open-falcon/common/model"
	"log"
	"sync"
)

// 单位 : K
type Win32_OperatingSystem struct {
	FreePhysicalMemory     int
	TotalVisibleMemorySize int
	FreeVirtualMemory      int
}

var (
	memStatQuery   string
	memStatHistory []Win32_OperatingSystem
	memStatLock    = new(sync.RWMutex)
)

func init() {
	var dst []Win32_OperatingSystem
	memStatQuery = wmi.CreateQuery(&dst, "")
}

func UpdateMemStat() error {
	memStatLock.Lock()
	defer memStatLock.Unlock()
	err := wmi.Query(memStatQuery, &memStatHistory)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func MemMetrics() []*model.MetricValue {
	err := UpdateMemStat()
	if err != nil {
		log.Println(err)
		return nil
	}

	return []*model.MetricValue{
		GaugeValue("win.mem.memtotal", memStatHistory[0].TotalVisibleMemorySize),
		GaugeValue("win.mem.memused", memStatHistory[0].TotalVisibleMemorySize-memStatHistory[0].FreePhysicalMemory),
		GaugeValue("win.mem.memfree", memStatHistory[0].FreePhysicalMemory),
		GaugeValue("win.mem.swapfree", memStatHistory[0].FreeVirtualMemory),
	}

}
