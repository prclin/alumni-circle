package util

import (
	"github.com/gin-gonic/gin"
	_error "github.com/prclin/alumni-circle/error"
	"github.com/prclin/alumni-circle/global"
	"github.com/prclin/alumni-circle/model"
)

func GetTokenClaims(context *gin.Context) (*model.TokenClaims, error) {
	//获取token
	token, err := context.Cookie("token")
	if err != nil {
		global.Logger.Debug(err)
		return nil, _error.TokenNotFoundError
	}

	//解析token
	claims, err := ParseToken(token)
	if err != nil {
		global.Logger.Debug(err)
		return nil, _error.InvalidTokenError
	}
	return claims, nil
}
