package mongoDTO

import (
	"github.com/saleh-ghazimoradi/Projectopher/internal/domain"
)

type RankingDTO struct {
	RankingValue int    `bson:"ranking_value"`
	RankingName  string `bson:"ranking_name"`
}

func FromRankingCoreToDTO(input *domain.Ranking) *RankingDTO {
	return &RankingDTO{
		RankingValue: input.RankingValue,
		RankingName:  input.RankingName,
	}
}

func FromRankingDTOToCore(input *RankingDTO) *domain.Ranking {
	return &domain.Ranking{
		RankingValue: input.RankingValue,
		RankingName:  input.RankingName,
	}
}
