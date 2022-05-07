package queries

import (
	"github.com/Xacnio/aptms-backend/app/models"
	"github.com/jmoiron/sqlx"
	"log"
)

type CDistQueries struct {
	*sqlx.DB
}

func (q *CDistQueries) GetCities() ([]models.City, error) {
	result := make([]models.City, 0)
	query := `SELECT id as city_id, city_name FROM cities ORDER BY city_name`
	udb := q.Unsafe()
	rows, err := udb.Queryx(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		city := new(models.City)
		err := rows.StructScan(&city)
		if err != nil {
			log.Fatalln(err)
		}
		result = append(result, *city)
	}
	return result, nil
}

func (q *CDistQueries) GetDistricts(cityId uint8) ([]models.District, error) {
	result := make([]models.District, 0)
	query := `SELECT id as district_id, district_name FROM districts WHERE city_id = $1 ORDER BY district_name`
	udb := q.Unsafe()
	rows, err := udb.Queryx(query, cityId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		dist := new(models.District)
		err := rows.StructScan(&dist)
		if err != nil {
			log.Fatalln(err)
		}
		result = append(result, *dist)
	}
	return result, nil
}
