package model

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID        int32      `gorm:"primarykey"`
	CreateeAt *time.Time `gorm:"column:add_time"`
	UpdateAt  *time.Time `gorm:"column：update_time"`
	DeletedAt gorm.DeletedAt
	IsDeleted bool
}

//在这里定义的时候发现的问题，在这里写备注的时候，我们需要注意的是，我们comment后面没有冒号，直接一个单引号就可以
type User struct {
	BaseModel
	Mobile   string     `gorm:"index:idx_mobile;unique;type:varchar(11);not null comment '手机号码'"`
	Password string     `gorm:"type:varchar(100);not null comment:'密码'"`
	NickName string     `gorm:"type:varchar(20);comment '昵称'"`
	Birthday *time.Time //这里是指针，这里使用指针是为了保存零值
	Gender   string     `gorm:"column:gender;default:male;type:varchar(6) comment'female 表示女，male表示男'"`
	Role     int        `gorm:"column:role;default:1;type:int comment '1表示普通用户，2表示管理员'"`
}
