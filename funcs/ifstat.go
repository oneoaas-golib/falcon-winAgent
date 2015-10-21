package funcs

import (
	"fmt"
	"github.com/StackExchange/wmi"
	"github.com/open-falcon/common/model"
	"log"
	"sync"
)

// 单位 : B
type Win32_PerfRawData_Tcpip_NetworkInterface struct {
	Name                   string
	BytesReceivedPersec    int
	BytesSentPersec        int
	BytesTotalPersec       int
	CurrentBandwidth       int
	PacketsPersec          int
	PacketsSentPersec      int
	PacketsReceivedPersec  int
	PacketsReceivedUnknown int
}

var (
	ifStatQuery   string
	ifStatHistory []Win32_PerfRawData_Tcpip_NetworkInterface
	ifStatLock    = new(sync.RWMutex)
)

func init() {
	var dst []Win32_PerfRawData_Tcpip_NetworkInterface
	ifStatQuery = wmi.CreateQuery(&dst, "")
}

func UpdateIfStat() error {
	ifStatLock.Lock()
	defer ifStatLock.Unlock()
	err := wmi.Query(ifStatQuery, &ifStatHistory)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func NetMetrics() (L []*model.MetricValue) {
	err := UpdateIfStat()

	if err != nil {
		log.Println(err)
		return
	}

	for _, ifstat := range ifStatHistory {

		tags := fmt.Sprintf("iface=%s", ifstat.Name)
		L = append(L, GaugeValue("win.net.bytes.recieve.persec", ifstat.BytesReceivedPersec, tags))
		L = append(L, GaugeValue("win.net.bytes.send.persec", ifstat.BytesSentPersec, tags))
		L = append(L, GaugeValue("win.net.bytes.total.persec", ifstat.BytesTotalPersec, tags))
		L = append(L, GaugeValue("win.net.bytes.speed.persect", ifstat.CurrentBandwidth, tags))
	}
	return
}
