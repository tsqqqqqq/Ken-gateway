package test

import (
	"ken-gateway/internal/ken-gateway/consul"
	"ken-gateway/internal/pkg/common"
	"testing"
)

func TestSetYaml(t *testing.T) {
	config := &common.Config[consul.Service]{FilePath: "../config/consul.yaml"}
	service := new(consul.Service)

	err := config.SetYaml(service)

	if service.ServiceId != "ken-gateway" || err != nil {
		t.Fatalf("文件赋值错误 失败原因：%s", err)
	}
}
