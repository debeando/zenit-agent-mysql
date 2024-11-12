package main

import (
	"fmt"
	"strings"
)

type Query struct {
	Key       string
	Name      string
	Statement string
	UnPivot   bool
	Value     string
}

var Queries = []Query{
	Query{
		Name:      "mysql_variables",
		Statement: "SHOW GLOBAL VARIABLES",
		Key:       "Variable_name",
		Value:     "Value",
	},
	Query{
		Name:      "mysql_status",
		Statement: "SHOW GLOBAL STATUS",
		Key:       "Variable_name",
		Value:     "Value",
	},
	Query{
		Name:      "mysql_replica",
		Statement: "SHOW REPLICA STATUS",
		UnPivot:   true,
	},
	Query{
		Statement: "SELECT name, count FROM information_schema.innodb_metrics WHERE status='enabled'",
		Name:      "mysql_innodb",
		Key:       "name",
		Value:     "count",
	},
	Query{
		Name: "mysql_performance_schema",
		Statement: fmt.Sprintf(`
		SELECT
            ifnull(SCHEMA_NAME, 'NONE') AS SCHEMA_NAME,
            DIGEST,
            DIGEST_TEXT,
            COUNT_STAR,
            SUM_TIMER_WAIT,
            SUM_ERRORS,
            SUM_WARNINGS,
            SUM_ROWS_AFFECTED,
            SUM_ROWS_SENT,
            SUM_ROWS_EXAMINED,
            SUM_CREATED_TMP_DISK_TABLES,
            SUM_CREATED_TMP_TABLES,
            SUM_SORT_MERGE_PASSES,
            SUM_SORT_ROWS,
            SUM_NO_INDEX_USED
        FROM performance_schema.events_statements_summary_by_digest
        WHERE SCHEMA_NAME NOT IN ('mysql', 'performance_schema', 'information_schema')
          AND last_seen > DATE_SUB(NOW(), INTERVAL %d SECOND)
        ORDER BY SUM_TIMER_WAIT DESC;
        `, int(getInterval().Seconds())),
		UnPivot: true,
	},
	Query{
		Name: "mysql_query_latency",
		Statement: fmt.Sprintf(`
		SELECT
            ifnull(SCHEMA_NAME, 'NONE') AS SCHEMA_NAME,
            sum(count_star) AS count_star,
            round(avg_timer_wait/1000000, 0) AS avg_time_us
        FROM performance_schema.events_statements_summary_by_digest
        WHERE SCHEMA_NAME NOT IN ('mysql', 'performance_schema', 'information_schema')
          AND last_seen > DATE_SUB(NOW(), INTERVAL %d SECOND)
        GROUP BY SCHEMA_NAME
        ORDER BY avg_time_us DESC;
        `, int(getInterval().Seconds())),
		UnPivot: true,
	},
}

func (q *Query) Beautifier() string {
	q.Statement = strings.ReplaceAll(q.Statement, "\r\n", " ")
	q.Statement = strings.ReplaceAll(q.Statement, "\n", " ")
	q.Statement = strings.ReplaceAll(q.Statement, "\t", " ")
	q.Statement = strings.ReplaceAll(q.Statement, "  ", "")
	q.Statement = strings.Trim(q.Statement, " ")

	return q.Statement
}
