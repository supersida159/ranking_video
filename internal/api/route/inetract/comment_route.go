package route

import (
	handler "ranking_video/internal/api/handler/interact"
	"ranking_video/pkg/utils/app_context"

	"github.com/gin-gonic/gin"
)

func CommentRoute(r *gin.RouterGroup, appCtx app_context.AppContext) {
	commentHandler := handler.NewCommentEventHandler(appCtx)

	r.POST("/comment", commentHandler.Add)
	r.GET("/comment", commentHandler.List)
	r.DELETE("/comment/:id", commentHandler.SoftDelete)
}
