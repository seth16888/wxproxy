package biz

import (
	"github.com/seth16888/wxcommon/hc"
)

func NewHttpClient() *hc.Client {
	return hc.NewClient(hc.DefaultTimeout, hc.DefaultIdleConnTimeout, hc.CommonCheckRedirect)
}
