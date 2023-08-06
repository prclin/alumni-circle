package global

import (
	"github.com/prclin/alumni-circle/config"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	Configuration *config.Configuration
	Logger        *zap.SugaredLogger
	Datasource    *gorm.DB
)
