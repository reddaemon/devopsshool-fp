package usecase

import (
	"context"
	"encoding/xml"
	"final-project/internal/models"
	"final-project/internal/repository"
	"fmt"
	"log"
	"net/http"
	"time"

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
