package players

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/thethan/fantasydraftroom/internal/mysql"
	"sync"
)

func NewMysqlRepository(connector *mysql.Connector, log log.Logger) MysqlRepository {
	return MysqlRepository{
		connector: connector,
		log:       log,
	}
}

type MysqlRepository struct {
	connector *mysql.Connector
	log       log.Logger
}


func (r MysqlRepository) GetDefaultPlayerRank(ctx context.Context, wg *sync.WaitGroup, draftID int, resultsChan chan <- Results) (error) {
	defer wg.Done()

	db := r.connector.Connect()
	defer db.Close()
	stmt, err := db.Prepare("SELECT p.id, COALESCE(p.nfl_rank, 9000) FROM fdr_players p JOIN fdr_leagues l ON p.league_id = l.id JOIN  fdr_drafts d ON l.id = d.league_id AND d.id = ?")
	if err != nil {
		return err
	}
	results, err := stmt.Query(draftID)
	if err != nil {
		return err
	}
	for results.Next() {
		var res Results

		err = results.Scan(&res.PlayerID, &res.PlayerIndex)
		if err != nil {
			return err
		}
		resultsChan <- res
	}

	return nil
}

func (r MysqlRepository) GetUserPlayerRank(ctx context.Context, wg *sync.WaitGroup, draftID int, userID int, resultsChan chan <- Results) error {
	panic("implement me... GetUserPlayerRank")
}

func (r MysqlRepository) GetDraftResults(ctx context.Context, wg *sync.WaitGroup, draftID int, resultsChan chan <- Results) error {
	defer wg.Done()

	db := r.connector.Connect()
	defer db.Close()
	stmt, err := db.Prepare("SELECT player_id, draft_order FROM fdr_draft_results WHERE draft_id = ? ORDER BY id DESC ")
	if err != nil {
		return err
	}
	results, err := stmt.Query(draftID)
	if err != nil {
		return err
	}
	for results.Next() {
		var res Results

		err = results.Scan(&res.PlayerID, &res.PlayerIndex)
		if err != nil {
			return err
		}
		resultsChan <- res
	}

	return nil
}