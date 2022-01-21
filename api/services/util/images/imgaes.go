package images

import (
	"api/services/util/log"
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"math"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type Image struct {
	Width  int
	Height int
	MediaType string
	B64Buffer []byte
}

func FormatBase64Images(base64img string) (Image, error) {
	f := Image{}
	var img, err = url.QueryUnescape(base64img)
	if err != nil {
		log.Error("QueryUnescape error", err)
		return f, err
	}
	//log.Debug("QueryUnescape", img)
	b64data := img[strings.Index(img, ",")+1:]
	var data = img[:strings.IndexByte(img, ';')]
	var mediaType = data[strings.IndexByte(data, ':')+1:]

	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(b64data))
	var config image.Config
	if mediaType == "image/png" {
		config, err = png.DecodeConfig(reader)
	} else {
		config, err = jpeg.DecodeConfig(reader)
	}
	if err != nil {
		log.Error("image DecodeConfig error", err, config)
		return f, err
	}

	b64Buffer, err := base64.StdEncoding.DecodeString(b64data)
	if err != nil {
		log.Error("Cannot decode b64", err)
		return f, err
	}

	f.Width = config.Width
	f.Height = config.Height
	f.MediaType = mediaType
	f.B64Buffer = b64Buffer
	return f, nil
}

func Resize(data Image) (image.Image, error) {
	var img image.Image
	var err error
	if data.MediaType == "image/png" {
		img, err = png.Decode(bytes.NewBuffer(data.B64Buffer))
	} else {
		img, err = jpeg.Decode(bytes.NewBuffer(data.B64Buffer))
	}
	if err != nil {
		log.Debug("image Decode Error", err)
		return nil, err
	}

	width, height := getResize(data.Width, data.Height)
	m := resize.Thumbnail(width, height, img, resize.Lanczos3)
	return m, nil
}

func CreateImageFile(data image.Image, filename string) error {
	if err := os.MkdirAll(filepath.Dir(filename), 0666); err != nil {
		log.Debug("create file path error", err)
		return err
	}

	out, err := os.Create(filename)
	if err != nil {
		log.Debug("create file error", err)
		return err
	}
	defer out.Close()

	// write new image to file
	err = jpeg.Encode(out, data, nil)
	if err != nil {
		log.Debug("image out file error", err)
		return err
	}
	return nil
}

func ResizeIcon(data image.Image) (image.Image, error) {
	m := resize.Thumbnail(100, 100, data, resize.Lanczos3)
	return m, nil
}

func TangentCircle(img image.Image) image.Image {
	srcImg, err := ResizeIcon(img)
	if err != nil {
		log.Error("Resize Image")
	}
	w := srcImg.Bounds().Max.X - srcImg.Bounds().Min.X

	h := srcImg.Bounds().Max.Y - srcImg.Bounds().Min.Y
	d := w
	if w > h {
		d = h
	}
	maskImg := circleMask(d)
	dstImg := image.NewRGBA(image.Rect(0,0,d,d))
	draw.DrawMask(dstImg, srcImg.Bounds().Add(image.Pt(0,0)), srcImg, image.Pt((w-d)/2,(h-d)/2), maskImg,image.Pt(0,0),draw.Src)
	return dstImg
}

func circleMask(d int) image.Image{
	img := image.NewRGBA(image.Rect(0,0,d,d))
	for x := 0; x < d; x++ {
		for y := 0; y < d; y++ {
			dis := math.Sqrt(math.Pow(float64(x-d/2), 2) + math.Pow(float64(y-d/2), 2))
			if dis > float64(d) / 2 {
				img.Set(x, y, color.RGBA{255, 255, 255, 0})
			} else {
				img.Set(x, y, color.RGBA{0, 0, 255, 255})
			}
		}
	}
	return img
}

func getResize(sWidth int, sHeight int) (uint, uint) {
	var width, heigth uint
	if sWidth > sHeight {
		if sWidth > 1024 {
			width = 1024
			heigth = 1024
		} else {
			width = uint(sWidth)
			heigth = uint(sWidth)
		}
	} else {
		if sHeight > 1024 {
			width = 1024
			heigth = 1024
		} else {
			width = uint(sHeight)
			heigth = uint(sHeight)
		}
	}
	return width, heigth
}

func GetImageUrl(image string) string {
	//host := middleware.GetHostname()
	return fmt.Sprintf("/static/images/product/%s", image)
}

func GetImageFromFilePath(filePath string) (image.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	image, _, err := image.Decode(f)
	return image, err
}