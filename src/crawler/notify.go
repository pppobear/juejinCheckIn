package crawler

import (
	"net/http"
	"net/url"

	"autoSignIn/src/config"
)

func NotifyChanify(msg string) {
	_, _ = http.PostForm(config.Cfg.Chanify.Url+config.Cfg.Chanify.Token, url.Values{"text": {msg}})
}
