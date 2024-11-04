package main

import (
	"context"
	"fmt"
	"time"

	"github.com/debeando/go-common/env"
	"github.com/debeando/go-common/log"
	"github.com/influxdata/influxdb-client-go/v2"
)

type InfluxDB struct {
	Connection influxdb2.Client
	Host       string
	Port       uint16
	Token      string
	Bucket     string
}

var influxDB InfluxDB

func init() {
	influxDB.Host = env.Get("INFLUXDB_HOST", "http://127.0.0.1")
	influxDB.Port = env.GetUInt16("INFLUXDB_PORT", 8086)
	influxDB.Token = env.Get("INFLUXDB_TOKEN", "")
	influxDB.Bucket = env.Get("INFLUXDB_BUCKET", "debeando")
}

func (i *InfluxDB) ServerURL() string {
	return fmt.Sprintf("%s:%d", i.Host, i.Port)
}

func (i *InfluxDB) New() {
	i.Connection = influxdb2.NewClient(i.ServerURL(), i.Token)
}

func (i *InfluxDB) Write(hostname, measurement, key string, value interface{}) {
	go func() {
		log.DebugWithFields(fmt.Sprintf("MySQL:%s", measurement), log.Fields{
			"hostname": hostname,
			key:        value,
		})

		err := i.Connection.WriteAPIBlocking("debeando", i.Bucket).WritePoint(
			context.Background(),
			influxdb2.NewPointWithMeasurement(measurement).
				AddTag("_hostname", hostname).
				AddField(key, value).
				SetTime(time.Now()),
		)
		if err != nil {
			log.Error(err.Error())
		}
	}()
}

func (i *InfluxDB) Close() {
	if i.Connection != nil {
		i.Connection.Close()
		i.Connection = nil
	}
}
