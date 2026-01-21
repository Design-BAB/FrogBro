//Author: Design-BAB
//Date: 12/12/2025
//Description: It is my happy garden game project. The goal is to reach 268 lines of code
//notes: start watching after 19:19

package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	Width        = 1000
	Height       = 800
	MoveDistance = 5
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
	Texture rl.Texture2D
	//this is the collision box``
	rl.Rectangle // This gives Actor all the fields of rl.Rectangle (X, Y, Width, Height)
	Speed        float32
}

func newActor(texture rl.Texture2D, x, y float32) *Actor {
	return &Actor{Texture: texture, Rectangle: rl.Rectangle{X: x, Y: y, Width: float32(texture.Width), Height: float32(texture.Height)}, Speed: MoveDistance}
}

func draw(background rl.Texture2D, tiles []rl.Vector2) {
	rl.BeginDrawing()

	rl.ClearBackground(rl.RayWhite)
	for _, tile := range tiles {
		rl.DrawTextureV(background, tile, rl.White)
	}
	rl.DrawText("Congrats! You created your first window!", 190, 200, 20, rl.LightGray)

	rl.EndDrawing()
}

func main() {
	rl.InitWindow(Width, Height, "raylib [core] example - basic window")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	background := rl.LoadTexture("images/Background/Blue.png")
	defer rl.UnloadTexture(background)
	var tiles []rl.Vector2
	tiles = getBackground(background)
	for !rl.WindowShouldClose() {
		draw(background, tiles)
	}
}
