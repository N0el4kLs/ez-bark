package ezbark

import (
	"fmt"
	"testing"
)

func TestNotice(t *testing.T) {
	// use map[string]interface{} to store your own data
	send := make(map[string]interface{})
	options := NewOptions()

	// Use DefaultConf to load default config in ~/.config/ezbark/config.yml
	// or you can use SetSend options to load bark server or devices
	options.DefaultConf()

	fmt.Printf("%#v\n", options)

	send["title"] = "title"
	send["body"] = "title"
	//send["server"] = "http://127.0.0.1:7070"
	//send["key"] = "abcdefg"

	options = options.SetSend(send)
	options.Notice()
}
