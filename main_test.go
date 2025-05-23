package main

import (
	"testing"

	"github.com/RLVxx001/gin-auto/auto"
)

func TestMain(m *testing.M) {
	a := auto.NewAutoApi()
	a.AutoApiCommon = auto.AutoApiCommon{
		Tags:      "测试tags",
		UrlPre:    "/test1",
		UrlPrePre: "/test2",
		Version:   "v1",
	}
	a.ControllerTemplates = []*auto.ControllerTemplate{
		{
			FuncName: "GetUserInfo",
			ControllerTemplateSwagger: auto.ControllerTemplateSwagger{
				Summary:     "测试Summary",
				Description: "测试Description",
				ReqName:     "TestReq",
				Url:         "/tt",
				RequestType: "post",
			},
		},
		{
			FuncName: "GGG",
			ControllerTemplateSwagger: auto.ControllerTemplateSwagger{
				Summary:     "测试Summary22222",
				Description: "测试Description2222",
				ReqName:     "GGGReq",
				RespName:    "GGGResp",
				Url:         "/tt2222",
				RequestType: "get",
			},
		},
	}
	a.Classes = []*auto.AutoApiClass{
		{
			Context: "type TestReq struct{}",
			Name:    "TestReq",
		},
		{
			Context: "type GGGReq struct{}",
			Name:    "GGGReq",
		},
		{
			Context: "type GGGResp struct{}",
			Name:    "GGGResp",
		},
	}
	a.InsertContext()
}
