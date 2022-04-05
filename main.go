package main

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/lib/pq" //the underscore triggers the lexer to not remove the import if not used
	"github.com/rtpa25/ecomm-api-go/api"
	db "github.com/rtpa25/ecomm-api-go/db/sqlc"
	"github.com/rtpa25/ecomm-api-go/utils"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/supertokens"
)

var store db.Store
var config utils.Config

func main() {
	var err error
	config, err = utils.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config vars", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal(err.Error(), "Cannot connect to database")
	}

	store = db.NewStore(conn)
	err = supertokens.Init(supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "https://try.supertokens.io",
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "SuperTokens Demo App",
			APIDomain:     config.ServerAddress,
			WebsiteDomain: config.WebsiteAddress,
		},
		RecipeList: []supertokens.Recipe{
			session.Init(nil),
			emailpassword.Init(&epmodels.TypeInput{
				SignUpFeature: &epmodels.TypeInputSignUp{
					FormFields: []epmodels.TypeInputFormField{
						{
							ID: "username",
						},
					},
				},
				Override: &epmodels.OverrideStruct{
					APIs: func(originalImplementation epmodels.APIInterface) epmodels.APIInterface {
						originalSignupPost := *originalImplementation.SignUpPOST
						*originalImplementation.SignUpPOST = func(formFields []epmodels.TypeFormField, options epmodels.APIOptions, userContext supertokens.UserContext) (epmodels.SignUpPOSTResponse, error) {
							var username string
							for _, formField := range formFields {
								if formField.ID == "username" {
									username = formField.Value
								}
							}
							//insert the user into the personal db
							res, err := originalSignupPost(formFields, options, userContext)
							store.CreateUser(context.Background(), db.CreateUserParams{
								Username: username,
								IsAdmin:  false,
							})
							return res, err
						}
						return originalImplementation
					},
				},
			}),
		},
	})

	if err != nil {
		log.Fatal("Could not start SuperTokens")
	}
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("Could not loadup server", err)
	}

	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Fatal("Cannot start server")
	}
}
