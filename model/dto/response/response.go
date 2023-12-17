package response

type ResponseWithSimpleData[T any] struct {
	Data   *T         `json:"data"`
	Status StatusCode `json:"status"`
	Errors *[]Error   `json:"errors"`
}

type ResponseWithArrayData[T any] struct {
	Data   *[]T       `json:"data"`
	Status StatusCode `json:"status"`
	Errors *[]Error   `json:"errors"`
}
