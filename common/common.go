package common

import (
	"encoding/json"
	"github.com/google/uuid"
)

const Version = "1.0.0"

// FrameReq 请求消息结构体
type FrameReq struct {
	Guid    string `json:"guid"`
	Version string `json:"version"`
	Code    string `json:"code"`
	Data    string `json:"data"`
}

type FrameRsp struct {
	Guid          string `json:"guid"`
	Version       string `json:"version"`
	Code          string `json:"code"`
	State         int    `json:"state"`
	ResultMessage string `json:"resultMessage"`
	Data          string `json:"data"`
}

var (
	ReqHeartBeat = "MCCS100" //心跳命令请求
	RspHeartBeat = "MCCR100" //心跳命令回复

	ReqLogin = "MCCS101" //登录命令请求
	RspLogin = "MCCR101" //登录命令回复

	RsqConfig = "MCCR102" //接收业务数据
	ReqConfig = "MCCS102" //请求业务数据

	ReqState = "MCCS103" //状态上传请求
	RspState = "MCCR103" //状态上传回复

	ReqGetConfig = "MCCS104" //获取配置请求
	RspGetConfig = "MCCR104" //获取配置回复

	ReqNetState = "MCCS105" //外网状态上传请求
	RspNetState = "MCCR105" //外网状态上传回复
)

type DataReqHeartBeat struct {
	Code string `json:"code"`
}

type DataRspHeartBeat struct {
	Code string `json:"code"`
}

func SetHeartBeat() ([]byte, error) {
	data := DataReqHeartBeat{Code: ReqHeartBeat}
	dataStr, err := json.Marshal(data)
	if err != nil {
		return nil, err
	} else {
		req := FrameReq{
			Guid:    uuid.New().String(),
			Version: Version,
			Code:    ReqHeartBeat,
			Data:    string(dataStr)}
		reqStr, err1 := json.Marshal(req)
		if err1 != nil {
			return nil, err1
		} else {
			return reqStr, nil
		}
	}
}

func GetHeartBeat(data []byte) (DataRspHeartBeat, error) {
	rsp := DataRspHeartBeat{
		Code: ""}
	err := json.Unmarshal(data, &rsp)
	return rsp, err
}

// DataReqLogin MCCS101
type DataReqLogin struct {
	Code        string `json:"code"`
	EquipNumber string `json:"equipNumber"`
	EquipIp     string `json:"equipIp"`
	EquipType   string `json:"equipType"`
	SoftVersion string `json:"softVersion"`
	DataVersion string `json:"dataVersion"`
}

func SetReqLogin(data DataReqLogin) ([]byte, error) {
	dataStr, err := json.Marshal(data)
	if err != nil {
		return nil, err
	} else {
		req := FrameReq{
			Guid:    uuid.New().String(),
			Version: Version,
			Code:    ReqLogin,
			Data:    string(dataStr)}
		reqStr, err1 := json.Marshal(req)
		if err1 != nil {
			return nil, err1
		} else {
			return reqStr, nil
		}
	}
}

func GetReqLogin(data []byte) (DataReqLogin, error) {
	req := DataReqLogin{
		Code: ""}
	err := json.Unmarshal(data, &req)
	return req, err
}

// DataRspLogin MCCR101
type DataRspLogin struct {
	Code    string `json:"code"`
	State   int    `json:"state"` //1:成功 2:失败
	Message string `json:"message"`
}

func SetRspLogin(data DataRspLogin) ([]byte, error) {
	dataStr, err := json.Marshal(data)
	if err != nil {
		return nil, err
	} else {
		req := FrameReq{
			Guid:    uuid.New().String(),
			Version: Version,
			Code:    RspLogin,
			Data:    string(dataStr)}
		reqStr, err1 := json.Marshal(req)
		if err1 != nil {
			return nil, err1
		} else {
			return reqStr, nil
		}
	}
}

func GetRspLogin(data []byte) (DataRspLogin, error) {
	rsp := DataRspLogin{
		Code: ""}
	err := json.Unmarshal(data, &rsp)
	return rsp, err
}

type IntersectionBaseSettingEntity struct {
	FlagEast    int     `json:"flagEast" db:"flagEast"`
	FlagSouth   int     `json:"flagSouth" db:"flagSouth"`
	FlagWest    int     `json:"flagWest" db:"flagWest"`
	FlagNorth   int     `json:"flagNorth" db:"flagNorth"`
	DeltaXEast  float64 `json:"deltaXEast" db:"deltaXEast"`
	DeltaYEast  float64 `json:"deltaYEast" db:"deltaYEast"`
	DeltaXSouth float64 `json:"deltaXSouth" db:"deltaXSouth"`
	DeltaYSouth float64 `json:"deltaYSouth" db:"deltaYSouth"`
	DeltaXWest  float64 `json:"deltaXWest" db:"deltaXWest"`
	DeltaYWest  float64 `json:"deltaYWest" db:"deltaYWest"`
	DeltaXNorth float64 `json:"deltaXNorth" db:"deltaXNorth"`
	DeltaYNorth float64 `json:"deltaYNorth" db:"deltaYNorth"`
	WidthX      float64 `json:"widthX" db:"widthX"`
	WidthY      float64 `json:"widthY" db:"widthY"`
}

// IntersectionEntity 所属路口信息
type IntersectionEntity struct {
	Guid                    string                        `json:"guid" db:"guid"`
	Name                    string                        `json:"name" db:"name"`
	Type                    int                           `json:"type" db:"type"` //路口的类型：1=十字形，2=X形，3=T形，4=Y形
	PlatId                  string                        `json:"platId" db:"platId"`
	XLength                 float64                       `json:"xLength" db:"xLength"`
	YLength                 float64                       `json:"yLength" db:"yLength"`
	LaneNumber              int                           `json:"laneNumber" db:"laneNumber"`
	Latitude                string                        `json:"latitude" db:"latitude"`
	Longitude               string                        `json:"longitude" db:"longitude"`
	IntersectionBaseSetting IntersectionBaseSettingEntity `json:"intersectionBaseSetting" db:"intersectionBaseSetting"`
}

// BaseSettingEntity 核心板基础设置
type BaseSettingEntity struct {
	City                    string `json:"city" db:"city"`
	IsUploadToPlatform      int    `json:"isUploadToPlatform" db:"isUploadToPlatform"`
	Is4Gmodel               int    `json:"is4Gmodel" db:"is4Gmodel"`
	IsIllegalCapture        int    `json:"isIllegalCapture" db:"isIllegalCapture"`
	IsPrintIntersectionName int    `json:"isPrintIntersectionName" db:"isPrintIntersectionName"`
	Remarks                 string `json:"remarks" db:"remarks"`
	FilesServicePath        string `json:"filesServicePath" db:"filesServicePath"`
	FilesServicePort        int    `json:"filesServicePort" db:"filesServicePort"`
	MainDNS                 string `json:"mainDNS" db:"mainDNS"`
	AlternateDNS            string `json:"alternateDNS" db:"alternateDNS"`
	PlatformTcpPath         string `json:"platformTcpPath" db:"platformTcpPath"`
	PlatformTcpPort         int    `json:"platformTcpPort" db:"platformTcpPort"`
	PlatformHttpPath        string `json:"platformHttpPath" db:"platformHttpPath"`
	PlatformHttpPort        int    `json:"platformHttpPort" db:"platformHttpPort"`
	SignalMachinePath       int    `json:"signalMachinePath" db:"signalMachinePath"`
	IsUseSignalMachine      int    `json:"isUseSignalMachine" db:"isUseSignalMachine"`
	NtpServerPath           string `json:"ntpServerPath" db:"ntpServerPath"`
	IllegalPlatformAddress  string `json:"illegalPlatformAddress" db:"illegalPlatformAddress"`
	MainboardIp             string `json:"mainboardIp" db:"mainboardIp"`
	MainboardPort           int    `json:"mainboardPort" db:"mainboardPort"`
}

// FusionParaSetting 融合参数设置
type FusionParaSetting struct {
	RepateX     float64 `json:"repateX" db:"repateX"`
	WidthX      float64 `json:"widthX" db:"widthX"`
	WidthY      float64 `json:"widthY" db:"widthY"`
	Xmax        float64 `json:"xmax" db:"xmax"`
	Ymax        float64 `json:"ymax" db:"ymax"`
	Gatetx      float64 `json:"gatetx" db:"gatetx"`
	Gatety      float64 `json:"gatety" db:"gatety"`
	Gatex       float64 `json:"gatex" db:"gatex"`
	Gatey       float64 `json:"gatey" db:"gatey"`
	Time_flag   int     `json:"time_flag" db:"time_flag"`
	Angle_value int     `json:"angle_value" db:"angle_value"`
}

// AssociatedEquip 关联设备
type AssociatedEquip struct {
	EquipType int    `json:"equipType" db:"equipType"`
	EquipCode string `json:"equipCode" db:"equipCode"`
}

// DataRspConfig 接收业务数据 MCCR102
type DataRspConfig struct {
	Code              string             `json:"code"`
	DataVersion       string             `json:"dataVersion"`
	IntersectionInfo  IntersectionEntity `json:"intersectionInfo"`
	Index             int                `json:"index"`
	BaseSetting       BaseSettingEntity  `json:"baseSetting"`
	FusionParaSetting FusionParaSetting  `json:"fusionParaSetting"`
	AssociatedEquips  []AssociatedEquip  `json:"associatedEquips"`
}

func SetRspConfig(data DataRspConfig) ([]byte, error) {
	dataStr, err := json.Marshal(data)
	if err != nil {
		return nil, err
	} else {
		req := FrameReq{
			Guid:    uuid.New().String(),
			Version: Version,
			Code:    RsqConfig,
			Data:    string(dataStr)}
		reqStr, err1 := json.Marshal(req)
		if err1 != nil {
			return nil, err1
		} else {
			return reqStr, nil
		}
	}
}

func GetRspConfig(data []byte) (DataRspConfig, error) {
	rsp := DataRspConfig{
		Code: ""}
	err := json.Unmarshal(data, &rsp)
	return rsp, err
}

// DataReqConfig MCCS102
type DataReqConfig struct {
	Code    string `json:"code"`
	State   int    `json:"state"`
	Message string `json:"message"`
}

func SetReqConfig(data DataReqConfig) ([]byte, error) {
	dataStr, err := json.Marshal(data)
	if err != nil {
		return nil, err
	} else {
		req := FrameReq{
			Guid:    uuid.New().String(),
			Version: Version,
			Code:    ReqConfig,
			Data:    string(dataStr)}
		reqStr, err1 := json.Marshal(req)
		if err1 != nil {
			return nil, err1
		} else {
			return reqStr, nil
		}
	}
}

func GetReqConfig(data []byte) (DataReqConfig, error) {
	rsp := DataReqConfig{
		Code: ""}
	err := json.Unmarshal(data, &rsp)
	return rsp, err
}

// DataReqState MCCS103
type DataReqState struct {
	Code  string `json:"code"`
	State int    `json:"state"`
}

func SetReqState(data DataReqState) ([]byte, error) {
	dataStr, err := json.Marshal(data)
	if err != nil {
		return nil, err
	} else {
		req := FrameReq{
			Guid:    uuid.New().String(),
			Version: Version,
			Code:    ReqState,
			Data:    string(dataStr)}
		reqStr, err1 := json.Marshal(req)
		if err1 != nil {
			return nil, err1
		} else {
			return reqStr, nil
		}
	}
}

func GetReqState(data []byte) (DataReqState, error) {
	rsp := DataReqState{
		Code: ""}
	err := json.Unmarshal(data, &rsp)
	return rsp, err
}

// DataRspState MCCR103
type DataRspState struct {
	Code    string `json:"code"`
	State   int    `json:"state"`
	Message string `json:"message"`
}

func SetRspState(data DataRspState) ([]byte, error) {
	dataStr, err := json.Marshal(data)
	if err != nil {
		return nil, err
	} else {
		req := FrameReq{
			Guid:    uuid.New().String(),
			Version: Version,
			Code:    RspState,
			Data:    string(dataStr)}
		reqStr, err1 := json.Marshal(req)
		if err1 != nil {
			return nil, err1
		} else {
			return reqStr, nil
		}
	}
}

func GetRspState(data []byte) (DataRspState, error) {
	rsp := DataRspState{
		Code: ""}
	err := json.Unmarshal(data, &rsp)
	return rsp, err
}

// DataRspGetConfig 发送主动获取配置后的回复
type DataRspGetConfig struct {
	Code              string             `json:"code"`
	DataVersion       string             `json:"dataVersion"`
	IntersectionInfo  IntersectionEntity `json:"intersectionInfo"`
	Index             int                `json:"index"`
	BaseSetting       BaseSettingEntity  `json:"baseSetting"`
	FusionParaSetting FusionParaSetting  `json:"fusionParaSetting"`
	AssociatedEquips  []AssociatedEquip  `json:"associatedEquips"`
}

func SetRspGetConfig(data DataRspGetConfig) ([]byte, error) {
	dataStr, err := json.Marshal(data)
	if err != nil {
		return nil, err
	} else {
		req := FrameReq{
			Guid:    uuid.New().String(),
			Version: Version,
			Code:    RspGetConfig,
			Data:    string(dataStr)}
		reqStr, err1 := json.Marshal(req)
		if err1 != nil {
			return nil, err1
		} else {
			return reqStr, nil
		}
	}
}

func GetRspGetConfig(data []byte) (DataRspGetConfig, error) {
	rsp := DataRspGetConfig{
		Code: ""}
	err := json.Unmarshal(data, &rsp)
	return rsp, err
}

// DataReqGetConfig 主动发送获取配置
type DataReqGetConfig struct {
	Code          string `json:"code"`          //MCCS104
	MainBoardGuid string `json:"mainBoardGuid"` //核心板guid
}

func SetReqGetConfig(data DataReqGetConfig) ([]byte, error) {
	dataStr, err := json.Marshal(data)
	if err != nil {
		return nil, err
	} else {
		req := FrameReq{
			Guid:    uuid.New().String(),
			Version: Version,
			Code:    ReqGetConfig,
			Data:    string(dataStr)}
		reqStr, err1 := json.Marshal(req)
		if err1 != nil {
			return nil, err1
		} else {
			return reqStr, nil
		}
	}
}

func GetReqGetConfig(data []byte) (DataReqGetConfig, error) {
	rsp := DataReqGetConfig{
		Code: ""}
	err := json.Unmarshal(data, &rsp)
	return rsp, err
}

type DataReqNetState struct {
	Code          string `json:"code"` //MCCS105
	Total         int    `json:"total"`
	Success       int    `json:"success"`
	MainBoardGuid string `json:"mainBoardGuid"`
}

func SetReqNetState(data DataReqNetState) ([]byte, error) {
	dataStr, err := json.Marshal(data)
	if err != nil {
		return nil, err
	} else {
		req := FrameReq{
			Guid:    uuid.New().String(),
			Version: Version,
			Code:    ReqNetState,
			Data:    string(dataStr)}
		reqStr, err1 := json.Marshal(req)
		if err1 != nil {
			return nil, err1
		} else {
			return reqStr, nil
		}
	}
}

func GetReqNetState(data []byte) (DataReqNetState, error) {
	rsp := DataReqNetState{
		Code: ""}
	err := json.Unmarshal(data, &rsp)
	return rsp, err
}

type DataRspNetState struct {
	Code    string `json:"code"` //MCCR105
	State   int    `json:"state"`
	Message string `json:"message"`
}

func SetRspNetState(data DataRspNetState) ([]byte, error) {
	dataStr, err := json.Marshal(data)
	if err != nil {
		return nil, err
	} else {
		req := FrameReq{
			Guid:    uuid.New().String(),
			Version: Version,
			Code:    RspNetState,
			Data:    string(dataStr)}
		reqStr, err1 := json.Marshal(req)
		if err1 != nil {
			return nil, err1
		} else {
			return reqStr, nil
		}
	}
}

func GetRspNetState(data []byte) (DataRspNetState, error) {
	rsp := DataRspNetState{
		Code: ""}
	err := json.Unmarshal(data, &rsp)
	return rsp, err
}
