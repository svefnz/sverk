package scripts

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"regexp"
	"sverk/pkg/notify"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/viper"
)

func Enshan() {
	headers := map[string]string{
		"User-Agent": "Mozilla/5.0 (Windows NT 6.3; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.5359.125 Safari/537.36",
		"Cookie":     viper.GetString("enshan.cookie"),
	}

	client := &http.Client{
		Timeout: viper.GetDuration("global.timeout") * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	req, err := http.NewRequest("GET", "https://www.right.com.cn/FORUM/home.php?mod=spacecp&ac=credit&showcredit=1", nil)
	if err != nil {
		notify.Bark("恩山论坛", "服务内部错误 "+err.Error())
		return
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		notify.Bark("恩山论坛", "服务内部错误 "+err.Error())
		return
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		notify.Bark("恩山论坛", "服务内部错误 "+err.Error())
		return
	}

	matches := regexp.MustCompile(`恩山币:\s*(\d+)\s*币`).FindStringSubmatch(doc.Find(".creditl.mtm.bbda.cl .cl").Text())
	var credits string
	if len(matches) > 1 {
		credits = matches[1]
	} else {
		notify.Bark("恩山论坛", "服务内部错误 获取恩山币信息失败")
		return
	}

	last := doc.Find(".bm.bw0 table tbody tr:nth-child(2) td:nth-child(4)").Text()
	lastFormat, err := time.Parse("2006-01-02 15:04", last)
	if err != nil {
		notify.Bark("恩山论坛", "服务内部错误 "+err.Error())
		return
	}

	now := time.Now()
	if now.Year() == lastFormat.Year() && now.Month() == lastFormat.Month() && now.Day() == lastFormat.Day() {
		notify.Bark("恩山论坛", fmt.Sprintf("[x] 重复签到\n上次签到时间: %s\n当前恩山币: %s", last, credits))
		return
	} else {
		notify.Bark("恩山论坛", fmt.Sprintf("[✓] 签到成功\n上次签到时间: %s\n当前恩山币: %s", last, credits))
	}
}
