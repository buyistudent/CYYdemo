package repo

import (
	"codeup.aliyun.com/6145b2b428003bdc3daa97c8/go-simba/go-simba-pkg.git/mysql8"
	"codeup.aliyun.com/6145b2b428003bdc3daa97c8/operation-platform/op-mno.git/src/domain/mnoSimDomain"
	"codeup.aliyun.com/6145b2b428003bdc3daa97c8/operation-platform/op-mno.git/src/infras_adapter/db"
	"context"
	"database/sql"
	"github.com/pkg/errors"

	"github.com/google/wire"
)

var _ mnoSimDomain.SimRepo = (*simRepo)(nil)

var SimRepositorySet = wire.NewSet(NewSimRepo)

type simRepo struct {
	db mysql8.DBEngine
}

func NewSimRepo(db mysql8.DBEngine) mnoSimDomain.SimRepo {
	return &simRepo{
		db: db,
	}
}

// 插入ESim 记录流水
func (o *simRepo) SaveBillESim(ctx context.Context, bill *mnoSimDomain.ESimBill) (*mnoSimDomain.ESimBill, error) {
	context.WithValue(context.Background(), "db", o.db)
	conn := o.db.GetDB()
	queries := db.New(conn)
	//开启事务
	tx, e := conn.Begin()
	if e != nil {
		return nil, errors.Wrap(e, "repo SaveBillESim ERROR")
	}
	var err error
	params := db.SaveEsimParams{
		ID:     bill.ID,
		Iccid:  bill.Iccid,
		Imsi:   bill.Imsi,
		Msisdn: bill.Msisdn,
		CreateTime: sql.NullTime{
			Time:  bill.CreateTime,
			Valid: true,
		},
		Json: sql.NullString{
			String: bill.Json,
			Valid:  true,
		},
		Status: sql.NullString{
			String: bill.Status,
			Valid:  true,
		},
	}
	_, err = queries.WithTx(tx).SaveEsim(ctx, params)
	if err != nil {
		tx.Rollback()
		return nil, errors.Wrap(err, "repo SaveEsim ERROR")
	}
	tx.Commit()
	return bill, nil

}

func (s *simRepo) UpdateESimBill(ctx context.Context, bill *mnoSimDomain.ESimBill) error {

	conn := s.db.GetDB()
	queries := db.New(conn)
	params := db.UpdateEsimEffectiveTimeParams{
		ID: bill.ID,
		Status: sql.NullString{
			String: bill.Status,
			Valid:  true,
		},
		EffectiveTime: sql.NullTime{
			Time:  bill.EffectiveTime,
			Valid: true,
		},
	}
	err := queries.UpdateEsimEffectiveTime(ctx, params)
	if err != nil {
		return errors.Wrap(err, "repo UpdateESimBill ERROR")
	}
	return nil
}
