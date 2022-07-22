package env

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

func LoadDotEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Panicf("Error loading .env file, err: %v", err)
	}
}

func lookupEnv(key string) (value string) {
	value, found := os.LookupEnv(key)
	if !found {
		log.Fatalf("%s not exists, please specify it", key)
	}

	return value
}

func GetVars() (string, string, int64) {
	botAPI := lookupEnv("TELEGRAM_API")
	moodleToken := lookupEnv("MOODLE_TOKEN")
	chatId := lookupEnv("CHAT_ID")

	chatIdInt64, err := strconv.ParseInt(chatId, 10, 64)
	if err != nil {
		log.Panicf("Can't parse chatId to int64, err: %v", err)
	}

	return botAPI, moodleToken, chatIdInt64
}
