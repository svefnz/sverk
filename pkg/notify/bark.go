package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

type barkPayload struct {
	Title string `json:"title"`
	// Subtitle  string `json:"subtitle"`
	Body      string `json:"body"`
	DeviceKey string `json:"device_key"`
	Icon      string `json:"icon"`
}

func Bark(title, content string) error {
	secretKey := viper.GetString("bark.secretKey")

	if secretKey == "" {
		return fmt.Errorf("请配置 Bark 的 SecretKey")
	}

	appName := viper.GetString("global.appName")
	if appName != "" && title == "" {
		title = appName
	} else if appName != "" && title != "" {
		title = fmt.Sprintf("%s - %s", appName, title)
	}

	payload := barkPayload{
		Title:     title,
		Body:      content,
		DeviceKey: secretKey,
		Icon:      viper.GetString("global.icon"),
	}

	jd, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://api.day.app/push", bytes.NewBuffer(jd))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	client := &http.Client{
		Timeout: viper.GetDuration("global.timeout") * time.Second,
	}
	_, err = client.Do(req)
	if err != nil {
		return err
	}
	// defer resp.Body.Close()
	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	return err
	// }

	// fmt.Printf("[.] 推送完成: %s\n", string(body))
	return nil
}
