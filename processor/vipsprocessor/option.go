package vipsprocessor

import (
	"go.uber.org/zap"
	"strings"
)

type Option func(v *VipsProcessor)

func WithFilter(name string, filter FilterFunc) Option {
	return func(v *VipsProcessor) {
		v.Filters[name] = filter
	}
}

func WithDisableBlur(disabled bool) Option {
	return func(v *VipsProcessor) {
		v.DisableBlur = disabled
	}
}

func WithDisableFilters(filters ...string) Option {
	return func(v *VipsProcessor) {
		for _, raw := range filters {
			splits := strings.Split(raw, ",")
			for _, name := range splits {
				name = strings.TrimSpace(name)
				if len(name) > 0 {
					v.DisableFilters = append(v.DisableFilters, name)
				}
			}
		}
	}
}

func WithMaxFilterOps(num int) Option {
	return func(v *VipsProcessor) {
		if num != 0 {
			v.MaxFilterOps = num
		}
	}
}

func WithConcurrencyLevel(num int) Option {
	return func(v *VipsProcessor) {
		if num > 0 {
			v.ConcurrencyLevel = num
		}
	}
}

func WithLogger(logger *zap.Logger) Option {
	return func(v *VipsProcessor) {
		if logger != nil {
			v.Logger = logger
		}
	}
}

func WithDebug(debug bool) Option {
	return func(v *VipsProcessor) {
		v.Debug = debug
	}
}
