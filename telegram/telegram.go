package telegram

type User struct {
	Id                      int64  `json:"id" bson:"id"`
	IsBot                   bool   `json:"is_bot" bson:"is_bot"`
	FirstName               string `json:"first_name" bson:"first_name"`
	LastName                string `json:"last_name" bson:"last_name"`
	Username                string `json:"username" bson:"username"`
	LanguageCode            string `json:"language_code" bson:"language_code"`
	IsPremium               bool   `json:"is_premium" bson:"is_premium"`
	AddedToAttachmentMenu   bool   `json:"added_to_attachment_menu" bson:"added_to_attachment_menu"`
	CanJoinGroups           bool   `json:"can_join_groups" bson:"can_join_groups"`
	CanReadAllGroupMessages bool   `json:"can_read_all_group_messages" bson:"can_read_all_group_messages"`
	SupportsInlineQueries   bool   `json:"supports_inline_queries" bson:"supports_inline_queries"`
}

type ChatPhoto struct {
	SmallFileId       string `json:"small_file_id"`
	SmallFileUniqueId string `json:"small_file_unique_id"`
	BigFileId         string `json:"big_file_id"`
	BigFileUniqueId   string `json:"big_file_unique_id"`
}

type File struct {
	FileId       string `json:"file_id"`
	FileUniqueId string `json:"file_unique_id"`
	FileSize     string `json:"file_size"`
	FilePath     string `json:"file_path"`
}

type MaskPosition struct {
	Point  string `json:"point"`
	XShift string `json:"x_shift"`
	YShift string `json:"y_shift"`
	Scale  string `json:"scale"`
}

type PhotoSize struct {
	FileId       string `json:"file_id"`
	FileUniqueId string `json:"file_unique_id"`
	Width        int64  `json:"width"`
	Height       int64  `json:"height"`
	FileSize     int64  `json:"file_size"`
}

type Sticker struct {
	FileId           string       `json:"file_id"`
	FileUniqueId     string       `json:"file_unique_id"`
	Type             string       `json:"type"`
	Width            int64        `json:"width"`
	Height           int64        `json:"height"`
	IsAnimated       bool         `json:"is_animated"`
	IsVideo          bool         `json:"is_video"`
	Thumb            PhotoSize    `json:"thumb"`
	Emoji            string       `json:"emoji"`
	SetName          string       `json:"set_name"`
	PremiumAnimation File         `json:"premium_animation"`
	MaskPosition     MaskPosition `json:"mask_position"`
	CustomEmojiId    string       `json:"custom_emoji_id"`
	FileSize         int64        `json:"file_size"`
}

type ChatPermissions struct {
	CanSendMessages       bool `json:"can_send_messages"`
	CanSendMediaMessages  bool `json:"can_send_media_messages"`
	CanSendPolls          bool `json:"can_send_polls"`
	CanSendOtherMessages  bool `json:"can_send_other_messages"`
	CanAddWebPagePreviews bool `json:"can_add_web_page_previews"`
	CanChangeInfo         bool `json:"can_change_info"`
	CanInviteUsers        bool `json:"can_invite_users"`
	CanPinMessages        bool `json:"can_pin_messages"`
	CanManageTopics       bool `json:"can_manage_topics"`
}

type Location struct {
	Longitude            float64 `json:"longitude"`
	Latitude             float64 `json:"latitude"`
	HorizontalAccuracy   float64 `json:"horizontal_accuracy"`
	LivePeriod           int64   `json:"live_period"`
	Heading              int64   `json:"heading"`
	ProximityAlertRadius int64   `json:"proximity_alert_radius"`
}

type ChatLocation struct {
	Location Location `json:"location"`
	Address  string   `json:"address"`
}

type Chat struct {
	Id                                 int64           `json:"id"`
	Type                               string          `json:"type"`
	Title                              string          `json:"title"`
	Username                           string          `json:"username"`
	FirstName                          string          `json:"first_name"`
	LastName                           string          `json:"last_name"`
	IsForum                            bool            `json:"is_forum"`
	Photo                              ChatPhoto       `json:"photo"`
	ActiveUsernames                    []string        `json:"active_usernames"`
	EmojiStatusCustomEmojiId           string          `json:"emoji_status_custom_emoji_id"`
	Bio                                string          `json:"bio"`
	HasPrivateForwards                 bool            `json:"has_private_forwards"`
	HasRestrictedVoiceAndVideoMessages bool            `json:"has_restricted_voice_and_video_messages"`
	JoinToSendMessages                 bool            `json:"join_to_send_messages"`
	JoinByRequest                      bool            `json:"join_by_request"`
	Description                        string          `json:"description"`
	InviteLink                         string          `json:"invite_link"`
	PinnedMessage                      *Message        `json:"pinned_message"`
	Permissions                        ChatPermissions `json:"permissions"`
	SlowModeDelay                      int64           `json:"slow_mode_delay"`
	MessageAutoDeleteTime              int64           `json:"message_auto_delete_time"`
	HasProtectedContent                bool            `json:"has_protected_content"`
	StickerSetName                     string          `json:"sticker_set_name"`
	CanSetStickerSet                   bool            `json:"can_set_sticker_set"`
	LinkedChatId                       int64           `json:"linked_chat_id"`
	Location                           ChatLocation    `json:"location"`
}

type WebAppData struct {
	ButtonText string `json:"button_text"`
	Data       string `json:"data"`
}

type Message struct {
	MessageId            int64           `json:"message_id"`
	MessageThreadId      int64           `json:"message_thread_id"`
	From                 User            `json:"from" bson:"from"`
	SenderChat           Chat            `json:"sender_chat"`
	Date                 int64           `json:"date"`
	Chat                 Chat            `json:"chat"`
	ForwardFrom          User            `json:"forward_from"`
	ForwardFromChat      Chat            `json:"forward_from_chat"`
	ForwardFromMessageId int64           `json:"forward_from_message_id"`
	ForwardSignature     string          `json:"forward_signature"`
	ForwardSenderName    string          `json:"forward_sender_name"`
	ForwardDate          int64           `json:"forward_date"`
	IsTopicMessage       bool            `json:"is_topic_message"`
	IsAutomaticForward   bool            `json:"is_automatic_forward"`
	ReplyToMessage       *Message        `json:"reply_to_message"`
	ViaBot               User            `json:"via_bot"`
	EditDate             int64           `json:"edit_date"`
	Sticker              Sticker         `json:"sticker"`
	Text                 string          `json:"text"`
	Entities             []MessageEntity `json:"entities"`
	NewChatMembers       []User          `json:"new_chat_members"`
	LeftChatMember       User            `json:"left_chat_member"`
	WebAppData           WebAppData      `json:"web_app_data"`
}

type CallbackQuery struct {
	Id              string  `json:"id"`
	From            User    `json:"from"`
	Message         Message `json:"message"`
	InlineMessageId string  `json:"inline_message_id"`
	ChatInstance    string  `json:"chat_instance"`
	Data            string  `json:"data"`
	GameShortName   string  `json:"game_short_name"`
}

type Update struct {
	UpdateId          int64         `json:"update_id"`
	Message           Message       `json:"message"`
	EditedMessage     Message       `json:"edited_message"`
	ChannelPost       Message       `json:"channel_post"`
	EditedChannelPost Message       `json:"edited_channel_post"`
	CallbackQuery     CallbackQuery `json:"callback_query"`
}

type MessageEntity struct {
	Type          string `json:"type"`
	Offset        int64  `json:"offset"`
	Length        int64  `json:"length"`
	Url           int64  `json:"url"`
	User          User   `json:"user"`
	Language      int64  `json:"language"`
	CustomEmojiId string `json:"custom_emoji_id"`
}

type WebAppInfo struct {
	Url string `json:"url,omitempty"`
}

type LoginUrl struct {
	Url                string `json:"url"`
	ForwardText        string `json:"forward_text"`
	BotUsername        string `json:"bot_username"`
	RequestWriteAccess bool   `json:"request_write_access"`
}
type InlineKeyboardButton struct {
	Text                         string      `json:"text"`
	Url                          string      `json:"url,omitempty"`
	CallbackData                 string      `json:"callback_data,omitempty"`
	WebApp                       *WebAppInfo `json:"web_app,omitempty"`
	LoginUrl                     *LoginUrl   `json:"login_url,omitempty"`
	SwitchInlineQuery            string      `json:"switch_inline_query,omitempty"`
	SwitchInlineQueryCurrentChat string      `json:"switch_inline_query_current_chat,omitempty"`
	Pay                          bool        `json:"pay,omitempty"`
}
type KeyboardButtonPollType struct {
	Type string `json:"type,omitempty"`
}
type KeyboardButton struct {
	Text            string                  `json:"text"`
	RequestContact  bool                    `json:"request_contact,omitempty"`
	RequestLocation bool                    `json:"request_location,omitempty"`
	RequestPoll     *KeyboardButtonPollType `json:"request_poll,omitempty"`
	WebApp          *WebAppInfo             `json:"web_app,omitempty"`
}
type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard,omitempty"`
}
type ReplyKeyboardMarkup struct {
	Keyboard              [][]KeyboardButton `json:"keyboard"`
	IsPersistent          bool               `json:"is_persistent"`
	ResizeKeyboard        bool               `json:"resize_keyboard"`
	OneTimeKeyboard       bool               `json:"one_time_keyboard"`
	InputFieldPlaceholder string             `json:"input_field_placeholder"`
	Selective             bool               `json:"selective"`
}
type replyMarkup interface {
	InlineKeyboardMarkup | ReplyKeyboardMarkup
}
type chatId interface {
	int64 | int | string
}
type SendMessageRequestBody[T replyMarkup, Q chatId] struct {
	ChatId                   Q               `json:"chat_id"`
	MessageThreadId          int64           `json:"message_thread_id"`
	Text                     string          `json:"text"`
	ParseMode                string          `json:"parse_mode"`
	Entities                 []MessageEntity `json:"entities"`
	DisableWebPagePreview    bool            `json:"disable_web_page_preview"`
	DisableNotification      bool            `json:"disable_notification"`
	ProtectContent           bool            `json:"protect_content"`
	ReplyToMessageId         int64           `json:"reply_to_message_id"`
	AllowSendingWithoutReply bool            `json:"allow_sending_without_reply"`
	ReplyMarkup              T               `json:"reply_markup,omitempty"`
}
type DeleteMessageRequestBody struct {
	ChatId    int64 `json:"chat_id"`
	MessageId int64 `json:"message_id"`
}
type EditMessageTextRequestBody struct {
	ChatId                int64                `json:"chat_id"`
	MessageId             int64                `json:"message_id"`
	InlineMessageId       string               `json:"inline_message_id"`
	Text                  string               `json:"text"`
	ParseMode             string               `json:"parse_mode"`
	Entities              MessageEntity        `json:"entities"`
	DisableWebPagePreview bool                 `json:"disable_web_page_preview"`
	ReplyMarkup           InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}
type EditMessageReplyMarkupRequestBody struct {
	ChatId          int64                `json:"chat_id"`
	MessageId       int64                `json:"message_id"`
	InlineMessageId string               `json:"inline_message_id"`
	ReplyMarkup     InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}
type AnswerCallbackQueryRequestBody struct {
	CallbackQueryId string `json:"callback_query_id"`
	Text            string `json:"text"`
	ShowAlert       bool   `json:"show_alert"`
}

type ButtonTelegrmBot interface {
	InlineKeyboardButton | KeyboardButton
}
type ButtonBot[T ButtonTelegrmBot] struct {
	Row    int
	Col    int
	Button T
}
