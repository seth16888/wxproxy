package bootstrap

import (
	"github.com/seth16888/wxproxy/internal/di"
	"github.com/seth16888/wxproxy/internal/server"
)

func StartApp() error {
  return server.Start(di.DI)
}
