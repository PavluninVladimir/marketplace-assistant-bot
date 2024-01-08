package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var conDB *gorm.DB

type Bot interface {
	addBot()
	//addUser()
	ListBots()
	//listUsers()
}

type ConnectDB interface {
	connDb()
}

type TelegramBot struct {
	Id    string
	Name  string
	Token string
	gorm.Model
}

func (bot *TelegramBot) addBot() {
	conDB.Create(&bot)
}

func (bot *TelegramBot) ListBots() []TelegramBot {
	connDB()
	var tbs []TelegramBot
	conDB.Find(&tbs)
	return tbs
}

func connDB() {
	var err error
	conDB, err = gorm.Open(sqlite.Open("../bot.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	//err = conDB.AutoMigrate(&TelegramBot{})
	//if err != nil {
	//	return
	//}
}

//func AddBot() {

//
//	// Create
//	db.Create(&Bot{Name: "D42", Token: "sf"})
//}
