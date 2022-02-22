package models

//参数
type Param struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Gas      string `json:"gas"`
	GasPrice string `json:"gasPrice"`
	Data     string `json:"data"`
	Value    string `json:"value"`
	Token    string `json:"token"`
}
