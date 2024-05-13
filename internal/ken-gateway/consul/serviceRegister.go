package consul

import (
	"errors"
	capi "github.com/hashicorp/consul/api"
	"strconv"
)

type ServiceRegister interface {
	Register(serviceInstance *Instance) error
	DeRegister() error
}

type LocalServiceRegister struct {
	Instances       map[string]*Instance
	Client          *capi.Client
	CurrentInstance *Instance
}

// NewLocalServiceRegister 创建本地实例
// @host: 服务注册中心地址
// @port: 服务注册中心端口
func NewLocalServiceRegister(host string, port int) *LocalServiceRegister {
	config := capi.DefaultConfig()
	config.Address = host + ":" + strconv.Itoa(port)

	client, err := capi.NewClient(config)
	if err != nil {
		panic(err)
	}

	return &LocalServiceRegister{
		Instances:       make(map[string]*Instance),
		Client:          client,
		CurrentInstance: nil,
	}
}

// Register 向consul 注册中心注册服务
// @instance: 服务实例
func (l *LocalServiceRegister) Register(instance *Instance) error {
	register := new(capi.AgentServiceRegistration)
	register.Name = instance.Service.ServiceId

	// 配置实例心跳
	register.Check = &capi.AgentServiceCheck{
		HTTP:                           instance.Check.Http,
		Timeout:                        instance.Check.Timeout,
		Interval:                       instance.Check.Interval,
		DeregisterCriticalServiceAfter: instance.Check.DeregisterCriticalServiceAfter,
	}

	err := l.Client.Agent().ServiceRegister(register)
	if err != nil {
		return err
	}
	l.CurrentInstance = instance
	l.Instances[instance.Service.ServiceId] = instance
	return nil
}

// DeRegister 向consul 注册中心注销服务
func (l *LocalServiceRegister) DeRegister() error {
	if l.CurrentInstance == nil {
		return errors.New("当前服务不存在")
	}

	err := l.Client.Agent().ServiceDeregister(l.CurrentInstance.Service.InstanceId)
	if err != nil {
		return err
	}
	l.CurrentInstance = nil
	return nil
}
