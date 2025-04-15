package route

import (
	handler "ranking_video/internal/api/handler/interact"
	"ranking_video/pkg/utils/app_context"

	"github.com/gin-gonic/gin"
)

func LikeRoute(r *gin.RouterGroup, appCtx app_context.AppContext) {
	likeHandler := handler.NewLikeEventHandler(appCtx)

	r.POST("/like", likeHandler.Add)
	r.GET("/like", likeHandler.List)
	r.DELETE("/like/:id", likeHandler.SoftDelete)
}
