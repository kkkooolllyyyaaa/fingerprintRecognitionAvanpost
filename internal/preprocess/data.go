package preprocess

type Data struct {
	Filename  string
	KeyPoints map[int]int
}

func NewData(keyPoints map[int]int, filename string) *Data {
	return &Data{
		KeyPoints: keyPoints,
		Filename:  filename,
	}
}
