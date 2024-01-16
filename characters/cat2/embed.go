package cat2

import (
	_ "embed"
)

var (
	//go:embed Attack.png
	Attack_png []byte

	//go:embed Death.png
	Death_png []byte

	//go:embed Hurt.png
	Hurt_png []byte

	//go:embed Idle.png
	Idle_png []byte

	//go:embed Walk.png
	Walk_png []byte
)
