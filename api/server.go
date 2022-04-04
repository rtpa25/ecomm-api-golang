package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/rtpa25/ecomm-api-go/db/sqlc"
	"github.com/rtpa25/ecomm-api-go/utils"
)

//server set's up HTTP routing for our banking service
type Server struct {
	config utils.Config
	store  db.Store
	router *gin.Engine
}

//starts the server
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.GET("/hi", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "Hi there")
	})
	router.POST("/addProduct", server.addProduct)              //route to add products
	router.DELETE("/deleteProduct", server.deleteProduct)      //route to delete product
	router.PATCH("/updateProduct", server.updateProduct)       //route to update product
	router.POST("/addCategory", server.addCategory)            //route to add a category
	router.DELETE("/deleteCategory", server.deleteCategory)    //route to delete a category
	router.GET("/listAllCategories", server.listAllCategories) //route to list all categories
	router.POST("/addSize", server.addSize)                    //route to add a size
	router.DELETE("/deleteSize", server.deleteSize)            //route to delete a size
	router.GET("/listAllSizes", server.listAllSizes)           //route to list all size
	router.GET("/listProducts", server.listProducts)           //route to fetch all products in a paginated fashion
	router.GET("/getProduct", server.getProduct)               //route to get a single product details by id
	server.router = router
}

//newServer creates a new http server and set's up routing
func NewServer(config utils.Config, store db.Store) (*Server, error) {
	server := &Server{
		store:  store,
		config: config,
	}

	server.setupRouter()
	return server, nil
}

//error helper
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
