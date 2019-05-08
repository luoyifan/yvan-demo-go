package mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type DbClient struct {
	db *sql.DB
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func OpenMySql(url string, maxOpenConn int, maxIdleConn int) (*DbClient, error) {

	db, err := sql.Open("mysql", url)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConn)
	db.SetMaxIdleConns(maxIdleConn)
	_ = db.Ping()

	dbClient := DbClient{
		db: db,
	}
	return &dbClient, nil
}

func (self DbClient) execute(query string, args ...interface{}) (sql.Result, error) {
	stmt, err := self.db.Prepare(query)
	if err != nil {
		return nil, err
	}

	return stmt.Exec(args...)
}

func (self DbClient) query(query string, args ...interface{}) (*[]map[string]interface{}, error) {
	stmt, err := self.db.Prepare(query)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	list := make([]map[string]interface{}, len(columns))
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}
		record := make(map[string]interface{})
		for i, col := range values {
			if col != nil {
				record[columns[i]] = col
			}
		}
		list = append(list, record)
	}
	return &list, nil
}
