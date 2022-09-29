package handler

import (
	"Completed_Server/user_srv/global"
	"Completed_Server/user_srv/model"
	"Completed_Server/user_srv/proto"
	"context"
	"crypto/md5"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"strings"
	"time"
)

type UserServer struct {
}

func ModelToRsponse(user model.User) proto.UserInfoResponse {
	//TODO 用户服务_将查询出来的数据进行转换
	userInfoRsp := proto.UserInfoResponse{
		Id:       user.ID,
		Password: user.Password,
		NickName: user.NickName,
		Gender:   user.Gender,
		Mobile:   user.Mobile,
		Role:     int32(user.Role),
	}
	if user.Birthday != nil {
		userInfoRsp.Birthday = uint64(user.Birthday.Unix())
	}
	return userInfoRsp
}
func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	//TODO 用户服务_分页函数
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func (u *UserServer) GetUserList(ctx context.Context, info *proto.PageInfo) (*proto.UserListRespons, error) {
	//TODO 用户服务_获取用户列表
	var users []model.User
	result := global.DB.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	rsp := &proto.UserListRespons{}
	rsp.Total = int32(result.RowsAffected)

	//对数据进行分页
	global.DB.Scopes(Paginate(int(info.Pn), int(info.PSize))).Find(&users)
	for _, user := range users {
		userInfoRsp := ModelToRsponse(user)
		rsp.Data = append(rsp.Data, &userInfoRsp)
	}
	return rsp, nil
}

func (u *UserServer) GetUserByMobile(ctx context.Context, request *proto.MobileRequest) (*proto.UserInfoResponse, error) {
	//TODO 用户服务_通过手机号进行用户查询
	var user model.User
	result := global.DB.Where(&model.User{Mobile: request.Mobile}).First(&user)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	userInfoRsp := ModelToRsponse(user)
	return &userInfoRsp, nil
}

func (u *UserServer) GetUserById(ctx context.Context, req *proto.IdRequest) (*proto.UserInfoResponse, error) {
	//TODO 用户服务_通过id进行用户查询
	var user model.User
	result := global.DB.First(&user, req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	userInfoRsp := ModelToRsponse(user)
	return &userInfoRsp, nil
}

func (u *UserServer) CreateUser(ctx context.Context, userinfo *proto.CreateUserinfo) (*proto.UserInfoResponse, error) {
	//TODO 用户服务_创建用户
	//首先查询用户是否存在
	var user model.User
	user.Mobile = userinfo.Mobile
	result := global.DB.Where(user).Find(&user)
	if result.RowsAffected == 1 {
		return nil, status.Errorf(codes.AlreadyExists, "用户已存在")
	}
	//没有找到用户那么就进行创建
	user.Mobile = userinfo.Mobile
	user.NickName = userinfo.NickName
	//密码加密
	options := &password.Options{16, 100, 32, md5.New}
	salt, encodedPwd := password.Encode("generic password", options)
	//这里使用的是$进行分割，对后面验证的时候使用$进行取出
	user.Password = fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)
	//进行创建
	result = global.DB.Create(&user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	//将用户信息返回回去
	userinfoRsp := ModelToRsponse(user)
	return &userinfoRsp, nil
}

func (u *UserServer) UpdateUser(ctx context.Context, info *proto.UpdateUserInfo) (*proto.Empty, error) {
	//TODO 用户服务_用户个人中心信息更新
	//首先进行用户查询
	var user model.User
	user.ID = info.Id
	reslut := global.DB.First(&user)
	if reslut.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	//查到用户信息后进行更改
	birthDay := time.Unix(int64(info.Birthday), 0)
	user.NickName = info.NickName
	user.Birthday = &birthDay
	user.Gender = info.Gender
	reslut = global.DB.Save(&user)
	if reslut.Error != nil {
		return nil, status.Errorf(codes.Internal, reslut.Error.Error())
	}
	/**
	这里需要注意不能直接返回nil，要使用对应的结构体
	*/
	return &proto.Empty{}, nil
}

func (u *UserServer) CheckPassword(ctx context.Context, info *proto.PasswordCheckInfo) (*proto.CheckReponse, error) {
	//TODO 用户服务_密码检查
	options := &password.Options{16, 100, 32, md5.New}
	passwordInfo := strings.Split(info.EncrytedPassword, "$")
	check := password.Verify(info.Password, passwordInfo[2], passwordInfo[3], options)
	return &proto.CheckReponse{Success: check}, nil
}
