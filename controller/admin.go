package controller

import (
	"encoding/json"
	"fmt"
	"github.com/SherryProgrammer/go_evnconfig/lib"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/SherryProgrammer/SherryGateway/dao"
	"github.com/SherryProgrammer/SherryGateway/dto"
	"github.com/SherryProgrammer/SherryGateway/middleware"
	"github.com/SherryProgrammer/SherryGateway/public"
)

// 官方 三方 本地
type AdminController struct {
}

func AdminRegister(group *gin.RouterGroup) {
	adminLogin := &AdminController{}
	group.GET("/admin_info", adminLogin.AdminInfo)
	//更新的用POST更好
	group.POST("/change_pwd", adminLogin.ChangePwd)

}

//action

// AdminInfo godoc
// @Summary 管理员信息
// @Description 管理员信息
// @Tags 管理员接口
// @ID /admin/admin_info
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=dto.AdminInfoOutput} "success"
// @Router /admin/admin_info [get]
func (adminlogin *AdminController) AdminInfo(c *gin.Context) {
	sess := sessions.Default(c)
	sessInfo := sess.Get(public.AdminSessionInfoKey)
	//sessInfoStr := sessInfo.(string)
	adminSessionInfo := &dto.AdminSessionInfo{}
	if err := json.Unmarshal([]byte(fmt.Sprint(sessInfo)), adminSessionInfo); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	} //拿到信息

	//1.读取sessionkey对应的json转化为结构体
	//2.取出数据然后封装输出结构体

	//Avatar       string    `json:"avatar"`
	//Introduction string    `json:"introduction"`
	//Roles        []string  `json:"roles"`
	out := &dto.AdminInfoOutput{
		ID:           adminSessionInfo.ID,
		UserName:     adminSessionInfo.UserName,
		LoginTime:    adminSessionInfo.LoginTime,
		Avatar:       "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif",
		Introduction: "I am a super administrator",
		Roles:        []string{"admin"},
	} //返回结构体
	middleware.ResponseSuccess(c, out)
}

// ChangePwd godoc
// @Summary 修改密码
// @Description 修改密码
// @Tags 管理员接口
// @ID /admin/change_pwd
// @Accept  json
// @Produce  json
// @Param polygon body dto.ChangePwdInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /admin/change_pwd [post]
func (changePwd *AdminController) ChangePwd(c *gin.Context) {
	params := &dto.ChangePwdInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	//1.session读取用户信息到结构体 sessInfo
	//2.sessInfo.ID读取数据库信息 adminInfo
	//3.params.password+adminInfo.salt sah256 saltPassword
	//4.saltPassword==>adminInfo.password 执行数据保存

	//读取用户信息到结构体
	sess := sessions.Default(c)
	sessInfo := sess.Get(public.AdminSessionInfoKey)
	//sessInfoStr := sessInfo.(string)
	adminSessionInfo := &dto.AdminSessionInfo{}
	if err := json.Unmarshal([]byte(fmt.Sprint(sessInfo)), adminSessionInfo); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}
	//从数据库中读取adninInfo
	tx, err := lib.GetGormPool("default")

	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	adminInfo := &dao.Admin{}
	adminInfo, err = adminInfo.Find(c, tx, &dao.Admin{UserName: adminSessionInfo.UserName})
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	//生成新密码
	saltPassword := public.GenSaltPassword(adminInfo.Salt, params.Password)
	adminInfo.Password = saltPassword
	//执行数据保存
	if err := adminInfo.Save(c, tx); err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}

	middleware.ResponseSuccess(c, "")
}
