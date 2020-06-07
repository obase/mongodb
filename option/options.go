package option

import (
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"time"
)

type Options struct {

	/*
	- CollectionOptions
	*/
	// The read concern to use for operations executed on the Collection. The default value is nil, which means that
	// the read concern of the database used to configure the Collection will be used.
	ReadConcern *readconcern.ReadConcern

	// The write concern to use for operations executed on the Collection. The default value is nil, which means that
	// the write concern of the database used to configure the Collection will be used.
	WriteConcern *writeconcern.WriteConcern

	// The read preference to use for operations executed on the Collection. The default value is nil, which means that
	// the read preference of the database used to configure the Collection will be used.
	ReadPreference *readpref.ReadPref

	// The BSON registry to marshal and unmarshal documents for operations executed on the Collection. The default value
	// is nil, which means that the registry of the database used to configure the Collection will be used.
	Registry *bsoncodec.Registry

	/*
	- FineOptions
	- FindOneOptions
	- CountOptions
	- EstimatedDocumentCountOptions
	*/
	// If true, an operation on a sharded cluster can return partial results if some shards are down rather than
	// returning an error. The default value is false.
	AllowPartialResults *bool

	// The maximum number of documents to be included in each batch returned by the server.
	BatchSize *int32

	// Specifies a collation to use for string comparisons during the operation. This option is only valid for MongoDB
	// versions >= 3.4. For previous server versions, the driver will return an error if this option is used. The
	// default value is nil, which means the default collation of the collection will be used.
	Collation *options.Collation

	// A string that will be included in server logs, profiling logs, and currentOp queries to help trace the operation.
	// The default is the empty string, which means that no comment will be included in the logs.
	Comment *string

	// Specifies the type of cursor that should be created for the operation. The default is NonTailable, which means
	// that the cursor will be closed by the server when the last batch of documents is retrieved.
	CursorType *options.CursorType

	// The index to use for the aggregation. This should either be the index name as a string or the index specification
	// as a document. The default value is nil, which means that no hint will be sent.
	Hint interface{}

	// The maximum number of documents to return. The default value is 0, which means that all documents matching the
	// filter will be returned. A negative limit specifies that the resulting documents should be returned in a single
	// batch. The default value is 0.
	Limit *int64

	// A document specifying the exclusive upper bound for a specific index. The default value is nil, which means that
	// there is no maximum value.
	Max interface{}

	// The maximum amount of time that the server should wait for new documents to satisfy a tailable cursor query.
	// This option is only valid for tailable await cursors (see the CursorType option for more information) and
	// MongoDB versions >= 3.2. For other cursor types or previous server versions, this option is ignored.
	MaxAwaitTime *time.Duration

	// The maximum amount of time that the query can run on the server. The default value is nil, meaning that there
	// is no time limit for query execution.
	MaxTime *time.Duration

	// A document specifying the inclusive lower bound for a specific index. The default value is 0, which means that
	// there is no minimum value.
	Min interface{}

	// If true, the cursor created by the operation will not timeout after a period of inactivity. The default value
	// is false.
	NoCursorTimeout *bool

	// This option is for internal replication use only and should not be set.
	OplogReplay *bool

	// A document describing which fields will be included in the documents returned by the operation. The default value
	// is nil, which means all fields will be included.
	Projection interface{}

	// If true, the documents returned by the operation will only contain fields corresponding to the index used. The
	// default value is false.
	ReturnKey *bool

	// If true, a $recordId field with a record identifier will be included in the documents returned by the operation.
	// The default value is false.
	ShowRecordID *bool

	// The number of documents to skip before adding documents to the result. The default value is 0.
	Skip *int64

	// If true, the cursor will not return a document more than once because of an intervening write operation. This
	// option has been deprecated in MongoDB version 4.0. The default value is false.
	Snapshot *bool

	// A document specifying the order in which documents should be returned.
	Sort interface{}
	/*
	-
	*/
}

/*--------------------------------转换其他选项------------------------------------*/
func FindOptions(opts ...*Options) []*options.FindOptions {
	if len(opts) == 0 {
		return nil
	}
	return nil
}

func FindOneOptions(opts ...*Options) []*options.FindOneOptions {
	if len(opts) == 0 {
		return nil
	}
	return nil
}
