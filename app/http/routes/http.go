package routes

import (
	"gin-tpl/app/http/controllers"
	"github.com/gin-gonic/gin"
)

func RegRoute(r *gin.Engine) {
	r.GET("/test", controllers.Home{}.Index)
}
