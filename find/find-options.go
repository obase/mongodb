package find

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func Options(fn func(opts *options.FindOptions)) *options.FindOptions {
	opts := options.Find()
	fn(opts)
	return opts
}

func AllowPartialResults(b bool) *options.FindOptions {
	return options.Find().SetAllowPartialResults(b)
}

func BatchSize(n int32) *options.FindOptions {
	return options.Find().SetBatchSize(n)
}

func Collation(c *options.Collation) *options.FindOptions {
	return options.Find().SetCollation(c)
}

// 参考: https://docs.mongodb.com/manual/reference/collation-locales-defaults/#collation-languages-locales
func Locale(locale string) *options.FindOptions {
	return options.Find().SetCollation(&options.Collation{
		Locale: locale,
	})
}

func SimpleLocale() *options.FindOptions {
	return options.Find().SetCollation(&options.Collation{
		Locale: "simple",
	})
}

func ChineseLocale() *options.FindOptions {
	return options.Find().SetCollation(&options.Collation{
		Locale: "zh",
	})
}

func CursorType(c options.CursorType) *options.FindOptions {
	return options.Find().SetCursorType(c)
}

func Hint(hint interface{}) *options.FindOptions {
	return options.Find().SetHint(hint)
}

func Skip(skip int64) *options.FindOptions {
	return options.Find().SetSkip(skip)
}

func Limit(limit int64) *options.FindOptions {
	return options.Find().SetLimit(limit)
}

func Max(max interface{}) *options.FindOptions {
	return options.Find().SetMax(max)
}

func MaxAwaitTime(d time.Duration) *options.FindOptions {
	return options.Find().SetMaxAwaitTime(d)
}

func MaxTime(d time.Duration) *options.FindOptions {
	return options.Find().SetMaxTime(d)
}

func Min(min interface{}) *options.FindOptions {
	return options.Find().SetMin(min)
}

func NoCursorTimeout(b bool) *options.FindOptions {
	return options.Find().SetNoCursorTimeout(b)
}

func Projection(projection interface{}) *options.FindOptions {
	return options.Find().SetProjection(projection)
}

func Project(fs ...string) *options.FindOptions {
	ps := make(bson.D, len(fs))
	for i, f := range fs {
		ps[i] = bson.E{Key: f, Value: 1}
	}
	return options.Find().SetProjection(ps)
}

func Sort(sort interface{}) *options.FindOptions {
	return options.Find().SetSort(sort)
}

func Asc(name string) *options.FindOptions {
	return options.Find().SetSort(bson.M{name: 1})
}

func Desc(name string) *options.FindOptions {
	return options.Find().SetSort(bson.M{name: -1})
}

func Page(skip int64, limit int64, sorts ...string) *options.FindOptions {
	opts := options.Find().SetSkip(skip).SetLimit(limit)
	if len(sorts) > 0 {
		smap := make(bson.M, len(sorts))
		for _, s := range sorts {
			switch s[0] {
			case '-':
				smap[s[1:]] = -1
			case '+':
				smap[s[1:]] = 1
			default:
				smap[s] = 1
			}
		}
		opts.SetSort(smap)
	}
	return opts
}
