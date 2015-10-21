package funcs

import (
	"fmt"
	"github.com/StackExchange/wmi"
	"github.com/open-falcon/common/model"
	"log"
	"sync"
)

// 单位 : B
type Win32_LogicalDisk struct {
	Name         string
	FileSystem   string
	Size         int
	FreeSpace    int
	Availability int
}

var (
	diskStatQuery   string
	diskStatHistory []Win32_LogicalDisk
	diskStatLock    = new(sync.RWMutex)
)

func init() {
	var dst []Win32_LogicalDisk
	diskStatQuery = wmi.CreateQuery(&dst, "")
}

func UpdateDiskStat() error {
	diskStatLock.Lock()
	defer diskStatLock.Unlock()
	err := wmi.Query(diskStatQuery, &diskStatHistory)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func DeviceMetrics() (L []*model.MetricValue) {
	err := UpdateDiskStat()

	if err != nil {
		log.Println(err)
		return
	}

	for _, diskstat := range diskStatHistory {

		freePercent := (1.0 * diskstat.FreeSpace) / diskstat.Size
		tags := fmt.Sprintf("mount=%s,fstype=%s", diskstat.Name, diskstat.FileSystem)
		L = append(L, GaugeValue("win.df.bytes.total", diskstat.Size, tags))
		L = append(L, GaugeValue("win.df.bytes.used", diskstat.Size-diskstat.FreeSpace, tags))
		L = append(L, GaugeValue("win.df.bytes.free", diskstat.FreeSpace, tags))
		L = append(L, GaugeValue("win.df.bytes.used.percent", 100-freePercent, tags))
		L = append(L, GaugeValue("win.df.bytes.free.percent", freePercent, tags))
	}
	return
}
