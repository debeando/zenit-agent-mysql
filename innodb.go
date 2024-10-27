package main

import (
	"context"
	"regexp"
	"strings"
	"time"

	"github.com/debeando/go-common/log"
	"github.com/debeando/go-common/mysql"
	"github.com/influxdata/influxdb-client-go/v2"
)

const SQLInnoDB = "SHOW ENGINE INNODB STATUS"

var Sections = map[string][]string{
	"BACKGROUND THREAD": {
		`(srv_master_thread loops: (?P<thread_loops_srv_active>\d+) srv_active, (?P<thread_loops_srv_shutdown>\d+) srv_shutdown, (?P<thread_loops_srv_idle>\d+) srv_idle)\n`,
		`(srv_master_thread log flush and writes: (?P<thread_log_flush_and_writes>\d+))\n`,
	},
	"SEMAPHORES": {
		`(OS WAIT ARRAY INFO: reservation count (?P<semaphores_reservation_count>\d+))\n`,
		`(OS WAIT ARRAY INFO: signal count (?P<semaphores_signal_count>\d+))\n`,
		`(RW-shared spins (?P<semaphores_shared_spins>\d+), rounds (?P<semaphores_shared_rounds>\d+), OS waits (?P<semaphores_shared_waits>\d+))\n`,
		`(RW-excl spins (?P<semaphores_excl_spins>\d+), rounds (?P<semaphores_excl_rounds>\d+), OS waits (?P<semaphores_excl_waits>\d+))\n`,
		`(RW-sx spins (?P<semaphores_sx_spins>\d+), rounds (?P<semaphores_sx_rounds>\d+), OS waits (?P<semaphores_sx_waits>\d+))\n`,
		`(Spin rounds per wait: (?P<semaphores_spin_rounds_per_wait>([0-9]*[.])?[0-9]+) RW-shared, (?P<semaphores_spin_rounds_per_wait_rw_shared>([0-9]*[.])?[0-9]+) RW-excl, (?P<semaphores_spin_rounds_per_wait_rw_excl>([0-9]*[.])?[0-9]+) RW-sx)\n`,
	},
	"LATEST FOREIGN KEY ERROR": {
		`((?P<latest_foreign_key_error_time>\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2})\s+\d+\s+Transaction:)\n`,
	},
	"LATEST DETECTED DEADLOCK": {
		`((?P<latest_detecetd_deadlock_time>\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2})\s+\d+)\n`,
	},
	"TRANSACTIONS": {
		`(Trx id counter (?P<transactions_trx_id>\d+))\n`,
		`(Purge done for trx's n:o < (?P<transactions_trx_purged>\d+) undo n:o < (?P<transactions_trx_undo>\d+) state: running but idle)\n`,
		`(History list length (?P<transactions_history_list_lengt>\d+))\n`,
	},
	"FILE I/O": {
		`(Pending flushes \(fsync\) log: (?P<file_io_pending_flushes>\d+); buffer pool: (?P<file_io_pending_buffer_pool>\d+))\n`,
		`((?P<file_io_os_file_reads>\d+) OS file reads, (?P<file_io_os_file_writes>\d+) OS file writes, (?P<file_io_os_fsyncs>\d+) OS fsyncs)\n`,
		`((?P<file_io_reads_per_sec>([0-9]*[.])?[0-9]+) reads/s, (?P<file_io_avg_bytes_per_read>\d+) avg bytes/read, (?P<file_io_writes_per_se>([0-9]*[.])?[0-9]+) writes/s, (?P<file_io_fsyncs_per_sec>([0-9]*[.])?[0-9]+) fsyncs/s)\n`,
	},
	"INSERT BUFFER AND ADAPTIVE HASH INDEX": {
		`(Ibuf: size (?P<insert_buffer_size>\d+), free list len (?P<insert_buffer_free_list_len>\d+), seg size (?P<insert_buffer_seg_size>\d+), (?P<insert_buffer_merges>\d+) merges)\n`,
		`(merged operations:\n insert (?P<insert_buffer_merged_operations_insert>\d+), delete mark (?P<insert_buffer_merged_operations_delete_mark>\d+), delete (?P<insert_buffer_merged_operations_delete>\d+))\n`,
		`(discarded operations:\n insert (?P<insert_buffer_discarded_operations_insert>\d+), delete mark (?P<insert_buffer_discarded_operations_delete_mark>\d+), delete (?P<insert_buffer_discarded_operations_delete>\d+))\n`,
		`[^\>]*`,
		`((?P<insert_buffer_hash_searches_per_sec>([0-9]*[.])?[0-9]+) hash searches/s, (?P<insert_buffer_non_hash_searches_per_sec>([0-9]*[.])?[0-9]+) non-hash searches/s)\n`,
	},
	"LOG": {
		`(Log sequence number\s+(?P<log_sequence_number>\d+))\n`,
		`(Log buffer assigned up to\s+(?P<log_buffer_assigned_up_to>\d+))\n`,
		`(Log buffer completed up to\s+(?P<log_buffer_completed_up_to>\d+))\n`,
		`(Log written up to\s+(?P<log_written_up_to>\d+))\n`,
		`(Log flushed up to\s+(?P<log_flushed_up_to>\d+))\n`,
		`(Added dirty pages up to\s+(?P<log_added_dirty_page_up_to>\d+))\n`,
		`(Pages flushed up to\s+(?P<log_pages_flushed_up_to>\d+))\n`,
		`(Last checkpoint at\s+(?P<log_last_checkpoint_at>\d+))\n`,
		`(Log minimum file id is\s+(?P<log_minimum_file_id_is>\d+))\n`,
		`(Log maximum file id is\s+(?P<log_maximum_file_id_is>\d+))\n`,
		`((?P<log_ios_done>\d+) log i/o's done, (?P<log_ios_per_sec>([0-9]*[.])?[0-9]+) log i/o's/second)\n`,
	},
	"BUFFER POOL AND MEMORY": {
		`(Total large memory allocated\s+(?P<buffer_pool_total_large_memory_allocated>\d+))\n`,
		`(Dictionary memory allocated\s+(?P<buffer_pool_dictionary_memory_allocated>\d+))\n`,
		`(Buffer pool size\s+(?P<buffer_pool_buffer_pool_size>\d+))\n`,
		`(Free buffers\s+(?P<buffer_pool_free_buffers>\d+))\n`,
		`(Database pages\s+(?P<buffer_pool_database_pages>\d+))\n`,
		`(Old database pages\s+(?P<buffer_pool_old_database_pages>\d+))\n`,
		`(Modified db pages\s+(?P<buffer_pool_modified_db_pages>\d+))\n`,
		`(Pending reads\s+(?P<buffer_pool_pending_reads>\d+))\n`,
		`(Pending writes: LRU (?P<buffer_pool_pending_writes_lru>\d+), flush list (?P<buffer_pool_pending_writes_flush_list>\d+), single page (?P<buffer_pool_pending_writes_single_page>\d+))\n`,
		`(Pages made young (?P<buffer_pool_young_page>\d+), not young (?P<buffer_pool_not_young_page>\d+))\n`,
		`((?P<buffer_pool_youngs_per_sec>([0-9]*[.])?[0-9]+) youngs/s, (?P<buffer_pool_not_youngs_per_sec>([0-9]*[.])?[0-9]+) non-youngs/s)\n`,
		`(Pages read (?P<buffer_pool_pages_read>\d+), created (?P<buffer_pool_pages_created>\d+), written (?P<buffer_pool_pages_written>\d+))\n`,
		`((?P<buffer_pool_pages_read_per_sec>([0-9]*[.])?[0-9]+) reads/s, (?P<buffer_pool_pages_created_per_sec>([0-9]*[.])?[0-9]+) creates/s, (?P<buffer_pool_pages_written_per_sec>([0-9]*[.])?[0-9]+) writes/s)\n`,
		`(Buffer pool hit rate (?P<buffer_pool_hit_rate_min>\d+) / (?P<buffer_pool_hit_rate_max>\d+), young-making rate (?P<buffer_pool_young_making_rate_min>\d+) / (?P<buffer_pool_young_making_rate_max>\d+) not (?P<buffer_pool_not_young_making_rate_min>\d+) / (?P<buffer_pool_not_young_making_rate_max>\d+))\n`,
		`(Pages read ahead (?P<buffer_pool_pages_read_ahead_per_sec>([0-9]*[.])?[0-9]+)/s, evicted without access (?P<buffer_pool_pages_evicted_without_access_per_sec>([0-9]*[.])?[0-9]+)/s, Random read ahead (?P<buffer_pool_pages_random_read_ahead_per_sec>([0-9]*[.])?[0-9]+)/s)\n`,
		`(LRU len:\s+(?P<buffer_pool_lru_len>\d+), unzip_LRU len:\s+(?P<buffer_pool_unzip_lru_le>\d+))\n`,
		`(I/O sum\[(?P<buffer_pool_io_sum>\d+)\]:cur\[(?P<buffer_pool_io_cur>\d+)\], unzip sum\[(?P<buffer_pool_io_unzip_sum>\d+)\]:cur\[(?P<buffer_pool_io_unzip_cur>\d+)\])\n`,
	},
	"ROW OPERATIONS": {
		`((?P<row_operation_queries_inside_innodb>\d+) queries inside InnoDB, (?P<row_operation_queries_in_queue>\d+) queries in queue)\n`,
		`((?P<row_operation_read_views_open_inside_innodb>\d+) read views open inside InnoDB)\n`,
		`[^\>]*`,
		`(Number of rows inserted (?P<row_operation_inserted>\d+), updated (?P<row_operation_updated>\d+), deleted (?P<row_operation_deleted>\d+), read (?P<row_operation_read>\d+))\n`,
		`((?P<row_operation_inserts_per_sec>([0-9]*[.])?[0-9]+) inserts/s, (?P<row_operation_updates_per_sec>([0-9]*[.])?[0-9]+) updates/s, (?P<row_operation_deletes_per_sec>([0-9]*[.])?[0-9]+) deletes/s, (?P<row_operation_reads_per_sec>([0-9]*[.])?[0-9]+) reads/s)\n`,
		`(Number of system rows inserted (?P<row_operation_system_inserted>\d+), updated (?P<row_operation_system_updated>\d+), deleted (?P<row_operation_system_deleted>\d+), read (?P<row_operation_system_read>\d+))\n`,
		`((?P<row_operation_system_inserts_per_sec>([0-9]*[.])?[0-9]+) inserts/s, (?P<row_operation_system_updates_per_sec>([0-9]*[.])?[0-9]+) updates/s, (?P<row_operation_system_deletes_per_sec>([0-9]*[.])?[0-9]+) deletes/s, (?P<row_operation_system_reads_per_sec>([0-9]*[.])?[0-9]+) reads/s)\n`,
	},
}

func CollectInnoDB() {
	MySQLConn.Connect()
	MySQLConn.FetchAll(SQLInnoDB, func(row map[string]string) {
		ParseSection(row["Status"], "BACKGROUND THREAD")
		ParseSection(row["Status"], "SEMAPHORES")
		ParseSection(row["Status"], "LATEST FOREIGN KEY ERROR")
		ParseSection(row["Status"], "LATEST DETECTED DEADLOCK")
		ParseSection(row["Status"], "TRANSACTIONS")
		ParseSection(row["Status"], "FILE I/O")
		ParseSection(row["Status"], "INSERT BUFFER AND ADAPTIVE HASH INDEX")
		ParseSection(row["Status"], "LOG")
		ParseSection(row["Status"], "BUFFER POOL AND MEMORY")
		ParseSection(row["Status"], "ROW OPERATIONS")
	})
}

func ParseSection(in, name string) {
	writer := InfluxDBConn.WriteAPIBlocking("debeando", InfluxDBBucket)
	expression := strings.Join(Sections[name], "")
	pattern := regexp.MustCompile(expression)
	matches := pattern.FindAllStringSubmatch(in, -1)
	keys := pattern.SubexpNames()

	if len(matches) != 1 {
		return
	}

	for index, key := range keys {
		if key == "" {
			continue
		}

		if value, ok := mysql.ParseNumberValue(matches[0][index]); ok {
			log.DebugWithFields("MySQL InnoDB", log.Fields{
				"hostname": MySQLHost,
				key:        value,
			})

			p := influxdb2.NewPointWithMeasurement("mysql_innodb").
				AddTag("_hostname", MySQLHost).
				AddField(key, value).
				SetTime(time.Now())

			err := writer.WritePoint(context.Background(), p)
			if err != nil {
				log.Error(err.Error())
			}
		}
	}
}
