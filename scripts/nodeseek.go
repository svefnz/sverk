package scripts

import (
	"io"
	"net/http"
	"strings"
	"sverk/pkg/notify"
	"time"

	"github.com/spf13/viper"
)

func Nodeseek() {
	url := "https://www.nodeseek.com/api/attendance?random=false"
	headers := map[string]string{
		"Accept":             "*/*",
		"Accept-Encoding":    "gzip, deflate, br, zstd",
		"Accept-Language":    "zh-CN,zh;q=0.9,en;q=0.8",
		"Content-Length":     "0",
		"Origin":             "https://www.nodeseek.com",
		"Referer":            "https://www.nodeseek.com/board",
		"Sec-CH-UA":          "'Chromium';v='134', 'Not:A-Brand';v='24', 'Google Chrome';v='134'",
		"Sec-CH-UA-Mobile":   "?0",
		"Sec-CH-UA-Platform": "'Windows'",
		"Sec-Fetch-Dest":     "empty",
		"Sec-Fetch-Mode":     "cors",
		"Sec-Fetch-Site":     "same-origin",
		"User-Agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36",
		"Cookie":             viper.GetString("nodeseek.cookie"),
	}

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		notify.Bark("Nodeseek", "服务内部错误 "+err.Error())
		return
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	client := &http.Client{
		Timeout: viper.GetDuration("global.timeout") * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		notify.Bark("Nodeseek", "服务内部错误 "+err.Error())
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		notify.Bark("Nodeseek", "服务内部错误 "+err.Error())
		return
	}
	bodyStr := string(body)
	if resp.StatusCode == 200 {
		notify.Bark("Nodeseek", "[✓] 签到成功\n"+bodyStr)
	} else if resp.StatusCode == 500 && strings.Contains(string(body), "今天已完成签到") {
		notify.Bark("Nodeseek", "[x] 重复签到\n"+bodyStr)
	} else {
		notify.Bark("Nodeseek", "[x] 签到失败\n"+bodyStr)
	}
}
