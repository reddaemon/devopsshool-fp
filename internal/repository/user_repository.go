package repository

import (
	"context"
	"errors"
	"final-project/internal/models"
	"log"

	"os"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"golang.org/x/crypto/bcrypt"
)

func (i *Instance) InsertUser(ctx context.Context, user *models.User) error {

	//TODO validate user

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	query := `INSERT INTO user 
	(email, password)
	 VALUES ($1, $2, $3)`

	commandTag, err := i.Db.Exec(ctx, query,
		user.Email,
		user.Password,
		user.Token,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			log.Println(pgErr.Message)
			log.Println(pgErr.Code)
			return err
		}
	}
	log.Println(commandTag.String())
	log.Println(commandTag.RowsAffected())

	// Create new JWT token for the newly registered account
	UserID, _ := uuid.NewRandom()
	tk := &models.Token{UserID: UserID}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("SECRET")))
	user.Token = tokenString

	user.Password = "" // delete password

	// TODO

	return nil
}
