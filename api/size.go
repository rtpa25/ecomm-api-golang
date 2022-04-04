package api

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type addSizeParams struct {
	Name string `json:"name" binding:"required"`
}

func (server *Server) addSize(ctx *gin.Context) {
	var req addSizeParams

	err := ctx.ShouldBindJSON(&req)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	size, err := server.store.AddSize(ctx, req.Name)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, size)
}

type deleteSizeParams struct {
	ID string `json:"id" binding:"required"`
}

func (server *Server) deleteSize(ctx *gin.Context) {
	var req deleteSizeParams

	err := ctx.ShouldBindJSON(&req)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	intId, err := strconv.Atoi(req.ID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	//this does not throw an error if id is not found
	err = server.store.DeleteSize(ctx, int32(intId))

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, map[string]string{
		"message": "size deleted successfully",
	})
}

func (server *Server) listAllSizes(ctx *gin.Context) {
	categories, err := server.store.ListAllSizes(ctx)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, categories)
}
