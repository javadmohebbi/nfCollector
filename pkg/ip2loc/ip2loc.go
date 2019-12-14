package ip2loc

import (
	"bytes"
	"nfCollector/pkg/cnf"
	"os/exec"
	"strings"
)


type IP2Loc struct {
	CountryLong			string
	CountryShort		string
	State 				string
	City				string
	Lat					string
	Long				string
}


// Run IP2Location CMD
func Run(ip string) (IP2Loc, error){
	i2l := IP2Loc{CountryLong: "N/A", CountryShort: "N/A", State: "N/A", City: "N/A", Lat: "-1", Long: "-1"}
	conf, err := cnf.ReadConfig()
	if err != nil {
		return i2l, err
	}
	command := conf.IP2Location.Cmd
	local := conf.IP2Location.Local

	cmd := exec.Command(command, "-i", ip, "-c", "tab", "-local", local)


	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()

	if err != nil {
		return i2l, err
	}
	output :=  out.String()

	s := strings.Split(output, "\t")
	if len(s) >= 7 {
		i2l.CountryShort = strings.Replace(s[1], "\"", "", -1)
		i2l.CountryLong = strings.Replace(s[2], "\"", "", -1)
		i2l.State = strings.Replace(s[3], "\"", "", -1) //s[3]
		i2l.City = strings.Replace(s[4], "\"", "", -1) //s[4]
		i2l.Lat = s[6]
		i2l.Long = strings.Replace(s[7], "\n", "", -1)
	}
	//fmt.Println(i2l.Lat, i2l.Long)
	return i2l, nil
}
