package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"

	commonmodels "github.com/koderover/zadig/pkg/microservice/aslan/core/common/repository/models"
	"github.com/koderover/zadig/pkg/microservice/aslan/core/system/service"
	internalhandler "github.com/koderover/zadig/pkg/shared/handler"
	e "github.com/koderover/zadig/pkg/tool/errors"
	"github.com/koderover/zadig/pkg/tool/log"
)

func ListExternalLinks(c *gin.Context) {
	ctx := internalhandler.NewContext(c)
	defer func() { internalhandler.JSONResponse(c, ctx) }()

	ctx.Resp, ctx.Err = service.ListExternalLinks(ctx.Logger)
}

func CreateExternalLink(c *gin.Context) {
	ctx := internalhandler.NewContext(c)
	defer func() { internalhandler.JSONResponse(c, ctx) }()

	args := new(commonmodels.ExternalLink)
	data, err := c.GetRawData()
	if err != nil {
		log.Errorf("CreateExternalLink c.GetRawData() err : %s", err)
	}
	if err = json.Unmarshal(data, args); err != nil {
		log.Errorf("CreateExternalLink json.Unmarshal err : %s", err)
	}
	internalhandler.InsertOperationLog(c, ctx.UserName, "", "新增", "系统配置-快捷链接", fmt.Sprintf("name:%s url:%s", args.Name, args.URL), string(data), ctx.Logger)

	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data))

	if err := c.ShouldBindJSON(&args); err != nil {
		ctx.Err = e.ErrInvalidParam.AddDesc("invalid externalLink args")
		return
	}
	args.UpdateBy = ctx.UserName

	ctx.Err = service.CreateExternalLink(args, ctx.Logger)
}

func UpdateExternalLink(c *gin.Context) {
	ctx := internalhandler.NewContext(c)
	defer func() { internalhandler.JSONResponse(c, ctx) }()

	args := new(commonmodels.ExternalLink)
	data, err := c.GetRawData()
	if err != nil {
		log.Errorf("UpdateExternal c.GetRawData() err : %s", err)
	}
	if err = json.Unmarshal(data, args); err != nil {
		log.Errorf("UpdateExternal json.Unmarshal err : %s", err)
	}
	internalhandler.InsertOperationLog(c, ctx.UserName, "", "更新", "系统配置-快捷链接", fmt.Sprintf("name:%s url:%s", args.Name, args.URL), string(data), ctx.Logger)

	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data))

	if err := c.ShouldBindJSON(&args); err != nil {
		ctx.Err = e.ErrInvalidParam.AddDesc("invalid externalLink args")
		return
	}
	args.UpdateBy = ctx.UserName

	ctx.Err = service.UpdateExternalLink(c.Param("id"), args, ctx.Logger)
}

func DeleteExternalLink(c *gin.Context) {
	ctx := internalhandler.NewContext(c)
	defer func() { internalhandler.JSONResponse(c, ctx) }()

	internalhandler.InsertOperationLog(c, ctx.UserName, "", "删除", "系统配置-快捷链接", fmt.Sprintf("id:%s", c.Param("id")), "", ctx.Logger)
	ctx.Err = service.DeleteExternalLink(c.Param("id"), ctx.Logger)
}
