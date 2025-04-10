package service

import (
	"time"
)

type AuthService struct {
	Token  string
	expire time.Duration
}

func (authService *AuthService) isValid() bool {
	return authService.expire.Seconds() < 3 // 如果小于三秒的话，那等三秒再刷新token 
}

func (authService *AuthService) GetToken() string {
	if authService.isValid() {
		return authService.Token
	}
	// TODO: 获取token的操作
	return ""
}




