package response

type StatusCode string

const (
	Success       StatusCode = "Success"
	BadRequest    StatusCode = "BadRequest"
	InternalError StatusCode = "InternalServerError"
	NotFound      StatusCode = "NotFound"
)
