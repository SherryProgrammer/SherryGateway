package controller

import (
	"encoding/json"
	"github.com/SherryProgrammer/SherryGateway/dao"
	"github.com/SherryProgrammer/SherryGateway/dto"
	"github.com/SherryProgrammer/SherryGateway/middleware"
	"github.com/SherryProgrammer/SherryGateway/public"
	"github.com/SherryProgrammer/go_evnconfig/lib"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"time"
)

type AdminLoginController struct {
}

func AdminLoginRegister(group *gin.RouterGroup) {
	adminLogin := &AdminLoginController{}
	group.POST("/login", adminLogin.AdminLogin)
	//登录的用POST更好
	group.GET("/logout", adminLogin.AdminLoginOut)
}

//action

// AdminLogin godoc
// @Summary 管理员登录
// @Description 管理员登录
// @Tags 管理员接口
// @ID /admin_login/login
// @Accept  json
// @Produce  json
// @Param polygon body dto.AdminLoginInput true "body"
// @Success 200 {object} middleware.Response{data=dto.AdminLoginOutput} "success"
// @Router /admin_login/login [post]
func (adminlogin *AdminLoginController) AdminLogin(c *gin.Context) {
	//定义一个用户名和密码的结构体 用json去请求
	//怎么转化 首先定义一个dto的结构体 admin_login
	params := &dto.AdminLoginInput{}
	//校验
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 1001, err)
		return
	}
	//1.params.username 取得管理员信息 admininfo
	//2.adminifo.salt+params.Password sha256 =>saltPassword
	//3.saltPassword==admininfo.password

	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	admin := &dao.Admin{}
	admin, err = admin.LoginCheck(c, tx, params)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	//设置session
	sessInfo := &dto.AdminSessionInfo{
		ID:        admin.Id,
		UserName:  admin.UserName,
		LoginTime: time.Now(),
	}
	sessBts, err := json.Marshal(sessInfo)
	if err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}

	sess := sessions.Default(c)
	sess.Set(public.AdminSessionInfoKey, string(sessBts))
	sess.Save()

	out := &dto.AdminLoginOutput{Token: params.UserName} //返回结构体
	middleware.ResponseSuccess(c, out)
}

// AdminLoginOut godoc
// @Summary 管理员退出
// @Description 管理员退出
// @Tags 管理员接口
// @ID /admin_login/logout
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /admin_login/logout [get]
func (adminlogin *AdminLoginController) AdminLoginOut(c *gin.Context) {

	sess := sessions.Default(c)
	sess.Delete(public.AdminSessionInfoKey)
	sess.Save()
	middleware.ResponseSuccess(c, "")
}
