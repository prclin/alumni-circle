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
	"net/url"
)

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
	request.Header.SetMethod("GET")
	path, err := url.JoinPath(global.Configuration.MBTI.API.Sheet, sheetName)
	if err != nil {
		global.Logger.Debug(err)
		return qs, err
	}
	request.Header.SetRequestURI(path)
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
