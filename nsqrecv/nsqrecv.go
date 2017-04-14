package nsqrecv

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"log"

	"github.com/nsqio/go-nsq"
)

var versionInfo = `nsqrecv V1.3.0 20170213`

// Conf nsqrecv main config struct
type Conf struct {
	Channel         string
	Topic           string
	NSQLookupHttpIP string
	MaxInFlight     int
}

// Handler main handler interface
type Handler interface {
	HandleMessage(message *nsq.Message) error
	Stop()
}

// HandleMessage handle message, a interface func which will be called by NSQ auto
// This func will be called serially
// interface needed by NSQ
// func (th *TailHandler) HandleMessage(m *nsq.Message) error {
// 	// message handle start
// 	payment.OrderHandle(m.Body, th.MessagesShown)
// 	th.MessagesShown++

// 	_, err := os.Stdout.Write(m.Body)
// 	if err != nil {
// 		log.Fatalf("ERROR: failed to write to os.Stdout - %s", err)
// 	}
// 	_, err = os.Stdout.WriteString("\n")
// 	if err != nil {
// 		log.Fatalf("ERROR: failed to write to os.Stdout - %s", err)
// 	}

// 	return nil
// }

// RecvStart nsq server main func
// this func will blocking process, if some err occurred, this function whill
// close the while process (using log.Fatal(err) func)
func RecvStart(conf *Conf, handler Handler) {
	if (conf.Channel == "") || (conf.Topic == "") || (conf.NSQLookupHttpIP == "") {
		log.Fatal("conf argumment error!", *conf)
	}

	cfg := nsq.NewConfig()
	// cfg.HeartbeatInterval = 5 * time.Second

	var lookupdHTTPAddrs []string
	lookupdHTTPAddrs = append(lookupdHTTPAddrs, conf.NSQLookupHttpIP)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	cfg.UserAgent = fmt.Sprintf("go-nsq/%s", nsq.VERSION)
	cfg.MaxInFlight = conf.MaxInFlight //max number of messages to allow in flight, default = 1, min = 0

	consumer, err := nsq.NewConsumer(conf.Topic, conf.Channel, cfg)
	if err != nil {
		log.Fatal(err)
	}

	consumer.AddHandler(handler)

	err = consumer.ConnectToNSQLookupds(lookupdHTTPAddrs)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case <-consumer.StopChan:
			handler.Stop() // second stop
			return
		case <-sigChan:
			consumer.Stop() //using graceful stop, first stop
		}
	}
}
