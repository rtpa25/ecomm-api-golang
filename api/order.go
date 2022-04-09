package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	db "github.com/rtpa25/ecomm-api-go/db/sqlc"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/session"
)

type createOrderRequestParams struct {
	Quantity     int32  `json:"quantity"`
	Address      string `json:"address"`
	ProdcutID    int32  `json:"prodcut_id"`
	SelectedSize string `json:"selected_size"`
}

func (server *Server) createOrder(ctx *gin.Context) {
	var req createOrderRequestParams

	err := ctx.ShouldBindJSON(&req)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	sess := session.GetSessionFromRequestContext(ctx.Request.Context())
	user, err := emailpassword.GetUserByID(sess.GetUserID())
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	userFromSelfDB, err := server.store.GetUserByEmail(ctx, user.Email)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	userId := userFromSelfDB.ID

	createdOrder, err := server.store.AddOrder(ctx, db.AddOrderParams{
		Quantity:     req.Quantity,
		UserID:       userId,
		Address:      req.Address,
		ProdcutID:    req.ProdcutID,
		SelectedSize: req.SelectedSize,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, createdOrder)
}

func (server *Server) getSelfOrder(ctx *gin.Context) {
	sess := session.GetSessionFromRequestContext(ctx.Request.Context())
	user, err := emailpassword.GetUserByID(sess.GetUserID())
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	userFromSelfDB, err := server.store.GetUserByEmail(ctx, user.Email)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	userId := userFromSelfDB.ID

	orders, err := server.store.GetOrdersForUser(ctx, userId)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, orders)
}

func (server *Server) getOrderForAnyUser(ctx *gin.Context) {
	userId := ctx.Request.URL.Query().Get("id")
	if userId == "" {
		ctx.JSON(http.StatusBadRequest, map[string]string{
			"message": "please pass the required userId in the url",
		})
		return
	}

	intUserId, err := strconv.Atoi(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	orders, err := server.store.GetOrdersForUser(ctx, int32(intUserId))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, orders)
}

type updateSelfOrderReqParams struct {
	ID           int32  `json:"id"`
	Quantity     int32  `json:"quantity"`
	SelectedSize string `json:"selected_size"`
	Address      string `json:"address"`
}

func (server *Server) updateSelfOrder(ctx *gin.Context) {
	var req updateSelfOrderReqParams

	err := ctx.ShouldBindJSON(&req)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	sess := session.GetSessionFromRequestContext(ctx.Request.Context())
	user, err := emailpassword.GetUserByID(sess.GetUserID())
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	userFromSelfDB, err := server.store.GetUserByEmail(ctx, user.Email)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	userId := userFromSelfDB.ID

	order, err := server.store.GetOrderById(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if order.UserID != userId {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	updatedOrder, err := server.store.UpdateOrderForUser(ctx, db.UpdateOrderForUserParams{
		ID:           req.ID,
		Quantity:     req.Quantity,
		SelectedSize: req.SelectedSize,
		Address:      req.Address,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, updatedOrder)
}

type deleteSelfOrderRequestParams struct {
	OrderId int32 `json:"order_id"`
}

func (server *Server) deleteSelfOrder(ctx *gin.Context) {
	var req deleteSelfOrderRequestParams

	err := ctx.ShouldBindJSON(&req)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	sess := session.GetSessionFromRequestContext(ctx.Request.Context())
	user, err := emailpassword.GetUserByID(sess.GetUserID())
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	userFromSelfDB, err := server.store.GetUserByEmail(ctx, user.Email)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	userId := userFromSelfDB.ID

	order, err := server.store.GetOrderById(ctx, req.OrderId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if order.UserID != userId {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	err = server.store.DeleteOrderById(ctx, req.OrderId)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, map[string]string{
		"message": "order deleted successfully",
	})
}
