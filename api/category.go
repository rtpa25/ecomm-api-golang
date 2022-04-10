package api

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

//handler to add a catergory to the server -- admin only
type addCategoryParams struct {
	Name string `json:"name" binding:"required"`
}

func (server *Server) addCategory(ctx *gin.Context) {
	var req addCategoryParams

	err := ctx.ShouldBindJSON(&req)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	category, err := server.store.AddCategory(ctx, req.Name)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, category)
}

//handler to delete a catergory to the server -- admin only
type deleteCategoryParams struct {
	ID string `json:"id" binding:"required"`
}

func (server *Server) deleteCategory(ctx *gin.Context) {
	var req deleteCategoryParams

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
	err = server.store.DeleteCategory(ctx, int32(intId))

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, map[string]string{
		"message": "category deleted successfully",
	})
}

//handler to list all cateagories
func (server *Server) listAllCategories(ctx *gin.Context) {
	categories, err := server.store.ListAllCategory(ctx)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, categories)
}
