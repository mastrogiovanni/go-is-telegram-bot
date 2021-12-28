package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Key struct {
	Public   string `bson:"public,omitempty"`
	Type     string `bson:"type,omitempty"`
	Secret   string `bson:"secret,omitempty"`
	Encoding string `bson:"encoding,omitempty"`
}

type Identity struct {
	Id  string `bson:"id,omitempty"`
	Key Key    `bson:"key,omitempty"`
}

func getRootIdentity(bot *tgbotapi.BotAPI, chatID int64) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URL")))

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			bot.Send(tgbotapi.NewMessage(chatID, "Cannot disconnect"))
			// panic(err)
		}
	}()

	collection := client.Database("integration-service-db").Collection("identity-keys")

	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		// log.Fatal(err)
		bot.Send(tgbotapi.NewMessage(chatID, "Error reading from database: cannot find"))
		return
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {

		var identity Identity

		err := cur.Decode(&identity)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(chatID, "Error reading from database: cannot decode"))
			// log.Fatal(err)
			return
		}

		fmt.Printf("%v\n", identity.Key.Secret)
		msg := tgbotapi.NewMessage(chatID, identity.Key.Secret)
		bot.Send(msg)

	}
	if err := cur.Err(); err != nil {
		// log.Fatal(err)
		bot.Send(tgbotapi.NewMessage(chatID, "Error reading from database: cannot close db"))
	}
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_KEY"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			getRootIdentity(bot, update.Message.Chat.ID)
		}
	}
}
