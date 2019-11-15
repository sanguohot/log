package child

import (
	"github.com/sanguohot/log/v2"
	"time"
)

func Test() {
	for {
		log.Sugar.Debug(">>>>>>>>>>>> sugar debug")
		log.Sugar.Info(">>>>>>>>>>>> sugar info")
		log.Sugar.Warn(">>>>>>>>>>>> sugar warn")
		log.Sugar.Error(">>>>>>>>>>>> sugar error")

		log.Logger.Debug(">>>>>>>>>>>> logger debug")
		log.Logger.Info(">>>>>>>>>>>> logger info")
		log.Logger.Warn(">>>>>>>>>>>> logger warn")
		log.Logger.Error(">>>>>>>>>>>> logger error")
		time.Sleep(time.Second * 10)
	}
}
