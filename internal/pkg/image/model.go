package image

type Image struct {
	Thumbnail string `xorm:"thumbnail" json:"thumbnail"`
	Origin    string `xorm:"origin" json:"origin"`
}

func ParseImage(image string) *Image {
	return &Image{Origin: image}
}

func ParseImages(image []string) (results []*Image) {
	for _, i := range image {
		results = append(results, ParseImage(i))
	}
	return
}
