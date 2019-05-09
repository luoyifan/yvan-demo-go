package mysql

import (
	"github.com/siddontang/go-log/log"
	"testing"
)

const (
	DB_Driver = "b2b_5:123456@tcp(10.3.5.30:3306)/b2b_third?charset=utf8"
)

var dbClient *DbClient

func init() {
	dbClient, _ = OpenMySql(DB_Driver, 50, 1000)
}

func TestOpenMySql(t *testing.T) {
	t.Log(dbClient.Ping())
}

func Benchmark_DbClient_Query(b *testing.B) {
	for i := 0; i < b.N; i++ { //use b.N for looping
		dbClient.Query(`SELECT
activity_id, 
activity_name, 
activity_content, 
create_by,
create_at
FROM tb_activity where site_id=?`, "1111111112")
	}
}

func TestDbClient_Query(t *testing.T) {
	reader, err := dbClient.Query(`SELECT
activity_id, 
activity_name, 
activity_content, 
create_by,
create_at
FROM tb_activity where site_id=?`, "1111111112")
	if err != nil {
		panic(err)
	}

	for i := 0; i < reader.RowLength(); i++ {

		for _, c := range reader.Columns {
			v, _ := reader.GetStringByName(i, c)
			log.Infof("row[%d][%s]=%s", i, c, v)
		}

		log.Info("---")
	}
}
