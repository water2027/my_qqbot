package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"qqbot/dto"
	"qqbot/utils"
	"strconv"
	"time"
)

type AuthService struct {
	Token        string        `json:"-"`
	expire       time.Duration `json:"-"`
	AppId        string        `json:"appId"`
	ClientSecret string        `json:"clientSecret"`
}

func NewAuthService() *AuthService {
	appId := os.Getenv("APP_ID")
	clientSecret := os.Getenv("APP_SECRET")
	return &AuthService{
		Token:        "",
		expire:       0,
		AppId:        appId,
		ClientSecret: clientSecret,
	}
}

var AuthHelper *AuthService

func (authService *AuthService) isValid() bool {
	return authService.Token != "" && authService.expire.Seconds() < 3 // 如果小于三秒的话，那等三秒再刷新token
}

func (authService *AuthService) GetToken() string {
	if authService.isValid() {
		return authService.Token
	}
	err := authService.refreshToken()
	if err != nil {
		fmt.Println(err.Error())
		time.Sleep(5 * time.Second)
		return authService.GetToken()
	}
	return fmt.Sprintf("QQBot %s", authService.Token)
}

func (authService *AuthService) refreshToken() error {
	resp, err := utils.NetHelper.POST("https://bots.qq.com/app/getAppAccessToken", authService)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(data))

	var result dto.AuthResponse
	err = json.Unmarshal(data, &result)
	if err != nil {
		return err
	}
	authService.Token = result.AccessToken
	exp, err := strconv.ParseUint(result.ExpiresIn, 10, 64)
	if err != nil {
		fmt.Println(result.ExpiresIn)
		return err
	}

	authService.expire = time.Duration(exp * uint64(time.Second))

	return nil
}

func init() {
	AuthHelper = NewAuthService()
	// 先获取一次token，避免第一次请求的时候需要等待
	_ = AuthHelper.refreshToken()
	fmt.Println("初始化完成")
}