package queries

import (
	"fmt"
	"github.com/Xacnio/aptms-backend/app/models"
	"github.com/jmoiron/sqlx"
	"log"
	"strconv"
	"strings"
)

type UserQueries struct {
	*sqlx.DB
}

func (q *UserQueries) GetUserByEmailPass(email, password string) (models.User, error) {
	user := models.User{}
	query := `SELECT * FROM users WHERE email LIKE $1 AND password LIKE $2`
	err := q.Get(&user, query, email, password)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (q *UserQueries) GetUserById(userid uint) (models.User, error) {
	user := models.User{}
	query := `SELECT * FROM users WHERE id = $1`
	err := q.Get(&user, query, userid)
	if err != nil {
		fmt.Println(err)
		return user, err
	}
	return user, nil
}

func (q *UserQueries) GetUsersA(page, count uint8) ([]models.ManagerUser, error) {
	result := make([]models.ManagerUser, 0)
	query := `
	SELECT 
	    u1.*, CASE WHEN u2.id > 0 THEN CONCAT(u2.name, ' ', u2.surname) ELSE '-' END as created_by_name, 
	    COALESCE(string_agg(a.id || '-' || a.name, ',' ORDER BY a.id),'') as buildings 
	FROM 
	    users as u1 
	    	LEFT JOIN users as u2 ON u1.created_by = u2.id 
	        LEFT JOIN user_building_membership as uap ON uap.user_id = u1.id AND uap.rank = 1 
	        LEFT JOIN buildings as a ON a.id = uap.building_id 
	GROUP BY u1.id,u2.id 
	LIMIT $2 OFFSET $1`
	rows, err := q.Queryx(query, page, count)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		test := struct {
			User      models.User `db:",prefix=u1."`
			Buildings string      `db:"buildings"`
		}{}
		err := rows.StructScan(&test)
		if err != nil {
			log.Fatalln(err)
		}
		buildings := make([]models.ManagerUserBuilding, 0)
		for _, value := range strings.Split(test.Buildings, ",") {
			values := strings.SplitN(value, "-", 2)
			if len(values) == 2 {
				id, _ := strconv.ParseInt(values[0], 10, 32)
				buildings = append(buildings, models.ManagerUserBuilding{ID: int(id), Name: values[1]})
			}
		}
		result = append(result, models.ManagerUser{User: test.User, Buildings: buildings})
	}
	return result, nil
}

func (q *UserQueries) CreateUserA(user *models.User) error {
	query := `
	WITH e AS(
	    INSERT INTO users (email, name, surname, phone_number, type, created_by) 
	        VALUES ($1, $2, $3, $4, $5, $6) 
	        ON CONFLICT (email) DO NOTHING RETURNING id
	) 
	SELECT * FROM e UNION SELECT id FROM users WHERE email LIKE $1`
	row := q.QueryRowx(query, user.Email, user.Name, user.Surname, user.PhoneNumber, user.Type, user.CreatedBy)
	if row.Err() != nil {
		return row.Err()
	}
	err := row.Scan(&user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (q *UserQueries) CreateUserBuildingPerm(userId uint, buildingIds []uint, rank int) {
	for _, buildingId := range buildingIds {
		query := `INSERT INTO user_building_membership (user_id, building_id, rank) VALUES ($1, $2, $3) ON CONFLICT (user_id, building_id) DO NOTHING`
		q.Exec(query, userId, buildingId, rank)
	}
}

func (q *UserQueries) UpdateUserAptPerm(userId uint, buildingIds []uint, rank int) {
	q.Exec("UPDATE user_building_membership SET rank = 0 WHERE user_id = $1", userId)
	for _, buildingId := range buildingIds {
		query := `INSERT INTO user_building_membership (user_id, building_id, rank) VALUES ($1, $2, $3) ON CONFLICT (user_id, building_id) DO UPDATE SET rank = $3`
		q.Exec(query, userId, buildingId, rank)
	}
}

func (q *UserQueries) UpdateUserA(user models.User) error {
	query := `UPDATE users SET email = $1, name = $2, surname = $3, phone_number = $4 WHERE id = $5`
	_, err := q.Exec(query, user.Email, user.Name, user.Surname, user.PhoneNumber, user.ID)
	if err != nil {
		return err
	}
	return nil
}
