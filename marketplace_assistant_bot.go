package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

//purchase price

type GroupProducts struct {
	NameGroup     string  `bson:"name_group"`
	PurchasePrice float64 `bson:"purchase_price"`
}
type ProductSetting struct {
	Cost          float64         `bson:"cost"`
	GroupProducts []GroupProducts `bson:"group_products"`
}

type OzonSetting struct {
	ClientId       string         `bson:"client_id"`
	Token          string         `bson:"token"`
	ProductSetting ProductSetting `bson:"product_setting"`
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
	AwaitingRegistration Status = iota // ожидает регистрации,
	AcceptanceInProgress               // идёт приёмка,
	AwaitingApprove                    // ожидает подтверждения,
	AwaitingPackaging                  // ожидает упаковки,
	AwaitingDeliver                    // ожидает отгрузки,
	Arbitration                        // арбитраж,
	ClientArbitration                  // клиентский арбитраж доставки,
	Delivering                         // доставляется,
	DriverPickup                       // у водителя,
	Delivered                          // доставлено,
	Cancelled                          // отменено,
	NotAccepted                        // не принят на сортировочном центре,
	SentBySeller                       // отправлено продавцом.
)

func (s Status) String() string {
	return [...]string{"awaiting_registration", "acceptance_in_progress", "awaiting_approve", "awaiting_packaging",
		"awaiting_deliver", "arbitration",
		"client_arbitration", "delivering", "driver_pickup", "delivered", "cancelled", "not_accepted",
		"sent_by_seller"}[s]
}

type CommandBot int

const (
	// SetClientIdOzonSetting Команда сохранения ClientId Ozon в параметры клиента
	SetClientIdOzonSetting CommandBot = iota
	GenReportToday
	GenReportYesterday
	GenReportArbitraryDate
)

func (c CommandBot) String() string {
	return [...]string{
		"setclientidozonsetting",
		"Сформировать отчет за сегодня",
		"Сформировать отчет за вчера",
		"Сформировать отчет за произвольную дату",
	}[c]
}

type СonsolidatedReportFBO struct {
	TotalCount                        int
	CancelledTotalCount               int
	SumCount                          decimal.Decimal
	SumWithoutCommission              decimal.Decimal
	SumWithoutCommissionPurchasePrice decimal.Decimal
	products                          map[string]int
	CancelledProducts                 map[string]int
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
	deleteMessage(body interface{}) bool
}

func SendMessageToBot(bot SendMessageBot, body interface{}) {
	bot.sendMessage(body)
}

// answerCallbackQueryToBot Реакция на нажатие кнопки под сообщение
func answerCallbackQueryToBot(bot AnswerCallbackQueryBot, body interface{}) {
	bot.answerCallbackQuery(body)
}

func DeleteMessageToBot(bot DeleteMessageBot, body interface{}) {
	bot.deleteMessage(body)
}

func EditMessageTextToBot(bot EditMessageTextBot, body interface{}) {
	bot.editMessageText(body)
}

type TelegramBot struct{}

type ReportMarketplace interface {
	orderSummaryReport(userId int64, filter FilterFbo) СonsolidatedReportFBO
}

type UserRepository interface {
	getOzonSetting(id int64) (*OzonSetting, error)
}
type Marketplace interface {
	ReportMarketplace
}

type OzonMarketplace struct{}

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
		log.Panic("Token telegram бота не обнаружен")
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
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./public"))))
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
				Text:   "ClientId успешно сохранен.",
				ReplyMarkup: InlineKeyboardMarkup{
					InlineKeyboard: CreateButtonsBot[InlineKeyboardButton]([]ButtonBot[InlineKeyboardButton]{
						{Row: 1, Col: 1, Button: InlineKeyboardButton{Text: "ClientId", CallbackData: "/setclientidozonsetting"}},
						{Row: 1, Col: 2, Button: InlineKeyboardButton{Text: "Token", CallbackData: "/settokenozonsetting"}},
						{Row: 2, Col: 1, Button: InlineKeyboardButton{Text: "Назад", CallbackData: "/backsettings"}},
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
				Text:   "Token успешно сохранен.",
				ReplyMarkup: InlineKeyboardMarkup{CreateButtonsBot[InlineKeyboardButton]([]ButtonBot[InlineKeyboardButton]{
					{Row: 1, Col: 1, Button: InlineKeyboardButton{Text: "ClientId", CallbackData: "/setclientidozonsetting"}},
					{Row: 1, Col: 2, Button: InlineKeyboardButton{Text: "Token", CallbackData: "/settokenozonsetting"}},
					{Row: 2, Col: 1, Button: InlineKeyboardButton{Text: "Назад", CallbackData: "/backsettings"}},
				})},
			}
			SendMessageToBot(&sm, smm)
		}
	}
	if Cash[m.Message.From.Id+m.Message.Chat.Id].LastCommand == "/setcostozon" {
		Cash[m.Message.From.Id+m.Message.Chat.Id] = DataCash{LastCommand: ""}
		cost, err := strconv.ParseFloat(m.Message.Text, 64)
		if err != nil {
			//panic(err)
			//TODO отправить сообщение об ошибки боту
		}
		coll := clientMongo.Database("MyInfantBotDB").Collection("bot_users")
		update := bson.D{{"$set", bson.D{{"telegram_user.settings.ozon_setting.product_setting.cost", cost}}}}
		filter := bson.D{{"telegram_user.user.id", mes.From.Id}}
		opts := options.Update().SetUpsert(true)
		_, err = coll.UpdateOne(context.TODO(), filter, update, opts)
		if err != nil {
			panic(err)
		} else {
			sm := TelegramBot{}
			smm := SendMessageRequestBody[InlineKeyboardMarkup, int64]{
				ChatId: m.Message.Chat.Id,
				Text:   "Token успешно сохранен.",
				ReplyMarkup: InlineKeyboardMarkup{CreateButtonsBot[InlineKeyboardButton]([]ButtonBot[InlineKeyboardButton]{
					{Row: 1, Col: 1, Button: InlineKeyboardButton{Text: "ClientId", CallbackData: "/setclientidozonsetting"}},
					{Row: 1, Col: 2, Button: InlineKeyboardButton{Text: "Token", CallbackData: "/settokenozonsetting"}},
					{Row: 2, Col: 1, Button: InlineKeyboardButton{Text: "Назад", CallbackData: "/backsettings"}},
				})},
			}
			SendMessageToBot(&sm, smm)
		}
	}
	if strings.Contains(Cash[m.Message.From.Id+m.Message.Chat.Id].LastCommand, "/setpurchaseprice") {
		productName := strings.Split(Cash[m.Message.From.Id+m.Message.Chat.Id].LastCommand, "-")[1]
		Cash[m.Message.From.Id+m.Message.Chat.Id] = DataCash{LastCommand: ""}
		cost, err := strconv.ParseFloat(m.Message.Text, 64)
		if err != nil {
			//panic(err)
			//TODO отправить сообщение об ошибки боту
		}
		coll := clientMongo.Database("MyInfantBotDB").Collection("bot_users")
		update := bson.D{{"$set", bson.D{{"telegram_user.settings.ozon_setting.product_setting.group_products.$[elem].purchase_price", cost}}}}
		filter := bson.D{{"telegram_user.user.id", mes.From.Id}}
		opts := options.Update().SetArrayFilters(options.ArrayFilters{
			Filters: []interface{}{bson.D{
				{"elem.name_group", productName},
			}},
		})
		_, err = coll.UpdateOne(context.TODO(), filter, update, opts)
		if err != nil {
			panic(err)
		} else {
			sm := TelegramBot{}
			smm := SendMessageRequestBody[InlineKeyboardMarkup, int64]{
				ChatId: m.Message.Chat.Id,
				Text:   "Token успешно сохранен.",
				ReplyMarkup: InlineKeyboardMarkup{CreateButtonsBot[InlineKeyboardButton]([]ButtonBot[InlineKeyboardButton]{
					{Row: 1, Col: 1, Button: InlineKeyboardButton{Text: "ClientId", CallbackData: "/setclientidozonsetting"}},
					{Row: 1, Col: 2, Button: InlineKeyboardButton{Text: "Token", CallbackData: "/settokenozonsetting"}},
					{Row: 2, Col: 1, Button: InlineKeyboardButton{Text: "Назад", CallbackData: "/backsettings"}},
				})},
			}
			SendMessageToBot(&sm, smm)
		}
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
			Text:   "Добро пожаловать! Чтобы использовать бота необходимо его настроить",
			ReplyMarkup: InlineKeyboardMarkup{CreateButtonsBot[InlineKeyboardButton]([]ButtonBot[InlineKeyboardButton]{
				{Row: 1, Col: 1, Button: InlineKeyboardButton{Text: "Перейти к настройкам?", CallbackData: "/settings"}},
			})},
		}
		SendMessageToBot(&sm, smm)
		return
	}
	if m.CallbackQuery.Data == "/settings" {
		answerCallbackQueryToBot(&bot, AnswerCallbackQueryRequestBody{CallbackQueryId: m.CallbackQuery.Id})
		smm := EditMessageTextRequestBody{
			ChatId:    m.CallbackQuery.Message.Chat.Id,
			MessageId: m.CallbackQuery.Message.MessageId,
			Text:      "Выберите, пожалуйста маркетплейс который вы бы хотели настроить.",
			ReplyMarkup: InlineKeyboardMarkup{CreateButtonsBot[InlineKeyboardButton]([]ButtonBot[InlineKeyboardButton]{
				{Row: 1, Col: 1, Button: InlineKeyboardButton{Text: "OZON", CallbackData: "/ozonsetting"}},
				{Row: 2, Col: 1, Button: InlineKeyboardButton{Text: "Назад", CallbackData: "/backsettings"}},
			})},
		}
		EditMessageTextToBot(&bot, smm)
	}
	if mes.Text == "/settings" {
		smm := SendMessageRequestBody[InlineKeyboardMarkup, int64]{
			ChatId: mes.Chat.Id,
			Text:   "Выберите, пожалуйста маркетплейс который вы бы хотели настроить.",
			ReplyMarkup: InlineKeyboardMarkup{CreateButtonsBot[InlineKeyboardButton]([]ButtonBot[InlineKeyboardButton]{
				{Row: 1, Col: 1, Button: InlineKeyboardButton{Text: "OZON", CallbackData: "/ozonsetting"}},
			})},
		}
		SendMessageToBot(&bot, smm)
	}
	if m.CallbackQuery.Data == "/ozonsetting" {
		answerCallbackQueryToBot(&bot, AnswerCallbackQueryRequestBody{CallbackQueryId: m.CallbackQuery.Id})
		sm := TelegramBot{}
		smm := EditMessageTextRequestBody{
			ChatId:    m.CallbackQuery.Message.Chat.Id,
			MessageId: m.CallbackQuery.Message.MessageId,
			Text:      "Для получения данных из OZON seller необходимо указать ClientId и Token. Их можно получить в личном кабинете продавца.",
			ReplyMarkup: InlineKeyboardMarkup{CreateButtonsBot[InlineKeyboardButton]([]ButtonBot[InlineKeyboardButton]{
				{Row: 1, Col: 1, Button: InlineKeyboardButton{Text: "ClientId", CallbackData: "/setclientidozonsetting"}},
				{Row: 1, Col: 2, Button: InlineKeyboardButton{Text: "Token", CallbackData: "/settokenozonsetting"}},
				{Row: 3, Col: 1, Button: InlineKeyboardButton{Text: "Настройка локального ценообразования", CallbackData: "/settinglocalpricing"}},
				{Row: 4, Col: 1, Button: InlineKeyboardButton{Text: "Проверка подключения к Ozon Seller", CallbackData: "/testconnectozonseller"}},
				{Row: 5, Col: 1, Button: InlineKeyboardButton{Text: "Назад", CallbackData: "/backsettings"}},
			})},
		}
		EditMessageTextToBot(&sm, smm)
	}
	if m.CallbackQuery.Data == "/settinglocalpricing" {
		answerCallbackQueryToBot(&bot, AnswerCallbackQueryRequestBody{CallbackQueryId: m.CallbackQuery.Id})
		sm := TelegramBot{}
		smm := EditMessageTextRequestBody{
			ChatId:    m.CallbackQuery.Message.Chat.Id,
			MessageId: m.CallbackQuery.Message.MessageId,
			Text:      "Для получения данных из OZON seller необходимо указать ClientId и Token. Их можно получить в личном кабинете продавца.",
			ReplyMarkup: InlineKeyboardMarkup{CreateButtonsBot[InlineKeyboardButton]([]ButtonBot[InlineKeyboardButton]{
				{Row: 1, Col: 1, Button: InlineKeyboardButton{Text: "Внести % сборов OZON", CallbackData: "/setcostozon"}},
				{Row: 2, Col: 1, Button: InlineKeyboardButton{Text: "Указать закупочную цену групп товаров", CallbackData: "/settingpurchaseprice"}},
			})},
		}
		EditMessageTextToBot(&sm, smm)
	}
	if m.CallbackQuery.Data == "/settingpurchaseprice" {
		answerCallbackQueryToBot(&bot, AnswerCallbackQueryRequestBody{CallbackQueryId: m.CallbackQuery.Id})
		set, _ := UserDB{}.getOzonSetting(m.CallbackQuery.From.Id)
		var buttons []ButtonBot[InlineKeyboardButton]
		for i, gp := range set.ProductSetting.GroupProducts {
			text := fmt.Sprintf("%s (Цена: %s)", gp.NameGroup, decimal.NewFromFloat(gp.PurchasePrice).StringFixed(2))
			buttons = append(buttons, ButtonBot[InlineKeyboardButton]{
				Row:    i + 1,
				Col:    1,
				Button: InlineKeyboardButton{Text: text, CallbackData: "/setpurchaseprice-" + gp.NameGroup},
			})
		}

		sm := TelegramBot{}
		smm := EditMessageTextRequestBody{
			ChatId:      m.CallbackQuery.Message.Chat.Id,
			MessageId:   m.CallbackQuery.Message.MessageId,
			Text:        "Для получения данных из OZON seller необходимо указать ClientId и Token. Их можно получить в личном кабинете продавца.",
			ReplyMarkup: InlineKeyboardMarkup{CreateButtonsBot[InlineKeyboardButton](buttons)},
		}
		EditMessageTextToBot(&sm, smm)
	}
	if strings.Contains(m.CallbackQuery.Data, "/setpurchaseprice") {
		answerCallbackQueryToBot(&bot, AnswerCallbackQueryRequestBody{CallbackQueryId: m.CallbackQuery.Id})
		Cash[m.CallbackQuery.From.Id+m.CallbackQuery.Message.Chat.Id] = DataCash{LastCommand: m.CallbackQuery.Data}
		sm := TelegramBot{}
		smm := SendMessageRequestBody[InlineKeyboardMarkup, int64]{
			ChatId: m.CallbackQuery.Message.Chat.Id,
			Text:   "ОК. Пришлите, пожалуйста себистоимость товара.",
		}
		SendMessageToBot(&sm, smm)
	}
	if m.CallbackQuery.Data == "/setcostozon" {
		answerCallbackQueryToBot(&bot, AnswerCallbackQueryRequestBody{CallbackQueryId: m.CallbackQuery.Id})
		Cash[m.CallbackQuery.From.Id+m.CallbackQuery.Message.Chat.Id] = DataCash{LastCommand: "/setcostozon"}
		sm := TelegramBot{}
		smm := SendMessageRequestBody[InlineKeyboardMarkup, int64]{
			ChatId: m.CallbackQuery.Message.Chat.Id,
			Text:   "ОК. Пришлите, пожалуйста % расходом на услуги OZON.",
		}
		SendMessageToBot(&sm, smm)
	}
	if m.CallbackQuery.Data == "/settokenozonsetting" {
		answerCallbackQueryToBot(&bot, AnswerCallbackQueryRequestBody{CallbackQueryId: m.CallbackQuery.Id})
		Cash[m.CallbackQuery.From.Id+m.CallbackQuery.Message.Chat.Id] = DataCash{LastCommand: "/settokenozonsetting"}
		sm := TelegramBot{}
		smm := SendMessageRequestBody[InlineKeyboardMarkup, int64]{
			ChatId: m.CallbackQuery.Message.Chat.Id,
			Text:   "ОК. Пришлите, пожалуйста Token для бота.",
		}
		SendMessageToBot(&sm, smm)
	}
	if m.CallbackQuery.Data == "/setclientidozonsetting" {
		answerCallbackQueryToBot(&bot, AnswerCallbackQueryRequestBody{CallbackQueryId: m.CallbackQuery.Id})
		Cash[m.CallbackQuery.From.Id+m.CallbackQuery.Message.Chat.Id] = DataCash{LastCommand: "/setclientidozonsetting"}
		sm := TelegramBot{}
		smm := SendMessageRequestBody[InlineKeyboardMarkup, int64]{
			ChatId: m.CallbackQuery.Message.Chat.Id,
			Text:   "ОК. Пришлите, пожалуйста ClientID для бота.",
		}
		SendMessageToBot(&sm, smm)
	}
	if m.CallbackQuery.Data == "/backsettings" {
		answerCallbackQueryToBot(&bot, AnswerCallbackQueryRequestBody{CallbackQueryId: m.CallbackQuery.Id})
		sm := TelegramBot{}
		smm := EditMessageTextRequestBody{
			ChatId:    m.CallbackQuery.Message.Chat.Id,
			MessageId: m.CallbackQuery.Message.MessageId,
			Text:      "Настроить бота?",
			ReplyMarkup: InlineKeyboardMarkup{CreateButtonsBot[InlineKeyboardButton]([]ButtonBot[InlineKeyboardButton]{
				{Row: 1, Col: 1, Button: InlineKeyboardButton{Text: "Да", CallbackData: "/settings"}},
			})},
		}
		EditMessageTextToBot(&sm, smm)
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
				{Row: 1, Col: 1, Button: KeyboardButton{Text: GenReportArbitraryDate.String(), WebApp: &WebAppInfo{
					Url: "https://bot.my-infant.com/static/",
				}}},
				{Row: 2, Col: 1, Button: KeyboardButton{Text: GenReportToday.String()}},
				{Row: 2, Col: 2, Button: KeyboardButton{Text: GenReportYesterday.String()}},
			}),
				ResizeKeyboard: true},
		})
	}
	if mes.Text == GenReportToday.String() {
		var marketplace Marketplace = &OzonMarketplace{}
		filter := FilterFbo{
			Since:  time.Now().Truncate(24 * time.Hour).UTC().Add(-(4 * time.Hour)).Format(time.RFC3339),
			Status: "",
			To:     time.Now().Truncate(24 * time.Hour).UTC().Add(-(4 * time.Hour)).Add(24 * time.Hour).Format(time.RFC3339),
		}

		SendMessageToBot(&bot, SendMessageRequestBody[InlineKeyboardMarkup, int64]{
			ChatId:    mes.Chat.Id,
			ParseMode: "HTML", //TODO приминить паттерн стратегия
			Text:      printOrderSummaryReport(marketplace.orderSummaryReport(m.Message.From.Id, filter)),
		})
	}
	if mes.Text == GenReportYesterday.String() {
		var marketplace Marketplace = &OzonMarketplace{}
		filter := FilterFbo{
			Since:  time.Now().Truncate(24 * time.Hour).UTC().Add(-(4 * time.Hour)).Add(-(24 * time.Hour)).Format(time.RFC3339),
			Status: "",
			To:     time.Now().Truncate(24 * time.Hour).UTC().Add(-(4 * time.Hour)).Add(24 * time.Hour).Add(-(24 * time.Hour)).Format(time.RFC3339),
		}
		SendMessageToBot(&bot, SendMessageRequestBody[InlineKeyboardMarkup, int64]{
			ChatId:    mes.Chat.Id,
			ParseMode: "HTML", //TODO приминить паттерн стратегия
			Text:      printOrderSummaryReport(marketplace.orderSummaryReport(m.Message.From.Id, filter)),
		})
	}
	if mes.WebAppData.ButtonText == GenReportArbitraryDate.String() {
		var marketplace Marketplace = &OzonMarketplace{}
		data := strings.Split(mes.WebAppData.Data, "::")
		from, _ := time.Parse("2006-01-02", data[0])
		to, _ := time.Parse("2006-01-02", data[1])
		filter := FilterFbo{
			Since:  from.Truncate(24 * time.Hour).UTC().Add(-(4 * time.Hour)).Add(-(24 * time.Hour)).Format(time.RFC3339),
			Status: "",
			To:     to.Truncate(24 * time.Hour).UTC().Add(-(4 * time.Hour)).Add(24 * time.Hour).Add(-(24 * time.Hour)).Format(time.RFC3339),
		}
		SendMessageToBot(&bot, SendMessageRequestBody[InlineKeyboardMarkup, int64]{
			ChatId:    mes.Chat.Id,
			ParseMode: "HTML", //TODO приминить паттерн стратегия
			Text:      printOrderSummaryReport(marketplace.orderSummaryReport(m.Message.From.Id, filter)),
		})
	}
	log.Printf("Рассылка сообщения %v", m)
	_, err := fmt.Fprint(w, "Hello, World!11111")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func printOrderSummaryReport(c СonsolidatedReportFBO) string {
	mess := "<b>Статистика продаж за день OZON FBO:</b>\n\n"
	mess += fmt.Sprintf("    <b>Количество заказов: %d</b> \n", c.TotalCount)
	mess += "\n"
	for value, key := range c.products {
		mess += fmt.Sprintf("        <i>%s: <b>%d</b></i> \n", value, key)
	}
	mess += "\n"
	if c.CancelledTotalCount > 0 {
		mess += fmt.Sprintf("    <b>Количество отмененных заказов: %d</b>\n", c.CancelledTotalCount)
		mess += "\n"
		for value, key := range c.CancelledProducts {
			mess += fmt.Sprintf("        <i>%s: <b>%d</b></i> \n", value, key)
		}
		mess += "\n"
	}
	mess += "------------------------------------------\n"
	mess += fmt.Sprintf("    <b>Итого количество: %d</b>\n", c.TotalCount-c.CancelledTotalCount)
	mess += fmt.Sprintf("    <b>Итого сумма: %s</b>\n", c.SumCount.StringFixed(2))
	mess += fmt.Sprintf("    <b>Итого сумма без комиссии OZON: %s</b>\n", c.SumWithoutCommission.StringFixed(2))
	mess += fmt.Sprintf("    <b>Итого доход: %s</b>\n",
		c.SumWithoutCommissionPurchasePrice.StringFixed(2))
	return mess
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

func (t *TelegramBot) deleteMessage(body interface{}) bool {
	client := &http.Client{}
	requestBody, err := json.Marshal(&body)
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
func (m UserDB) getOzonSetting(id int64) (*OzonSetting, error) {
	coll := clientMongo.Database("MyInfantBotDB").Collection("bot_users")
	opts := options.FindOne().SetProjection(bson.D{{"telegram_user.settings.ozon_setting", 1}, {"_id", 0}})
	filter := bson.D{{"telegram_user.user.id", id}}
	err := coll.FindOne(context.TODO(), filter, opts).Decode(&m)
	return &m.TelegramUser.Settings.OzonSetting, err
}

func (m UserDB) setProductGroupSetting(userId int64, name string) {
	var userDB UserDB
	coll := clientMongo.Database("MyInfantBotDB").Collection("bot_users")
	filter := bson.D{{"telegram_user.user.id", userId}}
	err := coll.FindOne(context.TODO(), filter).Decode(&userDB)
	if err != nil {
		//TODO написать обработку ошибки
	}

	pl := userDB.TelegramUser.Settings.OzonSetting.ProductSetting.GroupProducts
	isFindProduct := findIndex[GroupProducts](pl, func(ps GroupProducts) bool {
		if ps.NameGroup == name {
			return true
		}
		return false
	})
	if isFindProduct == -1 {
		newPl := append(pl, GroupProducts{NameGroup: name, PurchasePrice: 0})
		update := bson.D{{"$set", bson.D{{"telegram_user.settings.ozon_setting.product_setting.group_products", newPl}}}}
		opts := options.Update().SetUpsert(true)
		_, err := coll.UpdateOne(context.TODO(), filter, update, opts)
		if err != nil {
			//TODO написать обработку ошибки
		}
	}
}

func fboListHandler(userId int64, body ListBodyRequestFBO) (*ListResponseFBO, error) {
	client := &http.Client{}
	requestBody, err := json.Marshal(body)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	r, err := http.NewRequest(
		"POST", urlOzon+"/v2/posting/fbo/list",
		bytes.NewBuffer(requestBody),
	)
	var userDb UserRepository = UserDB{}
	setting, err := userDb.getOzonSetting(userId)
	if err != nil {
		panic(err)
	}

	r.Header.Set("Client-Id", setting.ClientId)
	r.Header.Set("Api-Key", setting.Token)
	r.Header.Set("content-type", "application/json")

	response, err := client.Do(r)
	defer response.Body.Close()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var l ListResponseFBO
	b, err := io.ReadAll(response.Body)
	err = json.Unmarshal(b, &l)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &l, nil
}

func (m *OzonMarketplace) orderSummaryReport(userId int64, filter FilterFbo) СonsolidatedReportFBO {
	defer fmt.Println("sss")
	crfbo := СonsolidatedReportFBO{}
	var userDb UserRepository = UserDB{}
	setting, err := userDb.getOzonSetting(userId)
	if err != nil {
		return crfbo
	}
	cost := setting.ProductSetting.Cost
	mp := make(map[string]float64)
	for _, iteam := range setting.ProductSetting.GroupProducts {
		mp[iteam.NameGroup] = iteam.PurchasePrice
	}
	crfbo.CancelledProducts = make(map[string]int)
	var bb = make(map[string]int)
	replacer := strings.NewReplacer("Получешки Colibri ", "", "Полупальцы Colibri ", "")
	limit := 1000
	offset := 0
	var resp = ListResponseFBO{}
	for {
		response, err := fboListHandler(userId, ListBodyRequestFBO{
			Dir:    "ASC",
			Filter: filter,
			Limit:  int64(limit),
			Offset: int64(offset),
		})
		resp.Result = append(resp.Result, response.Result...)
		if err != nil {
			panic(err)
		}
		if len(response.Result) == 0 {
			break
		}
		offset += limit
	}
	var pp float64
	for _, aa := range resp.Result {
		for _, product := range aa.Products {
			crfbo.TotalCount += product.Quantity
			user := UserDB{}
			// TODO массовое изменение товаров или горутину
			user.setProductGroupSetting(userId, replacer.Replace(product.Name))
			if aa.Status != Cancelled.String() {
				pp += mp[replacer.Replace(product.Name)]
				if price, err := strconv.ParseFloat(product.Price, 64); err == nil {
					crfbo.SumCount = decimal.NewFromFloat(crfbo.SumCount.InexactFloat64() + price)
				}
				bb[replacer.Replace(product.Name)] += product.Quantity
			} else {
				crfbo.CancelledTotalCount += product.Quantity
				crfbo.CancelledProducts[replacer.Replace(product.Name)] += product.Quantity
			}

		}
	}
	crfbo.products = bb
	crfbo.SumWithoutCommission = decimal.NewFromFloat(crfbo.SumCount.InexactFloat64() - ((cost / 100) * crfbo.SumCount.InexactFloat64()))
	crfbo.SumWithoutCommissionPurchasePrice = decimal.NewFromFloat(crfbo.SumWithoutCommission.InexactFloat64() - pp)
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
