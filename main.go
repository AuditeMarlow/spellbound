package main

import (
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/disintegration/imaging"

	"image"
	"os"
	"os/exec"
)

var (
	blurSigma      = 5.0
	spellboundPath = ".spellbound"
	lockFile       = "lock.png"
)

func main() {
	home := os.Getenv("HOME")
	backgroundPath := home + "/" + spellboundPath + "/" + lockFile

	screenshot, err := captureScreen()
	if err != nil {
		panic(err)
	}

	blurredImage := imaging.Blur(screenshot, blurSigma)
	err = imaging.Save(blurredImage, backgroundPath)
	if err != nil {
		panic(err)
	}

	startLock(backgroundPath)
}

func startLock(path string) {
	err := exec.Command("i3lock", "-n", "-i"+path).Run()
	if err != nil {
		panic(err)
	}
}

func screenRect() (image.Rectangle, error) {
	c, err := xgb.NewConn()
	if err != nil {
		return image.Rectangle{}, err
	}
	defer c.Close()

	screen := xproto.Setup(c).DefaultScreen(c)
	x := screen.WidthInPixels
	y := screen.HeightInPixels

	return image.Rect(0, 0, int(x), int(y)), nil
}

func captureScreen() (*image.RGBA, error) {
	r, e := screenRect()
	if e != nil {
		return nil, e
	}
	return captureRect(r)
}

func captureRect(rect image.Rectangle) (*image.RGBA, error) {
	c, err := xgb.NewConn()
	if err != nil {
		return nil, err
	}
	defer c.Close()

	screen := xproto.Setup(c).DefaultScreen(c)
	x, y := rect.Dx(), rect.Dy()
	xImg, err := xproto.GetImage(c, xproto.ImageFormatZPixmap, xproto.Drawable(screen.Root), int16(rect.Min.X), int16(rect.Min.Y), uint16(x), uint16(y), 0xffffffff).Reply()
	if err != nil {
		return nil, err
	}

	data := xImg.Data
	for i := 0; i < len(data); i += 4 {
		data[i], data[i+2], data[i+3] = data[i+2], data[i], 255
	}

	img := &image.RGBA{data, 4 * x, image.Rect(0, 0, x, y)}
	return img, nil
}
