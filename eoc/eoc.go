package eoc

import (
	"EocCommon/common"
	"EocCommon/db"
	"bufio"
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"io"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

/**
eoc通信
开启3个协程 接收 心跳发送 状态上传
*/

type Eoc struct {
	Ip    string //eoc服务器ip
	Port  int    //eoc服务器端口号
	conn  *tls.Conn
	Run   bool
	State int //eoc连接状态

	ServerPemFilePath string //"./server.pem"
}

var (
	NotConnect = 1
	Connect    = 2
	Login      = 3
)

func (e *Eoc) Open() error {
	e.State = NotConnect
	server := e.Ip + ":" + strconv.Itoa(e.Port)
	//起一个openssl 连接
	roots := x509.NewCertPool()
	rootPEM, err := ioutil.ReadFile(e.ServerPemFilePath)
	if err != nil {
		fmt.Println("Unable to read cert.pem")
	}

	ok := roots.AppendCertsFromPEM(rootPEM)
	if !ok {
		panic("failed to parse root certificate")
	}

	conn, err1 := tls.Dial("tcp", server, &tls.Config{
		RootCAs: roots,
		InsecureSkipVerify: true,
	})
	if err1 != nil {
		panic("failed to connect: " + err1.Error())
	} else {
		e.conn = conn
		e.Run = true
	}
	fmt.Println("eoc tcp open")
	e.State = Connect
	return err

	return nil
}

func (e *Eoc) Close() error {
	e.conn.Close()
	e.Run = false

	return nil
}

func GetEquipNumber() (string, error) {
	db, err := sqlx.Open("sqlite3", db.CLParkingDbPath)
	if err != nil {
		fmt.Println("sqlite3 open fail err", err)
		return "0123456789", err
	}
	defer db.Close()
	sqlCmd := "select UName from  CL_ParkingArea"
	row := db.QueryRowx(sqlCmd)
	if row.Err() != nil {
		return "0123456789", row.Err()
	}
	var result string
	err1 := row.Scan(&result)
	if err1 != nil {
		fmt.Printf(err1.Error())
		return result, err1
	}
	return result, nil
}

func SetEquipNumber(equipNumber string) error {
	db, err := sqlx.Open("sqlite3", db.CLParkingDbPath)
	if err != nil {
		fmt.Println("sqlite3 open fail err", err)
		return err
	}
	defer db.Close()
	sqlCmd := "update CL_ParkingArea set UName ='%s'"
	_, err1 := db.Exec(sqlCmd, equipNumber)
	if err1 != nil {
		fmt.Printf(err1.Error())
		return err1
	}
	return nil
}

func GetEquipIp() (string, error) {

	retStr := ""
	//1.执行shell指令
	shell := "/home/nvidianx/bin/get_nx_net_info"
	cmd := exec.Command("/bin/bash", "-c", shell)
	result, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("cmd %s exec fail:%v\n", cmd.String(), err.Error())
		return "", err
	}
	/*
		"* get_double_net_info $ip_type $curip $curmask $curgateway $eth1_ip_type $eth1_curip $eth1_curmask $eth1_curgateway $curmaindns $curslavedns $curcloudip $curcloudport $curdevicenum $cur_city $cur_mac $protocol_version *"
	*/
	//2.从结果中获取配置
	//从结果中获取含有get_double_net_info的那一行 按空格分开为 ip_type curip curmask curgateway eth1_ip_type eth1_curip eth1_curmask eth1_curgateway curmaindns curslavedns curcloudip curcloudport curdevicenum cur_city cur_mac protocol_version
	rd := bufio.NewReader(bytes.NewReader(result))
	contents := make([]string, 20)
	isFind := false
	for !isFind {
		line, _, err1 := rd.ReadLine()
		if err1 == io.EOF {
			isFind = false
			break
		}
		str := string(line)
		if strings.Index(str, "get_double_net_info") >= 0 {
			copy(contents, strings.Split(str, " "))
			isFind = true
			break
		}
	}
	if isFind {
		if len(contents) >= 18 {

			retStr = contents[3]
			retStr += "["
			retStr += contents[7]
			retStr += "]"

		}
	}
	return retStr, nil
}

// SendLogin 连接成功后，发送登录命令
func (e *Eoc) SendLogin() error {

	equipNumber, _ := GetEquipNumber()
	equipIp, _ := GetEquipIp()

	db.OpenConfigDb(db.ConfigDbPath)
	defer db.CloseConfigDb()
	var dataVersion string
	err := db.GetDataVersionInLoginInfo(&dataVersion)
	if err != nil {
		dataVersion = ""
	}
	login := common.DataReqLogin{
		Code:        common.ReqLogin,
		EquipNumber: equipNumber,
		EquipIp:     equipIp,
		EquipType:   "nx",
		SoftVersion: common.Version,
		DataVersion: dataVersion,
	}
	ori, err1 := common.SetReqLogin(login)
	if err1 != nil {
		fmt.Println("err:", err1.Error())
		return err1
	} else {
		fmt.Println("原文:", string(ori))
		//原文加×
		plain := append(ori, '*')

		_, err2 := e.conn.Write(plain)
		if err2 != nil {
			fmt.Println("login send fail:", err2.Error())
			return err2
		}
	}
	return nil
}

func (e *Eoc) ProcessRsp(rsp string) error {
	frame := common.FrameRsp{
		Code: "",
	}
	err := json.Unmarshal([]byte(rsp), &frame)
	if err != nil {
		fmt.Println("json unmarshal err ", err.Error())
		return err
	}
	switch frame.Code {
	case common.RspHeartBeat:
		fmt.Println("eoc 心跳")
	case common.RspLogin:
		fmt.Println("eoc 登录回复")
		login := common.DataRspLogin{}
		err1 := json.Unmarshal([]byte(frame.Data), login)
		if err1 != nil {
			fmt.Println("json unmarshal err ", err1.Error())
		}
		fmt.Printf("State:%d,Message:%s", login.State, login.Message)
		if login.State == Connect {
			//设置状态为登录成功，可以发送心跳
			fmt.Println("登录成功，可以发送心跳")
			e.State = Login
		}

	case common.RsqConfig:
		fmt.Println("eoc 配置回复")
		fmt.Println("eoc config rsp:", frame.Data)
		config := common.DataRspConfig{}
		err1 := json.Unmarshal([]byte(frame.Data), config)
		if err1 != nil {
			fmt.Println("json unmarshal err ", err1.Error())
		} else {
			var state = 1
			var message = ""
			//将下发的配置进行本地存储
			//1.先将dataVersion存到数据库
			db.OpenConfigDb(db.ConfigDbPath)

			err2 := db.SetDataVersionInLoginInfo(config.DataVersion)
			if err2 != nil {
				message += err2.Error()
			}
			db.CloseConfigDb()
			//2.将所需配置存到对应的数据库
			//2.1 写入设备编码到数据库
			err3 := SetEquipNumber(config.AssociatedEquips[0].EquipCode)
			if err3 != nil {
				message += err3.Error()
			}
			//2.2 将融合参量写入数据库
			err4 := db.OpenConfigDb(db.ConfigDbPath)
			if err4 != nil {
				message += err4.Error()
			}
			defer db.CloseConfigDb()
			err5 := db.SetFusionPara(config.FusionParaSetting)
			if err5 != nil {
				message += err5.Error()
			}

			//	配置解析正确后，本地存储后，发送配置请求
			if message == "" {
				message = "配置成功"
			}
			req := common.DataReqConfig{
				Code:    common.ReqConfig,
				State:   state,
				Message: message,
			}
			ori, err6 := common.SetReqConfig(req)
			if err6 != nil {
				fmt.Println("err:", err6.Error())
			} else {
				fmt.Println("原文:", string(ori))
				//原文加×
				plain := append(ori, '*')

				_, err7 := e.conn.Write(plain)
				if err7 != nil {
					fmt.Println("config req send fail:", err7.Error())
				}
			}
		}

	case common.RspState:
		fmt.Println("eoc 状态回复")
		state := common.DataRspState{}
		err1 := json.Unmarshal([]byte(frame.Data), &state)
		if err1 != nil {
			fmt.Println("json unmarshal err:", err1.Error())
		} else {
			fmt.Printf("state rsp state:%d,message:%s", state.State, state.Message)
		}

	default:
		fmt.Println("eoc 未知命令", frame.Code)
	}

	return nil
}

func (e *Eoc) ThreadReceive() {
	fmt.Println("eoc ThreadReceive")
	for e.Run {
		content := make([]byte, 1024*1024*2)
		n, err := e.conn.Read(content)
		if err != nil {
			fmt.Println("eoc sock receive err:", err)
			e.Run = false
		}
		fmt.Println("eoc receive:", string(content[:n]))
		//根据×分割命令
		rsps := strings.Split(string(content[:n]), "*")
		//逐条进行解析执行
		for k, v := range rsps {
			fmt.Println("解析第", k, "条命令:", v)
			e.ProcessRsp(v)
		}
	}
}

func (e *Eoc) ThreadHeartBeat() {
	fmt.Println("eoc ThreadHeartBeat")
	for e.Run {
		if e.State == Login {
			//发送心跳
			ori, err := common.SetHeartBeat()
			if err != nil {
				fmt.Println("heartbeat err:", err.Error())
			} else {
				fmt.Println("原文:", string(ori))
				//原文加×
				plain := append(ori, '*')

				_, err2 := e.conn.Write(plain)
				if err2 != nil {
					fmt.Println("heartbeat send fail:", err2.Error())
					e.State = NotConnect
				}
			}
			time.Sleep(time.Duration(30) * time.Second) //30s sleep
		}
	}
}

func (e *Eoc) ThreadSendState() {
	fmt.Println("eoc ThreadSendState")
	for e.Run {
		if e.State == Login {
			//发送状态上传
			state := common.DataReqState{
				Code:  common.ReqState,
				State: 1,
			}

			ori, err := common.SetReqState(state)
			if err != nil {
				fmt.Println("reqState err:", err.Error())
			} else {
				fmt.Println("原文:", string(ori))
				//原文加×
				plain := append(ori, '*')

				_, err2 := e.conn.Write(plain)
				if err2 != nil {
					fmt.Println("reqState send fail:", err2.Error())
					e.State = NotConnect
				}
			}
			time.Sleep(time.Duration(60) * time.Second) //60s sleep
		}
	}
}

func (e *Eoc) StartLocalBusiness() {
		go e.ThreadReceive()
		go e.ThreadHeartBeat()
		go e.ThreadSendState()
}

func (e *Eoc) BusinessKeep() {
	for true {
		time.Sleep(time.Duration(60) * time.Second) //60s sleep
		if e.State == NotConnect {
			fmt.Println("进入eoc 重连")
			e.Close()
			time.Sleep(time.Duration(1) * time.Second) //1s sleep
			err := e.Open()
			if err != nil {
				fmt.Println("eoc open fail,err", err.Error())
			} else {
				if e.Run && e.State == Connect {
					//发生登录请求
					for e.State != Login {
						err1 := e.SendLogin()
						if err1 == nil {
							e.State = Login
						}
						time.Sleep(time.Duration(1) * time.Second) //10s sleep
					}
					//开启本地业务
					e.StartLocalBusiness()
				}
			}
		}
	}
}
