package option

import (
	"context"
	"time"
)

type Options struct {
	RegistryAddress []string
	Name            string
	Port            int
	RegisterTTL     time.Duration
	Metadata        map[string]string
	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

type Option func(option *Options)

// Context specifies a context for the service.
// Can be used to signal shutdown of the service.
// Can be used for extra option values.
func Context(ctx context.Context) Option {
	return func(o *Options) {
		o.Context = ctx
	}
}

// eureka地址
func RegistryAddress(addr... string) Option {
	return func(o *Options) {
		if len(addr) == 0 {
			o.RegistryAddress = addr
		}
	}
}

// 实例名
func Name(name string) Option {
	return func(o *Options) {
		if name != "" {
			o.Name = name
		}
	}
}

// 实例端口
func Port(port int) Option {
	return func(o *Options) {
		if port > 0 {
			o.Port = port
		}
	}
}

// 实例续约时间与心跳时间
func RegisterTTL(ttl time.Duration) Option {
	return func(o *Options) {
		if ttl > time.Duration(0) {
			o.RegisterTTL = ttl
		}
	}
}

// meta info
func Metadata(metaData map[string]string) Option {
	return func(o *Options) {
		o.Metadata = metaData
	}
}