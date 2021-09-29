// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package main

import (
	"context"
	"github.com/google/wire"
	"go-chat/app/cache"
	"go-chat/app/http/handler"
	"go-chat/app/http/handler/api/v1"
	"go-chat/app/http/handler/open"
	"go-chat/app/http/handler/ws"
	"go-chat/app/http/router"
	"go-chat/app/repository"
	"go-chat/app/service"
	"go-chat/config"
	"go-chat/connect"
)

import (
	_ "go-chat/app/validator"
)

// Injectors from wire.go:

func Initialize(ctx context.Context, conf *config.Config) *Service {
	client := connect.RedisConnect(ctx, conf)
	smsCodeCache := &cache.SmsCodeCache{
		Redis: client,
	}
	smsService := &service.SmsService{
		SmsCodeCache: smsCodeCache,
	}
	db := connect.MysqlConnect(conf)
	userRepository := &repository.UserRepository{
		DB: db,
	}
	common := &v1.Common{
		SmsService: smsService,
		UserRepo:   userRepository,
	}
	userService := &service.UserService{
		Repo: userRepository,
	}
	authToken := &cache.AuthTokenCache{
		Redis: client,
	}
	redisLock := &cache.RedisLock{
		Redis: client,
	}
	auth := &v1.Auth{
		Conf:           conf,
		UserService:    userService,
		SmsService:     smsService,
		AuthTokenCache: authToken,
		RedisLock:      redisLock,
	}
	user := &v1.User{
		UserRepo:   userRepository,
		SmsService: smsService,
	}
	download := &v1.Download{}
	upload := &v1.Upload{
		Conf: conf,
	}
	index := &open.Index{}
	wsClient := &cache.WsClient{
		Redis: client,
		Conf:  conf,
	}
	clientService := &service.ClientService{
		WsClient: wsClient,
	}
	webSocket := &ws.WebSocket{
		ClientService: clientService,
	}
	handlerHandler := &handler.Handler{
		Common:   common,
		Auth:     auth,
		User:     user,
		Download: download,
		Upload:   upload,
		Index:    index,
		Ws:       webSocket,
	}
	engine := router.NewRouter(conf, handlerHandler)
	server := connect.NewHttp(conf, engine)
	serverRunID := cache.NewServerRun(client)
	socketService := &service.SocketService{
		Conf:        conf,
		ServerRunID: serverRunID,
	}
	mainService := &Service{
		HttpServer:   server,
		SocketServer: socketService,
	}
	return mainService
}

// wire.go:

var providerSet = wire.NewSet(connect.RedisConnect, connect.MysqlConnect, connect.NewHttp, router.NewRouter, cache.NewServerRun, wire.Struct(new(cache.WsClient), "*"), wire.Struct(new(cache.AuthTokenCache), "*"), wire.Struct(new(cache.SmsCodeCache), "*"), wire.Struct(new(cache.RedisLock), "*"), wire.Struct(new(v1.Common), "*"), wire.Struct(new(v1.Auth), "*"), wire.Struct(new(v1.User), "*"), wire.Struct(new(v1.Upload), "*"), wire.Struct(new(v1.Download), "*"), wire.Struct(new(open.Index), "*"), wire.Struct(new(ws.WebSocket), "*"), wire.Struct(new(handler.Handler), "*"), wire.Struct(new(repository.UserRepository), "*"), wire.Struct(new(service.ClientService), "*"), wire.Struct(new(service.UserService), "*"), wire.Struct(new(service.SocketService), "*"), wire.Struct(new(service.SmsService), "*"), wire.Struct(new(Service), "*"))
