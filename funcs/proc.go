package funcs

import (
	"github.com/StackExchange/wmi"
	"github.com/ZeaLoVe/falcon-winAgent/g"
	"github.com/open-falcon/common/model"
	"log"
	"strings"
	"sync"
)

type Win32_Process struct {
	Name        string
	CommandLine string
}

var (
	procStatQuery   string
	procStatHistory []Win32_Process
	procStatLock    = new(sync.RWMutex)
)

func init() {
	var dst []Win32_Process
	procStatQuery = wmi.CreateQuery(&dst, "")
}

func UpdateProcStat() error {
	procStatLock.Lock()
	defer procStatLock.Unlock()
	err := wmi.Query(procStatQuery, &procStatHistory)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func ProcMetrics() (L []*model.MetricValue) {

	reportProcs := g.ReportProcs()
	sz := len(reportProcs)
	if sz == 0 {
		return
	}

	err := UpdateProcStat()
	//wmi获取数据时候转化nil到string会抛出错误，这个错误忽略掉
	if err != nil && !strings.Contains(err.Error(), "unsupported type") {
		log.Println(err)
		return
	}

	pslen := len(procStatHistory)

	for tags, m := range reportProcs {
		cnt := 0
		for i := 0; i < pslen; i++ {
			if is_a(&procStatHistory[i], m) {
				cnt++
			}
		}

		L = append(L, GaugeValue("proc.num", cnt, tags))
	}

	return
}

func is_a(p *Win32_Process, m map[int]string) bool {
	// only one kv pair
	for key, val := range m {
		if key == 1 {
			// name
			if val != p.Name {
				return false
			}
		} else if key == 2 {
			// cmdline
			if !strings.Contains(p.CommandLine, val) {
				return false
			}
		}
	}
	return true
}
