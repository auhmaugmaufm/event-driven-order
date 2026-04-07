package dto

type PaginationRequest struct {
	Limit  int    `json:"limit" query:"limit"`
	Offset int    `query:"offset"`
	Status string `query:"status"`
}

type PaginationResponse struct {
	Data        interface{} `json:"data"`
	TotalItems  int64       `json:"total_items"`
	TotalPages  int         `json:"total_pages"`
	CurrentPage int         `json:"current_page"`
}
