package route

import (
	handler "ranking_video/internal/api/handler/interact"
	"ranking_video/pkg/utils/app_context"

	"github.com/gin-gonic/gin"
)

func ShareRoute(r *gin.RouterGroup, appCtx app_context.AppContext) {
	shareHandler := handler.NewShareEventHandler(appCtx)

	r.POST("/share", shareHandler.Add)
	r.GET("/share", shareHandler.List)
	r.DELETE("/share/:id", shareHandler.SoftDelete)
}
