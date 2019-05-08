package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/siddontang/go-mysql/client"
	"github.com/siddontang/go-mysql/mysql"
	"github.com/siddontang/go-mysql/replication"
)

type TableMeta struct {
	ColumnNames []string
}

func Execute() {
	conn, _ := client.Connect("47.99.62.170:3306", "ent", "test2018", "ent")
	_ = conn.Ping()
	conn.Execute(`replace into id_value(id_code, current_value) values
(?, ?)`, "go_test6", 10)
}

func GetTableNames(tableName string) []string {
	conn, _ := client.Connect("47.99.62.170:3306", "ent", "test2018", "ent")
	r, _ := conn.Execute("desc " + tableName)
	var ls []string
	for rowIndex := range r.Values {
		s, _ := r.GetString(rowIndex, 0)
		ls = append(ls, s)
	}
	return ls
}

func Start() {
	cfg := replication.BinlogSyncerConfig{
		ServerID: 198,
		Flavor:   "mysql",
		Host:     "47.99.62.170",
		Port:     3306,
		User:     "canal",
		Password: "canal",
	}

	syncer := replication.NewBinlogSyncer(cfg)
	streamer, _ := syncer.StartSync(mysql.Position{Name: "mysql-bin.000025", Pos: 11527237})
	metaCache := map[string]TableMeta{}

	for {
		ev, _ := streamer.GetEvent(context.Background())
		// Dump event
		var r bool
		var tp string
		switch ev.Header.EventType {
		case replication.WRITE_ROWS_EVENTv0:
		case replication.WRITE_ROWS_EVENTv1:
		case replication.WRITE_ROWS_EVENTv2:
			r = true
			tp = "INSERT"
			break

		case replication.UPDATE_ROWS_EVENTv0:
		case replication.UPDATE_ROWS_EVENTv1:
		case replication.UPDATE_ROWS_EVENTv2:
			r = true
			tp = "UPDATE"
			break

		case replication.DELETE_ROWS_EVENTv0:
		case replication.DELETE_ROWS_EVENTv1:
		case replication.DELETE_ROWS_EVENTv2:
			r = true
			tp = "DELETE"
			break

		default:
			r = false
		}

		if r {
			rowsEvent := ev.Event.(*replication.RowsEvent)
			tableName := string(rowsEvent.Table.Schema) + "." + string(rowsEvent.Table.Table)

			_, ok := metaCache[tableName]
			if !ok {
				//没有该表格的元数据
				metaCache[tableName] = TableMeta{
					ColumnNames: GetTableNames(tableName),
				}
			}

			//从缓存取元数据
			meta := metaCache[tableName]

			rowsLength := len(rowsEvent.Rows)
			data := make([]map[string]interface{}, rowsLength)

			//构造多个行
			for rowIndex, colsArray := range rowsEvent.Rows {

				//构造每行的列
				cols := make([]map[string]interface{}, len(colsArray))
				for colIndex, colValue := range colsArray {
					cols[colIndex] = map[string]interface{}{
						"name":   meta.ColumnNames[colIndex],
						"index":  colIndex,
						"before": "unknow",
						"after":  colValue,
					}
				}

				data[rowIndex] = map[string]interface{}{
					"database": string(rowsEvent.Table.Schema),
					"type":     tp,
					"table":    tableName,
					"cols":     cols,
				}
			}

			event := map[string]interface{}{
				"size":        rowsLength,
				"data":        data,
				"retry_count": 0,
			}
			json.Marshal(event)
			//fmt.Println(jsonString)
			fmt.Println(tp, " ", tableName, " rows:", rowsLength)
		}
	}
}
