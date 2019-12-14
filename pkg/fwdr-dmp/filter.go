package fwdr_dmp

import (
	"nfCollector/pkg/utl"
	"strings"
	"unicode/utf8"
)

type Filter struct {
	Version		string
	Exporter 	string
	SrcIP		string
	SrcPort		string
	DstIP		string
	DstPort		string
	Proto 		string
}

type FilteredMetric struct {
	FlowVersion			string	`header:"Version"`
	NFSender			string	`header:"NF Exporter"`
	Bytes				string	`header:"Bytes"`
	Packets				string	`header:"Packets"`
	SrcIP				string	`header:"Src IP"`
	DstIP				string	`header:"Dst IP"`
	ProtoName			string	`header:"Proto Name"`
	SrcPort				string	`header:"Src Port"`
	SrcPortName			string	`header:"SRC Port Name"`
	DstPort				string	`header:"Dst Port"`
	DstPortName			string	`header:"Dst Port Name"`
	NextHop				string	`header:"Next Hop"`
	TCPFlags			string	`header:"TCP Flags"`
}





// Make new Filter
func NewFilter(ver string, exp string, srcIP string, srcPort string, dstIP string, dstPort string, proto string) *Filter {
	f := &Filter{
		Version:	ver,
		Exporter:	exp,
		SrcIP:		srcIP,
		SrcPort:	srcPort,
		DstIP:		dstIP,
		DstPort:	dstPort,
		Proto:		proto,
	}
	return f
}


func (f *Filter) PrepareFilteredMetrics(originalMetrics []utl.Metric) []FilteredMetric {
	isOk := false
	var filteredMetric []FilteredMetric
	for _, m := range originalMetrics {
		isOk = true

		// VERSION
		if f.Version != "*" {
			if !FilterWildCard(strings.ToLower("*" + f.Version), strings.ToLower(m.FlowVersion)) {
				isOk = false
			}
		}

		// Exporter
		if f.Exporter != "*" {
			if !FilterWildCard(strings.ToLower(f.Exporter), strings.ToLower(m.NFSender)) {
				isOk = false
			}
		}

		// SrcIP
		if f.SrcIP != "*" {
			if !FilterWildCard(strings.ToLower(f.SrcIP), strings.ToLower(m.SrcIP)) {
				isOk = false
			}
		}

		// SrcPort
		if f.SrcPort != "*" {
			if !FilterWildCard(strings.ToLower(f.SrcPort), strings.ToLower(m.SrcPort)) {
				isOk = false
			}
		}

		// DstIP
		if f.DstIP != "*" {
			if !FilterWildCard(strings.ToLower(f.DstIP), strings.ToLower(m.DstIP)) {
				isOk = false
			}
		}

		// DstPort
		if f.DstPort != "*" {
			if !FilterWildCard(strings.ToLower(f.DstPort), strings.ToLower(m.DstPort)) {
				isOk = false
			}
		}

		// Proto
		if f.Proto != "*" {
			if !FilterWildCard(strings.ToLower(f.Proto), strings.ToLower(m.ProtoName)) {
				isOk = false
			}
		}


		if isOk {
			fm := FilteredMetric{}
			fm.FlowVersion = m.FlowVersion
			fm.NFSender = m.NFSender
			fm.Bytes = m.Bytes
			fm.Packets = m.Packets
			fm.SrcIP = m.SrcIP
			fm.SrcPort = m.SrcPort
			fm.SrcPortName = m.SrcPortName
			fm.DstIP = m.DstIP
			fm.DstPort = m.DstPort
			fm.DstPortName = m.DstPortName
			fm.ProtoName = m.ProtoName
			fm.NextHop = m.NextHop
			fm.TCPFlags = m.TCPFlags
			filteredMetric = append(filteredMetric, fm)
		}
	}
	return filteredMetric
}

func FilterWildCard(pattern string, name string) bool {
Pattern:
	for len(pattern) > 0 {
		var star bool
		var chunk string
		star, chunk, pattern = scanChunk(pattern)
		if star && chunk == "" {
			// Trailing * matches rest of string.
			return true
		}
		// Look for match at current position.
		t, ok := matchChunk(chunk, name)
		// if we're the last chunk, make sure we've exhausted the name
		// otherwise we'll give a false result even if we could still match
		// using the star
		if ok && (len(t) == 0 || len(pattern) > 0) {
			name = t
			continue
		}
		if star {
			// Look for match skipping i+1 bytes.
			for i := 0; i < len(name); i++ {
				t, ok := matchChunk(chunk, name[i+1:])
				if ok {
					// if we're the last chunk, make sure we exhausted the name
					if len(pattern) == 0 && len(t) > 0 {
						continue
					}
					name = t
					continue Pattern
				}
			}
		}
		return false
	}
	return len(name) == 0
}


// scanChunk gets the next segment of pattern, which is a non-star string
// possibly preceded by a star.
func scanChunk(pattern string) (star bool, chunk, rest string) {
	for len(pattern) > 0 && pattern[0] == '*' {
		pattern = pattern[1:]
		star = true
	}
	inrange := false
	var i int
Scan:
	for i = 0; i < len(pattern); i++ {
		switch pattern[i] {
		case '*':
			if !inrange {
				break Scan
			}
		}
	}
	return star, pattern[0:i], pattern[i:]
}

func matchChunk(chunk, s string) (rest string, ok bool) {
	for len(chunk) > 0 {
		if len(s) == 0 {
			return
		}
		switch chunk[0] {
		case '?':
			_, n := utf8.DecodeRuneInString(s)
			s = s[n:]
			chunk = chunk[1:]
		default:
			if chunk[0] != s[0] {
				return
			}
			s = s[1:]
			chunk = chunk[1:]
		}
	}
	return s, true
}