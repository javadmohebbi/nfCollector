package nfv6

import (
	"github.com/landoop/tableprinter"
	"github.com/tehmaze/netflow/netflow6"
	"net"
	"nfCollector/pkg/utl"
	"nfCollector/pkg/utl/proto"
	"nfCollector/pkg/utl/service"
	"os"

	"fmt"
)

func Prepare(addr string, p *netflow6.Packet) []utl.Metric{
	nfExporter, _, _ := net.SplitHostPort(addr)
	var metrics []utl.Metric
	var met utl.Metric
	for _, r := range p.Records {
		met = utl.Metric{OutBytes: "0", InBytes: "0", OutPacket: "0", InPacket: "0", NFSender: nfExporter}
		met.FlowVersion = "Netflow-V6"
		met.Direction = "Unsupported"
		met.First = fmt.Sprintf("%v", r.First)
		met.Last = fmt.Sprintf("%v", r.Last)
		met.Protocol = fmt.Sprintf("%v", r.Protocol)
		met.ProtoName = proto.ProtoToName(met.Protocol)
		met.Bytes = fmt.Sprintf("%v", r.Bytes)
		met.Packets = fmt.Sprintf("%v", r.Packets)
		met.TCPFlags = fmt.Sprintf("%v", r.TCPFlags)
		met.SrcAs = fmt.Sprintf("%v", r.SrcAS)
		met.DstAs = fmt.Sprintf("%v", r.DstAS)
		met.SrcMask = fmt.Sprintf("%v", r.SrcMask)
		met.DstMask = fmt.Sprintf("%v", r.DstMask)
		met.NextHop = fmt.Sprintf("%v", r.NextHop)

		met.SrcIP = fmt.Sprintf("%v", r.SrcAddr)
		met.DstIP = fmt.Sprintf("%v", r.DstAddr)

		met.SrcPort = fmt.Sprintf("%v", r.SrcPort)
		met.SrcPortName = service.GetPortName(met.SrcPort, met.ProtoName)

		met.DstPort = fmt.Sprintf("%v", r.DstPort)
		met.DstPortName = service.GetPortName(met.DstPort, met.ProtoName)

		metrics = append(metrics, met)

	}
	return metrics
}



// Dump - Print it as table
func Dump(metrics []utl.Metric) {
	printer := tableprinter.New(os.Stdout)
	printer.BorderTop, printer.BorderBottom, printer.BorderLeft, printer.BorderRight = true, true, true, true
	printer.CenterSeparator = "│"
	printer.ColumnSeparator = "│"
	printer.RowSeparator = "─"
	printer.Print(metrics)
}
