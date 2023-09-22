package controllers

import (
	"gin-tpl/app"
	"gin-tpl/app/http/models"
	"gin-tpl/app/utils"
	"github.com/gin-gonic/gin"
)

type Home struct {
}

func (h Home) Index(ctx *gin.Context) {
	app := ctx.MustGet("app").(*app.App)

	var ret []models.User
	err := app.DBW().Find(&ret).Error

	if err != nil {
		utils.Error(ctx, utils.WithMsg(err.Error()))
	}

	utils.Success(ctx, utils.WithData(ret))
}
