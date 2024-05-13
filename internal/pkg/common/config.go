package common

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Configuration[T any] interface {
	SetYaml(obj T)
}

type Config[T any] struct {
	FilePath string
}

// SetYaml 读取yaml文件并赋值
func (c *Config[T]) SetYaml(obj *T) error {
	data, err := os.ReadFile(c.FilePath)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, &obj)
	return err
}
