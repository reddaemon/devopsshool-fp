package repository

import (
	"encoding/json"
	"final-project/internal/models"
	"log"
	"strings"
	"time"

	"context"
	"errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/go-redis/redis/v8"
)

type Instance struct {
	Db *pgxpool.Pool
	redisConn *redis.Client
	
}

func NewInstance(db *pgxpool.Pool, rc *redis.Client) *Instance {
	return &Instance{Db: db, redisConn: rc}
}

func (i *Instance) Insert(ctx context.Context, currency []models.ValCurs) error {
	const (
		layoutISO = "02.01.2006"
		layoutUS  = "02-Jan-2006"
	)

	for number := range currency {
		log.Println("RECORD:", currency[number].Valute)
		for _, item := range currency[number].Valute {
			log.Println("DATE: ", currency[number].Date)
			log.Println("VALUTE ID: ", item.ID)
			log.Println("NUMCODE: ", item.NumCode)
			log.Println("CHARCODE: ", item.CharCode)
			log.Println("NOMINAL: ", item.Nominal)
			log.Println("VALUE: ", item.Value)
			log.Println("NAME: ", item.Name)

			date := currency[number].Date
			log.Println("date: ", date)
			convertdatetmp, err := time.Parse(layoutISO, date)
			log.Println("convertdatetmp: ", convertdatetmp)
			if err != nil {
				log.Println(err)
			}
			convertdate := convertdatetmp.Format(layoutUS)
			log.Println("date before insert: ", convertdate)
			value := item.Value
			convertvalue := strings.Replace(value, ",", ".", 1)
			query := `INSERT INTO currency 
	(date_of_request, valute_id, numcode, charcode, nominal, value, name)
	 VALUES ($1, $2, $3, $4, $5, $6, $7)`
			commandTag, err := i.Db.Exec(ctx, query,
				convertdate,
				item.ID,
				item.NumCode,
				item.CharCode,
				item.Nominal,
				convertvalue,
				item.Name,
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
		}
	}

	return nil
}

func (i *Instance) Select(ctx context.Context) ([]models.ValCurs, error) {
	query := `SELECT date_of_request::text, json_agg(json_build_object(
		'valute_id', valute_id,
		'numcode',  numcode::text, 
		'charcode', charcode::text, 
		'nominal', nominal::text,
		'value', value::text,
		'name', name::text)) FROM currency GROUP BY date_of_request ORDER by date_of_request;`

	var valcurs []models.ValCurs

	rows, err := i.Db.Query(ctx, query)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			log.Println(pgErr.Message)
			log.Println(pgErr.Code)
			return nil, err
		}
	}

	defer rows.Close()

	for rows.Next() {
		currency := models.ValCurs{}
		var currencytmp string
		rows.Scan(&currency.Date, &currencytmp)

		err := json.Unmarshal([]byte(currencytmp), &currency.Valute)
		if err != nil {
			log.Printf("%#v", err)
		}

		valcurs = append(valcurs, currency)
	}

	return valcurs, nil

}
