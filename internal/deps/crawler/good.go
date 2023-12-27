package crawler

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

func ToGoodFromString(gf GoodFrom) string {
	return toGoodFromString[gf]
}
