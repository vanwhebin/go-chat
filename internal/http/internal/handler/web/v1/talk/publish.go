package talk

import (
	"github.com/gin-gonic/gin/binding"
	"go-chat/api/pb/message/v1"
	"go-chat/internal/pkg/ichat"
	"go-chat/internal/service"
)

type Publish struct {
	mapping map[string]func(ctx *ichat.Context) error

	talkAuthService *service.TalkAuthService
	messageService  *service.MessageService
}

func NewPublish(talkAuthService *service.TalkAuthService, messageService *service.MessageService) *Publish {
	return &Publish{talkAuthService: talkAuthService, messageService: messageService}
}

type PublishBaseMessageRequest struct {
	Type     string    `json:"type" binding:"required"`
	Receiver *Receiver `json:"receiver" binding:"required"`
}

// Receiver 接受者信息
type Receiver struct {
	TalkType   int `json:"talk_type" binding:"required,gt=0"`   // 对话类型 1:私聊 2:群聊
	ReceiverId int `json:"receiver_id" binding:"required,gt=0"` // 好友ID或群ID
}

// Publish 发送消息接口
func (c *Publish) Publish(ctx *ichat.Context) error {

	params := &PublishBaseMessageRequest{}
	if err := ctx.Context.ShouldBindBodyWith(params, binding.JSON); err != nil {
		return ctx.InvalidParams(err)
	}

	// 权限验证
	if err := c.talkAuthService.IsAuth(ctx.Ctx(), &service.TalkAuthOpt{
		TalkType:   params.Receiver.TalkType,
		UserId:     ctx.UserId(),
		ReceiverId: params.Receiver.ReceiverId,
	}); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return c.transfer(ctx, params.Type)
}

// 文本消息
func (c *Publish) onSendText(ctx *ichat.Context) error {

	params := &message.TextMessageRequest{}
	if err := ctx.Context.ShouldBindBodyWith(params, binding.JSON); err != nil {
		return ctx.InvalidParams(err)
	}

	err := c.messageService.SendText(ctx.Ctx(), ctx.UserId(), params)
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(nil)
}

// nolint 图片消息
func (c *Publish) onSendImage(ctx *ichat.Context) error {

	params := &message.ImageMessageRequest{}
	if err := ctx.Context.ShouldBindBodyWith(params, binding.JSON); err != nil {
		return ctx.InvalidParams(err)
	}

	err := c.messageService.SendImage(ctx.Ctx(), ctx.UserId(), params)
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(nil)
}

// nolint 文件消息
func (c *Publish) onSendFile(ctx *ichat.Context) error {

	params := &message.FileMessageRequest{}
	if err := ctx.Context.ShouldBindBodyWith(params, binding.JSON); err != nil {
		return ctx.InvalidParams(err)
	}

	err := c.messageService.SendFile(ctx.Ctx(), ctx.UserId(), params)
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(nil)
}

// 代码消息
func (c *Publish) onSendCode(ctx *ichat.Context) error {

	params := &message.CodeMessageRequest{}
	if err := ctx.Context.ShouldBindBodyWith(params, binding.JSON); err != nil {
		return ctx.InvalidParams(err)
	}

	err := c.messageService.SendCode(ctx.Ctx(), ctx.UserId(), params)
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(nil)
}

// 位置消息
func (c *Publish) onSendLocation(ctx *ichat.Context) error {

	params := &message.LocationMessageRequest{}
	if err := ctx.Context.ShouldBindBodyWith(params, binding.JSON); err != nil {
		return ctx.InvalidParams(err)
	}

	err := c.messageService.SendLocation(ctx.Ctx(), ctx.UserId(), params)
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(nil)
}

// 转发消息
func (c *Publish) onSendForward(ctx *ichat.Context) error {

	params := &message.ForwardMessageRequest{}
	if err := ctx.Context.ShouldBindBodyWith(params, binding.JSON); err != nil {
		return ctx.InvalidParams(err)
	}

	err := c.messageService.SendForward(ctx.Ctx(), ctx.UserId(), params)
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(nil)
}

// 表情消息
func (c *Publish) onSendEmoticon(ctx *ichat.Context) error {

	params := &message.EmoticonMessageRequest{}
	if err := ctx.Context.ShouldBindBodyWith(params, binding.JSON); err != nil {
		return ctx.InvalidParams(err)
	}

	err := c.messageService.SendEmoticon(ctx.Ctx(), ctx.UserId(), params)
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(nil)
}

// 投票消息
func (c *Publish) onSendVote(ctx *ichat.Context) error {

	params := &message.VoteMessageRequest{}
	if err := ctx.Context.ShouldBindBodyWith(params, binding.JSON); err != nil {
		return ctx.InvalidParams(err)
	}

	if len(params.Options) <= 1 {
		return ctx.InvalidParams("options 选项必须大于1！")
	}

	if len(params.Options) > 6 {
		return ctx.InvalidParams("options 选项不能超过6个！")
	}

	err := c.messageService.SendVote(ctx.Ctx(), ctx.UserId(), params)
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(nil)
}

// 名片消息
func (c *Publish) onSendCard(ctx *ichat.Context) error {

	params := &message.LocationMessageRequest{}
	if err := ctx.Context.ShouldBindBodyWith(params, binding.JSON); err != nil {
		return ctx.InvalidParams(err)
	}

	err := c.messageService.SendLocation(ctx.Ctx(), ctx.UserId(), params)
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(nil)
}

func (c *Publish) transfer(ctx *ichat.Context, typeValue string) error {

	if c.mapping == nil {
		publishMapping := make(map[string]func(ctx *ichat.Context) error)
		publishMapping["text"] = c.onSendText
		publishMapping["code"] = c.onSendCode
		publishMapping["location"] = c.onSendLocation
		publishMapping["emoticon"] = c.onSendEmoticon
		publishMapping["vote"] = c.onSendVote
		publishMapping["image"] = c.onSendImage
		publishMapping["file"] = c.onSendFile
		publishMapping["card"] = c.onSendCard
		publishMapping["forward"] = c.onSendForward
		c.mapping = publishMapping
	}

	if call, ok := c.mapping[typeValue]; ok {
		return call(ctx)
	}

	return nil
}