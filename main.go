package main

import (
	// "bytes" - will be used for loading embedded imgs
	"fmt"
	//  "image" - will be used for loading embedded imgs and rects
	// "image/color"
	"log"

    // TODO: Replace with correct URL
	// "example.com/ebiten-playground/resources/sprites/cat1"
	// "example.com/ebiten-playground/resources/sprites/dog1"
    "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
)

func panicIfErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func getGofontFace() font.Face {
	gofont, err := opentype.Parse(goregular.TTF)
	panicIfErr(err)
    gofontFace, err := opentype.NewFace(gofont, &opentype.FaceOptions{
		Size:    36,
		DPI:     72,
		Hinting: font.HintingVertical,
	})
	panicIfErr(err)

    return gofontFace
}

const (
	screenWidth  = 320
	screenHeight = 240
)

type Game struct {
	count       int
	initialized bool
}

func (g *Game) Init() {
	g.initialized = true
    gofontFace := getGofontFace()
    textRect, adv := font.BoundString(gofontFace, title)
    h := int(textRect.Max.Y - textRect.Min.Y)
    w := int(textRect.Max.X - textRect.Min.X)
    /**LOG*/ fmt.Println("str width", w, "-height", h, "/adv", adv)

    text.Measure()
}

func (g *Game) Update() error {
	if !g.initialized {
		g.Init()
	}
	g.count++
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Setup font
    //gofontFace := getGofontFace()
	// Background
	screen.DrawImage(bgImg, nil)

	// Text
    //strLen, adv := font.BoundString(gofontFace, title)
    // /**LOG*/ fmt.Println("strLen", strLen, "/adv", adv)
	
    //text.Draw(bgImg, title, gofontFace, (screenWidth / 2) - (strLen / 2), screenHeight/4, color.RGBA{80, 40, 220, 255})

	// Draw frames for sprites
	// g.player.DrawFrame(screen, g.count)
	// g.cpu.DrawFrame(screen, g.count)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

var (
	bgImg *ebiten.Image
    title  = "Cats Vs. Dogs"
)

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle(title)
	game := &Game{}

	var err error
	bgImg, _, err = ebitenutil.NewImageFromFile("resources/bg1.png")
	panicIfErr(err)

	err = ebiten.RunGame(game)
	panicIfErr(err)
}
