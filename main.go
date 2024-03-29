package main

import (
	"bytes"
	"fmt"
	"image"
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

type CharState string

const (
	CharStateAttack CharState = "Attack"
	CharStateDeath  CharState = "Death"
	CharStateHurt   CharState = "Hurt"
	CharStateIdle   CharState = "Idle"
	CharStateWalk   CharState = "Walk"
)

type Character struct {
	hp          int
	remainingHp int
	alive       bool
	atk         int
	pts         int
	pos         Position
	state       CharState
	sprites     map[CharState]Sprite
	spriteCtr   int
	ctrStart    int
	isCpu       bool
}

func NewCharacter(pos Position, isCpu bool) *Character {
	sprites := make(map[CharState]Sprite)
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

	sprites[CharStateAttack] = NewSprite(4, ebiten.NewImageFromImage(atkImg))
	sprites[CharStateDeath] = NewSprite(4, ebiten.NewImageFromImage(deathImg))
	sprites[CharStateHurt] = NewSprite(2, ebiten.NewImageFromImage(hurtImg))
	sprites[CharStateIdle] = NewSprite(4, ebiten.NewImageFromImage(idleImg))
	sprites[CharStateWalk] = NewSprite(6, ebiten.NewImageFromImage(walkImg))

	return &Character{
		hp:          100,
		remainingHp: 100,
		alive:       true,
		atk:         50,
		pts:         0,
		pos:         pos,
		sprites:     sprites,
		state:       CharStateIdle,
		isCpu:       isCpu,
	}
}

func (c *Character) ChangeState(state CharState, count int) {
	c.state = state
	spriteCtr := c.CurrentSprite().frameCount
	if state == CharStateHurt {
		spriteCtr *= 2
	}
	c.spriteCtr = spriteCtr
	c.ctrStart = count
}

func (c *Character) Idle(count int) {
	c.state = CharStateIdle
	c.ctrStart = count
}

func (c *Character) StartAttack(count int) {
	c.ChangeState(CharStateAttack, count)
}

func (c *Character) TakeDamage(count int, dmg int) {
	c.ChangeState(CharStateHurt, count)
	if c.remainingHp-dmg <= 0 {
		c.remainingHp = 0
	} else {
		c.remainingHp -= dmg
	}
	fmt.Println("Char Remaining HP: ", c.remainingHp)
}

var (
	frameHeight = 48
	frameWidth  = 48
)

func (p *Character) CurrentSprite() Sprite {
	return p.sprites[p.state]
}

func (p *Character) GetImgOpts() *ebiten.DrawImageOptions {
	op := &ebiten.DrawImageOptions{}
	if p.isCpu {
		op.GeoM.Scale(-2, 2)
		op.GeoM.Translate(p.pos.x+float64(frameWidth*2), p.pos.y)
	} else {
		op.GeoM.Scale(2, 2)
		op.GeoM.Translate(p.pos.x, p.pos.y)
	}
	return op
}

func (c *Character) DrawFrame(screen *ebiten.Image, count int) {
	// Move rectangle to new frame of sprite
	pos := ((count - c.ctrStart) / 8) % c.CurrentSprite().frameCount
	if pos == c.CurrentSprite().frameCount-1 && c.state == CharStateDeath {
		c.alive = false
	}
	if !c.alive {
		// Hold on last frame of death animation
		pos = c.CurrentSprite().frameCount - 1
	}
	sx := pos * frameWidth
	rect := image.Rect(sx, 0, sx+frameWidth, frameHeight)
	img := c.CurrentSprite().img.SubImage(rect).(*ebiten.Image)
	op := c.GetImgOpts()

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
}

func (g *Game) SetupBattle() {
	g.state = BATTLE
	g.player = NewCharacter(Position{
		x: (screenWidth / 3) - float64(frameWidth),
		y: (screenHeight / 2) - float64(frameHeight),
	}, false)
	g.cpu = NewCharacter(Position{
		x: (screenWidth / 3 * 2) - float64(frameWidth),
		y: (screenHeight / 2) - float64(frameHeight),
	}, true)
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

	// Pause
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		resetBgImg()

		if g.state == PAUSED {
			fmt.Println("STATE: paused --to-> main menu")
			temp := g.lastState
			g.lastState = g.state
			g.state = temp
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
			g.SetupBattle()
		}
	case BATTLE:
		fmt.Println("Player alive?", g.player.alive)
		fmt.Println("CPU alive?", g.cpu.alive)
		if !g.player.alive {
			resetBgImg()
			g.lastState = g.state
			fmt.Printf("STATE: %s --to-> %s (count %d)\n", g.state.ToString(), LOST.ToString(), g.count)
			g.state = LOST
		} else if !g.cpu.alive {
			resetBgImg()
			g.lastState = g.state
			fmt.Printf("STATE: %s --to-> %s (count %d)\n", g.state.ToString(), WON.ToString(), g.count)
			g.state = WON
		}

		if inpututil.IsKeyJustPressed(ebiten.KeySpace) && g.player.state != CharStateAttack {
			fmt.Println("ATK start")
			g.player.StartAttack(g.count)
		} else if g.player.state == CharStateAttack && g.player.spriteCtr == 0 {
			if g.cpu.state != CharStateHurt && g.cpu.state != CharStateDeath {
				g.cpu.TakeDamage(g.count, g.player.atk)
				resetBgImg()
			}
		} else if g.cpu.state == CharStateHurt && g.cpu.spriteCtr == 0 && g.cpu.remainingHp == 0 {
			g.cpu.ChangeState(CharStateDeath, g.count)
		}

		if g.player.state != CharStateIdle && g.player.spriteCtr == 0 {
			fmt.Println("Action finish")
			g.player.Idle(g.count)
		}
		if g.cpu.state != CharStateIdle && g.cpu.spriteCtr == 0 {
			fmt.Println("CPU:Action finish")
			g.cpu.Idle(g.count)
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
		heading, subtext := "- Battle! -", "Press SPACE to attack!"
		smallFont := getSmallFont()
		textRect, _ := font.BoundString(smallFont, heading)
		w, h := (textRect.Max.X - textRect.Min.X).Round(), (textRect.Max.Y - textRect.Min.Y).Round()
		x, y := (screenWidth/2)-(w/2), h*2
		text.Draw(bgImg, heading, smallFont, x, y, colors.Red)

		// Subtext
		textRect, _ = font.BoundString(smallFont, subtext)
		w = (textRect.Max.X - textRect.Min.X).Round()
		x, y = (screenWidth/2)-(w/2), h*4
		text.Draw(bgImg, subtext, smallFont, x, y, colors.Black)

		// Health Bars
		// Player HP
		playerHpStr := fmt.Sprintf("HP %d/%d", g.player.remainingHp, g.player.hp)
		textRect, _ = font.BoundString(smallFont, playerHpStr)
		w = (textRect.Max.X - textRect.Min.X).Round()
		x, y = (screenWidth/3)-(w/2), h*6
		text.Draw(bgImg, playerHpStr, smallFont, x, y, colors.Grey)
		// CPU HP
		cpuHpStr := fmt.Sprintf("HP %d/%d", g.cpu.remainingHp, g.cpu.hp)
		textRect, _ = font.BoundString(smallFont, cpuHpStr)
		w = (textRect.Max.X - textRect.Min.X).Round()
		x, y = ((screenWidth/3)*2)-(w/2), h*6
		text.Draw(bgImg, cpuHpStr, smallFont, x, y, colors.Grey)

		// Sprites
		g.player.DrawFrame(screen, g.count)
		switch g.player.state {
		case CharStateAttack:
			if (g.count-g.player.ctrStart)%8 == 0 {
				g.player.spriteCtr -= 1
			}
			break
		}

		g.cpu.DrawFrame(screen, g.count)
		switch g.cpu.state {
		case CharStateHurt:
			if (g.count-g.cpu.ctrStart)%8 == 0 {
				fmt.Printf("CPU: ct DECREMENT %d -> %d\n", g.cpu.spriteCtr, g.cpu.spriteCtr-1)
				g.cpu.spriteCtr -= 1
			}
			break
		}
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
	White     color.Color
	Grey      color.Color
}{
	Primary:   color.RGBA{80, 220, 20, 255},
	Secondary: color.RGBA{60, 120, 20, 255},
	Green:     color.RGBA{0, 255, 0, 255},
	Red:       color.RGBA{255, 0, 0, 255},
	Black:     color.Black,
	White:     color.White,
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
