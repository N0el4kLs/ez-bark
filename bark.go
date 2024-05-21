package ezbark

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/imroc/req/v3"
)

// Options struct to store data in need
type Options struct {
	TestSwitch bool
	Send
	Conf
}

// Send post data
type Send struct {
	Messages []string
	Body     string `json:"body" yaml:"body"`
	Title    string `json:"title" yaml:"title"`
	Badge    int    `json:"badge" yaml:"badge"`
	Icon     string `json:"icon" yaml:"icon"`
	Group    string `json:"group" yaml:"group"`
	Url      string `json:"url" yaml:"url"`
	Sound    string `json:"sound" yaml:"sound"`
	Key      string `json:"device_key" yaml:"key"`
}

// Conf store config data
type Conf struct {
	BarkServer string `yaml:"server"`
	DeviceKey  string `yaml:"devicekey"`
}

// NewOptions create base options
func NewOptions() *Options {
	return &Options{}
}

// DefaultConf load default config from ~/.config/ezbark/config.yml
func (o *Options) DefaultConf() *Options {
	conf, _ := os.UserHomeDir()
	confFile := filepath.Join(conf, ".config", "ezbark", "config.yml")
	fmt.Println(confFile)
	confFileObj, _ := os.Open(confFile)
	defer confFileObj.Close()

	yaml.NewDecoder(confFileObj).Decode(&o.Conf)
	return o
}

// SetSend set filed
func (o *Options) SetSend(send map[string]interface{}) *Options {
	for key, value := range send {
		switch strings.ToLower(key) {
		case "title":
			v := value.(string)
			o.Title = v
		case "body":
			v := value.(string)
			o.Messages = append(o.Messages, v)
		case "badge":
			v := value.(int)
			o.Badge = v
		case "icon":
			v := value.(string)
			o.Icon = v
		case "group":
			v := value.(string)
			o.Group = v
		case "url":
			v := value.(string)
			o.Url = v
		case "sound":
			v := value.(string)
			o.Sound = v
		case "server":
			v := value.(string)
			o.Conf.BarkServer = v
		case "key":
			v := value.(string)
			o.Conf.DeviceKey = v
		}
	}
	return o
}

// Notice send notification
func (o *Options) Notice() error {
	if (o.Conf.DeviceKey == "" && o.Send.Key == "") || o.BarkServer == "" {
		log.Fatalln("Config file err.Please check the bark server or device key.")
		os.Exit(0)
	}
	if o.Conf.DeviceKey != "" && o.Send.Key == "" {
		o.Send.Key = o.Conf.DeviceKey
	}
	err := notice(*o)
	if err != nil {
		return err
	}
	return nil
}

// Notice command line main logic
func notice(o Options) error {
	for _, v := range o.Send.Messages {
		o.Send.Body = v
		send, err := json.Marshal(o.Send)
		if err != nil {
			log.Fatalln("Data struct err")
			return err
		}

		client := req.C().Post().SetBodyJsonBytes(send)
		if strings.HasSuffix(o.BarkServer, "/") {
			o.BarkServer = o.BarkServer[:len(o.BarkServer)-1]
		}
		resp, err := client.Post(o.BarkServer + "/push")
		if err != nil {
			log.Fatalln("Server maybe down")
			return err
		}
		fmt.Println(resp)
	}
	return nil
}
