# package mongo
mongo客户端

# Installation
- go get
```
go get -u github.com/globalsign/mgo
go get -u github.com/obase/mongo
```
- go mod
```
go mod edit -require=github.com/obase/mongo@latest
```

# Configuration
```
mongo:
  -
    # 引用的key(必需)
    key: test
    # 地址(必需). 多值用逗号分隔
    address: "127.0.0.1:27017"
    # DB名字(必需)
    database: jx3robot
    # 用户名(可选)
    username:
    # 密码(可选)
    password:
    # 授权, 默认与database相同
    source:
    # 模式(可选).primary | primaryPreferred | secondary | secondaryPreferred | nearest | eventual | monotonic | strong, 默认为strong
    model: "Strong"
    # 安全(可选).默认值{"W":0, "WMode":"majority", "RMode":"", "WTimeout":0, "FSync":false, "J":false}
    safe: {"W":0, "WMode":"majority", "RMode":"majority", "WTimeout":0, "FSync":false, "J":false}
    # 连接超时(可选). 默认10秒
    connectTimeout: "10s"
    # Keepalive设置(可选). 默认无
    keepalive:
    # 读超时(可选). 默认无
    readTimeout:
    # 写超时(可选). 默认无
    writeTimeout:
    # 连接池最小数量(可选)
    minPoolSize: 0
    # 连接池最大数量(可选)
    maxPoolSize: 0
    # 连接池最大等待毫秒(可选). 默认0阻塞
    maxPoolWaitTimeMS: 0
    # 连接池最大空闲毫秒(可选)
    maxPoolIdleTimeMS: 0
    default: true
```

# Index
- type Bulk
```
type Bulk interface {
	Insert(docs ...interface{})
	Upsert(pairs ...interface{})
	RemoveOne(selectors ...interface{})
	RemoveAll(selectors ...interface{})
	UpdateOne(pairs ...interface{})
	UpdateAll(pairs ...interface{})
}
```
批量操作行为

- type BulkFunc
```
type BulkFunc func(bk Bulk, args ...interface{}
```
批量操作函数

- type SessionFunc
```
type SessionFunc func(se *mgo.Session, args ...interface{}) (interface{}, error)
```
会话回调函数

- type CollectionFunc
```
type CollectionFunc func(cl *mgo.Collection, args ...interface{}) (interface{}, error)
```
集合回调函数


- type Mongo interface
```
type Mongo interface {
	Count(c string) (n int, err error)
	Indexes(c string) (indexes []mgo.Index, err error)
	EnsureIndex(c string, index mgo.Index) error
	EnsureIndexKey(c string, key ...string) error
	DropIndex(c string, key ...string) error
	DropIndexName(c string, name string) error

	// For whole document
	FindOne(c string, ret interface{}, query interface{}) (bool, error)
	FindAll(c string, ret interface{}, query interface{}, sort ...string) error
	FindRange(c string, ret interface{}, query interface{}, skip uint32, limit uint32, sort ...string) error
	FindPage(c string, tot *uint32, ret interface{}, query interface{}, skip uint32, limit uint32, sort ...string) error
	FindDistinct(c string, ret interface{}, query interface{}, key string, sort ...string) error
	FindId(c string, ret interface{}, id interface{}) (bool, error)
	// Find And Select
	SelectOne(c string, ret interface{}, query interface{}, projection interface{}) (bool, error)
	SelectAll(c string, ret interface{}, query interface{}, projection interface{}, sort ...string) error
	SelectRange(c string, ret interface{}, query interface{}, projection interface{}, skip uint32, limit uint32, sort ...string) error
	SelectPage(c string, tot *uint32, ret interface{}, query interface{}, projection interface{}, skip uint32, limit uint32, sort ...string) error
	SelectDistinct(c string, ret interface{}, query interface{}, projection interface{}, key string, sort ...string) error
	SelectId(c string, ret interface{}, id interface{}, projection interface{}) (bool, error)
	// FindAndModify
	FindAndUpdate(c string, ret interface{}, query interface{}, update interface{}) (updated int, err error)              // return old doucument
	FindAndUpsert(c string, ret interface{}, query interface{}, upsert interface{}) (upsertedId interface{}, err error)   // return old doucument
	FindAndRemove(c string, ret interface{}, query interface{}) (removed int, err error)                                  // return old doucument
	FindAndUpdateRN(c string, ret interface{}, query interface{}, update interface{}) (updated int, err error)            // return new doucument
	FindAndUpsertRN(c string, ret interface{}, query interface{}, upsert interface{}) (upsertedId interface{}, err error) // return new doucument

	Insert(c string, docs ...interface{}) error
	RemoveOne(c string, selector interface{}) (bool, error)
	RemoveAll(c string, selector interface{}) (removed int, err error)
	RemoveId(c string, id interface{}) (bool, error)
	UpdateOne(c string, selector interface{}, update interface{}) (bool, error)
	UpdateAll(c string, selector interface{}, update interface{}) (updated int, err error)
	UpdateId(c string, id interface{}, update interface{}) (bool, error)
	UpsertOne(c string, selector interface{}, update interface{}) (upsertedId interface{}, err error)
	UpsertId(c string, id interface{}, update interface{}) (upsertedId interface{}, err error)
	RunBulk(c string, f BulkFunc, args ...interface{}) (matched int, modified int, err error)
	RunCollection(c string, f CollectionFunc, args ...interface{}) (interface{}, error)

	DBCount(d string, c string) (n int, err error)
	DBIndexes(d string, c string) (indexes []mgo.Index, err error)
	DBEnsureIndex(d string, c string, index mgo.Index) error
	DBEnsureIndexKey(d string, c string, key ...string) error
	DBDropIndex(d string, c string, key ...string) error
	DBDropIndexName(d string, c string, name string) error

	// For whole document
	DBFindOne(d string, c string, ret interface{}, query interface{}) (bool, error)
	DBFindAll(d string, c string, ret interface{}, query interface{}, sort ...string) error
	DBFindRange(d string, c string, ret interface{}, query interface{}, skip uint32, limit uint32, sort ...string) error
	DBFindPage(d string, c string, tot *uint32, ret interface{}, query interface{}, skip uint32, limit uint32, sort ...string) error
	DBFindDistinct(d string, c string, ret interface{}, query interface{}, key string, sort ...string) error
	DBFindId(d string, c string, ret interface{}, id interface{}) (bool, error)
	// Find And Select
	DBSelectOne(d string, c string, ret interface{}, query interface{}, projection interface{}) (bool, error)
	DBSelectAll(d string, c string, ret interface{}, query interface{}, projection interface{}, sort ...string) error
	DBSelectRange(d string, c string, ret interface{}, query interface{}, projection interface{}, skip uint32, limit uint32, sort ...string) error
	DBSelectPage(d string, c string, tot *uint32, ret interface{}, query interface{}, projection interface{}, skip uint32, limit uint32, sort ...string) error
	DBSelectDistinct(d string, c string, ret interface{}, query interface{}, projection interface{}, key string, sort ...string) error
	DBSelectId(d string, c string, ret interface{}, id interface{}, projection interface{}) (bool, error)
	// FindAndModify
	DBFindAndUpdate(d string, c string, ret interface{}, query interface{}, update interface{}) (updated int, err error)              // return old doucument
	DBFindAndUpsert(d string, c string, ret interface{}, query interface{}, upsert interface{}) (upsertedId interface{}, err error)   // return old doucument
	DBFindAndRemove(d string, c string, ret interface{}, query interface{}) (removed int, err error)                                  // return old doucument
	DBFindAndUpdateRN(d string, c string, ret interface{}, query interface{}, update interface{}) (updated int, err error)            // return new doucument
	DBFindAndUpsertRN(d string, c string, ret interface{}, query interface{}, upsert interface{}) (upsertedId interface{}, err error) // return new doucument

	DBInsert(d string, c string, docs ...interface{}) error
	DBRemoveOne(d string, c string, selector interface{}) (bool, error)
	DBRemoveAll(d string, c string, selector interface{}) (removed int, err error)
	DBRemoveId(d string, c string, id interface{}) (bool, error)
	DBUpdateOne(d string, c string, selector interface{}, update interface{}) (bool, error)
	DBUpdateAll(d string, c string, selector interface{}, update interface{}) (updated int, err error)
	DBUpdateId(d string, c string, id interface{}, update interface{}) (bool, error)
	DBUpsertOne(d string, c string, selector interface{}, update interface{}) (upsertedId interface{}, err error)
	DBUpsertId(d string, c string, id interface{}, update interface{}) (upsertedId interface{}, err error)
	DBRunBulk(d string, c string, f BulkFunc, args ...interface{}) (matched int, modified int, err error)
	DBRunCollection(d string, c string, f CollectionFunc, args ...interface{}) (interface{}, error)

	RunSession(f SessionFunc, args ...interface{}) (interface{}, error)
}
```
Mongo行为抽象接口

- func Get
```
func Get(name string) Mongo
```

获取配置中特定名称的实例
# Examples
```
func FindOperationPage(ctx context.Context, pg *Query) (data *OperationPage, err error) {

	data = new(OperationPage)

	query := make(map[string]interface{})
	// 操作类型
	if pg.Type != "" && pg.Type != "全部" {
		query["type"] = pg.Type
	}
	// 是否批量操作
	if pg.Batch != Batch_Nil {
		if pg.Batch == Batch_Yes {
			query["batch"] = true
		} else {
			query["batch"] = false
		}
	}
	// 开始时间~结束时间
	if pg.BeginTime != 0 || pg.EndTime != 0 {
		ctime := make(map[string]interface{})
		if pg.BeginTime != 0 {
			ctime["$gte"] = pg.BeginTime
		}
		if pg.EndTime != 0 {
			ctime["$lte"] = pg.EndTime
		}
		query["ctime"] = bson.M(ctime)
	}

	var or []bson.M
	// 操作人
	if pg.Operator != "" {
		or = append(or, bson.M{"username": bson.M{"$regex": pg.Operator}}, bson.M{"operator": bson.M{"$regex": pg.Operator}})
	}
	// 关键词
	if ln := len(pg.Keywords); ln > 0 {
		for _, kw := range pg.Keywords {
			if len(kw) > 0 {
				or = append(or, bson.M{
					"target": bson.M{"$regex": kw},
				})
			}
		}
	}

	if len(or) > 0 {
		query["$or"] = or
	}

	err = mdb.FindPage(AuditCollection, &data.Total, &data.Rows, query, uint32(pg.PageIndex*pg.PageSize), uint32(pg.PageSize), "-ctime")
	return
}

```