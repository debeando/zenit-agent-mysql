package main

import (
	"github.com/debeando/go-common/log"
	"github.com/debeando/go-common/mysql"
)

const SQLInnoDB = "SELECT name, count FROM INFORMATION_SCHEMA.INNODB_METRICS WHERE status='enabled'"

func CollectInnoDB() {
	MySQLConn.Connect()
	MySQLConn.FetchAll(SQLInnoDB, func(row map[string]string) {
		if value, ok := mysql.ParseNumberValue(row["count"]); ok {
			log.DebugWithFields("MySQL InnoDB", log.Fields{
				"hostname":  Hostname,
				row["name"]: value,
			})

			InfluxDBWrite("mysql_innodb", row["name"], value)
		}
	})
}
