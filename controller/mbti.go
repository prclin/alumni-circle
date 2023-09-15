package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/prclin/alumni-circle/core"
	"github.com/prclin/alumni-circle/global"
	"github.com/prclin/alumni-circle/model"
	"github.com/prclin/alumni-circle/service"
	"github.com/prclin/alumni-circle/util"
	"sort"
)

func init() {
	mbti := core.ContextRouter.Group("/mbti")
	mbti.GET("/question-sheet/:sheet_name", GetQuestionSheet)
	mbti.POST("/answer", PostQuestionSheetAnswer)
}

// PostQuestionSheetAnswer 提交mbti测试答案
func PostQuestionSheetAnswer(context *gin.Context) {
	//获取cookie
	cookie, err := context.Cookie("token")
	if err != nil {
		global.Logger.Debug(err)
		model.Client(context)
		return
	}
	claims, err := util.ParseToken(cookie)
	if err != nil {
		global.Logger.Debug(err)
		model.Client(context)
		return
	}
	//获取参数
	var body struct {
		SheetName string `json:"sheet_name" binding:"required"`
		Answers   []struct {
			Answer int `json:"answer" binding:"required,min=1,max=5"`
			CalcId int `json:"calc_id" binding:"required,min=1"`
		} `json:"answers" binding:"required"`
	}
	err = context.ShouldBindJSON(&body)
	if err != nil {
		global.Logger.Debug(err)
		model.Client(context)
		return
	}

	//构造参数
	answer := model.SheetAnswer{
		Sheet:   body.SheetName,
		Answers: make([]int, 0, len(body.Answers)),
	}
	sort.Slice(body.Answers, func(i, j int) bool {
		return body.Answers[i].CalcId < body.Answers[j].CalcId
	})
	for _, value := range body.Answers {
		answer.Answers = append(answer.Answers, value.Answer)
	}
	//提交
	result, err := service.SubmitQuestionSheet(claims.Id, answer)
	if err != nil {
		global.Logger.Debug(err)
		model.Server(context)
		return
	}
	model.Ok(context, result)

}

// GetQuestionSheet 获取mbti测试题
func GetQuestionSheet(context *gin.Context) {
	sheetName := context.Param("sheet_name")
	if !util.Ternary(sheetName == "Standard", true, util.Ternary(sheetName == "Jungus", true, false)) {
		global.Logger.Debug("unsupported question sheet")
		model.Client(context)
		return
	}
	sheet, err := service.GetMBTIQuestionSheet(sheetName)
	if err != nil {
		model.Server(context)
		return
	}
	model.Ok(context, sheet)
}
