package main

import (
	"fmt"
	"time"

	"github.com/debeando/go-common/env"
	"github.com/debeando/go-common/log"
	"github.com/debeando/go-common/mysql"
	"github.com/influxdata/influxdb-client-go/v2"
)

var Debug string
var Interval time.Duration
var InfluxDBConn influxdb2.Client
var InfluxDBHost string
var InfluxDBPort uint16
var InfluxDBToken string
var InfluxDBBucket string
var MySQLConn *mysql.Connection
var MySQLDSN string
var MySQLHost string
var MySQLPassword string
var MySQLPort uint16
var MySQLTimeout uint8
var MySQLUser string

func init() {
	Debug = env.Get("DEBUG", "true")
	Interval = time.Duration(env.GetInt("INTERVAL", 3)) * time.Second

	InfluxDBHost = env.Get("INFLUXDB_HOST", "http://127.0.0.1")
	InfluxDBPort = env.GetUInt16("INFLUXDB_PORT", 8086)
	InfluxDBToken = env.Get("INFLUXDB_TOKEN", "")
	InfluxDBBucket = env.Get("INFLUXDB_BUCKET", "debeando")

	MySQLHost = env.Get("MYSQL_HOST", "127.0.0.1")
	MySQLPassword = env.Get("MYSQL_PASSWORD", "monitoring")
	MySQLPort = env.GetUInt16("MYSQL_PORT", 3306)
	MySQLTimeout = env.GetUInt8("MYSQL_TIMEOUT", 10)
	MySQLUser = env.Get("MYSQL_USER", "monitoring")
	MySQLDSN = (&mysql.MySQL{
		Host:     MySQLHost,
		Password: MySQLPassword,
		Port:     MySQLPort,
		Timeout:  MySQLTimeout,
		Username: MySQLUser,
	}).DSN()

	if Debug == "true" {
		log.SetLevel(log.DebugLevel)
	}
}

func main() {
	log.Info("Start DeBeAndo Zenit Agent for MySQL")
	log.DebugWithFields("Environment Variables", log.Fields{
		"DEBUG":           Debug,
		"INFLUXDB_BUCKET": InfluxDBBucket,
		"INFLUXDB_HOST":   InfluxDBHost,
		"INFLUXDB_PORT":   InfluxDBPort,
		"INFLUXDB_TOKEN":  InfluxDBToken,
		"INTERVAL":        Interval,
		"MYSQL_HOST":      MySQLHost,
		"MYSQL_PASSWORD":  MySQLPassword,
		"MYSQL_PORT":      MySQLPort,
		"MYSQL_TIMEOUT":   MySQLTimeout,
		"MYSQL_USER":      MySQLUser,
	})

	InfluxDBConn = influxdb2.NewClient(fmt.Sprintf("%s:%d", InfluxDBHost, InfluxDBPort), InfluxDBToken)
	defer InfluxDBConn.Close()
	MySQLConn = mysql.New(MySQLHost, MySQLDSN)
	defer MySQLConn.Close()

	for {
		CollectVariables()
		CollectStatus()
		CollectInnoDB()
		CollectReplica()
		log.Debug("Wait until next collect metrics.")
		time.Sleep(Interval)
	}
}
