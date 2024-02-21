package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math/rand"
	"strings"

	"image/png"
	"os"
	"strconv"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/font/gofont/goitalic"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

const vivant = "https://www.netflix.com/jp/title/81726701"

type VivantQR struct {
	cfg *Config
}

func NewVivantQR(cfg *Config) *VivantQR { return &VivantQR{cfg: cfg} }

func (v *VivantQR) Encode() ([]string, error) {
	HMAC := hmac.New(sha256.New, []byte(v.cfg.SecretKey))
	HMAC.Write([]byte(vivant))
	sig := HMAC.Sum(nil)
	octal := bytesToOctalString(sig)
	dummy := bytesToOctalString([]byte(string(rune(genData(5)))))
	return separate(octal, dummy, v.cfg.Order), nil
}

func (v *VivantQR) Decode(content []string) (string, error) {
	HMAC := hmac.New(sha256.New, []byte(v.cfg.SecretKey))
	HMAC.Write([]byte(vivant))
	sig := HMAC.Sum(nil)

	var correctNums string
	for _, o := range strconv.Itoa(v.cfg.Order) {
		idx, _ := strconv.Atoi(string(o))
		correctNums += content[idx]
	}
	e, err := octalStringToBytes(correctNums)
	if err != nil {
		return "", err
	}
	if !hmac.Equal(e, sig) {
		return "", fmt.Errorf("faild to decode")
	}
	return vivant, nil
}

func (v *VivantQR) EncodeRawData(encrypted []string) []string {
	var res []string
	for i := 2; i <= 12; i += 2 {
		a := strings.Join(encrypted[i-2:i], " ")
		res = append(res, a)
	}
	return res
}

func (v *VivantQR) FormatRawData(content string) []string {
	res := make([]string, 0, 10)
	spliteV := strings.Split(string(content), " ")
	for _, v := range spliteV {
		v = strings.TrimSpace(v)
		if len(v) == 0 {
			continue
		}
		res = append(res, v)
	}
	return res
}

func (v *VivantQR) Output(savePath string, texts []string) error {
	backGroundPath := "./background.png"
	file, err := os.Open(backGroundPath)
	if err != nil {
		return err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	bounds := img.Bounds()
	m := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	draw.Draw(m, m.Bounds(), img, bounds.Min, draw.Src)

	textWidth := basicfont.Face7x13.Width
	lines := 6
	verticalSpacing := bounds.Dy() / (lines + 1)

	for i := 0; i <= lines-1; i += 1 {
		x := (bounds.Dx() - textWidth) / 22
		y := (i + 1) * verticalSpacing
		t := texts[i]
		if err := drawTxt(m, x, y, t); err != nil {
			return err
		}
	}
	output, err := os.Create(savePath)
	if err != nil {
		return err
	}
	defer output.Close()

	png.Encode(output, m)
	return nil
}

func drawTxt(img *image.RGBA, x, y int, text string) error {
	f, err := opentype.Parse(goitalic.TTF)
	if err != nil {
		return err
	}
	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    20,
		DPI:     120,
		Hinting: font.Hinting(font.WeightThin),
	})
	if err != nil {
		return err
	}
	point := fixed.Point26_6{X: fixed.Int26_6(x * 64), Y: fixed.Int26_6(y * 64)}
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.Black),
		Face: face,
		Dot:  point,
	}
	d.DrawString(text)
	return nil
}

func separate(octal string, dummy string, order int) []string {
	var orderStr string = strconv.Itoa(order)
	num := len(octal) / 10
	remainder := len(octal) % 10
	var s = make([]string, 12)
	for i, v := range orderStr {
		idx, _ := strconv.Atoi(string(v))
		startIndex := i * num
		endIndex := (i + 1) * num
		if i == 9 {
			endIndex += remainder
		}
		s[idx] = octal[startIndex:endIndex]
	}
	return s
}

func genData(num int64) int {
	rand.New(rand.NewSource(num))
	var res string
	for _, v := range rand.Perm(int(num)) {
		res += strconv.Itoa(v)
	}
	a, _ := strconv.Atoi(res)
	return a
}
