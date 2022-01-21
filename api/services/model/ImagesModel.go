package model

import (
	"api/services/util/upload"
	"strings"
)

func HandleProductImages(images []string) ([]string, error) {
	var filename []string
	for _, v := range images {
		if len(v) != 0 {
			if strings.Index(v, ".") < 0 {
				name, err := upload.ProductImage(v)
				if err != nil {
					return nil, err
				}
				filename = append(filename, name)
			} else {
				filename = append(filename, v)
			}
		}
	}
	return filename, nil
}
