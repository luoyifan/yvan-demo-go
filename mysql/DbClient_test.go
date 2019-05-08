package mysql

import (
	"fmt"
	"testing"
)

const (
	DB_Driver = "b2b_5:123456@tcp(10.3.5.30:3306)/b2b_third?charset=utf8"
)

func TestOpenMySql(t *testing.T) {
	dbClient, err := OpenMySql(DB_Driver, 50, 1000)
	if err != nil {
		panic(err)
	}

	rows, err := dbClient.query(`SELECT
activity_id, 
activity_name, 
activity_content, 
create_by,
create_at
FROM tb_activity where site_id =?`, "1111111112")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\r\n", rows)
}
