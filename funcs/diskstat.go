package funcs

import (
	"fmt"
	"github.com/StackExchange/wmi"
	"github.com/open-falcon/common/model"
	"log"
	"strings"
	"sync"
)

// 单位 : B
type Win32_LogicalDisk struct {
	Name      string
	Size      int
	FreeSpace int
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

	if err != nil && !strings.Contains(err.Error(), "(<nil>)") {
		log.Println(err)
		return
	}

	for _, diskstat := range diskStatHistory {
		//光驱此处收集到的为0，需要过滤掉
		//		log.Println("diskSize:", diskstat.Size)
		if diskstat.Size == 0.0 {
			continue
		}
		freePercent := 100.0 * float64(diskstat.FreeSpace) / float64(diskstat.Size)
		tags := fmt.Sprintf("mount=%s", diskstat.Name)
		L = append(L, GaugeValue("win.df.bytes.total", diskstat.Size, tags))
		L = append(L, GaugeValue("win.df.bytes.used", diskstat.Size-diskstat.FreeSpace, tags))
		L = append(L, GaugeValue("win.df.bytes.free", diskstat.FreeSpace, tags))
		L = append(L, GaugeValue("win.df.bytes.used.percent", 100.0-freePercent, tags))
		L = append(L, GaugeValue("win.df.bytes.free.percent", freePercent, tags))
	}
	return
}
