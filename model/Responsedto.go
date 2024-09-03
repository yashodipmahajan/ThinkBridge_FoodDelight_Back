package model

type ResponseDto struct {
	LoggedInAdmin Admin  `json:"loggedinAdmin"`
	Token         string `json:"token"`
}
