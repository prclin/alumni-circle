package global

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gorilla/websocket"
	"github.com/prclin/alumni-circle/config"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	Configuration     *config.Configuration
	Logger            *zap.SugaredLogger
	Datasource        *gorm.DB
	RedisClient       *redis.Client
	OSSClient         *oss.Client
	WebsocketUpgrader *websocket.Upgrader
)
