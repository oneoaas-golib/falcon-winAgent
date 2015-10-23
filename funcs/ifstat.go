package funcs

import (
	"fmt"
	"github.com/StackExchange/wmi"
	"github.com/ZeaLoVe/falcon-winAgent/g"
	"github.com/open-falcon/common/model"
	"sync"
)

// 单位 : B
type Win32_PerfRawData_Tcpip_NetworkInterface struct {
	Name                  string
	BytesReceivedPersec   int
	BytesSentPersec       int
	BytesTotalPersec      int
	CurrentBandwidth      int
	PacketsPersec         int
	PacketsSentPersec     int
	PacketsReceivedPersec int
}

func (w Win32_PerfRawData_Tcpip_NetworkInterface) String() string {
	return fmt.Sprintf("%s-%d-%d-%d", w.Name, w.BytesReceivedPersec, w.BytesSentPersec, w.BytesTotalPersec)
}

var (
	ifStatQuery   string
	ifStatHistory []Win32_PerfRawData_Tcpip_NetworkInterface
	ifStatCurrent []Win32_PerfRawData_Tcpip_NetworkInterface
	ifStatLock    = new(sync.RWMutex)
)

func init() {
	var dst []Win32_PerfRawData_Tcpip_NetworkInterface
	ifStatQuery = wmi.CreateQuery(&dst, "")
}

func UpdateIfStat() error {
	ifStatLock.Lock()
	defer ifStatLock.Unlock()
	ifStatHistory = ifStatCurrent
	err := wmi.Query(ifStatQuery, &ifStatCurrent)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func NetMetrics() (L []*model.MetricValue) {

	interval := g.Config().Transfer.Interval
	for i := 0; i < len(ifStatHistory); i++ {

		tags := fmt.Sprintf("iface=%s", ifStatHistory[i].Name)
		L = append(L, GaugeValue("win.net.if.in.bytes", (ifStatCurrent[i].BytesReceivedPersec-ifStatHistory[i].BytesReceivedPersec)/interval, tags))
		L = append(L, GaugeValue("win.net.if.out.bytes", (ifStatCurrent[i].BytesSentPersec-ifStatHistory[i].BytesSentPersec)/interval, tags))
		L = append(L, GaugeValue("win.net.if.total.bytes", (ifStatCurrent[i].BytesTotalPersec-ifStatHistory[i].BytesTotalPersec)/interval, tags))
		L = append(L, GaugeValue("win.net.if.in.packets", (ifStatCurrent[i].PacketsReceivedPersec-ifStatHistory[i].PacketsReceivedPersec)/interval, tags))
		L = append(L, GaugeValue("win.net.if.out.packets", (ifStatCurrent[i].PacketsSentPersec-ifStatHistory[i].PacketsSentPersec)/interval, tags))
		L = append(L, GaugeValue("win.net.if.total.packets", (ifStatCurrent[i].PacketsPersec-ifStatHistory[i].PacketsPersec)/interval, tags))
		L = append(L, GaugeValue("win.net.if.bandwidth", ifStatHistory[i].CurrentBandwidth, tags))
	}
	return
}
