package hokm4

import (
	"context"
	"encoding/json"
	"time"

	cards "github.com/arian-nj/master-card/back/internal/card"
	"github.com/arian-nj/master-card/back/internal/socket"
	"github.com/arian-nj/master-card/back/pkg/hokm4engine"
	"github.com/arian-nj/master-card/back/sqldb"
	"github.com/jackc/pgx/v5/pgtype"
)

func (game *Room) RunGame() {
	for i := range 5 { // run tricks
		err := game.RunTrick(i)
		if err != nil {
			game.Logger.Error(err.Error())
			return
		}
		if game.TeamOneTrickScore >= SETTING_WINNING_TRICK_SCORE || game.TeamTwoTrickScore >= SETTING_WINNING_TRICK_SCORE {
			break
		}
	}
	err := game.TheEnd()
	if err != nil {
		game.Logger.Error(err.Error())
	}
}

func (game *Room) RunTrick(trick_number int) error {
	var err error

	HakemIndex := game.DeclareHakemIndex(trick_number)

	trickRow, err := game.Queries.InsertTrick(context.Background(), sqldb.InsertTrickParams{
		GameID:     game.ID,
		HakemIndex: HakemIndex,
	})
	if err != nil {
		game.Logger.Error(err.Error())
		return err
	}
	game.CurrentTrick = game.NewTrick(trickRow.TrickID, HakemIndex)
	game.Tricks = append(game.Tricks, game.CurrentTrick)
	game.CurrentTrick.TurnStarterIndex = game.CurrentTrick.HakemIndex

	for _, p := range game.GetHumanPlayers() {
		p.AllEvents = []*socket.Event{}
		err := game.SendGameData(TypeNewTrick, p)
		if err != nil {
			game.Logger.Error(err.Error())
			return err
		}
	}

	for _, p := range game.Players {
		p.SetCards([]cards.Card{})
	}

	time.Sleep(500 * time.Millisecond) //bad fix: some times TypeNewTrick get proccesed after TypeNewCard in sendCard(5,...)

	allCards := cards.NewAllCards()
	allCards, err = game.sendCards(5, allCards)
	if err != nil {
		game.Logger.Error(err.Error())
		return err
	}

	game.CurrentTrick.Hokm, err = game.WaitToChooseHokm()
	if err != nil {
		game.Logger.Error(err.Error())
		return err
	}

	for _, p := range game.GetHumanPlayers() { // update hokm data
		err = game.SendGameData(TypeNewHokm, p)
		if err != nil {
			game.Logger.Error(err.Error())
			return err
		}
	}

	// send rest of cards
	allCards, err = game.sendCards(4, allCards)
	if err != nil {
		game.Logger.Error(err.Error())
		return err
	}
	_, err = game.sendCards(4, allCards)
	if err != nil {
		game.Logger.Error(err.Error())
		return err
	}
	time.Sleep(SETTING_WAIT_BEFORE_STARTING_NEW_TURN)

	err = game.RunTurns()
	if err != nil {
		game.Logger.Error(err.Error())
		return err
	}

	// notify winners and end the trick
	for _, p := range game.GetHumanPlayers() { // update hokm data
		err := game.SendGameData(TypeEndTrick, p)
		if err != nil {
			game.Logger.Error(err.Error())
			return err
		}
	}

	err = game.Queries.UpdateTrickScores(context.Background(), sqldb.UpdateTrickScoresParams{
		TeamOneTricksScore: game.TeamOneTrickScore,
		TeamTwoTricksScore: game.TeamTwoTrickScore,
		ID:                 game.ID,
	})
	if err != nil {
		game.Logger.Error(err.Error())
	}

	err = game.SaveGameStateData()
	if err != nil {
		game.Logger.Error(err.Error())
	}
	return nil
}

func (game *Room) RunTurns() error {
	for range 13 {
		err := game.RunTurn()
		if err != nil {
			game.Logger.Error(err.Error())
			return err
		}
		if game.CurrentTrick.TeamOneTurnScore >= SETTING_WINNIG_TURN_SCORE || game.CurrentTrick.TeamTwoTurnScore >= SETTING_WINNIG_TURN_SCORE {
			if game.CurrentTrick.TeamOneTurnScore >= SETTING_WINNIG_TURN_SCORE {
				game.TeamOneTrickScore += 1
				game.CurrentTrick.WinnerTeam = hokm4engine.TeamOne
			} else {
				game.TeamTwoTrickScore += 1
				game.CurrentTrick.WinnerTeam = hokm4engine.TeamTwo
			}
			break
		}
	}
	return nil
}

func (game *Room) RunTurn() error {
	// new Turn
	game.CurrentTrick.CurrentTurn = hokm4engine.NewTurn()

	// game starts
	for _, p := range game.GetHumanPlayers() {
		err := game.SendGameData(TypeTurnStart, p)
		if err != nil {
			game.Logger.Debug(err.Error())
		}
	}

	// Actual game

	// Playing Order
	to_play_order := []hokm4engine.PlayerInterface{}
	before_ward := []hokm4engine.PlayerInterface{}
	after_ward := []hokm4engine.PlayerInterface{}
	var starterFound = false

	for player_index, p := range game.Players {
		if player_index == game.CurrentTrick.TurnStarterIndex {
			starterFound = true
			to_play_order = append(to_play_order, p)
			continue
		}

		if starterFound {
			after_ward = append(after_ward, p)

		} else {
			before_ward = append(before_ward, p)

		}
	}
	to_play_order = append(to_play_order, after_ward...)
	to_play_order = append(to_play_order, before_ward...)

	// Wait for each player to play
	for _, playing_player := range to_play_order {
		cardIndex, err := game.WaitForPlayerToPlayCard(playing_player)
		if err != nil {
			game.Logger.Error(err.Error())
			return err
		}
		currentTurn := game.CurrentTrick.CurrentTurn
		new_card_player := hokm4engine.NewPlayerCardPlayed(playing_player, playing_player.GetCards()[cardIndex])

		// Brodcast played card
		b_data, err := json.Marshal(new_card_player)
		if err != nil {
			game.Logger.Error(err.Error())
			return err
		}

		turn_played_event := socket.NewEvent(TypeTurnPlayed, socket.EventMessage(b_data))
		for _, player := range game.Players {
			if player != playing_player {
				player.AddToEgress(turn_played_event, true)
			}
		}

		currentTurn.CardsPlayed = append(currentTurn.CardsPlayed, new_card_player)
		plcards := playing_player.GetCards()
		playing_player.SetCards(append(plcards[:cardIndex], plcards[cardIndex+1:]...))
	}

	// Decide who wins Turn
	Winner := game.WhoWinsLastTurn()
	if Winner.Player.GetTeamID() == hokm4engine.TeamOne {
		game.CurrentTrick.TeamOneTurnScore += 1
	} else {
		game.CurrentTrick.TeamTwoTurnScore += 1
	}

	for playerIndex, player := range game.Players {
		if player == Winner.Player {
			game.CurrentTrick.TurnStarterIndex = playerIndex
		}
	}

	time.Sleep(SETTING_BEFORE_END_TURN_MESSAGE_SLEEP_TIME)

	// End The Turn
	for _, p := range game.GetHumanPlayers() {
		err := game.SendGameData(TypeTurnEnd, p)
		if err != nil {
			game.Logger.Debug(err.Error())
		}
	}

	data_byte, err := json.Marshal(game.CurrentTrick.CurrentTurn.CardsPlayed)
	if err != nil {
		game.Logger.Error(err.Error())
		return err
	}
	_, err = game.Queries.InsertTurn(context.Background(), sqldb.InsertTurnParams{
		Moves:   string(data_byte),
		TrickID: game.CurrentTrick.ID,
	})
	if err != nil {
		game.Logger.Error(err.Error())
		return err
	}
	err = game.Queries.UpdateTurnScores(context.Background(), sqldb.UpdateTurnScoresParams{
		TeamOneTurnScore: game.CurrentTrick.TeamOneTurnScore,
		TeamTwoTurnScore: game.CurrentTrick.TeamTwoTurnScore,
		TrickID:          game.CurrentTrick.ID,
	})
	if err != nil {
		game.Logger.Error(err.Error())

	}
	return err

}

func (game *Room) AddCoins(hplayer *HumanPlayer) error {
	coin_to_add := 0
	if hplayer.BetAmount == BET_NO_MONEY {
		coin_to_add = BET_AMOUNT_ONE
	} else if hplayer.BetAmount == BET_AMOUNT_ONE {
		coin_to_add = BET_AMOUNT_ONE_WIN
	}

	err := game.Queries.AddCoinToPerson(context.Background(), sqldb.AddCoinToPersonParams{
		Coin: coin_to_add,
		ID:   hplayer.UserId,
	})
	return err
}

func (game *Room) TheEnd() error {
	for _, p := range game.GetHumanPlayers() {
		delete(game.Lobby.UserGames, p.UserId)
	}

	for _, p := range game.Players {
		p.AddToEgress(socket.NewEvent(TypeTheEnd, socket.EventMessage("")), true)
	}

	game.Queries.UpdateHokm4Endstamp(context.Background(), sqldb.UpdateHokm4EndstampParams{
		EndStamp: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		ID:       game.ID,
	})

	// Statics
	var winner_team hokm4engine.Team
	if game.TeamOneTrickScore > game.TeamTwoTrickScore {
		winner_team = hokm4engine.TeamOne
	} else {
		winner_team = hokm4engine.TeamTwo
	}

	TeamOneTurnScores := 0
	TeamTwoTurnScores := 0

	for _, trick := range game.Tricks {
		TeamOneTurnScores += trick.TeamOneTurnScore
		TeamTwoTurnScores += trick.TeamTwoTurnScore
	}

	for _, humanPlayer := range game.GetHumanPlayers() {

		updateStaticsParams := sqldb.UpdateUserStatisticsParams{
			UserID: humanPlayer.UserId,
		}

		if winner_team == humanPlayer.GetTeamID() {
			updateStaticsParams.Win = 1
		} else {
			updateStaticsParams.Lose = 1
		}

		if updateStaticsParams.Win == 1 {
			err := game.AddCoins(humanPlayer)
			if err != nil {
				game.Logger.Error(err.Error())
				return err
			}
		}

		if humanPlayer.GetTeamID() == hokm4engine.TeamOne {
			updateStaticsParams.TotalTricksWon = game.TeamOneTrickScore
			updateStaticsParams.TotalTricksLost = game.TeamTwoTrickScore
			updateStaticsParams.TotalTurnsWon = TeamOneTurnScores
			updateStaticsParams.TotalTurnsLost = TeamTwoTurnScores
		} else {
			updateStaticsParams.TotalTurnsWon = TeamTwoTurnScores
			updateStaticsParams.TotalTurnsLost = TeamOneTurnScores
			updateStaticsParams.TotalTricksWon = game.TeamTwoTrickScore
			updateStaticsParams.TotalTricksLost = game.TeamOneTrickScore
		}

		// err := game.Queries.InsertHokm4Statistic(context.Background(), insertStatisticsParams)
		// if err != nil {
		// 				game.Logger.Error(err.Error())
		// return err
		// }

		err := game.Queries.UpdateUserStatistics(context.Background(), updateStaticsParams)
		if err != nil {
			game.Logger.Error(err.Error())
			return err
		}
	}

	time.Sleep(5 * time.Second)
	return nil
}
