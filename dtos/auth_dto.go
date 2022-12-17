package dtos

type MemberSignIn struct {
	Id       string `json:"id" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type DoorayMember struct {
	Id                   string `json:"id"`
	UserCode             string `json:"userCode"`
	Name                 string `json:"name"`
	ExternalEmailAddress string `json:"externalEmailAddress"`
}

type GoogleMember struct {
	Id      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
	Hd      string `json:"hd"`
}
