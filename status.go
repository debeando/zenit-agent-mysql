package main

import (
	"github.com/debeando/go-common/log"
	"github.com/debeando/go-common/mysql"
)

const SQLStatus = "SHOW GLOBAL STATUS"

func CollectStatus() {
	m := mysql.New(MySQLHost, MySQLDSN)
	m.Connect()
	m.FetchAll(SQLStatus, func(row map[string]string) {
		if value, ok := mysql.ParseNumberValue(row["Value"]); ok {
			log.DebugWithFields("MySQL Status", log.Fields{
				"hostname":           MySQLHost,
				row["Variable_name"]: value,
			})
		}
	})
	m.Close()
}
