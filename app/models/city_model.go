package models

import "database/sql"

type City struct {
	ID   uint8          `json:"cityId"  db:"city_id"`
	Name sql.NullString `json:"cityName"  db:"city_name"`
}

type District struct {
	ID   uint16         `json:"districtId"  db:"district_id"`
	Name sql.NullString `json:"districtName"  db:"district_name"`
}
