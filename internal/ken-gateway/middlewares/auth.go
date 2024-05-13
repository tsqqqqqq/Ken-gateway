package middlewares

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"ken-gateway/internal/ken-gateway/consul"
	"net/http"
	"slices"
	"strings"
)

var whiteList = []string{
	"/ping",
}

const superAdmin = "5fd724eeea7d3d5746f036d166715a8e"
const salt = "ken"

func Auth(localInstance *consul.LocalServiceRegister) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 白名单过滤
		if slices.Contains(whiteList, c.Request.RequestURI) {
			c.Next()
			return
		}
		// 判断header中有没有token
		user, token := c.Request.Header.Get("X-USERID"), c.Request.Header.Get("X-TOKEN")
		// 判断是否是超级管理员
		if IsSuperAdmin(token) {
			c.Next()
			return
		}
		if user == "" {
			// 获取登录url
			url, err := loginUrl(localInstance, c.Request.RequestURI)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"message": err,
				})
			}
			// TODO user 跳转登录接口
			c.Redirect(http.StatusTemporaryRedirect, url)
			c.Abort()
		}

		// TODO token 校验token是否过期
		if token == "" {
			// TODO token 跳转登录接口
		}

		c.Next()
	}

}

func loginUrl(localInstance *consul.LocalServiceRegister, uri string) (string, error) {
	params := strings.Split(uri, "?")[1:]
	query := strings.Join(params, "&")

	// 获取登录服务
	serverName := fmt.Sprintf("%s_%s", "kensrail", "login")
	instance := localInstance.GetService(serverName)
	if instance == nil {
		return "", errors.New("服务不存在")
	}
	url := fmt.Sprintf("http://%s:%s/api/%s/%s", instance.GetHost(), instance.GetPort(), "v1", "login"+"?"+query)
	return url, nil
}

func IsSuperAdmin(token string) bool {
	// TODO 校验token
	h := md5.New()
	io.WriteString(h, token+salt)
	auth := fmt.Sprintf("%x", h.Sum(nil))
	return superAdmin == auth
}
