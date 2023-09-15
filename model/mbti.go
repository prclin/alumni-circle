package model

import "time"

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

// SheetAnswer mbti测试答案提交参数
//
// Sheet和Answers必须填写，其它选填
type SheetAnswer struct {
	Sheet      string `json:"sheet"`
	Answers    []int  `json:"answers"`
	TimeLength int    `json:"time_length"`
	IP         string `json:"ip"`
	Nickname   string `json:"nickname"`
}

// SubmitResult mbti测试提交结果
type SubmitResult struct {
	Id                string     `json:"id"`
	Sheet             string     `json:"sheet"`
	EightValues       []float32  `json:"eight_values"`
	EightDescriptions []string   `json:"eight_descriptions"`
	TopMbtiResults    [][]string `json:"top_mbti_results"`
	TopMbtiData       []MBTIData `json:"top_mbti_data"`
	CreatedAt         time.Time  `json:"created_at"`
}

type MBTIData struct {
	Name             string          `json:"name"`
	Chinese          string          `json:"chinese"`
	Title            string          `json:"title"`
	PrimaryText      string          `json:"primary_text"`
	DescriptionText  string          `json:"description_text"`
	Features         string          `json:"features"`
	Image            string          `json:"image"`
	RecommendbookSet []Recommendbook `json:"recommendbook_set"`
	PublicfigureSet  []Publicfigure  `json:"publicfigure_set"`
}

type Publicfigure struct {
	Index  int    `json:"index"`
	Name   string `json:"name"`
	Title  string `json:"title"`
	Dictum string `json:"dictum"`
	Image  string `json:"image"`
}

type Recommendbook struct {
	Image string `json:"image"`
	Index int    `json:"index"`
	Name  string `json:"name"`
}
