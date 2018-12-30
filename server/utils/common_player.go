package utils

type CommonMoviePlayer struct {
	action chan MoviePlayerAction
}

type MoviePlayerActionType string

const (
	PLAY  = MoviePlayerActionType("play")
	STOP  = MoviePlayerActionType("stop")
	PAUSE = MoviePlayerActionType("pause")
)

type MoviePlayerAction struct {
	Action MoviePlayerActionType `json:"action"`
}

func NewCommonMoviePlayer() *CommonMoviePlayer {
	wp := CommonMoviePlayer{}
	wp.action = make(chan MoviePlayerAction)
	return &wp
}

func (wp *CommonMoviePlayer) Receiver(fn func(MoviePlayerAction)) {
	for v := range wp.action {
		fn(v)
	}
}

func (wp *CommonMoviePlayer) Sender(a MoviePlayerAction) {
	wp.action <- a
}
