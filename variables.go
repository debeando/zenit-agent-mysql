package main

import (
	"github.com/debeando/go-common/log"
	"github.com/debeando/go-common/mysql"
)

const SQLVariables = "SHOW GLOBAL VARIABLES"

func CollectVariables() {
	MySQLConn.Connect()
	MySQLConn.FetchAll(SQLVariables, func(row map[string]string) {
		if value, ok := mysql.ParseNumberValue(row["Value"]); ok {
			log.DebugWithFields("MySQL Variables", log.Fields{
				"hostname":           Hostname,
				row["Variable_name"]: value,
			})

			InfluxDBWrite("mysql_variables", row["Variable_name"], value)
		}
	})
}
