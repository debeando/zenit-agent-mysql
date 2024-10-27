package main

import (
	"time"

	"github.com/debeando/go-common/env"
	"github.com/debeando/go-common/log"
	"github.com/debeando/go-common/mysql"
)

var Debug string
var Interval time.Duration
var MySQLDSN string
var MySQLHost string
var MySQLPassword string
var MySQLPort uint16
var MySQLTimeout uint8
var MySQLUser string

func init() {
	Debug = env.Get("DEBUG", "true")
	Interval = time.Duration(env.GetInt("INTERVAL", 3)) * time.Second
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
		"DEBUG":          Debug,
		"INTERVAL":       Interval,
		"MYSQL_HOST":     MySQLHost,
		"MYSQL_PORT":     MySQLPort,
		"MYSQL_USER":     MySQLUser,
		"MYSQL_PASSWORD": MySQLPassword,
		"MYSQL_TIMEOUT":  MySQLTimeout,
	})

	for {
		CollectStatus()
		log.Debug("Wait until next collect metrics.")
		time.Sleep(Interval)
	}
}
