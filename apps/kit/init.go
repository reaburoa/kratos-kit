package kit

func Init(serviceName string, ops ...KitOptions) error {
	kitOps := &kitOptions{
		serviceName: serviceName,
	}
	clog.InitLogger()

	env.SetServiceName(kitOps.serviceName)
	metrics.InitMetrics(kitOps.serviceName)

	config.InitConfig()

	for _, o := range ops {
		o(kitOps)
	}
	// 监听退出信号
	go kitOps.waitingShutdown()
	return nil
}
