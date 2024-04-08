package controller

import (
	"github.com/SherryProgrammer/SherryGateway/dao"
	"github.com/SherryProgrammer/SherryGateway/dto"
	"github.com/SherryProgrammer/SherryGateway/middleware"
	"github.com/SherryProgrammer/SherryGateway/public"
	"github.com/SherryProgrammer/go_evnconfig/lib"
	"github.com/gin-gonic/gin"
)

type ServiceController struct {
}

func ServiceRegister(group *gin.RouterGroup) {
	service := &ServiceController{}
	group.GET("/service_list", service.ServiceList)
}

// ServiceList godoc
// @Summary 服务列表
// @Description 服务列表
// @Tags 服务管理
// @ID /service/service_list
// @Accept  json
// @Produce  json
// @Param info query string true "关键词"
// @Param page_size query int true "每页个数"
// @Param page_no query int true "当前页数"
// @Success 200 {object} middleware.Response{data=dto.ServiceListOutput} "success"
// @Router /service/service_list [get]
func (service *ServiceController) ServiceList(c *gin.Context) {
	params := &dto.ServiceListInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}
	//从数据库中读取adninInfo
	tx, err := lib.GetGormPool("default")

	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	serviceInfo := &dao.ServiceInfo{}
	list, total, err := serviceInfo.PageList(c, tx, params)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	outlist := []dto.ServiceListItemOutput{}
	for _, listItem := range list {
		sericeDetail, err := listItem.ServiceDetail(c, tx, &listItem)
		if err != nil {
			middleware.ResponseError(c, 2003, err)
			return
		}
		serviceAddr := ""
		clusterIP := lib.GetStringConf("base.cluster.cluster_ip")
		clusterPort := lib.GetStringConf("base.cluster.cluster_port")
		clusterSSLPort := lib.GetStringConf("base.cluster.cluster_ssl_port")

		if sericeDetail.Info.LoadType == public.LoadTypeHTTP && sericeDetail.HTTP.Rule == public.HTTPRuleTypPrefixURL {
			serviceAddr = clusterIP + clusterPort + sericeDetail.HTTP.Rule
		}

		//http后缀接入1.clusterIP+clusterPort+path
		//http域名接入2.domain
		//tcp\grpc接入3.clusterIP+servicePort

		outItem := dto.ServiceListItemOutput{
			ID:          listItem.ID,
			ServiceName: listItem.ServiceName,
			ServiceDesc: listItem.ServiceDesc,
		}
		outlist = append(outlist, outItem)
	}

	out := &dto.ServiceListOutput{
		Total: total,
		List:  outlist,
	} //返回结构体
	middleware.ResponseSuccess(c, out)
}
