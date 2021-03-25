package controller

import (
	"context"
	"errors"
	"sync"

	"ElasticView/engine/es"
	"ElasticView/engine/logs"
	"ElasticView/platform-basic-libs/my_error"
	"ElasticView/platform-basic-libs/response"

	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic"
)

type EsIndexController struct {
	BaseController
}

//创建索引
func (this EsIndexController) CreateAction(ctx *gin.Context) {
	esIndexInfo := es.EsIndexInfo{}
	err = ctx.Bind(&esIndexInfo)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	esClinet, err := es.GetEsClient(esIndexInfo.EsConnect)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	if esIndexInfo.IndexName == "" {
		this.Error(ctx, my_error.NewBusiness(es.ParmasNullError, es.IndexNameNullError))
		return
	}
	logs.Logger.Sugar().Infof("esIndexInfo.Types", esIndexInfo.Types)
	if esIndexInfo.Types == "update" {
		res, err := esClinet.IndexPutSettings(esIndexInfo.IndexName, esIndexInfo.Settings)
		if err != nil {
			this.Error(ctx, err)
			return
		}
		this.Success(ctx, response.OperateSuccess, res)
		return
	} else {
		res, err := esClinet.(*es.EsClientV6).Client.CreateIndex(esIndexInfo.IndexName).BodyJson(map[string]interface{}{
			"settings": esIndexInfo.Settings,
		}).Do(context.Background())
		if err != nil {
			this.Error(ctx, err)
			return
		}
		this.Success(ctx, response.OperateSuccess, res)
		return
	}
	return
}

//创建索引
func (this EsIndexController) GetSettingsAction(ctx *gin.Context) {
	esIndexInfo := es.EsIndexInfo{}
	err = ctx.Bind(&esIndexInfo)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	esClinet, err := es.GetEsClient(esIndexInfo.EsConnect)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	if esIndexInfo.IndexName == "" {
		this.Error(ctx, my_error.NewBusiness(es.ParmasNullError, es.IndexNameNullError))
		return
	}

	res, err := esClinet.(*es.EsClientV6).Client.IndexGetSettings(esIndexInfo.IndexName).Do(context.Background())
	if err != nil {
		this.Error(ctx, err)
		return
	}

	this.Success(ctx, response.OperateSuccess, res[esIndexInfo.IndexName].Settings)
	return
}

//获取别名
func (this EsIndexController) GetAliasAction(ctx *gin.Context) {
	esAliasInfo := es.EsAliasInfo{}
	err = ctx.Bind(&esAliasInfo)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	esClinet, err := es.GetEsClient(esAliasInfo.EsConnect)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	if esAliasInfo.IndexName == "" {
		this.Error(ctx, my_error.NewBusiness(es.ParmasNullError, es.IndexNameNullError))
		return
	}

	aliasRes, err := esClinet.(*es.EsClientV6).Client.Aliases().Index(esAliasInfo.IndexName).Do(context.TODO())

	this.Success(ctx, response.OperateSuccess, aliasRes.Indices[esAliasInfo.IndexName].Aliases)
	return
}

func (this EsIndexController) OperateAliasAction(ctx *gin.Context) {
	esAliasInfo := es.EsAliasInfo{}
	err = ctx.Bind(&esAliasInfo)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	esClinet, err := es.GetEsClient(esAliasInfo.EsConnect)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	const Add = 1
	const Delete = 2
	const MoveToAnotherIndex = 3
	const PatchAdd = 4
	var res interface{}
	switch esAliasInfo.Types {
	case Add:
		if esAliasInfo.IndexName == "" {
			this.Error(ctx, my_error.NewBusiness(es.ParmasNullError, es.IndexNameNullError))
			return
		}
		res, err = esClinet.(*es.EsClientV6).Client.Alias().Add(esAliasInfo.IndexName, esAliasInfo.AliasName).Do(context.TODO())
	case Delete:
		if esAliasInfo.IndexName == "" {
			this.Error(ctx, my_error.NewBusiness(es.ParmasNullError, es.IndexNameNullError))
			return
		}
		res, err = esClinet.(*es.EsClientV6).Client.Alias().Remove(esAliasInfo.IndexName, esAliasInfo.AliasName).Do(context.TODO())
	case MoveToAnotherIndex:
		res, err = esClinet.(*es.EsClientV6).Client.Alias().Action(elastic.NewAliasAddAction(esAliasInfo.AliasName).Index(esAliasInfo.NewIndexList...)).Do(context.TODO())
	case PatchAdd:
		if esAliasInfo.IndexName == "" {
			this.Error(ctx, my_error.NewBusiness(es.ParmasNullError, es.IndexNameNullError))
			return
		}
		wg := sync.WaitGroup{}
		NewAliasNameListLen := len(esAliasInfo.NewAliasNameList)
		if len(esAliasInfo.NewAliasNameList) > 10 {
			err = errors.New("别名列表数量不能大于10")
			break
		} else {
			wg.Add(NewAliasNameListLen)
			for _, aliasName := range esAliasInfo.NewAliasNameList {
				go func(aliasName string) {
					defer wg.Done()
					res, err = esClinet.(*es.EsClientV6).Client.Alias().
						Add(esAliasInfo.IndexName, aliasName).
						Do(context.TODO())
				}(aliasName)
			}
			wg.Wait()
		}
	default:
		err = es.ReqParmasValid
	}

	if err != nil {
		this.Error(ctx, err)
		return
	}

	this.Success(ctx, response.OperateSuccess, res)
	return
}

func (this EsIndexController) ReindexAction(ctx *gin.Context) {
	esReIndexInfo := es.EsReIndexInfo{}
	err = ctx.Bind(&esReIndexInfo)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	esClinet, err := es.GetEsClient(esReIndexInfo.EsConnect)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	reindex := esClinet.(*es.EsClientV6).Client.Reindex()
	urlValues := esReIndexInfo.UrlValues
	if urlValues.WaitForActiveShards != "" {
		reindex = reindex.WaitForActiveShards(urlValues.WaitForActiveShards)
	}
	if urlValues.Slices != 0 {
		reindex = reindex.Slices(urlValues.Slices)
	}
	if urlValues.Refresh != "" {
		reindex = reindex.Refresh(urlValues.Refresh)
	}
	if urlValues.Timeout != "" {
		reindex = reindex.Timeout(urlValues.Refresh)
	}
	if urlValues.RequestsPerSecond != 0 {
		reindex = reindex.RequestsPerSecond(urlValues.RequestsPerSecond)
	}
	if urlValues.WaitForCompletion != nil {
		reindex = reindex.WaitForCompletion(*urlValues.WaitForCompletion)
	}

	res, err := reindex.Body(esReIndexInfo.Body).Do(context.Background())
	if err != nil {
		this.Error(ctx, err)
		return
	}
	this.Success(ctx, response.OperateSuccess, res)
}

func (this EsIndexController) ReindexListAction(ctx *gin.Context) {

}

func (this EsIndexController) IndexNamesAction(ctx *gin.Context) {
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
	indexNames, err := esClinet.(*es.EsClientV6).Client.IndexNames()
	if err != nil {
		this.Error(ctx, err)
		return
	}
	this.Success(ctx, response.SearchSuccess, indexNames)
	return
}
