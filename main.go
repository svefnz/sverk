package main

import (
	"fmt"
	"os"
	"sverk/pkg/notify"
	"sverk/scripts"

	"github.com/robfig/cron/v3"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	service string
)

func init() {
	pflag.StringVarP(&service, "service", "s", "hifini", "服务名称")
}

func main() {
	// init config
	viper.SetConfigFile("./conf/config.yml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	args := os.Args
	// start once
	if len(args) > 1 && args[1] == "start" {
		scripts.Hifini()
		return
	} else if len(args) > 1 && args[1] == "serve" {
		c := cron.New(cron.WithSeconds())

		// hifini
		if viper.GetBool("hifini.enable") {
			_, err = c.AddFunc(viper.GetString("hifini.cron"), func() { scripts.Hifini() })
			if err != nil {
				notify.Bark("Hifini", "服务内部错误 "+err.Error())
				return
			}
			fmt.Printf("[service] Hifini - [cron] %s\n", viper.GetString("hifini.cron"))
		}

		// v2ex
		if viper.GetBool("v2ex.enable") {
			_, err = c.AddFunc(viper.GetString("v2ex.cron"), func() { scripts.V2ex() })
			if err != nil {
				notify.Bark("V2ex", "服务内部错误 "+err.Error())
				return
			}
			fmt.Printf("[service] V2ex - [cron] %s\n", viper.GetString("v2ex.cron"))
		}

		// nodeseek
		if viper.GetBool("nodeseek.enable") {
			_, err = c.AddFunc(viper.GetString("nodeseek.cron"), func() { scripts.Nodeseek() })
			if err != nil {
				notify.Bark("Nodeseek", "服务内部错误 "+err.Error())
				return
			}
			fmt.Printf("[service] Nodeseek - [cron] %s\n", viper.GetString("nodeseek.cron"))
		}

		c.Start()
		notify.Bark("服务已启动", "")
		select {}
	}

	// flag
	pflag.Parse()

	// if service == "hifini" {
	// 	scripts.Hifini()
	// } else if service == "v2ex" {
	// 	scripts.V2ex()
	// } else if service == "nodeseek" {
	// 	scripts.Nodeseek()
	// } else if service == "node" {
	// 	fmt.Println("未知服务")
	// }
	switch service {
	case "hifini":
		scripts.Hifini()
	case "v2ex":
		scripts.V2ex()
	case "nodeseek":
		scripts.Nodeseek()
	case "enshan":
		scripts.Enshan()
	default:
		fmt.Println("未知服务")
	}
}
