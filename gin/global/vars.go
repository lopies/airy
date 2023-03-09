package global

import (
	"github.com/MQEnergy/gin-framework/config"
	"github.com/MQEnergy/gin-framework/pkg/lib"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	clientv3 "go.etcd.io/etcd/client/v3"
	"gorm.io/gorm"
)

var (
	DB     *gorm.DB         // Mysql数据库
	Logger *lib.Logger      // 日志
	ETCD   *clientv3.Client //ETCD
	Redis  *redis.Client    // redis连接池
	Router *gin.Engine      // 路由
	Cfg    *config.Conf     // yaml配置
)
