package main

import (
	"flag"
	"fmt"
	"github.com/bitly/go-nsq"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const updateDuration = 1 * time.Second

var fatalError error

func fatal(e error) {
	fmt.Println(e)
	flag.PrintDefaults()
	fatalError = e
}

func main() {
	defer func() {
		if fatalError != nil {
			os.Exit(1)
		}
	}()

	log.Println("Connecting to database...")
	db, err := mgo.Dial("localhost")
	if err != nil {
		fatal(err)
		return
	}
	defer func() {
		log.Println("Closing database connection...")
		db.Close()
	}()

	pollData := db.DB("ballots").C("polls")

	var countLock sync.Mutex
	var counts map[string]int

	log.Println("Connecting to nsq...")
	q, err := nsq.NewConsumer("votes", "counter", nsq.NewConfig())
	if err != nil {
		fatal(err)
		return
	}

	q.AddHandler(nsq.HandlerFunc(func(m *nsq.Message) error {
		countLock.Lock()
		defer countLock.Unlock()
		if counts == nil {
			counts = make(map[string]int)
		}
		vote := string(m.Body)
		counts[vote]++
		return nil
	}))

	if err := q.ConnectToNSQLookupd("localhost:4161"); err != nil {
		fatal(err)
		return
	}

	ticker := time.NewTicker(updateDuration)
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM)
	for {
		select {
		case <-ticker.C:
			doCount(&countLock, &counts, pollData)
		case <-termChan:
			ticker.Stop()
			q.Stop()
		case <-q.StopChan:
			return
		}
	}
}

func doCount(countLock *sync.Mutex, counts *map[string]int, pollData *mgo.Collection) {
	countLock.Lock()
	defer countLock.Unlock()
	if len(*counts) == 0 {
		log.Println("No new votes, skipping database update")
		return
	}

	log.Println("Updating database...")
	log.Println(*counts)

	ok := true
	for option, count := range *counts {
		sel := bson.M{
			"options": bson.M{
				"$in": []string{option},
			}}
		up := bson.M{
			"$inc": bson.M{
				"results." + option: count,
			},
		}
		if _, err := pollData.UpdateAll(sel, up); err != nil {
			log.Println("failed to update: ", err)
			ok = false
		}
	}

	if ok {
		log.Println("Finished updating database...")
		*counts = nil
	}
}
