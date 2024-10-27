package main

import (
	"context"
	"time"

	"github.com/debeando/go-common/log"
	"github.com/debeando/go-common/mysql"
	"github.com/influxdata/influxdb-client-go/v2"
)

const SQLVariables = "SHOW GLOBAL VARIABLES"

func CollectVariables() {
	w := InfluxDBConn.WriteAPIBlocking("debeando", InfluxDBBucket)

	MySQLConn.Connect()
	MySQLConn.FetchAll(SQLVariables, func(row map[string]string) {
		if value, ok := mysql.ParseNumberValue(row["Value"]); ok {
			log.DebugWithFields("MySQL Variables", log.Fields{
				"hostname":           MySQLHost,
				row["Variable_name"]: value,
			})

			p := influxdb2.NewPointWithMeasurement("mysql_variables").
				AddTag("_hostname", MySQLHost).
				AddField(row["Variable_name"], value).
				SetTime(time.Now())

			err := w.WritePoint(context.Background(), p)
			if err != nil {
				log.Error(err.Error())
			}
		}
	})
}
