package option

import (
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Options struct {
	options.CollectionOptions
	options.FindOptions
	options.FindOneOptions
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
