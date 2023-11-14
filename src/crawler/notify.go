package crawler

import (
	"autoSignIn/src/config"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func NotifyChanify(msg string) {
	c, _ := http.PostForm(config.Cfg.Chanify.Url+config.Cfg.Chanify.Token, url.Values{"text": {msg}})
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(c.Body)
	body, err := ioutil.ReadAll(c.Body)
	if err != nil {
		log.Fatalf("failed read body: %v", err)
	}
	fmt.Println(string(body))
}
