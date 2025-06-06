package config

import (
	"fmt"
	"path"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/welltop-cn/common/utils/env"
)

type Config struct {
}

func loadLocalConfig() config.Config {
	rootPath, err := env.GetProjectPath()
	if err != nil {
		panic("get root path " + err.Error())
	}
	configPath := path.Join(rootPath, fmt.Sprintf("configs/%s", env.GetRuntimeEnv()))

	yamlResource := fmt.Sprintf("%s/config.yaml", configPath)

	c := config.New(config.WithSource(file.NewSource(yamlResource)))
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	return c
}

func InitConfig() {
	if env.IsDebug() {
		conf := loadLocalConfig()
		config.SetConfig(conf)
		return
	}
}
