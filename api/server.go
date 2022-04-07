package api

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	db "github.com/rtpa25/ecomm-api-go/db/sqlc"
	"github.com/rtpa25/ecomm-api-go/utils"
	"github.com/supertokens/supertokens-golang/supertokens"
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
	router := gin.New()

	router.Use(gin.Recovery())

	// CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "DELETE", "PUT", "OPTIONS"},
		AllowHeaders:     append([]string{"content-type"}, supertokens.GetAllCORSHeaders()...),
		MaxAge:           1 * time.Minute,
		AllowCredentials: true,
	}))

	// Adding the SuperTokens middleware
	router.Use(func(c *gin.Context) {
		supertokens.Middleware(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			c.Next()
		})).ServeHTTP(c.Writer, c.Request)
		// we call Abort so that the next handler in the chain is not called, unless we call Next explicitly
		c.Abort()
	})

	router.GET("/hi", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "Hi there")
	})
	router.POST("/addProduct", server.verifySession(nil), server.checkAuthority(), server.addProduct) //route to add products
	router.DELETE("/deleteProduct", server.deleteProduct)                                             //route to delete product
	router.PATCH("/updateProduct", server.updateProduct)                                              //route to update product
	router.POST("/addCategory", server.addCategory)                                                   //route to add a category
	router.DELETE("/deleteCategory", server.deleteCategory)                                           //route to delete a category
	router.GET("/listAllCategories", server.listAllCategories)                                        //route to list all categories
	router.POST("/addSize", server.addSize)                                                           //route to add a size
	router.DELETE("/deleteSize", server.deleteSize)                                                   //route to delete a size
	router.GET("/listAllSizes", server.listAllSizes)                                                  //route to list all size
	router.GET("/listProducts", server.listProducts)                                                  //route to fetch all products in a paginated fashion
	router.GET("/getProduct", server.getProduct)                                                      //route to get a single product details by id
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
