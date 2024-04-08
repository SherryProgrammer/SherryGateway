package dto

import (
	"github.com/SherryProgrammer/SherryGateway/public"
	"github.com/gin-gonic/gin"
	"time"
)

type AdminInfoOutput struct {
	ID           int       `json:"id"`
	UserName     string    `json:"username"`
	LoginTime    time.Time `json:"login_time"`
	Avatar       string    `json:"avatar"`
	Introduction string    `json:"introduction"`
	Roles        []string  `json:"roles"`
}

type ChangePwdInput struct {
	Password string `json:"password" form:"password" comment:"密码" example:"123456" validate:"required"` //密码
}

func (param *ChangePwdInput) BindValidParam(c *gin.Context) error {
	//绑定成功输出一个 new 赋值 不成功输出一个error
	//校验过程
	return public.DefaultGetValidParams(c, param)
}
