package common

import "FxService/model/dto/response"

func GetSimpleResponse[T any](data *T, statusCode response.StatusCode, errors *[]response.Error) response.ResponseWithSimpleData[T] {
	return response.ResponseWithSimpleData[T]{
		Data:   data,
		Status: statusCode,
		Errors: errors,
	}
}

func GetArrayResponse[T any](data *[]T, statusCode response.StatusCode, errors *[]response.Error) response.ResponseWithArrayData[T] {
	return response.ResponseWithArrayData[T]{
		Data:   data,
		Status: statusCode,
		Errors: errors,
	}
}
