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

var ALL = bson.M{}

// 组合已有客户端,直接支持相关方法, 另再提供若干方便方法
type Client struct {
	*mongo.Client
	DB                string
	collectionOptions *options.CollectionOptions
	ALL               bson.M
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
		ALL:    ALL, // 快捷引用
	}
	return
}

func (cc *Client) CollectionOptions(opts *options.CollectionOptions) *Client {
	cc.collectionOptions = opts
	return cc
}

func (cc *Client) Close() (err error) {
	if cc.Client != nil {
		err = cc.Client.Disconnect(nil)
	}
	return
}

func (cc *Client) ListDatabaseNames(filter interface{}) ([]string, error) {
	return cc.Client.ListDatabaseNames(nil, filter)
}

func (cc *Client) ListCollectionNames(filter interface{}) ([]string, error) {
	return cc.Client.Database(cc.DB).ListCollectionNames(nil, filter)
}

func (cc *Client) Collection(cl string, opts ...*options.CollectionOptions) *mongo.Collection {
	if cc.collectionOptions != nil {
		opts = append(opts, cc.collectionOptions)
	}
	return cc.Database(cc.DB).Collection(cl, opts...)
}

func (cc *Client) Count(cl string, filters ...interface{}) (ret int64, err error) {
	if len(filters) == 0 {
		return cc.Database(cc.DB).Collection(cl, cc.collectionOptions).EstimatedDocumentCount(nil)
	} else {
		return cc.Database(cc.DB).Collection(cl, cc.collectionOptions).CountDocuments(nil, filters[0])
	}

}

func (cc *Client) FindId(cl string, id interface{}, ret interface{}, opts ...*options.FindOneOptions) (not bool, err error) {
	err = cc.Database(cc.DB).Collection(cl, cc.collectionOptions).FindOne(nil, bson.M{"_id": id}, opts...).Decode(ret)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			not = true
			err = nil
		}
	}
	return
}

func (cc *Client) FindOne(cl string, filter interface{}, ret interface{}, opts ...*options.FindOneOptions) (not bool, err error) {
	err = cc.Database(cc.DB).Collection(cl, cc.collectionOptions).FindOne(nil, filter, opts...).Decode(ret)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			not = true
			err = nil
		}
	}
	return
}

func (cc *Client) Find(cl string, filter interface{}, ret interface{}, opts ...*options.FindOptions) (err error) {
	cur, err := cc.Database(cc.DB).Collection(cl, cc.collectionOptions).Find(nil, filter, opts...)
	if err == nil {
		err = cur.All(nil, ret)
	}
	return
}

func (cc *Client) FindWith(cl string, filter interface{}, with func(cur *mongo.Cursor) error, opts ...*options.FindOptions) (err error) {
	cur, err := cc.Database(cc.DB).Collection(cl, cc.collectionOptions).Find(nil, filter, opts...)
	if err == nil {
		defer cur.Close(nil)
		err = with(cur)
	}
	return
}

func (cc *Client) Distinct(cl string, fieldName string, filter interface{}, opts ...*options.DistinctOptions) (ret []interface{}, err error) {
	ret, err = cc.Database(cc.DB).Collection(cl, cc.collectionOptions).Distinct(nil, fieldName, filter, opts...)
	return
}

func (cc *Client) FindIdAndUpdate(cl string, id interface{}, update interface{}, ret interface{}, opts ...*options.FindOneAndUpdateOptions) (not bool, err error) {
	result := cc.Database(cc.DB).Collection(cl, cc.collectionOptions).FindOneAndUpdate(nil, bson.M{"_id": id}, update, opts...)
	if ret != nil {
		err = result.Decode(&ret)
		if err == mongo.ErrNoDocuments {
			not = true
			err = nil
		}
	}
	return
}

func (cc *Client) FindIdAndReplace(cl string, id interface{}, replace interface{}, ret interface{}, opts ...*options.FindOneAndReplaceOptions) (not bool, err error) {
	result := cc.Database(cc.DB).Collection(cl, cc.collectionOptions).FindOneAndReplace(nil, bson.M{"_id": id}, replace, opts...)
	if ret != nil {
		err = result.Decode(&ret)
		if err == mongo.ErrNoDocuments {
			not = true
			err = nil
		}
	}
	return
}

func (cc *Client) FindIdAndDelete(cl string, id interface{}, ret interface{}, opts ...*options.FindOneAndDeleteOptions) (not bool, err error) {
	result := cc.Database(cc.DB).Collection(cl, cc.collectionOptions).FindOneAndDelete(nil, bson.M{"_id": id}, opts...)
	if ret != nil {
		err = result.Decode(&ret)
		if err == mongo.ErrNoDocuments {
			not = true
			err = nil
		}
	}
	return
}

func (cc *Client) FindOneAndUpdate(cl string, filter interface{}, update interface{}, ret interface{}, opts ...*options.FindOneAndUpdateOptions) (not bool, err error) {
	result := cc.Database(cc.DB).Collection(cl, cc.collectionOptions).FindOneAndUpdate(nil, filter, update, opts...)
	if ret != nil {
		err = result.Decode(&ret)
		if err == mongo.ErrNoDocuments {
			not = true
			err = nil
		}
	}
	return
}

func (cc *Client) FindOneAndReplace(cl string, filter interface{}, replace interface{}, ret interface{}, opts ...*options.FindOneAndReplaceOptions) (not bool, err error) {
	result := cc.Database(cc.DB).Collection(cl, cc.collectionOptions).FindOneAndReplace(nil, filter, replace, opts...)
	if ret != nil {
		err = result.Decode(&ret)
		if err == mongo.ErrNoDocuments {
			not = true
			err = nil
		}
	}
	return
}

func (cc *Client) FindOneAndDelete(cl string, filter interface{}, ret interface{}, opts ...*options.FindOneAndDeleteOptions) (not bool, err error) {
	result := cc.Database(cc.DB).Collection(cl, cc.collectionOptions).FindOneAndDelete(nil, filter, opts...)
	if ret != nil {
		err = result.Decode(&ret)
		if err == mongo.ErrNoDocuments {
			not = true
			err = nil
		}
	}
	return
}

func (cc *Client) InsertOne(cl string, doc interface{}, opts ...*options.InsertOneOptions) (result *mongo.InsertOneResult, err error) {
	result, err = cc.Database(cc.DB).Collection(cl, cc.collectionOptions).InsertOne(nil, doc, opts...)
	if err != nil {
		return
	}
	return
}

func (cc *Client) InsertMany(cl string, docs []interface{}, opts ...*options.InsertManyOptions) (result *mongo.InsertManyResult, err error) {
	result, err = cc.Database(cc.DB).Collection(cl, cc.collectionOptions).InsertMany(nil, docs, opts...)
	if err != nil {
		return
	}
	return
}

func (cc *Client) ReplaceId(cl string, id interface{}, replace interface{}, opts ...*options.ReplaceOptions) (result *mongo.UpdateResult, err error) {
	result, err = cc.Database(cc.DB).Collection(cl, cc.collectionOptions).ReplaceOne(nil, bson.M{"_id": id}, replace, opts...)
	return
}

func (cc *Client) ReplaceOne(cl string, filter interface{}, replace interface{}, opts ...*options.ReplaceOptions) (result *mongo.UpdateResult, err error) {
	result, err = cc.Database(cc.DB).Collection(cl, cc.collectionOptions).ReplaceOne(nil, filter, replace, opts...)
	return
}

func (cc *Client) UpdateId(cl string, id interface{}, update interface{}, opts ...*options.UpdateOptions) (result *mongo.UpdateResult, err error) {
	result, err = cc.Database(cc.DB).Collection(cl, cc.collectionOptions).UpdateOne(nil, bson.M{"_id": id}, update, opts...)
	return
}

func (cc *Client) UpdateOne(cl string, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (result *mongo.UpdateResult, err error) {
	result, err = cc.Database(cc.DB).Collection(cl, cc.collectionOptions).UpdateOne(nil, filter, update, opts...)
	return
}

func (cc *Client) UpdateMany(cl string, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (result *mongo.UpdateResult, err error) {
	result, err = cc.Database(cc.DB).Collection(cl, cc.collectionOptions).UpdateMany(nil, filter, update, opts...)
	return
}

func (cc *Client) DeleteId(cl string, id interface{}, update interface{}, opts ...*options.DeleteOptions) (result *mongo.DeleteResult, err error) {
	result, err = cc.Database(cc.DB).Collection(cl, cc.collectionOptions).DeleteOne(nil, bson.M{"_id": id}, opts...)
	return
}

func (cc *Client) DeleteOne(cl string, filter interface{}, update interface{}, opts ...*options.DeleteOptions) (result *mongo.DeleteResult, err error) {
	result, err = cc.Database(cc.DB).Collection(cl, cc.collectionOptions).DeleteOne(nil, filter, opts...)
	return
}

// 必须注意: empty filter会删除整个集合数据
func (cc *Client) DeleteMany(cl string, filter interface{}, update interface{}, opts ...*options.DeleteOptions) (result *mongo.DeleteResult, err error) {
	result, err = cc.Database(cc.DB).Collection(cl, cc.collectionOptions).DeleteMany(nil, filter, opts...)
	return
}

func (cc *Client) Aggregate(cl string, pipeline interface{}, ret interface{}, opts ...*options.AggregateOptions) (err error) {
	cur, err := cc.Database(cc.DB).Collection(cl, cc.collectionOptions).Aggregate(nil, pipeline, opts...)
	if err == nil {
		err = cur.All(nil, ret)
	}
	return
}

func (cc *Client) AggregateWith(cl string, pipeline interface{}, with func(cur *mongo.Cursor), opts ...*options.AggregateOptions) (err error) {
	cur, err := cc.Database(cc.DB).Collection(cl, cc.collectionOptions).Aggregate(nil, pipeline, opts...)
	if err == nil {
		defer cur.Close(nil)
		with(cur)
	}
	return
}

func (cc *Client) BulkWrite(cl string, models []mongo.WriteModel, opts ...*options.BulkWriteOptions) (result *mongo.BulkWriteResult, err error) {
	if len(models) == 0 {
		return
	}
	result, err = cc.Database(cc.DB).Collection(cl, cc.collectionOptions).BulkWrite(nil, models, opts...)
	return
}
