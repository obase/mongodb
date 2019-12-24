package mongodb

import (
	"fmt"
	"github.com/obase/conf"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const CKEY = "mongo"

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
			compressors, _ := conf.ElemStringSlice(config, "compressors")
			safe, _ := getSafe(conf.ElemMap(config, "safe"))
			mode, _ := getMode(conf.Elem(config, "mode"))
			keepalive, _ := conf.ElemDuration(config, "keepalive")
			connectTimeout, _ := conf.ElemDuration(config, "connectTimeout")
			readTimeout, _ := conf.ElemDuration(config, "readTimeout")
			writeTimeout, _ := conf.ElemDuration(config, "writeTimeout")
			minPoolSize, _ := conf.ElemInt(config, "minPoolSize")
			maxPoolSize, _ := conf.ElemInt(config, "maxPoolSize")
			maxPoolWaitTimeMS, _ := conf.ElemInt(config, "maxPoolWaitTimeMS")
			maxPoolIdleTimeMS, _ := conf.ElemInt(config, "maxPoolIdleTimeMS")

			if err := Setup(key, &Config{
				Address:           address,
				Database:          database,
				Username:          username,
				Password:          password,
				Source:            source,
				Compressors:       compressors,
				Safe:              safe,
				Mode:              mode,
				Keepalive:         keepalive,
				ConnectTimeout:    connectTimeout,
				ReadTimeout:       readTimeout,
				WriteTimeout:      writeTimeout,
				MinPoolSize:       minPoolSize,
				MaxPoolSize:       maxPoolSize,
				MaxPoolWaitTimeMS: maxPoolWaitTimeMS,
				MaxPoolIdleTimeMS: maxPoolIdleTimeMS,
			}); err != nil {
				panic(err)
			}
		}
	}
}

func getMode(val interface{}, ok bool) (readpref.Mode, bool) {
	if !ok {
		return readpref.PrimaryMode, false
	}
	switch val := val.(type) {
	case string:
		ret, err := readpref.ModeFromString(val)
		if err != nil {
			return readpref.PrimaryMode, false
		}
		return ret, true
	case int:
		return readpref.Mode(val), true
	case int64:
		return readpref.Mode(val), true
	}
	panic("unsupport mode type: " + fmt.Sprint(val))
}

func getSafe(val map[string]interface{}, ok bool) (*Safe, bool) {
	safe := &Safe{
		WMode: "majority",
	}
	for k, v := range val {
		switch k {
		case "W", "w":
			safe.W = conf.ToInt(v)
		case "WMode", "wmode":
			safe.WMode = conf.ToString(v)
		case "RMode", "rmode":
			safe.RMode = conf.ToString(v)
		case "WTimeout", "wtimeout":
			safe.WTimeout = conf.ToInt(v)
		case "FSync", "fsync":
			safe.FSync = conf.ToBool(v)
		case "J", "j":
			safe.J = conf.ToBool(v)
		}
	}
	return safe, true
}
