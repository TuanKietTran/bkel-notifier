package model

type UserInfo struct {
	UserId int `json:"userid"`
}

type Message struct {
	Id             int    `json:"id"`
	UserFrom       string `json:"userfromfullname"`
	Text           string `json:"text"`
	IsNotification int    `json:"notification"`
}

type MsgResponse struct {
	Messages []Message `json:"messages"`
}
