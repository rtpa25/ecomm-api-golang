package api

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/session"
)

func checkAuthority(server Server) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sess := session.GetSessionFromRequestContext(ctx)
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
		if userFromSelfDB.IsAdmin == false {
			ctx.Abort()
		}
		ctx.Next()
	}
}
