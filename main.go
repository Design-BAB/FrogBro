// Author: Design-BAB
// Date: 3/21/2026
// Description: It is my platform game! The goal is to reach 268 lines of code
// Notes: My next step should be to continue working on collisions

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
	BlockTextureX  = 96
	BlockTextureY  = 0
	AnimationDelay = 5 // Frames to wait before changing sprite
)

func newVector(x, y int32) rl.Vector2 {
	return rl.Vector2{X: float32(x), Y: float32(y)}
}

func getBackground(background rl.Texture2D) []rl.Vector2 {
	//this creates an array of positions for the background tile
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
	Frame   rl.Rectangle
	Texture rl.Texture2D
}

// This newBlock constructor takes x, y, and size, then calls getBlock to set the texture.
func newBlock(texture rl.Texture2D, x, y, size int) *Block {
	//this x and y is the location on the .png of where the block design is
	frame := getBlock(BlockTextureX, BlockTextureY, size)
	//this x and y is where it goes in the actual game
	whereItGoes := rl.Rectangle{X: float32(x), Y: float32(y), Width: float32(size), Height: float32(size)}
	return &Block{Rectangle: whereItGoes, Texture: texture, Frame: frame}
}

func getBlock(x, y, size int) rl.Rectangle {
	theBlockLocation := rl.NewRectangle(float32(x), float32(y), float32(size), float32(size))
	return theBlockLocation
}

func drawBlock(blockToDraw *Block) {
	pos := rl.NewVector2(blockToDraw.X, blockToDraw.Y)
	rl.DrawTextureRec(blockToDraw.Texture, blockToDraw.Frame, pos, rl.White)
}

func update(player *Actor, frog map[string]rl.Texture2D, blocks []*Block) {
	player.handleMove()
	player.updateAnimation()

	//This is for Gravity
	player.Yvel += float32(min(1.0, float64(player.FallCount)/FPS*Gravity))
	player.FallCount += 1

	// Apply velocity
	if player.Xvel != 0 {
		player.X += player.Xvel
		if player.Texture != frog["run"] {
			player.Texture = frog["run"]
		}
	} else if player.Xvel == 0 && player.Texture == frog["run"] {
		player.Texture = frog["normal"]
	}

	handleVerticalCollision(player, blocks)

	player.X += player.Xvel
	player.Y += player.Yvel

	//collision with the window
	player.X = rl.Clamp(player.X, 0.0, Width-player.Width)
	player.Y = rl.Clamp(player.Y, 0.0, Height-player.Height)
}

func handleVerticalCollision(player *Actor, blocks []*Block) {
	for _, block := range blocks {
		//if collide and if Yvel > 0...then set to zero
		if rl.CheckCollisionRecs(player.Rectangle, block.Rectangle) {
			if player.Yvel > 0 {
				player.Yvel = 0
				player.FallCount = 0
			} else if player.Yvel < 0 {
				player.Yvel = 0
				//if they hit their head, reverse and make them fall down
				player.Yvel *= -1
			}
		}
	}
}

func draw(background rl.Texture2D, tiles []rl.Vector2, blocks []*Block, player *Actor) {
	rl.BeginDrawing()
	rl.ClearBackground(rl.RayWhite)

	// Draw background tiles
	for _, tile := range tiles {
		rl.DrawTextureV(background, tile, rl.White)
	}

	//draw blocks
	for _, block := range blocks {
		drawBlock(block)
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

	//this deals with the background
	background := rl.LoadTexture("images/Background/Yellow.png")
	defer rl.UnloadTexture(background)
	//this creates an array of tiles, with their location set
	tiles := getBackground(background)

	theFrogTextures := map[string]rl.Texture2D{
		"run":    rl.LoadTexture("images/run.png"),
		"normal": rl.LoadTexture("images/idle.png"),
	}
	for _, texture := range theFrogTextures {
		defer rl.UnloadTexture(texture)
	}

	blockTexture := rl.LoadTexture("images/Terrain.png")
	defer rl.UnloadTexture(blockTexture)
	//this is where we set the blocks, in the future this can be recorded in a JSON file
	blocks := []*Block{
		newBlock(blockTexture, 0, 700, 32),
		newBlock(blockTexture, 32, 700, 32),
		newBlock(blockTexture, 64, 700, 32),
		newBlock(blockTexture, 96, 700, 32),
	}
	// Specify the frame width and height for your sprite sheet
	// For example, if each frame is 32x32 pixels:
	player := newActor(theFrogTextures["run"], 32, 32, 100, 100)

	// Game loop
	for !rl.WindowShouldClose() {
		update(player, theFrogTextures, blocks)
		draw(background, tiles, blocks, player)
	}
}
