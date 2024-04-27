package consumer

import (
	"encoding/json"
	"fmt"
	nats "github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"github.com/nats-io/stan.go/pb"
	"log"
	"order-service/internal/cache"
	"order-service/internal/config"
	"order-service/internal/consumer/model"
	"order-service/internal/consumer/tools"
	"order-service/internal/storage"
	"os"
	"os/signal"
)

type Consumer struct {
	subscriber stan.Subscription
	storage    *storage.Storage
	cache      *cache.Cache
}

func NewConsumer(s *storage.Storage, c *cache.Cache) *Consumer {
	cons := &Consumer{
		storage: s,
		cache:   c,
	}

	if err := cons.FillCache(); err != nil {
		log.Fatal("Populate cache error")
	}

	return cons
}

func (c *Consumer) AddOrder(content json.RawMessage) error {
	var m model.Model

	if err := json.Unmarshal(content, &m); err != nil {
		log.Println("json")
		return err
	}

	if err := tools.CheckModel(m); !err {
		return fmt.Errorf("Invalid model")
	}

	if err := c.storage.AddOrder(*m.OrderUid, content); err != nil {
		return err
	}

	c.cache.Add(*m.OrderUid, content)

	return nil
}

func (c *Consumer) MsgHandler(msg *stan.Msg) {
	if err := c.AddOrder(msg.Data); err != nil {
		log.Println(err.Error())
	}
}

func (c *Consumer) StartSubscribe(cfg config.NatsConfig) {
	nc, err := nats.Connect(cfg.URL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	sc, err := stan.Connect(cfg.Cluster, cfg.Client, stan.NatsConn(nc),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Fatalf("Connection lost, reason: %v", reason)
		}))
	if err != nil {
		log.Fatalf("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, cfg.URL)
	}
	log.Printf("Connected to %s clusterID: [%s] clientID: [%s]\n", cfg.URL, "test-cluster", "stan-pub")

	c.subscriber, err = sc.Subscribe(cfg.Subject, c.MsgHandler, stan.StartAt(pb.StartPosition_NewOnly))
	if err != nil {
		sc.Close()
		log.Fatal(err)
	}

	log.Printf("Listening on [%s], clientID=[%s], qgroup=[%s] durable=[%s]\n", cfg.Subject, "stan-pub", "qgroup", "durable")

	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for range signalChan {
			log.Printf("\nReceived an interrupt, unsubscribing and closing connection...\n\n")
			c.subscriber.Close()
			sc.Close()
			cleanupDone <- true
		}
	}()
	<-cleanupDone

	//go func() {
	//	for range ctx.Done() {
	//		log.Println("Exit NATS")
	//		c.subscriber.Unsubscribe()
	//		c.subscriber.Close()
	//		sc.Close()
	//	}
	//}()
}

func (c *Consumer) FillCache() (err error) {
	var ids []string
	var cont []json.RawMessage

	if ids, cont, err = c.storage.FillCache(); err != nil {
		return err
	}

	c.cache.Fill(ids, cont)

	return nil
}
