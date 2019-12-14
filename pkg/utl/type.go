package utl

import (
	"github.com/fln/nf9packet"
)



type TemplateCache map[string]*nf9packet.TemplateRecord

type Metric struct {
	FlowVersion			string	`header:"Version"`
	NFSender			string	`header:"NF Exporter"`
	Last				string	//`header:"Last"`
	First				string	//`header:"First"`
	Bytes				string	`header:"Bytes"`
	Packets				string	`header:"Packets"`
	InBytes				string	//`header:"Bytes In"`
	InPacket			string	//`header:"Packets In"`
	OutBytes			string	//`header:"Bytes Out"`
	OutPacket			string	//`header:"Packets Out"`
	InEthernet			string	//`header:"In Eth"`
	OutEthernet			string	//`header:"Out Eth"`
	SrcIP				string	`header:"Src IP"`
	SrcIp2lCountryShort	string	//`header:"sCountry_S"`
	SrcIp2lCountryLong	string	//`header:"sCountry_L"`
	SrcIp2lState		string	//`header:"sState"`
	SrcIp2lCity			string	//`header:"sCity"`
	SrcIp2lLat			string	//`header:"sLat"`
	SrcIp2lLong			string	//`header:"sLong"`
	DstIP				string	`header:"Dst IP"`
	DstIp2lCountryShort	string	//`header:"dCountry_S"`
	DstIp2lCountryLong	string	//`header:"dCountry_L"`
	DstIp2lState		string	//`header:"dState"`
	DstIp2lCity			string	//`header:"dCity"`
	DstIp2lLat			string	//`header:"dLat"`
	DstIp2lLong			string	//`header:"dLong"`
	Protocol			string	//`header:"Proto"`
	ProtoName			string	`header:"Proto Name"`
	SrcToS				string	//`header:"SrcToS"`
	SrcPort				string	`header:"Src Port"`
	SrcPortName			string	`header:"SRC Port Name"`
	DstPort				string	`header:"Dst Port"`
	DstPortName			string	`header:"Dst Port Name"`
	FlowSamplerId		string	//`header:"FlowSampleId"`
	VendorPROPRIETARY	string	//`header:"VendorPROPRIETARY"`
	NextHop				string	`header:"Next Hop"`
	DstMask				string	//`header:"DstMask"`
	SrcMask				string	//`header:"SrcMask"`
	TCPFlags			string	`header:"TCP Flags"`
	Direction			string	//`header:"Direction"`
	DstAs				string	//`header:"DstAs"`
	SrcAs				string	//`header:"SrcAs"`

}
