package main

import (
	"fmt"
	"time"

	"github.com/debeando/go-common/env"
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
	i.Connection = influxdb2.NewClientWithOptions(
		i.ServerURL(),
		i.Token,
		influxdb2.DefaultOptions().SetBatchSize(100),
	)
}

func (i *InfluxDB) Write(metrics Metrics) {
	writeAPI := i.Connection.WriteAPI("debeando", i.Bucket)

	for _, metric := range metrics {
		point := influxdb2.NewPoint(
			metric.Measurement,
			metric.TagsToMap(),
			metric.FieldsToMap(),
			time.Now(),
		)

		writeAPI.WritePoint(point)
	}

	writeAPI.Flush()
}

func (i *InfluxDB) Close() {
	if i.Connection != nil {
		i.Connection.Close()
		i.Connection = nil
	}
}
