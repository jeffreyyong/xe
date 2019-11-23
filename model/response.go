package model

type ConvertResp struct {
	From           string `json:"from"`
	To             string `json:"to"`
	Rate           string `json:"rate"`
	TrendFromWeek  string `json:"trendLastWeek"`
	Recommendation string `json:"recommendation"`
}
