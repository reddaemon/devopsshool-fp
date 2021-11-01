package models

import (
	"encoding/xml"
)

type ValCurs struct {
	XMLName xml.Name `xml:"ValCurs" json:"-"`
	Text    string   `xml:",chardata" json:"-"`
	Date    string   `xml:"Date,attr"`
	Name    string   `xml:"name,attr" json:"-"`
	Valute  []Valute `xml:"Valute"`
}

type Valute struct {
	Text     string `xml:",chardata" json:"-"`
	ID       string `xml:"ID,attr" json:"Valute_ID"`
	NumCode  string `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Nominal  string `xml:"Nominal"`
	Name     string `xml:"Name"`
	Value    string `xml:"Value"`
}
