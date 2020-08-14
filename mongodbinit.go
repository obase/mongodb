package mongodb

import (
	"fmt"
	"github.com/obase/conf"
	"strings"
)

const CKEY = "mongodb"

// 对接conf.yml, 读取原redis相关配置
func init() {

	configs, ok := conf.GetSlice(CKEY)
	if !ok || len(configs) == 0 {
		return
	}

	for _, config := range configs {
		if key, ok := conf.ElemString(config, "key"); ok {

			address, _ := conf.ElemStringSlice(config, "address")
			database, _ := conf.ElemString(config, "database")
			username, _ := conf.ElemString(config, "username")
			password, _ := conf.ElemString(config, "password")
			source, _ := conf.ElemString(config, "source")
			readPreference, _ := GetReadPreference(conf.Elem(config, "readPreference"))
			readConcern, _ := GetReadConcern(conf.Elem(config, "readConcern"))
			writeConcern, _ := GetWriteConcern(conf.Elem(config, "writeConcern"))

			direct, _ := conf.ElemBool(config, "direct")
			replicaSet, _ := conf.ElemString(config, "replicaSet")
			keepalive, _ := conf.ElemDuration(config, "keepalive")
			connectTimeout, _ := conf.ElemDuration(config, "connectTimeout")
			serverSelectionTimeout, _ := conf.ElemDuration(config, "serverSelectionTimeout")
			socketTimeout, _ := conf.ElemDuration(config, "socketTimeout")
			heartbeatInterval, _ := conf.ElemDuration(config, "heartbeatInterval")
			localThreshold, _ := conf.ElemDuration(config, "localThreshold")

			minPoolSize, _ := conf.ElemInt64(config, "minPoolSize")
			maxPoolSize, _ := conf.ElemInt64(config, "maxPoolSize")
			maxConnIdleTime, _ := conf.ElemDuration(config, "maxConnIdleTime")

			compressors, _ := conf.ElemStringSlice(config, "compressors")
			zlibLevel, _ := conf.ElemInt(config, "zlibLevel")
			zstdLevel, _ := conf.ElemInt(config, "zstdLevel")
			retryReads, _ := conf.ElemBool(config, "retryReads")
			retryWrites, _ := conf.ElemBool(config, "retryWrites")

			if err := Setup(key, &Config{
				Address:                address,
				Database:               database,
				Username:               username,
				Password:               password,
				Source:                 source,
				ReadPreference:         readPreference,
				ReadConcern:            readConcern,
				WriteConcern:           writeConcern,
				Direct:                 direct,
				ReplicaSet:             replicaSet,
				Keepalive:              keepalive,
				ConnectTimeout:         connectTimeout,
				ServerSelectionTimeout: serverSelectionTimeout,
				SocketTimeout:          socketTimeout,
				HeartbeatInterval:      heartbeatInterval,
				LocalThreshold:         localThreshold,
				MinPoolSize:            uint64(minPoolSize),
				MaxPoolSize:            uint64(maxPoolSize),
				MaxConnIdleTime:        maxConnIdleTime,
				Compressors:            compressors,
				ZlibLevel:              zlibLevel,
				ZstdLevel:              zstdLevel,
				RetryReads:             retryReads,
				RetryWrites:            retryWrites,
			}); err != nil {
				panic(err)
			}
		}
	}
}

func GetReadPreference(val interface{}, ok bool) (*ReadPreference, bool) {
	switch val := val.(type) {
	case nil:
		return nil, true
	case string:
		val = strings.TrimSpace(val)
		if val == "" {
			return nil, true
		}
		return &ReadPreference{RMode: val}, true
	case map[string]interface{}:
		return &ReadPreference{
			RMode:         conf.ToString(val["RMode"]),
			RTagSet:       conf.ToStringMap(val["RTagSet"]),
			RMaxStateness: conf.ToDuration(val["RMaxStateness"]),
		}, true
	case map[interface{}]interface{}:
		ret := new(ReadPreference)
		for k, v := range val {
			switch conf.ToString(k) {
			case "RMode":
				ret.RMode = conf.ToString(v)
			case "RTagSet":
				ret.RTagSet = conf.ToStringMap(v)
			case "RMaxStateness":
				ret.RMaxStateness = conf.ToDuration(v)
			}
		}
		return ret, true
	default:
		panic(fmt.Sprintf("invalid value for read preference: %v", val))
	}
}

func GetReadConcern(val interface{}, ok bool) (*ReadConcern, bool) {
	switch val := val.(type) {
	case nil:
		return nil, true
	case string:
		val = strings.TrimSpace(val)
		if val == "" {
			return nil, true
		}
		return &ReadConcern{Level: val}, true
	case map[string]interface{}:
		return &ReadConcern{
			Level: conf.ToString(val["Level"]),
		}, true
	case map[interface{}]interface{}:
		ret := new(ReadConcern)
		for k, v := range val {
			switch conf.ToString(k) {
			case "Level":
				ret.Level = conf.ToString(v)
			}
		}
		return ret, true
	default:
		panic(fmt.Sprintf("invalid value for read concern: %v", val))
	}
}

func GetWriteConcern(val interface{}, ok bool) (*WriteConcern, bool) {
	switch val := val.(type) {
	case nil:
		return nil, true
	case string:
		val = strings.TrimSpace(val)
		if val == "" {
			return nil, true
		}
		switch strings.ToLower(val) {
		case WriteConcern_majority:
			return &WriteConcern{WMajority: true}, true
		case WriteConcern_w0:
			return &WriteConcern{W: 0}, true
		case WriteConcern_w1:
			return &WriteConcern{W: 1}, true
		case WriteConcern_w2:
			return &WriteConcern{W: 2}, true
		case WriteConcern_w3:
			return &WriteConcern{W: 3}, true
		default:
			panic(fmt.Sprintf("invalid value for write concern: %v", val))
		}
	case map[string]interface{}:
		return &WriteConcern{
			WMajority: conf.ToBool(val["WMajority"]),
			W:         conf.ToInt(val["W"]),
			J:         conf.ToBool(val["J"]),
			WTagSet:   conf.ToString(val["WTagSet"]),
			WTimeout:  conf.ToDuration(val["WTimeout"]),
		}, true
	case map[interface{}]interface{}:
		ret := new(WriteConcern)
		for k, v := range val {
			switch conf.ToString(k) {
			case "J":
				ret.J = conf.ToBool(v)
			case "W":
				ret.W = conf.ToInt(v)
			case "WMajority":
				ret.WMajority = conf.ToBool(v)
			case "WTagSet":
				ret.WTagSet = conf.ToString(v)
			case "WTimeout":
				ret.WTimeout = conf.ToDuration(v)
			}
		}
		return ret, true
	default:
		panic(fmt.Sprintf("invalid value for read concern: %v", val))
	}
}
