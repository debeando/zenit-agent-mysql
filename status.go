package main

import (
	"context"
	"time"

	"github.com/debeando/go-common/log"
	"github.com/debeando/go-common/mysql"
	"github.com/influxdata/influxdb-client-go/v2"
)

const SQLStatus = "SHOW GLOBAL STATUS"

func CollectStatus() {
	w := InfluxDBConn.WriteAPIBlocking("debeando", InfluxDBBucket)

	MySQLConn.Connect()
	MySQLConn.FetchAll(SQLStatus, func(row map[string]string) {
		if value, ok := mysql.ParseNumberValue(row["Value"]); ok {
			log.DebugWithFields("MySQL Status", log.Fields{
				"hostname":           MySQLHost,
				row["Variable_name"]: value,
			})

			p := influxdb2.NewPointWithMeasurement("mysql_status").
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
