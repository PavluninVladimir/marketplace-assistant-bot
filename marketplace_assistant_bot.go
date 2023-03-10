package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	clientMongo      *mongo.Client
	urlOzon          string
	urlTelegramBot   string
	tokenTelegramBot string
)

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
	Text                         string     `json:"text"`
	Url                          string     `json:"url"`
	CallbackData                 string     `json:"callback_data"`
	WebApp                       WebAppInfo `json:"web_app"`
	LoginUrl                     LoginUrl   `json:"login_url"`
	SwitchInlineQuery            string     `json:"switch_inline_query"`
	SwitchInlineQueryCurrentChat string     `json:"switch_inline_query_current_chat"`
	Pay                          bool       `json:"pay"`
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
type FilterFbo struct {
	Since  string `json:"since"`
	Status string `json:"status"`
	To     string `json:"to"`
}
type WithFbo struct {
	AnalyticsData bool `json:"analytics_data"`
	FinancialData bool `json:"financial_data"`
}
type ListBodyRequestFBO struct {
	Dir      string    `json:"dir"`
	Filter   FilterFbo `json:"filter"`
	Limit    int64     `json:"limit"`
	Offset   int64     `json:"offset"`
	Translit bool      `json:"translit"`
	With     WithFbo   `json:"with"`
}
type ListRequestFBO struct {
	Body ListBodyRequestFBO
}
type ListResponseFBO struct {
	Result []struct {
		OrderId        int       `json:"order_id"`
		OrderNumber    string    `json:"order_number"`
		PostingNumber  string    `json:"posting_number"`
		Status         string    `json:"status"`
		CancelReasonId int       `json:"cancel_reason_id"`
		CreatedAt      time.Time `json:"created_at"`
		InProcessAt    time.Time `json:"in_process_at"`
		Products       []struct {
			Sku          int           `json:"sku"`
			Name         string        `json:"name"`
			Quantity     int           `json:"quantity"`
			OfferId      string        `json:"offer_id"`
			Price        string        `json:"price"`
			DigitalCodes []interface{} `json:"digital_codes"`
			CurrencyCode string        `json:"currency_code"`
		} `json:"products"`
		AnalyticsData struct {
			Region               string `json:"region"`
			City                 string `json:"city"`
			DeliveryType         string `json:"delivery_type"`
			IsPremium            bool   `json:"is_premium"`
			PaymentTypeGroupName string `json:"payment_type_group_name"`
			WarehouseId          int64  `json:"warehouse_id"`
			WarehouseName        string `json:"warehouse_name"`
			IsLegal              bool   `json:"is_legal"`
		} `json:"analytics_data"`
		FinancialData struct {
			Products []struct {
				CommissionAmount     float64     `json:"commission_amount"`
				CommissionPercent    int         `json:"commission_percent"`
				Payout               float64     `json:"payout"`
				ProductId            int         `json:"product_id"`
				CurrencyCode         string      `json:"currency_code"`
				OldPrice             int         `json:"old_price"`
				Price                int         `json:"price"`
				TotalDiscountValue   int         `json:"total_discount_value"`
				TotalDiscountPercent float64     `json:"total_discount_percent"`
				Actions              []string    `json:"actions"`
				Picking              interface{} `json:"picking"`
				Quantity             int         `json:"quantity"`
				ClientPrice          string      `json:"client_price"`
				ItemServices         struct {
					MarketplaceServiceItemFulfillment                float64 `json:"marketplace_service_item_fulfillment"`
					MarketplaceServiceItemPickup                     int     `json:"marketplace_service_item_pickup"`
					MarketplaceServiceItemDropoffPvz                 int     `json:"marketplace_service_item_dropoff_pvz"`
					MarketplaceServiceItemDropoffSc                  int     `json:"marketplace_service_item_dropoff_sc"`
					MarketplaceServiceItemDropoffFf                  int     `json:"marketplace_service_item_dropoff_ff"`
					MarketplaceServiceItemDirectFlowTrans            int     `json:"marketplace_service_item_direct_flow_trans"`
					MarketplaceServiceItemReturnFlowTrans            int     `json:"marketplace_service_item_return_flow_trans"`
					MarketplaceServiceItemDelivToCustomer            int     `json:"marketplace_service_item_deliv_to_customer"`
					MarketplaceServiceItemReturnNotDelivToCustomer   int     `json:"marketplace_service_item_return_not_deliv_to_customer"`
					MarketplaceServiceItemReturnPartGoodsCustomer    int     `json:"marketplace_service_item_return_part_goods_customer"`
					MarketplaceServiceItemReturnAfterDelivToCustomer int     `json:"marketplace_service_item_return_after_deliv_to_customer"`
				} `json:"item_services"`
			} `json:"products"`
			PostingServices struct {
				MarketplaceServiceItemFulfillment                int `json:"marketplace_service_item_fulfillment"`
				MarketplaceServiceItemPickup                     int `json:"marketplace_service_item_pickup"`
				MarketplaceServiceItemDropoffPvz                 int `json:"marketplace_service_item_dropoff_pvz"`
				MarketplaceServiceItemDropoffSc                  int `json:"marketplace_service_item_dropoff_sc"`
				MarketplaceServiceItemDropoffFf                  int `json:"marketplace_service_item_dropoff_ff"`
				MarketplaceServiceItemDirectFlowTrans            int `json:"marketplace_service_item_direct_flow_trans"`
				MarketplaceServiceItemReturnFlowTrans            int `json:"marketplace_service_item_return_flow_trans"`
				MarketplaceServiceItemDelivToCustomer            int `json:"marketplace_service_item_deliv_to_customer"`
				MarketplaceServiceItemReturnNotDelivToCustomer   int `json:"marketplace_service_item_return_not_deliv_to_customer"`
				MarketplaceServiceItemReturnPartGoodsCustomer    int `json:"marketplace_service_item_return_part_goods_customer"`
				MarketplaceServiceItemReturnAfterDelivToCustomer int `json:"marketplace_service_item_return_after_deliv_to_customer"`
			} `json:"posting_services"`
		} `json:"financial_data"`
		AdditionalData []interface{} `json:"additional_data"`
	} `json:"result"`
}
type OzonSetting struct {
	ClientId string `bson:"client_id"`
	Token    string `bson:"token"`
}
type Settings struct {
	OzonSetting OzonSetting `bson:"ozon_setting"`
}
type TelegramUser struct {
	NameBot  string   `bson:"name_bot"`
	User     User     `bson:"user"`
	Chats    []Chat   `bson:"chats"`
	Settings Settings `bson:"settings"`
}
type UserDB struct {
	Id           primitive.ObjectID `bson:"_id"`
	TelegramUser TelegramUser       `bson:"telegram_user"`
}

type Status int

const (
	AwaitingRegistration Status = iota // ?????????????? ??????????????????????,
	AcceptanceInProgress               // ???????? ??????????????,
	AwaitingApprove                    // ?????????????? ??????????????????????????,
	AwaitingPackaging                  // ?????????????? ????????????????,
	AwaitingDeliver                    // ?????????????? ????????????????,
	Arbitration                        // ????????????????,
	ClientArbitration                  // ???????????????????? ???????????????? ????????????????,
	Delivering                         // ????????????????????????,
	DriverPickup                       // ?? ????????????????,
	Delivered                          // ????????????????????,
	Cancelled                          // ????????????????,
	NotAccepted                        // ???? ???????????? ???? ?????????????????????????? ????????????,
	SentBySeller                       // ???????????????????? ??????????????????.
)

func (s Status) String() string {
	return [...]string{"awaiting_registration", "acceptance_in_progress", "awaiting_approve", "awaiting_packaging",
		"awaiting_deliver", "arbitration",
		"client_arbitration", "delivering", "driver_pickup", "delivered", "cancelled", "not_accepted",
		"sent_by_seller"}[s]
}

type CommandBot int

const (
	// SetClientIdOzonSetting ?????????????? ???????????????????? ClientId Ozon ?? ?????????????????? ??????????????
	SetClientIdOzonSetting CommandBot = iota
	GenReportToday
	GenReportYesterday
)

func (c CommandBot) String() string {
	return [...]string{
		"setclientidozonsetting",
		"???????????????????????? ?????????? ???? ??????????????",
		"???????????????????????? ?????????? ???? ??????????",
	}[c]
}

type ??onsolidatedReportFBO struct {
	TotalCount           int
	CancelledTotalCount  int
	SumCount             decimal.Decimal
	SumWithoutCommission decimal.Decimal
	products             map[string]int
	CancelledProducts    map[string]int
}

type SendMessageBot interface {
	sendMessage(body interface{}) bool
}

type AnswerCallbackQueryBot interface {
	answerCallbackQuery(body interface{}) bool
}

type EditMessageTextBot interface {
	editMessageText(body interface{}) bool
}

type DeleteMessageBot interface {
	deleteMessage()
}

func SendMessageToBot(bot SendMessageBot, body interface{}) {
	bot.sendMessage(body)
}

// answerCallbackQueryToBot ?????????????? ???? ?????????????? ???????????? ?????? ??????????????????
func answerCallbackQueryToBot(bot AnswerCallbackQueryBot, body interface{}) {
	bot.answerCallbackQuery(body)
}

func deleteMessageToBot(bot DeleteMessageBot) {
	bot.deleteMessage()
}

func editMessageTextToBot(bot EditMessageTextBot, body interface{}) {
	bot.editMessageText(body)
}

type TelegramBot struct{}

type DataCash struct {
	LastCommand string
}

var Cash map[int64]DataCash

func main() {
	Cash = make(map[int64]DataCash)
	//TelegramBot =
	urlOzon = os.Getenv("URL_OZON")
	if urlOzon == "" {
		urlOzon = "https://api-seller.ozon.ru"
		log.Printf("Defaulting to ury %s", urlOzon)
	}

	urlTelegramBot = os.Getenv("URL_TELEGRAM_BOT")
	if urlTelegramBot == "" {
		urlTelegramBot = "https://api.telegram.org/bot"
		log.Printf("Defaulting to ury %s", urlTelegramBot)
	}

	tokenTelegramBot = os.Getenv("TOKEN_TELEGRAM_BOT")
	if tokenTelegramBot == "" {
		log.Panic("Token telegram ???????? ???? ??????????????????")
	}

	mongodbUry := os.Getenv("MONGODB_URY")
	if mongodbUry == "" {
		mongodbUry = "mongodb://localhost:27017"
		log.Printf("Defaulting to ury %s", mongodbUry)
	}
	clientMongo = connectMongoDB(mongodbUry)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8181"
		log.Printf("Defaulting to port %s", port)
	}

	router := mux.NewRouter()
	router.HandleFunc("/webhooks", webHooks)
	router.HandleFunc("/", indexHandler)
	http.Handle("/", router)

	log.Printf("Listening on port %s", port)
	log.Printf("Open http://localhost:%s in the browser", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func connectMongoDB(applyUrI string) *mongo.Client {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(applyUrI))
	if err != nil {
		panic(err)
	}
	return client
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	fmt.Println("sad")
	_, err := fmt.Fprint(w, "Hello, World!")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func webHooks(w http.ResponseWriter, r *http.Request) {
	var m Update
	json.NewDecoder(r.Body).Decode(&m)
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.Encode(&m.Message.Text)
	mes := &m.Message
	bot := TelegramBot{}
	if Cash[m.Message.From.Id+m.Message.Chat.Id].LastCommand == "/setclientidozonsetting" {
		Cash[m.Message.From.Id+m.Message.Chat.Id] = DataCash{LastCommand: ""}
		coll := clientMongo.Database("MyInfantBotDB").Collection("bot_users")
		update := bson.D{{"$set", bson.D{{"telegram_user.settings.ozon_setting.client_id", m.Message.Text}}}}
		filter := bson.D{{"telegram_user.user.id", mes.From.Id}}
		opts := options.Update().SetUpsert(true)
		_, err := coll.UpdateOne(context.TODO(), filter, update, opts)
		if err != nil {
			panic(err)
		} else {
			smm := SendMessageRequestBody[InlineKeyboardMarkup, int64]{
				ChatId: m.Message.Chat.Id,
				Text:   "ClientId ?????????????? ????????????????.",
				ReplyMarkup: InlineKeyboardMarkup{
					InlineKeyboard: CreateButtonsBot[InlineKeyboardButton]([]ButtonBot[InlineKeyboardButton]{
						{Row: 1, Col: 1, Button: InlineKeyboardButton{Text: "ClientId", CallbackData: "/setclientidozonsetting"}},
						{Row: 1, Col: 2, Button: InlineKeyboardButton{Text: "Token", CallbackData: "/settokenozonsetting"}},
						{Row: 2, Col: 1, Button: InlineKeyboardButton{Text: "??????????", CallbackData: "/backsettings"}},
					})},
			}
			SendMessageToBot(&bot, smm)
		}
	}
	if Cash[m.Message.From.Id+m.Message.Chat.Id].LastCommand == "/settokenozonsetting" {
		Cash[m.Message.From.Id+m.Message.Chat.Id] = DataCash{LastCommand: ""}
		coll := clientMongo.Database("MyInfantBotDB").Collection("bot_users")
		update := bson.D{{"$set", bson.D{{"telegram_user.settings.ozon_setting.token", m.Message.Text}}}}
		filter := bson.D{{"telegram_user.user.id", mes.From.Id}}
		opts := options.Update().SetUpsert(true)
		_, err := coll.UpdateOne(context.TODO(), filter, update, opts)
		if err != nil {
			panic(err)
		} else {
			sm := TelegramBot{}
			smm := SendMessageRequestBody[InlineKeyboardMarkup, int64]{
				ChatId: m.Message.Chat.Id,
				Text:   "Token ?????????????? ????????????????.",
				ReplyMarkup: InlineKeyboardMarkup{CreateButtonsBot[InlineKeyboardButton]([]ButtonBot[InlineKeyboardButton]{
					{Row: 1, Col: 1, Button: InlineKeyboardButton{Text: "ClientId", CallbackData: "/setclientidozonsetting"}},
					{Row: 1, Col: 2, Button: InlineKeyboardButton{Text: "Token", CallbackData: "/settokenozonsetting"}},
					{Row: 2, Col: 1, Button: InlineKeyboardButton{Text: "??????????", CallbackData: "/backsettings"}},
				})},
			}
			SendMessageToBot(&sm, smm)
		}
	}
	if mes.LeftChatMember.Id != 0 {
		log.Printf("Left")
		mes.deleteMessage()
	}
	if mes.NewChatMembers != nil {
		log.Printf("New")
		mes.deleteMessage()
		//mes.sendMessage("?????????? ???????????????????? @" + mes.NewChatMembers[0].Username)
	}
	if mes.Text == "/start" {
		var user UserDB
		coll := clientMongo.Database("MyInfantBotDB").Collection("bot_users")
		filter := bson.D{{"telegram_user.user.id", mes.From.Id}}
		err := coll.FindOne(context.TODO(), filter).Decode(&user)
		if err != nil {
			fmt.Println(err)
		}
		if user.Id.IsZero() {
			userDB := UserDB{
				Id: primitive.NewObjectID(),
				TelegramUser: TelegramUser{
					User:  m.Message.From,
					Chats: []Chat{m.Message.Chat},
				}}
			_, err := coll.InsertOne(context.TODO(), userDB)
			if err != nil {
				panic(err)
			}
		} else {
			if result := findIndex[Chat](user.TelegramUser.Chats, func(c Chat) bool {
				if c.Id == m.Message.Chat.Id {
					return true
				}
				return false
			}); result < 0 {
				user.TelegramUser.Chats = append(user.TelegramUser.Chats, m.Message.Chat)
			}
			update := bson.D{{"$set", user}}
			opts := options.Update().SetUpsert(true)
			ud, err := coll.UpdateByID(context.TODO(), user.Id, update, opts)
			if err != nil {
				panic(err)
			}
			fmt.Println(ud)
		}
		sm := TelegramBot{}
		smm := SendMessageRequestBody[InlineKeyboardMarkup, int64]{
			ChatId: m.Message.Chat.Id,
			Text:   "?????????? ????????????????????! ?????????? ???????????????????????? ???????? ???????????????????? ?????? ??????????????????",
			ReplyMarkup: InlineKeyboardMarkup{CreateButtonsBot[InlineKeyboardButton]([]ButtonBot[InlineKeyboardButton]{
				{Row: 1, Col: 1, Button: InlineKeyboardButton{Text: "?????????????? ?? ?????????????????????", CallbackData: "/settings"}},
			})},
		}
		SendMessageToBot(&sm, smm)
	}
	if m.CallbackQuery.Data == "/settings" {
		answerCallbackQueryToBot(&bot, AnswerCallbackQueryRequestBody{CallbackQueryId: m.CallbackQuery.Id})
		smm := EditMessageTextRequestBody{
			ChatId:    m.CallbackQuery.Message.Chat.Id,
			MessageId: m.CallbackQuery.Message.MessageId,
			Text:      "????????????????, ???????????????????? ?????????????????????? ?????????????? ???? ???? ???????????? ??????????????????.",
			ReplyMarkup: InlineKeyboardMarkup{CreateButtonsBot[InlineKeyboardButton]([]ButtonBot[InlineKeyboardButton]{
				{Row: 1, Col: 1, Button: InlineKeyboardButton{Text: "OZON", CallbackData: "/ozonsetting"}},
				{Row: 2, Col: 1, Button: InlineKeyboardButton{Text: "??????????", CallbackData: "/backsettings"}},
			})},
		}
		editMessageTextToBot(&bot, smm)
	}
	if m.CallbackQuery.Data == "/ozonsetting" {
		answerCallbackQueryToBot(&bot, AnswerCallbackQueryRequestBody{CallbackQueryId: m.CallbackQuery.Id})
		sm := TelegramBot{}
		smm := EditMessageTextRequestBody{
			ChatId:    m.CallbackQuery.Message.Chat.Id,
			MessageId: m.CallbackQuery.Message.MessageId,
			Text:      "?????? ?????????????????? ???????????? ???? OZON seller ???????????????????? ?????????????? ClientId ?? Token. ???? ?????????? ???????????????? ?? ???????????? ???????????????? ????????????????.",
			ReplyMarkup: InlineKeyboardMarkup{CreateButtonsBot[InlineKeyboardButton]([]ButtonBot[InlineKeyboardButton]{
				{Row: 1, Col: 1, Button: InlineKeyboardButton{Text: "ClientId", CallbackData: "/setclientidozonsetting"}},
				{Row: 1, Col: 2, Button: InlineKeyboardButton{Text: "Token", CallbackData: "/settokenozonsetting"}},
				{Row: 2, Col: 1, Button: InlineKeyboardButton{Text: "???????????????? ?????????????????????? ?? Ozon Seller", CallbackData: "/testconnectozonseller"}},
				{Row: 3, Col: 1, Button: InlineKeyboardButton{Text: "??????????", CallbackData: "/backsettings"}},
			})},
		}
		editMessageTextToBot(&sm, smm)
	}
	if m.CallbackQuery.Data == "/settokenozonsetting" {
		answerCallbackQueryToBot(&bot, AnswerCallbackQueryRequestBody{CallbackQueryId: m.CallbackQuery.Id})
		Cash[m.CallbackQuery.From.Id+m.CallbackQuery.Message.Chat.Id] = DataCash{LastCommand: "/settokenozonsetting"}
		sm := TelegramBot{}
		smm := SendMessageRequestBody[InlineKeyboardMarkup, int64]{
			ChatId: m.CallbackQuery.Message.Chat.Id,
			Text:   "????. ????????????????, ???????????????????? Token ?????? ????????.",
		}
		SendMessageToBot(&sm, smm)
	}
	if m.CallbackQuery.Data == "/setclientidozonsetting" {
		answerCallbackQueryToBot(&bot, AnswerCallbackQueryRequestBody{CallbackQueryId: m.CallbackQuery.Id})
		Cash[m.CallbackQuery.From.Id+m.CallbackQuery.Message.Chat.Id] = DataCash{LastCommand: "/setclientidozonsetting"}
		sm := TelegramBot{}
		smm := SendMessageRequestBody[InlineKeyboardMarkup, int64]{
			ChatId: m.CallbackQuery.Message.Chat.Id,
			Text:   "????. ????????????????, ???????????????????? ClientID ?????? ????????.",
		}
		SendMessageToBot(&sm, smm)
	}
	if m.CallbackQuery.Data == "/backsettings" {
		answerCallbackQueryToBot(&bot, AnswerCallbackQueryRequestBody{CallbackQueryId: m.CallbackQuery.Id})
		sm := TelegramBot{}
		smm := EditMessageTextRequestBody{
			ChatId:    m.CallbackQuery.Message.Chat.Id,
			MessageId: m.CallbackQuery.Message.MessageId,
			Text:      "?????????????????? ?????????",
			ReplyMarkup: InlineKeyboardMarkup{CreateButtonsBot[InlineKeyboardButton]([]ButtonBot[InlineKeyboardButton]{
				{Row: 1, Col: 1, Button: InlineKeyboardButton{Text: "????", CallbackData: "/settings"}},
			})},
		}
		editMessageTextToBot(&sm, smm)
	}
	if m.CallbackQuery.Data == "/testconnectozonseller" {
		var user UserDB
		coll := clientMongo.Database("MyInfantBotDB").Collection("bot_users")
		filter := bson.D{{"telegram_user.user.id", m.CallbackQuery.From.Id}}
		err := coll.FindOne(context.TODO(), filter).Decode(&user)
		if err != nil {
			panic(err)
		}
		answerCallbackQueryToBot(&bot, AnswerCallbackQueryRequestBody{
			CallbackQueryId: m.CallbackQuery.Id,
			Text:            checkAuthOzonSeller(user.TelegramUser.Settings.OzonSetting.ClientId, user.TelegramUser.Settings.OzonSetting.Token),
			ShowAlert:       false,
		})
		SendMessageToBot(&bot, SendMessageRequestBody[ReplyKeyboardMarkup, int64]{
			ChatId: m.CallbackQuery.Message.Chat.Id,
			Text:   "sdfsf",
			ReplyMarkup: ReplyKeyboardMarkup{Keyboard: CreateButtonsBot[KeyboardButton]([]ButtonBot[KeyboardButton]{
				{Row: 1, Col: 1, Button: KeyboardButton{Text: GenReportToday.String()}},
				{Row: 1, Col: 1, Button: KeyboardButton{Text: GenReportYesterday.String()}},
			}),
				ResizeKeyboard: true},
		})
	}
	if mes.Text == GenReportToday.String() {
		count := countFBO("")
		mess := "<b>???????????????????? ???????????? ???? ???????? OZON FBO:</b>\n\n"
		mess += fmt.Sprintf("    <b>???????????????????? ??????????????: %d</b> \n", count.TotalCount-count.CancelledTotalCount)
		mess += "\n"
		for value, key := range count.products {
			mess += fmt.Sprintf("        <i>%s: <b>%d</b></i> \n", value, key)
		}
		mess += "\n"
		if count.CancelledTotalCount > 0 {
			mess += fmt.Sprintf("    <b>???????????????????? ???????????????????? ??????????????: %d</b>\n", count.CancelledTotalCount)
			mess += "\n"
			for value, key := range count.CancelledProducts {
				mess += fmt.Sprintf("        <i>%s: <b>%d</b></i> \n", value, key)
			}
			mess += "\n"
		}
		mess += "------------------------------------------\n"
		mess += fmt.Sprintf("    <b>?????????? ????????????????????: %d</b>\n", count.TotalCount)
		mess += fmt.Sprintf("    <b>?????????? ??????????: %s</b>\n", count.SumCount.StringFixed(2))
		mess += fmt.Sprintf("    <b>?????????? ?????????? ?????? ???????????????? OZON: %s</b>\n", count.SumWithoutCommission.StringFixed(2))
		SendMessageToBot(&bot, SendMessageRequestBody[InlineKeyboardMarkup, int64]{
			ChatId:    mes.Chat.Id,
			ParseMode: "HTML",
			Text:      mess,
		})
	}
	if mes.Text == GenReportYesterday.String() {
		count := countYesterdayFBO("")
		mess := "<b>???????????????????? ???????????? ???? ???????? OZON FBO:</b>\n\n"
		mess += fmt.Sprintf("    <b>???????????????????? ??????????????: %d</b> \n", count.TotalCount-count.CancelledTotalCount)
		mess += "\n"
		for value, key := range count.products {
			mess += fmt.Sprintf("        <i>%s: <b>%d</b></i> \n", value, key)
		}
		mess += "\n"
		if count.CancelledTotalCount > 0 {
			mess += fmt.Sprintf("    <b>???????????????????? ???????????????????? ??????????????: %d</b>\n", count.CancelledTotalCount)
			mess += "\n"
			for value, key := range count.CancelledProducts {
				mess += fmt.Sprintf("        <i>%s: <b>%d</b></i> \n", value, key)
			}
			mess += "\n"
		}
		mess += "------------------------------------------\n"
		mess += fmt.Sprintf("    <b>?????????? ????????????????????: %d</b>\n", count.TotalCount)
		mess += fmt.Sprintf("    <b>?????????? ??????????: %s</b>\n", count.SumCount.StringFixed(2))
		mess += fmt.Sprintf("    <b>?????????? ?????????? ?????? ???????????????? OZON: %s</b>\n", count.SumWithoutCommission.StringFixed(2))
		SendMessageToBot(&bot, SendMessageRequestBody[InlineKeyboardMarkup, int64]{
			ChatId:    mes.Chat.Id,
			ParseMode: "HTML",
			Text:      mess,
		})
	}
	log.Printf("???????????????? ?????????????????? %v", m)
	_, err := fmt.Fprint(w, "Hello, World!11111")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

type buttonTelegrmBot interface {
	InlineKeyboardButton | KeyboardButton
}
type ButtonBot[T buttonTelegrmBot] struct {
	Row    int
	Col    int
	Button T
}

func CreateButtonsBot[Q buttonTelegrmBot](b []ButtonBot[Q]) [][]Q {
	sort.Slice(b, func(i, j int) bool {
		return b[i].Row < b[j].Row
	})
	temp := make(map[int][]Q)
	var keys []int
	for _, iteam := range b {
		if temp[iteam.Row] == nil {
			keys = append(keys, iteam.Row)
		}
		temp[iteam.Row] = append(temp[iteam.Row], iteam.Button)

	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	matrix := make([][]Q, len(temp))

	i := 0
	for index, value := range keys {
		matrix[index] = temp[value]
		i++
	}

	for i := range matrix {
		if matrix[i] != nil {
			matrix[i] = matrix[i]
		}
	}

	return matrix
}

func (m *Message) deleteMessage() bool {
	client := &http.Client{}
	requestBody, err := json.Marshal(map[string]int64{
		"chat_id":    (*m).Chat.Id,
		"message_id": (*m).MessageId,
	})
	if err != nil {
		log.Fatalln(err)
		return false
	}
	req, err := http.NewRequest(
		"POST", urlTelegramBot+tokenTelegramBot+"/deleteMessage",
		bytes.NewBuffer(requestBody),
	)
	req.Header.Set("content-type", "application/json")
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func (t *TelegramBot) editMessageText(body interface{}) bool {
	client := &http.Client{}
	command := "editMessageText"
	if _, ok := body.(EditMessageReplyMarkupRequestBody); ok {
		command = "editMessageReplyMarkup"
	}
	requestBody, err := json.Marshal(&body)
	if err != nil {
		log.Fatalln(err)
		return false
	}
	req, err := http.NewRequest(
		"POST", urlTelegramBot+tokenTelegramBot+"/"+command,
		bytes.NewBuffer(requestBody),
	)
	req.Header.Set("content-type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer resp.Body.Close()
	return true
}

func (t *TelegramBot) sendMessage(body interface{}) bool {
	client := &http.Client{}
	requestBody, err := json.Marshal(&body)
	if err != nil {
		log.Fatalln(err)
		return false
	}
	req, err := http.NewRequest(
		"POST", urlTelegramBot+tokenTelegramBot+"/sendMessage",
		bytes.NewBuffer(requestBody),
	)
	req.Header.Set("content-type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer resp.Body.Close()
	return true
}

func (t *TelegramBot) answerCallbackQuery(body interface{}) bool {
	client := &http.Client{}
	requestBody, err := json.Marshal(&body)
	if err != nil {
		log.Fatalln(err)
		return false
	}
	req, err := http.NewRequest(
		"POST", urlTelegramBot+tokenTelegramBot+"/answerCallbackQuery",
		bytes.NewBuffer(requestBody),
	)
	req.Header.Set("content-type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer resp.Body.Close()
	return true
}

func checkAuthOzonSeller(clientId string, token string) string {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", urlOzon+"/v1/actions", nil)

	req.Header.Set("Client-Id", clientId)
	req.Header.Set("Api-Key", token)
	req.Header.Set("content-type", "application/json")

	resp, _ := client.Do(req)
	defer resp.Body.Close()
	return resp.Status
}

func countFBO(status string) ??onsolidatedReportFBO {
	var os UserDB
	//sss, _ := primitive.ObjectIDFromHex("198710657")
	coll := clientMongo.Database("MyInfantBotDB").Collection("bot_users")
	opts := options.FindOne().SetProjection(bson.D{{"telegram_user.settings.ozon_setting", 1}, {"_id", 0}})
	filter := bson.D{{"telegram_user.user.id", 198710657}}
	err1 := coll.FindOne(context.TODO(), filter, opts).Decode(&os)
	if err1 != nil {
		panic(err1)
	}
	var body ListResponseFBO
	cancelled := Cancelled
	crfbo := ??onsolidatedReportFBO{}
	crfbo.CancelledProducts = make(map[string]int)
	//tt := time.Now().Truncate(24 * time.Hour).UTC().Add(-(4 * time.Hour)).Add(-(24 * time.Hour))
	tt := time.Now().Truncate(24 * time.Hour).UTC().Add(-(4 * time.Hour))
	client := &http.Client{}
	requestBody, err := json.Marshal(ListBodyRequestFBO{
		Dir: "ASC",
		Filter: FilterFbo{
			Since:  tt.Format("2006-01-02T15:04:05Z"),
			Status: status,
			To:     tt.Add(24 * time.Hour).Format("2006-01-02T15:04:05Z"),
		},
		Limit:  100,
		Offset: 0,
	})
	if err != nil {
		log.Fatalln(err)
		return crfbo
	}
	req, err := http.NewRequest(
		"POST", urlOzon+"/v2/posting/fbo/list",
		bytes.NewBuffer(requestBody),
	)

	req.Header.Set("Client-Id", os.TelegramUser.Settings.OzonSetting.ClientId)
	req.Header.Set("Api-Key", os.TelegramUser.Settings.OzonSetting.Token)
	req.Header.Set("content-type", "application/json")

	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println(err)
		return crfbo
	}
	b, err := io.ReadAll(resp.Body)
	err = json.Unmarshal(b, &body)
	if err != nil {
		fmt.Println(err)
		return crfbo
	}
	var bb = make(map[string]int)
	replacer := strings.NewReplacer("?????????????????? Colibri ", "", "???????????????????? Colibri ", "")
	for _, aa := range body.Result {
		for _, product := range aa.Products {
			crfbo.TotalCount += product.Quantity
			if aa.Status != cancelled.String() {
				if price, err := strconv.ParseFloat(product.Price, 32); err == nil {
					crfbo.SumCount = decimal.Sum(crfbo.SumCount, decimal.NewFromFloat(price))
				}
				bb[replacer.Replace(product.Name)] += product.Quantity
			} else {
				crfbo.CancelledTotalCount += product.Quantity
				crfbo.CancelledProducts[replacer.Replace(product.Name)] += product.Quantity
			}
		}
	}
	crfbo.products = bb
	crfbo.SumWithoutCommission = decimal.NewFromFloat(crfbo.SumCount.InexactFloat64() - (0.27 * crfbo.SumCount.InexactFloat64()))
	return crfbo
}

func countYesterdayFBO(status string) ??onsolidatedReportFBO {
	var os UserDB
	//sss, _ := primitive.ObjectIDFromHex("198710657")
	coll := clientMongo.Database("MyInfantBotDB").Collection("bot_users")
	opts := options.FindOne().SetProjection(bson.D{{"telegram_user.settings.ozon_setting", 1}, {"_id", 0}})
	filter := bson.D{{"telegram_user.user.id", 198710657}}
	err1 := coll.FindOne(context.TODO(), filter, opts).Decode(&os)
	if err1 != nil {
		panic(err1)
	}
	var body ListResponseFBO
	cancelled := Cancelled
	crfbo := ??onsolidatedReportFBO{}
	crfbo.CancelledProducts = make(map[string]int)
	tt := time.Now().Truncate(24 * time.Hour).UTC().Add(-(4 * time.Hour)).Add(-(24 * time.Hour))
	client := &http.Client{}
	requestBody, err := json.Marshal(ListBodyRequestFBO{
		Dir: "ASC",
		Filter: FilterFbo{
			Since:  tt.Format("2006-01-02T15:04:05Z"),
			Status: status,
			To:     tt.Add(24 * time.Hour).Format("2006-01-02T15:04:05Z"),
		},
		Limit:  100,
		Offset: 0,
	})
	if err != nil {
		log.Fatalln(err)
		return crfbo
	}
	req, err := http.NewRequest(
		"POST", urlOzon+"/v2/posting/fbo/list",
		bytes.NewBuffer(requestBody),
	)

	req.Header.Set("Client-Id", os.TelegramUser.Settings.OzonSetting.ClientId)
	req.Header.Set("Api-Key", os.TelegramUser.Settings.OzonSetting.Token)
	req.Header.Set("content-type", "application/json")

	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println(err)
		return crfbo
	}
	b, err := io.ReadAll(resp.Body)
	err = json.Unmarshal(b, &body)
	if err != nil {
		fmt.Println(err)
		return crfbo
	}
	var bb = make(map[string]int)
	replacer := strings.NewReplacer("?????????????????? Colibri ", "", "???????????????????? Colibri ", "")
	for _, aa := range body.Result {
		for _, product := range aa.Products {
			crfbo.TotalCount += product.Quantity
			if aa.Status != cancelled.String() {
				if price, err := strconv.ParseFloat(product.Price, 32); err == nil {
					crfbo.SumCount = decimal.Sum(crfbo.SumCount, decimal.NewFromFloat(price))
				}
				bb[replacer.Replace(product.Name)] += product.Quantity
			} else {
				crfbo.CancelledTotalCount += product.Quantity
				crfbo.CancelledProducts[replacer.Replace(product.Name)] += product.Quantity
			}
		}
	}
	crfbo.products = bb
	crfbo.SumWithoutCommission = decimal.NewFromFloat(crfbo.SumCount.InexactFloat64() - (0.27 * crfbo.SumCount.InexactFloat64()))
	return crfbo
}

func findIndex[T any](obj []T, f func(e T) (result bool)) int {
	result := -1
	for i, entity := range obj {
		if f(entity) {
			return i
		}
	}
	return result
}
