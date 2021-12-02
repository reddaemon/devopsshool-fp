package repository

import (
	"context"
	"errors"
	"final-project/internal/models"
	"log"
	"strings"

	"os"

	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"golang.org/x/crypto/bcrypt"
)

func (i *Instance) PreCheck(ctx context.Context, user *models.User) bool {
	if !strings.Contains(user.Email, "@") {
		log.Println("Invalid email address")
		return false
	}

	query := `SELECT email, password FROM users WHERE email=$1`
	rows, err := i.Db.Query(ctx, query,
		user.Email,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			log.Println(pgErr.Message)
			log.Println(pgErr.Code)
		}
	}

	defer rows.Close()

	for rows.Next() {
		temp := &models.User{}
		rows.Scan(&temp.Email, &temp.Password)
		if temp.Email != "" {
			log.Println("Email is already exists")
			return false
		}
		err := bcrypt.CompareHashAndPassword([]byte(temp.Password), []byte(user.Password))
		if err != nil {
			log.Printf("invalid password: %s", err)
			return false
		}
	}

	return true

}

func (i *Instance) CheckCredentials(ctx context.Context, user *models.User) (bool, error) {
	log.Printf("userBeforeQuery: %#v", user)
	temp := &models.User{}

	query := `SELECT email, passwordhash FROM users WHERE email=$1`
	err := i.Db.QueryRow(ctx, query,
		user.Email,
	).Scan(&temp.Email, &temp.Password)
	if err != nil {
		return false, err
	}

	if temp.Email != user.Email {
		log.Printf("invalid email: %s", temp.Email)
		return false, errors.New("email doesn't exist")
	}
	err = bcrypt.CompareHashAndPassword([]byte(temp.Password), []byte(user.Password))
	if err != nil {
		log.Printf("invalid password: %s", err)
		return false, err
	}

	return true, nil

}

func (i *Instance) AddUser(ctx context.Context, user *models.User) error {

	if ok := i.PreCheck(ctx, user); !ok {
		return errors.New("invalid user check")
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	query := `INSERT INTO users 
	(email, password)
	 VALUES ($1, $2)`

	commandTag, err := i.Db.Exec(ctx, query,
		user.Email,
		user.Password,
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

	return nil
}

func (I *Instance) CreateAuth(ctx context.Context, userid int64, td *models.TokenDetails) error {
	at := time.Unix(td.AtExpires, 0) // converting Unix to UTC
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	err := I.redisConn.Set(ctx, td.AccessUuid, strconv.Itoa(int(userid)), at.Sub(now)).Err()
	if err != nil {
		log.Println("error set access uuid key")
		return err
	}

	err = I.redisConn.Set(ctx, td.RefreshUuid, strconv.Itoa(int(userid)), rt.Sub(now)).Err()
	if err != nil {
		log.Println("error set refresh uuid key")
		return err
	}
	return nil
}

func (I *Instance) FetchAuth(ctx context.Context, authD *models.AccessDetails) (uint64, error) {
	userid, err := I.redisConn.Get(ctx, authD.AccessUuid).Result()
	if err != nil {
		return 0, err
	}
	userID,_ := strconv.ParseUint(userid, 10, 64)
	return userID, nil
	
}

func (I *Instance) DeleteAuth(ctx context.Context, givenUuid string) (int64, error) {
	deleted, err := I.redisConn.Del(ctx, givenUuid).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}