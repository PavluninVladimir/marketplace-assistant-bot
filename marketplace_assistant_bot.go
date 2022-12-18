package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
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
	MessageId            int64    `json:"message_id"`
	MessageThreadId      int64    `json:"message_thread_id"`
	From                 User     `json:"from" bson:"from"`
	SenderChat           Chat     `json:"sender_chat"`
	Date                 int64    `json:"date"`
	Chat                 Chat     `json:"chat"`
	ForwardFrom          User     `json:"forward_from"`
	ForwardFromChat      Chat     `json:"forward_from_chat"`
	ForwardFromMessageId int64    `json:"forward_from_message_id"`
	ForwardSignature     string   `json:"forward_signature"`
	ForwardSenderName    string   `json:"forward_sender_name"`
	ForwardDate          int64    `json:"forward_date"`
	IsTopicMessage       bool     `json:"is_topic_message"`
	IsAutomaticForward   bool     `json:"is_automatic_forward"`
	ReplyToMessage       *Message `json:"reply_to_message"`
	ViaBot               User     `json:"via_bot"`
	EditDate             int64    `json:"edit_date"`
	Sticker              Sticker  `json:"sticker"`
	Text                 string   `json:"text"`
	NewChatMembers       []User   `json:"new_chat_members"`
	LeftChatMember       User     `json:"left_chat_member"`
}

type Update struct {
	UpdateId          int64   `json:"update_id"`
	Message           Message `json:"message"`
	EditedMessage     Message `json:"edited_message"`
	ChannelPost       Message `json:"channel_post"`
	EditedChannelPost Message `json:"edited_channel_post"`
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
	Url string `json:"url"`
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

type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton
}

type replyMarkup interface {
	InlineKeyboardMarkup
}

type chatId interface {
	int64 | int | string
}

type SendMessageRequestBody[T replyMarkup, Q chatId] struct {
	ChatId                   Q      `json:"chat_id"`
	MessageThreadId          int64  `json:"message_thread_id"`
	Text                     string `json:"text"`
	ParseMode                string `json:"parse_mode"`
	Entities                 []MessageEntity
	DisableWebPagePreview    bool  `json:"disable_web_page_preview"`
	DisableNotification      bool  `json:"disable_notification"`
	ProtectContent           bool  `json:"protect_content"`
	ReplyToMessageId         int64 `json:"reply_to_message_id"`
	AllowSendingWithoutReply bool  `json:"allow_sending_without_reply"`
	ReplyMarkup              T     `json:"reply_markup"`
}

var urlListFBO = "https://api-seller.ozon.ru/v2/posting/fbo/list"

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

var clientMongo *mongo.Client

type СonsolidatedReportFBO struct {
	TotalCount           int
	CancelledTotalCount  int
	SumCount             decimal.Decimal
	SumWithoutCommission decimal.Decimal
	products             map[string]int
	CancelledProducts    map[string]int
}

func main() {
	var err error
	clientMongo, err = mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8181"
		log.Printf("Defaulting to port %s", port)
	}

	router := mux.NewRouter()
	router.HandleFunc("/webhooks", webHooks)
	router.HandleFunc("/", indexHandler)
	http.Handle("/", router)

	//go test()
	log.Printf("Listening on port %s", port)
	log.Printf("Open http://localhost:%s in the browser", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
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
	if mes.LeftChatMember.Id != 0 {
		log.Printf("Left")
		mes.deleteMessage()
	}
	if mes.NewChatMembers != nil {
		log.Printf("New")
		mes.deleteMessage()
		mes.sendMessage("Добро пожаловать @" + mes.NewChatMembers[0].Username)
	}
	if mes.Text == "/start" {
		coll := clientMongo.Database("MyInfantBotDB").Collection("bot_users")
		result, er := coll.InsertOne(context.TODO(), mes.From)
		log.Println(er)
		fmt.Println(result)
	}
	if mes.Text == "/gettodayfbo" {
		count := countFBO("")
		mess := "<b>Статистика продаж за день OZON FBO:</b>\n\n"
		mess += fmt.Sprintf("    <b>Количество заказов: %d</b> \n", count.TotalCount-count.CancelledTotalCount)
		mess += "\n"
		for value, key := range count.products {
			mess += fmt.Sprintf("        <i>%s: <b>%d</b></i> \n", value, key)
		}
		mess += "\n"
		if count.CancelledTotalCount > 0 {
			mess += fmt.Sprintf("    <b>Количество отмененных заказов: %d</b>\n", count.CancelledTotalCount)
			mess += "\n"
			for value, key := range count.CancelledProducts {
				mess += fmt.Sprintf("        <i>%s: <b>%d</b></i> \n", value, key)
			}
			mess += "\n"
		}
		mess += "------------------------------------------\n"
		mess += fmt.Sprintf("    <b>Итого количество: %d</b>\n", count.TotalCount)
		mess += fmt.Sprintf("    <b>Итого сумма: %s</b>\n", count.SumCount.StringFixed(2))
		mess += fmt.Sprintf("    <b>Итого сумма без коммисии OZON: %s</b>\n", count.SumWithoutCommission.StringFixed(2))
		mes.sendMessage(mess)
	}
	log.Printf("Рассылка сообщения %v", m.Message)
	_, err := fmt.Fprint(w, "Hello, World!11111")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
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
		"POST", "https://api.telegram.org/bot5851892989:AAE5Q4QFNipbx67qGumr7pcjAUfMWWQDUNQ/deleteMessage",
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

func (m *Message) sendMessage(text string) bool {
	client := &http.Client{}
	requestBody, err := json.Marshal(SendMessageRequestBody[InlineKeyboardMarkup, int64]{
		ChatId:    (*m).Chat.Id,
		Text:      text,
		ParseMode: "HTML",
	})
	if err != nil {
		log.Fatalln(err)
		return false
	}
	req, err := http.NewRequest(
		"POST", "https://api.telegram.org/bot5851892989:AAE5Q4QFNipbx67qGumr7pcjAUfMWWQDUNQ/sendMessage",
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

func countFBO(status string) СonsolidatedReportFBO {
	var body ListResponseFBO
	cancelled := Cancelled
	crfbo := СonsolidatedReportFBO{}
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
		"POST", "https://api-seller.ozon.ru/v2/posting/fbo/list",
		bytes.NewBuffer(requestBody),
	)

	req.Header.Set("Client-Id", "391200")
	req.Header.Set("Api-Key", "b04db513-226a-4b36-8fa4-edc54d253566")
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
	replacer := strings.NewReplacer("Получешки Colibri ", "", "Полупальцы Colibri ", "")
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
	crfbo.SumWithoutCommission = decimal.NewFromFloat(crfbo.SumCount.InexactFloat64() - (0.26 * crfbo.SumCount.InexactFloat64()))
	return crfbo
}
