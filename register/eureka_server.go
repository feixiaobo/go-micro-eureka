package register

import (
	"context"
	"fmt"
	"github.com/feixiaobo/go-micro-eureka/option"
	"github.com/feixiaobo/go-plugins/registry/eureka"
	"github.com/feixiaobo/go-plugins/server/http"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client/selector"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/server"
	"net"
	"time"
)

type Server struct {
	opts option.Options
}

func EurekaServer(opts ...option.Option) Server {
	return newServer(opts...)
}

func newServer(opts ...option.Option) Server {
	options := newOptions(opts...)

	return Server{
		opts: options,
	}
}

func newOptions(opts ...option.Option) option.Options {
	opt := option.Options{
		Context:   context.Background(),
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

func (s *Server) Start() {
	go register(s)
}

func register(s *Server) {
	opts := s.opts

	if len(opts.RegistryAddress) == 0 {
		panic("the register address is required")
	}
	registerCenter := eureka.NewRegistry(
		registry.Addrs(opts.RegistryAddress...),
	)

	name := opts.Name
	if name == "" {
		panic("the server name is required")
	}
	ip := getLocalIP()
	port := opts.Port
	if port == 0 {
		panic("the server port is required")
	}
	ttl := opts.RegisterTTL
	if ttl == time.Duration(0) {
		ttl = time.Second*30
	}

	addr := fmt.Sprintf("%s:%d", ip, port)
	instanceId := fmt.Sprintf("%s:%s:%d", ip, name, port)

	metaMap := opts.Metadata
	metaMap["instanceId"] = instanceId

	ser := http.NewServer(
		server.Metadata(metaMap),
		server.Id(instanceId),
		server.Registry(registerCenter),
		server.Address(addr),
		server.Name(name),
		server.Advertise(addr),
	)

	selector := selector.NewSelector(
		selector.Registry(registerCenter),
		selector.SetStrategy(selector.RoundRobin),
	)

	service := micro.NewService(
		micro.Name(name),
		micro.Registry(registerCenter),
		micro.Server(ser),
		micro.Address(addr),
		micro.Selector(selector),
		micro.RegisterInterval(ttl),
	)

	service.Init()
	service.Run()
}

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	panic("Unable to determine local IP address (non loopback). Exiting.")
}
