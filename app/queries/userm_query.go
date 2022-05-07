package queries

import (
	"database/sql"
	"errors"
	"github.com/Xacnio/aptms-backend/app/models"
	"log"
)

func (q *UserQueries) MGetUser(buildingId int, userId int) (models.User, error) {
	query := `
	SELECT 
    	u.id, u.email, u.name, u.surname, u.phone_number, u.type, u.created_by, uap.rank 
	FROM 
	     user_building_membership as uap 
	         LEFT JOIN users as u ON uap.user_id = u."id" 
	WHERE uap.building_id = $1 AND uap.user_id = $2`
	row := q.QueryRowx(query, buildingId, userId)
	if row.Err() != nil {
		return models.User{}, row.Err()
	}
	user := models.User{}
	err := row.StructScan(&user)
	if err != nil {
		log.Fatalln(err)
	}
	return user, nil
}

func (q *UserQueries) MGetUsersType1(buildingId int) ([]models.User, error) {
	result := make([]models.User, 0)
	query := `
	SELECT 
		u.id, u.email, u.name, u.surname, u.phone_number, u.type, u.created_by, uap.rank 
	FROM 
		user_building_membership as uap 
		    LEFT JOIN users as u ON uap.user_id = u."id" 
	WHERE uap.building_id = $1`
	rows, err := q.Queryx(query, buildingId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		user := new(models.User)
		err := rows.StructScan(&user)
		if err != nil {
			log.Fatalln(err)
		}
		result = append(result, *user)
	}
	return result, nil
}

func (q *UserQueries) MGetUsersType2(buildingId int) ([]models.BuildingUser, error) {
	result := make([]models.BuildingUser, 0)
	query := `
	SELECT 
		u.id, u.email, u.name, u.surname, u.phone_number, u.type, u.created_by, x.*, y.*
	FROM 
		user_building_membership as uap 
		LEFT JOIN users as u ON u.id = uap.user_id 
		LEFT JOIN LATERAL (
		    SELECT 
		    	string_agg(b.letter || '-' || f."number", ', ' ORDER BY b.letter,f."number") AS "owned" 
		    FROM 
		        flats as f 
		        	LEFT JOIN blocks as b ON f.block_id = b.id 
		    WHERE f.owner_id = u.id AND f.building_id = $1 
		    GROUP BY f.owner_id
		) x ON true
		LEFT JOIN LATERAL (
		    SELECT 
		    	string_agg(b.letter || '-' || f."number", ', ' ORDER BY b.letter,f."number") AS rented 
		    FROM 
		        flats as f 
		            LEFT JOIN blocks as b ON f.block_id = b.id 
		    WHERE f.tenant_id = u.id AND f.building_id = $1 
		    GROUP BY f.tenant_id
		) y ON true
	WHERE 
		uap.building_id = $1
	GROUP BY 
		uap.user_id,u.id,x."owned",y.rented
	ORDER BY
		uap.user_id
	`
	rows, err := q.Queryx(query, buildingId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		dbData := struct {
			User        models.User    `db:",prefix=u."`
			OwnedFlats  sql.NullString `db:"owned"`
			RentedFlats sql.NullString `db:"rented"`
		}{}
		buildingUser := new(models.BuildingUser)
		err := rows.StructScan(&dbData)
		if err != nil {
			log.Fatalln(err)
		}
		buildingUser.User = dbData.User
		if dbData.OwnedFlats.Valid {
			buildingUser.OwnedFlats = dbData.OwnedFlats.String
		} else {
			buildingUser.OwnedFlats = "Yok"
		}
		if dbData.RentedFlats.Valid {
			buildingUser.RentedFlats = dbData.RentedFlats.String
		} else {
			buildingUser.RentedFlats = "Yok"
		}
		result = append(result, *buildingUser)
	}
	return result, nil
}

func (q *UserQueries) MCreateUserA(user *models.User) error {
	query := `
	WITH e AS(
	    INSERT INTO users (email, password, name, surname, phone_number, type, created_by) VALUES ($1, $2, $3, $4, $5, $6, $7) 
	    ON CONFLICT (email) DO NOTHING 
	        RETURNING id
	) 
	SELECT * FROM e UNION SELECT id FROM users WHERE email LIKE $1`
	row := q.QueryRowx(query, user.Email, user.Password, user.Name, user.Surname, user.PhoneNumber, user.Type, user.CreatedBy)
	if row.Err() != nil {
		return row.Err()
	}
	err := row.Scan(&user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (q *UserQueries) MCreateUserAptPerm(userId uint, buildingId uint) {
	query := `INSERT INTO user_building_membership (user_id, building_id, rank) VALUES ($1, $2, $3) ON CONFLICT (user_id, building_id) DO NOTHING`
	q.Exec(query, userId, buildingId, 1)
}

func (q *UserQueries) MUserCheckFlat(buildingId, userId uint) (bool, error) {
	query := `SELECT COUNT(*) AS count FROM flats WHERE building_id = $1 AND (owner_id = $2 OR tenant_id = $2)`
	row := q.QueryRowx(query, buildingId, userId)
	if row.Err() != nil {
		return false, row.Err()
	}
	count := 0
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

func (q *UserQueries) MKickUser(buildingId, userId uint) error {
	query := `DELETE FROM user_building_membership WHERE building_id = $1 AND user_id = $2`
	result, err := q.Exec(query, buildingId, userId)
	if err != nil {
		return err
	}
	count, _ := result.RowsAffected()
	if count == 0 {
		return errors.New("invalid user")
	}
	return nil
}
