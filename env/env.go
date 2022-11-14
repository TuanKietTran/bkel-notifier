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

func GetVars() (botAPI string, chatID int64, oldMoodleToken string, newMoodleToken string, mongoConnStr string) {
	botAPI = lookupEnv("TELEGRAM_API")
	stringChatID := lookupEnv("CHAT_ID")

	oldMoodleToken = lookupEnv("OLD_MOODLE_TOKEN")
	newMoodleToken = lookupEnv("NEW_MOODLE_TOKEN")
	mongoConnStr = lookupEnv("MONGO_CONNECTION_STRING")

	chatID, err := strconv.ParseInt(stringChatID, 10, 64)
	if err != nil {
		log.Panicf("Can't parse chatID to int64, err: %v", err)
	}

	return botAPI, chatID, oldMoodleToken, newMoodleToken, mongoConnStr
}
