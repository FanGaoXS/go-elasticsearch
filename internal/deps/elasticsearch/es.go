package elasticsearch

type ES interface{}

const (
	ESListenAddr = "http://localhost:9200"
	Index        = "my_index"
)

func New() (ES, error) {
	return &esImpl{}, nil
}

type esImpl struct{}
