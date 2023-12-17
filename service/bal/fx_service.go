package bal

import (
	"FxService/model/dto/request"
	"FxService/model/dto/response"
	"FxService/model/entity"
	"FxService/service/common"
	"FxService/service/dal"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Fx_service struct {
	DbService dal.DBService[entity.ForexData]
}

func getForexDtoFromEntity(result entity.ForexData) *response.ForexDataResponse {
	return &response.ForexDataResponse{
		Id:                           result.ID.Hex(),
		TenantId:                     result.TenantID,
		BankId:                       result.BankID,
		BaseCurrency:                 result.BaseCurrency,
		TargetCurrency:               result.TargetCurrency,
		Tier:                         result.Tier,
		DirectIndirectFlag:           result.DirectIndirectFlag,
		Multiplier:                   result.Multiplier,
		BuyRate:                      result.BuyRate,
		SellRate:                     result.SellRate,
		TolerancePercentage:          result.TolerancePercentage,
		EffectiveDate:                result.EffectiveDate,
		ExpirationDate:               result.ExpirationDate,
		ContractRequirementThreshold: result.ContractRequirementThreshold,
	}
}

func (s *Fx_service) CreateForexData(forexData request.CreateForexDataRequest) response.ResponseWithSimpleData[response.ForexDataResponse] {
	var dbObject = entity.ForexData{
		ID:                           primitive.NewObjectID(),
		Tier:                         forexData.Tier,
		DirectIndirectFlag:           forexData.DirectIndirectFlag,
		Multiplier:                   forexData.Multiplier,
		BuyRate:                      forexData.BuyRate,
		SellRate:                     forexData.SellRate,
		TolerancePercentage:          forexData.TolerancePercentage,
		EffectiveDate:                forexData.EffectiveDate,
		ExpirationDate:               forexData.ExpirationDate,
		ContractRequirementThreshold: forexData.ContractRequirementThreshold,
		TenantID:                     forexData.TenantId,
		BankID:                       forexData.BankId,
		BaseCurrency:                 forexData.BaseCurrency,
		TargetCurrency:               forexData.TargetCurrency,
		CreatedDate:                  time.Now(),
		DocVersion:                   1,
		UpdatedDate:                  time.Now(),
	}
	common.Logger.Info("Create a forex record started")
	result, err := s.DbService.CreateOne(dbObject)
	common.Logger.Info("Create a forex record ended")
	if err != nil {
		common.Logger.Errorf("Error in creating a new Record. Exception:%v", err)
		e := &[]response.Error{
			{Code: "FAILURE", Message: "Unable to create record", Details: "Unable to create record due to some exception."},
		}
		return common.GetSimpleResponse[response.ForexDataResponse](nil, response.InternalError, e)
	}

	return common.GetSimpleResponse[response.ForexDataResponse](getForexDtoFromEntity(result), response.Success, nil)
}

func (s *Fx_service) GetForexRateById(id string) response.ResponseWithSimpleData[response.ForexDataResponse] {
	objectId, _ := primitive.ObjectIDFromHex(id)
	result, err := s.DbService.GetOne(bson.D{{Key: "_id", Value: objectId}})

	if err != nil {
		common.Logger.Errorf("Error in retriving forex rate by id. Exception:%v", err)
		e := &[]response.Error{
			{Code: "DATA_NOT_FOUND", Message: "No record found", Details: "No record found"},
		}
		return common.GetSimpleResponse[response.ForexDataResponse](nil, response.NotFound, e)
	}

	return common.GetSimpleResponse[response.ForexDataResponse](getForexDtoFromEntity(result), response.Success, nil)
}

func (s *Fx_service) DeleteForexRateById(id string) response.ResponseWithSimpleData[response.ForexDataResponse] {
	objectId, _ := primitive.ObjectIDFromHex(id)
	result, err := s.DbService.DeleteOne(bson.D{{Key: "_id", Value: objectId}})

	if err != nil || result == 0 {
		e := &[]response.Error{
			{Code: "DATA_NOT_FOUND", Message: "No record Deleted", Details: "No record Deleted"},
		}
		common.Logger.Errorf("Error in Deleting forex rate by id. Exception:%v", err)
		return common.GetSimpleResponse[response.ForexDataResponse](nil, response.NotFound, e)
	}

	return common.GetSimpleResponse[response.ForexDataResponse](nil, response.Success, nil)
}

func (s *Fx_service) GetForexRateByFilter(tenantId int, bankId int, baseCurrency string, targetCurrency string) response.ResponseWithArrayData[response.ForexDataResponse] {
	result, err := s.DbService.Get(bson.D{
		{Key: "tenantId", Value: tenantId},
		{Key: "bankId", Value: bankId},
		{Key: "baseCurrency", Value: baseCurrency},
		{Key: "targetCurrency", Value: targetCurrency},
	})

	if err != nil {
		e := &[]response.Error{
			{Code: "DATA_NOT_FOUND", Message: "No record found", Details: "No record found"},
		}
		common.Logger.Errorf("Error in retriving forex rate by filters. Exception:%v", err)
		return common.GetArrayResponse[response.ForexDataResponse](nil, response.NotFound, e)
	}

	var data []response.ForexDataResponse

	for _, item := range result {
		data = append(data, *getForexDtoFromEntity(item))
	}

	return common.GetArrayResponse[response.ForexDataResponse](&data, response.Success, nil)

}

func (s *Fx_service) UpdateForexRateById(id string, body request.UpdateForexDataRequest) response.ResponseWithSimpleData[response.ForexDataResponse] {
	objectId, _ := primitive.ObjectIDFromHex(id)

	updateDocument := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "sellRate", Value: body.SellRate},
			{Key: "buyRate", Value: body.BuyRate},
			{Key: "expirationDate", Value: body.ExpirationDate},
			{Key: "effectiveDate", Value: body.EffectiveDate},
			{Key: "tolerancePercentage", Value: body.TolerancePercentage},
			{Key: "multiplier", Value: body.Multiplier},
			{Key: "directIndirectFlag", Value: body.DirectIndirectFlag},
			{Key: "contractRequirementThreshold", Value: body.ContractRequirementThreshold},
			{Key: "updatedDate", Value: time.Now()},
		}},
	}

	result, err := s.DbService.UpdateOne(updateDocument, bson.D{{Key: "_id", Value: objectId}})

	if err != nil {
		e := &[]response.Error{
			{Code: "DATA_NOT_FOUND", Message: "No record Updated", Details: "No record Updated"},
		}
		common.Logger.Errorf("Error in updating forex rate by id. Exception:%v", err)
		return common.GetSimpleResponse[response.ForexDataResponse](nil, response.NotFound, e)
	}

	return common.GetSimpleResponse[response.ForexDataResponse](getForexDtoFromEntity(result), response.Success, nil)
}

func (s *Fx_service) GetConvertedRate(tenantId int, bankId int, amount float64, baseCurrency string, targetCurrency string, tier string) response.ResponseWithSimpleData[response.ConversionResponse] {
	filter := bson.D{
		{Key: "tenantId", Value: tenantId},
		{Key: "bankId", Value: bankId},
		{Key: "baseCurrency", Value: baseCurrency},
		{Key: "targetCurrency", Value: targetCurrency},
		{Key: "tier", Value: tier},
	}

	result, err := s.DbService.GetOne(filter)

	if err != nil {
		e := &[]response.Error{
			{Code: "DATA_NOT_FOUND", Message: "No record found", Details: "No record found"},
		}
		common.Logger.Errorf("Error in retriving and converting forex rate. Exception:%v", err)
		return common.GetSimpleResponse[response.ConversionResponse](nil, response.NotFound, e)
	}

	resp := response.ConversionResponse{
		Amount:          amount,
		ConvertedAmount: amount * result.BuyRate,
		BaseCurrency:    baseCurrency,
		TargetCurrency:  targetCurrency,
		InitiatedOn:     int64(time.Nanosecond),
		Rate:            result.BuyRate,
	}
	return common.GetSimpleResponse[response.ConversionResponse](&resp, response.Success, nil)
}
