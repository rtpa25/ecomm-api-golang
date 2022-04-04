package api

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	db "github.com/rtpa25/ecomm-api-go/db/sqlc"
)

// ------------------ handler to add a product {admin only} ---------------------- //
type addProductRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	ImageUrl    string `json:"image_url" binding:"required"`
	ImageID     string `json:"image_id" binding:"required"`
	Price       string `json:"price" binding:"required"`
}

func (server *Server) addProduct(ctx *gin.Context) {
	var req addProductRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.AddProductParams{
		Name:        req.Name,
		Description: req.Description,
		ImageUrl:    req.ImageUrl,
		ImageID:     req.ImageID,
		Price:       req.Price,
	}

	product, err := server.store.AddProduct(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, product)
}

// ------------------ handler to delete a product {admin only} ---------------------- //
type deleteProductParams struct {
	Id string `json:"id" binding:"required"`
}

func (server *Server) deleteProduct(ctx *gin.Context) {
	var req deleteProductParams

	err := ctx.ShouldBindJSON(&req)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	intId, err := strconv.Atoi(req.Id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = server.store.DeleteProduct(ctx, int32(intId))

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, map[string]string{
		"message": "product deleted successfully",
	})
}

// ------------------ handler to update a product {admin only} ---------------------- //
type updateProductParams struct {
	ID          string `json:"id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Price       string `json:"price" binding:"required"`
	ImageUrl    string `json:"image_url" binding:"required"`
	ImageID     string `json:"image_id" binding:"required"`
}

func (server *Server) updateProduct(ctx *gin.Context) {
	var req updateProductParams

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

	arg := db.UpdateProductParams{
		ID:          int32(intId),
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		ImageUrl:    req.ImageUrl,
		ImageID:     req.ImageID,
	}

	product, err := server.store.UpdateProduct(ctx, arg)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, product)
}

// ------------------ handler to list all products with pagination ---------------------- //
type listProductParams struct {
	Limit  string `json:"limit"`  //the number of results you want
	Offset string `json:"offset"` //the number of entries you want to ommit before the result set
}

func (server *Server) listProducts(ctx *gin.Context) {
	var req listProductParams

	err := ctx.ShouldBindJSON(&req)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	intLimit, err := strconv.Atoi(req.Limit)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	intOffset, err := strconv.Atoi(req.Offset)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.ListProductsParams{
		Limit:  int32(intLimit),
		Offset: int32(intOffset),
	}

	products, err := server.store.ListProducts(ctx, arg)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, products)
}

// ------------------ handler to get a specific product  ---------------------- //
type getProductParams struct {
	Id string `json:"id" binding:"required"`
}

func (server *Server) getProduct(ctx *gin.Context) {
	var req getProductParams

	err := ctx.ShouldBindJSON(&req)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	intId, err := strconv.Atoi(req.Id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	product, err := server.store.GetProduct(ctx, int32(intId))

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, product)
}
