package controllers

import (
	"FxService/model/dto/request"
	"FxService/model/dto/response"
	"FxService/model/entity"
	"FxService/model/validation"
	"FxService/service/bal"
	"FxService/service/common"
	"FxService/service/dal"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

var dbService = &dal.MongoDbService[entity.ForexData]{}

var fxService = &bal.Fx_service{
	DbService: dbService,
}

func CreateForexRate(c *gin.Context) {
	if body, err := common.ValidateAndReturnBody[request.CreateForexDataRequest](c); err == nil {
		c.IndentedJSON(http.StatusOK, fxService.CreateForexData(body))
	}
}

func GetForexRateById(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, fxService.GetForexRateById(c.Param("id")))
}

func DeleteForexRateById(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, fxService.DeleteForexRateById(c.Param("id")))
}

func GetForexRateByFilter(c *gin.Context) {
	tenantId, _ := strconv.Atoi(c.Query("tenantId"))
	bankId, _ := strconv.Atoi(c.Query("bankId"))
	baseCurrency := c.Query("baseCurrency")
	targetCurrency := c.Query("targetCurrency")
	c.IndentedJSON(http.StatusOK, fxService.GetForexRateByFilter(tenantId, bankId, baseCurrency, targetCurrency))
}

func UpdateForexRateById(c *gin.Context) {
	if body, err := common.ValidateAndReturnBody[request.UpdateForexDataRequest](c); err == nil {
		c.IndentedJSON(http.StatusOK, fxService.UpdateForexRateById(c.Param("id"), body))
	}
}

func GetConvertedRate(c *gin.Context) {
	tenantId, _ := strconv.Atoi(c.Query("tenantId"))
	bankId, _ := strconv.Atoi(c.Query("bankId"))
	amount, _ := strconv.ParseFloat(c.Query("amount"), 64)
	baseCurrency := c.Query("baseCurrency")
	targetCurrency := c.Query("targetCurrency")
	tier := c.Query("tier")

	c.IndentedJSON(http.StatusOK, fxService.GetConvertedRate(tenantId, bankId, amount, baseCurrency, targetCurrency, tier))
}

func AddRoutes(e *gin.Engine) {
	e.GET("/api/forexrates",
		common.ParamValidationMiddleware[response.ForexDataResponse]([]validation.ValidationRule{
			{ParamName: "tenantId", Required: true, ParamType: "int"},
			{ParamName: "bankId", Required: true, ParamType: "int"},
			{ParamName: "baseCurrency", Required: true, ParamType: "string"},
			{ParamName: "targetCurrency", Required: true, ParamType: "string"},
		}),
		GetForexRateByFilter)
	e.GET("/api/forexrates/:id", GetForexRateById)
	e.GET("/api/convert",
		common.ParamValidationMiddleware[response.ConversionResponse]([]validation.ValidationRule{
			{ParamName: "tenantId", Required: true, ParamType: "int"},
			{ParamName: "bankId", Required: true, ParamType: "int"},
			{ParamName: "amount", Required: true, ParamType: "int"},
			{ParamName: "baseCurrency", Required: true, ParamType: "string"},
			{ParamName: "targetCurrency", Required: true, ParamType: "string"},
			{ParamName: "tier", Required: true, ParamType: "string"},
		}),
		GetConvertedRate)
	e.POST("/api/forexrates", CreateForexRate)
	e.DELETE("/api/forexrates/:id", DeleteForexRateById)
	e.PUT("/api/forexrates/:id", UpdateForexRateById)
}
