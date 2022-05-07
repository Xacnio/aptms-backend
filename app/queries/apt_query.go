package queries

import (
	"github.com/Xacnio/aptms-backend/app/models"
	"github.com/jmoiron/sqlx"
	"log"
)

type AptQueries struct {
	*sqlx.DB
}

func (q *AptQueries) GetBuildings(page uint16, count uint8) ([]models.Building, error) {
	if page < 0 {
		page = 0
	}
	page = page * uint16(count)
	if count < 1 {
		count = 1
	} else if count > 200 {
		count = 200
	}
	result := make([]models.Building, 0)
	query := `
	SELECT 
		buildings.*, cities.city_name, districts.district_name 
	FROM 
	   	buildings 
	    	LEFT JOIN cities ON cities.id = buildings.city_id 
	   		LEFT JOIN districts ON districts.city_id = buildings.city_id AND districts.id = buildings.district_id 
	LIMIT $2 
	OFFSET $1`
	udb := q.Unsafe()
	rows, err := udb.Queryx(query, page, count)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		building := new(models.Building)
		err := rows.StructScan(&building)
		if err != nil {
			log.Fatalln(err)
		}
		result = append(result, *building)
	}
	return result, nil
}

func (q *AptQueries) GetBuildingById(id uint16) (models.Building, error) {
	var result models.Building
	query := `
	SELECT 
		buildings.*, cities.city_name, districts.district_name 
	FROM 
	    buildings 
	    	LEFT JOIN cities ON cities.id = buildings.city_id 
	        LEFT JOIN districts ON districts.city_id = buildings.city_id AND districts.id = buildings.district_id 
	WHERE 
	    buildings.id = $1`
	udb := q.Unsafe()
	row := udb.QueryRowx(query, id)
	if row.Err() != nil {
		return result, row.Err()
	}
	row.StructScan(&result)
	return result, nil
}

func (q *AptQueries) UpdateBuilding(building models.Building) error {
	query := `
	UPDATE 
	    buildings 
	SET 
	    name=$2,city_id=$3,district_id=$4,phone_number=$5,tax_number=$6,address=$7 
	WHERE 
		id = $1`
	_, err := q.Exec(query, building.ID, building.Name, building.City.ID, building.District.ID, building.PhoneNumber, building.TaxNumber, building.Address)
	if err != nil {
		return err
	}
	return nil
}

func (q *AptQueries) GetBlocks(buildingId int) ([]models.Block, error) {
	result := make([]models.Block, 0)
	query := `SELECT * FROM blocks WHERE building_id = $1 LIMIT 50`
	rows, err := q.Queryx(query, buildingId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		block := new(models.Block)
		err := rows.StructScan(&block)
		if err != nil {
			log.Fatalln(err)
		}
		result = append(result, *block)
	}
	return result, nil
}

func (q *AptQueries) NewBuildingBlock(block models.Block) error {
	query := `INSERT INTO blocks (building_id, letter, d_number) VALUES ($1, $2, $3)`
	_, err := q.Exec(query, block.BuildingID, block.Letter, block.DNumber)
	if err != nil {
		return err
	}
	return nil
}

func (q *AptQueries) NewBuilding(building models.Building) error {
	query := `INSERT INTO buildings (name, city_id, district_id, phone_number, tax_number, address) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := q.Exec(query, building.Name, building.City.ID, building.District.ID, building.PhoneNumber, building.TaxNumber, building.Address)
	if err != nil {
		return err
	}
	return nil
}

func (q *AptQueries) DeleteBuilding(buildingId int) error {
	q.Exec("DELETE FROM user_building_membership WHERE building_id = $1", buildingId)
	q.Exec("DELETE FROM revenues WHERE building_id = $1", buildingId)
	q.Exec("DELETE FROM flats WHERE building_id = $1", buildingId)
	q.Exec("DELETE FROM blocks WHERE building_id = $1", buildingId)
	_, err := q.Exec("DELETE FROM buildings WHERE id = $1", buildingId)
	if err != nil {
		return err
	}
	return nil
}

func (q *AptQueries) DeleteBuildingBlock(buildingId int, letter string) error {
	_, err := q.Exec("DELETE FROM blocks WHERE building_id = $1 AND letter LIKE $2", buildingId, letter)
	if err != nil {
		return err
	}
	return nil
}

func (q *AptQueries) SystemInfo() models.SystemInfo {
	row := q.QueryRowx(`
		SELECT 
			(SELECT COUNT(*) FROM buildings) as total_buildings,
			(SELECT COUNT(*) FROM flats) as total_flats,
			(SELECT COUNT(*) FROM users) as total_users
		`)
	data := models.SystemInfo{}
	if row.Err() != nil {
		return data
	}
	row.StructScan(&data)
	return data
}
