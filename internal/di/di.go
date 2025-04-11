package di

import (
	"github.com/seth16888/wxcommon/hc"
	"github.com/seth16888/wxcommon/logger"
	"github.com/seth16888/wxcommon/redis"

	"github.com/seth16888/wxproxy/internal/biz"
	"github.com/seth16888/wxproxy/internal/config"

	"github.com/seth16888/wxproxy/internal/service"
	"go.uber.org/zap"
)

var DI *Container

type Container struct {
	Conf *config.Bootstrap
	Log  *zap.Logger
	Svc  *service.MPProxyService
  Redis *redis.RedisClient
}

func NewContainer(configFile string) *Container {
  conf:= config.ReadConfigFromFile(configFile)
  log := logger.InitLogger(conf.Log)

  redis.ConnectRedis(conf.Redis.Addr, conf.Redis.Username,
    conf.Redis.Password, conf.Redis.DB, log)

  hc := hc.NewClient(hc.DefaultTimeout, hc.DefaultIdleConnTimeout, hc.CommonCheckRedirect)

  uc := biz.NewMPProxyUsecase(hc, log)

  svc := service.NewMPProxyService(uc, log)

	DI = &Container{
    Conf: conf,
    Log: log,
    Svc: svc,
    Redis: redis.Redis,
  }
	return DI
}
