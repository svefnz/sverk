package scripts

import (
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sverk/pkg/notify"
	"time"

	"github.com/spf13/viper"
)

func Hifini() {
	var (
		signInURL = "https://www.hifini.com/sg_sign.htm"
		headers   = map[string]string{
			"Cookie":             viper.GetString("hifini.cookie"),
			"authority":          "www.hifini.com",
			"accept":             "text/plain, */*; q=0.01",
			"accept-language":    "zh-CN,zh;q=0.9",
			"origin":             "https://www.hifini.com",
			"referer":            "https://www.hifini.com/",
			"sec-ch-ua":          "'Not.A/Brand';v='8', 'Chromium';v='114', 'Google Chrome';v='114'",
			"sec-ch-ua-mobile":   "?0",
			"sec-ch-ua-platform": "'macOS'",
			"sec-fetch-dest":     "empty",
			"sec-fetch-mode":     "cors",
			"sec-fetch-site":     "same-origin",
			"user-agent":         "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/135.0.0.0 Safari/537.36",
			"x-requested-with":   "XMLHttpRequest",
		}
		client = http.Client{
			Timeout: viper.GetDuration("global.timeout") * time.Second,
		}
	)

	signVal := getSignVal()
	if signVal == "" {
		notify.Bark("Hifini", "[x] 签到失败\n无法获取到 sign 参数")
		return
	}

	formData := url.Values{}
	formData.Set("sign", signVal)

	req, err := http.NewRequest("POST", signInURL, strings.NewReader(formData.Encode()))
	if err != nil {
		notify.Bark("Hifini", "服务内部错误 "+err.Error())
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		notify.Bark("Hifini", "服务内部错误 "+err.Error())
		return
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		notify.Bark("Hifini", "服务内部错误 "+err.Error())
		return
	}

	rspText := strings.ReplaceAll(strings.ReplaceAll(string(body), "\n", ""), " ", "")
	// fmt.Println(rspText)

	if strings.Contains(rspText, "成功签到") {
		notify.Bark("Hifini", "[✓] 签到成功\n"+rspText)
	} else if strings.Contains(rspText, "签过") {
		notify.Bark("Hifini", "[-] 重复签到\n"+rspText)
	} else {
		notify.Bark("Hifini", "[x] 签到失败\n"+rspText)
	}
}

func getSignVal() string {
	headers := map[string]string{
		"accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
		"accept-language":           "zh-CN,zh;q=0.9",
		"cookie":                    viper.GetString("hifini.cookie"),
		"priority":                  "u=0, i",
		"referer":                   "https://www.hifini.com/sg_sign.htm",
		"sec-ch-ua":                 "'Chromium';v='124', 'Google Chrome';v='124', 'Not-A.Brand';v='99'",
		"sec-ch-ua-mobile":          "?0",
		"sec-ch-ua-platform":        "'Windows'",
		"sec-fetch-dest":            "document",
		"sec-fetch-mode":            "navigate",
		"sec-fetch-site":            "same-origin",
		"sec-fetch-user":            "?1",
		"upgrade-insecure-requests": "1",
		"user-agent":                "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36",
	}

	client := &http.Client{
		Timeout: viper.GetDuration("global.timeout") * time.Second,
	}

	req, err := http.NewRequest("GET", "https://www.hifini.com/sg_sign.htm", nil)
	if err != nil {
		notify.Bark("Hifini", "服务内部错误 "+err.Error())
		return ""
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		notify.Bark("Hifini", "服务内部错误 "+err.Error())
		return ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		notify.Bark("Hifini", "服务内部错误 "+err.Error())
		return ""
	}

	responseText := string(body)

	pattern := `var sign = "([0-9a-f]+)"`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(responseText)

	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}
