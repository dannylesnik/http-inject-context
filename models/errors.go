package models

//Error - HTTP Error Response
type Error struct {
	Error   string `json:"error"`
	Message string `json:"msg"`
	Code    int16  `json:"code"`
}