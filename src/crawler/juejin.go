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

func sendRequest(url, method string, postData io.Reader) *juejinResp {
	req, err := http.NewRequest(method, url, postData)
	if err != nil {
		log.Fatalf("failed NewRequest: %v", err)
	}
	for k, v := range Header {
		req.Header.Set(k, v)
	}
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
	resp := sendRequest(globalBigUrl, "POST", strings.NewReader(""))
	postData := map[string]interface{}{
		"lottery_history_id": resp.Data.(map[string]interface{})["lotteries"].([]interface{})[0].(map[string]interface{})["history_id"]}
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

func RunTask() string {
	return fmt.Sprintf(`
【掘金】
当前矿石数量: %s
自动签到结果: %s
自动抽奖结果: %s
沾幸运值结果: %s`,
		getCurPoint(), checkIn(), lottery(), dipLucky())
}
