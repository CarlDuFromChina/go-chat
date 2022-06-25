package v1

import (
	"strconv"
	"time"

	"go-chat/internal/http/internal/dto/web"
	"go-chat/internal/pkg/ichat"
	"go-chat/internal/pkg/jwtutil"
	"go-chat/internal/service/note"

	"go-chat/config"
	"go-chat/internal/cache"
	"go-chat/internal/entity"
	"go-chat/internal/service"
)

type Auth struct {
	config             *config.Config
	userService        *service.UserService
	smsService         *service.SmsService
	session            *cache.Session
	redisLock          *cache.RedisLock
	talkMessageService *service.TalkMessageService
	ipAddressService   *service.IpAddressService
	talkSessionService *service.TalkSessionService
	noteClassService   *note.ArticleClassService
}

func NewAuthHandler(
	config *config.Config,
	userService *service.UserService,
	smsService *service.SmsService,
	session *cache.Session,
	redisLock *cache.RedisLock,
	talkMessageService *service.TalkMessageService,
	ipAddressService *service.IpAddressService,
	talkSessionService *service.TalkSessionService,
	noteClassService *note.ArticleClassService,
) *Auth {
	return &Auth{
		config:             config,
		userService:        userService,
		smsService:         smsService,
		session:            session,
		redisLock:          redisLock,
		talkMessageService: talkMessageService,
		ipAddressService:   ipAddressService,
		talkSessionService: talkSessionService,
		noteClassService:   noteClassService,
	}
}

// Login 登录接口
func (c *Auth) Login(ctx *ichat.Context) error {

	params := &web.AuthLoginRequest{}
	if err := ctx.Context.ShouldBindJSON(params); err != nil {
		return ctx.InvalidParams(err)
	}

	user, err := c.userService.Login(params.Mobile, params.Password)
	if err != nil {
		return ctx.BusinessError(err.Error())
	}

	ip := ctx.Context.ClientIP()

	address, _ := c.ipAddressService.FindAddress(ip)

	_, _ = c.talkSessionService.Create(ctx.Context.Request.Context(), &service.TalkSessionCreateOpts{
		UserId:     user.Id,
		TalkType:   entity.ChatPrivateMode,
		ReceiverId: 4257,
		IsBoot:     true,
	})

	// 推送登录消息
	_ = c.talkMessageService.SendLoginMessage(ctx.Context.Request.Context(), &service.LoginMessageOpts{
		UserId:   user.Id,
		Ip:       ip,
		Address:  address,
		Platform: params.Platform,
		Agent:    ctx.Context.GetHeader("user-agent"),
	})

	return ctx.Success(&web.AuthLoginResponse{
		Type:        "Bearer",
		AccessToken: c.token(user.Id),
		ExpiresIn:   int(c.config.Jwt.ExpiresTime),
	})
}

// Register 注册接口
func (c *Auth) Register(ctx *ichat.Context) error {

	params := &web.RegisterRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	// 验证短信验证码是否正确
	if !c.smsService.CheckSmsCode(ctx.Context.Request.Context(), entity.SmsRegisterChannel, params.Mobile, params.SmsCode) {
		return ctx.InvalidParams("短信验证码填写错误！")
	}

	_, err := c.userService.Register(&service.UserRegisterOpts{
		Nickname: params.Nickname,
		Mobile:   params.Mobile,
		Password: params.Password,
		Platform: params.Platform,
	})
	if err != nil {
		return ctx.BusinessError(err.Error())
	}

	c.smsService.DeleteSmsCode(ctx.Context.Request.Context(), entity.SmsRegisterChannel, params.Mobile)

	return ctx.Success(nil)
}

// Logout 退出登录接口
func (c *Auth) Logout(ctx *ichat.Context) error {

	c.toBlackList(ctx)

	return ctx.Success(nil)
}

// Refresh Token 刷新接口
func (c *Auth) Refresh(ctx *ichat.Context) error {

	c.toBlackList(ctx)

	return ctx.Success(&web.AuthRefreshResponse{
		Type:        "Bearer",
		AccessToken: c.token(jwtutil.GetUid(ctx.Context)),
		ExpiresIn:   int(c.config.Jwt.ExpiresTime),
	})
}

// Forget 账号找回接口
func (c *Auth) Forget(ctx *ichat.Context) error {

	params := &web.ForgetRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	// 验证短信验证码是否正确
	if !c.smsService.CheckSmsCode(ctx.Context.Request.Context(), entity.SmsForgetAccountChannel, params.Mobile, params.SmsCode) {
		return ctx.InvalidParams("短信验证码填写错误！")
	}

	if _, err := c.userService.Forget(&service.UserForgetOpts{
		Mobile:   params.Mobile,
		Password: params.Password,
		SmsCode:  params.SmsCode,
	}); err != nil {
		return ctx.BusinessError(err.Error())
	}

	c.smsService.DeleteSmsCode(ctx.Context.Request.Context(), entity.SmsForgetAccountChannel, params.Mobile)

	return ctx.Success(nil)
}

func (c *Auth) token(uid int) string {

	expiresAt := time.Now().Add(time.Second * time.Duration(c.config.Jwt.ExpiresTime)).Unix()

	// 生成登录凭证
	token := jwtutil.GenerateToken("api", c.config.Jwt.Secret, &jwtutil.Options{
		ExpiresAt: expiresAt,
		Id:        strconv.Itoa(uid),
	})

	return token
}

// 设置黑名单
func (c *Auth) toBlackList(ctx *ichat.Context) {
	info := ctx.Context.GetStringMapString("jwt")

	expiresAt, _ := strconv.Atoi(info["expires_at"])

	ex := expiresAt - int(time.Now().Unix())

	// 将 session 加入黑名单
	_ = c.session.SetBlackList(ctx.Context.Request.Context(), info["session"], ex)
}