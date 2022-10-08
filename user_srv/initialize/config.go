package initialize

import (
	"Completed_Server/user_srv/global"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}

func InitConfig() {
	Viper := viper.New()
	Viper.SetConfigFile("config.yaml")
	err := Viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	/**
	  使用结构体进行接收
	*/
	if err := Viper.Unmarshal(global.SerConfig); err != nil {
		panic(err)
	}
	zap.S().Infof("Config init success: %s", global.SerConfig)
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file channed:", e.Name)
		_ = Viper.ReadInConfig()
		_ = Viper.Unmarshal(global.SerConfig)
		fmt.Println(global.SerConfig)
	})
}
