package utils

type Pagination struct {
	Page         int32 `json:"page"`
	Limit        int32 `json:"limit"`
	TotalRecords int32 `json:"total_records"`
	TotalPages   int32 `json:"total_pages"`
	HasNext      bool  `json:"has_next"`
	HasPrev      bool  `json:"has_prev"`
}

func NewPagination(page, limit, totalRecords int32) *Pagination {
	if page <= 0 {
		page = 1
	}

	if limit <= 0 {
		limit = int32(GetIntEnv("LIMIT_ITEM_ON_PER_PAGE", 10))
	}

	totalPages := (totalRecords + limit - 1) / limit

	return &Pagination{
		Page:         page,
		Limit:        limit,
		TotalRecords: totalRecords,
		TotalPages:   totalPages,
		HasNext:      page < totalPages,
		HasPrev:      page > 1,
	}
}

func NewPaginationResponse(data any, page, limit, totalRecords int32) map[string]any {
	return map[string]any{
		"data":       data,
		"pagination": NewPagination(page, limit, totalRecords),
	}
}
