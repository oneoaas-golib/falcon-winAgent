package funcs

import (
	"github.com/StackExchange/wmi"
	"github.com/open-falcon/common/model"
	"log"
	"strings"
	"sync"
)

type Win32_Processor struct {
	LoadPercentage int
}

var (
	cpuStatQuery   string
	cpuStatHistory []Win32_Processor
	cpuStatLock    = new(sync.RWMutex)
)

func init() {
	var dst []Win32_Processor
	cpuStatQuery = wmi.CreateQuery(&dst, "")
}

func UpdateCpuStat() error {
	cpuStatLock.Lock()
	defer cpuStatLock.Unlock()
	err := wmi.Query(cpuStatQuery, &cpuStatHistory)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func CpuBusy() int {
	cpuStatLock.RLock()
	defer cpuStatLock.RUnlock()
	return cpuStatHistory[0].LoadPercentage
}

func CpuIdle() int {
	cpuStatLock.RLock()
	defer cpuStatLock.RUnlock()
	return 100 - cpuStatHistory[0].LoadPercentage
}

func CpuMetrics() []*model.MetricValue {
	err := UpdateCpuStat()
	if err != nil && !strings.Contains(err.Error(), "(<nil>)") {
		log.Println(err)
		return nil
	}

	idle := GaugeValue("win.cpu.idle", CpuIdle())
	busy := GaugeValue("win.cpu.busy", CpuBusy())
	return []*model.MetricValue{idle, busy}
}
