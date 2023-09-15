package service

import (
	"encoding/json"
	"errors"
	"github.com/prclin/alumni-circle/dao"
	"github.com/prclin/alumni-circle/global"
	"github.com/prclin/alumni-circle/model"
	"github.com/prclin/alumni-circle/util"
	"github.com/redis/go-redis/v9"
	"github.com/valyala/fasthttp"
	"net/http"
	"net/url"
)

//SubmitQuestionSheet 提交测试答案
//
//暂时未实现重复提交判断，重复提交时应获取已提交的结果
func SubmitQuestionSheet(accountId uint64, answer model.SheetAnswer) (*model.SubmitResult, error) {
	//创建请求
	request := fasthttp.AcquireRequest()
	request.Header.SetMethod(http.MethodPost)
	request.SetRequestURI(global.Configuration.MBTI.API.Submit)
	request.Header.Set("Authorization", "Token "+global.Configuration.MBTI.Token)
	request.Header.SetContentType("application/json")
	request.SetBodyString(string(util.IgnoreError(json.Marshal(&answer))))
	//创建响应
	response := fasthttp.AcquireResponse()
	//执行请求
	if err := fasthttp.Do(request, response); err != nil || response.StatusCode() != 200 {
		err = util.Ternary(err == nil, errors.New(string(response.Body())), err)
		global.Logger.Debug(err)
		return nil, err
	}
	body := struct {
		Success bool                `json:"success"`
		Result  *model.SubmitResult `json:"result"`
	}{
		Result: &model.SubmitResult{},
	}
	//读取响应
	if err := json.Unmarshal(response.Body(), &body); err != nil {
		global.Logger.Debug(err)
		return nil, err
	}
	//保存结果id,暂时存储到account_info的extra字段
	tx := global.Datasource.Begin()
	defer tx.Commit()
	accountInfoDao := dao.NewAccountInfoDao(tx)
	err := accountInfoDao.UpdateMBTIResultIdById(accountId, body.Result.Id)
	if err != nil {
		tx.Rollback()
		global.Logger.Debug(err)
		return nil, err
	}
	return body.Result, nil
}

func GetMBTIQuestionSheet(sheetName string) (model.QuestionSheet, error) {
	var qs model.QuestionSheet
	//获取缓存中的题目
	sheet, err := dao.GetString("MBTI:sheet:" + sheetName)
	//获取成功
	if err == nil {
		if err := json.Unmarshal([]byte(sheet), &qs); err != nil {
			global.Logger.Debug(err)
			return qs, err
		}
		return qs, nil
	}

	//获取失败
	if err != nil && err != redis.Nil {
		global.Logger.Debug(err)
		return qs, err
	}

	//创建请求
	request := fasthttp.AcquireRequest()
	request.Header.SetMethod(http.MethodGet)
	path, err := url.JoinPath(global.Configuration.MBTI.API.Sheet, sheetName)
	if err != nil {
		global.Logger.Debug(err)
		return qs, err
	}
	request.SetRequestURI(path)
	request.Header.Set("Authorization", "Token "+global.Configuration.MBTI.Token)
	//创建响应
	response := fasthttp.AcquireResponse()
	//执行请求
	if err := fasthttp.Do(request, response); err != nil || response.StatusCode() != 200 {
		err = util.Ternary(err == nil, errors.New(string(response.Body())), err)
		global.Logger.Debug(err)
		return qs, err
	}

	//读取响应
	if err := json.Unmarshal(response.Body(), &qs); err != nil {
		global.Logger.Debug(err)
		return qs, err
	}
	//缓存题单
	if err := dao.SetString("MBTI:sheet:"+sheetName, string(response.Body()), 0); err != nil {
		global.Logger.Info(err)
	}
	return qs, nil
}
