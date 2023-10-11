package global

import (
	dysmsapi "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/prclin/alumni-circle/config"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	Configuration *config.Configuration
	Logger        *zap.SugaredLogger
	Datasource    *gorm.DB
	RedisClient   *redis.Client
	OSSClient     *oss.Client
	SMSClient     *dysmsapi.Client
)
