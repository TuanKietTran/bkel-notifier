package main

import (
	"bkmail/env"
	"bkmail/forwarder"
	"bkmail/forwarder/adapter"
	"bkmail/moodle"
	"context"
	"crypto/tls"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
)

const LogFileName = "bkel-fetching.log"

const newMoodleHostname = "e-learning.hcmut.edu.vn"

//const oldMoodleHostname = "e-learning-old.hcmut.edu.vn"

func main() {
	env.LoadDotEnv()
	logFile := openFileForLogging(LogFileName)

	_, _, _, _, mongoConnStr := env.GetVars()

	tlsCfg := generateTLSConfigs()
	credential := options.Credential{
		AuthMechanism: "MONGODB-X509",
	}

	mongoClientOpts := options.Client().ApplyURI(mongoConnStr).SetAuth(credential).SetTLSConfig(tlsCfg)

	mongoClient, err := mongo.Connect(context.TODO(), mongoClientOpts)
	if err != nil {
		log.Fatalf("Cannot connect to Mongo")
	}

	err = mongoClient.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		log.Fatalf("Cannot ping to server")
	}

	coll := mongoClient.Database("subscription").Collection("bkmail")
	cursor, err := coll.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Fatalf("Cannot get all docs")
	}

	for cursor.Next(context.TODO()) {
		result := bson.M{}
		err := cursor.Decode(&result)
		if err != nil {
			log.Printf("Cannot parsed 1 record, skip")
		}

		log.Printf("%v", result)
		var email string = fmt.Sprint(result["email"])
		moodleToken := fmt.Sprint(result["bkel_token"])

		log.Println("Bot created successfully")

		// New Moodle
		newMoodleClient := moodle.NewClient(newMoodleHostname, moodleToken)

		newMessages := newMoodleClient.FetchMessages()

		for _, msg := range newMessages {
			if msg.IsNotification == 1 {
				newMoodleClient.MarkNotificationAsRead(msg.Id)
			} else {
				newMoodleClient.MarkChatAsRead(msg.Id)
			}
			parsedMsg := adapter.RenderMessage(msg)

			log.Printf("Received message from calendar_forward %s", msg.UserFrom)
			forwarder.SendMailSimple(parsedMsg, email)
			log.Printf("Handled message from calendar_forward: %v", msg.Id)
		}
	}

	_ = logFile.Close()
}

func openFileForLogging(filename string) *os.File {
	newLogFile, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Can't open log file, err: %v", err)
	}

	log.SetOutput(newLogFile)
	return newLogFile
}

func generateTLSConfigs() *tls.Config {
	cert, err := tls.LoadX509KeyPair("X509-cert-591008662129803129.pem", "X509-cert-591008662129803129.pem")
	if err != nil {
		log.Fatalf("Cannot load X509 file, err: %v", err)

	}

	cfg := &tls.Config{Certificates: []tls.Certificate{cert}}

	return cfg
}
