package register

import (
	"fmt"
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
	opts Options
}

func EurekaServer(opts ...Option) Server {
	return newServer(opts...)
}

func newServer(opts ...Option) Server {
	options := newOptions(opts...)

	options.Client = &clientWrapper{
		options.Client,
		metadata.Metadata{
			HeaderPrefix + "From-Service": options.Server.Options().Name,
		},
	}

	return &service{
		opts: options,
	}
}

func Start()  {
	go registry()
}

func register() {
	registerCenter := eureka.NewRegistry(
		registry.Addrs("http://localhost:8761/eureka"),
	)

	ip := getLocalIP()
	addr := fmt.Sprintf("%s:%d", ip, 9101)
	instanceId := fmt.Sprintf("%s:%s:%d", ip, "wemall", 9101)

	metaMap := make(map[string]string)
	metaMap["instanceId"] = instanceId

	ser := http.NewServer(
		server.Metadata(metaMap),
		server.Id(instanceId),
		server.Registry(registerCenter),
		server.Address(addr),
		server.Name("wemall"),
		server.Advertise(addr),
	)

	selector := selector.NewSelector(
		selector.Registry(registerCenter),
		selector.SetStrategy(selector.RoundRobin),
	)

	service := micro.NewService(
		micro.Name("wemall"),
		micro.Registry(registerCenter),
		micro.Server(ser),
		micro.Address(addr),
		micro.Selector(selector),
		micro.RegisterInterval(time.Second*30),
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
