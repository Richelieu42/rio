package main

import (
	"github.com/gin-gonic/gin"
	"github.com/richelieu42/go-scales/src/core/strKit"
	"github.com/richelieu42/rio"
	"github.com/richelieu42/rio/src/manager"
	"github.com/sirupsen/logrus"
	"net/http"
)

type (
	TestListener struct {
	}
)

func (t *TestListener) OnIllegalRequest(ctx *gin.Context) {
	text := strKit.Format("Request(Proto: %s, Method: %s, Connection: %s, Upgrade: %s, RequestURI: %s) is illegal.",
		ctx.Request.Proto, ctx.Request.Method, ctx.Request.Header["Connection"], ctx.Request.Header["Upgrade"], ctx.Request.RequestURI)
	ctx.String(http.StatusOK, text)
}

func (t *TestListener) OnFailureToUpgrade(ctx *gin.Context, err error) {
	logrus.Errorf("Fail to upgrade, error: %v", err)
}

func (t *TestListener) OnHandshake(c *manager.Channel) {
	logrus.Infof("Channel(id: %s) is established.", c.GetId())
}

func (t *TestListener) OnMessage(c *manager.Channel, messageType int, data []byte) {
	logrus.Infof("Channel(id: %s) receives a message(type: %d, data: %s)", c.GetId(), messageType, string(data))
}

func (t *TestListener) OnClose(c *manager.Channel, code int, text string) {
	logrus.Infof("Channel(id: %s) is closed with code(%d) and text(%s).", c.GetId(), code, text)
}

func main() {
	listener := &TestListener{}
	handler, err := rio.NewGinHandler(listener)
	if err != nil {
		logrus.Panic(err)
		return
	}

	r := gin.Default()
	r.GET("/ping", handler)
	_ = r.Run(":80")
}
