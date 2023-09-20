package controllers

import (
	"gin-tpl/app"
	"gin-tpl/app/http/models"
	"github.com/gin-gonic/gin"
)

type Home struct {
}

func (h Home) Index(ctx *gin.Context) {
	app := ctx.MustGet("app").(*app.App)

	var ret []models.User
	err := app.DBW().Find(&ret).Error

	if err != nil {
		ctx.JSON(500, gin.H{
			"code": 500,
			"msg":  err.Error(),
		})
	}

	ctx.JSON(200, gin.H{
		"code": 200,
		"msg":  "success",
		"data": ret,
	})
}
