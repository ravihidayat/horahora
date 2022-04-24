package models

import (
	"strings"

	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	_ "github.com/horahoradev/horahora/user_service/protocol"
	videoproto "github.com/horahoradev/horahora/video_service/protocol"
	"github.com/jmoiron/sqlx"
)

// TODO: test suite for recommender implementations with precision and recall for sample dataset
type Recommender interface {
	GetRecommendations(userID int64) ([]*videoproto.VideoRec, error)
}

// Dumb recommender system, computes expected rating value for user from a video's tags
// and orders by sum
// No more train otomads??? (please)
type BayesianTagSum struct {
	db *sqlx.DB
}

func NewBayesianTagSum(db *sqlx.DB) BayesianTagSum {
	return BayesianTagSum{
		db: db,
	}
}

func (b *BayesianTagSum) GetRecommendations(uid int64) ([]*videoproto.VideoRec, error) {
	// Videos which have been viewed and not rated are implicitly rated 0
	// left join from video scores returns some random videos by default
	sql := "WITH tag_ratings AS (select tag, coalesce(avg(ratings.rating), 0.00) - 2.5 AS tag_score from videos INNER JOIN tags ON videos.id = tags.video_id LEFT JOIN ratings ON ratings.video_id = videos.id WHERE (ratings.user_id = $1 OR ratings.user_id is null) AND videos.views > 0 GROUP BY tag), " +
		"video_scores AS (SELECT videos.id, coalesce(sum(tag_score), 0.00) AS video_score from videos LEFT JOIN tags ON tags.video_id = videos.id LEFT JOIN tag_ratings ON tag_ratings.tag = tags.tag WHERE videos.transcoded = true GROUP BY videos.id) " +
		"SELECT videos.id, title, newLink from video_scores INNER JOIN videos ON video_scores.id = videos.id WHERE videos.is_deleted = false AND videos.id NOT IN (SELECT video_id FROM ratings WHERE ratings.user_id = $1) ORDER BY video_score, videos.views DESC limit 10;"
	rows, err := b.db.Query(sql, uid)
	if err != nil {
		return nil, err
	}

	var ret []*videoproto.VideoRec
	for rows.Next() {
		vid := videoproto.VideoRec{}
		err = rows.Scan(&vid.VideoID, &vid.VideoTitle, &vid.ThumbnailLoc)
		if err != nil {
			return nil, err
		}

		// I should stop doing this...
		vid.ThumbnailLoc = strings.Replace(vid.ThumbnailLoc, ".mpd", ".thumb", 1)

		ret = append(ret, &vid)
	}

	return ret, nil
}
