package api

import (
	"acct-mgmt/db"
	"acct-mgmt/errors"
	"acct-mgmt/model"
	"acct-mgmt/utils"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var r *gin.Engine

func TestMain(m *testing.M) {
	if err := godotenv.Load("../.env"); err != nil {
		panic("Failed to load .env file")
	}

	collection, err := db.NewDB()
	if err != nil {
		panic("Failed to connect to database")
	}
	r = NewAPI(collection).NewRouter()

	os.Exit(m.Run())
}

type args struct {
	username string
	password string
	attempts int
}

type expected struct {
	statusCode int
	success    bool
	reason     string
}

func TestSignup(t *testing.T) {
	existingUsername := utils.RandomStringAlphanumeric(3, 32)

	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "valid_account",
			args: args{
				username: existingUsername,
				password: utils.RandomStringAlphanumeric(8, 32),
			},
			expected: expected{
				statusCode: http.StatusCreated,
				success:    true,
				reason:     "",
			},
		}, {
			name: "invalid_account_username_already_exists",
			args: args{
				username: existingUsername,
				password: utils.RandomStringAlphanumeric(8, 32),
			},
			expected: expected{
				statusCode: http.StatusBadRequest,
				success:    false,
				reason:     errors.ErrUsernameAlreadyExists.Error(),
			},
		}, {
			name: "invalid_account_username_short",
			args: args{
				username: utils.RandomStringAlphanumeric(1, 2),
				password: utils.RandomStringAlphanumeric(8, 32),
			},
			expected: expected{
				statusCode: http.StatusBadRequest,
				success:    false,
				reason:     errors.ErrUsernameLengthTooShort.Error(),
			},
		}, {
			name: "invalid_account_username_long",
			args: args{
				username: utils.RandomStringAlphanumeric(33, 64),
				password: utils.RandomStringAlphanumeric(8, 32),
			},
			expected: expected{
				statusCode: http.StatusBadRequest,
				success:    false,
				reason:     errors.ErrUsernameLengthTooLong.Error(),
			},
		}, {
			name: "invalid_account_password_short",
			args: args{
				username: utils.RandomStringAlphanumeric(3, 32),
				password: utils.RandomStringAlphanumeric(1, 7),
			},
			expected: expected{
				statusCode: http.StatusBadRequest,
				success:    false,
				reason:     errors.ErrPasswordLengthTooShort.Error(),
			},
		}, {
			name: "invalid_account_password_long",
			args: args{
				username: utils.RandomStringAlphanumeric(3, 32),
				password: utils.RandomStringAlphanumeric(33, 64),
			},
			expected: expected{
				statusCode: http.StatusBadRequest,
				success:    false,
				reason:     errors.ErrPasswordLengthTooLong.Error(),
			},
		}, {
			name: "invalid_account_password_format",
			args: args{
				username: utils.RandomStringAlphanumeric(3, 32),
				password: utils.RandomStringLetter(8, 32),
			},
			expected: expected{
				statusCode: http.StatusBadRequest,
				success:    false,
				reason:     errors.ErrPasswordInvalidFormat.Error(),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Logf("username: %q and password: %q", test.args.username, test.args.password)

			w, resp, err := performRequest(r, "POST", "/v1/signup", gin.H{
				"username": test.args.username,
				"password": test.args.password,
			})
			assert.NoError(t, err)
			assert.Equal(t, test.expected.statusCode, w.Code)
			assert.Equal(t, test.expected.success, resp.Success)
			assert.Equal(t, test.expected.reason, resp.Reason)
		})
	}
}

func TestLogin(t *testing.T) {
	testUsername := utils.RandomStringAlphanumeric(3, 32)
	testPassword := utils.RandomStringAlphanumeric(8, 32)

	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "success",
			args: args{
				username: testUsername,
				password: testPassword,
				attempts: 1,
			},
			expected: expected{
				statusCode: http.StatusOK,
				success:    true,
				reason:     "",
			},
		}, {
			name: "failed_wrong_username",
			args: args{
				username: "unexistingUsername",
				password: testPassword,
				attempts: 1,
			},
			expected: expected{
				statusCode: http.StatusBadRequest,
				success:    false,
				reason:     errors.ErrNoAccountFound.Error(),
			},
		}, {
			name: "failed_wrong_password",
			args: args{
				username: testUsername,
				password: utils.RandomStringAlphanumeric(8, 32),
				attempts: 1,
			},
			expected: expected{
				statusCode: http.StatusBadRequest,
				success:    false,
				reason:     errors.ErrIncorrectCredential.Error(),
			},
		}, {
			name: "failed_attempts_not_exceeded",
			args: args{
				username: testUsername,
				password: utils.RandomStringAlphanumeric(8, 32),
				attempts: 4,
			},
			expected: expected{
				statusCode: http.StatusBadRequest,
				success:    false,
				reason:     errors.ErrIncorrectCredential.Error(),
			},
		}, {
			name: "failed_attempts_exceeded",
			args: args{
				username: testUsername,
				password: utils.RandomStringAlphanumeric(8, 32),
				attempts: 3, // Already attempted to log in 5 (1+4) times.
			},
			expected: expected{
				statusCode: http.StatusBadRequest,
				success:    false,
				reason:     errors.ErrTooManyAttempts.Error(),
			},
		},
	}

	// Create an account and test different login scenarios with that account.
	_, _, err := performRequest(r, "POST", "/v1/signup", gin.H{
		"username": testUsername,
		"password": testPassword,
	})
	assert.NoError(t, err)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Logf("username: %q and password: %q", test.args.username, test.args.password)

			for i := 0; i < test.args.attempts; i++ {
				w, resp, err := performRequest(r, "POST", "/v1/login", gin.H{
					"username": test.args.username,
					"password": test.args.password,
				})
				assert.NoError(t, err)
				assert.Equal(t, test.expected.statusCode, w.Code)
				assert.Equal(t, test.expected.success, resp.Success)
				assert.Equal(t, test.expected.reason, resp.Reason)
			}
		})
	}
}

func performRequest(r http.Handler, method, path string, payload interface{}) (*httptest.ResponseRecorder, *model.ResponsePayload, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, nil, err
	}

	req, err := http.NewRequest(method, path, bytes.NewReader(body))
	if err != nil {
		return nil, nil, err
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var resp model.ResponsePayload
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		return w, &resp, err
	}

	return w, &resp, nil
}
