package main

import (
	"encoding/base64"
	"encoding/hex"
	"image"
	"image/color"
	"image/draw"
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

type VivantQR struct {
	cfg *Config
}

func (v *VivantQR) Encrpto() ([]string, error) {
	key, err := hex.DecodeString(v.cfg.SecretKey)
	if err != nil {
		return nil, err
	}
	iv, encrypted, err := Encrypt([]byte(vivant), key)
	if err != nil {
		return nil, err
	}
	encodedIV := base64.StdEncoding.EncodeToString(iv)
	encodedEncrypted := base64.StdEncoding.EncodeToString(encrypted)
	dummy := genDummy(16)
	return separate(encodedIV, encodedEncrypted, dummy, v.cfg.Order), nil
}

func (v *VivantQR) Decrypt(content []string) ([]byte, error) {
	key, err := hex.DecodeString(v.cfg.SecretKey)
	if err != nil {
		return nil, err
	}
	var correctNums = make([]string, 0, 11)
	for _, o := range strconv.Itoa(v.cfg.Order) {
		idx, _ := strconv.Atoi(string(o))
		correctNums = append(correctNums, content[idx])
	}
	iv := correctNums[len(correctNums)-1] + correctNums[len(correctNums)-2]
	ivByte, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return nil, err
	}
	notDummyEncrypted := strings.Join(correctNums[:8], "")
	nde, err := base64.StdEncoding.DecodeString(notDummyEncrypted)
	if err != nil {
		return nil, err
	}
	decrypted, err := decrypt(nde, key, ivByte)
	if err != nil {
		return nil, err
	}
	return decrypted, nil
}
func (v *VivantQR) Output(backGroundPath, savePath string, texts []string) error {
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
		Size:    10,
		DPI:     100,
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

func separate(iv string, encryptedTxt string, dummy string, order int) []string {
	ivSepa := []string{iv[len(iv)/2:], iv[:len(iv)/2]}
	var orderStr string = strconv.Itoa(order)
	num := len(encryptedTxt) / 8
	var s = make([]string, 12)
	for i, v := range orderStr[:len(orderStr)-2] {
		idx, _ := strconv.Atoi(string(v))
		startIndex := i * num
		endIndex := (i + 1) * num
		s[idx] = encryptedTxt[startIndex:endIndex]
	}
	for i, v := range orderStr[len(orderStr)-2:] {
		idx, _ := strconv.Atoi(string(v))
		s[idx] = ivSepa[i]
	}
	s[10] = dummy[len(dummy)/2:]
	s[11] = dummy[:len(dummy)/2]
	return s
}
