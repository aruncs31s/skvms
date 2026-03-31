package dto

type ListResponse[T any] struct {
	List       []T   `json:"list"`
	TotalCount int64 `json:"total_count"`
}
