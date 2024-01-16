package main

import (
	"fmt"
	"image/color"

	"github.com/STLnick/catsvsdogs/characters"
	"github.com/STLnick/catsvsdogs/globals"
	"github.com/STLnick/catsvsdogs/utils"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
)

func getGofontFace(size float64) font.Face {
	gofont, err := opentype.Parse(goregular.TTF)
	utils.PanicIfErr(err)
	gofontFace, err := opentype.NewFace(gofont, &opentype.FaceOptions{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingVertical,
	})
	utils.PanicIfErr(err)

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
	utils.PanicIfErr(err)
}

var ()

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
	player      *characters.Character
	cpu         *characters.Character
}

func (g *Game) Init() {
	g.initialized = true
	g.state = MAIN_MENU
}

func (g *Game) SetupBattle() {
	g.state = BATTLE
	g.player = characters.NewCharacter("cat2", characters.Position{
		X: (globals.ScreenWidth / 3) - float64(globals.FrameWidth),
		Y: (globals.ScreenHeight / 2) - float64(globals.FrameHeight),
	}, false)
	g.cpu = characters.NewCharacter("dog2", characters.Position{
		X: (globals.ScreenWidth / 3 * 2) - float64(globals.FrameWidth),
		Y: (globals.ScreenHeight / 2) - float64(globals.FrameHeight),
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
		fmt.Println("Player alive?", g.player.Alive)
		fmt.Println("CPU alive?", g.cpu.Alive)
		if !g.player.Alive {
			resetBgImg()
			g.lastState = g.state
			fmt.Printf("STATE: %s --to-> %s (count %d)\n", g.state.ToString(), LOST.ToString(), g.count)
			g.state = LOST
		} else if !g.cpu.Alive {
			resetBgImg()
			g.lastState = g.state
			fmt.Printf("STATE: %s --to-> %s (count %d)\n", g.state.ToString(), WON.ToString(), g.count)
			g.state = WON
		}

		if inpututil.IsKeyJustPressed(ebiten.KeySpace) && g.player.State != characters.CharStateAttack {
			fmt.Println("ATK start")
			g.player.StartAttack(g.count)
		} else if g.player.State == characters.CharStateAttack && g.player.SpriteCtr == 0 {
			if g.cpu.State != characters.CharStateHurt && g.cpu.State != characters.CharStateDeath {
				g.cpu.TakeDamage(g.count, g.player.Atk)
				resetBgImg()
			}
		} else if g.cpu.State == characters.CharStateHurt && g.cpu.SpriteCtr == 0 && g.cpu.RemainingHp == 0 {
			g.cpu.ChangeState(characters.CharStateDeath, g.count)
		}

		if g.player.State != characters.CharStateIdle && g.player.SpriteCtr == 0 {
			fmt.Println("Action finish")
			g.player.Idle(g.count)
		}
		if g.cpu.State != characters.CharStateIdle && g.cpu.SpriteCtr == 0 {
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
		x, y := (globals.ScreenWidth/2)-(w/2), (globals.ScreenHeight/4)+(h/2)
		text.Draw(bgImg, title, titleFont, x, y, colors.Primary)
		// Subtext
		subtext := "* Press ENTER to start a Battle! *"
		smallFont := getSmallFont()
		textRect, _ = font.BoundString(smallFont, subtext)
		w = (textRect.Max.X - textRect.Min.X).Round()
		x, y = (globals.ScreenWidth/2)-(w/2), (globals.ScreenHeight/3)+h
		text.Draw(bgImg, subtext, smallFont, x, y, colors.Secondary)
	case PAUSED:
		heading, subtext := "(PAUSED)", "press Cmd+W to exit or ESC to unpause"
		baseFont := getBaseFont()
		textRect, _ := font.BoundString(baseFont, heading)
		w, h := (textRect.Max.X - textRect.Min.X).Round(), (textRect.Max.Y - textRect.Min.Y).Round()
		x, y := (globals.ScreenWidth/2)-(w/2), (globals.ScreenHeight/4)+(h/2)
		text.Draw(bgImg, heading, baseFont, x, y, colors.Primary)
		// Subtext
		smallFont := getSmallFont()
		textRect, _ = font.BoundString(smallFont, subtext)
		w, h = (textRect.Max.X - textRect.Min.X).Round(), (textRect.Max.Y - textRect.Min.Y).Round()
		x, y = (globals.ScreenWidth/2)-(w/2), (globals.ScreenHeight/3)+(h/2)
		text.Draw(bgImg, subtext, smallFont, x, y, colors.Secondary)
	case BATTLE:
		heading, subtext := "- Battle! -", "Press SPACE to attack!"
		smallFont := getSmallFont()
		textRect, _ := font.BoundString(smallFont, heading)
		w, h := (textRect.Max.X - textRect.Min.X).Round(), (textRect.Max.Y - textRect.Min.Y).Round()
		x, y := (globals.ScreenWidth/2)-(w/2), h*2
		text.Draw(bgImg, heading, smallFont, x, y, colors.Red)

		// Subtext
		textRect, _ = font.BoundString(smallFont, subtext)
		w = (textRect.Max.X - textRect.Min.X).Round()
		x, y = (globals.ScreenWidth/2)-(w/2), h*4
		text.Draw(bgImg, subtext, smallFont, x, y, colors.Black)

		// Health Bars
		// Player HP
		playerHpStr := fmt.Sprintf("HP %d/%d", g.player.RemainingHp, g.player.Hp)
		textRect, _ = font.BoundString(smallFont, playerHpStr)
		w = (textRect.Max.X - textRect.Min.X).Round()
		x, y = (globals.ScreenWidth/3)-(w/2), h*6
		text.Draw(bgImg, playerHpStr, smallFont, x, y, colors.Grey)
		// CPU HP
		cpuHpStr := fmt.Sprintf("HP %d/%d", g.cpu.RemainingHp, g.cpu.Hp)
		textRect, _ = font.BoundString(smallFont, cpuHpStr)
		w = (textRect.Max.X - textRect.Min.X).Round()
		x, y = ((globals.ScreenWidth/3)*2)-(w/2), h*6
		text.Draw(bgImg, cpuHpStr, smallFont, x, y, colors.Grey)

		// Sprites
		g.player.DrawFrame(screen, g.count)
		switch g.player.State {
		case characters.CharStateAttack:
			if (g.count-g.player.CtrStart)%8 == 0 {
				g.player.SpriteCtr -= 1
			}
			break
		}

		g.cpu.DrawFrame(screen, g.count)
		switch g.cpu.State {
		case characters.CharStateHurt:
			if (g.count-g.cpu.CtrStart)%8 == 0 {
				fmt.Printf("CPU: ct DECREMENT %d -> %d\n", g.cpu.SpriteCtr, g.cpu.SpriteCtr-1)
				g.cpu.SpriteCtr -= 1
			}
			break
		}
	case WON:
		heading, subtext := "VICTORY!", "Press enter to return to main menu"
		baseFont := getBaseFont()
		textRect, _ := font.BoundString(baseFont, heading)
		w, h := (textRect.Max.X - textRect.Min.X).Round(), (textRect.Max.Y - textRect.Min.Y).Round()
		x, y := (globals.ScreenWidth/2)-(w/2), (globals.ScreenHeight/4)+(h/2)
		text.Draw(bgImg, heading, baseFont, x, y, colors.Green)
		// Subtext
		smallFont := getSmallFont()
		textRect, _ = font.BoundString(smallFont, subtext)
		w, h = (textRect.Max.X - textRect.Min.X).Round(), (textRect.Max.Y - textRect.Min.Y).Round()
		x, y = (globals.ScreenWidth/2)-(w/2), (globals.ScreenHeight/3)+(h/2)
		text.Draw(bgImg, subtext, smallFont, x, y, colors.Secondary)
	case LOST:
		heading, subtext := "You are dead", "Press enter to return to main menu"
		baseFont := getBaseFont()
		textRect, _ := font.BoundString(baseFont, heading)
		w, h := (textRect.Max.X - textRect.Min.X).Round(), (textRect.Max.Y - textRect.Min.Y).Round()
		x, y := (globals.ScreenWidth/2)-(w/2), (globals.ScreenHeight/4)+(h/2)
		text.Draw(bgImg, heading, baseFont, x, y, colors.Black)
		// Subtext
		smallFont := getSmallFont()
		textRect, _ = font.BoundString(smallFont, subtext)
		w, h = (textRect.Max.X - textRect.Min.X).Round(), (textRect.Max.Y - textRect.Min.Y).Round()
		x, y = (globals.ScreenWidth/2)-(w/2), (globals.ScreenHeight/3)+(h/2)
		text.Draw(bgImg, subtext, smallFont, x, y, colors.Secondary)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return globals.ScreenWidth, globals.ScreenHeight
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
	utils.PanicIfErr(err)
	err = ebiten.RunGame(game)
	utils.PanicIfErr(err)

	fmt.Println("---- end ----")
}
