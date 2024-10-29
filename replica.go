package main

import (
	"github.com/debeando/go-common/log"
	"github.com/debeando/go-common/mysql"
)

const SQLReplica = "SHOW REPLICA STATUS"

func CollectReplica() {
	MySQLConn.Connect()
	MySQLConn.FetchAll(SQLReplica, func(row map[string]string) {
		for column := range row {
			if value, ok := mysql.ParseNumberValue(row[column]); ok {
				log.DebugWithFields("MySQL Replica", log.Fields{
					"hostname": Hostname,
					column:     value,
				})
				InfluxDBWrite("mysql_replica", column, value)
			}
		}
	})
}
