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
			Password:         "df ",
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

func TestUpdateUserInfo() {
	_, err := userClient.UpdateUser(context.Background(), &proto.UpdateUserInfo{
		Id:       11,
		NickName: "lzz1222",
		Gender:   "2",
	})
	if err != nil {
		panic(err)
	} else {
		fmt.Println("更新成功")
	}
}

//测试通过手机号进行查找用户
func TestByMobileGetUserinfo() {
	mobile, err := userClient.GetUserByMobile(context.Background(), &proto.MobileRequest{
		Mobile: "19182031524",
	})
	if err != nil {
		panic(err)
	} else {
		fmt.Println(mobile)
	}
}

//通过id进行用户的查找
func TestByIdGetUserinfo() {
	id, err := userClient.GetUserById(context.Background(), &proto.IdRequest{Id: 11})
	if err != nil {
		panic(err)
	}
	fmt.Println(id)
}

/**
测试全部通过
*/
func main() {
	Init()
	//返回用户列表和检查密码接口
	//GetUserList()
	//测试创建用户
	//TestCreatUser()
	//测试更新用户
	//TestUpdateUserInfo()
	//测试通过mobile进行查找
	//TestByMobileGetUserinfo()
	//测试使用id进行用户查找
	//TestByIdGetUserinfo()
	err := conn.Close()
	if err != nil {
		return
	}
}
