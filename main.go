package main

import (
	"EocCommon/common"
	"EocCommon/db"
	"EocCommon/eoc"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"sync"
	"time"
)

/**
与eoc通信的程序
通信方式 TCP 端口6526
每条命令发送的结尾添加×作为结束分割符
心跳间隔 30s
断线重连间隔60s
通信内容 加密方式Tls

详细的通信内容见common文件

*/

var e eoc.Eoc

func main() {

	var showVersion bool
	var ServerIp string //10.100.32.128
	var ServerPort int  //6526

	var ServerPemFilePath string

	var dbNetPath string
	var dbConfigPath string
	var dbCLParkingPath string

	//读取传入的参数
	flag.BoolVar(&showVersion, "v", false, "显示版本号")
	flag.StringVar(&ServerIp, "ip", "", "默认为空，从数据库读取ip，也可以自己设置")
	flag.IntVar(&ServerPort, "port", 0, "默认为空,从数据库读取port，也可以自己设置")
	flag.StringVar(&ServerPemFilePath, "pem", "./eoc.pem", "默认为./eoc.pem，也可以自己设置")

	flag.StringVar(&dbConfigPath, "dbConfigPath", "/home/nvidianx/bin/eocConfig.db", "默认为/home/nvidianx/bin/eocConfig.db，也可以自己设置")
	flag.StringVar(&dbNetPath, "dbNetPath", "/home/nvidianx/bin/RoadsideParking.db", "默认为/home/nvidianx/bin/RoadsideParking.db，也可以自己设置")
	flag.StringVar(&dbCLParkingPath, "dbCLParkingPath", "/home/nvidianx/bin/CLParking.db", "默认为/home/nvidianx/bin/CLParking.db，也可以自己设置")
	flag.Parse()
	if showVersion {
		fmt.Println("version:", common.Version)
		os.Exit(0)
	}

	db.ConfigDbPath = dbConfigPath
	db.ServerNetDbPath = dbNetPath
	db.CLParkingDbPath = dbCLParkingPath

	//打开数据库获取网络配置
	err := db.OpenServerNetDB(dbNetPath)
	if err != nil {
		fmt.Println(err)
	}
	serverNet := db.ServerNet{}
	if db.ServerNetDbIsOpen {

		err1 := db.GetServerNet(&serverNet)
		if err1 != nil {
			fmt.Println(err1)
		}
	}

	if ServerIp == "" {
		ServerIp = serverNet.IP
		fmt.Println("ip:", ServerIp)
	}
	if ServerPort == 0 {
		ServerPort = serverNet.Port
		fmt.Println("port:", ServerPort)
	}

	e.Ip = ServerIp
	e.Port = ServerPort
	e.ServerPemFilePath = ServerPemFilePath

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		fmt.Println("eoc thread")

		err := e.Open()
		if err != nil {
			fmt.Println("eoc open fail,err", err.Error())
		} else {
			if e.State == eoc.Connect {
				//发生登录请求
				for e.State != eoc.Login {
					err1 := e.SendLogin()
					if err1 == nil {
						e.State = eoc.Login
						fmt.Println("登录发送成功")
					}
					time.Sleep(time.Duration(1) * time.Second) //10s sleep
				}
			}

			//开启本地业务
			e.StartLocalBusiness()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("eoc keep thread")

		e.BusinessKeep()
	}()

	wg.Wait()
}
