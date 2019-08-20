package players

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
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

func (r MysqlRepository) GetDefaultPlayerRank(ctx context.Context, wg *sync.WaitGroup, draftID int, resultsChan chan<- Results) (error) {
	defer wg.Done()

	db := r.connector.Connect()
	defer db.Close()
	q := fmt.Sprintf("SELECT p.id, COALESCE(p.nfl_rank, %d) FROM fdr_players p JOIN fdr_leagues l ON p.league_id = l.id JOIN  fdr_drafts d ON l.id = d.league_id AND d.id = ? ORDER BY nfl_rank", RankNotAvailable)
	stmt, err := db.Prepare(q)
	if err != nil {
		return err
	}
	results, err := stmt.Query(draftID)
	if err != nil {
		return err
	}

	rankNotAvailable := make([]Results, 0)
	for results.Next() {
		var res Results

		err = results.Scan(&res.PlayerID, &res.PlayerIndex)
		if err != nil {
			return err
		}
		if res.PlayerIndex == PlayerIndex(RankNotAvailable) {
			rankNotAvailable = append(rankNotAvailable, res)
			continue
		}
		resultsChan <- res
	}

	for idx := range rankNotAvailable {
		resultsChan <- rankNotAvailable[idx]
	}

	return nil
}

func (r MysqlRepository) GetUserPlayerRank(ctx context.Context, wg *sync.WaitGroup, draftID int, userID int, resultsChan chan<- Results) error {
	defer wg.Done()

	db := r.connector.Connect()
	defer db.Close()
	stmt, err := db.Prepare("SELECT fupp.player_id, fdr.id from fdr_user_player_preferences fupp LEFT JOIN fdr_draft_results fdr ON fupp.player_id = fdr.player_id AND fupp.draft_id = fdr.draft_id  WHERE fupp.draft_id = ? AND fupp.user_id = ? order by preference_order asc")
	if err != nil {
		return err
	}
	results, err := stmt.Query(draftID, userID)
	if err != nil {
		return err
	}

	for results.Next() {
		var res Results
		var draftResID sql.NullInt64
		err = results.Scan(&res.PlayerID, &draftResID)
		if err != nil {
			return err
		}
		// the draft result id is not valid, add it to the channel
		// Not valid means that they have not been drafted
		if !draftResID.Valid {
			resultsChan <- res
		}
	}

	return nil
}

func (r MysqlRepository) GetDraftResults(ctx context.Context, wg *sync.WaitGroup, draftID int, resultsChan chan<- Results) error {
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


func (r MysqlRepository) RemoveFromListIfNotIn(ctx context.Context, draftID int, userID int, playerList []int) error {
	db := r.connector.Connect()

	stmt, err := db.Prepare("DELETE FROM fdr_user_player_preferences WHERE user_id = ? and draft_id = ?")
	if err != nil {
		_ =level.Error(r.log).Log("error:", err)

		return err
	}

	_, err = stmt.Query( userID, draftID)
	if err != nil {
		return err
	}
	return nil

}


func (r MysqlRepository) SaveUserPlayerPreference(ctx context.Context, wg *sync.WaitGroup, draftID int, userID int, playerID int, preferenceOrder int) error {
	defer wg.Done()
	db := r.connector.Connect()

	stmt, err := db.Prepare("SELECT id FROM fdr_user_player_preferences WHERE draft_id = ? AND user_id = ? AND player_id = ? ORDER BY id DESC LIMIT 1")
	if err != nil {
		_ =level.Error(r.log).Log("error:", err)

		return err
	}

	results, err := stmt.Query(draftID, userID, playerID)
	if err != nil {
		return err
	}
	var id int
	if results.Next() {
		err = results.Scan(&id)
		_ =level.Debug(r.log).Log("Found id:", id)
	}
	if id == 0 {
		insertStmt, err := db.Prepare("INSERT INTO fdr_user_player_preferences (draft_id, user_id, player_id, preference_order, created_at, updated_at) VALUES (?, ?, ?, ?, NOW(), NOW())")
		if err != nil {
			_ =level.Error(r.log).Log("error:", err)

			return err
		}

		_, err = insertStmt.Exec(draftID, userID, playerID, preferenceOrder)
		if err != nil {
			return err
		}

		return nil
	}

	// update if result if already exists
	updateStmt, err := db.Prepare("UPDATE fdr_user_player_preferences SET preference_order = ? where id = ?")
	if err != nil {
		return err
	}
	_, err = updateStmt.Exec(preferenceOrder, id)
	if err != nil {
		return err
	}
	return nil

}
