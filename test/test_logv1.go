package main

import (
	"github.com/sanguohot/log"
	"go.uber.org/zap"
	"time"
)

func main() {
	for {
		log.Sugar.Debug(">>>>>>>>>>>> sugar debug")
		log.Sugar.Info(">>>>>>>>>>>> sugar info")
		log.Sugar.Warn(">>>>>>>>>>>> sugar warn")
		log.Sugar.Error(">>>>>>>>>>>> sugar error")

		log.Logger.Debug(">>>>>>>>>>>> logger debug")
		log.Logger.Info(">>>>>>>>>>>> logger info")
		log.Logger.Warn(">>>>>>>>>>>> logger warn")
		log.Logger.Error(">>>>>>>>>>>> logger error", zap.String("logger", "v1"))
		time.Sleep(time.Second * 10)
	}
}
