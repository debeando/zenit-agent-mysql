package main

import (
    "time"

    "github.com/debeando/go-common/cast"
    "github.com/debeando/go-common/env"
    "github.com/debeando/go-common/log"
)

var debug string
var interval time.Duration

func init() {
    debug = env.Get("DEBUG", "true")
    interval = time.Duration(cast.StringToInt(env.Get("INTERVAL", "3")))

    if debug == "true" {
        log.SetLevel(log.DebugLevel)
    }
}

func main() {    
    log.Info("Start DeBeAndo Zenit Agent for MySQL")
    log.DebugWithFields("Environment Variables", log.Fields{
            "DEBUG": debug,
            "INTERVAL": interval,
            "MYSQL_HOST": env.Get("MYSQL_HOST", "127.0.0.1"),
            "MYSQL_PORT": env.Get("MYSQL_PORT", "3306"),
            "MYSQL_USER": env.Get("MYSQL_USER", "monitoring"),
            "MYSQL_PASSWORD": env.Get("MYSQL_PASSWORD", "monitoring"),
    })

    for {
        log.Debug("Wait until next collect metrics.")
        time.Sleep(interval * time.Second)
    }
}
