package main

import (
	"time"

	"github.com/debeando/go-common/env"
	"github.com/debeando/go-common/log"
	"github.com/debeando/go-common/mysql"
)

var Debug bool
var Hostname string
var Interval time.Duration

func init() {
	Debug = env.GetBool("DEBUG", true)
	Interval = time.Duration(env.GetInt("INTERVAL", 3)) * time.Second
	Hostname = env.Get("HOSTNAME", MySQL.Host)

	if Debug {
		log.SetLevel(log.DebugLevel)
	}
}

func main() {
	log.Info("Start DeBeAndo Zenit Agent for MySQL")
	log.DebugWithFields("Environment Variables", log.Fields{
		"DEBUG":           Debug,
		"HOSTNAME":        Hostname,
		"INFLUXDB_BUCKET": influxDB.Bucket,
		"INFLUXDB_HOST":   influxDB.Host,
		"INFLUXDB_PORT":   influxDB.Port,
		"INFLUXDB_TOKEN":  influxDB.Token,
		"INTERVAL":        Interval,
		"MYSQL_HOST":      MySQL.Host,
		"MYSQL_PASSWORD":  MySQL.Password,
		"MYSQL_PORT":      MySQL.Port,
		"MYSQL_TIMEOUT":   MySQL.Timeout,
		"MYSQL_USER":      MySQL.Username,
	})

	influxDB.New()
	defer influxDB.Close()
	MySQL.Connection = mysql.New(MySQL.Host, MySQL.DSN())
	defer MySQL.Connection.Close()

	for {
		MySQL.Connection.Connect()

		for _, metric := range Metrics {
			MySQL.Connection.FetchAll(metric.Query, func(row map[string]string) {
				if metric.Iterate {
					for column, value := range row {
						if valueParsed, ok := mysql.ParseNumberValue(value); ok {
							influxDB.Write(Hostname, metric.Name, column, valueParsed)
						}
					}
				} else if valueParsed, ok := mysql.ParseNumberValue(row[metric.Value]); ok {
					influxDB.Write(Hostname, metric.Name, row[metric.Key], valueParsed)
				}
			})
		}

		log.Debug("Wait until next collect metrics.")
		time.Sleep(Interval)
	}
}
