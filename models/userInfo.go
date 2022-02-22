package models

type UserInfo struct {
	Login     string `json:"login"`
	Id     int64 `json:"id"`
	Node_Id     string `json:"node_id"`
	Message     string `json:"message"`
}
