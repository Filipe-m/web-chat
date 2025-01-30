package user_test

import (
	"errors"
	"testing"
	"time"
	"web-chat/internal/user"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestRepository_Create(t *testing.T) {
	cases := []struct {
		name     string
		req      *user.User
		mockFunc func(mock sqlmock.Sqlmock)
		wantErr  error
	}{
		{
			name: "should create user successfully",
			req: &user.User{
				Name:  "John Doe",
				Email: "johndoe@example.com",
			},
			mockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO users").
					WithArgs("John Doe", sqlmock.AnyArg(), "johndoe@example.com").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: nil,
		},
		{
			name: "should return error on insert failure",
			req: &user.User{
				Name:  "Jane Doe",
				Email: "janedoe@example.com",
			},
			mockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO users").
					WithArgs("Jane Doe", sqlmock.AnyArg(), "janedoe@example.com").
					WillReturnError(errors.New("insert error"))
			},
			wantErr: errors.New("insert error"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			repo := user.NewRepository(db)

			tc.mockFunc(mock)

			err = repo.Create(tc.req)

			assert.Equal(t, tc.wantErr, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
func TestRepository_GetUserByID(t *testing.T) {
	cases := []struct {
		name     string
		userID   int
		mockFunc func(mock sqlmock.Sqlmock)
		wantUser *user.User
		wantErr  error
	}{
		{
			name:   "should return user successfully",
			userID: 1,
			mockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "password", "email", "created_at", "updated_at"}).
						AddRow(1, "Jane Doe", "password123", "janedoe@example.com", time.Now(), time.Now()))
			},
			wantUser: &user.User{
				ID:    1,
				Name:  "Jane Doe",
				Email: "janedoe@example.com",
			},
			wantErr: nil,
		},
		{
			name:   "should return error if user not found",
			userID: 99,
			mockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT").
					WithArgs(99).
					WillReturnError(errors.New("user not found"))
			},
			wantUser: nil,
			wantErr:  errors.New("user not found"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			repo := user.NewRepository(db)
			tc.mockFunc(mock)

			savedUser := user.User{ID: tc.userID}
			err = repo.GetUserByID(&savedUser)

			if tc.wantErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.wantUser.Name, savedUser.Name)
				assert.Equal(t, tc.wantUser.Email, savedUser.Email)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRepository_GetUserByEmail(t *testing.T) {
	cases := []struct {
		name     string
		email    string
		mockFunc func(mock sqlmock.Sqlmock)
		wantUser *user.User
		wantErr  error
	}{
		{
			name:  "should return user successfully",
			email: "alice@example.com",
			mockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, name, email, password, created_at, updated_at FROM users WHERE email = ?").
					WithArgs("alice@example.com").
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "password", "created_at", "updated_at"}).
						AddRow(1, "Alice Doe", "alice@example.com", "password123", time.Now(), time.Now()))
			},
			wantUser: &user.User{
				ID:    1,
				Name:  "Alice Doe",
				Email: "alice@example.com",
			},
			wantErr: nil,
		},
		{
			name:  "should return error if user not found",
			email: "notfound@example.com",
			mockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, name, email, password, created_at, updated_at FROM users WHERE email = ?").
					WithArgs("notfound@example.com").
					WillReturnError(errors.New("user not found"))
			},
			wantUser: nil,
			wantErr:  errors.New("user not found"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			repo := user.NewRepository(db)
			tc.mockFunc(mock)

			savedUser, err := repo.GetUserByEmail(tc.email)

			if tc.wantErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.wantUser.Name, savedUser.Name)
				assert.Equal(t, tc.wantUser.Email, savedUser.Email)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
