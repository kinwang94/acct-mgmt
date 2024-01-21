package api

import (
	"acct-mgmt/model"
	"acct-mgmt/services"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type API struct {
	collection *mongo.Collection
}

func NewAPI(collection *mongo.Collection) *API {
	return &API{collection: collection}
}

// TODO:
// 1. HTTP status code confirmation
// 2. Error structure design

// CreateAccountHandler godoc
// @Summary		Create an account
// @Description Create an account with the required username and password.
// @Description	The username must meet the following criteria:
// @Description	- Minimum length of 3 characters and maximum length of 32 characters.
// @Description
// @Description	The password must meet the following criteria:
// @Description	- Minimum length of 8 characters and maximum length of 32 characters.
// @Description	- Must contain at least 1 uppercase letter, 1 lowercase letter, and 1 number.
// @Tags		account
// @Accept		json
// @Produce		json
// @Param		account		body		model.RequestPayload		true	"Account credential"
// @Success		201			{object}	model.ResponsePayload
// @Failure		400			{object}	model.ResponsePayload
// @Router		/v1/signup [post]
func (a *API) CreateAccountHandler(c *gin.Context) {
	var req model.RequestPayload

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, &model.ResponsePayload{
			Success: false,
			Reason:  err.Error(),
		})
		return
	}

	if err := services.CreateAccount(context.Background(), a.collection, &model.Account{
		Username:       req.Username,
		Password:       req.Password,
		FailedAttempts: 0,
	}); err != nil {
		c.JSON(http.StatusBadRequest, &model.ResponsePayload{
			Success: false,
			Reason:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, model.ResponsePayload{
		Success: true,
		Reason:  "",
	})
}

// VerifyCredentialHandler godoc
// @Summary		Verify account credential
// @Description	Verify the provided account credentials.
// @Description If verification failed more than five times, the user is required to wait one minute before attempting again.
// @Tags		account
// @Accept		json
// @Produce		json
// @Param		account		body		model.RequestPayload		true	"Account credential"
// @Success		200			{object}	model.ResponsePayload
// @Failure		400			{object}	model.ResponsePayload
// @Router		/v1/login [post]
func (a *API) VerifyCredentialHandler(c *gin.Context) {
	var req model.RequestPayload

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, &model.ResponsePayload{
			Success: false,
			Reason:  err.Error(),
		})
		return
	}

	if err := services.VerifyCredential(context.Background(), a.collection, &model.Account{
		Username: req.Username,
		Password: req.Password,
	}); err != nil {
		c.JSON(http.StatusBadRequest, &model.ResponsePayload{
			Success: false,
			Reason:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &model.ResponsePayload{
		Success: true,
		Reason:  "",
	})
}
