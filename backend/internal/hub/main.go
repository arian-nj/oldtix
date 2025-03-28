package hub

type Hub interface {
	Init()
	OnPlayerJoin()
	OnPlayerLeave()
	MatchLoop()
}
