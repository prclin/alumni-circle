package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/prclin/alumni-circle/core"
	. "github.com/prclin/alumni-circle/global"
	"github.com/prclin/alumni-circle/model"
	"github.com/prclin/alumni-circle/service"
	"net/http"
	"strconv"
)

func init() {
	auth := core.ContextRouter.Group("/auth")
	auth.POST("/sign_up", EmailSignUp) //邮箱注册
	auth.PUT("/sign_in", EmailSignIn)  //邮箱登录
	auth.POST("/api", PostAPI)
	auth.PUT("/api/:id", PutAPI)
	auth.DELETE("/api/:id", DeleteAPI)
	auth.GET("/api/list", GetAPIList)
}

// GetAPIList 获取接口列表
func GetAPIList(c *gin.Context) {
	//获取分页
	var query model.Pagination
	err := c.ShouldBindQuery(&query)
	if err != nil {
		Logger.Debug(err)
		model.Client(c)
		return
	}
	apiList, err := service.GetAPIList(query)
	if err != nil {
		Logger.Debug(err)
		model.Server(c)
		return
	}
	model.Ok(c, apiList)
}

// DeleteAPI 删除接口
func DeleteAPI(c *gin.Context) {
	//获取id
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		Logger.Debug(err)
		model.Client(c)
		return
	}
	//删除接口
	err = service.DeleteAPI(uint32(id))
	if err != nil {
		Logger.Debug(err)
		model.Server(c)
		return
	}

	model.Write(c, model.Response[any]{Code: 200, Message: "删除成功"})
}

// PutAPI 修改接口
func PutAPI(c *gin.Context) {
	//获取id
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		Logger.Debug(err)
		model.Client(c)
		return
	}
	//获取参数
	var body struct {
		Name        string  `json:"name" binding:"required"`
		Method      string  `json:"method" binding:"required"`
		Path        string  `json:"path" binding:"required"`
		Description string  `json:"description"  binding:"omitempty"`
		State       *uint8  `json:"state" binding:"required"`
		Extra       *string `json:"extra" binding:"omitempty"`
	}
	err = c.ShouldBindJSON(&body)
	if err != nil {
		Logger.Debug(err)
		model.Client(c)
		return
	}
	//更新
	api, err := service.UpdateAPI(model.TAPI{
		Id:          uint32(id),
		Name:        body.Name,
		Method:      body.Method,
		Path:        body.Path,
		Description: body.Description,
		State:       *body.State,
		Extra:       body.Extra,
	})
	if err != nil {
		Logger.Debug(err)
		model.Server(c)
		return
	}
	model.Ok(c, api)
}

// PostAPI 创建接口
func PostAPI(c *gin.Context) {
	//获取参数
	var body struct {
		Name        string  `json:"name" binding:"required"`
		Method      string  `json:"method" binding:"required"`
		Path        string  `json:"path" binding:"required"`
		Description string  `json:"description"  binding:"omitempty"`
		State       *uint8  `json:"state" binding:"required"`
		Extra       *string `json:"extra" binding:"omitempty"`
	}
	err := c.ShouldBindJSON(&body)
	if err != nil {
		Logger.Debug(err)
		model.Client(c)
		return
	}
	//创建
	api, err := service.CreateAPI(model.TAPI{
		Name:        body.Name,
		Method:      body.Method,
		Path:        body.Path,
		Description: body.Description,
		State:       *body.State,
		Extra:       body.Extra,
	})
	if err != nil {
		Logger.Debug(err)
		model.Server(c)
		return
	}
	model.Ok(c, api)
}

/*
EmailSignUp 账户注册（邮箱）
参数：邮箱 验证码 密码
*/
func EmailSignUp(c *gin.Context) {
	//参数结构
	var body struct {
		Email    string `form:"email" binding:"required"`
		Captcha  string `form:"captcha" binding:"required"`
		Password string `form:"password" binding:"required"`
	}

	//参数绑定
	err := c.ShouldBindJSON(&body)
	if err != nil {
		Logger.Debug(err)
		model.Client(c)
		return
	}

	//注册逻辑
	account := model.TAccount{
		Email:    body.Email,
		Password: body.Password,
	}
	res := service.EmailSignUp(account, body.Captcha)
	model.Write(c, res)
}

/*
EmailSignIn 账户登录（邮箱）

参数：邮箱 密码
*/
func EmailSignIn(c *gin.Context) {
	var body struct {
		Email    string `form:"email" binding:"required"`
		Password string `form:"password" binding:"required"`
	}
	//参数绑定
	err := c.ShouldBindJSON(&body)
	if err != nil {
		Logger.Debug(err)
		model.Client(c)
		return
	}

	//登录逻辑
	res := service.EmailSignIn(body.Email, body.Password)
	//登录成功回写cookie
	if res.Code == http.StatusOK {
		c.SetCookie("token", *res.Data, -1, "/", "*", false, false)
	}
	model.Write(c, res)
}
