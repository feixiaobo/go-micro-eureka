package client

import (
	"context"
	http2 "github.com/feixiaobo/go-plugins/client/http"
	"github.com/google/martian/log"
	"github.com/micro/go-micro/client"
	client2 "github.com/micro/go-micro/client"
	"github.com/micro/go-micro/client/selector"
	"github.com/micro/go-micro/registry"
)

var httpClient client.Client

func Call(serviceName, path string, req interface{}, opts ...client.CallOption) (res interface{}, err error) {
	request := httpClient.NewRequest(serviceName, path, req)
	err = httpClient.Call(context.Background(), request, res, opts...)
	if err != nil {
		log.Errorf("call server error", err)
		return nil, err
	}
	return res, err
}

func InitClient(register *registry.Registry, s *selector.Selector, retries int) *client.Client {
	httpClient = http2.NewClient(
		client2.Retries(retries),
		client2.Registry(*register),
		client2.ContentType("application/json"),
		client2.Selector(*s),
	)
	return &httpClient
}
