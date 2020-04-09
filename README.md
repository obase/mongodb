# package mongodb
集成官方驱动mongo-go-driver

# Installation
- go get
```
go get -u go.mongodb.org/mongo-driver
go get -u github.com/obase/mongodb
```
- go mod
```
go mod edit -require=github.com/obase/mongodb@latest
```

# Configuration
```
mongodb:
  -
    # 主键(string).多值用逗号分隔
    key: test
    # 地址(string).多值用逗号分隔
    address: 10.11.165.44:27017
    # 数据库(string).默认数据库
    database: jx3robot
    # 用户名(string).
    username: admin
    # 密码(string).
    password: mongo@kingsoft.com
    # 授权源(string), 默认admin
    source: admin
    # 读优先(string或{RMode string; RTagSet map[string]string; RMaxStateness time.Duration}): primary | primarypreferred | secondary | secondarypreferred | nearest, 默认由server端决定
    readPreference: primary
    # 读安全(string或{Level string}): majority | local, 默认由server端决定
    readConcern: majority
    # 写安全(string或{WMajority bool; W int; J bool; WTagSet string; WTimeout time.Duration}): majority | w3 | w2 | w1 | w0, 默认由server端决定
    writeConcern: majority
    # 是否直连(bool), true | false, 默认false
    direct: false
    # 指定副本
    replicaSet:
    # keepalive(time.Duration), 默认300秒
    keepalive: 300s
    # 连接超时(time.Duration), 默认10秒
    connectTimeout: 10s
    # 服务端选择超时(time.Duration), 默认30秒
    serverSelectionTimeout: 30s
    # 读写超时(time.Duration), 默认0表示永远不超时
    socketTimeout: 0
    # 心跳间隔(time.Duration), 默认10秒
    heartbeatInterval: 10s
    # 延迟窗口(time.Duration), 默认15毫秒
    localThreshold: 15ms
    # 最小连接数(uint64), 默认0
    minPoolSize: 0
    # 最大连接数(uint64), 默认0表示无限
    maxPoolSize: 0
    # 最大连接闲置(time.Duration), 默认0表示无限
    maxConnIdleTime: 0
    # 压缩通信(string或[]string), snappy(3.4) | zlib(3.6) | zstd(4.2)
    compressors: "snappy"
    # ZLIB级别(int)
    zlibLevel:
    # ZSTD级别(int)
    zstdLevel:
    # 读重试3.6+(bool)
    retryReads:
    # 写重试3.6+(bool)
    retryWrites:

```

# Index
- const
```
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
```
配置常量

- type Config
```
type Config struct {
	// 连接URL, 格式为[mongodb://][user:pass@]host1[:port1][,host2[:port2],...][/database][?options]
	// 常用选项
	Address        []string        `json:"address" bson:"address" yaml:"address"`
	Database       string          `json:"database" bson:"database" yaml:"database"`                   // 默认DB
	Username       string          `json:"username" bson:"username" yaml:"username"`                   // 用户名
	Password       string          `json:"password" bson:"password" yaml:"password"`                   // 密码
	Source         string          `json:"source" bson:"source" yaml:"source"`                         // 授权DB, 默认为admin
	ReadPreference *ReadPreference `json:"readPreference" bson:"readPreference" yaml:"readPreference"` // 读优先级, primary|primaryPreferred|secondary|secondaryPreferred|nearest
	ReadConcern    *ReadConcern    `json:"readConcern" bson:"readConcern" yaml:"readConcern"`          // 读影响， available|local|majority|linerizable
	WriteConcern   *WriteConcern   `json:"writeConcern" bson:"writeConcern" yaml:"writeConcern"`       // 写影响

	// 连接管理
	Direct                 bool          `json:"direct" bson:"direct" yaml:"direct"`                                                 //是否直接
	ReplicaSet             string        `json:"replicaSet" bson:"replicaSet" yaml:"replicaSet"`                                     //数据集
	Keepalive              time.Duration `json:"keepalive" bson:"keepalive" yaml:"keepalive"`                                        //默认300秒
	ConnectTimeout         time.Duration `json:"connectTimeout" bson:"connectTimeout" yaml:"connectTimeout"`                         //连接超时. 默认为10秒
	ServerSelectionTimeout time.Duration `json:"serverSelectionTimeout" bson:"serverSelectionTimeout" yaml:"serverSelectionTimeout"` // 服务端选择超时, 默认为30秒
	SocketTimeout          time.Duration `json:"socketTimeout" bson:"socketTimeout" yaml:"socketTimeout"`                            //写超时, 默认为0, 表示阻塞
	HeartbeatInterval      time.Duration `json:"heartbeatInterval" bson:"heartbeatInterval" yaml:"heartbeatInterval"`                // 心跳间隔, 默认为10秒
	LocalThreshold         time.Duration `json:"localThreshold" bson:"localThreshold" yaml:"localThreshold"`                         // 延迟窗口, 默认为15毫秒

	// 连接池管理
	MinPoolSize     uint64        `json:"minPoolSize" bson:"minPoolSize" yaml:"minPoolSize"`             // 连接池最小连接数
	MaxPoolSize     uint64        `json:"maxPoolSize" bson:"maxPoolSize" yaml:"maxPoolSize"`             // 连接池最大连接数
	MaxConnIdleTime time.Duration `json:"maxConnIdleTime" bson:"maxConnIdleTime" yaml:"maxConnIdleTime"` // 连接池最大空闲时间

	// 压缩通信
	Compressors []string `json:"compressors" bson:"compressors" yaml:"compressors"` // 压缩算法, snappy(3.4), zlib(3.6), zstd(4.2)
	ZlibLevel   int      `json:"zlibLevel" bson:"zlibLevel" yaml:"zlibLevel"`       // zlib压缩等级
	ZstdLevel   int      `json:"zstdLevel" bson:"zstdLevel" yaml:"zstdLevel"`       // zstd压缩等级

	// 重试机制
	RetryReads  bool `json:"retryReads" bson:"retryReads" yaml:"retryReads"`    // 重试读(3.6)
	RetryWrites bool `json:"retryWrites" bson:"retryWrites" yaml:"retryWrites"` // 重试写(3.6)
}
```
客户端配置

- type ReadPreference
```
type ReadPreference struct {
	RMode         string            // 读模式, primary | primaryPreferred | secondary | secondaryPreferred | nearest
	RTagSet       map[string]string // 读标签, 支持 k1:v1,k2:v2,...的格式
	RMaxStateness time.Duration     // specify a maxinum replication lag for reads from secondaries in a replica set
}
```
读优先配置

- type ReadConcern
```
type ReadConcern struct {
	Level string `json:"level" bson:"level" yaml:"level"`
}
```
读安全配置

- type WriteConcern
```
type WriteConcern struct {
	J         bool          `json:"J" bson:"J" yaml:"J"`                         // write operations are written to the journal
	W         int           `json:"W" bson:"W" yaml:"W"`                         // write operations propagate to the specified number of mongod instances
	WMajority bool          `json:"WMajority" bson:"WMajority" yaml:"WMajority"` // write operations propagate to the majority of mongod instances
	WTagSet   string        `json:"WTagSet" bson:"WTagSet" yaml:"WTagSet"`       // write operations propagate to the specified mongod instance
	WTimeout  time.Duration `json:"WTimeout" bson:"WTimeout" yaml:"WTimeout"`    // specifies a time limit for the write concern
}
```
写安全配置

- type Client
```
type Client struct {
	*mongo.Client
	DB string
}
```

- func NewClient
```$xslt
func NewClient(opt *Config) (ret *Client, err error) {
```
根据配置创建新客户端

- func Setup
```
func Setup(key string, cnf *Config) (err error) 
```
初始安装客户端,然后使用Get()或Must()获取

- func Get
```
func Get(key string) *Client 
```
返回指定的客户端, 结果可能为空

- func Must
```
func Must(key string) *Client 
```
返回指定的客户端, 结果不能为空, 否则panic!