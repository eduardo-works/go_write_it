package image

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/eduardo-works/writeit/pkg/validator"
	webp "github.com/mayahiro/go-webp"
	"github.com/sergeymakinen/go-ico"
	"golang.org/x/image/bmp"
	"golang.org/x/image/draw"
	_ "golang.org/x/image/webp"
)

type encoder func(io.Writer, image.Image) error

var encoders = map[string]encoder{
	"jpg": func(w io.Writer, img image.Image) error {
		return jpeg.Encode(w, img, &jpeg.Options{Quality: 90})
	},
	"png": func(w io.Writer, img image.Image) error {
		return (&png.Encoder{
			CompressionLevel: png.BestCompression,
		}).Encode(w, img)
	},
	"gif": func(w io.Writer, img image.Image) error {
		return gif.Encode(w, img, &gif.Options{NumColors: 256})
	},
	"bmp": func(w io.Writer, img image.Image) error {
		return bmp.Encode(w, img)
	},
	"webp": func(w io.Writer, img image.Image) error {
		return webp.Encode(w, img, &webp.Options{
			Compression: webp.CompressionLossy,
			Quality:     90,
		})
	},
	"ico": encodeICO,
}

func Convert(source, extTarget string) (string, error) {
	source, err := validator.ValidateFile(source)
	if err != nil {
		return "", err
	}

	source, extSource := validateSource(source)
	if extSource == "" {
		return "", fmt.Errorf("the file %s doesn't have a valid image type", source)
	}

	requestedType := extTarget
	target, extTarget := validateTarget(source, extTarget)
	if extTarget == "" {
		return "", fmt.Errorf("invalid destination image type: %s", requestedType)
	}

	if extSource == extTarget {
		return "", fmt.Errorf("cannot convert to the same image type: %s", extSource)
	}

	return convertStandard(source, target, extTarget)
}

func validateSource(path string) (string, string) {
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(path), "."))

	if ext == "jpeg" {
		ext = "jpg"
	}

	if slices.Contains(
		[]string{"bmp", "gif", "ico", "jpg", "webp", "png"},
		ext,
	) {
		return path, ext
	}

	return path, ""
}

func validateTarget(path string, toType string) (string, string) {
	toType = strings.ToLower(strings.TrimPrefix(toType, "."))

	if toType == "jpeg" {
		toType = "jpg"
	}

	if _, exists := encoders[toType]; !exists {
		return "", ""
	}

	ext := filepath.Ext(path)
	path = strings.TrimSuffix(path, ext) + "." + toType

	return path, toType
}

func convertStandard(source, target, targetType string) (string, error) {
	encode, exists := encoders[targetType]
	if !exists {
		return "", nil
	}

	input, err := os.Open(source)
	if err != nil {
		return "", err
	}
	defer input.Close()

	img, _, err := image.Decode(input)
	if err != nil {
		return "", fmt.Errorf("cannot decode image: %w", err)
	}

	output, err := os.CreateTemp(filepath.Dir(target), ".writeit-image-*")
	if err != nil {
		return "", err
	}

	tempPath := output.Name()
	defer os.Remove(tempPath)

	if err := encode(output, img); err != nil {
		output.Close()
		return "", fmt.Errorf("cannot encode image: %w", err)
	}

	if err := output.Close(); err != nil {
		return "", fmt.Errorf("cannot close image: %w", err)
	}

	if err := os.Rename(tempPath, target); err != nil {
		return "", fmt.Errorf("cannot save image: %w", err)
	}

	return target, nil
}

func encodeICO(w io.Writer, img image.Image) error {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	if width == 0 || height == 0 {
		return fmt.Errorf("invalid image dimensions")
	}

	sizes := []int{16, 32, 48, 256}
	images := make([]image.Image, 0, len(sizes))

	for _, size := range sizes {
		canvas := image.NewNRGBA(image.Rect(0, 0, size, size))

		scale := min(
			float64(size)/float64(width),
			float64(size)/float64(height),
		)

		w := max(1, int(float64(width)*scale))
		h := max(1, int(float64(height)*scale))

		x := (size - w) / 2
		y := (size - h) / 2

		rect := image.Rect(x, y, x+w, y+h)

		draw.CatmullRom.Scale(canvas, rect, img, bounds, draw.Over, nil)
		images = append(images, canvas)
	}

	return ico.EncodeAll(w, images)
}
