package websocket

import (
	"fmt"

	"go-chat/app/pkg/im"
)

type DefaultChannelHandle struct {
}

func NewDefaultChannelHandle() *DefaultChannelHandle {
	return new(DefaultChannelHandle)
}

// Open 连接成功回调事件
func (d *DefaultChannelHandle) Open(client *im.Client) {
	fmt.Printf("[%s] 客户端已连接[%s] \n", client.Channel.Name, client.Uuid)
}

// Message 消息接收回调事件
func (d *DefaultChannelHandle) Message(message *im.RecvMessage) {
	fmt.Printf("[%s]消息通知 Client:%s ，Content: %s \n", message.Client.Channel.Name, message.Client.Uuid, message.Content)

	message.Client.Channel.SendMessage(&im.SendMessage{
		IsAll:   true,
		Clients: nil,
		Event:   "talk",
		Content: message.Content,
	})
}

// Close 客户端关闭回调事件
func (d *DefaultChannelHandle) Close(client *im.Client, code int, text string) {
	fmt.Printf("[%s] 客户端[%s] 已关闭 \n", client.Channel.Name, client.Uuid)
}