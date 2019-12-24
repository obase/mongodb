package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"net"
	"strings"
	"time"
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
			auth.AuthSource = opt.Database
		}
		opts.SetAuth(auth)
	}
	if len(opt.Compressors) > 0 {
		opts.SetCompressors(opt.Compressors)
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

	socketTimeout := opt.ReadTimeout
	if socketTimeout < opt.WriteTimeout {
		socketTimeout = opt.WriteTimeout
	}
	if socketTimeout > 0 {
		opts.SetSocketTimeout(socketTimeout)
	}

	if opt.MaxPoolIdleTimeMS > 0 {
		opts.SetMaxConnIdleTime(time.Duration(opt.MaxPoolIdleTimeMS) * time.Millisecond)
	}
	if opt.MaxPoolSize > 0 {
		opts.SetMaxPoolSize(uint64(opt.MaxPoolSize))
	}
	if opt.MinPoolSize > 0 {
		opts.SetMinPoolSize(uint64(opt.MinPoolSize))
	}
	if opt.MaxPoolWaitTimeMS > 0 {
		opts.SetServerSelectionTimeout(time.Duration(opt.MaxPoolWaitTimeMS) * time.Millisecond)
	}

	if opt.Mode > 0 {
		var pref *readpref.ReadPref
		if pref, err = readpref.New(opt.Mode); err != nil {
			return
		}
		opts.SetReadPreference(pref)
	}

	if opt.Safe != nil {

		if opt.Safe.RMode != "" {
			opts.SetReadConcern(readconcern.New(readconcern.Level(strings.ToLower(opt.Safe.RMode))))
		}

		var ops []writeconcern.Option
		if opt.Safe.J {
			ops = append(ops, writeconcern.J(true))
		}
		if opt.Safe.W > 0 {
			ops = append(ops, writeconcern.W(opt.Safe.W))
		}
		if strings.ToLower(opt.Safe.WMode) == "majority" {
			ops = append(ops, writeconcern.WMajority())
		}
		if opt.Safe.WTimeout > 0 {
			ops = append(ops, writeconcern.WTimeout(time.Duration(opt.Safe.WTimeout)*time.Millisecond))
		}
		opts.SetWriteConcern(writeconcern.New(ops...))
	}

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
