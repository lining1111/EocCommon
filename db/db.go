package db

import (
	"EocCommon/common"
	"fmt"
	"github.com/jmoiron/sqlx"
)

/**
数据库读写
*/

var ConfigDbPath = "./eocConfigDB/eocConfig.db"
var CLParkingDbPath = "./eocConfigDB/CLParking.db"

var ConfigDb *sqlx.DB
var ConfigDbIsOpen = false

func OpenConfigDb(path string) error {
	var err error
	ConfigDb, err = sqlx.Open("sqlite3", path)
	if err != nil {
		fmt.Printf("can not open db:%s,err:%v\n", path, err)
	}
	ConfigDbIsOpen = true
	return nil
}

func CloseConfigDb() {
	ConfigDb.Close()
	ConfigDbIsOpen = false
}

type LoginInfo struct {
	EquipNumber string `db:"equipNumber"`
	EquipIp     string `db:"equipIp"`
	EquipType   string `db:"equipType"`
	SoftVersion string `db:"softVersion"`
	DataVersion string `db:"dataVersion"`
}

func SetLoginInfo(info LoginInfo) error {
	ConfigDb.Exec("delete from loginInfo")
	_, err := ConfigDb.Exec("replace into loginInfo("+
		"equipNumber,"+
		"equipIp,"+
		"equipType,"+
		"softVersion,"+
		"dataVersion)"+
		"values(?,?,?,?,?)",
		info.EquipNumber,
		info.EquipIp,
		info.EquipType,
		info.SoftVersion,
		info.DataVersion)
	return err
}

func GetLoginInfo(info *LoginInfo) error {
	sqlCmd := "select * from  loginInfo"
	row := ConfigDb.QueryRowx(sqlCmd)
	if row.Err() != nil {
		return row.Err()
	}
	err := row.StructScan(info)
	if err != nil {
		fmt.Printf(err.Error())
		return err
	}
	return nil
}

func SetDataVersionInLoginInfo(dataVersion string) error {
	sqlCmd := "replace into loginInfo(id,dataVersion) values (1,'s')"
	_, err := ConfigDb.Exec(sqlCmd, dataVersion)
	return err
}

func GetDataVersionInLoginInfo(dataVersion *string) error {
	sqlCmd := "select dataVersion from  loginInfo"
	row := ConfigDb.QueryRowx(sqlCmd)
	if row.Err() != nil {
		return row.Err()
	}

	err := row.Scan(dataVersion)
	if err != nil {
		fmt.Printf(err.Error())
		return err
	}
	return nil
}

func SetFusionPara(fusionPara common.FusionParaSetting) error {
	ConfigDb.Exec("delete from fusionPara")
	_, err := ConfigDb.Exec("replace into fusionPara("+
		"repateX,"+
		"widthX,"+
		"widthY,"+
		"xmax,"+
		"ymax,"+
		"gatetx,"+
		"gatety,"+
		"gatex,"+
		"gatey,"+
		"time_flag,"+
		"angle_value)"+
		"values(?,?,?,?,?,?,?,?,?,?,?)",
		fusionPara.RepateX,
		fusionPara.WidthX,
		fusionPara.WidthY,
		fusionPara.Xmax,
		fusionPara.Ymax,
		fusionPara.Gatetx,
		fusionPara.Gatety,
		fusionPara.Gatex,
		fusionPara.Gatey,
		fusionPara.Time_flag,
		fusionPara.Angle_value)
	return err
}
