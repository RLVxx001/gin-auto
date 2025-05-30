# gin-auto

# 0.

项目中所有的变量都采用覆盖式的方法，例如你在文件扫描下流程3中设置了一个变量，在文件扫描下流程4中也设置了那么变量值将为流程4中设置的值，你可以当成有一个map在存储变量。

### 项目结构图

```Go
//项目结构图
-api
    -{{api_module_name}}.go
-app
    -{{api_module_name}}
        -controller
            {{controller_subject_name}}.go   【主文件】
            {{controller_sub_module_name}}.go 【子模块文件】
        -service
            {{service_subject_name}}.go        【主文件】
            {{service_sub_module_name}}.go    【子模块文件】
        -dao
            {{dao_subject_name}}.go        【主文件】
            {{dao_sub_module_name}}.go    【子模块文件】

具体文件的生成请依据_file_addr结尾的变量名称
```

### 关系文件（settings下的所有_relation结尾的文件）

##### 1.语法

```go
IF
xx=xxx  //可有多行也可为空（为空的话默认true），=左边为变量名称，右边为值（值可用{{}}表达式）
->
xx=xxxx //可有多行，=左边为变量名称，右边为值（值可用{{}}表达式） 【如果上述的判断全部为true】
```

##### 2.样例演示

```Go
//controller以及service新增数据模板中使用
IF
->
RespType=object
dao_resp_name={{dao_dir_name}}.{{resp_name}}
resp_desc=响应体
resp_detail=resp,err
resp_detail_result=resp
//service增量数据模板中使用
service_resp_type_context=(*{{dao_dir_name}}.{{respName}}, error)
service_resp_detail_context=nil,nil

IF
respName=
->
resp_type=string
dao_resp_name=string
resp_desc=succ
resp_detail=err
resp_detail_result="succ"
//service增量数据模板中使用
service_resp_type_context=error
service_resp_detail_context=nil
```

##### 3.注意事项

在settings中的所有的配置文档中均支持注释（//）以及变量引用（{{}}）

在以_relation结尾的关系文件中，所有的关系均持久化（也就是只有在模板引用值或自身关系文档中才会发生变更）格式以IF (条件判断)【可以为空则必定为true】 -> (符合所有条件后的赋值操作)

# 1.文件扫描

### 1.init（初始化）（系统完成）（包含把setting文件夹下所有的以_relation结尾的关系文件导入）

project_dir_file_addr（项目文件地址）

project_dir_package_addr（项目包地址）

### 2.setting_front（设置开头）

由用户来定夺每个参数都是什么

### 3.api_setting（api中用户端设置）

由用户来定夺每个参数都是什么

### 4.setting_back（设置结尾）

由用户来定夺每个参数都是什么

### 5.api_list（api文件中各个接口信息）

Summary//描述

FuncName//方法名

RequestType//请求类型[get|post|put]

Url//请求地址

ReqName//请求参数名称

RespName//返回参数名称

# 2.生成前必定要有的环境变量

### 规范：

```Go
//_file_addr 表示文件地址
//_package_addr 表示引用的包地址
//_name 表示文件名称或其名称
//_dir表示文件夹
```

### 项目相关

```Go
//project 该项目
//project_dir_file_addr=   //（不填的话会向上找go.mod文件所在位置得到此文件地址）
//project_dir_package_addr= //（不填的话会向上找go.mod文件中找到该信息）
```

### 模块相关

```Go
//module 模块
module_name=test
```

### 子模块相关

```Go
//sub_module 子模块（一般在.api文件中提到）
//sub_module_name=
```

### api路由相关

##### api-路由文件包相关

```Go
//api_dir 路由文件
api_dir_name=api
api_dir_file_addr={{project_dir_file_addr}}/api
api_dir_package_addr={{project_dir_package_addr}}/api
```

##### api-路由模块文件相关

```Go
//api_module
api_module_name={{module_name}}
api_module_file_addr={{api_dir_file_addr}}/{{api_module_name}}.go
```

##### api-路由模块文件相关-模板

```Go
//api_module_template（模板）
api_module_template_init_file_addr={{templates_file_addr}}/api_module_template_init //初始化模板
api_module_template_add_file_addr={{templates_file_addr}}/api_module_template_add //新增模板
```

### app业务相关

#### app文件包相关

```Go
//app_dir
app_dir_name=app
app_dir_file_addr={{project_dir_file_addr}}/app
app_dir_package_addr={{project_dir_package_addr}}/app
```

#### controller相关

##### controller-文件包相关

```Go
//controller_dir
controller_dir_name=controller
controller_dir_file_addr={{app_dir_file_addr}}/{{module_name}}/{{controller_dir_name}}
controller_dir_package_addr={{app_dir_package_addr}}/{{module_name}}/{{controller_dir_name}}
```

##### controller-主文件相关

```Go
//controller_subject
controller_subject_name=controller
controller_subject_file_addr={{controller_dir_file_addr}}/{{controller_subject_name}}.go
```

##### controller-主文件相关-模板

```Go
//controller_subject_template（模板）
controller_subject_template_init_file_addr={{templates_file_addr}}/controller_subject_template_init //初始化模板
//（在实现代码中没有用到）controller_subject_template_add_file_addr={{templates_file_addr}}/controller_subject_template_add //新增模板
```

##### controller-子模块文件相关

```Go
//controller_sub_module
controller_sub_module_name={{sub_module_name}}
controller_sub_module_file_addr={{controller_dir_file_addr}}/{{controller_sub_module_name}}.go
```

##### controller-子模块文件相关-模板

```Go
//controller_sub_module_template（模板）
controller_sub_module_template_init_file_addr={{templates_file_addr}}/controller_sub_module_template_init //初始化模板
controller_sub_module_template_add_file_addr={{templates_file_addr}}/controller_sub_module_template_add //新增模板
```

#### service相关

##### service-文件包相关

```Go
//service_dir
service_dir_name=service
service_dir_file_addr={{app_dir_file_addr}}/{{module_name}}/{{service_dir_name}}
service_dir_package_addr={{app_dir_package_addr}}/{{module_name}}/{{service_dir_name}}
```

##### service-主文件相关

```Go
//service_subject
service_subject_name=service
service_subject_file_addr={{service_dir_file_addr}}/{{service_subject_name}}.go
service_subject_logo=type Services interface
```

##### service-主文件相关-模板

```Go
//service_subject_template（模板）
service_subject_template_init_file_addr={{templates_file_addr}}/service_subject_template_init //初始化模板
service_subject_template_add_file_addr={{templates_file_addr}}/service_subject_template_add //新增模板
```

##### service-子模块文件相关

```Go
//service_sub_module
service_sub_module_name={{sub_module_name}}
service_sub_module_file_addr={{service_dir_file_addr}}/{{service_sub_module_name}}.go
```

##### service-子模块文件相关-模版

```Go
//service_sub_module_template（模板）
service_sub_module_template_init_file_addr={{templates_file_addr}}/service_sub_module_template_init //初始化模板
service_sub_module_template_add_file_addr={{templates_file_addr}}/service_sub_module_template_add //新增模板
```

#### dao相关

##### dao-文件包相关

```Go
//dao_dir
dao_dir_name=dao
dao_dir_file_addr={{app_dir_file_addr}}/{{module_name}}/{{dao_dir_name}}
dao_dir_package_addr={{app_dir_package_addr}}/{{module_name}}/{{dao_dir_name}}
```

##### dao-主文件相关

```Go
//dao_subject
dao_subject_name=dao
dao_subject_file_addr={{dao_dir_file_addr}}/{{dao_subject_name}}.go
```

##### dao-主文件相关-模板

```Go
//dao_subject_template（模板）
dao_subject_template_init_file_addr={{templates_file_addr}}/dao_subject_template_init //初始化模板
//（在实现代码中没有用到）dao_subject_template_add_file_addr={{templates_file_addr}}/dao_subject_template_add //新增模板
```

##### dao-子模块文件相关

```Go
//dao_sub_module
dao_sub_module_name=type_{{sub_module_name}}
dao_sub_module_file_addr={{dao_dir_file_addr}}/{{dao_sub_module_name}}.go
```

##### dao-子模块文件相关-模板

```Go
//dao_sub_module_template（模板）
dao_sub_module_template_init_file_addr={{templates_file_addr}}/dao_sub_module_template_init //初始化模板
//（在实现代码中没有用到）dao_sub_module_template_add_file_addr={{templates_file_addr}}/dao_sub_module_template_add //新增模板
```

# 3.生成时

总的生成流程为：

先去重（根据controller/sub_module的文件中判断是否以及存在funcName的函数存在则去重操作）

将应该生成的文件分别

### 0.初始化

在初始化时判断api_module_file_addr、controller_subject_file_addr、service_subject_file_addr、dao_subject_file_addr是否存在，不存在则生成初始化模板

### 1.api_module_file_addr中生成增量（采用寻找最后一个}增量） 【原文件不存在则会先创建对应的初始化模板】

### 2.controller_sub_module_file_addr中生成增量（采用push增量）【原文件不存在则会先创建对应的初始化模板】

### 3.service_sub_module_file_addr中生成增量（采用push增量）【原文件不存在则会先创建对应的初始化模板】

### 4.service_subject_file_addr中生成增量（依据service_subject_logo找到增量位置）【原文件不存在则会先创建对应的初始化模板】

### 5.dao_sub_module_file_addr中生成增量（采用push增量）【原文件不存在则会先创建对应的初始化模板】

在生成过程中均会在最后赋值模板值时才将【关系文件】使用。
