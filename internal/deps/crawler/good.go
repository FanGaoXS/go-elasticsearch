package crawler

import (
	"bytes"

	"fangaoxs.com/go-elasticsearch/internal/infras/errors"
)

type GoodFrom int

const (
	GoodFromInvalid GoodFrom = iota
	GoodFromJinDong
	GoodFromTaoBao
	GoodFromPinDuoDuo
)

var toGoodFromString = map[GoodFrom]string{
	GoodFromInvalid:   "invalid",
	GoodFromJinDong:   "jd.com",
	GoodFromTaoBao:    "taobao.com",
	GoodFromPinDuoDuo: "pinduoduo.com",
}

var toGoodFromID = map[string]GoodFrom{
	"invalid":       GoodFromInvalid,
	"jd.com":        GoodFromJinDong,
	"taobao.com":    GoodFromTaoBao,
	"pinduoduo.com": GoodFromPinDuoDuo,
}

func ToGoodFromString(gf GoodFrom) string {
	return toGoodFromString[gf]
}

func ToGoodFromID(s string) GoodFrom {
	return toGoodFromID[s]
}

func (g GoodFrom) String() string {
	return toGoodFromString[g]
}

func (g GoodFrom) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(ToGoodFromString(g))
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (g *GoodFrom) UnmarshalJSON(data []byte) error {
	s := string(data[1 : len(data)-1]) // trim `""`
	if _, ok := toGoodFromID[s]; !ok {
		return errors.Newf(errors.InvalidArgument, nil, "unsupported good from: %s", s)
	}
	*g = toGoodFromID[s]
	return nil
}
