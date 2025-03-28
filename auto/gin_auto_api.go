package auto

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/jinzhu/copier"
)

// 自动api模型
type AutoApi struct {
	ControllerTemplates []*ControllerTemplate
	AutoApiCommon
	Classs []*AutoApiClass
}
type AutoApiClass struct {
	Name    string
	Context string
}

func NewAutoApi() *AutoApi {
	return &AutoApi{
		ControllerTemplates: make([]*ControllerTemplate, 0),
		Classs:              make([]*AutoApiClass, 0),
	}
}
func (a *AutoApi) mkdirAll() {
	os.MkdirAll(co.AppUrl, os.ModePerm)
	os.MkdirAll(co.AppUrl+`/`+co.ControllerPackageName, os.ModePerm)
	os.MkdirAll(co.AppUrl+`/`+co.ServicePackageName, os.ModePerm)
	os.MkdirAll(co.AppUrl+`/`+co.DaoPackageName, os.ModePerm)
}

// 父类主文件初始化检查
func (a *AutoApi) parentInit() {
	a.mkdirAll()
	//api文件查看是否存在
	_, err := open(co.ApiUrl)
	if err != nil {
		s, _ := open(co.Api_init_templateUrl)
		write(co.ApiUrl, replace.SubReplaceAll(s, map[string]string{
			"name": cases.Title(language.English).String(replace.GetValue("name")),
		}))
	}

	//主dao文件是否存在
	_, err = open(co.ParentDaoUrl)
	if err != nil {
		s, _ := open(co.ParentDao_init_templateUrl)
		write(co.ParentDaoUrl, replace.ReplaceAll(s))
	}
	//主controller文件是否存在
	_, err = open(co.ParentControllerUrl)
	if err != nil {
		s, _ := open(co.ParentController_init_templateUrl)
		write(co.ParentControllerUrl, replace.ReplaceAll(s))
	}
	//主service文件是否存在
	_, err = open(co.ParentServiceUrl)
	if err != nil {
		s, _ := open(co.ParentService_init_templateUrl)
		write(co.ParentServiceUrl, replace.ReplaceAll(s))
	}
}
func (a *AutoApi) init() {
	//controller文件查看是否存在
	_, err := open(co.ControllerUrl)
	if err != nil {
		s, _ := open(co.Controller_init_templateUrl)
		write(co.ControllerUrl, replace.ReplaceAll(s))
	}
	//dao文件是否存在
	_, err = open(co.DaoUrl)
	if err != nil {
		s, _ := open(co.Dao_init_templateUrl)
		write(co.DaoUrl, replace.ReplaceAll(s))
	}
	//service文件是否存在
	_, err = open(co.ServiceUrl)
	if err != nil {
		s, _ := open(co.Service_init_templateUrl)
		write(co.ServiceUrl, replace.ReplaceAll(s))
	}
}

// 去重
func (a *AutoApi) removeRepetition() {
	//controller类去重
	s, _ := open(co.ControllerUrl)
	ControllerTemplates := make([]*ControllerTemplate, 0)
	for _, item := range a.ControllerTemplates {
		if !strings.Contains(s, item.FuncName) {
			ControllerTemplates = append(ControllerTemplates, item)
		}
	}
	a.ControllerTemplates = ControllerTemplates

	//dao类去重
	s, _ = open(co.DaoUrl)
	classes := make([]*AutoApiClass, 0)
	for _, item := range a.Classs {
		if !strings.Contains(s, item.Name) {
			classes = append(classes, item)
		}
	}
	a.Classs = classes
}

// 插入内容
func (a *AutoApi) InsertContext() {
	a.parentInit()
	a.init()
	a.removeRepetition()
	a.insertControllerContext()

	a.insertApiContext()
	a.insertDaoContext()
	a.insertServiceContext()
}

// 插入service内容
func (a *AutoApi) insertServiceContext() {
	toString := a.toServiceString()
	s, _ := open(co.ServiceUrl)
	s = strings.TrimSpace(s) + "\n\n" + toString
	write(co.ServiceUrl, s)
	fmt.Println(co.ServiceUrl)
	//插入parentService内容
	s, _ = open(co.ParentServiceUrl)
	index := strings.Index(s, co.ParentServiceLogo)
	if index == -1 {
		panic("parentService文件异常")
	}
	i := strings.Index(s[index:], "}")
	i += index
	toString = a.toParentServiceString()
	s = s[:i] + toString + s[i:]
	write(co.ParentServiceUrl, s)
	fmt.Println(co.ParentServiceUrl)
}

// parentService内容
func (a *AutoApi) toParentServiceString() string {
	noResp, err := open(co.ParentService_no_resp_templateUrl)
	if err != nil {
		panic(err)
	}
	resp, err := open(co.ParentService_resp_templateUrl)
	if err != nil {
		panic(err)
	}
	builder := strings.Builder{}
	for _, item := range a.ControllerTemplates {
		var s string
		if item.RespDetail == "err" {
			s = noResp
		} else {
			s = resp
		}
		builder.WriteString("\t")
		builder.WriteString(replace.SubReplaceAll(s, item))
		builder.WriteString("\n")
	}
	return builder.String()
}

// service内容
func (a *AutoApi) toServiceString() string {
	noResp, err := open(co.Service_no_resp_templateUrl)
	if err != nil {
		panic(err)
	}
	resp, err := open(co.Service_resp_templateUrl)
	if err != nil {
		panic(err)
	}
	builder := strings.Builder{}
	for _, item := range a.ControllerTemplates {
		var s string
		if item.RespDetail == "err" {
			s = noResp
		} else {
			s = resp
		}
		builder.WriteString(replace.SubReplaceAll(s, item))
		builder.WriteString("\n\n")
	}
	return builder.String()
}

// 插入dao内容
func (a *AutoApi) insertDaoContext() {
	toString := a.toDaoString()
	s, _ := open(co.DaoUrl)
	s = strings.TrimSpace(s) + "\n\n" + toString
	write(co.DaoUrl, s)
	fmt.Println(co.DaoUrl)
}

// dao内容
func (a *AutoApi) toDaoString() string {
	builder := strings.Builder{}
	for _, item := range a.Classs {
		builder.WriteString(fmt.Sprintf("type %s struct%s\n", item.Name, item.Context))
	}
	return builder.String()
}

// 插入api内容
func (a *AutoApi) insertApiContext() {
	toString := a.toApiString()
	s, _ := open(co.ApiUrl)
	index := strings.LastIndex(s, "}")
	if index == -1 {
		panic("api文件异常")
	}
	s = s[:index] + toString + s[index:]
	write(co.ApiUrl, s)
	fmt.Println(co.ApiUrl)
}

// api内容
func (a *AutoApi) toApiString() string {
	builder := strings.Builder{}
	s, _ := open(co.Api_sub_templateUrl)
	for _, item := range a.ControllerTemplates {
		st := "\t" + s + "\n"
		builder.WriteString(replace.SubReplaceAll(st, item, map[string]string{"requestType": strings.ToUpper(item.RequestType)}))
	}
	return builder.String()
}

// 插入controller内容
func (a *AutoApi) insertControllerContext() {
	toString := a.toControllerString()
	s, _ := open(co.ControllerUrl)
	s = strings.TrimSpace(s) + "\n\n" + toString
	write(co.ControllerUrl, s)
	fmt.Println(co.ControllerUrl)
}

// controller内容
func (a *AutoApi) toControllerString() string {
	builder := strings.Builder{}
	for _, item := range a.ControllerTemplates {
		item.AutoApiCommon = a.AutoApiCommon
		builder.WriteString(item.toString())
		builder.WriteString("\n\n")
	}
	return builder.String()
}

type AutoApiCommon struct {
	Tags      string //标签
	Version   string //版本
	UrlPre    string //地址前缀
	UrlPrePre string //地址前缀的前缀
}
type ControllerTemplate struct {
	FuncName       string //方法名称
	ShouldBindType string //绑定类型
	RespDetail     string //响应详情
	RespDetailName string //响应详情名称
	ControllerTemplateSwagger
	AutoApiCommon
}
type ControllerTemplateSwagger struct {
	Summary     string //描述
	Description string //详情
	RequestWay  string //请求方式[body|query]
	ReqName     string //请求参数名称
	RespType    string //返回类型
	RespName    string //返回参数名称
	DaoRespName string //dao+resp包返回参数名称
	RespDesc    string //返回参数描述
	Url         string //地址
	RequestType string //请求类型[get|post]
}
type common struct {
	AddClassSub                       string
	Name                              string
	SuccLogo                          string //成功
	DaoPackageName                    string //dao包名称
	ServicePackageName                string
	ControllerPackageName             string
	InitUserLogo                      string //初始化用户string
	ParentServiceLogo                 string //type Services interface
	AppUrl                            string //app包地址
	ApiUrl                            string //api包地址
	Api_sub_templateUrl               string
	Api_init_templateUrl              string
	ControllerUrl                     string //controller地址
	Controller_templateUrl            string //controller模板地址
	Controller_init_templateUrl       string
	ServiceUrl                        string
	Service_no_resp_templateUrl       string
	Service_resp_templateUrl          string
	Service_init_templateUrl          string
	DaoUrl                            string
	Dao_init_templateUrl              string
	ParentDaoUrl                      string
	ParentDao_init_templateUrl        string
	ParentControllerUrl               string
	ParentController_init_templateUrl string
	ParentServiceUrl                  string
	ParentService_init_templateUrl    string
	ParentService_resp_templateUrl    string
	ParentService_no_resp_templateUrl string
	AppPackagePre                     string
	Form_apiUrl                       string
}

var co = common{}

func (c *ControllerTemplate) init() {
	if c.RequestType == "get" {
		c.RequestWay = "query"
		c.ShouldBindType = "ShouldBindQuery"
	} else {
		c.RequestWay = "body"
		c.ShouldBindType = "ShouldBindJSON"
	}
	if c.RespName == "" {
		c.DaoRespName = "string"
		c.RespType = "string"
		c.RespDesc = co.SuccLogo
		c.RespDetail = "err"
		c.RespDetailName = `"` + co.SuccLogo + `"`
	} else {
		c.RespType = "object"
		c.RespDesc = "响应体"
		c.RespDetail = "resp,err"
		c.RespDetailName = "resp"
		c.DaoRespName = co.DaoPackageName + "." + c.RespName
	}
}
func open(url string) (string, error) {
	bytes, err := os.ReadFile(url)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(bytes)), nil
}
func write(url string, data string) {
	err := os.WriteFile(url, []byte(data), 0666)
	if err != nil {
		fmt.Printf("写入错误：err:%v\n", err)
		panic(err)
	}
}
func (c *ControllerTemplate) toString() string {
	template, _ := open(co.Controller_templateUrl)
	c.init()
	return replace.SubReplaceAll(template, c)
}

var replace *Replace

func init() {
	replace = NewReplace()
	replace.WriteFile("gin_auto_public_info.application")
	s, err := open("gin_auto_setting_info.application")
	if err != nil {
		panic(err)
	}
	replace.WriteFileContext(replace.ReplaceAll(s))
	replace.Copy(&co)
	A = NewAutoApi()
}

var A *AutoApi

// 从api文件中抓取api内容
func GetApi() {
	s, err := open(replace.GetValue("Form_apiUrl"))
	if err != nil {
		panic(err)
	}
	ty := 0
	for _, v := range strings.Split(s, "\n") {
		v = strings.TrimSpace(v)
		if v == "type (" || v == "type(" {
			ty = 1
			continue
		} else if v == "@server (" || v == "@server(" {
			ty = 2
			continue
		} else if v == "service {" || v == "service{" {
			ty = 3
			continue
		}
		switch ty {
		case 1:
			getApiType(v)
		case 2:
			getApiServer(v)
		case 3:
			getApiService(v)
		}
	}
}
func getApiType(s string) {
	s = strings.TrimSpace(s)
	if s == "" {
		return
	}
	if s[len(s)-1] == '{' {
		A.Classs = append(A.Classs, &AutoApiClass{
			Name:    s[:len(s)-1],
			Context: "{",
		})
		if replace.GetValue("addClassSub") != "" {
			A.Classs[len(A.Classs)-1].Context += "\n\t" + replace.GetValue("addClassSub")
		}
	} else if s[len(s)-1] == '}' {
		A.Classs[len(A.Classs)-1].Context += "\n" + s
	} else if len(A.Classs) != 0 && A.Classs[len(A.Classs)-1].
		Context[len(A.Classs[len(A.Classs)-1].Context)-1] != '}' {
		A.Classs[len(A.Classs)-1].Context += "\n\t" + s
	}
}
func getApiServer(s string) {
	split := strings.Split(s, ":")
	if len(split) == 2 {
		split[0] = strings.TrimSpace(split[0])
		split[1] = strings.Trim(strings.TrimSpace(split[1]), `"`)
		switch split[0] {
		case "Tags":
			A.AutoApiCommon.Tags = split[1]
		case "UrlPre":
			A.AutoApiCommon.UrlPre = split[1]
		case "UrlPrePre":
			A.AutoApiCommon.UrlPrePre = split[1]
		case "Version":
			A.AutoApiCommon.Version = split[1]
		}
	}
}
func getApiService(s string) {
	split := strings.Split(strings.TrimSpace(s), " ")
	list := []string{}
	for _, v := range split {
		if v != "" && v != " " {
			list = append(list, v)
		}
	}
	if len(list) == 2 {
		if list[0] == "@doc" {
			template := &ControllerTemplate{}
			template.Summary = strings.Trim(list[1], `"`)
			A.ControllerTemplates = append(A.ControllerTemplates, template)
		} else if list[0] == "@handler" {
			A.ControllerTemplates[len(A.ControllerTemplates)-1].FuncName = list[1]
		}
	} else if len(list) == 5 {
		A.ControllerTemplates[len(A.ControllerTemplates)-1].RequestType = list[0]
		A.ControllerTemplates[len(A.ControllerTemplates)-1].Url = list[1]
		A.ControllerTemplates[len(A.ControllerTemplates)-1].ReqName = list[2][1 : len(list[2])-1]
		A.ControllerTemplates[len(A.ControllerTemplates)-1].RespName = list[4][1 : len(list[4])-1]
	}
}

type Replace struct {
	relationModel RelationModel //0代表插入关系时替换，1代表主动替换
	dict          map[string]string
	relations     []*ReplaceRelation
}

func (r *Replace) Copy(v interface{}) {
	val := reflect.ValueOf(v)
	// 获取 Person 类型的反射类型
	typ := val.Type()
	elem := val.Elem()
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem() // 获取指针指向的元素类型
	}
	// 遍历结构体字段
	for i := 0; i < typ.NumField(); i++ {
		// 获取字段
		field := typ.Field(i)
		field1 := elem.FieldByName(field.Name)
		if field1.CanSet() && field1.Kind() == reflect.String {
			field1.SetString(r.GetValue(field.Name))
		}
	}
}
func (r *Replace) GetValue(k string) string {
	return r.dict[r.serialization(k)]
}

// kv结构
type ReplaceKV struct {
	k string
	v string
}

// 替换条件
type ReplaceRelation struct {
	old []*ReplaceKV //现有条件
	new []*ReplaceKV //导向
}

// 字符串序列号
func (r *Replace) serialization(s string) string {
	return strings.ToLower(strings.ReplaceAll(s, "_", ""))
}

// 带上子类不改变自身dict情况下
func (r *Replace) SubReplaceAll(s string, objs ...interface{}) string {
	cv := r.CV()
	for _, obj := range objs {
		cv.WriteInterface(obj)
	}
	return cv.ReplaceAll(s)
}

// 替换全部
func (r *Replace) ReplaceAll(s string) string {
	buffer := &bytes.Buffer{}
	l := 0
	for i := 0; i < len(s)-1; i++ {
		if s[i] == '}' && s[i+1] == '}' {
			for j := i - 1; j > l; j-- {
				if s[j] == '{' && s[j-1] == '{' {
					buffer.WriteString(s[l : j-1])
					buffer.WriteString(r.dict[r.serialization(s[j+1:i])])
					l = i + 2
					break
				}
			}
		}
	}
	buffer.WriteString(s[l:])
	return buffer.String()
}
func NewReplaceKVByContext(st string) (*ReplaceKV, error) {
	split := strings.Split(st, "=")
	if split != nil && len(split) == 2 {
		return &ReplaceKV{
			k: split[0],
			v: split[1],
		}, nil
	}
	return nil, errors.New("newReplaceKV | err")
}
func NewReplaceKV(k, v string) *ReplaceKV {
	return &ReplaceKV{
		k: k,
		v: v,
	}
}
func NewReplaceKVs(obj interface{}) []*ReplaceKV {
	marshal, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	mp := make(map[string]interface{})
	if err := json.Unmarshal(marshal, &mp); err != nil {
		panic(err)
	}
	kvs := make([]*ReplaceKV, 0)
	for k, v := range mp {
		st := ""
		switch s := v.(type) {
		case int:
			st = strconv.Itoa(s)
		case string:
			st = s
		case int64:
			st = strconv.FormatInt(s, 10)
		case int32:
			st = strconv.FormatInt(int64(s), 10)
		case float64:
			st = strconv.FormatFloat(s, 'f', -1, 64)
		case float32:
			st = strconv.FormatFloat(float64(s), 'f', -1, 64)
		default:
			continue
		}
		kvs = append(kvs, &ReplaceKV{
			k: k,
			v: st,
		})
	}
	return kvs
}

type RelationModel int

const (
	Number_0_RelationModel RelationModel = 0
	Number_1_RelationModel RelationModel = 1
)

type ReplaceOpts func(r *Replace)

func NewReplace(opts ...ReplaceOpts) *Replace {
	r := &Replace{
		dict:      make(map[string]string),
		relations: make([]*ReplaceRelation, 0),
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}
func (r *Replace) CV() *Replace {
	replace := &Replace{}
	copier.Copy(replace, r)
	return replace
}

// 写入文件结构键值对
func (r *Replace) WriteFileContext(context string) {
	split := strings.Split(context, "\n")
	kvs := make([]*ReplaceKV, 0)
	for _, v := range split {
		if kv, err := NewReplaceKVByContext(v); err == nil {
			kvs = append(kvs, kv)
		}
	}
	r.WriteReplaceKV(kvs...)
}

// 写入文件地址键值对
func (r *Replace) WriteFile(url string) {
	bytes, err := os.ReadFile(url)
	if err != nil {
		panic(err)
	}
	r.WriteFileContext(string(bytes))
}
func (r *Replace) WriteInterface(obj interface{}) {
	r.WriteReplaceKV(NewReplaceKVs(obj)...)
}
func (r *Replace) WriteReplaceKV(kvs ...*ReplaceKV) {
	for _, kv := range kvs {
		r.dict[r.serialization(kv.k)] = strings.TrimSpace(kv.v)
	}
}

// 写入文件结构关系键值对
func (r *Replace) WriteFileContextRelation(context string) {
	split := strings.Split(context, "\n")
	for _, v := range split {
		vs := strings.Split(v, "->")
		if len(vs) == 2 {
			vs1 := make([]*ReplaceKV, 0)
			for _, k := range strings.Split(vs[0], ",") {
				if kv, err := NewReplaceKVByContext(k); err == nil {
					vs1 = append(vs1, kv)
				}
			}
			vs2 := make([]*ReplaceKV, 0)
			for _, k := range strings.Split(vs[1], ",") {
				if kv, err := NewReplaceKVByContext(k); err == nil {
					vs2 = append(vs2, kv)
				}
			}
			r.WriteReplaceKVRelation(vs1, vs2)
		}
	}
}

// 写入文件地址关系键值对
func (r *Replace) WriteFileRelation(url string) {
	byts, err := os.ReadFile(url)
	if err != nil {
		panic(err)
	}
	r.WriteFileContextRelation(string(byts))
}
func (r *Replace) WriteReplaceKVRelation(old []*ReplaceKV, new []*ReplaceKV) {
	for _, kv := range old {
		kv.v = strings.TrimSpace(kv.v)
		kv.k = r.serialization(kv.k)
		fmt.Println(kv.k, kv.v)
	}
	for _, kv := range new {
		kv.v = strings.TrimSpace(kv.v)
		kv.k = r.serialization(kv.k)
		fmt.Println(kv.k, kv.v)
	}
	relation := &ReplaceRelation{
		old: old,
		new: new,
	}
	if r.relationModel == Number_0_RelationModel {
		r.runRelation(relation)
	} else {
		r.relations = append(r.relations, relation)
	}
}
func (r *Replace) runRelation(relation *ReplaceRelation) {
	// 判断关系条件是否满足
	for _, kv := range relation.old {
		if v, ok := r.dict[kv.k]; (!ok) || (ok && v != kv.v) {
			return
		}
	}
	// 满足条件，替换关系
	for _, kv := range relation.new {
		r.dict[kv.k] = kv.v
	}
}
func (r *Replace) EnableRelation() {
	for _, relation := range r.relations {
		r.runRelation(relation)
	}
	r.relations = r.relations[:0]
}
