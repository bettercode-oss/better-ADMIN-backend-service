package dtos

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

type PageResult struct {
	Result     interface{} `json:"result"`
	TotalCount int64       `json:"totalCount"`
}

const PageSize = 20

type Pageable struct {
	Page     int
	PageSize int
}

func (p Pageable) GetOffset() int {
	return (p.Page - 1) * p.PageSize
}

func NewPageableFromRequest(ctx *gin.Context) Pageable {
	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		page = 1
	}

	pageSize, err := strconv.Atoi(ctx.Query("pageSize"))
	if err != nil {
		pageSize = PageSize
	}

	return Pageable{
		Page:     page,
		PageSize: pageSize,
	}
}
