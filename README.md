# Netflow Collector
**NFCOL** Collects Netflow version 1, 5, 6, 7, 9 & IPFIX & stores them on **InfluxData** time-series DB (InfluxDB)


## Features
As I develop this tool for my personal usage at first step, It has a small features. Here is the list of current features.

- **Collect all available versions of Netflow**
  - **YAML** config file which help you to configure the whole application from a single file
  - Thanks to [tehmaze](https://github.com/tehmaze/netflow), We are using this package to decode Netflow traffics & we are able to decode 1, 5, 6, 7, 9 & IPFIX protocols.
  - You can forwarded decoded Netflow packets to another hosts by enabling **Packet Forwarders** in **nfc.yaml** config file
    - In addition to packet forwarder, **nfc-dump** is a client for forwarder to print a nice table for you. Also this tool can filter packets according to the provided arguments.
  - Exporting traffics to **InfluxDB** time-series database.
    - To enable this feature you need to configure **nfc.yaml** config file
    - Also, It can provide you a GEO Location information for the SOURCE & DESTINATION IP Addresses (both IPv4 & IPv6).
      - To get more information about IP2Location tool you can read this GITHUB Repo. README file [IP2Location](https://github.com/javadmohebbi/IP2Location)


## Usage
In order to use this tool you can download the compiled binaries or you can compile it for yourself.
You can download the compiled versions from the further links:
  - **Linux AMD64**: [Download Linux AMD64 Binaries](https://)
  - **Windows AMD64**: [Download Windows AMD64 Binaries](https://)
  - **MacOS AMD64**: [Download MacOS AMD64 Binaries](https://)

### Configuration
To config **nfcol** you need to provide configuration file in **yaml** format. ```*nix``` users must place this file in ```/etc/nfcol/nfc.yaml``` & ```windwos``` users must place it in ```ProgramFiles\Netflow-Collector\nfc.yaml```
***If you use installtion packages, It will create it for you automatically***

Your ```nfc.yaml``` file must look like this sample:
```
# # # # # # # # # # # # # # # # # #
#       Netflow Collector         #
#         Configuration           #
# # # # # # # # # # # # # # # # # #
**server**:
  # Listen Address
  address: 0.0.0.0

  # Listen UDP Port
  port: 6859

  # If true, nfc will write flow data into stdout
  dump: false

  # Activate forwarder
  forwarder: true

  # Host to forward - Can be separated by ; (semi-colon) eg: 127.0.0.1;192.168.100.1
  forwarderHost: 127.0.0.1

  # Forwarder UDP Port
  forwarderPort: 7161

# # # # # # # # # # # # # # # # # #
#    IP2Location Configuration    #
# # # # # # # # # # # # # # # # # #
ip2location:
  # IP2Location command path
  cmd: /usr/local/bin/ip2location

  # Path to Local GEO Database. Read more at:
  #      https://github.com/javadmohebbi/IP2Location#local-database-format
  local: /etc/ip2location/local.csv

# # # # # # # # # # # # # # # # # #
#     Exporter Configuration      #
# # # # # # # # # # # # # # # # # #
exporter:
  # Enable if it's true
  enable: false

  # Currently Only InfluxDB (1.x) supported
  type: influxdb


# # # # # # # # # # # # # # # # # #
#     InfluxDB Configuration      #
# # # # # # # # # # # # # # # # # #
influxDB:
  # InfluxDB Host
  host: 127.0.0.1

  # InfluxDB Port
  port: 8086

  # InfluxDB Username. Can be null
  username: #user

  # InfluxDB Password. Can be null
  password: #secret

  # InfluxDB Database
  # InfluxDB command example for creating database:
  #       CREATE DATABASE "netflowDB" WITH DURATION 10d REPLICATION 1 SHARD DURATION 1h NAME "nfc"
  database: netflowDB

  # Temp Dir for InfluxDB Metrics. MUST be ended with / (Linux) or \ (Windows)
  tmpDir: /tmp/nfcol/


# # # # # # # # # # # # # # # # # #
#          Measurements           #
# # # # # # # # # # # # # # # # # #
measurements:
  # Netflow Summary Measurement Name
  summaryProto: sum_proto

  # Netflow GEO Summary Measurement Name
  summaryProtoGeo: sum_proto_geo
```
