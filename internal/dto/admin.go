package dto

import "github.com/saleh-ghazimoradi/Projectopher/internal/helper"

type AdminReviewUpdateReq struct {
	AdminReview string `json:"admin_review"`
}

type AdminReviewResp struct {
	RankingName string `json:"ranking_name"`
	AdminReview string `json:"admin_review"`
}

func ValidateAdminReview(v *helper.Validator, req *AdminReviewUpdateReq) {
	v.Check(req.AdminReview != "", "adminReview", "adminReview must be provided")
}
