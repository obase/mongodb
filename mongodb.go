package mongodb

import "time"

/*
客户端选项
*/

type WriteConcernOption struct {
	J         bool          `json:"J" bson:"J" yaml:"J"`
	W         int           `json:"W" bson:"W" yaml:"W"`
	WMajority bool          `json:"WMajority" bson:"WMajority" yaml:"WMajority"`
	WTagSet   string        `json:"WTagSet" bson:"WTagSet" yaml:"WTagSet"`
	WTimeout  time.Duration `json:"WTimeout" bson:"WTimeout" yaml:"WTimeout"`
}

type Option struct {
	Key                     string             `json:"key" bson:"key" yaml:"key"`
	Address                 []string           `json:"address" bson:"address" yaml:"address"`
	Source                  string             `json:"source" bson:"source" yaml:"source"`
	Username                string             `json:"username" bson:"username" yaml:"username"`
	Password                string             `json:"password" bson:"password" yaml:"password"`
	AuthMechanism           string             `json:"authMechanism" bson:"authMechanism" yaml:"authMechanism"`
	AuthMechanismProperties map[string]string  `json:"authMechanismProperties" bson:"authMechanismProperties" yaml:"authMechanismProperties"`
	Compressors             []string           `json:"compressors" bson:"compressors" yaml:"compressors"` // zstd, zlib, snappy
	ConnectTimeout          time.Duration      `json:"connectTimeout" bson:"connectTimeout" yaml:"connectTimeout"`
	Keepalive               time.Duration      `json:"keepalive" bson:"keepalive" yaml:"keepalive"`
	ConnectDirect           bool               `json:"connectDirect" bson:"connectDirect" yaml:"connectDirect"` // 是否直接连接. 如果是仅限于connect_string指定的hosts,否则会尝试与集群交互发现
	HeartbeatInterval       time.Duration      `json:"heartbeatInterval" bson:"heartbeatInterval" yaml:"heartbeatInterval"`
	LocalThreshold          time.Duration      `json:"localThreshold" bson:"localThreshold" yaml:"localThreshold"`
	MaxConnIdleTime         time.Duration      `json:"maxConnIdleTime" bson:"maxConnIdleTime" yaml:"maxConnIdleTime"`
	MaxPoolSize             uint64             `json:"maxPoolSize" bson:"maxPoolSize" yaml:"maxPoolSize"`
	MinPoolSize             uint64             `json:"minPoolSize" bson:"minPoolSize" yaml:"minPoolSize"`
	ReadConcern             string             `json:"readConcern" bson:"readConcern" yaml:"readConcern"` // local, majority
	ReadPreference          string             `json:"readPreference" bson:"readPreference" yaml:"readPreference"`
	ReplicaSet              string             `json:"replicaSet" bson:"replicaSet" yaml:"replicaSet"`
	RetryReads              bool               `json:"retryReads" bson:"retryReads" yaml:"retryReads"`
	RetryWrites             bool               `json:"retryWrites" bson:"retryWrites" yaml:"retryWrites"`
	ServerSelectionTimeout  time.Duration      `json:"serverSelectionTimeout" bson:"serverSelectionTimeout" yaml:"serverSelectionTimeout"`
	SocketTimeout           time.Duration      `json:"socketTimeout" bson:"socketTimeout" yaml:"socketTimeout"`
	WriteConcern            WriteConcernOption `json:"writeConcern" bson:"writeConcern" yaml:"writeConcern"` // local, majority
}
