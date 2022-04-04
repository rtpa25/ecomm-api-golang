package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	db "github.com/rtpa25/ecomm-api-go/db/sqlc"
)

// ------------------ handler to add a product {admin only} ---------------------- //
type addProductRequestParams struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description" binding:"required"`
	ImageUrl    string   `json:"image_url" binding:"required"`
	ImageID     string   `json:"image_id" binding:"required"`
	Price       string   `json:"price" binding:"required"`
	CategoryIDs []string `json:"category_ids" binding:"required"`
	SizeIDs     []string `json:"size_ids" binding:"required"`
}

func (server *Server) addProduct(ctx *gin.Context) {
	var req addProductRequestParams

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.AddProductRequestParams{
		Name:        req.Name,
		Description: req.Description,
		ImageUrl:    req.ImageUrl,
		ImageID:     req.ImageID,
		Price:       req.Price,
		CategoryIDs: req.CategoryIDs,
		SizeIDs:     req.SizeIDs,
	}

	res, err := server.store.AddProductTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, res)
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

	//this does not throw an error if id is not found
	err = server.store.DeleteProduct(ctx, int32(intId))

	if err != nil {
		fmt.Println("THIS IS COMING HERE")
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

type listProductReponseParams struct {
	Reuslt []getProductReposneParams `json:"result"`
}

func (server *Server) listProducts(ctx *gin.Context) {
	var req listProductParams
	var res listProductReponseParams

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

	for _, product := range products {
		sizes, err := server.store.GetAvailableSizesOfAProduct(ctx, product.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		catagories, err := server.store.GetCategoriesOfAProduct(ctx, product.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		entry := getProductReposneParams{
			Product:        product,
			AvailableSizes: sizes,
			Catagories:     catagories,
		}

		res.Reuslt = append(res.Reuslt, entry)
	}

	ctx.JSON(http.StatusOK, res)
}

// ------------------ handler to get a specific product  ---------------------- //
type getProductParams struct {
	Id string `json:"id" binding:"required"`
}

type getProductReposneParams struct {
	Product        db.Product `json:"product"`
	AvailableSizes []string   `json:"available_sizes"`
	Catagories     []string   `json:"catagories"`
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

	catagories, err := server.store.GetCategoriesOfAProduct(ctx, product.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	sizes, err := server.store.GetAvailableSizesOfAProduct(ctx, product.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := getProductReposneParams{
		Product:        product,
		AvailableSizes: sizes,
		Catagories:     catagories,
	}

	ctx.JSON(http.StatusOK, res)
}
