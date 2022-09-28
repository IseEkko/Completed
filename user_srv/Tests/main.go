package main

import (
	"Completed_Server/user_srv/proto"
	"context"
	"fmt"
	"google.golang.org/grpc"
)

var userClient proto.UserClient
var conn *grpc.ClientConn

func Init() {
	var err error
	conn, err = grpc.Dial("127.0.0.1:50051", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	userClient = proto.NewUserClient(conn)
}

func GetUserList() {
	list, err := userClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    1,
		PSize: 10,
	})
	if err != nil {
		return
	}
	for _, v := range list.Data {
		password, err := userClient.CheckPassword(context.Background(), &proto.PasswordCheckInfo{
			Password:         "generic password",
			EncrytedPassword: v.Password,
		})
		if err != nil {
			return
		}
		if password.Success {
			fmt.Println("用户登录成功" + v.NickName)
		} else {
			panic("失败")
		}
	}
}
func TestCreatUser() {
	for i := 0; i < 10; i++ {
		_, err := userClient.CreateUser(context.Background(), &proto.CreateUserinfo{
			NickName: fmt.Sprintf("ceshi%d", i),
			PassWord: "ceshi123",
			Mobile:   fmt.Sprintf("%d", i),
		})
		if err != nil {
			panic(err)
		}
		fmt.Println("创建成功")
	}

}
func main() {
	Init()
	//返回用户列表和检查密码接口
	//GetUserList()
	//测试创建用户
	TestCreatUser()
	conn.Close()

}
