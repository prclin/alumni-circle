package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/prclin/alumni-circle/core"
	. "github.com/prclin/alumni-circle/global"
	"github.com/prclin/alumni-circle/model"
	"github.com/prclin/alumni-circle/service"
	"net/http"
	"regexp"
	"strconv"
)

func init() {
	auth := core.ContextRouter.Group("/auth")
	auth.POST("/sign_up", EmailSignUp)       //邮箱注册
	auth.PUT("/sign_in", EmailSignIn)        //邮箱登录
	auth.POST("/sign_in/phone", PhoneSignIn) //手机号码注册或登录
	auth.POST("/api", PostAPI)
	auth.PUT("/api/:id", PutAPI)
	auth.DELETE("/api/:id", DeleteAPI)
	auth.GET("/api/list", GetAPIList)
	auth.POST("/role", PostRole)
	auth.PUT("/role/:id", PutRole)
	auth.DELETE("/role/:id", DeleteRole)
	auth.GET("/role/list", GetRoleList)
	auth.POST("/api/allocation", PostAPIAllocation)
	auth.DELETE("/api/allocation", DeleteAPIAllocation)
	auth.POST("/role/allocation", PostRoleAllocation)
	auth.DELETE("/role/allocation", DeleteRoleAllocation)
}

var (
	// 手机号校验
	phoneRegexp *regexp.Regexp
)

// 初始全局化变量
func init() {
	//初始化phoneRegexp
	phoneReg, err := regexp.Compile("^((13[0-9])|(14[5|7])|(15([0-3]|[5-9]))|(17[013678])|(18[0,5-9]))\\d{8}$")
	if err != nil {
		Logger.Fatal(err)
	}
	phoneRegexp = phoneReg
}

// PhoneSignIn 手机号注册登录
func PhoneSignIn(context *gin.Context) {
	//获取参数
	var body struct {
		Phone    string `json:"phone" binding:"required,min=11,max=11"`
		Method   string `json:"method" binding:"required"`
		Password string `json:"password" binding:"omitempty,min=6,max=18"`
		Captcha  string `json:"captcha" binding:"omitempty,min=6,max=6"`
	}

	err := context.ShouldBindJSON(&body)
	if err != nil {
		Logger.Debug(err)
		model.Client(context)
		return
	}

	//手机号错误
	if !phoneRegexp.MatchString(body.Phone) {
		Logger.Debug(err)
		model.Client(context)
		return
	}
	var res model.Response[*string]
	switch body.Method {
	case "password":
		//手机密码登录
		res = service.PhonePasswordSignIn(body.Phone, body.Password)
		break
	case "captcha":
		//手机验证码登录
		res = service.PhoneCaptchaSignIn(body.Phone, body.Captcha)
		break
	}
	//登录成功回写cookie
	if res.Code == http.StatusOK {
		context.SetCookie("token", *res.Data, -1, "/", "*", false, false)
	}
	model.Write(context, res)
}

// DeleteRoleAllocation 解配角色
func DeleteRoleAllocation(c *gin.Context) {
	//获取参数
	var body struct {
		AccountIds []uint64 `json:"account_ids" binding:"required,min=1"`
		RoleIds    []uint32 `json:"role_ids" binding:"required,min=1"`
	}
	err := c.ShouldBindJSON(&body)
	if err != nil {
		Logger.Debug(err)
		model.Client(c)
		return
	}
	//分配
	err = service.RevokeRoleAllocation(body.AccountIds, body.RoleIds)
	if err != nil {
		Logger.Debug(err)
		model.Server(c)
		return
	}
	model.Write(c, model.Response[any]{Code: http.StatusOK, Message: "解配成功"})
}

// PostRoleAllocation  分配角色
func PostRoleAllocation(c *gin.Context) {
	//获取参数
	var body struct {
		AccountIds []uint64 `json:"account_ids" binding:"required,min=1"`
		RoleIds    []uint32 `json:"role_ids" binding:"required,min=1"`
	}
	err := c.ShouldBindJSON(&body)
	if err != nil {
		Logger.Debug(err)
		model.Client(c)
		return
	}
	//分配
	err = service.AllocateRole(body.AccountIds, body.RoleIds)
	if err != nil {
		Logger.Debug(err)
		model.Server(c)
		return
	}
	model.Write(c, model.Response[any]{Code: http.StatusOK, Message: "分配成功"})
}

// DeleteAPIAllocation 接口解配
func DeleteAPIAllocation(c *gin.Context) {
	//获取参数
	var body struct {
		RoleIds []uint32 `json:"role_ids" binding:"required,min=1"`
		APIIds  []uint32 `json:"api_ids"  binding:"required,min=1"`
	}
	err := c.ShouldBindJSON(&body)
	if err != nil {
		Logger.Debug(err)
		model.Client(c)
		return
	}
	//解配
	err = service.RevokeAPIAllocation(body.RoleIds, body.APIIds)
	if err != nil {
		Logger.Debug(err)
		model.Server(c)
		return
	}
	model.Write(c, model.Response[any]{Code: http.StatusOK, Message: "解配成功"})
}

// PostAPIAllocation 接口分配
func PostAPIAllocation(c *gin.Context) {
	//获取参数
	var body struct {
		RoleIds []uint32 `json:"role_ids" binding:"required,min=1"`
		APIIds  []uint32 `json:"api_ids"  binding:"required,min=1"`
	}
	err := c.ShouldBindJSON(&body)
	if err != nil {
		Logger.Debug(err)
		model.Client(c)
		return
	}
	//分配
	err = service.AllocateAPI(body.RoleIds, body.APIIds)
	if err != nil {
		Logger.Debug(err)
		model.Server(c)
		return
	}
	model.Write(c, model.Response[any]{Code: http.StatusOK, Message: "分配成功"})
}

// GetRoleList 获取角色列表
func GetRoleList(c *gin.Context) {
	//获取参数
	var query model.Pagination
	err := c.ShouldBindQuery(&query)
	if err != nil {
		Logger.Debug(err)
		model.Client(c)
		return
	}
	//获取角色列表
	roles, err := service.GetRoleList(query)
	if err != nil {
		Logger.Debug(err)
		model.Server(c)
		return
	}
	model.Ok(c, roles)
}

// DeleteRole 删除角色
func DeleteRole(c *gin.Context) {
	//获取id
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		Logger.Debug(err)
		model.Client(c)
		return
	}
	//删除
	err = service.DeleteRole(uint32(id))
	if err != nil {
		Logger.Debug(err)
		model.Server(c)
		return
	}
	model.Write(c, model.Response[any]{Code: http.StatusOK, Message: "删除成功"})
}

// PutRole 更新角色
func PutRole(c *gin.Context) {
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
		Identifier  string  `json:"identifier" binding:"required"`
		Description *string `json:"description" binding:"required"`
		State       *uint8  `json:"state" binding:"required"`
	}
	err = c.ShouldBindJSON(&body)
	if err != nil {
		Logger.Debug(err)
		model.Client(c)
		return
	}
	//更新
	role, err := service.UpdateRole(model.TRole{
		Id:          uint32(id),
		Name:        body.Name,
		Identifier:  body.Identifier,
		Description: *body.Description,
		State:       *body.State,
	})
	if err != nil {
		Logger.Debug(err)
		model.Server(c)
		return
	}
	model.Ok(c, role)
}

// PostRole 创建角色
func PostRole(c *gin.Context) {
	//获取参数
	var body struct {
		Name        string  `json:"name" binding:"required"`
		Identifier  string  `json:"identifier" binding:"required"`
		Description *string `json:"description" binding:"required"`
		State       *uint8  `json:"state" binding:"required"`
	}
	err := c.ShouldBindJSON(&body)
	if err != nil {
		Logger.Debug(err)
		model.Client(c)
		return
	}
	//创建
	role, err := service.CreateRole(model.TRole{
		Name:        body.Name,
		Identifier:  body.Identifier,
		Description: *body.Description,
		State:       *body.State,
	})
	if err != nil {
		Logger.Debug(err)
		model.Server(c)
		return
	}
	model.Ok(c, role)
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
