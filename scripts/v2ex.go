package scripts

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sverk/pkg/notify"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/viper"
)

func V2ex() {
	var (
		headers = map[string]string{
			"Cookie":             viper.GetString("v2ex.cookie"),
			"Accept":             "*/*",
			"Accept-Language":    "en,zh-CN;q=0.9,zh;q=0.8,ja;q=0.7,zh-TW;q=0.6",
			"cache-control":      "max-age=0",
			"pragma":             "no-cache",
			"Referer":            "https://www.v2ex.com/",
			"sec-ch-ua":          "'Chromium';v='124', 'Google Chrome';v='124', 'Not-A.Brand';v='99'",
			"sec-ch-ua-mobile":   "?0",
			"Sec-Ch-Ua-Platform": "Windows",
			"Sec-Fetch-Dest":     "document",
			"Sec-Fetch-Mode":     "navigate",
			"Sec-Fetch-Site":     "none",
			"Sec-Fetch-User":     "?1",
			"User-Agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36",
		}
		client = &http.Client{
			Timeout: viper.GetDuration("global.timeout") * time.Second,
		}
		isCheckined = false
	)

	req, err := http.NewRequest("GET", "https://www.v2ex.com/mission/daily", nil)
	if err != nil {
		notify.Bark("V2ex", "服务内部错误 "+err.Error())
		return
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		notify.Bark("V2ex", "服务内部错误 "+err.Error())
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		notify.Bark("V2ex", "服务内部错误 "+err.Error())
		return
	}

	if strings.Contains(string(body), "你要查看的页面需要先登录") {
		notify.Bark("V2ex", "[x] 签到失败\nCookie 已已失效")
		return
	} else if strings.Contains(string(body), "每日登录奖励已领取") {
		isCheckined = true
	}

	fmt.Println(regexp.MustCompile(`once=(\d+)`).FindStringSubmatch(string(body))[1])
	////
	if isCheckined {
		notify.Bark("V2ex", "[-] 重复签到")
		return
	}
	matches := regexp.MustCompile(`once=(\d+)`).FindStringSubmatch(string(body))
	fmt.Println(matches[1])
	if len(matches) > 1 {
		http.Get("https://www.v2ex.com/mission/daily/redeem?once=" + matches[1])
	} else {
		notify.Bark("V2ex", "[x] 签到失败\n未获取到 once 参数")
		return
	}
	////

	req, err = http.NewRequest("GET", "https://www.v2ex.com/balance", nil)
	if err != nil {
		notify.Bark("V2ex", "服务内部错误 "+err.Error())
		return
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err = client.Do(req)
	if err != nil {
		notify.Bark("V2ex", "服务内部错误 "+err.Error())
		return
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		notify.Bark("V2ex", "服务内部错误 "+err.Error())
		return
	}

	msg := doc.Find("#Main table tbody tr:nth-child(2) td:nth-child(5)").Text()
	balance := strings.Split(strings.TrimSpace(doc.Find(".balance_area.bigger").First().Text()), "  ")

	notify.Bark("V2ex", fmt.Sprintf("[✓] 签到成功 %s\n当前余额: 🟡%s\n⚪%s", msg, balance[0], balance[1]))
}
