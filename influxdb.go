package main

import (
	"context"
	"time"

	"github.com/debeando/go-common/log"
	"github.com/influxdata/influxdb-client-go/v2"
)

func InfluxDBWrite(measurement string, key string, value interface{}) {
	go func() {
		err := InfluxDBConn.WriteAPIBlocking("debeando", InfluxDBBucket).WritePoint(
			context.Background(),
			influxdb2.NewPointWithMeasurement(measurement).
				AddTag("_hostname", Hostname).
				AddField(key, value).
				SetTime(time.Now()),
		)
		if err != nil {
			log.Error(err.Error())
		}
	}()
}
