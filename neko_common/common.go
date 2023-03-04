package neko_common

import (
	"io"
	"net/http"
)

var Version_v2ray string = "N/A"
var Version_neko string = "N/A"

var Debug bool

// platform

var RunMode int

const (
	RunMode_Other = iota
	RunMode_NekoRay_Core
	RunMode_NekoBox_Core
	RunMode_NekoBoxForAndroid
)

var NB4A_GuiLogWriter io.Writer

// proxy

var GetProxyHttpClient func() *http.Client
