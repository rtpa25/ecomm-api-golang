package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/session"
)

func (server *Server) checkAuthority() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sess := session.GetSessionFromRequestContext(ctx.Request.Context())
		user, err := emailpassword.GetUserByID(sess.GetUserID())
		if err != nil {
			log.Fatal(err.Error())
			ctx.Abort()
		}
		userFromSelfDB, err := server.store.GetUserByEmail(ctx, user.Email)
		if err != nil {
			log.Fatal(err.Error())
			ctx.Abort()
		}
		if !userFromSelfDB.IsAdmin {
			ctx.JSON(http.StatusUnauthorized, "unauthorized")
			ctx.Abort()
		}
		ctx.Next()
	}
}
