package queries

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"io/ioutil"
)

type InitQueries struct {
	*sqlx.DB
}

func (DB *InitQueries) CheckTablesExist() {
	DB.loadSqlFile("create.sql")
}

func (DB *InitQueries) TestData() {
	DB.loadSqlFile("cities.sql")
	DB.loadSqlFile("tests.sql")

	/*DB.Exec(`INSERT INTO users (email, password, name, surname, "type")
	VALUES ('test@test.com', 'vKrtXaX67RHf5zwpEdA1HtT8FxrhRRu6krsIUaZOTzo=', 'ALPEREN', 'ÇETİN', 1);

	INSERT INTO users (email, password, name, surname, "type", created_by)
	VALUES ('test2@test.com', 'pTzBptE5nGU6zFQi7uLYXx1W2iuJiQm2OjtI1xVFpzI=', 'AD', 'SOYAD', 0, 1);`,
	)*/
}

func (DB *InitQueries) loadSqlFile(filename string) {
	body, err := ioutil.ReadFile("platform/database/sql/" + filename)
	if err == nil {
		_, err := DB.Exec(string(body))
		if err != nil {
			fmt.Println(err)
		}
	}
}
