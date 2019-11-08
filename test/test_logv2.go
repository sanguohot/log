package main

import (
	"github.com/sanguohot/log/v2"
	"time"
)

func main() {
	for {
		log.SugarV2.Debug(">>>>>>>>>>>> sugar debug")
		log.SugarV2.Info(">>>>>>>>>>>> sugar info")
		log.SugarV2.Warn(">>>>>>>>>>>> sugar warn")
		log.SugarV2.Error(">>>>>>>>>>>> sugar error")

		log.LoggerV2.Debug(">>>>>>>>>>>> logger debug")
		log.LoggerV2.Info(">>>>>>>>>>>> logger info")
		log.LoggerV2.Warn(">>>>>>>>>>>> logger warn")
		log.LoggerV2.Error(">>>>>>>>>>>> logger error")
		time.Sleep(time.Second * 10)
	}
}
