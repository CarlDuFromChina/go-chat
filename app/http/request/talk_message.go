package request

type TextMessageRequest struct {
	TalkType   int    `form:"talk_type" json:"talk_type" binding:"required,oneof=1 2" label:"talk_type"`
	ReceiverId int    `form:"receiver_id" json:"receiver_id" binding:"required,numeric" label:"receiver_id"`
	Text       string `form:"text" json:"text" binding:"required,len:65535" label:"text"`
}

type CodeMessageRequest struct {
	TalkType   int    `form:"talk_type" json:"talk_type" binding:"required,oneof=1 2" label:"talk_type"`
	ReceiverId int    `form:"receiver_id" json:"receiver_id" binding:"required,numeric" label:"receiver_id"`
	Lang       string `form:"lang" json:"lang" binding:"required,len:65535"`
	Code       string `form:"code" json:"code" binding:"required,len:65535"`
}

type ImageMessageRequest struct {
	TalkType   int    `form:"talk_type" json:"talk_type" binding:"required,oneof=1 2" label:"talk_type"`
	ReceiverId int    `form:"receiver_id" json:"receiver_id" binding:"required,numeric" label:"receiver_id"`
	Image      string `form:"image" json:"image" binding:"required,file"`
}

type FileMessageRequest struct {
	TalkType   int    `form:"talk_type" json:"talk_type" binding:"required,oneof=1 2" label:"talk_type"`
	ReceiverId int    `form:"receiver_id" json:"receiver_id" binding:"required,numeric" label:"receiver_id"`
	HashName   string `form:"hash_name" json:"hash_name" binding:"required"`
}

type VoteMessageRequest struct {
	ReceiverId int    `form:"receiver_id" json:"receiver_id" binding:"required,numeric" label:"receiver_id"`
	Mode       string `form:"hash_name" json:"hash_name" binding:"required,,oneof=0 1"`
	Title      string `form:"title" json:"title" binding:"required"`
	Options    string `form:"options" json:"options" binding:"required"`
}

type EmoticonMessageRequest struct {
	TalkType   int    `form:"talk_type" json:"talk_type" binding:"required,oneof=1 2" label:"talk_type"`
	ReceiverId int    `form:"receiver_id" json:"receiver_id" binding:"required,numeric" label:"receiver_id"`
	EmoticonId string `form:"emoticon_id" json:"emoticon_id" binding:"required,numeric"`
}

type ForwardMessageRequest struct {
	TalkType        int `form:"talk_type" json:"talk_type" binding:"required,oneof=1 2" label:"talk_type"`
	ReceiverId      int `form:"receiver_id" json:"receiver_id" binding:"required,numeric" label:"receiver_id"`
	ForwardMode     int `form:"forward_mode" json:"forward_mode" binding:"required,oneof=1 2"`
	RecordsIds      int `form:"records_ids" json:"records_ids" binding:"required,ids"`
	ReceiveUserIds  int `form:"receive_user_ids" json:"receive_user_ids" binding:"required,ids"`
	ReceiveGroupIds int `form:"receive_group_ids" json:"receive_group_ids" binding:"required,ids"`
}

type CardMessageRequest struct {
	TalkType   int    `form:"talk_type" json:"talk_type" binding:"required,oneof=1 2" label:"talk_type"`
	ReceiverId int    `form:"receiver_id" json:"receiver_id" binding:"required,numeric" label:"receiver_id"`
	EmoticonId string `form:"emoticon_id" json:"emoticon_id" binding:"required,numeric"`
}