package controller

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"ElasticView/engine/es"
	"ElasticView/model"
	"ElasticView/platform-basic-libs/jwt"
	"ElasticView/platform-basic-libs/response"

	"github.com/cch123/elasticsql"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic"
)

type EsController struct {
	BaseController
}

//Ping
func (this EsController) PingAction(ctx *gin.Context) {
	esConnect := es.EsConnect{}
	err = ctx.Bind(&esConnect)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	esClinet, err := es.GetEsClient(esConnect)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	data, _, err := esClinet.Ping()
	if err != nil {
		this.Error(ctx, err)
		return
	}
	this.Success(ctx, response.OperateSuccess, data)
}

//Elasticsearch状态
func (this EsController) CatAction(ctx *gin.Context) {

	esCat := es.EsCat{}
	err = ctx.Bind(&esCat)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	esClinet, err := es.GetEsClientByID(esCat.EsConnect)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	var data interface{}

	switch esCat.Cat {
	case "CatHealth":
		data, err = esClinet.(*es.EsClientV6).Client.CatHealth().Human(true).Do(context.Background())
	case "CatShards":
		data, err = esClinet.(*es.EsClientV6).Client.CatShards().Human(true).Do(context.Background())
	case "CatCount":
		data, err = esClinet.(*es.EsClientV6).Client.CatCount().Human(true).Do(context.Background())
	case "CatAllocation":
		data, err = esClinet.(*es.EsClientV6).Client.CatAllocation().Human(true).Do(context.Background())
	case "CatAliases":
		data, err = esClinet.(*es.EsClientV6).Client.CatAliases().Human(true).Do(context.Background())
	case "CatIndices":
		data, err = esClinet.(*es.EsClientV6).Client.CatIndices().Human(true).Do(context.Background())
	}

	if err != nil {
		this.Error(ctx, err)
		return
	}

	this.Success(ctx, response.SearchSuccess, data)
}

func (this EsController) RunDslAction(ctx *gin.Context) {
	esRest := es.EsRest{}
	err = ctx.Bind(&esRest)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	esClinet, err := es.GetEsClientByID(esRest.EsConnect)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	esRest.Method = strings.ToUpper(esRest.Method)
	if esRest.Method == "GET" {
		c, err := jwt.ParseToken(ctx.GetHeader("X-Token"))
		if err != nil {
			this.Error(ctx, err)
			return
		}

		gmDslHistoryModel := model.GmDslHistoryModel{
			Uid:    int(c.ID),
			Method: esRest.Method,
			Path:   esRest.Path,
			Body:   esRest.Body,
		}

		err = gmDslHistoryModel.Insert()

		if err != nil {
			this.Error(ctx, err)
			return
		}
	}

	res, err := esClinet.(*es.EsClientV6).Client.PerformRequest(context.TODO(), elastic.PerformRequestOptions{
		Method: esRest.Method,
		Path:   esRest.Path,
		Body:   esRest.Body,
	})

	if err != nil {
		this.Error(ctx, err)
		return
	}

	if res.StatusCode != 200 && res.StatusCode != 201 {
		this.Output(ctx, map[string]interface{}{
			"code": res.StatusCode,
			"msg":  fmt.Sprintf("请求异常! 错误码 :" + strconv.Itoa(res.StatusCode)),
			"data": res.Body,
		})
		return
	}

	this.Success(ctx, response.OperateSuccess, res.Body)
}

func (this EsController) SqlToDslAction(ctx *gin.Context) {
	sql := ctx.Request.FormValue("sql")
	dsl, table, err := elasticsql.ConvertPretty(sql)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	this.Success(ctx, "转换成功!", map[string]interface{}{
		"dsl":       dsl,
		"tableName": table,
	})
}
