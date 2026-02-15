// Author: Design-BAB
// Date: 12/12/2025
// Description: It is my platform game! The goal is to reach 268 lines of code
// Notes: My next step should be to create a Block struct with rl.Rectangle and Texture fields. -done
//Write a getBlock function that loads your terrain sprite sheet and extracts a specific tile (e.g., at position 96, 0).-in progress
//Create a newBlock constructor that takes x, y, and size, then calls getBlock to set the texture.
//Add a simple draw function for blocks.
//In main, create a slice of blocks ([]Block) to build your floor and platforms by positioning multiple blocks.
//Update your draw function to loop through and draw all blocks before the player.
//Later, add collision detection using rl.CheckCollisionRecs to detect when the player overlaps blocks and adjust their position accordingly

package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	Width          = 1000
	Height         = 800
	MoveDistance   = 5
	Gravity        = 1
	FPS            = 60
	AnimationDelay = 5 // Frames to wait before changing sprite
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

// Load a sprite sheet and split it into frames
// this function is called when a new actor is created
func splitSpriteSheet(spriteSheet rl.Texture2D, frameWidth, frameHeight int32) []rl.Rectangle {
	frames := []rl.Rectangle{}

	// Calculate how many frames fit horizontally
	frameCount := spriteSheet.Width / frameWidth

	// Create a rectangle for each frame
	for i := int32(0); i < frameCount; i++ {
		frame := rl.NewRectangle(
			float32(i*frameWidth), // X position in sprite sheet
			0,                     // Y position (0 if single row)
			float32(frameWidth),   // Width of frame
			float32(frameHeight),  // Height of frame
		)
		frames = append(frames, frame)
	}
	//this returns a slice of rectangles
	return frames
}

type Actor struct {
	rl.Rectangle
	Texture        rl.Texture2D
	Frames         []rl.Rectangle // Source rectangles for each frame
	CurrentFrame   int            // Which frame we're on
	AnimationCount int
	Xvel           float32
	Yvel           float32
	Direction      string
	FallCount      int
}

func newActor(texture rl.Texture2D, frameWidth, frameHeight int32, x, y float32) *Actor {
	frames := splitSpriteSheet(texture, frameWidth, frameHeight)

	return &Actor{Rectangle: rl.Rectangle{X: x, Y: y, Width: float32(frameWidth), Height: float32(frameHeight)}, Texture: texture, Frames: frames, Direction: "left"}
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

func (a *Actor) updateAnimation() {
	a.AnimationCount++

	if a.AnimationCount >= AnimationDelay {
		a.AnimationCount = 0
		a.CurrentFrame = (a.CurrentFrame + 1) % len(a.Frames) // Loop through frames, once it its the last frame, mod would make it go back to zero
	}
}

type Block struct {
	rl.Rectangle
	Texture rl.Texture2D
}

func getBlock(x, y, size int) {
}

func update(player *Actor, frog map[string]rl.Texture2D) {
	player.handleMove()
	player.updateAnimation()

	//This is for Gravity
	//player.Yvel += float32(min(1.0, float64(player.FallCount)/FPS*Gravity))
	//player.FallCount += 1

	// Apply velocity
	if player.Xvel != 0 {
		player.X += player.Xvel
		if player.Texture != frog["run"] {
			player.Texture = frog["run"]
		}
	} else if player.Xvel == 0 && player.Texture == frog["run"] {
		player.Texture = frog["normal"]
	}
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
	drawFrog(player)

	rl.EndDrawing()
}

func drawFrog(player *Actor) {
	// Get the source rectangle for the current frame
	src := player.Frames[player.CurrentFrame]

	dst := rl.NewRectangle(player.X, player.Y, player.Width, player.Height)
	origin := rl.NewVector2(0, 0)

	if player.Direction == "left" {
		// Flip horizontally by making source width negative
		src.Width = -src.Width
		// Shift the source rect start so it doesn't disappear
		src.X += player.Width
	}

	rl.DrawTexturePro(player.Texture, src, dst, origin, 0, rl.White)
}

func main() {
	rl.InitWindow(Width, Height, "Platformer Game")
	defer rl.CloseWindow()
	rl.SetTargetFPS(FPS)

	background := rl.LoadTexture("images/Background/Yellow.png")
	defer rl.UnloadTexture(background)
	tiles := getBackground(background)

	theFrogTextures := map[string]rl.Texture2D{
		"run":    rl.LoadTexture("images/run.png"),
		"normal": rl.LoadTexture("images/idle.png"),
	}

	for _, texture := range theFrogTextures {
		defer rl.UnloadTexture(texture)
	}

	// Specify the frame width and height for your sprite sheet
	// For example, if each frame is 32x32 pixels:
	player := newActor(theFrogTextures["run"], 32, 32, 100, 100)

	// Game loop
	for !rl.WindowShouldClose() {
		update(player, theFrogTextures)
		draw(background, tiles, player)
	}
}
