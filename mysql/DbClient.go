package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pingcap/errors"
	"github.com/siddontang/go/hack"
	"strconv"
	"time"
)

type DbClient struct {
	db *sql.DB
}

type DbReader struct {
	Columns     []string
	ColumnNames map[string]int
	Rows        [][]interface{}
	Cost        time.Duration
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

func (self *DbClient) Ping() error {
	return self.db.Ping()
}

func (self *DbClient) Execute(query string, args ...interface{}) (sql.Result, error) {
	stmt, err := self.db.Prepare(query)
	if err != nil {
		return nil, err
	}

	return stmt.Exec(args...)
}

func (self *DbClient) Query(query string, args ...interface{}) (*DbReader, error) {

	//log.Infof("query:%s, params:%v", query, args)

	t1 := time.Now()
	stmt, err := self.db.Prepare(query)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}
	cost := time.Since(t1)

	defer stmt.Close()
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	length := len(columns)
	columnNames := make(map[string]int, length)
	for i, n := range columns {
		columnNames[n] = i
	}

	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	var list [][]interface{}
	for rows.Next() {
		e := rows.Scan(scanArgs...)
		if e != nil {
			return nil, e
		}
		data := make([]interface{}, length)
		for i, v := range scanArgs {
			data[i] = *(v.(*interface{}))
		}
		list = append(list, data)
	}

	//log.Infof("rows:%d, columns:%d, cost:%s", len(list), length, cost)

	return &DbReader{
		Columns:     columns,
		Rows:        list,
		ColumnNames: columnNames,
		Cost:        cost,
	}, nil
}

func (r *DbReader) ColumnIndex(name string) (int, error) {
	if column, ok := r.ColumnNames[name]; ok {
		return column, nil
	} else {
		return 0, errors.Errorf("invalid field name %s", name)
	}
}

func (r *DbReader) RowLength() int {
	return len(r.Rows)
}

func (r *DbReader) ColumnLength() int {
	return len(r.Columns)
}

func (r *DbReader) GetValue(row int, column int) (interface{}, error) {
	if row >= len(r.Rows) || row < 0 {
		return nil, errors.Errorf("invalid row index %d", row)
	}

	if column >= len(r.Columns) || column < 0 {
		return nil, errors.Errorf("invalid column index %d", column)
	}

	return r.Rows[row][column], nil
}

func (r *DbReader) GetValueByName(row int, name string) (interface{}, error) {
	if column, err := r.ColumnIndex(name); err != nil {
		return nil, errors.Trace(err)
	} else {
		return r.GetValue(row, column)
	}
}

func (r *DbReader) IsNull(row, column int) (bool, error) {
	d, err := r.GetValue(row, column)
	if err != nil {
		return false, err
	}

	return d == nil, nil
}

func (r *DbReader) IsNullByName(row int, name string) (bool, error) {
	if column, err := r.ColumnIndex(name); err != nil {
		return false, err
	} else {
		return r.IsNull(row, column)
	}
}

func (r *DbReader) GetUint(row, column int) (uint64, error) {
	d, err := r.GetValue(row, column)
	if err != nil {
		return 0, err
	}

	switch v := d.(type) {
	case int:
		return uint64(v), nil
	case int8:
		return uint64(v), nil
	case int16:
		return uint64(v), nil
	case int32:
		return uint64(v), nil
	case int64:
		return uint64(v), nil
	case uint:
		return uint64(v), nil
	case uint8:
		return uint64(v), nil
	case uint16:
		return uint64(v), nil
	case uint32:
		return uint64(v), nil
	case uint64:
		return uint64(v), nil
	case float32:
		return uint64(v), nil
	case float64:
		return uint64(v), nil
	case string:
		return strconv.ParseUint(v, 10, 64)
	case []byte:
		return strconv.ParseUint(string(v), 10, 64)
	case nil:
		return 0, nil
	default:
		return 0, errors.Errorf("data type is %T", v)
	}
}

func (r *DbReader) GetUintByName(row int, name string) (uint64, error) {
	if column, err := r.ColumnIndex(name); err != nil {
		return 0, err
	} else {
		return r.GetUint(row, column)
	}
}

func (r *DbReader) GetInt(row, column int) (int64, error) {
	v, err := r.GetUint(row, column)
	if err != nil {
		return 0, err
	}

	return int64(v), nil
}

func (r *DbReader) GetIntByName(row int, name string) (int64, error) {
	v, err := r.GetUintByName(row, name)
	if err != nil {
		return 0, err
	}

	return int64(v), nil
}

func (r *DbReader) GetFloat(row, column int) (float64, error) {
	d, err := r.GetValue(row, column)
	if err != nil {
		return 0, err
	}

	switch v := d.(type) {
	case int:
		return float64(v), nil
	case int8:
		return float64(v), nil
	case int16:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case uint:
		return float64(v), nil
	case uint8:
		return float64(v), nil
	case uint16:
		return float64(v), nil
	case uint32:
		return float64(v), nil
	case uint64:
		return float64(v), nil
	case float32:
		return float64(v), nil
	case float64:
		return v, nil
	case string:
		return strconv.ParseFloat(v, 64)
	case []byte:
		return strconv.ParseFloat(string(v), 64)
	case nil:
		return 0, nil
	default:
		return 0, errors.Errorf("data type is %T", v)
	}
}

func (r *DbReader) GetFloatByName(row int, name string) (float64, error) {
	if column, err := r.ColumnIndex(name); err != nil {
		return 0, err
	} else {
		return r.GetFloat(row, column)
	}
}

func (r *DbReader) GetString(row, column int) (string, error) {
	d, err := r.GetValue(row, column)
	if err != nil {
		return "", err
	}

	switch v := d.(type) {
	case string:
		return v, nil
	case []byte:
		return hack.String(v), nil
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v), nil
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 64), nil
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	case nil:
		return "", nil
	default:
		return "", errors.Errorf("data type is %T", v)
	}
}

func (r *DbReader) GetStringByName(row int, name string) (string, error) {
	if column, err := r.ColumnIndex(name); err != nil {
		return "", err
	} else {
		return r.GetString(row, column)
	}
}
