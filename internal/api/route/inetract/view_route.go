package route

import (
	handler "ranking_video/internal/api/handler/interact"
	"ranking_video/pkg/utils/app_context"

	"github.com/gin-gonic/gin"
)

func ViewRoute(r *gin.RouterGroup, appCtx app_context.AppContext) {
	viewHandler := handler.NewViewEventHandler(appCtx)

	r.POST("/view", viewHandler.Add)
	r.GET("/view", viewHandler.List)
	r.DELETE("/view/:id", viewHandler.SoftDelete)
}
