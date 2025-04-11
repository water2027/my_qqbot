package randomimage

import (
	"io"
	"strings"
	"encoding/json"

	"qqbot/service/message"
	"qqbot/utils"
)

type RandomImageResponse struct {
	Success bool   `json:"success"`
	Type    string `json:"type"`
	Url     string `json:"url"`
}

func getRandomImage(msg *message.Message) error {
	if !msg.CanBeSet() {
		return nil
	}
	if !strings.HasPrefix(msg.GetRawContent(), "/随机图片") {
		return nil
	}
	resp, err := utils.NetHelper.GET("https://api.vvhan.com/api/wallpaper/acg?type=json")
	if err != nil {
		msg.SetContent("获取随机图片失败，请求失败")
		return nil
	}
	defer resp.Body.Close()

	bytesData, err := io.ReadAll(resp.Body)
	if err != nil {
		msg.SetContent("获取随机图片失败，读取数据失败")
		return nil
	}
	var randomImageResponse RandomImageResponse
	err = json.Unmarshal(bytesData, &randomImageResponse)
	if err != nil {
		msg.SetContent("获取随机图片失败，解析json失败")
		return nil
	}
	if !randomImageResponse.Success {
		msg.SetContent("获取随机图片失败，回复为false")
		return nil
	}
	media := message.NewMediaObject(message.ImageFile, randomImageResponse.Url, msg.GetRouteId())
	if media == nil {
		msg.SetContent("获取随机图片失败，上传图片失败")
		return nil
	}
	msg.SetMedia(media, "随机图片")
	return nil
}

func init() {
	message.MS.RegisterBeforeSendHook(message.BeforeSendHook{
		Priority: 255,
		Fn:       getRandomImage,
	})
}
