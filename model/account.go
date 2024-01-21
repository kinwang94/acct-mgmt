package model

import (
	"acct-mgmt/errors"
	"acct-mgmt/utils"
	"context"

	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Account struct {
	Username              string    `bson:"username"`
	Password              string    `bson:"password"`
	FailedAttempts        int       `bson:"failed_attempts"`
	LastFailedAttemptTime time.Time `bson:"last_failed_attempt_time"`
}

// Validate checks whether the username and password of the account are valid.
func (ac *Account) Validate() error {
	if err := utils.ValidateUsername(ac.Username); err != nil {
		return err
	}

	if err := utils.ValidatePassword(ac.Password); err != nil {
		return err
	}

	return nil
}

// FindFromDatabaseByName finds the account record from database using username as filter.
func (ac *Account) FindFromDatabaseByName(ctx context.Context, collection *mongo.Collection) (*Account, error) {
	var account *Account

	filter := bson.M{"username": ac.Username}
	if err := collection.FindOne(ctx, filter).Decode(&account); err != nil {
		return nil, err
	}

	return account, nil
}

// InsertToDatabase inserts a new account record into database.
func (ac *Account) InsertToDatabase(ctx context.Context, collection *mongo.Collection) error {
	if _, err := collection.InsertOne(ctx, ac); err != nil {
		return err
	}

	return nil
}

// CheckFailedAttempts checks whether the number of failed login attempts exceeds the limitation.
func (ac *Account) CheckFailedAttempts(ctx context.Context, collection *mongo.Collection) error {
	if ac.FailedAttempts >= 5 {
		if time.Since(ac.LastFailedAttemptTime) < time.Minute {
			return errors.ErrTooManyAttempts
		}
	}

	return nil
}

// UpdateFailedAttempts updates the number of failed attempts and the time of the last
// failed attempt in the account record.
func (ac *Account) UpdateFailedAttempts(ctx context.Context, collection *mongo.Collection, failedAttempts int) error {
	ac.FailedAttempts = failedAttempts
	if failedAttempts != 0 {
		ac.LastFailedAttemptTime = time.Now()
	}

	filter := bson.M{"username": ac.Username}
	update := bson.M{"$set": bson.M{
		"failed_attempts":          ac.FailedAttempts,
		"last_failed_attempt_time": ac.LastFailedAttemptTime,
	}}

	if _, err := collection.UpdateOne(ctx, filter, update); err != nil {
		return err
	}

	return nil
}
