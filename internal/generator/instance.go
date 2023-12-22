package generator

import (
	"promogen/internal/utils"

	"fmt"
	"github.com/zenthangplus/goccm"
	"io"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"sync"
)

type Instance struct {
	Client *http.Client
	Body   []byte

	Cookie string
}

type Payload map[string]any
type Header struct {
	Cookie          string `json:"cookie"`
	UserAgent       string `json:"user-agent"`
	Referer         string `json:"referer"`
	ContentLength   string `json:"content-length"`
	SecChUaPlatform string `json:"sec-ch-ua-platform"`
	SecChUa         string `json:"sec-ch-ua"`
}
type OperaResponseToken struct {
	Token string `json:"token"`
}

func CreateInstances() (instances []Instance, err error) {
	fmt.Println("initializing..")

	cfg, err := utils.LoadConfig("config.json")
	if err != nil {
		return nil, err
	}

	var mu sync.Mutex

	tp := &http.Transport{}
	if cfg.Proxy != "" {
		proxy, err := url.Parse("http://" + cfg.Proxy)
		if err != nil {
			return nil, err
		}
		tp.Proxy = http.ProxyURL(proxy)
	}

	routines := cfg.Cfg.Threads
	if cfg.Cfg.Limit > 0 {
		routines = cfg.Cfg.Limit
	}
	wg := goccm.New(routines)

	for i := 0; i < cfg.Cfg.Threads; i++ {
		wg.Wait()
		go func() {
			defer wg.Done()

			mu.Lock()
			instances = append(instances, Instance{
				//Cookie: utils.Cookies("discord.com"),
				Client: &http.Client{
					Transport: tp,
				},
			})
			mu.Unlock()
		}()
	}
	wg.WaitAllDone()

	return instances, nil
}

func (in *Instance) Request(method, url string, payload Payload, headers *Header) (*http.Response, error) {
	req, err := http.NewRequest(method, url, utils.Payload(payload))
	if err != nil {
		return nil, err
	}

	in.Header(*headers, req)

	resp, err := in.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Println(err)
			return
		}
	}(resp.Body)

	in.Body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (in *Instance) Header(headers Header, req *http.Request) {
	for h, o := range map[string]string{
		"accept":           "*/*",
		"accept-language":  "en-US,en-GB;q=0.9",
		"content-type":     "application/json",
		"sec-ch-ua-mobile": "?0",
		"sec-fetch-dest":   "empty",
		"sec-fetch-mode":   "cors",
		"sec-fetch-site":   "same-origin",
		"Sec-Ch-Ua":        `"Opera GX";v="106", "Chromium";v="120", "Not?A_Brand";v="8"`,
		"user-agent":       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	} {
		req.Header.Set(h, o)
	}

	v := reflect.ValueOf(headers)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i).Interface()

		if tag, ok := field.Tag.Lookup("json"); ok && value != "" {
			req.Header.Set(tag, fmt.Sprintf("%v", value))
		}
	}
}
