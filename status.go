package main

import (
	"github.com/debeando/go-common/log"
	"github.com/debeando/go-common/mysql"
)

const SQLStatus = "SHOW GLOBAL STATUS"

func CollectStatus() {
	MySQLConn.Connect()
	MySQLConn.FetchAll(SQLStatus, func(row map[string]string) {
		if value, ok := mysql.ParseNumberValue(row["Value"]); ok {
			log.DebugWithFields("MySQL Status", log.Fields{
				"hostname":           Hostname,
				row["Variable_name"]: value,
			})

			InfluxDBWrite("mysql_status", row["Variable_name"], value)
		}
	})
}
