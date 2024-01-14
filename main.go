package main

import (
	"bytes" // will be used for loading embedded imgs
	"fmt"
	"image" // will be used for loading embedded imgs and rects
	"image/color"
	"log"

	"github.com/STLnick/catsvsdogs/resources/cat1"
	"github.com/STLnick/catsvsdogs/resources/dog1"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
)

func panicIfErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func getGofontFace(size float64) font.Face {
	gofont, err := opentype.Parse(goregular.TTF)
	panicIfErr(err)
	gofontFace, err := opentype.NewFace(gofont, &opentype.FaceOptions{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingVertical,
	})
	panicIfErr(err)

	return gofontFace
}

func getBaseFont() font.Face {
	return getGofontFace(20)
}

func getSmallFont() font.Face {
	return getGofontFace(12)
}

func getTitleFont() font.Face {
	return getGofontFace(32)
}

func resetBgImg() {
	var err error
	bgImg.Clear()
	bgImg, _, err = ebitenutil.NewImageFromFile("resources/bg1.png")
	panicIfErr(err)
}

type Position struct {
	x float64
	y float64
}

type Sprite struct {
	frameCount int
	img        *ebiten.Image
}

func NewSprite(frameCount int, img *ebiten.Image) Sprite {
	return Sprite{
		frameCount,
		img,
	}
}

type Character struct {
	hp            int
	pts           int
	pos           Position
	sprites       map[string]Sprite
	currentSprite Sprite
}

func NewCharacter(pos Position, isCpu bool) *Character {
	sprites := make(map[string]Sprite)
	var (
		atkImg   image.Image
		deathImg image.Image
		hurtImg  image.Image
		idleImg  image.Image
		walkImg  image.Image
		err      error
	)

	if !isCpu {
		atkImg, _, err = image.Decode(bytes.NewReader(cat1.Attack_png))
		panicIfErr(err)
		deathImg, _, err = image.Decode(bytes.NewReader(cat1.Death_png))
		panicIfErr(err)
		hurtImg, _, err = image.Decode(bytes.NewReader(cat1.Hurt_png))
		panicIfErr(err)
		idleImg, _, err = image.Decode(bytes.NewReader(cat1.Idle_png))
		panicIfErr(err)
		walkImg, _, err = image.Decode(bytes.NewReader(cat1.Walk_png))
		panicIfErr(err)
	} else {
		atkImg, _, err = image.Decode(bytes.NewReader(dog1.Attack_png))
		panicIfErr(err)
		deathImg, _, err = image.Decode(bytes.NewReader(dog1.Death_png))
		panicIfErr(err)
		hurtImg, _, err = image.Decode(bytes.NewReader(dog1.Hurt_png))
		panicIfErr(err)
		idleImg, _, err = image.Decode(bytes.NewReader(dog1.Idle_png))
		panicIfErr(err)
		walkImg, _, err = image.Decode(bytes.NewReader(dog1.Walk_png))
		panicIfErr(err)
	}

	sprites["Attack"] = NewSprite(4, ebiten.NewImageFromImage(atkImg))
	sprites["Death"] = NewSprite(4, ebiten.NewImageFromImage(deathImg))
	sprites["Hurt"] = NewSprite(2, ebiten.NewImageFromImage(hurtImg))
	sprites["Idle"] = NewSprite(4, ebiten.NewImageFromImage(idleImg))
	sprites["Walk"] = NewSprite(6, ebiten.NewImageFromImage(walkImg))

	return &Character{
		hp:            100,
		pts:           0,
		pos:           pos,
		sprites:       sprites,
		currentSprite: sprites[testState],
	}
}

var (
	frameHeight = 48
	frameWidth  = 48
	testCounter = 0
	testState   = "Idle"
)

func (p *Character) GetImgOpts() *ebiten.DrawImageOptions {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(p.pos.x, p.pos.y)
	return op
}

func (p *Character) DrawFrame(screen *ebiten.Image, count int) {
	// Move rectangle to new frame of sprite
	pos := (count / 8) % p.currentSprite.frameCount
	sx := pos * frameWidth
	rect := image.Rect(sx, 0, sx+frameWidth, frameHeight)
	img := p.currentSprite.img.SubImage(rect).(*ebiten.Image)

	// Adjust position of image
	// op := &ebiten.DrawImageOptions{}
	// op.GeoM.Translate(p.pos.x, p.pos.y)
	op := p.GetImgOpts()

	screen.DrawImage(img, op)
}

const (
	screenWidth  = 320
	screenHeight = 240
)

type GameState int

const (
	MAIN_MENU GameState = iota
	PAUSED
	BATTLE
	WON
	LOST
)

func (gs GameState) ToString() string {
	switch gs {
	case MAIN_MENU:
		return "MAIN_MENU"
	case PAUSED:
		return "PAUSED"
	case BATTLE:
		return "BATTLE"
	case WON:
		return "WON"
	case LOST:
		return "LOST"
	default:
		panic("invalid game state")
	}
}

type Game struct {
	count       int
	initialized bool
	lastState   GameState
	state       GameState
	player      *Character
	cpu         *Character
}

func (g *Game) Init() {
	g.initialized = true
	g.state = MAIN_MENU
	g.player = NewCharacter(Position{x: 75, y: 100}, false)
	g.cpu = NewCharacter(Position{x: 150, y: 100}, true)
}

func (g *Game) Update() error {
	if !g.initialized {
		g.Init()
	}
	g.count++

	// Quit program
	if ebiten.IsKeyPressed(ebiten.KeyMeta) && ebiten.IsKeyPressed(ebiten.KeyW) {
		fmt.Println("Key: META && W ::: TERMINATING...")
		return ebiten.Termination
	}

	//resetBgImg()

	// Pause
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		resetBgImg()

		if g.state == PAUSED {
			fmt.Println("STATE: paused --to-> main menu")
			g.lastState = g.state
			g.state = MAIN_MENU
		} else {
			fmt.Println("STATE: main menu --to-> paused")
			g.lastState = g.state
			g.state = PAUSED
		}
	}

	switch g.state {
	case MAIN_MENU:
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			resetBgImg()
			g.lastState = g.state
			fmt.Printf("STATE: %s --to-> %s (count %d)\n", g.state.ToString(), BATTLE.ToString(), g.count)
			g.state = BATTLE
		}
	case BATTLE:
		// !!! TESTING
		if g.count > 240 {
			resetBgImg()
			g.lastState = g.state
			fmt.Printf("STATE: (!testing!) %s --to-> %s (count %d)\n", g.state.ToString(), WON.ToString(), g.count)
			g.count = 0
			//g.state = WON
			g.state = LOST
		}

		if g.player.hp <= 0 {
			resetBgImg()
			g.lastState = g.state
			fmt.Printf("STATE: %s --to-> %s (count %d)\n", g.state.ToString(), LOST.ToString(), g.count)
			g.state = LOST
		} else if g.cpu.hp <= 0 {
			resetBgImg()
			g.lastState = g.state
			fmt.Printf("STATE: %s --to-> %s (count %d)\n", g.state.ToString(), WON.ToString(), g.count)
			g.state = WON
		}
	case WON:
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			resetBgImg()
			g.lastState = g.state
			fmt.Printf("STATE: %s --to-> %s (count %d)\n", g.state.ToString(), MAIN_MENU.ToString(), g.count)
			g.state = MAIN_MENU
		}
	case LOST:
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			resetBgImg()
			g.lastState = g.state
			fmt.Printf("STATE: %s --to-> %s (count %d)\n", g.state.ToString(), MAIN_MENU.ToString(), g.count)
			g.state = MAIN_MENU
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw background
	screen.DrawImage(bgImg, nil)

	switch g.state {
	case MAIN_MENU:
		titleFont := getTitleFont()
		textRect, _ := font.BoundString(titleFont, title)
		w, h := (textRect.Max.X - textRect.Min.X).Round(), (textRect.Max.Y - textRect.Min.Y).Round()
		x, y := (screenWidth/2)-(w/2), (screenHeight/4)+(h/2)
		text.Draw(bgImg, title, titleFont, x, y, colors.Primary)
		// Subtext
		subtext := "* Press ENTER to start a Battle! *"
		smallFont := getSmallFont()
		textRect, _ = font.BoundString(smallFont, subtext)
		w = (textRect.Max.X - textRect.Min.X).Round()
		x, y = (screenWidth/2)-(w/2), (screenHeight/3)+h
		text.Draw(bgImg, subtext, smallFont, x, y, colors.Secondary)
	case PAUSED:
		heading, subtext := "(PAUSED)", "press Cmd+W to exit or ESC to unpause"
		baseFont := getBaseFont()
		textRect, _ := font.BoundString(baseFont, heading)
		w, h := (textRect.Max.X - textRect.Min.X).Round(), (textRect.Max.Y - textRect.Min.Y).Round()
		x, y := (screenWidth/2)-(w/2), (screenHeight/4)+(h/2)
		text.Draw(bgImg, heading, baseFont, x, y, colors.Primary)
		// Subtext
		smallFont := getSmallFont()
		textRect, _ = font.BoundString(smallFont, subtext)
		w, h = (textRect.Max.X - textRect.Min.X).Round(), (textRect.Max.Y - textRect.Min.Y).Round()
		x, y = (screenWidth/2)-(w/2), (screenHeight/3)+(h/2)
		text.Draw(bgImg, subtext, smallFont, x, y, colors.Secondary)
	case BATTLE:
		heading := "- Battle! -"
		smallFont := getSmallFont()
		textRect, _ := font.BoundString(smallFont, heading)
		w, h := (textRect.Max.X - textRect.Min.X).Round(), (textRect.Max.Y - textRect.Min.Y).Round()
		x, y := (screenWidth/2)-(w/2), h*2
		text.Draw(bgImg, heading, smallFont, x, y, colors.Red)

		// draw sprites
		g.player.DrawFrame(screen, g.count)
		g.cpu.DrawFrame(screen, g.count)
		// maybe text saying attack command?
	case WON:
		heading, subtext := "VICTORY!", "Press enter to return to main menu"
		baseFont := getBaseFont()
		textRect, _ := font.BoundString(baseFont, heading)
		w, h := (textRect.Max.X - textRect.Min.X).Round(), (textRect.Max.Y - textRect.Min.Y).Round()
		x, y := (screenWidth/2)-(w/2), (screenHeight/4)+(h/2)
		text.Draw(bgImg, heading, baseFont, x, y, colors.Green)
		// Subtext
		smallFont := getSmallFont()
		textRect, _ = font.BoundString(smallFont, subtext)
		w, h = (textRect.Max.X - textRect.Min.X).Round(), (textRect.Max.Y - textRect.Min.Y).Round()
		x, y = (screenWidth/2)-(w/2), (screenHeight/3)+(h/2)
		text.Draw(bgImg, subtext, smallFont, x, y, colors.Secondary)
	case LOST:
		heading, subtext := "You are dead", "Press enter to return to main menu"
		baseFont := getBaseFont()
		textRect, _ := font.BoundString(baseFont, heading)
		w, h := (textRect.Max.X - textRect.Min.X).Round(), (textRect.Max.Y - textRect.Min.Y).Round()
		x, y := (screenWidth/2)-(w/2), (screenHeight/4)+(h/2)
		text.Draw(bgImg, heading, baseFont, x, y, colors.Black)
		// Subtext
		smallFont := getSmallFont()
		textRect, _ = font.BoundString(smallFont, subtext)
		w, h = (textRect.Max.X - textRect.Min.X).Round(), (textRect.Max.Y - textRect.Min.Y).Round()
		x, y = (screenWidth/2)-(w/2), (screenHeight/3)+(h/2)
		text.Draw(bgImg, subtext, smallFont, x, y, colors.Secondary)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

var (
	bgImg *ebiten.Image
	title = "Cats Vs. Dogs"
)

var colors = struct {
	Primary   color.Color
	Secondary color.Color
	Green     color.Color
	Red       color.Color
	Black     color.Color
	Grey      color.Color
}{
	Primary:   color.RGBA{80, 220, 20, 255},
	Secondary: color.RGBA{60, 120, 20, 255},
	Green:     color.RGBA{0, 255, 0, 255},
	Red:       color.RGBA{255, 0, 0, 255},
	Black:     color.Black,
	Grey:      color.RGBA{20, 20, 20, 255},
}

func main() {
	fmt.Println("---- start ----")
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle(title)
	game := &Game{}

	var err error
	bgImg, _, err = ebitenutil.NewImageFromFile("resources/bg1.png")
	panicIfErr(err)
	err = ebiten.RunGame(game)
	panicIfErr(err)

	fmt.Println("---- end ----")
}
