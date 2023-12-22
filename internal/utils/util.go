package utils

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

const LogFormat = "\u001B[38;2;64;64;64m[\u001B[38;2;50;100;0m%dms\u001B[38;2;64;64;64m] [\u001B[38;2;50;100;0m+\u001B[38;2;64;64;64m]\u001B[39m Token: %s\n"
const CountFormat = "\u001B[38;2;64;64;64m[\u001B[38;2;100;200;0m%.2fs\u001B[38;2;64;64;64m] [\u001B[38;2;100;200;0m+\u001B[38;2;64;64;64m]\u001B[39m Codes Created: %d\n"

type Config struct {
	Cfg struct {
		Threads int `json:"threads"`
		Limit   int `json:"thread_limit"`
	} `json:"Configuration"`
	Proxy string `json:"proxy"`
}

func Payload(v any) *bytes.Buffer {
	data, err := json.Marshal(v)
	if err != nil {
		log.Println(err)
		return nil
	}
	return bytes.NewBuffer(data)
}

func Cookies(url string) (cookie string) {
	resp, err := http.Get("https://" + url)
	if err != nil {
		return Cookies(url)
	}
	defer resp.Body.Close()

	if resp.Cookies() != nil {
		for _, cookies := range resp.Cookies() {
			cookie += fmt.Sprintf("%s=%s; ", cookies.Name, cookies.Value)
		}
		return cookie
	} else {
		return Cookies(url)
	}
}

func WriteArrayToFile(path string, data []string) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND, 0660)
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()
	for i := 0; i < len(data); i++ {
		f.WriteString(data[i] + "\n")
	}
}

func LoadConfig(filename string) (Config, error) {
	var config Config
	conf, err := os.Open(filename)
	if err != nil {
		return config, err
	}
	defer conf.Close()

	if err = json.NewDecoder(conf).Decode(&config); err != nil {
		return config, err
	}

	return config, nil
}

func RandomID() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}
