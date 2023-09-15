package model

// QuestionSheet mbti题单
type QuestionSheet struct {
	Sheet     Sheet      `json:"sheet"`
	Questions []Question `json:"questions"`
}

// Question mbti题目
type Question struct {
	Title   string `json:"title"`
	Answers string `json:"answers"`
	CalcId  int64  `json:"calc_id"`
}

// Sheet mbti题单信息
type Sheet struct {
	Id              int64  `json:"id"`
	SheetName       string `json:"sheet_name"`
	SheetCode       string `json:"sheet_code"`
	Subtitle        string `json:"subtitle"`
	BannerSubtitle  string `json:"banner_subtitle"`
	DescriptionText string `json:"description_text"`
}
