package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/tag"
	"net"
)

type mongodbClient struct {
	*mongo.Client
	DB string
}

func newMongodbClient(opt *Config) (ret *mongodbClient, err error) {

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
	ret = &mongodbClient{
		Client: client,
		DB:     opt.Database,
	}
	return
}

func (client *mongodbClient) Count(c string) (int64, error) {
	return client.Client.Database(client.DB).Collection(c).EstimatedDocumentCount(context.Background())
}

func (client *mongodbClient) CountWith(c string, filter interface{}) (n int64, err error) {
	return client.Client.Database(client.DB).Collection(c).CountDocuments(context.Background(), filter)
}

func (client *mongodbClient) Close() error {
	return client.Client.Disconnect(context.Background())
}
