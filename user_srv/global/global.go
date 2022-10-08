package global

import (
	"Completed_Server/user_srv/Config"

	"gorm.io/gorm"
)

var (
	//mysql数据可的链接对象
	DB *gorm.DB
	//配置文件
	SerConfig *Config.ServerConfig = &Config.ServerConfig{}
)

//func init() {
//	dsn := "root:@tcp(127.0.0.1:3306)/completed_user_srv?charset=utf8mb4&parseTime=True&loc=Local"
//	newLogger := logger.New(
//		log.New(os.Stdout, "\r\n", log.LstdFlags),
//		logger.Config{
//			SlowThreshold: time.Second,
//			Colorful:      true,
//			LogLevel:      logger.Info,
//		})
//	//全局模式
//	var err error
//	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
//		Logger: newLogger,
//	})
//	if err != nil {
//		panic(err)
//	}
//}
