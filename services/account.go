package services

import (
	"acct-mgmt/errors"
	"acct-mgmt/model"
	"acct-mgmt/utils"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
)

// CreateAccount creates and stores the account record in the database.
func CreateAccount(ctx context.Context, collection *mongo.Collection, account *model.Account) error {
	if err := account.Validate(); err != nil {
		return err
	}

	// Check if the username already exists in the database.
	if result, err := account.FindFromDatabaseByName(ctx, collection); result != nil {
		return errors.ErrUsernameAlreadyExists
	} else if err != nil && err != mongo.ErrNoDocuments {
		return err
	}

	// Hash the password before adding to the database.
	hashedPassword, err := utils.HashPassword(account.Password)
	if err != nil {
		return err
	}
	account.Password = hashedPassword

	// Add the record to database.
	if err := account.InsertToDatabase(ctx, collection); err != nil {
		return err
	}

	return nil
}

// VerifyCredential verifies that the credential is correct and updates the failed attempt record.
func VerifyCredential(ctx context.Context, collection *mongo.Collection, account *model.Account) error {
	result, err := account.FindFromDatabaseByName(ctx, collection)
	if result == nil {
		return errors.ErrNoAccountFound
	}
	if err != nil {
		return err
	}

	// Check whether the number of failed login attempts exceeds the limitation.
	if err := result.CheckFailedAttempts(ctx, collection); err != nil {
		return err
	}

	if err := utils.VerifyPassword(account.Password, result.Password); err != nil {
		// Increase the number of failed attempts if verification failed.
		if err := result.UpdateFailedAttempts(ctx, collection, result.FailedAttempts+1); err != nil {
			log.Println("Failed to update account failed attempts: ", err)
		}

		return errors.ErrIncorrectCredential
	}

	// Reset the number of failed attempts if verification succeeded.
	if err := result.UpdateFailedAttempts(ctx, collection, 0 /* reset */); err != nil {
		log.Println("Failed to update account failed attempts: ", err)
	}

	return nil
}
