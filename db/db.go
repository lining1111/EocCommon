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
	sqlCmd := "replace into loginInfo(id,dataVersion) values (1,?)"
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

func GetFusionPara(fusionPara *common.FusionParaSetting) error {
	sqlCmd := "select * from  fusionPara"
	row := ConfigDb.QueryRowx(sqlCmd)
	if row.Err() != nil {
		return row.Err()
	}
	err := row.StructScan(fusionPara)
	if err != nil {
		fmt.Printf(err.Error())
		return err
	}
	return nil
}

func SetIntersectionEntity(intersectionEntity common.IntersectionEntity) error {
	ConfigDb.Exec("delete from intersectionEntity")
	_, err := ConfigDb.Exec("replace into intersectionEntity("+
		"guid,"+
		"name,"+
		"type,"+
		"plateId,"+
		"xLength,"+
		"yLength,"+
		"laneNumber,"+
		"latitude,"+
		"longitude)"+
		"values(?,?,?,?,?,?,?,?,?)",
		intersectionEntity.Guid,
		intersectionEntity.Name,
		intersectionEntity.Type,
		intersectionEntity.PlatId,
		intersectionEntity.XLength,
		intersectionEntity.YLength,
		intersectionEntity.LaneNumber,
		intersectionEntity.Latitude,
		intersectionEntity.Longitude)
	return err
}

func GetIntersectionEntity(intersectionEntity *common.IntersectionEntity) error {
	sqlCmd := "select * from  intersectionEntity"
	rows, err := ConfigDb.Query(sqlCmd)
	if err != nil {
		return err
	}
	for rows.Next() {
		err1 := rows.Scan(&intersectionEntity.Guid,
			&intersectionEntity.Name,
			&intersectionEntity.Type,
			&intersectionEntity.PlatId,
			&intersectionEntity.XLength,
			&intersectionEntity.YLength,
			&intersectionEntity.LaneNumber,
			&intersectionEntity.Latitude,
			&intersectionEntity.Longitude)
		if err1 != nil {
			fmt.Printf(err1.Error())
			return err1
		}
	}
	return nil
}

func SetIntersectionBaseSetting(intersectionBaseSetting common.IntersectionBaseSettingEntity) error {
	ConfigDb.Exec("delete from intersectionBaseSetting")
	_, err := ConfigDb.Exec("replace into intersectionBaseSetting("+
		"flagEast,"+
		"flagSouth,"+
		"flagWest,"+
		"flagNorth,"+
		"deltaXEast,"+
		"deltaYEast,"+
		"deltaXSouth,"+
		"deltaYSouth,"+
		"deltaXWest,"+
		"deltaYWest,"+
		"deltaXNorth,"+
		"deltaYNorth,"+
		"widthX,"+
		"widthY)"+
		"values(?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
		intersectionBaseSetting.FlagEast,
		intersectionBaseSetting.FlagSouth,
		intersectionBaseSetting.FlagWest,
		intersectionBaseSetting.FlagNorth,
		intersectionBaseSetting.DeltaXEast,
		intersectionBaseSetting.DeltaYEast,
		intersectionBaseSetting.DeltaXSouth,
		intersectionBaseSetting.DeltaYSouth,
		intersectionBaseSetting.DeltaXWest,
		intersectionBaseSetting.DeltaYWest,
		intersectionBaseSetting.DeltaXNorth,
		intersectionBaseSetting.DeltaYNorth,
		intersectionBaseSetting.WidthX,
		intersectionBaseSetting.WidthY)
	return err
}

func GetIntersectionBaseSetting(intersectionBaseSetting *common.IntersectionBaseSettingEntity) error {
	sqlCmd := "select * from  intersectionBaseSetting"
	row := ConfigDb.QueryRowx(sqlCmd)
	if row.Err() != nil {
		return row.Err()
	}
	err := row.StructScan(intersectionBaseSetting)
	if err != nil {
		fmt.Printf(err.Error())
		return err
	}
	return nil
}

func SetBaseSettingEntity(baseSettingEntity common.BaseSettingEntity) error {
	ConfigDb.Exec("delete from baseSettingEntity")
	_, err := ConfigDb.Exec("replace into baseSettingEntity("+
		"city,"+
		"isUploadToPlatform,"+
		"is4Gmodel,"+
		"isIllegalCapture,"+
		"isPrintIntersectionName,"+
		"remarks,"+
		"filesServicePath,"+
		"filesServicePort,"+
		"mainDNS,"+
		"alternateDNS,"+
		"platformTcpPath,"+
		"platformTcpPort,"+
		"platformHttpPath,"+
		"platformHttpPort,"+
		"signalMachinePath,"+
		"signalMachinePort,"+
		"isUseSignalMachine,"+
		"ntpServerPath,"+
		"illegalPlatformAddress,"+
		"mainboardIp,"+
		"mainboardPort)"+
		"values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
		baseSettingEntity.City,
		baseSettingEntity.IsUploadToPlatform,
		baseSettingEntity.Is4Gmodel,
		baseSettingEntity.IsIllegalCapture,
		baseSettingEntity.IsPrintIntersectionName,
		baseSettingEntity.Remarks,
		baseSettingEntity.FilesServicePath,
		baseSettingEntity.FilesServicePort,
		baseSettingEntity.MainDNS,
		baseSettingEntity.AlternateDNS,
		baseSettingEntity.PlatformTcpPath,
		baseSettingEntity.PlatformTcpPort,
		baseSettingEntity.PlatformHttpPath,
		baseSettingEntity.PlatformHttpPort,
		baseSettingEntity.SignalMachinePath,
		baseSettingEntity.SignalMachinePort,
		baseSettingEntity.IsUseSignalMachine,
		baseSettingEntity.NtpServerPath,
		baseSettingEntity.IllegalPlatformAddress,
		baseSettingEntity.FusionMainboardIp,
		baseSettingEntity.FusionMainboardPort)
	return err
}
func GetBaseSettingEntity(baseSettingEntity *common.BaseSettingEntity) error {
	sqlCmd := "select * from  baseSettingEntity"
	row := ConfigDb.QueryRowx(sqlCmd)
	if row.Err() != nil {
		return row.Err()
	}
	err := row.StructScan(baseSettingEntity)
	if err != nil {
		fmt.Printf(err.Error())
		return err
	}
	return nil
}

func SetAssociatedEquips(associatedEquips []common.AssociatedEquip) error {
	ConfigDb.Exec("delete from associatedEquip")
	var errRet error
	for _, v := range associatedEquips {
		_, errRet = ConfigDb.Exec("replace into associatedEquip("+
			"equipType,"+
			"equipCode)"+
			"values(?,?)",
			v.EquipType,
			v.EquipCode)
	}

	return errRet
}
func GetAssociatedEquips(associatedEquips *[]common.AssociatedEquip) error {
	sqlCmd := "select * from  associatedEquip"
	rows, err := ConfigDb.Query(sqlCmd)
	if err != nil {
		return err
	}
	for rows.Next() {
		var result common.AssociatedEquip
		err1 := rows.Scan(&result.EquipType,
			&result.EquipCode)
		if err1 != nil {
			fmt.Printf(err1.Error())
			return err1
		} else {
			*associatedEquips = append(*associatedEquips, result)
		}
	}
	return nil
}

func TestIntersectionEntity() {

	OpenConfigDb("./eocConfigDB/eocConfig.db")
	defer CloseConfigDb()
	in := common.IntersectionEntity{
		Name:   "Test",
		Type:   2,
		PlatId: "123",
	}
	SetIntersectionEntity(in)
	out := common.IntersectionEntity{}
	GetIntersectionEntity(&out)

}

func TestIntersectionBaseSetting() {

	OpenConfigDb("./eocConfigDB/eocConfig.db")
	defer CloseConfigDb()
	in := common.IntersectionBaseSettingEntity{
		FlagEast: 1,
		WidthX:   2.0,
	}
	SetIntersectionBaseSetting(in)
	out := common.IntersectionBaseSettingEntity{}
	GetIntersectionBaseSetting(&out)

}

func TestBaseSettingEntity() {

	OpenConfigDb("./eocConfigDB/eocConfig.db")
	defer CloseConfigDb()
	in := common.BaseSettingEntity{
		City: "beijing",
	}
	SetBaseSettingEntity(in)
	out := common.BaseSettingEntity{}
	GetBaseSettingEntity(&out)

}

func TestAssociatedEquip() {

	OpenConfigDb("./eocConfigDB/eocConfig.db")
	defer CloseConfigDb()
	in := []common.AssociatedEquip{
		{
			EquipCode: "123",
			EquipType: 1,
		},
		{
			EquipCode: "456",
			EquipType: 2,
		},
	}
	SetAssociatedEquips(in)
	out := make([]common.AssociatedEquip, 0)
	GetAssociatedEquips(&out)

}
