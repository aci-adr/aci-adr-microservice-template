package response

type Error struct {
	Code    string `json:"code,required"`
	Message string `json:"message"`
	Details string `json:"details"`
}
