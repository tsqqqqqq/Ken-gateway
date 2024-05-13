package consul

import (
	"ken-gateway/internal/pkg/common"
	"math/rand"
	"strconv"
	"time"
)

type ServiceInstance interface {
	GetInstanceId() string

	GetServiceId() string

	GetHost() string

	GetPort() string

	IsSecure() bool

	GetMetaData() map[string]string
}

type Service struct {
	InstanceId string            `json:"instanceId,omitempty" yaml:"instanceId"`
	ServiceId  string            `json:"serviceId,omitempty" yaml:"serviceId"`
	Host       string            `json:"host" json:"host,omitempty" yaml:"host"`
	Port       string            `json:"port" json:"port,omitempty" yaml:"port"`
	Secure     bool              `json:"secure,omitempty" yaml:"secure"`
	MetaData   map[string]string `json:"metaData,omitempty" yaml:"metaData"`
}

type Check struct {
	Http                           string `json:"http,omitempty" yaml:"http"`
	Timeout                        string `json:"timeout,omitempty" yaml:"timeout"`
	Interval                       string `json:"interval,omitempty" yaml:"interval"`
	DeregisterCriticalServiceAfter string `json:"deregisterCriticalServiceAfter,omitempty" yaml:"deregisterCriticalServiceAfter"`
}

type Instance struct {
	Service `json:"service" yaml:"service"`
	Check   `json:"check" yaml:"check"`
}

// NewServiceInstance 创建服务实例
func NewServiceInstance() (*Instance, error) {
	config := &common.Config[Instance]{FilePath: "config/consul.yaml"}
	instance := new(Instance)

	err := config.SetYaml(instance)

	instance.Service.InstanceId = instance.Service.ServiceId + "-" + strconv.FormatInt(time.Now().Unix(), 10) + "-" + strconv.Itoa(rand.Intn(9000)+1000)

	return instance, err
}

func (s *Service) GetInstanceId() string {
	return s.InstanceId
}

func (s *Service) GetServiceId() string {
	return s.ServiceId
}

func (s *Service) GetHost() string {
	return s.Host
}

func (s *Service) GetPort() string {
	return s.Port
}

func (s *Service) IsSecure() bool {
	return s.Secure
}

func (s *Service) GetMetaData() map[string]string {
	return s.MetaData
}
