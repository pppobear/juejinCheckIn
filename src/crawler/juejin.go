package crawler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"autoSignIn/src/config"
)

const (
	baseUrl      = "https://api.juejin.cn/growth_api/v1/"
	curPointUrl  = baseUrl + "get_cur_point"
	infoUrl      = baseUrl + "lottery_config/get"
	checkInUrl   = baseUrl + "check_in"
	lotteryUrl   = baseUrl + "lottery/draw"
	globalBigUrl = baseUrl + "lottery_history/global_big"
	dipLuckyUrl  = baseUrl + "lottery_lucky/dip_lucky"
)

type juejinResp struct {
	No   int         `json:"err_no"`
	Msg  string      `json:"err_msg"`
	Data interface{} `json:"data"`
}

func sendRequest(reqUrl, method string, postData io.Reader) *juejinResp {
	req, err := http.NewRequest(method, reqUrl, postData)
	if err != nil {
		log.Fatalf("failed NewRequest: %v", err)
	}
	for k, v := range Header {
		req.Header.Set(k, v)
	}
	q := req.URL.Query()
	cookies := strings.Split(config.Cfg.Cookies.JueJin, ";")
	var aid, uid string
	for _, cookie := range cookies {
		cookie = strings.TrimSpace(cookie)
		if strings.HasPrefix(cookie, "__tea_cookie_tokens_") {
			s := strings.Split(cookie, "=")
			aid = strings.TrimPrefix(s[0], "__tea_cookie_tokens_")
			if d1v, err := url.QueryUnescape(s[1]) ; err != nil {
				return nil
			} else {
				if d2v, err := url.QueryUnescape(d1v) ; err != nil {
					return nil
				} else {
					var d2vJson map[string]string
					json.Unmarshal([]byte(d2v), &d2vJson)
					uid = d2vJson["user_unique_id"]
				}
			}
			break
		}
	}
	q.Add("aid", aid)
	q.Add("uuid", uid)
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Cookie", config.Cfg.Cookies.JueJin)
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("failed Do Request: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalf("failed close body: %v", err)
		}
	}(resp.Body)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("failed read body: %v", err)
	}
	sResp := new(juejinResp)
	if err = json.Unmarshal(body, sResp); err != nil {
		return nil
	}
	return sResp
}

func getCurPoint() string {
	resp := sendRequest(curPointUrl, "GET", strings.NewReader(""))
	return fmt.Sprintf("%.f", resp.Data)
}

func hasLotteryFreeCount() (bool, error) {
	resp := sendRequest(infoUrl, "GET", strings.NewReader(""))
	if resp.Data == nil {
		return false, errors.New(resp.Msg)
	}
	return resp.Data.(map[string]interface{})["free_count"].(float64) != 0, nil
}

func checkIn() string {
	return sendRequest(checkInUrl, "POST", strings.NewReader("")).Msg
}

func lottery() string {
	toLottery, err := hasLotteryFreeCount()
	if err == nil {
		if toLottery {
			return sendRequest(lotteryUrl, "POST", strings.NewReader("")).Msg
		}
		return "今日已经抽奖"
	} else {
		return err.Error()
	}
}

func dipLucky() string {
	dipId := getDipluckId()
	if dipId == "" {
		return "获取要沾的id失败"
	}
	postData := map[string]interface{}{
		"lottery_history_id": dipId}
	if jsonBytes, err := json.Marshal(postData); err == nil {
		resp := sendRequest(dipLuckyUrl, "POST", bytes.NewReader(jsonBytes))
		if resp.Data == nil {
			return resp.Msg
		}
		return fmt.Sprintf("沾到喜气: %.f，当前幸运值: %.f",
			resp.Data.(map[string]interface{})["dip_value"],
			resp.Data.(map[string]interface{})["total_value"])
	}
	return "出错了"
}

func getDipluckId() string {
	postData := map[string]interface{}{"page_no": 1, "page_size": 5}
	if jsonBytes, err := json.Marshal(postData); err == nil {
		resp := sendRequest(globalBigUrl, "POST", bytes.NewReader(jsonBytes))
		if resp.Data == nil {
			return ""
		}
		return resp.Data.(map[string]interface{})["lotteries"].([]interface{})[0].(map[string]interface{})["history_id"].(string)
	}
	return ""
}

func RunTask() string {
	checkIn := checkIn()
	point := getCurPoint()
	lottery := lottery()
	lucky := dipLucky()
	return fmt.Sprintf(`
【掘金】
当前矿石数量: %s
自动签到结果: %s
自动抽奖结果: %s
沾幸运值结果: %s`,
		point, checkIn, lottery, lucky)
}
