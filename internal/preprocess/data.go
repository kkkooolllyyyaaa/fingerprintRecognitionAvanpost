package preprocess

type Data struct {
	bin       [][]bool
	keyPoints map[int]int
}

func NewData(keyPoints map[int]int) *Data {
	return &Data{
		keyPoints: keyPoints,
	}
}
