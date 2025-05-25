package hokm4

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	cards "github.com/arian-nj/master-card/back/internal/card"
	"github.com/arian-nj/master-card/back/internal/randutils"
	"github.com/arian-nj/master-card/back/internal/socket"
	"github.com/arian-nj/master-card/back/pkg/hokm4engine"
	"github.com/arian-nj/master-card/back/sqldb"
	"github.com/jackc/pgx/v5/pgtype"
)

func (game *GameState) DeclareHakemIndex(trick_number int) int {
	var HakemIndex int
	if trick_number == 0 { // if first trick
		// HakemIndex = randutils.GenerateRandomNumber(len(game.Players))
		HakemIndex = 0
	} else if game.CurrentTrick.WinnerTeam != game.Players[game.CurrentTrick.HakemIndex].GetTeamID() {
		HakemIndex = game.CurrentTrick.HakemIndex
		if HakemIndex < len(game.Players)-1 {
			HakemIndex += 1
		} else {
			HakemIndex = 0
		}
	}
	return HakemIndex
}

// have card
// same suite as first card if not first move
func (game *GameState) ValidateAndDoMove(player *HumanPlayer, played_card *cards.Card) (int, bool) {
	// game.Logger.Debug("new card " + card.String())
	cardIndex, haveCard := cards.IsCardInCards(played_card, &player.Cards)
	if !haveCard {
		game.Logger.Debug("not in hand")
		return 0, false
	}
	currentTurn := game.CurrentTrick.CurrentTurn
	if len(currentTurn.CardsPlayed) > 0 { // if not first move
		// game.Logger.Debug("first card " + currentTurn.CardsPlayed[0].Card.String())
		first_card_played_suite := currentTurn.CardsPlayed[0].Card.Suit
		if cards.HasSuit(first_card_played_suite, &player.Cards) {
			if played_card.Suit != first_card_played_suite {
				game.Logger.Debug("no allowed")
				return 0, false
			}
		}
	}

	return cardIndex, true
}

// hokm
// zamineh
// number

type GameStateOut struct {
	*GameState
	YourTeam hokm4engine.Team `json:"your_team"`
}

func (game *GameState) SendGameData(message_turn socket.EventType, p *HumanPlayer) error {
	gsOut := GameStateOut{
		GameState: game,
		YourTeam:  p.TeamId,
	}
	// send game data
	game_data, err := json.Marshal(gsOut)
	if err != nil {
		return err
	}
	p.AddToEgress(socket.NewEvent(message_turn, socket.EventMessage(game_data)), true)
	return nil
}

func (game *GameState) sendCards(number int, all_cards []cards.Card) ([]cards.Card, error) {
	var remaining_cards []cards.Card = all_cards

	for _, p := range game.Players {
		var randomCards []cards.Card
		var err error
		randomCards, remaining_cards, err = cards.GiveRandomCards(number, remaining_cards)
		if err != nil {
			game.Logger.Error(err.Error())
		}

		p.SetCards(append(p.GetCards(), randomCards...))

		humanPlayer, ok := p.(*HumanPlayer)
		if !ok {
			continue
		}

		err = game.sendCardsToEgres(randomCards, humanPlayer, TypeNewCard)
		if err != nil {
			return nil, err
		}
	}
	return remaining_cards, nil

}
func (game *GameState) sendCardsToEgres(randomCards []cards.Card, humanPlayer *HumanPlayer, eventType socket.EventType) error {
	var output struct {
		NewCards []cards.Card `json:"cards"`
	}
	output.NewCards = randomCards
	data_byte, err := json.Marshal(output)
	if err != nil {
		return err
	}
	humanPlayer.AddToEgress(socket.NewEvent(eventType, socket.EventMessage(data_byte)), true)
	return nil
}

func (game *GameState) WaitToChooseHokm() (cards.Suite, error) {

	var new_hokm cards.Suite
	hakemPlayer := game.Players[game.CurrentTrick.HakemIndex]
	// game.Logger.Info(fmt.Sprintln("hakem is ", hakem.PlayerUnique, game.CurrentTrick.HakemIndex))
	var choose_hokm_ticker *time.Ticker

	choose_hokm_ticker = time.NewTicker(SETTING_BOT_CHOOSE_HOKM_WAIT)
	humanHakemPlayer, isHakemHuman := hakemPlayer.(*HumanPlayer)
	if isHakemHuman && humanHakemPlayer.IsPlayng {
		choose_hokm_ticker = time.NewTicker(SETTING_PLAYER_CHOOSE_HOKM_WAIT)
	}

	defer choose_hokm_ticker.Stop()

Outerloop:
	for {
		select {
		case new_game_event := <-game.GameEventsCh:
			if !isHakemHuman {
				continue
			}

			if new_game_event.event.Type != TypePlayerSelectedHokm {
				continue
			}

			if new_game_event.Player.UserId != humanHakemPlayer.UserId {
				continue
			}

			hokm_data := new_game_event.event.Data
			if new_game_event.event.Data == nil {
				game.Logger.Info("no data")
				continue
			}

			hokm_int, err := strconv.Atoi(string(*hokm_data))
			if err != nil {
				game.Logger.Error(fmt.Sprintf("trying to set %s as hokm", string(*hokm_data)))
				continue
			}
			new_hokm = cards.Suite(hokm_int)
			// game.Logger.Info(fmt.Sprintf("new hokm is choosed by hakem %d ", hokm_int))
			break Outerloop
		case <-choose_hokm_ticker.C:
			rand_index := randutils.GenerateRandomNumber(4)
			new_hokm = cards.AllSuits[rand_index]
			// game.Logger.Info(fmt.Sprintf("new hokm is choosed by server %d ", int(new_hokm)))
			break Outerloop
		}
	}

	err := game.Queries.UpdateHokmTrick(context.Background(), sqldb.UpdateHokmTrickParams{
		Hokm:    pgtype.Int8{Int64: int64(game.CurrentTrick.Hokm), Valid: true},
		TrickID: game.CurrentTrick.ID,
	})

	return new_hokm, err
}

func (game *GameState) WaitForPlayerToPlayCard(playing_player hokm4engine.PlayerInterface) (cardIndex int, err error) {
	playing_player.AddToEgress(socket.NewEvent(TypeYourTurn, socket.EventMessage("")), true)

	var NewTicker *time.Ticker
	NewTicker = time.NewTicker(SETTING_BOT_PLAY_WAIT)

	PlayingHumanPlayer, ok := playing_player.(*HumanPlayer)
	if ok && PlayingHumanPlayer.IsPlayng {
		NewTicker = time.NewTicker(SETTING_PLAYER_PLAY_WAIT)
	}
	var card_played cards.Card
	cardIndex = -1

OuterLoop:
	for {
		select {
		case new_game_event := <-game.GameEventsCh:
			if new_game_event.event.Type != TypeTurnPlayed {
				// app.Logger.Debug("not same type")
				continue
			}
			if new_game_event.Player != playing_player {
				// app.Logger.Debug("not same user")
				continue
			}
			if new_game_event.event.Data == nil {
				// app.Logger.Debug("data is nil")
				continue
			}

			err := json.Unmarshal([]byte(*new_game_event.event.Data), &card_played)
			if err != nil {
				// app.Logger.Debug("can't marshal")
				game.Logger.Debug(err.Error())
				continue
			}

			newCardIndex, isValid := game.ValidateAndDoMove(new_game_event.Player, &card_played)

			// app.Logger.Debug(card_played.String())
			if !isValid {
				new_game_event.Player.AddToEgress(socket.NewEvent(TypeInvalidPlay, socket.EventMessage("")), true)
				// app.Logger.Debug("move not valid")
				continue
			}
			cardIndex = newCardIndex
			new_game_event.Player.AddToEgress(socket.NewEvent(TypeValidPlay, socket.EventMessage("")), true)
			break OuterLoop

		case <-NewTicker.C:
			NewTicker.Stop()
			cardIndex = game.BotPlayTurn(playing_player.GetCards())
			choosen_card_by_bot := playing_player.GetCards()[cardIndex]
			data_byte, err := json.Marshal(choosen_card_by_bot)
			if err != nil {
				return -1, err
			}
			playing_player.AddToEgress(socket.NewEvent(TypePlayTimeout, socket.EventMessage(data_byte)), true)
			break OuterLoop
		}
	}
	return cardIndex, nil
}
