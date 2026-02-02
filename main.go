//Author: Design-BAB
//Date: 12/12/2025
//Description: It is my happy garden game project. The goal is to reach 268 lines of code
//notes: start watching after 27:43

package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	Width        = 1000
	Height       = 800
	MoveDistance = 5
)

type GameState struct {
	IsOver bool
}

func newGameState() *GameState {
	return &GameState{}
}

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
	rl.Rectangle   // This gives Actor all the fields of rl.Rectangle (X, Y, Width, Height)
	Xvel           float32
	Yvel           float32
	Direction      string
	AnimationCount int
}

func newActor(texture rl.Texture2D, x, y float32) *Actor {
	return &Actor{Texture: texture, Rectangle: rl.Rectangle{X: x, Y: y, Width: float32(texture.Width), Height: float32(texture.Height)}, Direction: "left"}
}

func move(player *Actor, dx, dy float32) {
	player.X += dx
	player.Y += dy
}

func moveLeft(player *Actor, vel float32) {
	player.Xvel = -vel
	if player.Direction != "left" {
		player.Direction = "left"
		player.AnimationCount = 0
	}
}

func moveRight(player *Actor, vel float32) {
	player.Xvel = vel
	if player.Direction != "right" {
		player.Direction = "right"
		player.AnimationCount = 0
	}
}

func update(player *Actor) {
	move(player, player.Xvel, player.Yvel)

}

// this will act simular to getInput
func handleMove(player *Actor, yourGame *GameState) {
	if yourGame.IsOver == false {
		if rl.IsKeyDown(rl.KeyRight) {
			moveRight(player, MoveDistance)
		}
		if rl.IsKeyDown(rl.KeyLeft) {
			moveLeft(player, MoveDistance)
		}
	}
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
