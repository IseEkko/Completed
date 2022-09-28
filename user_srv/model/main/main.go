package main

import (
	"Completed_Server/user_srv/model"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

/***
我们如果直接使用md5自带的返回的是一个16大小的byte数组，
返回的时候我们需要的是string，所以使用的是
          _, _ = io.WriteString(Md5, code)
	      return hex.EncodeToString(Md5.Sum(nil))
只用纯的md5进行加密也是有缺陷的，彩虹表问题，也就是拿着已有的加密然后去做对比

解决这个问题的方式：
   加上salt进行解决
*/
func ginMd5(code string) string {
	Md5 := md5.New()
	_, _ = io.WriteString(Md5, code)
	return hex.EncodeToString(Md5.Sum(nil))
}

/***
用于同步数据表结构
*/
func main() {
	dsn := "root:@tcp(127.0.0.1:3306)/completed_user_srv?charset=utf8mb4&parseTime=True&loc=Local"
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			Colorful:      true,
			LogLevel:      logger.Info,
		})
	//全局模式
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}
	options := &password.Options{16, 100, 32, md5.New}
	/***
	增加安全性
	使用三段拼接
	*/
	salt, encodedPwd := password.Encode("generic password", options)
	//这里使用的是$进行分割，对后面验证的时候使用$进行取出
	newPassword := fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)
	for i := 0; i < 10; i++ {
		user := model.User{
			Mobile:   fmt.Sprintf("1918203152%d", i),
			Password: newPassword,
			NickName: fmt.Sprintf("bobby%d"),
		}
		db.Save(&user)
	}

	//_ = db.AutoMigrate(&model.User{})
}

/***
这里是验证盐值加密的手法
*/
func Md5Salt() {
	//这里返回的是两个值，第一个是盐值，第二个是生成之后的密码
	/***
	这里我们使用的包是：
	"github.com/anaskhan96/go-password-encoder"
	使用这个包的好处是，可以使用：password.Verify进行验证密码
	最后会返回一个bool值
	*/
	//salt, encodedPwd := password.Encode("generic password", nil)
	//check := password.Verify("generic password", salt, encodedPwd, nil)
	//fmt.Println(check) // true

	// Using custom options
	salt, encodedPwd := password.Encode("generic password", nil)
	options := &password.Options{16, 100, 32, md5.New}
	/***
	增加安全性
	使用三段拼接
	*/
	salt, encodedPwd = password.Encode("generic password", options)
	//这里使用的是$进行分割，对后面验证的时候使用$进行取出
	newPassword := fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)
	//这个newPassword就是最后存入数据库的部分
	fmt.Println(len(newPassword))
	//使用Strings进行分割
	passwordinfo := strings.Split(newPassword, "$")
	check := password.Verify("generic password", passwordinfo[2], passwordinfo[3], options)
	fmt.Println(check) // true
}
