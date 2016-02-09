package redis

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Info holds the data returned from redis' INFO command.
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

// ServerInfo holds the data inside the "# Server" portion
// of redis' INFO command.
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

// ClientsInfo holds the data inside the "# Clients" portion
// of redis' INFO command.
type ClientsInfo struct {
	ConnectedClients        int64 `redis:"connected_clients"`
	ClientLongestOutputList int64 `redis:"client_longest_output_list"`
	ClientBiggestInputBuf   int64 `redis:"client_biggest_input_buf"`
	BlockedClients          int64 `redis:"blocked_clients"`
}

// MemoryInfo holds the data inside the "# Memory" portion
// of redis' INFO command.
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

// PersistenceInfo holds the data inside the "# Persistence" portion
// of redis' INFO command.
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

// StatsInfo holds the data inside the "# Stats" portion
// of redis' INFO command.
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

// ReplicationInfo holds the data inside the "# Replication" portion
// of redis' INFO command.
type ReplicationInfo struct {
	Role                       string `redis:"role"`
	ConnectedSlaves            int64  `redis:"connected_slaves"`
	MasterReplOffset           int64  `redis:"master_repl_offset"`
	ReplBacklogActive          int64  `redis:"repl_backlog_active"`
	ReplBacklogSize            int64  `redis:"repl_backlog_size"`
	ReplBacklogFirstByteOffset int64  `redis:"repl_backlog_first_byte_offset"`
	ReplBacklogHistlen         int64  `redis:"repl_backlog_histlen"`
}

// CPUInfo holds the data inside the "# CPU" portion
// of redis' INFO command.
type CPUInfo struct {
	UsedCPUSys          float64 `redis:"used_cpu_sys"`
	UsedCPUUser         float64 `redis:"used_cpu_user"`
	UsedCPUSysChildren  float64 `redis:"used_cpu_sys_children"`
	UsedCPUUserChildren float64 `redis:"used_cpu_user_children"`
}

// ClusterInfo holds the data inside the "# Cluster" portion
// of redis' INFO command.
type ClusterInfo struct {
	ClusterEnabled bool `redis:"cluster_enabled"`
}

// KeyspaceInfo holds the data inside the "# Keyspace" portion
// of redis' INFO command.
type KeyspaceInfo map[string]Keyspace

// Keyspace is the data contained inside each member of the
// "# Keyspace" section.
type Keyspace struct {
	Keys    int64 `redis:"keys"`
	Expires int64 `redis:"expires"`
	AvgTTL  int64 `redis:"avg_ttl"`
}

// Set implements the Setter interface.
func (k KeyspaceInfo) Set(key, data string) error {

	if data == "" || key == "" {
		return errors.New("redigo: could not parse keyspace because key or data was nil")
	}

	parts := strings.Split(data, ",")

	if len(parts) == 0 {
		return errors.New("redigo: could not parse keyspace")
	}

	var ks Keyspace
	kss := reflect.ValueOf(&ks).Elem()
	t := kss.Type()
	var field reflect.Value
	for _, part := range parts {

		kv := strings.Split(part, "=")
		if len(kv) < 2 {
			return errors.New("redigo: could not parse key=value pairs")
		}

		i, err := strconv.ParseInt(kv[1], 10, 64)
		if err != nil {
			return err
		}

		for i := 0; i < kss.NumField(); i++ {
			tag := t.Field(i).Tag.Get("redis")
			if tag == kv[0] {
				field = kss.Field(i)
				break
			}
		}

		field.SetInt(i)
	}
	k[key] = ks
	return nil
}

// Setter is the interface that the maps inside type Info should
// implement. Setter should set the map's key with potentially parsed data.
// This allows for better generalization of arbitrary-sized (e.g.
// maps, slices) types.
// For instance, the Keyspace fields are the names of each database.
// The values of each field are comma-delimited k=v pairs. However,
// future fields (the major "# <name>" fields) may not use k=v
// pairs, so it makes no sense to hardcode it inside our ParseInfo
// function.
type Setter interface {
	Set(key, data string) error
}

var setter = reflect.TypeOf((*Setter)(nil)).Elem()

// ParseInfo is a helper that converts a command reply to a *Info.
// If err is not equal to nil, then ParseInfo returns nil, err.
// Otherwise, ParseInfo converts the reply to a *Info as follows:
//
// 	Reply type 	Result
// 	string		*Info, nil
// 	[]byte		*Info, nil
// 	nil			nil, ErrNil
// 	other		nil, error
func ParseInfo(reply interface{}, err error) (*Info, error) {
	// Return a pointer because Info is nearly 1KB in size.
	// Future implementations will only grow the struct.
	// I haven't profiled, but 1KB is pretty big.

	if err != nil {
		return nil, err
	}

	var s *bufio.Scanner
	switch reply := reply.(type) {
	case string:
		s = bufio.NewScanner(strings.NewReader(reply))
	case []byte:
		s = bufio.NewScanner(bytes.NewReader(reply))
	case nil:
		return nil, ErrNil
	case Error:
		return nil, reply
	default:
		return nil, fmt.Errorf("redigo: unexpected type for ParseInfo, got type %T", reply)
	}

	// curField is the current field inside the Info struct
	var curField reflect.Value
	var info Info
	val := reflect.ValueOf(&info).Elem()
	for s.Scan() {
		text := s.Text()

		// Skip empty lines.
		if text == "" || text == "\n" {
			continue
		}

		// Beginning of portion new field in Info struct.
		if strings.HasPrefix(text, "# ") {
			parts := strings.Split(text, "# ")
			if len(parts) < 2 {
				return nil, errors.New("redigo: line started with comment but didn't have field name")
			}

			curField = val.FieldByName(parts[1])
			t := curField.Type()

			// Initialize containers if needed.
			if t.Kind() == reflect.Map {
				curField.Set(reflect.MakeMap(t))
			}

			// Continue on to next subfields.
			continue
		}

		parts := strings.Split(text, ":")
		if len(parts) < 2 {
			return nil, errors.New("redigo: data could not be split on ':'")
		}

		switch t := curField.Type(); t.Kind() {
		case reflect.Struct:
			for i := 0; i < curField.NumField(); i++ {
				tag := t.Field(i).Tag.Get("redis")
				if tag == parts[0] {
					err := setField(curField.Field(i), parts[1])
					if err != nil {
						return nil, err
					}
					break
				}
			}
		case reflect.Map:
			if t.Implements(setter) {
				fn := curField.MethodByName("Set")
				ret := fn.Call([]reflect.Value{
					reflect.ValueOf(parts[0]),
					reflect.ValueOf(parts[1]),
				})
				if len(ret) != 1 {
					return nil, errors.New("redigo: cannot parse data")
				}
				if !ret[0].IsNil() {
					err, ok := ret[0].Interface().(error)
					if !ok {
						return nil, errors.New("redigo: bug where ret[0].Interface().(error) == false")
					}
					return nil, err
				}
			}

		}

	}
	return &info, s.Err()
}

func setField(field reflect.Value, val string) error {
	switch t := field.Type().Kind(); t {
	case reflect.String:
		field.SetString(val)
	case reflect.Int, reflect.Int8,
		reflect.Int16, reflect.Int32,
		reflect.Int64:
		i, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(i)
	case reflect.Uint, reflect.Uint8,
		reflect.Uint16, reflect.Uint32,
		reflect.Uint64:
		i, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return err
		}
		field.SetUint(i)
	case reflect.Bool:
		x, err := strconv.ParseBool(val)
		if err != nil {
			return err
		}
		field.SetBool(x)
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return err
		}
		field.SetFloat(f)
	default:
		return fmt.Errorf("redigo: cannot set field of type %T", t)
	}
	return nil
}
