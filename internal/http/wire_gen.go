// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"context"
	"github.com/google/wire"
	"go-chat/config"
	"go-chat/internal/cache"
	"go-chat/internal/dao"
	"go-chat/internal/dao/note"
	"go-chat/internal/http/internal/handler"
	"go-chat/internal/http/internal/handler/api/v1"
	"go-chat/internal/http/internal/handler/api/v1/article"
	"go-chat/internal/http/internal/handler/api/v1/contact"
	"go-chat/internal/http/internal/handler/api/v1/group"
	"go-chat/internal/http/internal/handler/api/v1/talk"
	"go-chat/internal/http/internal/router"
	"go-chat/internal/pkg/client"
	"go-chat/internal/pkg/filesystem"
	"go-chat/internal/provider"
	"go-chat/internal/service"
	note2 "go-chat/internal/service/note"
)

import (
	_ "github.com/urfave/cli/v2"
	_ "go-chat/internal/pkg/validation"
)

// Injectors from wire.go:

func Initialize(ctx context.Context, conf *config.Config) *Providers {
	redisClient := provider.NewRedisClient(ctx, conf)
	smsCodeCache := &cache.SmsCodeCache{
		Redis: redisClient,
	}
	smsService := service.NewSmsService(smsCodeCache)
	db := provider.NewMySQLClient(conf)
	baseDao := dao.NewBaseDao(db, redisClient)
	usersDao := dao.NewUserDao(baseDao)
	userService := service.NewUserService(usersDao)
	common := v1.NewCommonHandler(conf, smsService, userService)
	session := cache.NewSession(redisClient)
	redisLock := cache.NewRedisLock(redisClient)
	baseService := service.NewBaseService(db, redisClient)
	unreadTalkCache := cache.NewUnreadTalkCache(redisClient)
	lastMessage := cache.NewLastMessage(redisClient)
	talkVote := cache.NewTalkVote(redisClient)
	talkRecordsVoteDao := dao.NewTalkRecordsVoteDao(baseDao, talkVote)
	relation := cache.NewRelation(redisClient)
	groupMemberDao := dao.NewGroupMemberDao(baseDao, relation)
	sidServer := cache.NewSid(redisClient)
	wsClientSession := cache.NewWsClientSession(redisClient, conf, sidServer)
	filesystemFilesystem := filesystem.NewFilesystem(conf)
	splitUploadDao := dao.NewFileSplitUploadDao(baseDao)
	talkMessageService := service.NewTalkMessageService(baseService, conf, unreadTalkCache, lastMessage, talkRecordsVoteDao, groupMemberDao, sidServer, wsClientSession, filesystemFilesystem, splitUploadDao)
	httpClient := provider.NewHttpClient()
	clientHttpClient := client.NewHttpClient(httpClient)
	ipAddressService := service.NewIpAddressService(baseService, conf, clientHttpClient)
	talkSessionDao := dao.NewTalkSessionDao(baseDao)
	talkSessionService := service.NewTalkSessionService(baseService, talkSessionDao)
	articleClassDao := note.NewArticleClassDao(baseDao)
	articleClassService := note2.NewArticleClassService(baseService, articleClassDao)
	auth := v1.NewAuthHandler(conf, userService, smsService, session, redisLock, talkMessageService, ipAddressService, talkSessionService, articleClassService)
	user := v1.NewUserHandler(userService, smsService)
	talkService := service.NewTalkService(baseService, groupMemberDao)
	talkMessageForwardService := service.NewTalkMessageForwardService(baseService)
	splitUploadService := service.NewSplitUploadService(baseService, splitUploadDao, conf, filesystemFilesystem)
	contactDao := dao.NewContactDao(baseDao, relation)
	contactService := service.NewContactService(baseService, contactDao)
	groupMemberService := service.NewGroupMemberService(baseService, groupMemberDao)
	authPermissionService := service.NewAuthPermissionService(contactDao, groupMemberDao)
	message := talk.NewTalkMessageHandler(talkMessageService, talkService, talkRecordsVoteDao, talkMessageForwardService, splitUploadService, contactService, groupMemberService, authPermissionService)
	groupDao := dao.NewGroupDao(baseDao)
	groupService := service.NewGroupService(baseService, groupDao, groupMemberDao, relation)
	talkTalk := talk.NewTalkHandler(talkService, talkSessionService, redisLock, userService, wsClientSession, lastMessage, unreadTalkCache, contactService, groupService, authPermissionService)
	talkRecordsDao := &dao.TalkRecordsDao{
		BaseDao: baseDao,
	}
	talkRecordsService := service.NewTalkRecordsService(baseService, talkVote, talkRecordsVoteDao, filesystemFilesystem, groupMemberDao, talkRecordsDao)
	records := talk.NewTalkRecordsHandler(talkRecordsService, groupMemberService, filesystemFilesystem)
	emoticonDao := dao.NewEmoticonDao(baseDao)
	emoticonService := service.NewEmoticonService(baseService, emoticonDao, filesystemFilesystem)
	emoticon := v1.NewEmoticonHandler(emoticonService, filesystemFilesystem, redisLock)
	upload := v1.NewUploadHandler(conf, filesystemFilesystem, splitUploadService)
	groupGroup := group.NewGroupHandler(groupService, groupMemberService, talkSessionService, redisLock, contactService, userService)
	groupNoticeDao := &dao.GroupNoticeDao{
		BaseDao: baseDao,
	}
	groupNoticeService := service.NewGroupNoticeService(groupNoticeDao)
	groupNotice := group.NewGroupNoticeHandler(groupNoticeService, groupMemberService)
	contactContact := contact.NewContactHandler(contactService, wsClientSession, userService)
	contactApplyService := service.NewContactsApplyService(baseService)
	contactApply := contact.NewContactsApplyHandler(contactApplyService, userService)
	articleService := note2.NewArticleService(baseService)
	articleAnnexDao := note.NewArticleAnnexDao(baseDao)
	articleAnnexService := note2.NewArticleAnnexService(baseService, articleAnnexDao, filesystemFilesystem)
	articleArticle := article.NewArticleHandler(articleService, filesystemFilesystem, articleAnnexService)
	annex := article.NewAnnexHandler(articleAnnexService, filesystemFilesystem)
	class := article.NewClassHandler(articleClassService)
	articleTagService := note2.NewArticleTagService(baseService)
	tag := article.NewTagHandler(articleTagService)
	handlerHandler := &handler.Handler{
		Common:        common,
		Auth:          auth,
		User:          user,
		TalkMessage:   message,
		Talk:          talkTalk,
		TalkRecords:   records,
		Emoticon:      emoticon,
		Upload:        upload,
		Group:         groupGroup,
		GroupNotice:   groupNotice,
		Contact:       contactContact,
		ContactsApply: contactApply,
		Article:       articleArticle,
		ArticleAnnex:  annex,
		ArticleClass:  class,
		ArticleTag:    tag,
	}
	engine := router.NewRouter(conf, handlerHandler, session)
	httpServer := provider.NewHttpServer(conf, engine)
	providers := &Providers{
		Config: conf,
		Server: httpServer,
	}
	return providers
}

// wire.go:

var providerSet = wire.NewSet(provider.NewMySQLClient, provider.NewRedisClient, provider.NewHttpClient, provider.NewHttpServer, client.NewHttpClient, router.NewRouter, filesystem.NewFilesystem, cache.NewSession, cache.NewSid, cache.NewUnreadTalkCache, cache.NewRedisLock, cache.NewWsClientSession, cache.NewLastMessage, cache.NewTalkVote, cache.NewRoom, cache.NewRelation, wire.Struct(new(cache.SmsCodeCache), "*"), dao.NewBaseDao, dao.NewContactDao, dao.NewGroupMemberDao, dao.NewUserDao, dao.NewGroupDao, wire.Struct(new(dao.TalkRecordsDao), "*"), wire.Struct(new(dao.TalkRecordsCodeDao), "*"), wire.Struct(new(dao.TalkRecordsLoginDao), "*"), wire.Struct(new(dao.TalkRecordsFileDao), "*"), wire.Struct(new(dao.GroupNoticeDao), "*"), dao.NewTalkSessionDao, dao.NewEmoticonDao, dao.NewTalkRecordsVoteDao, dao.NewFileSplitUploadDao, note.NewArticleClassDao, note.NewArticleAnnexDao, service.NewBaseService, service.NewUserService, service.NewSmsService, service.NewTalkService, service.NewTalkMessageService, service.NewClientService, service.NewGroupService, service.NewGroupMemberService, service.NewGroupNoticeService, service.NewTalkSessionService, service.NewTalkMessageForwardService, service.NewEmoticonService, service.NewTalkRecordsService, service.NewContactService, service.NewContactsApplyService, service.NewSplitUploadService, service.NewIpAddressService, service.NewAuthPermissionService, note2.NewArticleService, note2.NewArticleTagService, note2.NewArticleClassService, note2.NewArticleAnnexService, v1.NewAuthHandler, v1.NewCommonHandler, v1.NewUserHandler, contact.NewContactHandler, contact.NewContactsApplyHandler, group.NewGroupHandler, group.NewGroupNoticeHandler, talk.NewTalkHandler, talk.NewTalkMessageHandler, v1.NewUploadHandler, v1.NewEmoticonHandler, talk.NewTalkRecordsHandler, article.NewAnnexHandler, article.NewArticleHandler, article.NewClassHandler, article.NewTagHandler, wire.Struct(new(handler.Handler), "*"), wire.Struct(new(Providers), "*"))
