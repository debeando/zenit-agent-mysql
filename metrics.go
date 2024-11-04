package main

type Metric struct {
	Query   string
	Name    string
	Key     string
	Value   string
	Iterate bool
}

var Metrics = []Metric{
	Metric{
		Query: "SHOW GLOBAL VARIABLES",
		Name:  "variables",
		Key:   "Variable_name",
		Value: "Value",
	},
	Metric{
		Query: "SHOW GLOBAL STATUS",
		Name:  "status",
		Key:   "Variable_name",
		Value: "Value",
	},
	Metric{
		Query:   "SHOW REPLICA STATUS",
		Name:    "replica",
		Iterate: true,
	},
	Metric{
		Query: "SELECT name, count FROM information_schema.innodb_metrics WHERE status='enabled'",
		Name:  "innodb",
		Key:   "name",
		Value: "count",
	},
}
