package dao

import (
	"github.com/SherryProgrammer/SherryGateway/dto"
	"github.com/SherryProgrammer/go_evnconfig/lib"
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	"sync"
)

type ServiceDetail struct {
	Info          *ServiceInfo   `json:"info" description:"基本信息"`
	HTTPRule      *HttpRule      `json:"http_rule" description:"http_rule"`
	TCPRule       *TcpRule       `json:"tcp_rule" description:"tcp_rule"`
	GRPCRule      *GrpcRule      `json:"grpc_rule" description:"grpc_rule"`
	LoadBalance   *LoadBalance   `json:"load_balance" description:"load_balance"`
	AccessControl *AccessControl `json:"access_control" description:"access_control"`
}

var ServiceManagerHandler *ServiceManager

func init() {
	ServiceManagerHandler = NewServiceManager()
}

type ServiceManager struct {
	ServiceMap   map[string]*ServiceDetail
	ServiceSlice []*ServiceDetail //遍历
	Locker       sync.RWMutex
	init         sync.Once
	err          error
}

func NewServiceManager() *ServiceManager {
	return &ServiceManager{
		ServiceMap:   map[string]*ServiceDetail{},
		ServiceSlice: []*ServiceDetail{},
		Locker:       sync.RWMutex{},
		init:         sync.Once{},
	}
}

func (s *ServiceManager) HTTPAccessMode(c *gin.Context) (*ServiceDetail, error) {
	//1,前缀匹配/abc ==> service.rule
	//2,前缀匹配 www.test.com ==> service.rule

	return nil, nil
}

func (s *ServiceManager) LoadOnce() error {
	s.init.Do(func() {
		//从db中分页读取基本
		serviceInfo := &ServiceInfo{}
		//从数据库中读取adninInfo
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		tx, err := lib.GetGormPool("default")
		if err != nil {
			s.err = err
			return
		}
		params := &dto.ServiceListInput{PageNo: 1, PageSize: 99999}
		list, _, err := serviceInfo.PageList(c, tx, params)
		if err != nil {
			s.err = err
			return
		}
		s.Locker.Lock()
		defer s.Locker.Unlock()
		for _, listItem := range list {
			serviceDetail, err := listItem.ServiceDetail(c, tx, &listItem)
			if err != nil {
				s.err = err
				return
			}
			s.ServiceMap[listItem.ServiceName] = serviceDetail
			s.ServiceSlice = append(s.ServiceSlice, serviceDetail)
		}
	})
	return s.err
}
