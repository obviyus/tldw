package form

type SummaryQuery struct {
	VideoID string `form:"timestamp" binding:"required"`
}
