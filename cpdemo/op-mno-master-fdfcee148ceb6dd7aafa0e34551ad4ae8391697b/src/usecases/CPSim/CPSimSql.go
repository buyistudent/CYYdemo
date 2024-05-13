package CPSim

import (
	"codeup.aliyun.com/6145b2b428003bdc3daa97c8/go-simba/go-simba-proto.git/gen"
	"codeup.aliyun.com/6145b2b428003bdc3daa97c8/operation-platform/op-mno.git/src/app/config"
	"codeup.aliyun.com/6145b2b428003bdc3daa97c8/operation-platform/op-mno.git/src/domain/CPSimDomain"
	"context"
	"database/sql"
	"fmt"
)

type CPSimService struct {
	cfg *config.Config
}

func NewSimService(cfg *config.Config) CPSimCase {
	return &CPSimService{
		cfg: cfg,
	}
}

func (s *CPSimService) GetCPSimStatus(ctx context.Context, req CPSimDomain.CPSimUsageDTO) (*gen.GetCustomerSimStatusResp, error) {
	connectionString := fmt.Sprintf("root:Simba123!@#@tcp(192.168.24.225:3306)/cp-rus-pre")
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		fmt.Println("connect sql error:", err)
	}
	defer db.Close()
	query := "SELECT iccid,id,status FROM t_cp_sim Where id=" + req.SimId
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("query sql error:", err)
	}
	defer rows.Close()
	resp := gen.SimStatus{}
	for rows.Next() {
		var id string
		var iccid string
		var status string
		err := rows.Scan(&id, &iccid, &status)
		if err != nil {
			return nil, err
		}
		switch status {
		case "1":
			status = "未生效"
		case "2":
			status = "可测试"
		case "3":
			status = "库存"
		case "4":
			status = "激活"
		case "5":
			status = "停卡"
		}
		resp.Status = status
		//成功
		return &gen.GetCustomerSimStatusResp{Code: "200", Message: "查询成功", Result: &resp}, err
	}
	//失败
	return &gen.GetCustomerSimStatusResp{Code: "500", Message: "查询失败"}, nil
}
