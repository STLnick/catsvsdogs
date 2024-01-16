package characters

import (
	"bytes"
	"fmt"
	"image"

	"github.com/STLnick/catsvsdogs/characters/cat1"
	"github.com/STLnick/catsvsdogs/characters/cat2"
	"github.com/STLnick/catsvsdogs/characters/dog1"
	"github.com/STLnick/catsvsdogs/characters/dog2"

	"github.com/STLnick/catsvsdogs/globals"

	"github.com/STLnick/catsvsdogs/utils"
	"github.com/hajimehoshi/ebiten/v2"
)

type Position struct {
	X float64
	Y float64
}

type Sprite struct {
	FrameCount int
	Img        *ebiten.Image
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
	Hp          int
	RemainingHp int
	Alive       bool
	Atk         int
	Pts         int
	Pos         Position
	State       CharState
	Sprites     map[CharState]Sprite
	SpriteCtr   int
	CtrStart    int
	IsCpu       bool
}

func getImages(character string) map[CharState]*ebiten.Image {
	var (
		atkImg   image.Image
		deathImg image.Image
		hurtImg  image.Image
		idleImg  image.Image
		walkImg  image.Image
		err      error
	)

    switch character {
    case "cat1":
		atkImg, _, err = image.Decode(bytes.NewReader(cat1.Attack_png))
		utils.PanicIfErr(err)
		deathImg, _, err = image.Decode(bytes.NewReader(cat1.Death_png))
		utils.PanicIfErr(err)
		hurtImg, _, err = image.Decode(bytes.NewReader(cat1.Hurt_png))
		utils.PanicIfErr(err)
		idleImg, _, err = image.Decode(bytes.NewReader(cat1.Idle_png))
		utils.PanicIfErr(err)
		walkImg, _, err = image.Decode(bytes.NewReader(cat1.Walk_png))
		utils.PanicIfErr(err)
        break;
    case "cat2":
		atkImg, _, err = image.Decode(bytes.NewReader(cat2.Attack_png))
		utils.PanicIfErr(err)
		deathImg, _, err = image.Decode(bytes.NewReader(cat2.Death_png))
		utils.PanicIfErr(err)
		hurtImg, _, err = image.Decode(bytes.NewReader(cat2.Hurt_png))
		utils.PanicIfErr(err)
		idleImg, _, err = image.Decode(bytes.NewReader(cat2.Idle_png))
		utils.PanicIfErr(err)
		walkImg, _, err = image.Decode(bytes.NewReader(cat2.Walk_png))
		utils.PanicIfErr(err)
        break;
    case "dog1":
		atkImg, _, err = image.Decode(bytes.NewReader(dog1.Attack_png))
		utils.PanicIfErr(err)
		deathImg, _, err = image.Decode(bytes.NewReader(dog1.Death_png))
		utils.PanicIfErr(err)
		hurtImg, _, err = image.Decode(bytes.NewReader(dog1.Hurt_png))
		utils.PanicIfErr(err)
		idleImg, _, err = image.Decode(bytes.NewReader(dog1.Idle_png))
		utils.PanicIfErr(err)
		walkImg, _, err = image.Decode(bytes.NewReader(dog1.Walk_png))
		utils.PanicIfErr(err)
        break;
    case "dog2":
		atkImg, _, err = image.Decode(bytes.NewReader(dog2.Attack_png))
		utils.PanicIfErr(err)
		deathImg, _, err = image.Decode(bytes.NewReader(dog2.Death_png))
		utils.PanicIfErr(err)
		hurtImg, _, err = image.Decode(bytes.NewReader(dog2.Hurt_png))
		utils.PanicIfErr(err)
		idleImg, _, err = image.Decode(bytes.NewReader(dog2.Idle_png))
		utils.PanicIfErr(err)
		walkImg, _, err = image.Decode(bytes.NewReader(dog2.Walk_png))
		utils.PanicIfErr(err)
        break;
    }

    return map[CharState]*ebiten.Image{
        CharStateAttack: ebiten.NewImageFromImage(atkImg),
        CharStateDeath: ebiten.NewImageFromImage(deathImg),
        CharStateHurt: ebiten.NewImageFromImage(hurtImg),
        CharStateIdle: ebiten.NewImageFromImage(idleImg),
        CharStateWalk: ebiten.NewImageFromImage(walkImg),
    }
}

func NewCharacter(character string, pos Position, isCpu bool) *Character {
    images := getImages(character)
    sprites := make(map[CharState]Sprite)
	sprites[CharStateAttack] = NewSprite(4, images[CharStateAttack])
	sprites[CharStateDeath] = NewSprite(4, images[CharStateDeath])
	sprites[CharStateHurt] = NewSprite(2, images[CharStateHurt])
	sprites[CharStateIdle] = NewSprite(4, images[CharStateIdle])
	sprites[CharStateWalk] = NewSprite(6, images[CharStateWalk])
	
	return &Character{
		Hp:          100,
		RemainingHp: 100,
		Alive:       true,
		Atk:         50,
		Pts:         0,
		Pos:         pos,
		Sprites:     sprites,
		State:       CharStateIdle,
		IsCpu:       isCpu,
	}
}

func (c *Character) ChangeState(state CharState, count int) {
	c.State = state
	spriteCtr := c.CurrentSprite().FrameCount
	if state == CharStateHurt {
		spriteCtr *= 2
	}
	c.SpriteCtr = spriteCtr
	c.CtrStart = count
}

func (c *Character) Idle(count int) {
	c.State = CharStateIdle
	c.CtrStart = count
}

func (c *Character) StartAttack(count int) {
	c.ChangeState(CharStateAttack, count)
}

func (c *Character) TakeDamage(count int, dmg int) {
	c.ChangeState(CharStateHurt, count)
	if c.RemainingHp-dmg <= 0 {
		c.RemainingHp = 0
	} else {
		c.RemainingHp -= dmg
	}
	fmt.Println("Char Remaining HP: ", c.RemainingHp)
}

func (p *Character) CurrentSprite() Sprite {
	return p.Sprites[p.State]
}

func (p *Character) GetImgOpts() *ebiten.DrawImageOptions {
	op := &ebiten.DrawImageOptions{}
	if p.IsCpu {
		op.GeoM.Scale(-2, 2)
		op.GeoM.Translate(p.Pos.X+float64(globals.FrameWidth*2), p.Pos.Y)
	} else {
		op.GeoM.Scale(2, 2)
		op.GeoM.Translate(p.Pos.X, p.Pos.Y)
	}
	return op
}

func (c *Character) DrawFrame(screen *ebiten.Image, count int) {
	// Move rectangle to new frame of sprite
	pos := ((count - c.CtrStart) / 8) % c.CurrentSprite().FrameCount
	if pos == c.CurrentSprite().FrameCount-1 && c.State == CharStateDeath {
		c.Alive = false
	}
	if !c.Alive {
		// Hold on last frame of death animation
		pos = c.CurrentSprite().FrameCount - 1
	}
	sx := pos * globals.FrameWidth
	rect := image.Rect(sx, 0, sx+globals.FrameWidth, globals.FrameHeight)
	img := c.CurrentSprite().Img.SubImage(rect).(*ebiten.Image)
	op := c.GetImgOpts()

	screen.DrawImage(img, op)
}
