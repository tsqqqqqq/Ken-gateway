package consul

import "strconv"

type ServiceDiscover interface {
	GetServices() ([]string, error)

	GetService(serviceId string) *Instance
}

func (l *LocalServiceRegister) GetServices() ([]string, error) {
	services, err := l.Client.Agent().Services()
	if err != nil {
		return nil, err
	}
	result := make([]string, len(services))
	for _, service := range services {
		if _, ok := l.Instances[service.Service]; ok {
			continue
		}
		l.Instances[service.Service] = &Instance{
			Service: Service{
				ServiceId:  service.Service,
				InstanceId: service.ID,
				Host:       service.Address,
				Port:       strconv.Itoa(service.Port),
				MetaData:   service.Meta,
			},
		}
		result = append(result, service.Service)
	}
	return result, nil
}

func (l *LocalServiceRegister) GetService(serviceId string) *Instance {
	instance, ok := l.Instances[serviceId]
	if ok {
		return instance
	}

	service, _, err := l.Client.Agent().Service(serviceId, nil)
	if err != nil {
		return nil
	}
	instance = &Instance{
		Service: Service{
			ServiceId:  service.Service,
			InstanceId: service.ID,
			Host:       service.Address,
			Port:       strconv.Itoa(service.Port),
			MetaData:   service.Meta,
		},
	}
	l.Instances[serviceId] = instance

	return instance
}
