package message

import (
	"encoding/json"
	"fmt"
	"io"

	"qqbot/service/auth"
	"qqbot/utils"
)

type MediaObject struct {
	FileUUID string `json:"file_uuid"`
	FileInfo string `json:"file_info"`
	TTL      int    `json:"ttl"`
}

type FileType int

const (
	ImageFile FileType = iota + 1
	VideoFile
	AudioFile
)

func NewMediaObject(fileType FileType, url, groupId string) *MediaObject {
	token := auth.AuthHelper.GetToken()
	resp, err := utils.NetHelper.POST(fmt.Sprintf("https://api.sgroup.qq.com/v2/groups/%s/files", groupId), map[string]interface{}{
		"file_type": fileType,
		"url":       url,
	}, utils.WithToken(token))
	if err != nil {
		fmt.Println("上传文件失败", err)
		return nil
	}
	defer resp.Body.Close()
	bytesData, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取文件失败", err)
		return nil
	}
	var mediaObject MediaObject
	err = json.Unmarshal(bytesData, &mediaObject)
	if err != nil {
		fmt.Println("解析json失败", err)
		return nil
	}
	if mediaObject.FileUUID == "" || mediaObject.FileInfo == "" {
		fmt.Println("上传文件失败", string(bytesData))
		return nil
	}
	return &mediaObject

}

func (m *MediaObject) ToStruct() interface{} {
	return map[string]interface{}{
		"file_uuid": m.FileUUID,
		"file_info": m.FileInfo,
	}
}
