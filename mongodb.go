package mongodb

import (
	"fmt"
	"github.com/obase/conf"
	"time"
)

const (
	ReadPreference_primary            = "primary"
	ReadPreference_primaryPreferred   = "primaryPreferred"
	ReadPreference_secondary          = "secondary"
	ReadPreference_secondaryPreferred = "secondaryPreferred"
	ReadPreference_nearest            = "nearest"

	ReadConcern_available   = "available"
	ReadConcern_local       = "local"
	ReadConcern_majority    = "majority"
	ReadConcern_linerizable = "linerizable"

	WriteConcern_majority = "majority" // 写到大多数primary结点
	WriteConcern_w0       = "w0"       // 写后不理
	WriteConcern_w1       = "w1"       // 写1台后返回
	WriteConcern_w2       = "w2"       // 写2台后返回
	WriteConcern_w3       = "w3"       // 写3台后返回

)

// 读优先
type ReadPreference struct {
	RMode         string            // 读模式, primary | primaryPreferred | secondary | secondaryPreferred | nearest
	RTagSet       map[string]string // 读标签, 支持 k1:v1,k2:v2,...的格式
	RMaxStateness time.Duration     // specify a maxinum replication lag for reads from secondaries in a replica set
}

// 读安全
type ReadConcern struct {
	Level string `json:"level" bson:"level" yaml:"level"`
}

// 写安全, 如果WMajority为true则忽略W, J表示journal
type WriteConcern struct {
	J         bool          `json:"J" yaml:"J"`                 // write operations are written to the journal
	W         int           `json:"W" yaml:"W"`                 // write operations propagate to the specified number of mongod instances
	WMajority bool          `json:"WMajority" yaml:"WMajority"` // write operations propagate to the majority of mongod instances
	WTagSet   string        `json:"WTagSet" yaml:"WTagSet"`     // write operations propagate to the specified mongod instance
	WTimeout  time.Duration `json:"WTimeout" yaml:"WTimeout"`   // specifies a time limit for the write concern
}

type Config struct {
	// 连接URL, 格式为[mongodb://][user:pass@]host1[:port1][,host2[:port2],...][/database][?options]
	// 常用选项
	Address        []string        `json:"address" yaml:"address"`
	Database       string          `json:"database" yaml:"database"`             // 默认DB
	Username       string          `json:"username" yaml:"username"`             // 用户名
	Password       string          `json:"password" yaml:"password"`             // 密码
	Source         string          `json:"source" yaml:"source"`                 // 授权DB, 默认为admin
	ReadPreference *ReadPreference `json:"readPreference" yaml:"readPreference"` // 读优先级, primary|primaryPreferred|secondary|secondaryPreferred|nearest
	ReadConcern    *ReadConcern    `json:"readConcern" yaml:"readConcern"`       // 读影响， available|local|majority|linerizable
	WriteConcern   *WriteConcern   `json:"writeConcern" yaml:"writeConcern"`     // 写影响

	// 连接管理
	Direct                 bool          `json:"direct" yaml:"direct"`                                 //是否直接
	ReplicaSet             string        `json:"replicaSet" yaml:"replicaSet"`                         //数据集
	Keepalive              time.Duration `json:"keepalive" yaml:"keepalive"`                           //默认300秒
	ConnectTimeout         time.Duration `json:"connectTimeout" yaml:"connectTimeout"`                 //连接超时. 默认为10秒
	ServerSelectionTimeout time.Duration `json:"serverSelectionTimeout" yaml:"serverSelectionTimeout"` // 服务端选择超时, 默认为30秒
	SocketTimeout          time.Duration `json:"socketTimeout" yaml:"socketTimeout"`                   //写超时, 默认为0, 表示阻塞
	HeartbeatInterval      time.Duration `json:"heartbeatInterval" yaml:"heartbeatInterval"`           // 心跳间隔, 默认为10秒
	LocalThreshold         time.Duration `json:"localThreshold" yaml:"localThreshold"`                 // 延迟窗口, 默认为15毫秒

	// 连接池管理
	MinPoolSize     uint64        `json:"minPoolSize" yaml:"minPoolSize"`         // 连接池最小连接数
	MaxPoolSize     uint64        `json:"maxPoolSize" yaml:"maxPoolSize"`         // 连接池最大连接数
	MaxConnIdleTime time.Duration `json:"maxConnIdleTime" yaml:"maxConnIdleTime"` // 连接池最大空闲时间

	// 压缩通信
	Compressors []string `json:"compressors" yaml:"compressors"` // 压缩算法, snappy(3.4), zlib(3.6), zstd(4.2)
	ZlibLevel   int      `json:"zlibLevel" yaml:"zlibLevel"`     // zlib压缩等级
	ZstdLevel   int      `json:"zstdLevel" yaml:"zstdLevel"`     // zstd压缩等级

	// 重试机制
	RetryReads  bool `json:"retryReads" yaml:"retryReads"`   // 重试读(3.6)
	RetryWrites bool `json:"retryWrites" yaml:"retryWrites"` // 重试写(3.6)
}

var (
	clients map[string]*Client = make(map[string]*Client)
)

func Setup(key string, cnf *Config) (err error) {

	keys := conf.ToStringSlice(key)
	for _, k := range keys {
		if _, ok := clients[k]; ok {
			return fmt.Errorf("duplicate mongodb client: %v", k)
		}
	}

	if cnf == nil {
		cnf = new(Config)
	}

	client, err := newClient(cnf)
	if err != nil {
		return
	}
	for _, k := range keys {
		clients[k] = client
	}
	return
}

func Get(key string) *Client {
	return clients[key]
}

func Must(key string) *Client {
	ret, ok := clients[key]
	if !ok {
		panic("invalid mongodb client: " + key)
	}
	return ret
}
