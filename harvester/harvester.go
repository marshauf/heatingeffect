package main

import (
	"fmt"
	"os"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/marshauf/heatingeffect/chillingeffects"
	"github.com/marshauf/heatingeffect/common"
	mgo "gopkg.in/mgo.v2"
)

func main() {
	log.Infof("os.Args: %+v", os.Args)
	app := cli.NewApp()
	app.Name = "harvester"
	app.Usage = "Harvests notices from chillingeffects.org"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "payload",
			Value: "",
			Usage: "Config to run this tool.",
		},
		cli.StringFlag{
			Name:  "e",
			Value: "",
			Usage: "Enviroment.",
		},
		cli.StringFlag{
			Name:  "id",
			Value: "",
			Usage: "Task ID.",
		},
		cli.StringFlag{
			Name:  "d",
			Value: "",
			Usage: "User writeable directory for temporary storage.",
		},
	}
	app.Action = func(c *cli.Context) {
		// Check parameters
		var (
			payloadFileName = c.String("payload")
		)
		log.Info("Check commandline parameters.")
		if len(payloadFileName) == 0 {
			log.Error("The payload parameter is empty.")
			cli.ShowAppHelp(c)
			return
		}

		// Get config
		log.Infof("Loading config file \"%s\".", payloadFileName)
		config, err := common.LoadConfig(payloadFileName)
		if err != nil {
			log.Errorf("LoadConfig: %s", err)
			return
		}

		// Set logging
		if config.RunMode == "debug" {
			log.SetLevel(log.DebugLevel)
		} else {
			log.SetLevel(log.WarnLevel)
		}
		log.Debugf("Config: %+v", config)

		// Connect to MongoDB
		log.Info("Connecting to MongoDB.")
		session, err := initMongoDB(config)
		if err != nil {
			log.Fatal(err)
		}

		log.Info("Harvesting notices and upserting to database.")
		harvest(config.IDRange.Low, config.IDRange.High, session)

		// ShutDown
		log.Info("Shutting down.")
		session.Close()
	}

	app.Run(os.Args)
}

func harvest(low, high int, session *mgo.Session) {
	ids := make(chan int)
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go work(ids, &wg, session)
	}
	for id := low; id <= high; id++ {
		ids <- id
	}
	close(ids)
	wg.Wait()
}

func work(ids <-chan int, wg *sync.WaitGroup, session *mgo.Session) {
	defer wg.Done()
	n := 0
	b := session.DB("").C("notices").Bulk()
	for id := range ids {
		notice, err := chillingeffects.RequestNotice(id)
		if err != nil {
			if err.Error() == "StatusCode: 404 Not Found" {
				continue
			}
			log.Printf("RequestNotice:%s", err)
			continue
		}
		b.Insert(notice)
		n++
		if n == 99 {
			_, err = b.Run()
			if err != nil {
				log.Printf("bulk.Run:%s", err)
			}
			b = session.DB("").C("notices").Bulk()
			n = 0
		}
	}
	if n > 0 {
		_, err := b.Run()
		if err != nil {
			log.Printf("bulk.Run:%s", err)
		}
	}
}

func initMongoDB(config *common.Config) (*mgo.Session, error) {
	if config.MongoDB == nil {
		return nil, fmt.Errorf("Config.MongoDB is nil")
	}
	dialInfo := &mgo.DialInfo{
		Addrs:    config.MongoDB.Addrs,
		Timeout:  config.MongoDB.Timeout,
		Database: config.MongoDB.Database,
		Username: config.MongoDB.Username,
		Password: config.MongoDB.Password,
	}
	return mgo.DialWithInfo(dialInfo)
}
