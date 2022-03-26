// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"context"
	"github.com/google/wire"
	"go-chat/internal/cache"
	"go-chat/internal/dao"
	"go-chat/internal/pkg/filesystem"
	"go-chat/internal/provider"
	"go-chat/internal/service"
	"go-chat/internal/websocket/internal/handler"
	"go-chat/internal/websocket/internal/process"
	"go-chat/internal/websocket/internal/process/handle"
	"go-chat/internal/websocket/internal/router"
)

// Injectors from wire.go:

func Initialize(ctx context.Context) *Providers {
	config := provider.NewConfig()
	client := provider.NewRedisClient(ctx, config)
	sidServer := cache.NewSid(client)
	wsClientSession := cache.NewWsClientSession(client, config, sidServer)
	clientService := service.NewClientService(wsClientSession)
	room := cache.NewRoom(client)
	db := provider.NewMySQLClient(config)
	baseService := service.NewBaseService(db, client)
	baseDao := dao.NewBaseDao(db, client)
	relation := cache.NewRelation(client)
	groupMemberDao := dao.NewGroupMemberDao(baseDao, relation)
	groupMemberService := service.NewGroupMemberService(baseService, groupMemberDao)
	defaultWebSocket := handler.NewDefaultWebSocket(client, config, clientService, room, groupMemberService)
	handlerHandler := &handler.Handler{
		DefaultWebSocket: defaultWebSocket,
	}
	session := cache.NewSession(client)
	engine := router.NewRouter(config, handlerHandler, session)
	websocketServer := provider.NewWebsocketServer(engine)
	redisLock := cache.NewRedisLock(client)
	clearGarbage := process.NewClearGarbage(client, redisLock, sidServer)
	heartbeat := process.NewImHeartbeat()
	server := process.NewServer(config, sidServer)
	talkVote := cache.NewTalkVote(client)
	talkRecordsVoteDao := dao.NewTalkRecordsVoteDao(baseDao, talkVote)
	filesystemFilesystem := filesystem.NewFilesystem(config)
	talkRecordsDao := dao.NewTalkRecordsDao(baseDao)
	talkRecordsService := service.NewTalkRecordsService(baseService, talkVote, talkRecordsVoteDao, filesystemFilesystem, groupMemberDao, talkRecordsDao)
	usersFriendsDao := dao.NewUsersFriendsDao(baseDao, relation)
	contactService := service.NewContactService(baseService, usersFriendsDao)
	subscribeConsume := handle.NewSubscribeConsume(config, wsClientSession, room, talkRecordsService, contactService)
	wsSubscribe := process.NewWsSubscribe(client, config, subscribeConsume)
	processProcess := process.NewProcess(clearGarbage, heartbeat, server, wsSubscribe)
	providers := &Providers{
		Config:   config,
		WsServer: websocketServer,
		Process:  processProcess,
	}
	return providers
}

// wire.go:

var providerSet = wire.NewSet(provider.NewConfig, provider.NewMySQLClient, provider.NewRedisClient, provider.NewWebsocketServer, router.NewRouter, process.NewProcess, process.NewClearGarbage, process.NewImHeartbeat, process.NewServer, process.NewWsSubscribe, handle.NewSubscribeConsume, cache.NewSession, cache.NewSid, cache.NewRedisLock, cache.NewWsClientSession, cache.NewRoom, cache.NewTalkVote, cache.NewRelation, dao.NewBaseDao, dao.NewTalkRecordsDao, dao.NewTalkRecordsVoteDao, dao.NewGroupMemberDao, dao.NewUsersFriendsDao, filesystem.NewFilesystem, service.NewBaseService, service.NewTalkRecordsService, service.NewClientService, service.NewGroupMemberService, service.NewContactService, handler.NewDefaultWebSocket, wire.Struct(new(handler.Handler), "*"), wire.Struct(new(Providers), "*"))
