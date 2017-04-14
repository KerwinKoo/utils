package nsqsend

import (
	"errors"
	"fmt"
	"log"
	
	"github.com/KerwinKoo/utils"
)

// SendToNSQ send order to NSQ server
func SendToNSQ(nsqServer, nsqTopic string, data []byte) error {
	nsqTopicServer := fmt.Sprintf("%s/put?topic=%s", nsqServer, nsqTopic)
	log.Println("nsq addr =", nsqTopicServer)

	responseBody, err := utils.PostBuffer2URL(data, nsqTopicServer, "application/json")

	// if responseBody = OK, return true
	// if err occured or response not OK, return false
	resultCheck := func() bool {
		if err != nil {
			log.Println("ERROR: nsq post error:", err)
			return false
		}

		if string(responseBody) != "OK" {
			log.Println("Response unknown:", responseBody)
			return false
		}
		log.Println("** NSQ response OK! **")
		return true
	}

	postResult := resultCheck()

	if postResult == false {
		errStr := "order post to NSQ failed"
		return errors.New(errStr)
	}

	return nil
}
