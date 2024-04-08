package dao

import (
	"errors"
	"time"

	"github.com/SherryProgrammer/SherryGateway/dto"
	"github.com/SherryProgrammer/SherryGateway/public"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Admin struct {
	Id       int       `json:"id" orm:"column(id);auto" description:"自增主键"`
	UserName string    `json:"user_name" orm:"column(user_name);size(191)" description:"管理员用户名"`
	Salt     string    `json:"salt" orm:"column(salt);size(191)" description:"盐"`
	Password string    `json:"Password" orm:"column(Password);size(191)" description:"密码"`
	CityId   int       `json:"city_id" orm:"column(city_id)" description:"城市id"`
	UserId   int64     `json:"user_id" orm:"column(user_id)" description:"操作人"`
	UpdateAt time.Time `json:"update_at" orm:"column(update_at);type(datetime)" description:"更新时间"`
	CreateAt time.Time `json:"create_at" orm:"column(create_at);type(datetime)" description:"创建时间"`
	IsDelete int       `json:"is_delete" orm:"column(is_delete);type(datetime)" description:"是否删除"`
}

func (t *Admin) TableName() string {
	return "gateway_admin"
}

func (t *Admin) LoginCheck(c *gin.Context, tx *gorm.DB, param *dto.AdminLoginInput) (*Admin, error) {

	adminInfo, err := t.Find(c, tx, (&Admin{UserName: param.UserName, IsDelete: 0}))
	if err != nil {
		return nil, errors.New("用户信息不存在")
	}
	//param.Password
	//adminInfo.Salt
	saltPassword := public.GenSaltPassword(adminInfo.Salt, param.Password)
	if adminInfo.Password != saltPassword {
		return nil, errors.New("密码错误，请重新输入")
	}
	return adminInfo, nil
}

// 定义一个查询的方法
// 1.params.username 取得管理员信息 admininfo
func (t *Admin) Find(c *gin.Context, tx *gorm.DB, search *Admin) (*Admin, error) {
	out := &Admin{}
	//不用那样封装
	err := tx.WithContext(c).Where(search).Find(out).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (t *Admin) Save(c *gin.Context, tx *gorm.DB) error {

	return tx.WithContext(c).Save(t).Error

}
