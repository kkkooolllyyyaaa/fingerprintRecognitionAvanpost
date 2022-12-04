package myimage

import "image"

type MyImage struct {
	Img      *image.Gray
	Filename string
}

func NewMyImage(gray *image.Gray, filename string) *MyImage {
	return &MyImage{
		Img:      gray,
		Filename: filename,
	}
}
