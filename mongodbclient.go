package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/tag"
	"net"
)

// 组合已有客户端,直接支持相关方法, 另再提供若干方便方法
type Client struct {
	*mongo.Client
	DB string
}

func newClient(opt *Config) (ret *Client, err error) {

	opts := options.Client()
	opts.SetHosts(opt.Address)
	if opt.Username != "" {
		var auth options.Credential
		auth.Username = opt.Username
		auth.Password = opt.Password
		auth.PasswordSet = true
		if opt.Source != "" {
			auth.AuthSource = opt.Source
		} else {
			auth.AuthSource = "admin"
		}
		opts.SetAuth(auth)
	}

	if opt.ReadPreference != nil {

		var rpopts []readpref.Option
		if size := len(opt.ReadPreference.RTagSet); size > 0 {
			var set tag.Set
			for k, v := range opt.ReadPreference.RTagSet {
				set = append(set, tag.Tag{
					Name:  k,
					Value: v,
				})
			}
			rpopts = append(rpopts, readpref.WithTagSets(set))
		}
		if opt.ReadPreference.RMaxStateness > 0 {
			rpopts = append(rpopts, readpref.WithMaxStaleness(opt.ReadPreference.RMaxStateness))
		}

		switch opt.ReadPreference.RMode {
		case ReadPreference_primary:
			opts.SetReadPreference(readpref.Primary())
		case ReadPreference_primaryPreferred:
			opts.SetReadPreference(readpref.PrimaryPreferred(rpopts...))
		case ReadPreference_secondary:
			opts.SetReadPreference(readpref.Secondary(rpopts...))
		case ReadPreference_secondaryPreferred:
			opts.SetReadPreference(readpref.SecondaryPreferred(rpopts...))
		case ReadPreference_nearest:
			opts.SetReadPreference(readpref.Nearest(rpopts...))
		}
	}

	if opt.ReadConcern != nil {
		switch opt.ReadConcern.Level {
		case ReadConcern_available:
			opts.SetReadConcern(readconcern.Available())
		case ReadConcern_local:
			opts.SetReadConcern(readconcern.Local())
		case ReadConcern_majority:
			opts.SetReadConcern(readconcern.Majority())
		case ReadConcern_linerizable:
			opts.SetReadConcern(readconcern.Linearizable())
		}
	}

	if opt.WriteConcern != nil {
		var wcopts []writeconcern.Option
		if opt.WriteConcern.WMajority {
			wcopts = append(wcopts, writeconcern.WMajority())
		} else {
			wcopts = append(wcopts, writeconcern.W(opt.WriteConcern.W))
		}

		if opt.WriteConcern.J {
			wcopts = append(wcopts, writeconcern.J(true))
		}
		if opt.WriteConcern.WTagSet != "" {
			wcopts = append(wcopts, writeconcern.WTagSet(opt.WriteConcern.WTagSet))
		}
		if opt.WriteConcern.WTimeout > 0 {
			wcopts = append(wcopts, writeconcern.WTimeout(opt.WriteConcern.WTimeout))
		}
		opts.SetWriteConcern(writeconcern.New(wcopts...))
	}

	opts.SetDirect(true)

	if opt.ReplicaSet != "" {
		opts.SetReplicaSet(opt.ReplicaSet)
	}

	if opt.Keepalive > 0 {
		var dialer net.Dialer
		dialer.KeepAlive = opt.Keepalive
		if opt.ConnectTimeout > 0 {
			dialer.Timeout = opt.ConnectTimeout
		}
		opts.SetDialer(&dialer)
	} else if opt.ConnectTimeout > 0 {
		opts.SetConnectTimeout(opt.ConnectTimeout)
	}

	if opt.ServerSelectionTimeout > 0 {
		opts.SetServerSelectionTimeout(opt.ServerSelectionTimeout)
	}

	if opt.SocketTimeout > 0 {
		opts.SetSocketTimeout(opt.SocketTimeout)
	}

	if opt.HeartbeatInterval > 0 {
		opts.SetHeartbeatInterval(opt.HeartbeatInterval)
	}

	if opt.LocalThreshold > 0 {
		opts.SetLocalThreshold(opt.LocalThreshold)
	}

	if opt.MinPoolSize > 0 {
		opts.SetMinPoolSize(opt.MinPoolSize)
	}
	if opt.MaxPoolSize > 0 {
		opts.SetMaxPoolSize(opt.MaxPoolSize)
	}
	if opt.MaxConnIdleTime > 0 {
		opts.SetMaxConnIdleTime(opt.MaxConnIdleTime)
	}

	if len(opt.Compressors) > 0 {
		opts.SetCompressors(opt.Compressors)
	}

	if opt.ZlibLevel > 0 {
		opts.SetZlibLevel(opt.ZlibLevel)
	}

	if opt.ZstdLevel > 0 {
		opts.SetZstdLevel(opt.ZstdLevel)
	}

	opts.SetRetryReads(opt.RetryReads)
	opts.SetRetryWrites(opt.RetryWrites)

	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		return
	}
	ret = &Client{
		Client: client,
		DB:     opt.Database,
	}
	return
}

func (cc *Client) Close() error {
	if cc.Client != nil {
		return cc.Client.Disconnect(nil)
	}
	return nil
}

func (cc *Client) Collection(cl string, opts ...*options.CollectionOptions) *mongo.Collection {
	return cc.Database(cc.DB).Collection(cl, opts...)
}

func (cc *Client) DBCollection(db string, cl string, opts ...*options.CollectionOptions) *mongo.Collection {
	return cc.Database(db).Collection(cl, opts...)
}

func (cc *Client) Count(cl string, filters ...interface{}) (ret int64, err error) {

	if len(filters) == 0 {
		return cc.Database(cc.DB).Collection(cl).EstimatedDocumentCount(nil)
	} else {
		return cc.Database(cc.DB).Collection(cl).CountDocuments(nil, filters[0])
	}

}

func (cc *Client) DBCount(db string, cl string, filters ...interface{}) (ret int64, err error) {

	if len(filters) == 0 {
		return cc.Database(db).Collection(cl).EstimatedDocumentCount(nil)
	} else {
		return cc.Database(db).Collection(cl).CountDocuments(nil, filters[0])
	}

}

func (cc *Client) FindId(ret interface{}, cl string, id interface{}) (not bool, err error) {
	err = cc.Database(cc.DB).Collection(cl).FindOne(nil, bson.M{"_id": id}).Decode(ret)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			not = true
			err = nil
		}
	}
	return
}

func (cc *Client) DBFindId(ret interface{}, db string, cl string, id interface{}) (not bool, err error) {
	err = cc.Database(db).Collection(cl).FindOne(nil, bson.M{"_id": id}).Decode(ret)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			not = true
			err = nil
		}
	}
	return
}

func (cc *Client) FindOne(ret interface{}, cl string, filter interface{}) (not bool, err error) {
	err = cc.Database(cc.DB).Collection(cl).FindOne(nil, filter).Decode(ret)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			not = true
			err = nil
		}
	}
	return
}

func (cc *Client) DBFindOne(ret interface{}, db string, cl string, filter interface{}) (not bool, err error) {
	err = cc.Database(db).Collection(cl).FindOne(nil, filter).Decode(ret)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			not = true
			err = nil
		}
	}
	return
}

func (cc *Client) FindAll(ret interface{}, cl string, filter interface{}, sorts ...bson.E) (err error) {
	var cur *mongo.Cursor
	if len(sorts) > 0 {
		cur, err = cc.Database(cc.DB).Collection(cl).Find(nil, filter, options.Find().SetSort(bson.D(sorts)))
	} else {
		cur, err = cc.Database(cc.DB).Collection(cl).Find(nil, filter)
	}
	if err == nil {
		err = cur.All(nil, ret)
	}
	return
}

func (cc *Client) DBFindAll(ret interface{}, db string, cl string, filter interface{}, sorts ...bson.E) (err error) {
	var cur *mongo.Cursor
	if len(sorts) > 0 {
		cur, err = cc.Database(db).Collection(cl).Find(nil, filter, options.Find().SetSort(bson.D(sorts)))
	} else {
		cur, err = cc.Database(db).Collection(cl).Find(nil, filter)
	}
	if err == nil {
		err = cur.All(nil, ret)
	}
	return
}

func (cc *Client) FindWith(fn func(*mongo.Cursor) error, cl string, filter interface{}, sorts ...bson.E) (err error) {
	var cur *mongo.Cursor
	if len(sorts) > 0 {
		cur, err = cc.Database(cc.DB).Collection(cl).Find(nil, filter, options.Find().SetSort(bson.D(sorts)))
	} else {
		cur, err = cc.Database(cc.DB).Collection(cl).Find(nil, filter)
	}
	if err == nil {
		defer cur.Close(nil)
		err = fn(cur)
	}
	return
}

func (cc *Client) DBFindWith(fn func(*mongo.Cursor) error, db string, cl string, filter interface{}, sorts ...bson.E) (err error) {
	var cur *mongo.Cursor
	if len(sorts) > 0 {
		cur, err = cc.Database(db).Collection(cl).Find(nil, filter, options.Find().SetSort(bson.D(sorts)))
	} else {
		cur, err = cc.Database(db).Collection(cl).Find(nil, filter)
	}
	if err == nil {
		defer cur.Close(nil)
		err = fn(cur)
	}
	return
}

func (cc *Client) FindRange(ret interface{}, cl string, filter interface{}, skip int64, limit int64, sorts ...bson.E) (err error) {

	var cur *mongo.Cursor
	if len(sorts) > 0 {
		cc.Database(cc.DB).Collection(cl).Find(nil, filter, options.Find().SetSkip(skip).SetLimit(limit).SetSort(bson.D(sorts)))
	} else {
		cc.Database(cc.DB).Collection(cl).Find(nil, filter, options.Find().SetSkip(skip).SetLimit(limit))
	}
	if err == nil {
		err = cur.All(nil, ret)
	}
	return
}
