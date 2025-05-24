package auto

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/jinzhu/copier"
)

// 自动api模型
type AutoApi struct {
	SubList []*SubAutoApi
	Co      map[string]string
	Classes []*AutoApiClass
}
type AutoApiClass struct {
	Name    string
	Context string
}

func NewAutoApi() *AutoApi {
	return &AutoApi{
		SubList: make([]*SubAutoApi, 0),
		Classes: make([]*AutoApiClass, 0),
		Co:      make(map[string]string),
	}
}

func (a *AutoApi) mkdirAll() {
	os.MkdirAll(co.ApiDirFileAddr, os.ModePerm)
	os.MkdirAll(co.AppDirFileAddr, os.ModePerm)
	os.MkdirAll(co.ControllerDirFileAddr, os.ModePerm)
	os.MkdirAll(co.ServiceDirFileAddr, os.ModePerm)
	os.MkdirAll(co.DaoDirFileAddr, os.ModePerm)
}

// checkStatIfNotExistsCreateTemplate 检查文件是否存在不存在即创建模板
func checkStatIfNotExistsCreateTemplate(fileAddr string, templateAddr string) {
	if !stat(fileAddr) {
		s, _ := open(templateAddr)
		//把替换后的内容写入文件
		write(fileAddr, replace.ReplaceAll(s))
	}
}

// init 初始化（包含创建初始化模板）
func (a *AutoApi) init() {
	a.mkdirAll()
	//api
	checkStatIfNotExistsCreateTemplate(co.ApiModuleFileAddr, co.ApiModuleTemplateInitFileAddr)
	//controller
	checkStatIfNotExistsCreateTemplate(co.ControllerSubjectFileAddr, co.ControllerSubjectTemplateInitFileAddr)
	checkStatIfNotExistsCreateTemplate(co.ControllerSubModuleFileAddr, co.ControllerSubModuleTemplateInitFileAddr)
	//service
	checkStatIfNotExistsCreateTemplate(co.ServiceSubjectFileAddr, co.ServiceSubjectTemplateInitFileAddr)
	checkStatIfNotExistsCreateTemplate(co.ServiceSubModuleFileAddr, co.ServiceSubModuleTemplateInitFileAddr)
	//dao
	checkStatIfNotExistsCreateTemplate(co.DaoSubjectFileAddr, co.DaoSubjectTemplateInitFileAddr)
	checkStatIfNotExistsCreateTemplate(co.DaoSubModuleFileAddr, co.DaoSubModuleTemplateInitFileAddr)
}

// 去重
func (a *AutoApi) removeRepetition() {

	//controller去重
	s, _ := open(co.ControllerSubModuleFileAddr)
	subList := make([]*SubAutoApi, 0)
	for _, item := range a.SubList {
		if !strings.Contains(s, item.FuncName+"()") {
			subList = append(subList, item)
		}
	}
	a.SubList = subList

	//dao类去重
	s, _ = open(co.DaoSubModuleFileAddr)
	classes := make([]*AutoApiClass, 0)
	for _, item := range a.Classes {
		if !strings.Contains(s, item.Name) {
			classes = append(classes, item)
		}
	}
	a.Classes = classes
}

// 插入内容
func (a *AutoApi) InsertContext() {
	a.init()
	a.removeRepetition()

	a.insertControllerSubModuleContext()
	a.insertApiModuleContext()
	a.insertDaoContext()
	a.insertServiceContext()
}

// 插入Service内容
func (a *AutoApi) insertServiceContext() {
	//ServiceSubModule内容
	a.insertServiceSubModuleContext()
	//ServiceSubject内容
	a.insertServiceSubjectContext()
}

// 插入ServiceSubModule内容
func (a *AutoApi) insertServiceSubModuleContext() {
	toString := a.toServiceSubModuleString()
	s, _ := open(co.ServiceSubModuleFileAddr)
	s = strings.TrimSpace(s) + "\n\n" + toString
	write(co.ServiceSubModuleFileAddr, s)
	fmt.Println(co.ServiceSubModuleFileAddr)
}

// 插入ServiceSubject内容
func (a *AutoApi) insertServiceSubjectContext() {
	//ServiceSubject内容
	s, _ := open(co.ServiceSubjectFileAddr)
	index := strings.Index(s, co.ServiceSubjectLogo)
	if index == -1 {
		panic("ServiceSubject文件异常")
	}
	i := strings.Index(s[index:], "}")
	i += index
	toString := a.toServiceSubjectString()
	s = s[:i] + toString + s[i:]
	write(co.ServiceSubjectFileAddr, s)
	fmt.Println(co.ServiceSubjectFileAddr)
}

// ServiceSubject内容
func (a *AutoApi) toServiceSubjectString() string {
	s, err := open(co.ServiceSubjectTemplateAddFileAddr)
	if err != nil {
		panic("toParentServiceString:" + err.Error())
	}
	builder := strings.Builder{}
	for _, item := range a.SubList {
		builder.WriteString("\t")
		builder.WriteString(replace.SubReplaceAll(s, item))
		builder.WriteString("\n")
	}
	return builder.String()
}

// ServiceSubModule内容
func (a *AutoApi) toServiceSubModuleString() string {
	s, err := open(co.ServiceSubModuleTemplateAddFileAddr)
	if err != nil {
		panic("service内容:" + err.Error())
	}
	builder := strings.Builder{}
	for _, item := range a.SubList {
		builder.WriteString(replace.SubReplaceAll(s, item))
		builder.WriteString("\n\n")
	}
	return builder.String()
}

// 插入dao内容
func (a *AutoApi) insertDaoContext() {
	toString := a.toDaoString()
	s, err := open(co.DaoSubModuleFileAddr)
	if err != nil {
		panic("insertDaoContext:" + err.Error())
	}
	s = strings.TrimSpace(s) + "\n\n" + toString
	write(co.DaoSubModuleFileAddr, s)
	fmt.Println(co.DaoSubModuleFileAddr)
}

// dao内容
func (a *AutoApi) toDaoString() string {
	builder := strings.Builder{}
	for _, item := range a.Classes {
		builder.WriteString(fmt.Sprintf("type %s struct%s\n", item.Name, item.Context))
	}
	return builder.String()
}

// 插入ApiModule内容
func (a *AutoApi) insertApiModuleContext() {
	toString := a.toApiModuleString()
	s, err := open(co.ApiModuleFileAddr)
	if err != nil {
		panic("insertApiContext:" + err.Error())
	}
	index := strings.LastIndex(s, "}")
	if index == -1 {
		panic("api文件异常")
	}
	s = s[:index] + toString + s[index:]
	write(co.ApiModuleFileAddr, s)
	fmt.Println(co.ApiModuleFileAddr)
}

// ApiModule内容
func (a *AutoApi) toApiModuleString() string {
	builder := strings.Builder{}
	s, err := open(co.ApiModuleTemplateAddFileAddr)
	if err != nil {
		panic("api内容:" + err.Error())
	}
	for _, item := range a.SubList {
		st := "\t" + s + "\n"
		builder.WriteString(replace.SubReplaceAll(st, item))
	}
	return builder.String()
}

// 插入ControllerSubModule内容
func (a *AutoApi) insertControllerSubModuleContext() {
	toString := a.toControllerSubModuleString()
	s, err := open(co.ControllerSubModuleFileAddr)
	if err != nil {
		panic("insertControllerContext:" + err.Error())
	}
	s = strings.TrimSpace(s) + "\n\n" + toString
	write(co.ControllerSubModuleFileAddr, s)
	fmt.Println(co.ControllerSubModuleFileAddr)
}

// ControllerSubModule内容
func (a *AutoApi) toControllerSubModuleString() string {
	s, err := open(co.ControllerSubModuleTemplateAddFileAddr)
	if err != nil {
		panic("controller内容:" + err.Error())
	}
	builder := strings.Builder{}
	for _, item := range a.SubList {
		builder.WriteString(replace.SubReplaceAll(s, item))
		builder.WriteString("\n\n")
	}
	return builder.String()
}

type SubAutoApi struct {
	Summary     string //描述
	FuncName    string //方法名称
	RequestType string //请求类型[get|post|put]
	Url         string //地址
	ReqName     string //请求参数名称
	RespName    string //返回参数描述
}

type Common struct {
	WorkDir           string // 工作目录，app目录下的pwd
	TemplatesFileAddr string // 模板文件目录
	SettingsFileAddr  string //配置文件目录

	//project 该项目
	ProjectDirFileAddr    string
	ProjectDirPackageAddr string

	//module 模块
	ModuleName string
	//sub_module 子模块（一般在.api文件中提到）
	SubModuleName string

	//api_dir
	ApiDirName        string
	ApiDirFileAddr    string
	ApiDirPackageAddr string
	//api_module
	ApiModuleName     string
	ApiModuleFileAddr string
	//api_module_template（模板）
	ApiModuleTemplateInitFileAddr string
	ApiModuleTemplateAddFileAddr  string

	//app_dir
	AppDirName        string
	AppDirFileAddr    string
	AppDirPackageAddr string

	//controller_dir
	ControllerDirName        string
	ControllerDirFileAddr    string
	ControllerDirPackageAddr string
	//controller_subject
	ControllerSubjectName     string
	ControllerSubjectFileAddr string
	//controller_subject_template（模板）
	ControllerSubjectTemplateInitFileAddr string
	ControllerSubjectTemplateAddFileAddr  string
	//controller_sub_module
	ControllerSubModuleName     string
	ControllerSubModuleFileAddr string
	//controller_sub_module_template（模板）
	ControllerSubModuleTemplateInitFileAddr string
	ControllerSubModuleTemplateAddFileAddr  string

	//service_dir
	ServiceDirName        string
	ServiceDirFileAddr    string
	ServiceDirPackageAddr string
	//service_subject
	ServiceSubjectName     string
	ServiceSubjectFileAddr string
	ServiceSubjectLogo     string
	//service_subject_template（模板）
	ServiceSubjectTemplateInitFileAddr string
	ServiceSubjectTemplateAddFileAddr  string
	//service_sub_module
	ServiceSubModuleName     string
	ServiceSubModuleFileAddr string
	//service_sub_module_template（模板）
	ServiceSubModuleTemplateInitFileAddr string
	ServiceSubModuleTemplateAddFileAddr  string

	//dao_dir
	DaoDirName        string
	DaoDirFileAddr    string
	DaoDirPackageAddr string
	//dao_subject
	DaoSubjectName     string
	DaoSubjectFileAddr string
	//dao_subject_template（模板）
	DaoSubjectTemplateInitFileAddr string
	DaoSubjectTemplateAddFileAddr  string
	//dao_sub_module
	DaoSubModuleName     string
	DaoSubModuleFileAddr string
	//dao_sub_module_template（模板）
	DaoSubModuleTemplateInitFileAddr string
	DaoSubModuleTemplateAddFileAddr  string
}

func InitTemplateDir(templateDir string) {
	//co.TemplateDir = templateDir
	//co.Controller_templateUrl = filepath.Join(co.TemplateDir, "controller_template")
	//co.Controller_init_templateUrl = filepath.Join(co.TemplateDir, "controller_init_template")
	//co.Service_no_resp_templateUrl = filepath.Join(co.TemplateDir, "service_no_resp_template")
	//co.Service_resp_templateUrl = filepath.Join(co.TemplateDir, "service_resp_template")
	//co.Service_init_templateUrl = filepath.Join(co.TemplateDir, "service_init_template")
	//co.Dao_init_templateUrl = filepath.Join(co.TemplateDir, "dao_init_template")
	//co.ParentDao_init_templateUrl = filepath.Join(co.TemplateDir, "parent_dao_init_template")
	//co.ParentController_init_templateUrl = filepath.Join(co.TemplateDir, "parent_controller_init_template")
	//co.ParentService_init_templateUrl = filepath.Join(co.TemplateDir, "parent_service_init_template")
	//co.ParentService_resp_templateUrl = filepath.Join(co.TemplateDir, "parent_service_resp_template")
	//co.ParentService_no_resp_templateUrl = filepath.Join(co.TemplateDir, "parent_service_no_resp_template")
	//co.Api_sub_templateUrl = filepath.Join(co.TemplateDir, "api_sub_template")
	//co.Api_init_templateUrl = filepath.Join(co.TemplateDir, "api_init_template")
}

func InitWorkDir(workDir string) {
	// 使用filepath.Join处理路径分隔符
	co.WorkDir = workDir
	//co.AppUrl = filepath.Join(co.WorkDir, "app")
	//co.ApiUrl = filepath.Join(co.AppUrl, "api.go")
	//co.ControllerPackageName = filepath.Join("app", "controller")
	//co.ServicePackageName = filepath.Join("app", "service")
	//co.DaoPackageName = filepath.Join("app", "dao")
	//co.ControllerUrl = filepath.Join(co.AppUrl, co.ControllerPackageName, "controller.go")
	//co.ServiceUrl = filepath.Join(co.AppUrl, co.ServicePackageName, "service.go")
	//co.DaoUrl = filepath.Join(co.AppUrl, co.DaoPackageName, "dao.go")
	//co.ParentDaoUrl = filepath.Join(co.AppUrl, co.DaoPackageName, "parent_dao.go")
	//co.ParentControllerUrl = filepath.Join(co.AppUrl, co.ControllerPackageName, "parent_controller.go")
	//co.ParentServiceUrl = filepath.Join(co.AppUrl, co.ServicePackageName, "parent_service.go")
}

var co = Common{}

func GetCommon() *Common {
	return &co
}

func (c *SubAutoApi) init() {
	//if c.RequestType == "get" {
	//	c.RequestWay = "query"
	//	c.ShouldBindType = "ShouldBindQuery"
	//} else {
	//	c.RequestWay = "body"
	//	c.ShouldBindType = "ShouldBindJSON"
	//}
	//if c.RespName == "" {
	//	c.DaoRespName = "string"
	//	c.RespType = "string"
	//	c.RespDesc = co.SuccLogo
	//	c.RespDetail = "err"
	//	c.RespDetailName = `"` + co.SuccLogo + `"`
	//} else {
	//	c.RespType = "object"
	//	c.RespDesc = "响应体"
	//	c.RespDetail = "resp,err"
	//	c.RespDetailName = "resp"
	//	c.DaoRespName = co.DaoPackageName + "." + c.RespName
	//}
}
func stat(url string) bool {
	_, err := os.Stat(url)
	return err == nil
}
func findGoModDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		goModPath := filepath.Join(dir, "go.mod")
		if stat(goModPath) {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("找不到包含 go.mod 的目录")
		}
		dir = parent
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

var replace *Replace

func Initialize(settingsFileAddr, templatesFileAddr string) {
	replace = NewReplace(WriteRelationModel(Number_1_RelationModel))
	goModDir, err := findGoModDir()
	if err == nil {
		//初始化项目路径以及包路径
		replace.WriteReplaceKV(NewReplaceKV("project_dir_file_addr", goModDir))
		s, err := open(filepath.Join(goModDir, "go.mod"))
		if err == nil {
			for _, v := range strings.Split(s, "\n") {
				if strings.HasPrefix(v, "module") {
					replace.WriteReplaceKV(NewReplaceKV("project_dir_package_addr", strings.TrimSpace(strings.Split(v, " ")[1])))
					break
				}
			}
		}
	}

	replace.WriteReplaceKV(
		NewReplaceKV("templates_file_addr", templatesFileAddr),
		NewReplaceKV("settings_file_addr", settingsFileAddr))

	co = Common{}
	replace.Copy(&co)

	A = NewAutoApi()

	HandlerAllRelation(settingsFileAddr)
}

// HandlerAllRelation 处理所有关系
func HandlerAllRelation(settingsFileAddr string) {
	// 打开目录
	dir, err := os.Open(settingsFileAddr)
	if err != nil {
		fmt.Println("打开目录失败:", err)
		return
	}
	defer dir.Close()

	// 读取目录内容
	files, err := dir.Readdir(-1)
	if err != nil {
		fmt.Println("读取目录失败:", err)
		return
	}

	// 遍历目录中的文件和子目录
	for _, file := range files {
		if (!file.IsDir()) && strings.HasSuffix(file.Name(), "_relation") {
			HandlerOneFileRelation(filepath.Join(settingsFileAddr, file.Name()))
		}
	}
}

// HandlerOneFileRelation 处理单个关系文件
func HandlerOneFileRelation(addr string) {
	replace.WriteFileRelation(addr)
	fmt.Printf("关系文件%v已被读入\n", addr)
}
func RunSettingFront() {
	runSetting(filepath.Join(co.SettingsFileAddr, "gin_auto_setting_front.application"))
}
func RunSettingEnd() {
	runSetting(filepath.Join(co.SettingsFileAddr, "gin_auto_setting_end.application"))
}
func runSetting(addr string) {
	replace.WriteFile(addr)
	replace.Copy(&co)
}

var A *AutoApi

// 从api文件中抓取api内容
func GetApi(apiURL string) {
	s, err := open(apiURL)
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
	//把公共区域写入文件
	replace.WriteInterface(A.Co)
}

func getApiType(s string) {
	s = strings.TrimSpace(s)
	if s == "" {
		return
	}
	if s[len(s)-1] == '{' {
		A.Classes = append(A.Classes, &AutoApiClass{
			Name:    s[:len(s)-1],
			Context: "{",
		})
	} else if s[len(s)-1] == '}' {
		A.Classes[len(A.Classes)-1].Context += "\n" + s
	} else if len(A.Classes) != 0 && A.Classes[len(A.Classes)-1].
		Context[len(A.Classes[len(A.Classes)-1].Context)-1] != '}' {
		A.Classes[len(A.Classes)-1].Context += "\n\t" + s
	}
}

func getApiServer(s string) {
	split := strings.Split(s, ":")
	if len(split) == 2 {
		split[0] = strings.TrimSpace(split[0])
		split[1] = strings.Trim(strings.TrimSpace(split[1]), `"`)
		A.Co[split[0]] = split[1]
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
			template := &SubAutoApi{}
			template.Summary = strings.Trim(list[1], `"`)
			A.SubList = append(A.SubList, template)
		} else if list[0] == "@handler" {
			A.SubList[len(A.SubList)-1].FuncName = list[1]
		}
	} else if len(list) == 5 {
		A.SubList[len(A.SubList)-1].RequestType = list[0]
		A.SubList[len(A.SubList)-1].Url = list[1]
		A.SubList[len(A.SubList)-1].ReqName = list[2][1 : len(list[2])-1]
		A.SubList[len(A.SubList)-1].RespName = list[4][1 : len(list[4])-1]
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
func (r *Replace) setValue(k, v string) {
	r.dict[r.serialization(k)] = v
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
	cv := r.CV()
	cv.EnableRelation()
	return cv.replaceAll(s)
}

// 替换全部
func (r *Replace) replaceAll(s string) string {
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
	split := strings.Split(deleteComment(st), "=")
	if split != nil && len(split) == 2 {
		return &ReplaceKV{
			k: split[0],
			v: split[1],
		}, nil
	}
	return nil, errors.New("newReplaceKV | err")
}

// deleteComment 删除注释
func deleteComment(st string) string {
	index := strings.Index(st, "//")
	if index != -1 {
		st = st[:index]
	}
	return st
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
func WriteRelationModel(relationModel RelationModel) ReplaceOpts {
	return func(r *Replace) {
		r.relationModel = relationModel
	}
}
func (r *Replace) CV() *Replace {
	re := &Replace{}
	copier.Copy(re, r)
	return re
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
		r.setValue(kv.k, r.replaceAll(strings.TrimSpace(kv.v)))
	}
}

// 写入文件结构关系键值对
func (r *Replace) WriteFileContextRelation(context string) {
	split := strings.Split(context, "\n")
	vs1 := make([]*ReplaceKV, 0)
	vs2 := make([]*ReplaceKV, 0)
	vsP := 0
	for _, v := range split {
		v = strings.TrimSpace(v)
		if v == "IF" {
			if len(vs1)+len(vs2) > 0 {
				r.WriteReplaceKVRelation(vs1, vs2)
			}
			vs1 = make([]*ReplaceKV, 0)
			vs2 = make([]*ReplaceKV, 0)
			vsP = 1
			continue
		}
		if v == "->" {
			vsP = 2
			continue
		}
		if strings.ContainsAny(v, "=") {
			if vsP == 1 {
				if kv, err := NewReplaceKVByContext(v); err == nil {
					vs1 = append(vs1, kv)
				}
			} else if vsP == 2 {
				if kv, err := NewReplaceKVByContext(v); err == nil {
					vs2 = append(vs2, kv)
				}
			}
		}
	}
	if len(vs1)+len(vs2) > 0 {
		r.WriteReplaceKVRelation(vs1, vs2)
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
		//fmt.Println(kv.k, kv.v)
	}
	for _, kv := range new {
		kv.v = strings.TrimSpace(kv.v)
		kv.k = r.serialization(kv.k)
		//fmt.Println(kv.k, kv.v)
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
		if r.GetValue(kv.k) != r.replaceAll(kv.v) {
			return
		}
	}
	// 满足条件，替换关系
	for _, kv := range relation.new {
		r.setValue(kv.k, r.replaceAll(kv.v))
	}
}

func (r *Replace) EnableRelation() {
	for _, relation := range r.relations {
		r.runRelation(relation)
	}
	r.relations = r.relations[:0]
}
