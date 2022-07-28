package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type ServerNet struct {
	IP   string
	Port int
}

var ServerNetDbPath = "./eocConfigDB/RoadsideParking.db"

var ServerNetDb *sqlx.DB
var ServerNetDbIsOpen = false

func OpenServerNetDB(path string) error {
	var err error
	ServerNetDb, err = sqlx.Open("sqlite3", path)
	if err != nil {
		fmt.Printf("can not open db:%s,err:%v\n", path, err)
	}
	ServerNetDbIsOpen = true
	return nil
}

func CloseServerNetDB() {
	ServerNetDb.Close()
	ServerNetDbIsOpen = false
}

func GetServerNet(serverNet *ServerNet) error {
	sqlCmd := "select CloudServerPath from TB_ParkingLot where ID=1"
	row := ServerNetDb.QueryRowx(sqlCmd)
	if row.Err() != nil {
		return row.Err()
	}
	err := row.Scan(&serverNet.IP)
	if err != nil {
		fmt.Printf(err.Error())
		return err
	}

	sqlCmd1 := "select CloudServerPort from TB_ParkingLot where ID=1"
	row1 := ServerNetDb.QueryRowx(sqlCmd1)
	if row1.Err() != nil {
		return row1.Err()
	}
	err1 := row1.Scan(&serverNet.Port)
	if err1 != nil {
		fmt.Printf(err1.Error())
		return err1
	}

	return nil
}
