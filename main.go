// Author: Design-BAB
// Date: 12/12/2025
// Description: It is my happy garden game project. The goal is to reach 268 lines of code
// notes: start watching after 35:08:
package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	Width        = 1000
	Height       = 800
	MoveDistance = 5
	Gravity      = 1
	FPS          = 60
)

func newVector(x, y int32) rl.Vector2 {
	return rl.Vector2{X: float32(x), Y: float32(y)}
}

func getBackground(background rl.Texture2D) []rl.Vector2 {
	var tiles []rl.Vector2
	for i := range Width/background.Width + 1 {
		for j := range Height/background.Height + 1 {
			tiles = append(tiles, newVector(i*background.Width, j*background.Height))
		}
	}
	return tiles
}

type Actor struct {
	rl.Rectangle
	Texture   rl.Texture2D
	Xvel      float32
	Yvel      float32
	Direction string
	FallCount int
}

func newActor(texture rl.Texture2D, x, y float32) *Actor {
	return &Actor{
		Rectangle: rl.Rectangle{
			X:      x,
			Y:      y,
			Width:  float32(texture.Width),
			Height: float32(texture.Height),
		},
		Texture:   texture,
		Direction: "left",
	}
}

func (a *Actor) handleMove() {
	a.Xvel = 0 // Reset horizontal velocity each frame

	if rl.IsKeyDown(rl.KeyRight) {
		a.Xvel = MoveDistance
		if a.Direction != "right" {
			a.Direction = "right"
		}
	}
	if rl.IsKeyDown(rl.KeyLeft) {
		a.Xvel = -MoveDistance
		if a.Direction != "left" {
			a.Direction = "left"
		}
	}
}

func update(player *Actor) {
	player.handleMove()
	//This is for Gravity
	player.Yvel += float32(min(1.0, float64(player.FallCount)/FPS*Gravity))
	player.FallCount += 1
	// Apply velocity
	player.X += player.Xvel
	player.Y += player.Yvel
}

func draw(background rl.Texture2D, tiles []rl.Vector2, player *Actor) {
	rl.BeginDrawing()
	rl.ClearBackground(rl.RayWhite)

	// Draw background tiles
	for _, tile := range tiles {
		rl.DrawTextureV(background, tile, rl.White)
	}

	// Draw player
	rl.DrawTexture(player.Texture, int32(player.X), int32(player.Y), rl.White)

	rl.EndDrawing()
}

func main() {
	rl.InitWindow(Width, Height, "Platformer Game")
	defer rl.CloseWindow()
	rl.SetTargetFPS(FPS)

	background := rl.LoadTexture("images/Background/Yellow.png")
	defer rl.UnloadTexture(background)
	tiles := getBackground(background)

	gopherTexture := rl.LoadTexture("images/Gopher.png")
	defer rl.UnloadTexture(gopherTexture)

	player := newActor(gopherTexture, 100, 100)

	// Game loop
	for !rl.WindowShouldClose() {
		update(player)
		draw(background, tiles, player)
	}
}
