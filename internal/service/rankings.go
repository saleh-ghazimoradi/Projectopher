package service

import (
	"context"
	"github.com/saleh-ghazimoradi/Projectopher/internal/domain"
	"github.com/saleh-ghazimoradi/Projectopher/internal/dto"
	"github.com/saleh-ghazimoradi/Projectopher/internal/repository"
)

type RankingService interface {
	GetRankings(ctx context.Context) ([]dto.Ranking, error)
}

type rankingService struct {
	rankRepository repository.RankingRepository
}

func (r *rankingService) GetRankings(ctx context.Context) ([]dto.Ranking, error) {
	rankings, err := r.rankRepository.GetRankings(ctx)
	if err != nil {
		return nil, err
	}

	response := make([]dto.Ranking, len(rankings))
	for i := range rankings {
		response[i] = *r.toRanking(&rankings[i])
	}
	
	return response, nil
}

func (r *rankingService) toRanking(ranking *domain.Ranking) *dto.Ranking {
	return &dto.Ranking{
		RankingValue: ranking.RankingValue,
		RankingName:  ranking.RankingName,
	}
}

func NewRankingsService(rankRepository repository.RankingRepository) RankingService {
	return &rankingService{
		rankRepository: rankRepository,
	}
}
