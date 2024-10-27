package main

import (
	"github.com/debeando/go-common/log"
	"github.com/debeando/go-common/mysql"
)

const SQLVariables = "SHOW GLOBAL VARIABLES"

func CollectVariables() {
	m := mysql.New(MySQLHost, MySQLDSN)
	m.Connect()
	m.FetchAll(SQLVariables, func(row map[string]string) {
		if value, ok := mysql.ParseNumberValue(row["Value"]); ok {
			log.DebugWithFields("MySQL Variables", log.Fields{
				"hostname":           MySQLHost,
				row["Variable_name"]: value,
			})
		}
	})
	m.Close()
}
