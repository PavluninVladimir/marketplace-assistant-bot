package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gqlgen"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"telegram"
	"time"
)

var (
	clientMongo      *mongo.Client
	urlOzon          string
	urlTelegramBot   string
	tokenTelegramBot string
)

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
	NameBot  string          `bson:"name_bot"`
	User     telegram.User   `bson:"user"`
	Chats    []telegram.Chat `bson:"chats"`
	Settings Settings        `bson:"settings"`
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
	//var configDefault = Config{
	//	Next:          nil,
	//	Done:          nil,
	//	Format:        "[${time}] ${status} - ${latency} ${method} ${path}\n",
	//	TimeFormat:    "15:04:05",
	//	TimeZone:      "Local",
	//	TimeInterval:  500 * time.Millisecond,
	//	Output:        os.Stdout,
	//	DisableColors: false,
	//}
	app := fiber.New()
	app.Use(
		logger.New(), // add Logger middleware
	)
	routeGQR := "/query"
	query, playground := gqlgen.GraphQLPlaygroundHandler(routeGQR)
	app.Get("/playground", adaptor.HTTPHandlerFunc(playground))
	app.Post("/query", adaptor.HTTPHandler(query))
	app.Listen(":" + port)

	//router := mux.NewRouter()
	//router.HandleFunc("/webhooks", webHooks)
	//http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./public"))))
	//router.HandleFunc("/", indexHandler)
	//http.Handle("/", router)

	//log.Printf("Listening on port %s", port)
	//log.Printf("Open http://localhost:%s in the browser", port)
	//log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
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
	var m telegram.Update
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
			smm := telegram.SendMessageRequestBody[telegram.InlineKeyboardMarkup, int64]{
				ChatId: m.Message.Chat.Id,
				Text:   "ClientId успешно сохранен.",
				ReplyMarkup: telegram.InlineKeyboardMarkup{
					InlineKeyboard: CreateButtonsBot[telegram.InlineKeyboardButton]([]telegram.ButtonBot[telegram.InlineKeyboardButton]{
						{Row: 1, Col: 1, Button: telegram.InlineKeyboardButton{Text: "ClientId", CallbackData: "/setclientidozonsetting"}},
						{Row: 1, Col: 2, Button: telegram.InlineKeyboardButton{Text: "Token", CallbackData: "/settokenozonsetting"}},
						{Row: 2, Col: 1, Button: telegram.InlineKeyboardButton{Text: "Назад", CallbackData: "/backsettings"}},
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
			smm := telegram.SendMessageRequestBody[telegram.InlineKeyboardMarkup, int64]{
				ChatId: m.Message.Chat.Id,
				Text:   "Token успешно сохранен.",
				ReplyMarkup: telegram.InlineKeyboardMarkup{CreateButtonsBot[telegram.InlineKeyboardButton]([]telegram.ButtonBot[telegram.InlineKeyboardButton]{
					{Row: 1, Col: 1, Button: telegram.InlineKeyboardButton{Text: "ClientId", CallbackData: "/setclientidozonsetting"}},
					{Row: 1, Col: 2, Button: telegram.InlineKeyboardButton{Text: "Token", CallbackData: "/settokenozonsetting"}},
					{Row: 2, Col: 1, Button: telegram.InlineKeyboardButton{Text: "Назад", CallbackData: "/backsettings"}},
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
			smm := telegram.SendMessageRequestBody[telegram.InlineKeyboardMarkup, int64]{
				ChatId: m.Message.Chat.Id,
				Text:   "Token успешно сохранен.",
				ReplyMarkup: telegram.InlineKeyboardMarkup{CreateButtonsBot[telegram.InlineKeyboardButton]([]telegram.ButtonBot[telegram.InlineKeyboardButton]{
					{Row: 1, Col: 1, Button: telegram.InlineKeyboardButton{Text: "ClientId", CallbackData: "/setclientidozonsetting"}},
					{Row: 1, Col: 2, Button: telegram.InlineKeyboardButton{Text: "Token", CallbackData: "/settokenozonsetting"}},
					{Row: 2, Col: 1, Button: telegram.InlineKeyboardButton{Text: "Назад", CallbackData: "/backsettings"}},
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
			smm := telegram.SendMessageRequestBody[telegram.InlineKeyboardMarkup, int64]{
				ChatId: m.Message.Chat.Id,
				Text:   "Token успешно сохранен.",
				ReplyMarkup: telegram.InlineKeyboardMarkup{CreateButtonsBot[telegram.InlineKeyboardButton]([]telegram.ButtonBot[telegram.InlineKeyboardButton]{
					{Row: 1, Col: 1, Button: telegram.InlineKeyboardButton{Text: "ClientId", CallbackData: "/setclientidozonsetting"}},
					{Row: 1, Col: 2, Button: telegram.InlineKeyboardButton{Text: "Token", CallbackData: "/settokenozonsetting"}},
					{Row: 2, Col: 1, Button: telegram.InlineKeyboardButton{Text: "Назад", CallbackData: "/backsettings"}},
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
					Chats: []telegram.Chat{m.Message.Chat},
				}}
			_, err := coll.InsertOne(context.TODO(), userDB)
			if err != nil {
				panic(err)
			}
		} else {
			if result := findIndex[telegram.Chat](user.TelegramUser.Chats, func(c telegram.Chat) bool {
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
		smm := telegram.SendMessageRequestBody[telegram.InlineKeyboardMarkup, int64]{
			ChatId: m.Message.Chat.Id,
			Text:   "Добро пожаловать! Чтобы использовать бота необходимо его настроить",
			ReplyMarkup: telegram.InlineKeyboardMarkup{CreateButtonsBot[telegram.InlineKeyboardButton]([]telegram.ButtonBot[telegram.InlineKeyboardButton]{
				{Row: 1, Col: 1, Button: telegram.InlineKeyboardButton{Text: "Перейти к настройкам?", CallbackData: "/settings"}},
			})},
		}
		SendMessageToBot(&sm, smm)
		return
	}
	if m.CallbackQuery.Data == "/settings" {
		answerCallbackQueryToBot(&bot, telegram.AnswerCallbackQueryRequestBody{CallbackQueryId: m.CallbackQuery.Id})
		smm := telegram.EditMessageTextRequestBody{
			ChatId:    m.CallbackQuery.Message.Chat.Id,
			MessageId: m.CallbackQuery.Message.MessageId,
			Text:      "Выберите, пожалуйста маркетплейс который вы бы хотели настроить.",
			ReplyMarkup: telegram.InlineKeyboardMarkup{CreateButtonsBot[telegram.InlineKeyboardButton]([]telegram.ButtonBot[telegram.InlineKeyboardButton]{
				{Row: 1, Col: 1, Button: telegram.InlineKeyboardButton{Text: "OZON", CallbackData: "/ozonsetting"}},
				{Row: 2, Col: 1, Button: telegram.InlineKeyboardButton{Text: "Назад", CallbackData: "/backsettings"}},
			})},
		}
		EditMessageTextToBot(&bot, smm)
	}
	if mes.Text == "/settings" {
		smm := telegram.SendMessageRequestBody[telegram.InlineKeyboardMarkup, int64]{
			ChatId: mes.Chat.Id,
			Text:   "Выберите, пожалуйста маркетплейс который вы бы хотели настроить.",
			ReplyMarkup: telegram.InlineKeyboardMarkup{CreateButtonsBot[telegram.InlineKeyboardButton]([]telegram.ButtonBot[telegram.InlineKeyboardButton]{
				{Row: 1, Col: 1, Button: telegram.InlineKeyboardButton{Text: "OZON", CallbackData: "/ozonsetting"}},
			})},
		}
		SendMessageToBot(&bot, smm)
	}
	if m.CallbackQuery.Data == "/ozonsetting" {
		answerCallbackQueryToBot(&bot, telegram.AnswerCallbackQueryRequestBody{CallbackQueryId: m.CallbackQuery.Id})
		sm := TelegramBot{}
		smm := telegram.EditMessageTextRequestBody{
			ChatId:    m.CallbackQuery.Message.Chat.Id,
			MessageId: m.CallbackQuery.Message.MessageId,
			Text:      "Для получения данных из OZON seller необходимо указать ClientId и Token. Их можно получить в личном кабинете продавца.",
			ReplyMarkup: telegram.InlineKeyboardMarkup{CreateButtonsBot[telegram.InlineKeyboardButton]([]telegram.ButtonBot[telegram.InlineKeyboardButton]{
				{Row: 1, Col: 1, Button: telegram.InlineKeyboardButton{Text: "ClientId", CallbackData: "/setclientidozonsetting"}},
				{Row: 1, Col: 2, Button: telegram.InlineKeyboardButton{Text: "Token", CallbackData: "/settokenozonsetting"}},
				{Row: 3, Col: 1, Button: telegram.InlineKeyboardButton{Text: "Настройка локального ценообразования", CallbackData: "/settinglocalpricing"}},
				{Row: 4, Col: 1, Button: telegram.InlineKeyboardButton{Text: "Проверка подключения к Ozon Seller", CallbackData: "/testconnectozonseller"}},
				{Row: 5, Col: 1, Button: telegram.InlineKeyboardButton{Text: "Назад", CallbackData: "/backsettings"}},
			})},
		}
		EditMessageTextToBot(&sm, smm)
	}
	if m.CallbackQuery.Data == "/settinglocalpricing" {
		answerCallbackQueryToBot(&bot, telegram.AnswerCallbackQueryRequestBody{CallbackQueryId: m.CallbackQuery.Id})
		sm := TelegramBot{}
		smm := telegram.EditMessageTextRequestBody{
			ChatId:    m.CallbackQuery.Message.Chat.Id,
			MessageId: m.CallbackQuery.Message.MessageId,
			Text:      "Для получения данных из OZON seller необходимо указать ClientId и Token. Их можно получить в личном кабинете продавца.",
			ReplyMarkup: telegram.InlineKeyboardMarkup{CreateButtonsBot[telegram.InlineKeyboardButton]([]telegram.ButtonBot[telegram.InlineKeyboardButton]{
				{Row: 1, Col: 1, Button: telegram.InlineKeyboardButton{Text: "Внести % сборов OZON", CallbackData: "/setcostozon"}},
				{Row: 2, Col: 1, Button: telegram.InlineKeyboardButton{Text: "Указать закупочную цену групп товаров", CallbackData: "/settingpurchaseprice"}},
			})},
		}
		EditMessageTextToBot(&sm, smm)
	}
	if m.CallbackQuery.Data == "/settingpurchaseprice" {
		answerCallbackQueryToBot(&bot, telegram.AnswerCallbackQueryRequestBody{CallbackQueryId: m.CallbackQuery.Id})
		set, _ := UserDB{}.getOzonSetting(m.CallbackQuery.From.Id)
		var buttons []telegram.ButtonBot[telegram.InlineKeyboardButton]
		for i, gp := range set.ProductSetting.GroupProducts {
			text := fmt.Sprintf("%s (Цена: %s)", gp.NameGroup, decimal.NewFromFloat(gp.PurchasePrice).StringFixed(2))
			buttons = append(buttons, telegram.ButtonBot[telegram.InlineKeyboardButton]{
				Row:    i + 1,
				Col:    1,
				Button: telegram.InlineKeyboardButton{Text: text, CallbackData: "/setpurchaseprice-" + gp.NameGroup},
			})
		}

		sm := TelegramBot{}
		smm := telegram.EditMessageTextRequestBody{
			ChatId:      m.CallbackQuery.Message.Chat.Id,
			MessageId:   m.CallbackQuery.Message.MessageId,
			Text:        "Для получения данных из OZON seller необходимо указать ClientId и Token. Их можно получить в личном кабинете продавца.",
			ReplyMarkup: telegram.InlineKeyboardMarkup{CreateButtonsBot[telegram.InlineKeyboardButton](buttons)},
		}
		EditMessageTextToBot(&sm, smm)
	}
	if strings.Contains(m.CallbackQuery.Data, "/setpurchaseprice") {
		answerCallbackQueryToBot(&bot, telegram.AnswerCallbackQueryRequestBody{CallbackQueryId: m.CallbackQuery.Id})
		Cash[m.CallbackQuery.From.Id+m.CallbackQuery.Message.Chat.Id] = DataCash{LastCommand: m.CallbackQuery.Data}
		sm := TelegramBot{}
		smm := telegram.SendMessageRequestBody[telegram.InlineKeyboardMarkup, int64]{
			ChatId: m.CallbackQuery.Message.Chat.Id,
			Text:   "ОК. Пришлите, пожалуйста себистоимость товара.",
		}
		SendMessageToBot(&sm, smm)
	}
	if m.CallbackQuery.Data == "/setcostozon" {
		answerCallbackQueryToBot(&bot, telegram.AnswerCallbackQueryRequestBody{CallbackQueryId: m.CallbackQuery.Id})
		Cash[m.CallbackQuery.From.Id+m.CallbackQuery.Message.Chat.Id] = DataCash{LastCommand: "/setcostozon"}
		sm := TelegramBot{}
		smm := telegram.SendMessageRequestBody[telegram.InlineKeyboardMarkup, int64]{
			ChatId: m.CallbackQuery.Message.Chat.Id,
			Text:   "ОК. Пришлите, пожалуйста % расходом на услуги OZON.",
		}
		SendMessageToBot(&sm, smm)
	}
	if m.CallbackQuery.Data == "/settokenozonsetting" {
		answerCallbackQueryToBot(&bot, telegram.AnswerCallbackQueryRequestBody{CallbackQueryId: m.CallbackQuery.Id})
		Cash[m.CallbackQuery.From.Id+m.CallbackQuery.Message.Chat.Id] = DataCash{LastCommand: "/settokenozonsetting"}
		sm := TelegramBot{}
		smm := telegram.SendMessageRequestBody[telegram.InlineKeyboardMarkup, int64]{
			ChatId: m.CallbackQuery.Message.Chat.Id,
			Text:   "ОК. Пришлите, пожалуйста Token для бота.",
		}
		SendMessageToBot(&sm, smm)
	}
	if m.CallbackQuery.Data == "/setclientidozonsetting" {
		answerCallbackQueryToBot(&bot, telegram.AnswerCallbackQueryRequestBody{CallbackQueryId: m.CallbackQuery.Id})
		Cash[m.CallbackQuery.From.Id+m.CallbackQuery.Message.Chat.Id] = DataCash{LastCommand: "/setclientidozonsetting"}
		sm := TelegramBot{}
		smm := telegram.SendMessageRequestBody[telegram.InlineKeyboardMarkup, int64]{
			ChatId: m.CallbackQuery.Message.Chat.Id,
			Text:   "ОК. Пришлите, пожалуйста ClientID для бота.",
		}
		SendMessageToBot(&sm, smm)
	}
	if m.CallbackQuery.Data == "/backsettings" {
		answerCallbackQueryToBot(&bot, telegram.AnswerCallbackQueryRequestBody{CallbackQueryId: m.CallbackQuery.Id})
		sm := TelegramBot{}
		smm := telegram.EditMessageTextRequestBody{
			ChatId:    m.CallbackQuery.Message.Chat.Id,
			MessageId: m.CallbackQuery.Message.MessageId,
			Text:      "Настроить бота?",
			ReplyMarkup: telegram.InlineKeyboardMarkup{CreateButtonsBot[telegram.InlineKeyboardButton]([]telegram.ButtonBot[telegram.InlineKeyboardButton]{
				{Row: 1, Col: 1, Button: telegram.InlineKeyboardButton{Text: "Да", CallbackData: "/settings"}},
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
		answerCallbackQueryToBot(&bot, telegram.AnswerCallbackQueryRequestBody{
			CallbackQueryId: m.CallbackQuery.Id,
			Text:            checkAuthOzonSeller(user.TelegramUser.Settings.OzonSetting.ClientId, user.TelegramUser.Settings.OzonSetting.Token),
			ShowAlert:       false,
		})
		SendMessageToBot(&bot, telegram.SendMessageRequestBody[telegram.ReplyKeyboardMarkup, int64]{
			ChatId: m.CallbackQuery.Message.Chat.Id,
			Text:   "sdfsf",
			ReplyMarkup: telegram.ReplyKeyboardMarkup{Keyboard: CreateButtonsBot[telegram.KeyboardButton]([]telegram.ButtonBot[telegram.KeyboardButton]{
				{Row: 1, Col: 1, Button: telegram.KeyboardButton{Text: GenReportArbitraryDate.String(), WebApp: &telegram.WebAppInfo{
					Url: "https://bot.my-infant.com/static/",
				}}},
				{Row: 2, Col: 1, Button: telegram.KeyboardButton{Text: GenReportToday.String()}},
				{Row: 2, Col: 2, Button: telegram.KeyboardButton{Text: GenReportYesterday.String()}},
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

		SendMessageToBot(&bot, telegram.SendMessageRequestBody[telegram.InlineKeyboardMarkup, int64]{
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
		SendMessageToBot(&bot, telegram.SendMessageRequestBody[telegram.InlineKeyboardMarkup, int64]{
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
		SendMessageToBot(&bot, telegram.SendMessageRequestBody[telegram.InlineKeyboardMarkup, int64]{
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

func CreateButtonsBot[Q telegram.ButtonTelegrmBot](b []telegram.ButtonBot[Q]) [][]Q {
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
	if _, ok := body.(telegram.EditMessageReplyMarkupRequestBody); ok {
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
