package domain

type Ranking struct {
	RankingValue int    `bson:"ranking_value"`
	RankingName  string `bson:"ranking_name"`
}
