tags : [[turn-base]] [[multiplayer]] [[game-design]]

source : https://herluf-ba.github.io/making-a-turn-based-multiplayer-game-in-rust-01-whats-a-turn-based-game-anyway.html?utm_source=chatgpt.com



the OG turn-based games: Chess

![[Pasted image 20241217020304.png]]
```rust
enum Piece {
  // White pawn, bishop, knight, rook, queen and king 
  WP, WB, WN, WR, WQ, WK,
  // Black pawn, bishop, knight, rook, queen and king 
  BP, BB, BN, BR, BQ, BK
}

struct Board(Vec<Vec<Option<Piece>>>);

let the_board_from_the_image = {
  use Piece::*;
  Board(vec![
    vec![Some(BR), None,     None,     Some(BQ), Some(BK), Some(BB), None,     Some(BR)],
    vec![Some(BP), Some(BB), Some(BP), Some(BN), None,     Some(BP), Some(BP), Some(BP)],
    vec![None,     Some(BP), None,     None,     None,     None,     None,     None   ],
    vec![None,     None,     None,     Some(BP), None,     None,     Some(WB), None   ],
    vec![None,     None,     None,     Some(WP), None,     None,     None,     None   ],
    vec![Some(WP), None,     None,     None,     None,     Some(WN), None,     None   ],
    vec![None,     Some(WP), Some(WQ), None,     Some(WP), Some(WP), Some(WP), Some(WP)],
    vec![Some(WR), None,     None,     None,     Some(WK), Some(WB), None,     Some(WR)],
  ])
};
```

How could we synchronize a game of chess between players playing online? A naive approach would be to just have a server receive moves from each player, calculate a new Board based on that move and then send this new Board to both players. That would work, and honestly, you could probably make a game with that approach as long as the size of the game's state, ie. the Board struct, isn't too large.

A more efficient, and frankly straight up cooler, solution is to just have the server send the _move the player performs_ rather than the state. Because _the game state can be determined from a sequence of actions_, the players themselves can determine the most recent board. This way the messages between the server and players become something like `GameEvent::MovePiece { piece: Piece::Rook, position: "a4" }` rather than `*entire board with all the pieces*`, which almost always is way smaller!

Some other less obvious consequences of thinking about game states as sequences of actions are:

- It's fairly straightforward to implement undo and redo functionality. Just pop or push actions onto the sequence.
- Saving and loading games can be done by simply saving and loading the sequence of actions.
- Having the actions players perform in the game, given a particular sequence of previous actions is _very useful_ if we want to train a machine learning algorithm to add interesting bots to our games.
## A pattern for online turn-based games

**A _reducer_ is a function that takes a state and an event and produces a new state**

**Events are just what we have been calling actions up until now.**


So the game will progress like this:

1. Player sends event to the server
2. The server inspects the event and verifies that it is valid and allowed
3. The server updates its own state using the reducer function
4. The server sends the valid event to all players
5. Each player uses the reducer function to update their states
6. Either someone has won or we go to 1.

![[update-state.gif]]

