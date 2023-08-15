package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prclin/alumni-circle/core"
	"github.com/prclin/alumni-circle/global"
	"github.com/prclin/alumni-circle/model/entity"
	"github.com/prclin/alumni-circle/model/request"
	"github.com/prclin/alumni-circle/model/response"
	"github.com/prclin/alumni-circle/service"
	"sort"
	"strconv"
)

// 注册路由
func init() {
	breaks := core.ContextRouter.Group("/break")
	breaks.POST("/create", CreateBreak)
	breaks.POST("/add/image", AddImageToBreak)
}

// 从Cookie获取用户id
// TODO:应放在controller/account中,暂存在此
func GetIdFromCookie(c *gin.Context) (id int, err error) {
	cookie, err := c.Cookie("id")
	if err != nil {
		return 0, err
	}
	id, err = strconv.Atoi(cookie)
	return
}

// CreateBreak 新建课间
func CreateBreak(c *gin.Context) {
	// 校验cookie
	accountId, err := GetIdFromCookie(c)
	if err != nil {
		response.NLI(c)
		return
	}
	// 校验JSON
	aBreak := new(entity.Break)
	if err := c.ShouldBindJSON(aBreak); err != nil {
		response.Client(c, err)
		return
	}
	// 校验break是否存在
	if service.BreakExist(aBreak.Id, accountId) {
		response.Client(c, fmt.Sprintf("break(id=%d) already exist", aBreak.Id))
		return
	}
	// 校验课间标题,内容不得为空
	if aBreak.Title == "" || aBreak.Content == "" {
		response.Client(c, "tile and content cannot be empty")
		return
	}
	// 创建课间
	aBreak.AccountId = accountId
	if err := service.CreateBreak(aBreak); err != nil {
		response.Server(c, err)
	}
	response.Ok(c, aBreak)
}

// AddImageToBreak 添加图片至课间
func AddImageToBreak(c *gin.Context) {
	// 校验cookie
	accountId, err := GetIdFromCookie(c)
	if err != nil {
		response.NLI(c)
		return
	}
	// 校验JSON
	vo := new(request.BreakAddImageVO)
	if err := c.ShouldBindJSON(vo); err != nil {
		response.Client(c, err)
		return
	}
	breakId := vo.BreakId
	bindingList := vo.BindingList
	// 校验break是否存在
	if !service.BreakExist(breakId, accountId) {
		response.Client(c, "break does not exist")
		return
	}
	// 校验binding列表
	length := len(bindingList)
	sizeLimit := global.Configuration.Limit.Size.PictureInBreak
	if length < 1 || length > sizeLimit {
		response.Client(c, "number of pictures should not exceed 9")
		return
	}
	sort.Slice(bindingList, func(i, j int) bool {
		return bindingList[i].Order < bindingList[j].Order
	})
	for i, binding := range bindingList {
		// binding.order应从0开始不间断
		if binding.Order != i {
			response.Client(c, fmt.Sprintf("image(id=%d) order slhould be %d but infact %d", binding.ImageId, i, binding.Order))
			return
		}
		// 校验image是否存在
		if !service.ImageExist(binding.ImageId) {
			response.Client(c, fmt.Sprintf("image(id=%d) does not exist", binding.ImageId))
			return
		}
	}
	// 添加图片至课间
	if err = service.AddImageToBreak(breakId, bindingList); err != nil {
		response.Server(c, err.Error())
		return
	}
	response.Ok(c, "success")
}
