package redis

import (
	"bufio"
	"bytes"
	"errors"
	"reflect"
	"strconv"
	"strings"
)

type Info struct {
	Server      ServerInfo
	Clients     ClientsInfo
	Memory      MemoryInfo
	Persistence PersistenceInfo
	Stats       StatsInfo
	Replication ReplicationInfo
	CPU         CPUInfo
	Cluster     ClusterInfo
	Keyspace    KeyspaceInfo
}

type ServerInfo struct {
	RedisVersion    string `redis:"redis_version"`
	RedisGitSHA1    string `redis:"redis_git_sha1"`
	RedisGitDirty   bool   `redis:"redis_git_dirty"`
	RedisBuildID    string `redis:"redis_build_id"`
	RedisMode       string `redis:"redis_mode"`
	OS              string `redis:"os"`
	ArchBits        byte   `redis:"arch_bits"`
	MultiplexingAPI string `redis:"multiplexing_api"`
	GCCVersion      string `redis:"gcc_version"`
	ProcessID       uint32 `redis:"process_id"`
	RunID           string `redis:"run_id"`
	TCPPort         string `redis:"tcp_port"`
	UptimeInSeconds int64  `redis:"uptime_in_seconds"`
	UptimeInDays    int64  `redis:"uptime_in_days"`
	Hz              string `redis:"hz"`
	LRUClock        int64  `redis:"lru_clock"`
	ConfigFile      string `redis:"config_file"`
}

type ClientsInfo struct {
	ConnectedClients        int64 `redis:"connected_clients"`
	ClientLongestOutputList int64 `redis:"client_longest_output_list"`
	ClientBiggestInputBuf   int64 `redis:"client_biggest_input_buf"`
	BlockedClients          int64 `redis:"blocked_clients"`
}

type MemoryInfo struct {
	UsedMemory            int64   `redis:"used_memory"`
	UsedMemoryHuman       string  `redis:"used_memory_human"`
	UsedMemoryRSS         int64   `redis:"used_memory_rss"`
	UsedMemoryPeak        int64   `redis:"used_memory_peak"`
	UsedMemoryPeakHuman   string  `redis:"used_memory_peak_human"`
	UsedMemoryLua         int64   `redis:"used_memory_lua"`
	MemFragmentationRatio float64 `redis:"mem_fragmentation_ratio"`
	MemAllocator          string  `redis:"mem_allocator"`
}

type PersistenceInfo struct {
	Loading                  bool   `redis:"loading"`
	RDBChangesSinceLastSave  int64  `redis:"rdb_changes_since_last_save"`
	RDBBgsaveInProgress      int64  `redis:"rdb_bgsave_in_progress"`
	RDBLastSaveTime          int64  `redis:"rdb_last_save_time"`
	RDBLastBGSaveStatus      string `redis:"rdb_last_bgsave_status"`
	RDBLastBGSaveTimeSec     int64  `redis:"rdb_last_bgsave_time_sec"`
	RDBCurrentBGSaveTimeSec  int64  `redis:"rdb_current_bgsave_time_sec"`
	AOFEnabled               bool   `redis:"aof_enabled"`
	AOFRewriteInProgress     bool   `redis:"aof_rewrite_in_progress"`
	AOFRewriteScheduled      bool   `redis:"aof_rewrite_scheduled"`
	AOFLastRewriteTimeSec    int64  `redis:"aof_last_rewrite_time_sec"`
	AOFCurrentRewriteTimeSec int64  `redis:"aof_current_rewrite_time_sec"`
	AOFLastBGRewriteStatus   string `redis:"aof_last_bgrewrite_status"`
	AOFLastWriteStatus       string `redis:"aof_last_write_status"`
}

type StatsInfo struct {
	TotalConnectionsReceived int64   `redis:"total_connections_received"`
	TotalCommandsProcessed   int64   `redis:"total_commands_processed"`
	InstantaneousOpsPerSec   int64   `redis:"instantaneous_ops_per_sec"`
	TotalNetInputBytes       int64   `redis:"total_net_input_bytes"`
	TotalNetOutputBytes      int64   `redis:"total_net_output_bytes"`
	InstantaneousInputKbps   float64 `redis:"instantaneous_input_kbps"`
	InstantaneousOutputKbps  float64 `redis:"instantaneous_output_kbps"`
	RejectedConnections      int64   `redis:"rejected_connections"`
	SyncFull                 bool    `redis:"sync_full"`
	SyncPartialOK            bool    `redis:"sync_partial_ok"`
	SyncPartialErr           bool    `redis:"sync_partial_err"`
	ExpiredKeys              int64   `redis:"expired_keys"`
	EvictedKeys              int64   `redis:"evicted_keys"`
	KeyspaceHits             int64   `redis:"keyspace_hits"`
	KeyspaceMisses           int64   `redis:"keyspace_misses"`
	PubsubChannels           int64   `redis:"pubsub_channels"`
	PubsubPatterns           int64   `redis:"pubsub_patterns"`
	LatestForkUsec           int64   `redis:"latest_fork_usec"`
	MigrateCachedSockets     bool    `redis:"migrate_cached_sockets"`
}

type ReplicationInfo struct {
	Role                       string `redis:"role"`
	ConnectedSlaves            int64  `redis:"connected_slaves"`
	MasterReplOffset           int64  `redis:"master_repl_offset"`
	ReplBacklogActive          int64  `redis:"repl_backlog_active"`
	ReplBacklogSize            int64  `redis:"repl_backlog_size"`
	ReplBacklogFirstByteOffset int64  `redis:"repl_backlog_first_byte_offset"`
	ReplBacklogHistlen         int64  `redis:"repl_backlog_histlen"`
}

type CPUInfo struct {
	UsedCPUSys          float64 `redis:"used_cpu_sys"`
	UsedCPUUser         float64 `redis:"used_cpu_user"`
	UsedCPUSysChildren  float64 `redis:"used_cpu_sys_children"`
	UsedCPUUserChildren float64 `redis:"used_cpu_user_children"`
}

type ClusterInfo struct {
	ClusterEnabled bool `redis:"cluster_enabled"`
}

type KeyspaceInfo map[string]Keyspace

func (k *Keyspace) Parse(data string) error {

	if data == "" {
		return errors.New("redigo: could not parse keyspace")
	}

	parts := strings.Split(data, ",")

	if len(parts) == 0 {
		return errors.New("redigo: could not parse keyspace")
	}

	// var ks Keyspace

	vf := reflect.ValueOf(k).Elem()
	t := vf.Type()
	var ff reflect.Value
	for _, part := range parts {

		kv := strings.Split(part, "=")

		i, err := strconv.ParseInt(kv[1], 10, 64)
		if err != nil {
			return err
		}

		for i := 0; i < vf.NumField(); i++ {
			tag := t.Field(i).Tag.Get("redis")
			if tag == kv[0] {
				ff = vf.Field(i)
				break
			}
		}

		ff.SetInt(i)
	}
	return nil
}

type Keyspace struct {
	Keys    int64 `redis:"keys"`
	Expires int64 `redis:"expires"`
	AvgTTL  int64 `redis:"avg_ttl"`
}

func ParseInfo(v interface{}, e error) (*Info, error) {
	if v == nil {
		return nil, errors.New("redigo: no data to parse")
	}

	var s *bufio.Scanner
	switch t := v.(type) {
	case string:
		s = bufio.NewScanner(strings.NewReader(t))
	case []byte:
		s = bufio.NewScanner(bytes.NewReader(t))
	default:
		return nil, errors.New("redigo: cannot parse data")
	}

	var vf, ff reflect.Value
	var info Info
	val := reflect.ValueOf(&info).Elem()
	for s.Scan() {
		text := s.Text()

		if text == "" {
			continue
		}

		// Beginning of portion new field in Info struct.
		if strings.HasPrefix(text, "# ") {
			parts := strings.Split(text, "# ")
			if len(parts) < 2 {
				return nil, errors.New("redigo: cannot parse data")
			}
			vf = val.FieldByNameFunc(func(s string) bool {
				return s == parts[1]
			})
			if vf.Type().Kind() == reflect.Map {
				vf = reflect.MakeMap(vf.Type())
				val.FieldByName(parts[1]).Set(vf)
			}
			continue
		}

		parts := strings.Split(text, ":")
		if len(parts) < 2 {
			return nil, errors.New("redigo: cannot parse data")
		}

		t := vf.Type()
		if t.Kind() == reflect.Struct {
			for i := 0; i < vf.NumField(); i++ {
				tag := t.Field(i).Tag.Get("redis")
				if tag == parts[0] {
					ff = vf.Field(i)
					break
				}
			}
		} else if t.Kind() == reflect.Map {
			rv := reflect.New(t.Elem())
			fn := rv.MethodByName("Parse")
			ret := fn.Call([]reflect.Value{reflect.ValueOf(parts[1])})
			if len(ret) != 1 {
				return nil, errors.New("redigo: cannot parse data")
			}
			if !ret[0].IsNil() {
				return nil, ret[0].Interface().(error)
			}
			vf.SetMapIndex(reflect.ValueOf(parts[0]), rv.Elem())

			continue
		}

		switch ff.Type().Kind() {
		case reflect.String:
			ff.SetString(parts[1])
		case reflect.Int, reflect.Int8,
			reflect.Int16, reflect.Int32, reflect.Int64:
			i, err := strconv.ParseInt(parts[1], 10, 64)
			if err != nil {
				return nil, err
			}
			ff.SetInt(i)
		case reflect.Uint, reflect.Uint8,
			reflect.Uint16, reflect.Uint32, reflect.Uint64:
			i, err := strconv.ParseUint(parts[1], 10, 64)
			if err != nil {
				return nil, err
			}
			ff.SetUint(i)
		case reflect.Bool:
			x, err := strconv.ParseBool(parts[1])
			if err != nil {
				return nil, err
			}
			ff.SetBool(x)
		case reflect.Float32, reflect.Float64:
			f, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				return nil, err
			}
			ff.SetFloat(f)
		default:
			return nil, errors.New("redigo: cannot parse data")
		}
	}
	return &info, s.Err()
}
