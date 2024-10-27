package main

import (
	"context"
	"time"

	"github.com/debeando/go-common/log"
	"github.com/debeando/go-common/mysql"
	"github.com/influxdata/influxdb-client-go/v2"
)

const SQLReplica = "SHOW REPLICA STATUS"

func CollectReplica() {
	w := InfluxDBConn.WriteAPIBlocking("debeando", InfluxDBBucket)

	MySQLConn.Connect()
	MySQLConn.FetchAll(SQLReplica, func(row map[string]string) {
		for column := range row {
			if value, ok := mysql.ParseNumberValue(row[column]); ok {
				log.DebugWithFields("MySQL Replica", log.Fields{
					"hostname": MySQLHost,
					column:     value,
				})

				p := influxdb2.NewPointWithMeasurement("mysql_variables").
					AddTag("_hostname", MySQLHost).
					AddField(column, value).
					SetTime(time.Now())

				err := w.WritePoint(context.Background(), p)
				if err != nil {
					log.Error(err.Error())
				}

			}
		}
	})
}
