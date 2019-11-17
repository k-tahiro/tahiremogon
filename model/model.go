package model

type (
	Code struct {
		ID   int    `json:"id"`
		Key  string `json:"key"`
		Code string `json:"code"`
	}

	TransmitResponse struct {
		Sucess bool `json:"sucess"`
		Result int  `json:"result"`
	}
)
