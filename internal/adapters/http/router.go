package httpadapter

import "github.com/gin-gonic/gin"

func NewRouter(handler *NoteHandler) *gin.Engine {
	router := gin.Default()

	notes := router.Group("/notes")
	notes.POST("", handler.Create)
	notes.GET("", handler.List)
	notes.PUT("/:id", handler.Update)
	notes.DELETE("/:id", handler.Delete)

	return router
}
