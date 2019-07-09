package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
)

//指针类型  和配置文件相关 不要改类型
var GlobalObj *globalClass

var workPath string

//应用参数 如分隔符 换行符

//参数和值的连接符号 如id=5还是id:5
const linkSymbol = "="

//换行符
const newLineSymbol = "\n"

//注释符
const annotationSymbol = "#"

type globalClass struct {
	IP string

	Port int

	NetProtocol string //网络协议 暂时tcp tcp4

	Name string //服务名

	MaxConnection uint32 //最大连接数 暂时未使用到

	MissionPoolCount uint32 //消息队列个数

	MissionItemCount uint32 //每个队列的最大容量
}

func InitLoad() {

	GlobalObj = &globalClass{
		IP:               "0.0.0.0",
		Port:             9999,
		NetProtocol:      "tcp",
		Name:             "my zframe",
		MaxConnection:    2,
		MissionPoolCount: 10,
		MissionItemCount: 4096,
	}

	GlobalObj.loadConfig()

	//fmt.Println(*GlobalObj)
	//fmt.Printf("%+v\n", *GlobalObj)

}
func InitPath(wPath string) {
	workPath = wPath
	InitLoad()
}

func (g *globalClass) loadConfig1() {

	path := "./conf/conf.json"

	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("配置文件错误或不存在，将使用默认配置,err:" + err.Error())
		return
	}

	err = json.Unmarshal(data, &GlobalObj)
	if err != nil {
		fmt.Println("配置文件错误或不存在，将使用默认配置,err:" + err.Error())
		return
		//panic("配置文件错误或不存在，将使用默认配置,err:" + err.Error())
	}

}

func (g *globalClass) loadConfig() {

	path := "conf/conf.cfg"

	if workPath != "" {
		path = workPath + path
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("配置文件错误或不存在，将使用默认配置,err:" + err.Error())
		return
	}

	str := string(data)

	lines := strings.Split(str, newLineSymbol) //所有的行

	for i := 0; i < len(lines); i++ {

		//读到每行数据

		line := lines[i]

		//跳过注释行
		line = strings.Trim(line, " ")
		if strings.HasPrefix(line, annotationSymbol) {
			continue
		}
		//跳过空行(没有:的行)
		if !strings.Contains(line, linkSymbol) {
			continue
		}
		//分割:
		res := strings.Split(line, linkSymbol)

		param := strings.Trim(res[0], " ") //参数
		val := strings.Trim(res[1], " ")   //值

		//fmt.Println(param, "***", val)

		handle(GlobalObj, param, val)

	}

}

func handle(args interface{}, param, val string) {

	v := reflect.ValueOf(args)
	v = v.Elem()

	n := v.NumField() //字段个数

	for i := 0; i < n; i++ {
		sf := v.Type().Field(i) //每个字段

		//字段名书写正确 忽略大小写
		//if sf.Name == param {
		if strings.ToLower(sf.Name) == strings.ToLower(param) {
			//fmt.Println(sf.Name, "字段存在", val, sf.Type.Name())
			//fmt.Println(sf.Type.Name())
			switch sf.Type.Name() {

			case "int":
				//res, err := strconv.Atoi(val)
				res, err := strconv.ParseInt(val, 10, 64)
				if err != nil {
					PanicError(param, val)
				}
				v.FieldByName(sf.Name).SetInt(res)
			case "uint32", "uint64":
				res, err := strconv.ParseUint(val, 10, 64)
				if err != nil {
					PanicError(param, val)
				}
				v.FieldByName(sf.Name).SetUint(res)
			case "string":
				v.FieldByName(sf.Name).SetString(val)

			case "bool":
				res, err := strconv.ParseBool(val)
				if err != nil {
					PanicError(param, val)
				}
				v.FieldByName(sf.Name).SetBool(res)
			case "float64", "float32":
				res, err := strconv.ParseFloat(val, 64)
				if err != nil {
					PanicError(param, val)
				}
				v.FieldByName(sf.Name).SetFloat(res)
			}
		}

	}

}

func PanicError(name string, value string) {
	panic("config file '" + name + "' parameter with wrong number : " + value)
}
