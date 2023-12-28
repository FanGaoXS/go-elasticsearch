package crawler

type BoardType int

const (
	BoardTypeInvalid BoardType = iota
	BoardTypeRealtime
	BoardTypeNovel
	BoardTypeMovie
	BoardTypeTeleplay
	BoardTypeCar
	BoardTypeGame
)

var toBoardTypeString = map[BoardType]string{
	BoardTypeInvalid:  "invalid",
	BoardTypeRealtime: "realtime",
	BoardTypeNovel:    "novel",
	BoardTypeMovie:    "movie",
	BoardTypeTeleplay: "teleplay",
	BoardTypeCar:      "car",
	BoardTypeGame:     "game",
}

func ToBoardTypeString(bt BoardType) string {
	return toBoardTypeString[bt]
}

type BoardFrom int

const (
	BoardFromInvalid BoardFrom = iota
	BoardFromDouBan
	BoardFromBaidu
)

var toBoardFromString = map[BoardFrom]string{
	BoardFromInvalid: "invalid",
	BoardFromDouBan:  "douban",
	BoardFromBaidu:   "baidu",
}

func ToBoardFromString(bf BoardFrom) string {
	return toBoardFromString[bf]
}

type Board struct {
	Title string `json:"title"`
	Desc  string `json:"desc"`
	Index string `json:"index"`
	Link  string `json:"link"`

	Type BoardType `json:"type"`
	From BoardFrom `json:"from"`
}
