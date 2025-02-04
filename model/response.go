package model

type ResponseTexttrans struct {
	LogId  uint64 `json:"log_id"`
	Result Result `json:"result"`
}

type Result struct {
	From        string        `json:"from"`
	To          string        `json:"to"`
	TransResult []TransResult `json:"trans_result"`
}

type TransResult struct {
	Dst string `json:"dst"`
	Src string `json:"src"`
}
