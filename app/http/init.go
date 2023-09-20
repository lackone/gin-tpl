package http

import (
	"gin-tpl/app"
	"gin-tpl/app/http/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

func NewServer(app *app.App) error {
	r := gin.Default()

	attachVar(r, app)
	useCORS(r)

	routes.RegRoute(r)

	err := r.Run(app.Config.Http.Listen)
	if err != nil {
		return err
	}

	return nil
}

// attachVar 附加变量
func attachVar(r *gin.Engine, app *app.App) {
	r.Use(func(ctx *gin.Context) {
		ctx.Set("app", app)
	})
}

// useCORS 使用CORS中间件
func useCORS(r *gin.Engine) {
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "PUT", "PATCH"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	}))
}
