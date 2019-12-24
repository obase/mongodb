package mongodb

import (
	"errors"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"strings"
	"time"
)

type Client interface {
	Close() error
	Count(c string) (int64, error)
	CountWith(c string, filter interface{}) (int64, error)
}

// Safe session safety mode. See SetSafe for details on the Safe type.
type Safe struct {
	W        int    // Min # of servers to ack before success
	WMode    string // Write mode for MongoDB 2.0+ (e.g. "majority")
	RMode    string // Read mode for MonogDB 3.2+ ("majority", "local", "linearizable")
	WTimeout int    // Milliseconds to wait for W before timing out
	FSync    bool   // Sync via the journal if present, or via data files sync otherwise
	J        bool   // Sync via the journal if present
}

type Config struct {
	// 连接URL, 格式为[mongodb://][user:pass@]host1[:port1][,host2[:port2],...][/database][?options]
	Address     []string
	Database    string
	Username    string
	Password    string
	Source      string
	Safe        *Safe
	Mode        readpref.Mode
	Compressors []string

	// 连接管理
	ConnectTimeout time.Duration //连接超时. 默认为10秒
	Keepalive      time.Duration //DialInfo.DialServer实现
	WriteTimeout   time.Duration //写超时, 默认为ConnectTimeout
	ReadTimeout    time.Duration //读超时, 默认为ConnectTimeout

	// 连接池管理
	MinPoolSize       int //对应DialInfo.MinPoolSize
	MaxPoolSize       int //对应DialInfo.PoolLimit
	MaxPoolWaitTimeMS int //对应DialInfo.PoolTimeout获取连接超时, 默认为0永不超时
	MaxPoolIdleTimeMS int //对应DialInfo.MaxIdleTimeMS
}

var (
	clients map[string]Client = make(map[string]Client)
)

func mergeConfig(opt *Config) *Config {
	if opt == nil {
		opt = new(Config)
	}
	return opt
}

func Setup(key string, cnf *Config) (err error) {

	client, err := newMongodbClient(mergeConfig(cnf))
	if err != nil {
		return
	}
	for _, k := range strings.Split(key, ",") {
		k = strings.TrimSpace(k)
		if _, ok := clients[k]; ok {
			return errors.New("duplicate mongodb client: " + k)
		}
		clients[k] = client
	}
	return
}

func Get(key string) Client {
	return clients[key]
}

func Must(key string) Client {
	ret, ok := clients[key]
	if !ok {
		panic("missing mongodb client: " + key)
	}
	return ret
}
