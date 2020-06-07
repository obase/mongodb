package mongodb

import (
	"context"
	"github.com/obase/mongodb/option"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/tag"
	"net"
)

var emptyFilter = bson.M{}

// 组合已有客户端,直接支持相关方法, 另再提供若干方便方法
type Client struct {
	*mongo.Client
	DB                string
	collectionOptions *options.CollectionOptions
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

	opts.SetDirect(opt.Direct) // FIXBUG: direct to whole cluster

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

func (cc *Client) CollectionOptions(opts *options.CollectionOptions) {
	cc.collectionOptions = opts
}

func (cc *Client) Close() error {
	if cc.Client != nil {
		return cc.Client.Disconnect(nil)
	}
	return nil
}

func (cc *Client) ListDatabaseNames(filter interface{}) ([]string, error) {
	if filter == nil {
		filter = emptyFilter
	}
	return cc.Client.ListDatabaseNames(nil, filter)
}

func (cc *Client) ListCollectionNames(filter interface{}) ([]string, error) {
	if filter == nil {
		filter = emptyFilter
	}
	return cc.Client.Database(cc.DB).ListCollectionNames(nil, filter)
}

func (cc *Client) Collection(cl string, opts ...*options.CollectionOptions) *mongo.Collection {
	if cc.collectionOptions != nil {
		opts = append(opts, cc.collectionOptions)
	}
	return cc.Database(cc.DB).Collection(cl, opts...)
}

func (cc *Client) Count(cl string, filter interface{}, opts ...*option.Options) (ret int64, err error) {

	if filter == nil {
		return cc.Database(cc.DB).Collection(cl, cc.collectionOptions).EstimatedDocumentCount(nil)
	} else {
		return cc.Database(cc.DB).Collection(cl, cc.collectionOptions).CountDocuments(nil, filter)
	}

}

func (cc *Client) FindId(cl string, id interface{}, ret interface{}, opts ...*option.Options) (not bool, err error) {
	err = cc.Database(cc.DB).Collection(cl, cc.collectionOptions).FindOne(nil, bson.M{"_id": id}).Decode(ret)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			not = true
			err = nil
		}
	}
	return
}

func (cc *Client) FindOne(cl string, filter interface{}, ret interface{}) (not bool, err error) {
	if filter == nil {
		filter = emptyFilter
	}
	err = cc.Database(cc.DB).Collection(cl, cc.collectionOptions).FindOne(nil, filter).Decode(ret)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			not = true
			err = nil
		}
	}
	return
}

func (cc *Client) FindAll(cl string, filter interface{}, ret interface{}, opts ...*options.FindOptions) (err error) {
	if filter == nil {
		filter = emptyFilter
	}
	cur, err := cc.Database(cc.DB).Collection(cl, cc.collectionOptions).Find(nil, filter, opts...)
	if err == nil {
		err = cur.All(nil, ret)
	}
	return
}

func (cc *Client) FindWith(cl string, filter interface{}, fn func(cur *mongo.Cursor) error, opts ...*options.FindOptions) (err error) {
	if filter == nil {
		filter = emptyFilter
	}
	cur, err := cc.Database(cc.DB).Collection(cl, cc.collectionOptions).Find(nil, filter, opts...)
	if err == nil {
		defer cur.Close(nil)
		err = fn(cur)
	}
	return
}

func (cc *Client) Distinct(cl string, fieldName string, filter interface{}, opts ...*options.DistinctOptions) (ret []interface{}, err error) {
	if filter == nil {
		filter = emptyFilter
	}
	ret, err = cc.Database(cc.DB).Collection(cl, cc.collectionOptions).Distinct(nil, fieldName, filter, opts...)
	return
}

func (cc *Client) FindIdAndUpdate(cl string, id interface{}, update interface{}) {
	cc.Database(cc.DB).Collection(cl, cc.collectionOptions).FindOneAndUpdate(nil, bson.M{"_id": id}, update)
}
