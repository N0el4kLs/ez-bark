package util

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"

	ezbark "github.com/N0el4kLs/ez-bark"
)

var (
	Body     string
	Title    string
	Badge    int
	Icon     string
	Group    string
	Url      string
	Sound    string
	Init     string
	ShowConf bool
	Test     bool
	initf    string
	file     string

	send    *ezbark.Send
	conf    *ezbark.Conf
	options *ezbark.Options

	Bodys = make([]string, 0, 1)
)

func newConf(server, key string) *ezbark.Conf {
	return &ezbark.Conf{
		BarkServer: server,
		DeviceKey:  key,
	}
}

// InitOptions Init command line options
func InitOptions() (*ezbark.Options, error) {
	flag.StringVar(&Body, "m", "", "Content of the message")
	flag.StringVar(&Title, "t", "", "Title of the message (font size would be larger than the body)")
	flag.IntVar(&Badge, "b", 1, "The number displayed next to App icon")
	flag.StringVar(&Icon, "i", "", "An url to the icon, available only on iOS 15 or later")
	flag.StringVar(&Group, "g", "", "The group of the notification")
	flag.StringVar(&Url, "u", "", "Url that will jump when click notification")
	flag.StringVar(&Sound, "s", "", "The sound of the notification")
	flag.StringVar(&Init, "init", "", "Init server address and device key (Format is server,devicekey, "+
		"Using the comma as a separators,like http://127.0.0.1,abcdefg)")
	flag.BoolVar(&ShowConf, "showconf", false, "show barkserver and devicekeys")
	flag.BoolVar(&Test, "test", false, "send a test message")
	flag.StringVar(&initf, "initf", "", "Generate message file")
	flag.StringVar(&file, "file", "", "load message config file")

	flag.Parse()

	// handler init options
	// if init options is not empty, then generate a config file
	if Init != "" {
		confSilce := strings.Split(Init, ",")
		server, deviceKey := confSilce[0], confSilce[1]
		globalConf := newConf(server, deviceKey)
		if err := checkGlobalConf(); err != nil {
			err = createGlobalConf()
			if err != nil {
				return nil, err
			}
		}
		globalConfObj, err := os.OpenFile(globalConfigFile(), os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Printf("Witre config file err")
			return nil, err
		}
		defer globalConfObj.Close()

		// Unmarshal yaml strcut into config.yml
		err = yaml.NewEncoder(globalConfObj).Encode(globalConf)
		if err != nil {
			log.Println("Generate template for config file err:", err)
			return nil, err
		}
		return nil, nil
	}

	err := checkGlobalConf()
	if err != nil {
		log.Fatalln("Can not find global config.yml.Please use init option to set config...")
		return nil, err
	}

	confFile := globalConfigFile()
	confFileObj, _ := os.Open(confFile)
	conf = &ezbark.Conf{}
	err = yaml.NewDecoder(confFileObj).Decode(conf)
	if err != nil {
		log.Println("Load config file error")
		return nil, err
	}

	if Sound != "" {
		Sound = Sound + ".caf"
	}
	if ShowConf {
		fmt.Println("BarkServer: ", conf.BarkServer)
		fmt.Println("Devicekey: ", conf.DeviceKey)
		os.Exit(0)
	}

	if initf != "" {
		messageConf := &ezbark.Send{}
		messageFile := filepath.Join(getCurrentPath(), initf)
		messageFileObj, _ := os.OpenFile(messageFile, os.O_CREATE|os.O_WRONLY, 0666)
		defer messageFileObj.Close()

		err := yaml.NewEncoder(messageFileObj).Encode(messageConf)
		if err != nil {
			log.Println("Generate message config file error")
			return nil, err
		}
		log.Println("Generate message config file successfully:", messageFile)
		os.Exit(0)
	}

	if hasStdin() {
		s := bufio.NewScanner(os.Stdin)
		bodyBuilder := strings.Builder{}
		lens := 0
		for s.Scan() {
			l, _ := bodyBuilder.WriteString(s.Text())
			bodyBuilder.WriteString("\n")
			if lens < 500 && lens+l+1 > 500 {
				Bodys = append(Bodys, bodyBuilder.String())
				bodyBuilder.Reset()
				lens = 0
			}
			lens += l + 1
		}
		// when the message is less than 500 characters
		if len(Bodys) == 0 {
			Bodys = append(Bodys, bodyBuilder.String())
		}
	}

	if Body == "" && !Test && !hasStdin() {
		log.Fatalln("Messages of the message is required")
	}
	if Body != "" {
		Bodys = append(Bodys, Body)
	}

	send = &ezbark.Send{
		Messages: Bodys,
		Title:    Title,
		Badge:    Badge,
		Icon:     Icon,
		Group:    Group,
		Url:      Url,
		Sound:    Sound,
	}

	if file != "" {
		messageFile := filepath.Join(getCurrentPath(), file)
		messageFileObj, err := os.Open(messageFile)
		defer messageFileObj.Close()
		if err != nil {
			log.Println("Load message config file error")
		}

		err = yaml.NewDecoder(messageFileObj).Decode(send)
		if err != nil {
			log.Println("Parse message config file error")
		}
	}

	send.Key = conf.DeviceKey
	options = &ezbark.Options{
		TestSwitch: Test,
		Send:       *send,
		Conf:       *conf,
	}

	return options, nil
}

func hasStdin() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return false
	}

	return (stat.Mode() & os.ModeCharDevice) == 0
}
