package router

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"ken-gateway/internal/ken-gateway/consul"
	"ken-gateway/internal/ken-gateway/middlewares"
	"ken-gateway/internal/pkg/common"
	"net/http"
	"strings"
)

type ConsulConfiguration struct {
	Consul struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"consul"`
}

var LocalInstance *consul.LocalServiceRegister

func InitRouter() *gin.Engine {
	r := gin.New()

	// 向注册中心注册网关服务
	config := &common.Config[ConsulConfiguration]{FilePath: "config/consul.yaml"}
	consulServer := new(ConsulConfiguration)

	err := config.SetYaml(consulServer)
	if err != nil {
		panic(err)
	}

	LocalInstance = consul.NewLocalServiceRegister(consulServer.Consul.Host, consulServer.Consul.Port)
	instance, _ := consul.NewServiceInstance()
	err = LocalInstance.Register(instance)
	if err != nil {
		panic(err)
	}
	// 关闭该服务
	defer LocalInstance.DeRegister()
	// 路由中间件，权限校验
	r.Use(middlewares.Auth(LocalInstance))

	// 注册路由
	registerRouter(r)

	return r
}

func health(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "ok",
	})
}

func registerRouter(e *gin.Engine) {
	e.GET("/ping", health)
	g := e.Group("/api")
	g.Any("/*proxyPath", proxy)
}

// proxy
// ctx.Request.URL 表示当前访问的全路由
// 规则为 /api/api版本/服务名/
func proxy(ctx *gin.Context) {
	// 解析路由信息
	path := ctx.Request.RequestURI

	paths := strings.Split(path, "/")

	version, serverName := paths[2], paths[3]

	requestPath := strings.Join(paths[3:], "/")

	data := make(map[string]any)

	err := ctx.ShouldBind(&data)
	if err != nil {
		fmt.Println("------>", err)
	}

	// 获取服务实例
	serverName = fmt.Sprintf("%s_%s", "kensrail", strings.ToLower(serverName))
	fmt.Println(serverName)
	instance := LocalInstance.GetService(serverName)
	fmt.Println(instance)
	if instance == nil {
		ctx.JSON(500, gin.H{
			"message": "服务不存在",
		})
		return
	}

	// TODO 根据实例返回的信息构造请求
	url := fmt.Sprintf("http://%s:%s/api/%s/%s", instance.GetHost(), instance.GetPort(), version, requestPath)

	// TODO 调用服务 并返回结果
	client := http.Client{}
	req, err := http.NewRequest(ctx.Request.Method, url, ctx.Request.Body)
	if err != nil {
		ctx.JSON(500, gin.H{
			"message": "服务不存在",
		})
		return
	}
	req.Header = ctx.Request.Header
	resp, err := client.Do(req)
	if err != nil {
		ctx.JSON(500, gin.H{
			"message": "服务不存在",
		})
		return
	}
	body := make(map[string]any)
	bytesBody, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(bytesBody, &body)
	ctx.JSON(resp.StatusCode, gin.H{
		"data": body,
	})
}

// TODO 根据实例返回的信息构造请求
