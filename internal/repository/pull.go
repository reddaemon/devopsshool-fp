package repository

import (
	"encoding/xml"
	"final-project/internal/models"
	"log"
	"net/http"

	"golang.org/x/net/html/charset"
)

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
