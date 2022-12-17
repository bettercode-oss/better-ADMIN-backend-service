package dtos

import (
	"fmt"
	"net/http"
	"strings"
)

type WebHookInformation struct {
	Id          uint   `json:"id"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type WebHookDetails struct {
	Id              uint            `json:"id"`
	Name            string          `json:"name"`
	Description     string          `json:"description"`
	WebHookCallSpec WebHookCallSpec `json:"webHookCallSpec"`
}

func (w *WebHookDetails) FillInWebHookCallSpec(httpRequest *http.Request, accessToken string) {
	url := fmt.Sprintf("%v://%v/api/web-hooks/%v/note", strings.ToLower(strings.Split(httpRequest.Proto, "/")[0]), httpRequest.Host, w.Id)
	spec := WebHookCallSpec{
		HttpRequestMethod: http.MethodPost,
		Url:               url,
		AccessToken:       accessToken,
		SampleRequest: fmt.Sprintf("curl -X %v %v -H \"Content-Type: application/json\" -H \"Authorization: Bearer %v\" -d '{\"text\":\"테스트 메시지 입니다.\"}'",
			http.MethodPost, url, accessToken),
	}

	w.WebHookCallSpec = spec
}

type WebHookCallSpec struct {
	HttpRequestMethod string `json:"httpRequestMethod"`
	Url               string `json:"url"`
	AccessToken       string `json:"accessToken"`
	SampleRequest     string `json:"sampleRequest"`
}

type WebHookMessage struct {
	Title string `json:"title"`
	Text  string `json:"text" binding:"required"`
}
