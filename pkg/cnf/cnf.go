package cnf

import (
	"log"
	"os"
	"runtime"

	"github.com/spf13/viper"
)

// Configurations struct
type Configurations struct {
	Server       ServerConfigurations
	IP2Location  IP2LocationConfigurations
	InfluxDB     InfluxDBConfigurations
	Exporter     ExporterConfigurations
	Measurements MeasurementsConfigurations
}

// ServerConfigurations struct
type ServerConfigurations struct {
	Address       string
	Port          int
	Forwarder     bool
	ForwarderHost string
	ForwarderPort int
	Dump          bool
}

// IP2LocationConfigurations struct
type IP2LocationConfigurations struct {
	Cmd   string
	Local string
}

// ExporterConfigurations struct
type ExporterConfigurations struct {
	Enable bool
	Type   string
}

// InfluxDBConfigurations struct
type InfluxDBConfigurations struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
	TmpDir   string
}

// MeasurementsConfigurations struct
type MeasurementsConfigurations struct {
	SummaryProto    string
	SummaryProtoGeo string
	SummarySrc      string
	SummaryDst      string
}

var (
	configPath      string = "/home/mj/go/src/nfCollector/configs/"
	configName      string = "nfc"
	configExtension string = "yaml"
)

const IsDev bool = false

// ReadConfig - ReadConfig & Return Configurations struct
func ReadConfig() (Configurations, error) {

	configPath = PrepareConfigPath()

	// Set the file name of the configurations file
	viper.SetConfigName(configName)
	// Set the path to look for the configurations file
	viper.AddConfigPath(configPath)
	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()
	viper.SetConfigType(configExtension)
	var conf Configurations

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading config file, %s", err)
		return conf, err
	}
	err := viper.Unmarshal(&conf)
	if err != nil {
		log.Printf("Unable to decode into struct, %v", err)
		return conf, err
	}
	return conf, err
}

// PrepareConfigPath - Config Dir
func PrepareConfigPath() string {
	if IsDev {
		return configPath
	}
	if runtime.GOOS == "windows" {
		return os.Getenv("ProgramFiles") + "\\Netflow-Collector\\"
	} else if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		return "/etc/nfcol/"
	}
	log.Fatal("Read configs from: ", configPath)
	return configPath
}
