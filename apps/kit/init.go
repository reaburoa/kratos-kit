package kit

import (
	"github.com/welltop-cn/common/cloud/config"
	"github.com/welltop-cn/common/utils/env"
	"github.com/welltop-cn/common/utils/log"
)

func Init(serviceName string, ops ...KitOptions) error {
	kitOps := &kitOptions{
		serviceName: serviceName,
	}
	log.InitLogger()

	env.SetServiceName(kitOps.serviceName)
	//	metrics.InitMetrics(kitOps.serviceName)

	config.InitConfig()

	for _, o := range ops {
		o(kitOps)
	}
	// 监听退出信号
	go kitOps.waitingShutdown()
	return nil
}
