// @title			Account Management
// @version			1.0
// @description     This is a account management server which can create an account and verify credential.
// @termsOfService  http://swagger.io/terms/

// @contact.name   	Kinwang
// @contact.url    	https://github.com/kinwang94
// @contact.email  	autumn4090@gmail.com

// @license.name  	Apache 2.0
// @license.url   	http://www.apache.org/licenses/LICENSE-2.0.html

package main

import (
	"acct-mgmt/api"
	"acct-mgmt/db"
	"log"

	_ "acct-mgmt/docs"
)

func main() {
	collection, err := db.NewDB()
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	api.NewAPI(collection).Run()
}
