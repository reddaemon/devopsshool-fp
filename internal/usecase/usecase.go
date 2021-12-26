package usecase

import (
	"context"
	"encoding/xml"
	"final-project/internal/models"
	"final-project/internal/repository"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/google/uuid"
	"golang.org/x/net/html/charset"
)

type UsecaseInterface interface {
	GetData()
}

type Usecase struct {
	repo *repository.Instance
}

func NewUseCase(repo *repository.Instance) *Usecase {
	return &Usecase{repo: repo}
}

func (u *Usecase) PullDataByPeriod(startdate string) {
	const (
		layoutISO = "02/01/2006"
		layoutUS  = "02-Jan-2006"
	)
	convertstartdatetmp, err := time.Parse(layoutISO, startdate)
	if err != nil {
		log.Println(err)
	}

	convertstrartdate := convertstartdatetmp.Format(layoutISO)
	end := convertstartdatetmp.AddDate(0, -1, 0)
	convertenddate := end.Format(layoutISO)
	log.Println(convertstrartdate, convertenddate)

	ctx := context.Background()

	for d := convertstartdatetmp; !d.Before(end); d = d.AddDate(0, 0, -1) {
		log.Printf("process date %#v", d.Format(layoutISO))
		url := fmt.Sprintf("http://www.cbr.ru/scripts/XML_daily.asp?date_req=%s", d.Format(layoutISO))
		parsed := ParseXml(url)
		err = u.repo.Insert(ctx, parsed)
		if err != nil {
			log.Printf("Insert error: %#v", err)
		}
	}
}

func ParseXml(url string) []models.ValCurs {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	var currency []models.ValCurs

	decoder := xml.NewDecoder(resp.Body)
	decoder.CharsetReader = charset.NewReaderLabel
	err = decoder.Decode(&currency)
	if err != nil {
		log.Fatalln(err)
	}

	return currency
}

func (u *Usecase) GetData() ([]models.ValCurs, error) {
	ctx := context.Background()
	result, err := u.repo.Select(ctx)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *Usecase) VerifyUser(ctx context.Context, user *models.User) (bool, error) {
	ok, err := u.repo.CheckCredentials(ctx, user)
	if !ok {
		log.Printf("invalid user creds: %s", err)
		return false, err
	}
	return true, nil
}

func (u *Usecase) AddUser(ctx context.Context, user *models.User) error {
	err := u.repo.AddUser(ctx, user)
	if err != nil {
		return err
	}
	return nil

}

func (u *Usecase) CreateToken(userid float64) (*models.TokenDetails, error) {
	td := &models.TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.AccessUuid = uuid.New().String()

	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUuid = uuid.New().String()
	var err error
	// Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["user_id"] = userid
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}
	// Creating Refresh Token
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = userid
	rtClaims["user_id"] = userid
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}
	return td, nil
}

func (u *Usecase) CreateAuth(ctx context.Context, userid float64, td *models.TokenDetails) error {
	err := u.repo.CreateAuth(ctx, userid, td)
	if err != nil {
		return err
	}
	return nil
}

func (u *Usecase) DeleteAuth(ctx context.Context, userid string) (int64, error) {
	userID, err := u.repo.DeleteAuth(ctx, userid)
	if err != nil {
		return 0, err
	}
	return userID, nil
}


