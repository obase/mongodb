package distinct

import (
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func Options(fn func(opts *options.DistinctOptions)) *options.DistinctOptions {
	opts := options.Distinct()
	fn(opts)
	return opts
}

func Collation(c *options.Collation) *options.DistinctOptions {
	return options.Distinct().SetCollation(c)
}

func Locale(locale string) *options.DistinctOptions {
	return options.Distinct().SetCollation(&options.Collation{
		Locale: locale,
	})
}

func SimpleLocale(locale string) *options.DistinctOptions {
	return options.Distinct().SetCollation(&options.Collation{
		Locale: "simple",
	})
}

func ChineseLocale(locale string) *options.DistinctOptions {
	return options.Distinct().SetCollation(&options.Collation{
		Locale: "zh",
	})
}

func MaxTime(d time.Duration) *options.DistinctOptions {
	return options.Distinct().SetMaxTime(d)
}
