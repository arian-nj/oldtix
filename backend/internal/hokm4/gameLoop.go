package hokm4

import (
	"context"
	"encoding/json"
	"time"

	cards "github.com/arian-nj/master-card/back/internal/card"
	"github.com/arian-nj/master-card/back/internal/socket"
	"github.com/arian-nj/master-card/back/sqldb"
)

func (game *GameState) AddHumanPlayerToGame(player *HumanPlayer, gameId int) error {
	_, err := game.Queries.InsertGamePlayer(context.Background(), sqldb.InsertGamePlayerParams{
		PlayerID: player.UserId,
		GameID:   gameId,
		Team:     int(player.TeamId),
	})
	game.Players = append(game.Players, player)
	return err
}

func (game *GameState) AddBotPlayerToGame(player *BotPlayer) {
	game.Players = append(game.Players, player)
}

func (app *ApplicationHokm4) MatchUsers() error {
	for {
		game, err := app.NewGameState()
		if err != nil {
			return err
		}

		var foundPlayers []*HumanPlayer

		for len(foundPlayers) < 2 {
			foundPl := <-app.Lobby.Queue
			if foundPl.Client.State != socket.OPEN {
				continue
			}
			foundPl.BackgroundSocketHandlers(game)
			foundPlayers = append(foundPlayers, foundPl)

			for index, p := range foundPlayers {
				if p.Client.State != socket.OPEN {
					foundPlayers = append(foundPlayers[:index], foundPlayers[index+1:]...)
				}
			}
		}

		for _, p := range foundPlayers {
			err = game.AddHumanPlayerToGame(p, game.ID)
			if err != nil {
				return err
			}
		}

		// foundPlayer.IsPlayng = true

		p3 := NewBotPlayer([]cards.Card{})
		game.AddBotPlayerToGame(p3)

		p4 := NewBotPlayer([]cards.Card{})
		game.AddBotPlayerToGame(p4)

		game.Players[0].SetTeamID(TeamOne)
		game.Players[1].SetTeamID(TeamTwo)
		game.Players[2].SetTeamID(TeamOne)
		game.Players[3].SetTeamID(TeamTwo)

		app.Lobby.Mu.Lock()
		for _, p := range game.GetHumanPlayers() {
			app.Lobby.UserGames[p.UserId] = game
		}
		app.Lobby.Mu.Unlock()

		for _, p := range game.GetHumanPlayers() {
			err := game.SendGameData(TypeMatchFound, p)
			if err != nil {
				return err
			}
		}

		app.RunGameInBackground(game)
	}
}

func (app *ApplicationHokm4) RunGameInBackground(game *GameState) {
	app.BackgroundTask(func() error {
		defer func() {
			for _, p := range game.GetHumanPlayers() {
				delete(app.Lobby.UserGames, p.UserId)
				p.Client.Close()
			}
		}()

		err := game.RunGame()
		app.Logger.Error("Game Loop Ended")
		if err != nil {
			app.ReportError(err)
		}
		return err
	})
}

func (game *GameState) RunGame() error {
	for i := range 5 { // run tricks
		err := game.RunTrick(i)
		if err != nil {
			return err
		}
		if game.TeamOneTrickScore >= SETTING_WINNING_TRICK_SCORE || game.TeamTwoTrickScore >= SETTING_WINNING_TRICK_SCORE {
			break
		}
	}
	return game.TheEnd()
}

func (game *GameState) RunTrick(trick_number int) error {
	var err error

	HakemIndex := game.DeclareHakemIndex(trick_number)

	game.CurrentTrick, err = game.NewTrick(HakemIndex)
	if err != nil {
		return err
	}
	game.Tricks = append(game.Tricks, game.CurrentTrick)
	game.CurrentTrick.TurnStarterIndex = game.CurrentTrick.HakemIndex

	for _, p := range game.GetHumanPlayers() {
		err := game.SendGameData(TypeNewTrick, p)
		if err != nil {
			return err
		}
	}

	for _, p := range game.Players {
		p.SetCards([]cards.Card{})
	}

	allCards := cards.NewAllCards()
	allCards, err = game.sendCards(5, allCards)
	if err != nil {
		return err
	}

	game.CurrentTrick.Hokm, err = game.WaitToChooseHokm()
	if err != nil {
		return err
	}

	for _, p := range game.GetHumanPlayers() { // update hokm data
		err = game.SendGameData(TypeGameData, p)
		if err != nil {
			return err
		}
	}

	// send rest of cards
	allCards, err = game.sendCards(4, allCards)
	if err != nil {
		return err
	}
	_, err = game.sendCards(4, allCards)
	if err != nil {
		return err
	}
	err = game.RunTurns()
	if err != nil {
		return err
	}

	// notify winners and end the game
	for _, p := range game.GetHumanPlayers() { // update hokm data
		err := game.SendGameData(TypeEndTrick, p)
		if err != nil {
			return err
		}
	}

	err = game.Queries.UpdateTrickScores(context.Background(), sqldb.UpdateTrickScoresParams{
		TeamOneTricksScore: game.TeamOneTrickScore,
		TeamTwoTricksScore: game.TeamTwoTrickScore,
		ID:                 game.ID,
	})
	if err != nil {
		return err
	}
	return nil
}

func (game *GameState) RunTurns() error {
	for range 13 {
		err := game.RunTurn()
		if err != nil {
			return err
		}
		if game.CurrentTrick.TeamOneTurnScore >= SETTING_WINNIG_TURN_SCORE || game.CurrentTrick.TeamTwoTurnScore >= SETTING_WINNIG_TURN_SCORE {
			if game.CurrentTrick.TeamOneTurnScore >= SETTING_WINNIG_TURN_SCORE {
				game.TeamOneTrickScore += 1
				game.CurrentTrick.WinnerTeam = TeamOne
			} else {
				game.TeamTwoTrickScore += 1
				game.CurrentTrick.WinnerTeam = TeamTwo
			}
			break
		}
	}
	return nil
}

func (game *GameState) RunTurn() error {
	// new Turn
	game.CurrentTrick.CurrentTurn = NewTurn()

	// game starts
	for _, p := range game.GetHumanPlayers() {
		err := game.SendGameData(TypeTurnStart, p)
		if err != nil {
			game.Logger.Debug(err.Error())
		}
	}

	// actual game
	to_play_order := []PlayerInterface{}
	before_ward := []PlayerInterface{}
	after_ward := []PlayerInterface{}
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

	for _, playing_player := range to_play_order {
		cardIndex, err := game.WaitForPlayerToPlayCard(playing_player)
		if err != nil {
			return err
		}
		currentTurn := game.CurrentTrick.CurrentTurn
		new_card_player := NewPlayerCardPlayed(playing_player, playing_player.GetCards()[cardIndex])

		// Brodcast played card
		b_data, err := json.Marshal(new_card_player)
		if err != nil {
			return err
		}

		turn_played_event := socket.NewEvent(TypeTurnPlayed, socket.EventMessage(b_data))
		for _, player := range game.Players {
			if player != playing_player {
				player.AddToEgress(turn_played_event)
			}
		}

		currentTurn.CardsPlayed = append(currentTurn.CardsPlayed, new_card_player)
		plcards := playing_player.GetCards()
		playing_player.SetCards(append(plcards[:cardIndex], plcards[cardIndex+1:]...))
	}

	// Decide who wins Turn
	Winner := game.WhoWins()
	if Winner.Player.GetTeamID() == TeamOne {
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
		return err
	}
	_, err = game.Queries.InsertTurn(context.Background(), sqldb.InsertTurnParams{
		Moves:   string(data_byte),
		TrickID: game.CurrentTrick.id,
	})
	if err != nil {
		return err
	}
	err = game.Queries.UpdateTurnScores(context.Background(), sqldb.UpdateTurnScoresParams{
		TeamOneTurnScore: game.CurrentTrick.TeamOneTurnScore,
		TeamTwoTurnScore: game.CurrentTrick.TeamTwoTurnScore,
		TrickID:          game.CurrentTrick.id,
	})
	return err

}

func (game *GameState) TheEnd() error {
	for _, p := range game.Players {
		p.AddToEgress(socket.NewEvent(TypeTheEnd, socket.EventMessage("")))
	}
	// Statics
	var winner_team Team
	if game.TeamOneTrickScore > game.TeamTwoTrickScore {
		winner_team = TeamOne
	} else {
		winner_team = TeamTwo
	}

	TeamOneTurnScores := 0
	TeamTwoTurnScores := 0

	for _, trick := range game.Tricks {
		TeamOneTurnScores += trick.TeamOneTurnScore
		TeamTwoTurnScores += trick.TeamTwoTurnScore
	}

	for _, humanPlayer := range game.GetHumanPlayers() {

		insertStatisticsParams := sqldb.InsertHokm4StatisticParams{
			MatchID:  game.ID,
			PersonID: humanPlayer.UserId,
		}
		if humanPlayer.GetTeamID() == TeamOne {
			insertStatisticsParams.TricksWon = game.TeamOneTrickScore
			insertStatisticsParams.TricksLost = game.TeamTwoTrickScore
			insertStatisticsParams.TurnsWon = TeamOneTurnScores
			insertStatisticsParams.TurnsLost = TeamTwoTurnScores
		} else {
			insertStatisticsParams.TurnsWon = TeamTwoTurnScores
			insertStatisticsParams.TurnsLost = TeamOneTurnScores
			insertStatisticsParams.TricksWon = game.TeamTwoTrickScore
			insertStatisticsParams.TricksLost = game.TeamOneTrickScore
		}

		insertStatisticsParams.IsWon = false
		if winner_team == humanPlayer.GetTeamID() {
			insertStatisticsParams.IsWon = true
		}

		err := game.Queries.InsertHokm4Statistic(context.Background(), insertStatisticsParams)
		if err != nil {
			return err
		}
		win := 0
		loss := 0
		if insertStatisticsParams.IsWon {
			win += 1
		} else {
			loss += 1
		}
		err = game.Queries.UpdateUserStatistics(context.Background(), sqldb.UpdateUserStatisticsParams{
			Wins:            win,
			Losses:          loss,
			TotalTricksWon:  insertStatisticsParams.TricksWon,
			TotalTricksLost: insertStatisticsParams.TricksLost,
			TotalTurnsWon:   insertStatisticsParams.TurnsWon,
			TotalTurnsLost:  insertStatisticsParams.TurnsLost,
			UserID:          humanPlayer.UserId,
		})
		if err != nil {
			return err
		}
	}

	time.Sleep(5 * time.Second)
	return nil
}
