# covid19
Small golang code to get corona virus information for Brazil and publish to InfluxDB database

What do you need:
- influxdb installation (without authentication for posting data points over http)
- golang 1.11 (with export GO111MODULES=on or go 1.12+)
- grafana installation
- influxdb configured as datasource for grafana

Attention: I don't give any guarantees that it works outside my machine ;-)

1) Create covid19 database in your influxdb installation
2) Build the app:
`go mod tidy`
and
`go build`
3) Run update.sh file gets the database from remote url and saves resources/db.json file.
4) Run the program
`INFLUXDB_URL=http://influxdb-url:influxdb-port ./covid19 update`
  
For dashboard works:
- import grafana/dashboard.json to your grafana installation
