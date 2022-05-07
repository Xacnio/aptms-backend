package queries

import (
	"database/sql"
	"github.com/Xacnio/aptms-backend/app/models"
	"log"
)

func (q *AptQueries) IsBuildingMember(userId, buildingId uint) (bool, error) {
	count := 0
	query := `SELECT COUNT(*) FROM user_building_membership WHERE user_id = $1 AND building_id = $2`
	row := q.QueryRowx(query, userId, buildingId)
	if row.Err() != nil {
		return false, row.Err()
	}
	err1 := row.Scan(&count)
	if err1 != nil {
		return false, err1
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

func (q *AptQueries) IsBuildingManager(userId, buildingId uint) (bool, error) {
	count := 0
	query := `SELECT COUNT(*) FROM user_building_membership WHERE user_id = $1 AND building_id = $2 AND rank = 1`
	row := q.QueryRowx(query, userId, buildingId)
	if row.Err() != nil {
		return false, row.Err()
	}
	err1 := row.Scan(&count)
	if err1 != nil {
		return false, err1
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

func (q *AptQueries) MGetBuildings(userId uint, page uint16, count uint8) ([]models.Building, error) {
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
	    a.*, cities.city_name, districts.district_name 
	FROM 
	     user_building_membership as uap 
	         RIGHT JOIN buildings as a ON uap.building_id = a.id 
	         LEFT JOIN cities ON cities."id" = a.city_id 
	         LEFT JOIN districts ON districts."id" = a.district_id 
	WHERE 
	      uap.user_id = $1 
	LIMIT $3 OFFSET $2`
	udb := q.Unsafe()
	rows, err := udb.Queryx(query, userId, page, count)
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

func (q *AptQueries) MGetFlats(buildingId int, page uint16, count uint8) ([]models.Flat, error) {
	if page < 0 {
		page = 0
	}
	page = page * uint16(count)
	if count < 1 {
		count = 1
	} else if count > 200 {
		count = 200
	}
	result := make([]models.Flat, 0)
	query := `
	SELECT 
	       f.*, b.letter "block_letter", concat(ow."name", ' ', ow.surname) "owname", concat(te."name", ' ',  te.surname) "tename" 
	FROM 
	     flats as f 
	         LEFT JOIN users as ow ON ow.id = f.owner_id 
	         LEFT JOIN blocks as b ON f.block_id = b.id 
	         LEFT JOIN users as te ON te.id = f.tenant_id 
	WHERE f.building_id = $1 
	ORDER BY b.letter,f.number 
	LIMIT $3 OFFSET $2`
	udb := q.Unsafe()
	rows, err := udb.Queryx(query, buildingId, page, count)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		flatDb := struct {
			ID          int
			Type        uint8
			BuildingID  int `db:"building_id"`
			BlockID     int `db:"block_id"`
			Number      string
			OwnerID     sql.NullInt32  `db:"owner_id"`
			TenantID    sql.NullInt32  `db:"tenant_id"`
			OwnerName   sql.NullString `db:"owname"`
			TenantName  sql.NullString `db:"tename"`
			BlockLetter sql.NullString `db:"block_letter"`
		}{}
		flat := new(models.Flat)
		err := rows.StructScan(&flatDb)
		if err != nil {
			log.Fatalln(err)
		}
		flat.Default()
		flat.ID = flatDb.ID
		flat.Type = flatDb.Type
		flat.BuildingID = flatDb.BuildingID
		flat.BlockID = flatDb.BlockID
		flat.Number = flatDb.Number
		if flatDb.BlockLetter.Valid {
			flat.BlockLetter = flatDb.BlockLetter.String
		}
		if flatDb.OwnerID.Valid {
			flat.OwnerID = int(flatDb.OwnerID.Int32)
		}
		if flatDb.TenantID.Valid {
			flat.TenantID = int(flatDb.TenantID.Int32)
		}
		if flatDb.OwnerName.Valid && len(flatDb.OwnerName.String) > 1 {
			flat.OwnerName = flatDb.OwnerName.String
		}
		if flatDb.TenantName.Valid && len(flatDb.TenantName.String) > 1 {
			flat.TenantName = flatDb.TenantName.String
		}
		result = append(result, *flat)
	}
	return result, nil
}

func (q *AptQueries) MGetBlocks(buildingId int) ([]models.Block, error) {
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

func (q *AptQueries) MNewBuildingFlat(flat models.Flat) error {
	var ownerId, tenantId sql.NullInt32
	if flat.OwnerID == 0 {
		ownerId.Valid = false
	} else {
		ownerId.Valid = true
		ownerId.Int32 = int32(flat.OwnerID)
	}
	if flat.TenantID == 0 {
		tenantId.Valid = false
	} else {
		tenantId.Valid = true
		tenantId.Int32 = int32(flat.TenantID)
	}
	query := `INSERT INTO flats (building_id, block_id, number, owner_id, tenant_id) VALUES ($1, $2, $3, $4, $5)`
	_, err := q.Exec(query, flat.BuildingID, flat.BlockID, flat.Number, ownerId, tenantId)
	if err != nil {
		return err
	}
	return nil
}

func (q *AptQueries) MDeleteBuildingFlat(buildingId, flatId uint) error {
	_, err := q.Exec("DELETE FROM flats WHERE building_id = $1 AND id = $2", buildingId, flatId)
	if err != nil {
		return err
	}
	return nil
}

func (q *AptQueries) MEditBuildingFlat(flat models.Flat) error {
	var ownerId, tenantId sql.NullInt32
	if flat.OwnerID == 0 {
		ownerId.Valid = false
	} else {
		ownerId.Valid = true
		ownerId.Int32 = int32(flat.OwnerID)
	}
	if flat.TenantID == 0 {
		tenantId.Valid = false
	} else {
		tenantId.Valid = true
		tenantId.Int32 = int32(flat.TenantID)
	}
	query := `UPDATE flats SET owner_id = $1, tenant_id = $2 WHERE id = $3 AND building_id = $4`
	_, err := q.Exec(query, ownerId, tenantId, flat.ID, flat.BuildingID)
	if err != nil {
		return err
	}
	return nil
}

func (q *AptQueries) MRevenuesList(buildingId int) ([]models.Revenue, error) {
	result := make([]models.Revenue, 0)
	query := `
	SELECT 
	    CONCAT(b.letter || '-' || f."number") as flat_name, r.* 
	FROM 
	     revenues as r 
	         LEFT JOIN flats as f ON f.id = r.flat_id 
	         LEFT JOIN blocks as b ON b.id = f.block_id 
	WHERE r.building_id = $1 
	ORDER BY r.rid DESC 
	LIMIT 50`
	rows, err := q.Queryx(query, buildingId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		revenue := new(models.Revenue)
		err := rows.StructScan(&revenue)
		if err != nil {
			log.Fatalln(err)
		}
		result = append(result, *revenue)
	}
	return result, nil
}

func (q *AptQueries) MNewRevenue(revenue *models.Revenue) error {
	query := `
	INSERT INTO 
	    revenues (building_id, flat_id, rid, total, paid_type, paid_time, payer_full_name, payer_email, payer_phone, paid_status, details) 
	VALUES ($1, $2, (select coalesce(max(rid) + 1,1) from revenues as rr where rr.building_id = $1), $3, 0, $4, $5, $6, $7, $8, $9) 
	RETURNING id, time, rid`
	row := q.QueryRowx(query, revenue.BuildingID, revenue.FlatID, revenue.Total, revenue.Paid.Time, revenue.Payer.FullName,
		revenue.Payer.Email, revenue.Payer.Phone, revenue.Paid.Status, revenue.Details)
	if row.Err() != nil {
		return row.Err()
	}
	if err := row.Scan(&revenue.ID, &revenue.Time, &revenue.RID); err != nil {
		return err
	}
	return nil
}

func (q *AptQueries) MEditRevenue(revenue *models.Revenue) error {
	query := `
	UPDATE revenues 
	SET paid_time = $1, payer_full_name = $2, payer_email = $3, payer_phone = $4, paid_status = $5, details = $6
	WHERE id = $7 AND paid_status = false
	RETURNING total`
	row := q.QueryRowx(query, revenue.Paid.Time, revenue.Payer.FullName, revenue.Payer.Email, revenue.Payer.Phone, revenue.Paid.Status, revenue.Details, revenue.ID)
	if row.Err() != nil {
		return row.Err()
	}
	if err := row.Scan(&revenue.Total); err != nil {
		return err
	}
	return nil
}

func (q *AptQueries) MBuildingDetails(buildingId int) (models.BuildingDetails, error) {
	query := `
	SELECT 
	    a.cash_amount, c.city_name as city_name, d.district_name as district_name,
	    (SELECT COUNT(*) as count FROM blocks WHERE building_id = $1) as block_count, 
		(SELECT COUNT(*) as count FROM flats WHERE building_id = $1) as flat_count,
		(SELECT COUNT(*) as count FROM user_building_membership WHERE building_id = $1) as user_count
	FROM 
	    buildings as a LEFT JOIN cities as c ON c.id = a.city_id LEFT JOIN districts as d ON d.id = a.district_id
	WHERE 
	    a.id = $1`
	row := q.QueryRowx(query, buildingId)
	data := models.BuildingDetails{}
	if row.Err() != nil {
		return data, row.Err()
	}
	err := row.StructScan(&data)
	if err != nil {
		return data, err
	}
	return data, err
}

func (q *AptQueries) MBuildingCashIncrease(buildingId int, amount float64) {
	query := `UPDATE buildings SET cash_amount = cash_amount + $2 WHERE id = $1`
	q.Exec(query, buildingId, amount)
}

func (q *AptQueries) MBuildingCashDecrease(buildingId int, amount float64) {
	query := `UPDATE buildings SET cash_amount = cash_amount - $2 WHERE id = $1`
	q.Exec(query, buildingId, amount)
}

func (q *AptQueries) MExpensesList(buildingId int) ([]models.Expense, error) {
	result := make([]models.Expense, 0)
	query := `SELECT r.* FROM expenses as r WHERE r.building_id = $1 ORDER BY r.eid DESC LIMIT 50`
	rows, err := q.Queryx(query, buildingId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		expense := new(models.Expense)
		err := rows.StructScan(&expense)
		if err != nil {
			log.Fatalln(err)
		}
		result = append(result, *expense)
	}
	return result, nil
}

func (q *AptQueries) MNewExpense(expense *models.Expense) error {
	query := `
	INSERT INTO 
	    expenses (building_id, eid, total, paid_type, paid_time, to_name, to_email, to_phone, paid_status, details) 
	VALUES ($1, (select coalesce(max(eid) + 1,1) from expenses as rr where rr.building_id = $1), $2, 0, $3, $4, $5, $6, $7, $8) 
	RETURNING id, time, eid`
	row := q.QueryRowx(query, expense.BuildingID, expense.Total, expense.Paid.Time, expense.To.FullName,
		expense.To.Email, expense.To.Phone, expense.Paid.Status, expense.Details)
	if row.Err() != nil {
		return row.Err()
	}
	if err := row.Scan(&expense.ID, &expense.Time, &expense.EID); err != nil {
		return err
	}
	return nil
}

func (q *AptQueries) MEditExpense(exp *models.Expense) error {
	query := `
	UPDATE expenses 
	SET paid_time = $1, to_name = $2, to_email = $3, to_phone = $4, paid_status = $5, details = $6 
	WHERE id = $7 AND paid_status = false
	RETURNING total`
	row := q.QueryRowx(query, exp.Paid.Time, exp.To.FullName, exp.To.Email, exp.To.Phone, exp.Paid.Status, exp.Details, exp.ID)
	if row.Err() != nil {
		return row.Err()
	}
	if err := row.Scan(&exp.Total); err != nil {
		return err
	}
	return nil
}
