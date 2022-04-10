package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"sort"
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
type listProductReponseParams struct {
	Reuslt []getProductReposneParams `json:"result"`
}

func (server *Server) listProducts(ctx *gin.Context) {
	var res listProductReponseParams

	val := ctx.Request.URL.Query()
	intLimit, err := strconv.Atoi(val["limit"][0])

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	intOffset, err := strconv.Atoi(val["offset"][0])

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	sortCriteria := val.Get("sort")

	filterCriteria := val.Get("filters")

	var filters map[string]string
	json.Unmarshal([]byte(filterCriteria), &filters)

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

	switch sortCriteria {
	case "Price (desc)":
		sort.Slice(products, func(i, j int) bool {
			intPricei, _ := strconv.Atoi(products[i].Price)
			intPricej, _ := strconv.Atoi(products[j].Price)
			return intPricei > intPricej
		})
	case "Price (asc)":
		sort.Slice(products, func(i, j int) bool {
			intPricei, _ := strconv.Atoi(products[i].Price)
			intPricej, _ := strconv.Atoi(products[j].Price)
			return intPricei < intPricej
		})
	case "Newest":
		sort.Slice(products, func(i, j int) bool {
			return products[i].CreatedAt.After(products[j].CreatedAt)
		})
	}

	for _, product := range products {
		sizes, err := server.store.GetAvailableSizesOfAProduct(ctx, product.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		hasSize := false
		for _, sizeOfProduct := range sizes {
			if filters["size"] != "" {
				if sizeOfProduct == filters["size"] {
					hasSize = true
				}
			}
		}

		if filters["size"] != "" {
			if !hasSize {
				break
			}
		}

		catagories, err := server.store.GetCategoriesOfAProduct(ctx, product.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		hasCatagory := false
		for _, catagoryOfProduct := range catagories {
			if filters["catagories"] != "" {
				if catagoryOfProduct == filters["catagories"] {
					hasCatagory = true
				}
			}
		}

		if filters["catagories"] != "" {
			if !hasCatagory {
				break
			}
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

type getProductReposneParams struct {
	Product        db.Product `json:"product"`
	AvailableSizes []string   `json:"available_sizes"`
	Catagories     []string   `json:"catagories"`
}

func (server *Server) getProduct(ctx *gin.Context) {
	id := ctx.Request.URL.Query().Get("id")

	var intId int
	var err error

	if id != "" {
		intId, err = strconv.Atoi(id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
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
