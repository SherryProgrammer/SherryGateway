package dto

import (
	"github.com/SherryProgrammer/SherryGateway/public"
	"github.com/gin-gonic/gin"
)

type TokensInput struct {
	//需要设置几类tag example可以在swagger文档里面生成一个默认值 comment错误输出的时候直接用这个名字去输出 validate校验是否一定要设置
	//json是指输出非时候用结构体转换成一个json的形式 form 是指输入的时候name 里面的值 json转化成结构体的形式
	//is_valid_username translation tag
	GrantType string `json:"grant_type" form:"grant_type" comment:"授权类型" example:"admin" validate:"required,client_credentials"` //授权类型
	Scope     string `json:"scope" form:"scope" comment:"权限范围" example:"read_write" validate:"required"`                         //权限范围
}

// 参数
func (param *TokensInput) BindValidParam(c *gin.Context) error {
	//绑定成功输出一个 new 赋值 不成功输出一个error
	//校验过程
	return public.DefaultGetValidParams(c, param)
}

type TokensOutput struct {
	AccessToken string `json:"access_token" form:"access_token" ` //授权类型
	ExpiresIn   int    `json:"expires_in" form:"expires_in" `     //授权类型
	TokenType   string `json:"token-type" form:"token-type" `     //授权类型
	Scope       string `json:"scope" form:"access_token" `        //授权类型
}

// 参数
func (param *TokensOutput) BindValidParam(c *gin.Context) error {
	//绑定成功输出一个 new 赋值 不成功输出一个error
	//校验过程
	return public.DefaultGetValidParams(c, param)
}
