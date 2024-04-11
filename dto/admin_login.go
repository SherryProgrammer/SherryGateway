package dto

import (
	"github.com/SherryProgrammer/SherryGateway/public"
	"github.com/gin-gonic/gin"
	"time"
)

// 输入参数
type AdminSessionInfo struct {
	ID        int       `json:"id"`
	UserName  string    `json:"username"`
	LoginTime time.Time `json:"login_time"`
}

type AdminLoginInput struct {
	//需要设置几类tag example可以在swagger文档里面生成一个默认值 comment错误输出的时候直接用这个名字去输出 validate校验是否一定要设置
	//json是指输出非时候用结构体转换成一个json的形式 form 是指输入的时候name 里面的值 json转化成结构体的形式
	//is_valid_username translation tag
	UserName string `json:"username" form:"username" comment:"姓名" example:"admin" validate:"required,valid_username"` //管理员用户名
	Password string `json:"password" form:"password" comment:"密码" example:"123456" validate:"required"`               //密码
}

// 参数
func (param *AdminLoginInput) BindValidParam(c *gin.Context) error {
	//绑定成功输出一个 new 赋值 不成功输出一个error
	//校验过程
	return public.DefaultGetValidParams(c, param)
}

type AdminLoginOutput struct {
	//需要设置几类tag example可以在swagger文档里面生成一个默认值 comment错误输出的时候直接用这个名字去输出 validate校验是否一定要设置
	//json是指输出非时候用结构体转换成一个json的形式 form 是指输入的时候name 里面的值 json转化成结构体的形式
	//因为是输出不需要验证 不用任何信息
	Token string `json:"token" form:"token" comment:"token" example:"token" validate:""` //Token
}
