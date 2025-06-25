package smsPoller

import (
	"domofon-api/pkg/huaweimodem"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"time"

	"domofon-api.gg/config"
)

type SMSPoller struct {
	modem        *huaweimodem.Device
	ticker       *time.Ticker
	lastSmsIds   []int
	lastSmsFile  string
	aliveSmsTime int
}

type SMS struct {
	Id      int
	Date    time.Time
	Phone   string
	Content string
}

type NewSMSEvent = func(SMS)

func New(modem *huaweimodem.Device, config *config.Config) *SMSPoller {
	poller := &SMSPoller{
		modem:        modem,
		lastSmsFile:  config.LastSmsFile,
		aliveSmsTime: config.SmsAliveTime,
	}

	err := poller.readDatabase()
	if err != nil {
		panic(err)
	}

	return poller
}

func (p *SMSPoller) readDatabase() error {
	file, err := os.Open(p.lastSmsFile)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Database not found")
			p.lastSmsIds = []int{}
			return nil
		}
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&p.lastSmsIds)
	return err
}

func (p *SMSPoller) writeDatabase() error {
	file, err := os.Create(p.lastSmsFile)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(p.lastSmsIds)
}

func (p *SMSPoller) poll(event NewSMSEvent) {
	smsList, err := p.modem.ReadSMSInbox()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, message := range smsList.Messages {
		if !slices.Contains(p.lastSmsIds, message.Index) {
			p.lastSmsIds = append(p.lastSmsIds, message.Index)
			go func() {
				err := p.writeDatabase()
				if err != nil {
					fmt.Println(err)
				}
			}()

			fmt.Printf("New SMS %v\n", message)

			date, err := time.Parse("2006-01-02 15:04:05", message.Date)
			if err != nil {
				fmt.Println(err)
				return
			}

			if time.Since(date).Seconds() > float64(p.aliveSmsTime) {
				fmt.Printf("SMS %d is too old\n", message.Index)
				continue
			}

			event(SMS{
				Id:      message.Index,
				Date:    date,
				Phone:   message.Phone,
				Content: message.Content,
			})
		}
	}
}

func (p *SMSPoller) Start(event NewSMSEvent) {
	p.ticker = time.NewTicker(5 * time.Second)
	go func() {
		for range p.ticker.C {
			p.poll(event)
		}
	}()
}
