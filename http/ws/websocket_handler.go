package ws

import (
	"better-admin-backend-service/adapters"
	"better-admin-backend-service/helpers"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

func WebSocketHandler(upgrader websocket.Upgrader) gin.HandlerFunc {
	fn := func(ctx *gin.Context) {
		ws, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			helpers.ErrorHelper().InternalServerError(ctx, err)
			return
		}

		webSocketId := ctx.Param("id")
		adapters.WebSocketAdapter().AddConnection(webSocketId, ws)

		ctx.Status(http.StatusOK)
	}

	return gin.HandlerFunc(fn)
}
