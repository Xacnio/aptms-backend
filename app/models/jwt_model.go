package models

type AccessDetails struct {
	AccessUuid string
	UserId     uint
}

type TokenDetails struct {
	AccessToken string
	AccessUuid  string
	AtExpires   int64
}

type AuthDetails struct {
	AccessToken string `json:"AccessToken"`
	AccessUuid  string `json:"AccessUuid"`
	UserID      string `json:"UserID"`
}

type AuthResponse struct {
	User        User        `json:"user"`
	AuthDetails AuthDetails `json:"tokens"`
}
