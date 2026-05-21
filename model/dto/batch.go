package dto

// BatchOperationDto 批量操作请求
type BatchOperationDto struct {
	IDs        []uint `json:"ids" form:"ids" validate:"required,min=1"`
	Action     string `json:"action" form:"action" validate:"required,oneof=enable disable delete extend reset"`
	Days       *int   `json:"days" form:"days" validate:"omitempty,gte=1,lte=365"`       // 续期天数
	QuotaGB    *int   `json:"quotaGB" form:"quotaGB" validate:"omitempty,gte=-1,lte=1024"` // 流量GB (-1=无限)
}

// BatchResetTrafficDto 批量重置流量
type BatchResetTrafficDto struct {
	IDs []uint `json:"ids" form:"ids" validate:"required,min=1"`
}
