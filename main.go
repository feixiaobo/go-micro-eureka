package go_micro_eureka

import (
	"github.com/feixiaobo/go-micro-eureka/option"
	"github.com/feixiaobo/go-micro-eureka/register"
)

func main() {
	server := register.EurekaServer(
			option.RegistryAddress("http://localhost:8761/eureka"),
			option.Name("wemall"),
		)
	server.Start()
}