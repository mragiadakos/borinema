package utils

type CommonMoviePlayer struct {
	action chan MoviePlayerAction
}

type MoviePlayerActionType string

const (
	PLAY                 = MoviePlayerActionType("play")
	STOP                 = MoviePlayerActionType("stop")
	PAUSE                = MoviePlayerActionType("pause")
	CURRENT_TIME         = MoviePlayerActionType("current_time")
	REQUEST_CURRENT_TIME = MoviePlayerActionType("request_current_time")
)

type MoviePlayerAction struct {
	Action    MoviePlayerActionType `json:"action"`
	Time      float64               `json:"time"`
	IsPlaying bool                  `json:"is_playing"`
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
