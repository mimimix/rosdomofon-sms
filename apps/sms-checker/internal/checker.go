package checker

import (
	"fmt"
	"log"
	"strings"

	"domofon-api/pkg/smsPoller"

	"domofon-api.gg/config"
	"github.com/imroc/req/v3"
)

func Start(poller *smsPoller.SMSPoller, config *config.Config) {
	poller.Start(func(sms smsPoller.SMS) {
		//fmt.Println("NewSMS", sms)

		if !strings.Contains(sms.Content, "domofon") {
			fmt.Println("Not domofon text")
			return
		}
		if !strings.Contains(sms.Content, config.ProtectionCode) {
			fmt.Println("Not protection code")
			return
		}

		// Create a new request client
		client := req.C()

		// Make GET request with query parameters
		resp, err := client.R().
			SetQueryParam("code", config.SecretKey).
			Get(fmt.Sprintf("http://domofonapi:%d/api/open", config.HttpPort))

		if err != nil {
			log.Printf("Error making request: %v\n", err)
			return
		}

		// Log response status and body
		log.Printf("Response status: %s\n", resp.Status)
		log.Printf("Response body: %s\n", resp.String())
	})
}
