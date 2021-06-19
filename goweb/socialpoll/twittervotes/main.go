package main

import (
	"github.com/bitly/go-nsq"
	"gopkg.in/mgo.v2"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	var stoplock sync.Mutex
	stop := false
	stopChan := make(chan struct{}, 1)
	signalChan := make(chan os.Signal, 1)
	go func() {
		<-signalChan
		stoplock.Lock()
		stop = true
		stoplock.Unlock()
		log.Println("Stopping...")
		stopChan <- struct{}{}
		closeCon()
		closeDb()
	}()
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT)

	if err := dialDb(); err != nil {
		log.Fatalln("failed to dial MongoDB:", err)
	}

	// start things
	votes := make(chan string)
	publisherStoppedChan := publishVotes(votes)
	twitterStoppedChan := startTwitterStream(stopChan, votes)
	go func() {
		time.Sleep(1 * time.Minute)
		closeCon()
		stoplock.Lock()
		if stop {
			stoplock.Unlock()
			return
		}
		stoplock.Unlock()
	}()
	<-twitterStoppedChan
	close(votes)
	<-publisherStoppedChan
}

var (
	db *mgo.Session
)

func dialDb() error {
	var err error
	log.Println("dialing mongodb: localhost")
	db, err = mgo.Dial("localhost")
	return err
}

func closeDb() {
	db.Close()
	log.Println("closed database connection")
}

type poll struct {
	Options []string
}

func loadOptions() ([]string, error) {
	var options []string
	iter := db.DB("ballots").C("polls").Find(nil).Iter()
	var p poll
	for iter.Next(&p) {
		options = append(options, p.Options...)
	}
	iter.Close()
	return options, iter.Err()
}

func publishVotes(votes <-chan string) <-chan struct{} {
	stopchan := make(chan struct{}, 1)
	pub, _ := nsq.NewProducer("localhost:4150", nsq.NewConfig())
	go func() {
		for vote := range votes {
			pub.Publish("votes", []byte(vote))
		}
		log.Println("Publisher: Stopping")
		pub.Stop()
		log.Println("Publisher: Stopped")
		stopchan <- struct{}{}
	}()
	return stopchan
}
