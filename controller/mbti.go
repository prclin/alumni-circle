package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/prclin/alumni-circle/core"
	"github.com/prclin/alumni-circle/global"
	"github.com/prclin/alumni-circle/model"
	"github.com/prclin/alumni-circle/service"
	"github.com/prclin/alumni-circle/util"
)

func init() {
	mbti := core.ContextRouter.Group("/mbti")
	mbti.GET("/question-sheet/:sheet_name", GetQuestionSheet)
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
