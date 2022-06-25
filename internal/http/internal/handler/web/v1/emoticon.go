package v1

import (
	"bytes"
	"fmt"
	"time"

	"go-chat/internal/cache"
	"go-chat/internal/entity"
	"go-chat/internal/http/internal/dto/apiv"
	api2 "go-chat/internal/http/internal/dto/web"
	"go-chat/internal/model"
	"go-chat/internal/pkg/ichat"
	"go-chat/internal/pkg/jwtutil"
	"go-chat/internal/pkg/utils"

	"go-chat/internal/pkg/filesystem"
	"go-chat/internal/pkg/sliceutil"
	"go-chat/internal/pkg/strutil"
	"go-chat/internal/service"
)

type Emoticon struct {
	fileSystem *filesystem.Filesystem
	service    *service.EmoticonService
	redisLock  *cache.RedisLock
}

func NewEmoticonHandler(
	service *service.EmoticonService,
	fileSystem *filesystem.Filesystem,
	redisLock *cache.RedisLock,
) *Emoticon {
	return &Emoticon{
		service:    service,
		fileSystem: fileSystem,
		redisLock:  redisLock,
	}
}

// CollectList 收藏列表
func (c *Emoticon) CollectList(ctx *ichat.Context) error {
	var (
		uid     = jwtutil.GetUid(ctx.Context)
		sys     = make([]*apiv.SysEmoticonResponse, 0)
		collect = make([]*apiv.EmoticonItem, 0)
	)

	if ids := c.service.Dao().GetUserInstallIds(uid); len(ids) > 0 {
		if items, err := c.service.Dao().FindByIds(ids); err == nil {
			for _, item := range items {
				data := &apiv.SysEmoticonResponse{
					EmoticonId: item.Id,
					Url:        item.Icon,
					Name:       item.Name,
					List:       make([]*apiv.EmoticonItem, 0),
				}

				if items, err := c.service.Dao().GetDetailsAll(item.Id, 0); err == nil {
					for _, item := range items {
						data.List = append(data.List, &apiv.EmoticonItem{
							MediaId: item.Id,
							Src:     item.Url,
						})
					}
				}

				sys = append(sys, data)
			}
		}
	}

	if items, err := c.service.Dao().GetDetailsAll(0, uid); err == nil {
		for _, item := range items {
			collect = append(collect, &apiv.EmoticonItem{
				MediaId: item.Id,
				Src:     item.Url,
			})
		}

	}

	return ctx.Success(entity.H{
		"sys_emoticon":     sys,
		"collect_emoticon": collect,
	})
}

// DeleteCollect 删除收藏表情包
func (c *Emoticon) DeleteCollect(ctx *ichat.Context) error {
	params := &api2.DeleteCollectRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	if err := c.service.DeleteCollect(jwtutil.GetUid(ctx.Context), sliceutil.ParseIds(params.Ids)); err != nil {
		return ctx.BusinessError(err.Error())
	}

	return ctx.Success(nil)
}

// Upload 上传自定义表情包
func (c *Emoticon) Upload(ctx *ichat.Context) error {
	file, err := ctx.Context.FormFile("emoticon")
	if err != nil {
		return ctx.InvalidParams("emoticon 字段必传！")
	}

	if !sliceutil.InStr(strutil.FileSuffix(file.Filename), []string{"png", "jpg", "jpeg", "gif"}) {
		return ctx.InvalidParams("上传文件格式不正确,仅支持 png、jpg、jpeg 和 gif")
	}

	// 判断上传文件大小（5M）
	if file.Size > 5<<20 {
		return ctx.InvalidParams("上传文件大小不能超过5M！")
	}

	stream, err := filesystem.ReadMultipartStream(file)
	if err != nil {
		return ctx.BusinessError("上传失败！")
	}

	size := utils.ReadFileImage(bytes.NewReader(stream))
	ext := strutil.FileSuffix(file.Filename)
	src := fmt.Sprintf("public/media/image/emoticon/%s/%s", time.Now().Format("20060102"), strutil.GenImageName(ext, size.Width, size.Height))
	if err = c.fileSystem.Default.Write(stream, src); err != nil {
		return ctx.BusinessError("上传失败！")
	}

	m := &model.EmoticonItem{
		UserId:     jwtutil.GetUid(ctx.Context),
		Describe:   "自定义表情包",
		Url:        c.fileSystem.Default.PublicUrl(src),
		FileSuffix: ext,
		FileSize:   int(file.Size),
	}

	if err := c.service.Db().Create(m).Error; err != nil {
		return ctx.BusinessError("上传失败！")
	}

	return ctx.Success(entity.H{
		"media_id": m.Id,
		"src":      m.Url,
	})
}

// SystemList 系统表情包列表
func (c *Emoticon) SystemList(ctx *ichat.Context) error {
	items, err := c.service.Dao().GetSystemEmoticonList()

	if err != nil {
		return ctx.BusinessError(err.Error())
	}

	ids := c.service.Dao().GetUserInstallIds(jwtutil.GetUid(ctx.Context))

	data := make([]*apiv.SysEmoticonList, 0, len(items))
	for _, item := range items {
		data = append(data, &apiv.SysEmoticonList{
			ID:     item.Id,
			Name:   item.Name,
			Icon:   item.Icon,
			Status: strutil.BoolToInt(sliceutil.InInt(item.Id, ids)), // 查询用户是否使用
		})
	}

	return ctx.Success(data)
}

// SetSystemEmoticon 添加或移除系统表情包
func (c *Emoticon) SetSystemEmoticon(ctx *ichat.Context) error {
	var (
		err    error
		params = &api2.SetSystemEmoticonRequest{}
		uid    = jwtutil.GetUid(ctx.Context)
		key    = fmt.Sprintf("sys-emoticon:%d", uid)
	)

	if err = ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	if !c.redisLock.Lock(ctx.Context, key, 5) {
		return ctx.BusinessError("请求频繁！")
	}
	defer c.redisLock.UnLock(ctx.Context, key)

	if params.Type == 2 {
		if err = c.service.RemoveUserSysEmoticon(uid, params.EmoticonId); err != nil {
			return ctx.BusinessError(err.Error())
		}

		return ctx.Success(nil)
	}

	// 查询表情包是否存在
	info, err := c.service.Dao().FindById(params.EmoticonId)
	if err != nil {
		return ctx.BusinessError(err.Error())
	}

	if err := c.service.AddUserSysEmoticon(uid, params.EmoticonId); err != nil {
		return ctx.BusinessError(err.Error())
	}

	items := make([]*apiv.EmoticonItem, 0)
	if list, err := c.service.Dao().GetDetailsAll(params.EmoticonId, 0); err == nil {
		for _, item := range list {
			items = append(items, &apiv.EmoticonItem{
				MediaId: item.Id,
				Src:     item.Url,
			})
		}
	}

	return ctx.Success(&apiv.SysEmoticonResponse{
		EmoticonId: info.Id,
		Url:        info.Icon,
		Name:       info.Name,
		List:       items,
	})
}