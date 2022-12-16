package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/richelieu42/go-scales/src/core/strKit"
	"github.com/richelieu42/rio/src/manager"
	"github.com/richelieu42/rio/src/rio"
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
	logrus.Errorf("Fail to upgrade request(Proto: %s, Method: %s, Connection: %s, Upgrade: %s, RequestURI: %s), error: %v",
		ctx.Request.Proto, ctx.Request.Method, ctx.Request.Header["Connection"], ctx.Request.Header["Upgrade"], ctx.Request.RequestURI,
		err)
}

func (t *TestListener) OnHandshake(c *manager.Channel) {
	logrus.Infof("Channel(id: %s) is established.", c.GetId())
	_ = c.PushMessage(websocket.TextMessage, []byte(strKit.Format("Hello, id of this channel is [%s].", c.GetId())))

	_ = c.PushMessage(websocket.TextMessage, []byte("后端主动关闭"))
	c.Close()
}

func (t *TestListener) OnMessage(c *manager.Channel, messageType int, data []byte) {
	logrus.Infof("Channel(id: %s) receives a message(type: %d, data: %s)", c.GetId(), messageType, string(data))
}

func (t *TestListener) OnCloseByFrontend(c *manager.Channel, code int, text string) {
	logrus.Infof("Channel(id: %s) is closed by frontend with code(%d) and text(%s).", c.GetId(), code, text)
}

func (t *TestListener) OnCloseByBackend(c *manager.Channel) {
	logrus.Infof("Channel(id: %s) is closed by backend.", c.GetId())
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
