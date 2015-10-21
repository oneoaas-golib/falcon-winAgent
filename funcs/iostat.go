package funcs

import (
	"github.com/StackExchange/wmi"
	"github.com/open-falcon/common/model"
	"log"
	"sync"
)

// 单位 : B
type Win32_PerfFormattedData_PerfDisk_LogicalDisk struct {
	Name                 string
	DiskReadBytesPerSec  int
	DiskWriteBytesPerSec int
}

var (
	ioStatQuery   string
	ioStatHistory []Win32_PerfFormattedData_PerfDisk_LogicalDisk
	ioStatLock    = new(sync.RWMutex)
)

func init() {
	var dst []Win32_PerfFormattedData_PerfDisk_LogicalDisk
	ioStatQuery = wmi.CreateQuery(&dst, "")
}

func UpdateIoStat() error {
	ioStatLock.Lock()
	defer ioStatLock.Unlock()
	err := wmi.Query(ioStatQuery, &ioStatHistory)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func DiskIOMetrics() (L []*model.MetricValue) {

	err := UpdateIoStat()
	if err != nil {
		log.Println(err)
		return
	}

	for _, iostat := range ioStatHistory {
		device := "device=" + iostat.Name

		L = append(L, CounterValue("win.disk.read.bytes.persec", iostat.DiskReadBytesPerSec, device))
		L = append(L, CounterValue("win.disk.write.bytes.persec", iostat.DiskWriteBytesPerSec, device))
	}
	return
}
