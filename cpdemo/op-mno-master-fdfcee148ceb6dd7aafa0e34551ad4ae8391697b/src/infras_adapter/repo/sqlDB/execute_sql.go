package sqlDB

import (
	"codeup.aliyun.com/6145b2b428003bdc3daa97c8/go-simba/go-simba-pkg.git/mysql8"
	"context"
	"database/sql"
	"github.com/google/wire"

	"github.com/pkg/errors"
	"golang.org/x/exp/slog"

	"strconv"
)

type ExcuteSql struct {
	db mysql8.DBEngine
}

var ExcuteSqlSet = wire.NewSet(NewExcuteSql)

func NewExcuteSql(db mysql8.DBEngine) ExcuteSql {
	return ExcuteSql{
		db: db,
	}
}

// 获取数据库连接
func (e *ExcuteSql) GetDBConnect() *sql.DB {
	return e.db.GetDB()
}

// 查询单条数据
func (e *ExcuteSql) GetRow(ctx context.Context, query string) (*map[string]interface{}, error) {
	slog.Info("DB query SQL :", query)
	if query == "" {
		return nil, errors.New("DB The Query Is Not Empty")
	}

	var dataMap map[string]interface{}

	// 建立连接
	connect := e.GetDBConnect()
	// 执行单条数据查询
	rows, err := connect.QueryContext(ctx, query)
	if err != nil {
		slog.Error("DB QueryContext ERROR,", err)
		return nil, err
	}
	cols, _ := rows.Columns()
	values := make([]sql.RawBytes, len(cols))
	scans := make([]interface{}, len(cols))
	for i := range values {
		scans[i] = &values[i]
	}
	var results []interface{}
	i := 0
	for rows.Next() {
		if err := rows.Scan(scans...); err != nil {
			return nil, err
		}
		data := make(map[string]interface{})
		for j, v := range values {
			key := cols[j]
			data[key] = string(v)
		}
		results = append(results, data)
		i++
	}
	if len(results) == 0 {
		slog.Error("DB The Query Data Does Not Exist,", query)
		return &dataMap, nil
	}
	dataMap = results[0].(map[string]interface{})
	return &dataMap, nil
}

func (e *ExcuteSql) GetRows(ctx context.Context, query string) (*[]map[string]interface{}, error) {
	slog.Info("DB query SQL :", query)
	if query == "" {
		return nil, errors.New("DB The Query Is Not Empty")
	}
	var dataMap []map[string]interface{}

	// 建立连接
	connect := e.GetDBConnect()
	// 执行单条数据查询
	rows, err := connect.QueryContext(ctx, query)
	if err != nil {
		slog.Error("DB QueryContext ERROR,", err)
		return nil, err
	}
	cols, _ := rows.Columns()
	values := make([]sql.RawBytes, len(cols))
	scans := make([]interface{}, len(cols))
	for i := range values {
		scans[i] = &values[i]
	}
	var results []interface{}
	i := 0
	for rows.Next() {
		if err := rows.Scan(scans...); err != nil {
			return nil, err
		}
		data := make(map[string]interface{})
		for j, v := range values {
			key := cols[j]
			data[key] = string(v)
		}
		results = append(results, data)
		i++
	}
	if len(results) == 0 {
		slog.Error("DB The Query Data Does Not Exist,", query)
		return &dataMap, nil
	}
	for k, _ := range results {
		dataMap = append(dataMap, results[k].(map[string]interface{}))
	}
	return &dataMap, nil
}

func (e *ExcuteSql) InsertRows(ctx context.Context, query string) error {
	slog.Info("DB query SQL :", query)
	if query == "" {
		return errors.New("DB The Query Is Not Empty")
	}
	// 建立连接
	connect := e.GetDBConnect()
	//开启事务
	tx, err := connect.Begin()
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, query)
	if err != nil {
		tx.Rollback()
		slog.Error("DB QueryContext ERROR,", err)
		return err
	}
	tx.Commit()
	return nil
}

func (e *ExcuteSql) GetRowsCount(ctx context.Context, query string) (*int32, error) {
	slog.Info("DB query SQL :", query)
	if query == "" {
		return nil, errors.New("DB The Query Is Not Empty")
	}

	var totalCount int32

	// 建立连接
	connect := e.GetDBConnect()
	// 执行单条数据查询
	rows, err := connect.QueryContext(ctx, query)
	if err != nil {
		slog.Error("DB QueryContext ERROR,", err)
		return nil, err
	}
	cols, _ := rows.Columns()
	values := make([]sql.RawBytes, len(cols))
	scans := make([]interface{}, len(cols))
	for i := range values {
		scans[i] = &values[i]
	}
	for rows.Next() {
		if err := rows.Scan(scans...); err != nil {
			return nil, err
		}
		for _, v := range values {
			num, _ := strconv.ParseInt(string(v), 10, 32)
			totalCount = int32(num)
		}
	}
	return &totalCount, nil
}
