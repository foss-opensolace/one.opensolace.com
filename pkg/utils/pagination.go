package utils

import (
	"math"

	"gorm.io/gorm"
)

type Pagination[T any] struct {
	Data []T `json:"data"`

	RecordsPerPage int `json:"records_per_page"`
	RecordsInPage  int `json:"records_in_page"`

	TotalRecords int `json:"total_records"`
	TotalPages   int `json:"total_pages"`

	PreviousPage *int `json:"previous_page"`
	CurrentPage  int  `json:"current_page"`
	NextPage     *int `json:"next_page"`
}

// Get a pagination with the amount per page
func Paginate[T any](query *gorm.DB, page, amount int) (*Pagination[T], error) {
	var totalRecords int64
	if err := query.Count(&totalRecords).Error; err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(amount)))

	if page < 1 {
		page = 1
	} else if page > totalPages {
		page = totalPages
	}

	if amount < 1 {
		amount = 1
	} else if amount > 100 {
		amount = 100
	}

	data := make([]T, 0, amount)

	if err := query.Limit(amount).Offset((page - 1) * amount).Scan(&data).Error; err != nil {
		return nil, err
	}

	pagination := &Pagination[T]{
		Data:           data,
		RecordsPerPage: amount,
		RecordsInPage:  len(data),
		TotalRecords:   int(totalRecords),
		TotalPages:     int(totalPages),
		PreviousPage:   ToPtr(page - 1),
		CurrentPage:    page,
		NextPage:       ToPtr(page + 1),
	}

	if *pagination.PreviousPage < 1 {
		pagination.PreviousPage = nil
	}

	if *pagination.NextPage > totalPages {
		pagination.NextPage = nil
	}

	return pagination, nil
}
