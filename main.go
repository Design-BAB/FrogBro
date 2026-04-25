// Author: Design-BAB
// Date: 3/23/2026

package main

import (
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	Width          = 800
	Height         = 672
	MoveDistance   = 5
	Gravity        = 1
	FPS            = 60
	BlockTextureX  = 96
	BlockTextureY  = 0
	BlockSize      = 32
	JumpHeight     = -5
	AnimationDelay = 5 // Frames to wait before changing sprite
	FlapDelay      = 9
	TextSize       = 20
)

type GameState struct {
	Score           int
	isOver          bool
	numberOfUpdates int
}

func newGame() *GameState {
	return &GameState{}
}

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
	for i := range frameCount {
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

type Fly struct {
	Texture rl.Texture2D
	rl.Rectangle
}

func newFly(texture rl.Texture2D, x, y float32) *Fly {
	return &Fly{Texture: texture, Rectangle: rl.Rectangle{X: x, Y: y, Width: float32(texture.Width), Height: float32(texture.Height)}}
}

type Object struct {
	Texture rl.Texture2D
	rl.Rectangle
}

type Level struct {
	Flies  [3]*Fly
	Salt   *Object
	Blocks []*Block
}

func newLevel(blocks []*Block, flies [3]*Fly, salt *Object) *Level {
	return &Level{Blocks: blocks, Flies: flies, Salt: salt}
}

// I noticed that this particular object deals pratically with int32 so, the axis is going to stay as int32 instead of float32
func newObject(texture rl.Texture2D, x, y int32) *Object {
	return &Object{Texture: texture, Rectangle: rl.Rectangle{X: float32(x), Y: float32(y), Width: float32(texture.Width), Height: float32(texture.Height)}}
}

type Actor struct {
	rl.Rectangle
	Texture        rl.Texture2D
	Frames         []rl.Rectangle // Source rectangles for each frame
	CurrentFrame   int            // Which frame we're on
	AnimationCount int
	Xvel           float32
	Yvel           float32
	FacingRight    bool
	FallCount      int
	JumpCount      int
}

func newActor(texture rl.Texture2D, frameWidth, frameHeight int32, x, y float32) *Actor {
	frames := splitSpriteSheet(texture, frameWidth, frameHeight)

	return &Actor{Rectangle: rl.Rectangle{X: x, Y: y, Width: float32(frameWidth), Height: float32(frameHeight)}, Texture: texture, Frames: frames}
}

func (frog *Actor) updateAnimation() {
	frog.AnimationCount++

	if frog.AnimationCount >= AnimationDelay {
		frog.AnimationCount = 0
		frog.CurrentFrame = (frog.CurrentFrame + 1) % len(frog.Frames) // Loop through frames, once it its the last frame, mod would make it go back to zero
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

func makeBlockRow(texture rl.Texture2D, startX, y, count int) []*Block {
	var blocks []*Block
	for i := range count {
		blocks = append(blocks, newBlock(texture, startX+i*BlockSize, y, BlockSize))
	}
	return blocks
}

func createLevel(blockTexture rl.Texture2D, flyTextures *[2]rl.Texture2D, saltTexture rl.Texture2D) *Level {
	blocks := []*Block{}

	// Tier 1 (y=640)
	blocks = append(blocks, makeBlockRow(blockTexture, 192, 640, 3)...)
	blocks = append(blocks, makeBlockRow(blockTexture, 608, 640, 3)...)
	// Tier 2 (y=544)
	blocks = append(blocks, makeBlockRow(blockTexture, 0, 544, 4)...)
	blocks = append(blocks, makeBlockRow(blockTexture, 416, 544, 5)...)
	// Tier 3 (y=448)
	blocks = append(blocks, makeBlockRow(blockTexture, 224, 448, 4)...)
	blocks = append(blocks, makeBlockRow(blockTexture, 640, 448, 3)...)
	// Tier 4 (y=352)
	blocks = append(blocks, makeBlockRow(blockTexture, 0, 352, 4)...)
	blocks = append(blocks, makeBlockRow(blockTexture, 448, 352, 4)...)
	// Tier 5 (y=256)
	blocks = append(blocks, makeBlockRow(blockTexture, 224, 256, 7)...)
	blocks = append(blocks, makeBlockRow(blockTexture, 640, 256, 2)...)
	// Goal platform
	blocks = append(blocks, makeBlockRow(blockTexture, 600, 160, 7)...)

	var flys [3]*Fly
	flys[0] = newFly(flyTextures[0], 200, 580)
	flys[1] = newFly(flyTextures[1], 300, 400)
	flys[2] = newFly(flyTextures[0], 100, 250)

	salt := newObject(saltTexture, 230, 256-saltTexture.Height)

	return newLevel(blocks, flys, salt)
}

func update(player *Actor, frog map[string]rl.Texture2D, level *Level, flyTextures *[2]rl.Texture2D, blockTexture, saltTexture rl.Texture2D, yourGame *GameState) {
	if yourGame.isOver == false {
		player.handleMove()
		player.updateAnimation()

		//This is for Gravity
		player.Yvel += float32(min(1.0, float64(player.FallCount)/FPS*Gravity))
		player.FallCount += 1

		// Apply both velocities
		player.X += player.Xvel
		player.Y += player.Yvel

		// Resolve collisions in a single pass
		handleCollision(player, level, yourGame)

		// Update texture based on movement
		if player.Xvel != 0 {
			if player.Texture != frog["run"] {
				player.Texture = frog["run"]
			}
		} else if player.Texture == frog["run"] {
			player.Texture = frog["normal"]
		}
		if yourGame.numberOfUpdates == FlapDelay {
			level.Flies = flap(level.Flies, flyTextures)
			yourGame.numberOfUpdates = 0
		} else {
			yourGame.numberOfUpdates += 1
		}
		//collision with the window
		player.X = rl.Clamp(player.X, 0.0, Width-player.Width)
		player.Y = rl.Clamp(player.Y, 0.0, Height-player.Height)
		// Reset jump when hitting the bottom edge
		if player.Y >= Height-player.Height {
			player.Yvel = 0
			player.FallCount = 0
			player.JumpCount = 0
		}
	} else {
		if rl.IsKeyPressed(rl.KeyX) {
			yourGame.isOver = false
			yourGame.Score = 0

			*level = *createLevel(blockTexture, flyTextures, saltTexture)

			player.X = 50
			player.Y = 70
			player.Xvel = 0
			player.Yvel = 0
			player.JumpCount = 0
			player.FallCount = 0
		}
	}
}

func (frog *Actor) handleMove() {
	frog.Xvel = 0 // Reset horizontal velocity each frame
	if rl.IsKeyDown(rl.KeyRight) {
		frog.Xvel = MoveDistance
		if frog.FacingRight == false {
			frog.FacingRight = true
		}
	}
	if rl.IsKeyDown(rl.KeyLeft) {
		frog.Xvel = -MoveDistance
		if frog.FacingRight {
			frog.FacingRight = false
		}
	}
	if rl.IsKeyPressed(rl.KeyUp) {
		if frog.JumpCount < 2 {
			frog.Yvel = JumpHeight
			frog.FallCount = 0
			frog.JumpCount++
		}
	}
}

func flap(flys [3]*Fly, textures *[2]rl.Texture2D) [3]*Fly {
	for _, fly := range flys {
		if fly.Texture == textures[0] {
			fly.Texture = textures[1]
		} else {
			fly.Texture = textures[0]
		}
	}
	return flys
}

func handleCollision(player *Actor, level *Level, yourGame *GameState) {
	for _, fly := range level.Flies {
		if rl.CheckCollisionRecs(player.Rectangle, fly.Rectangle) {
			//just gonna make it "disappear"
			fly.X = 900
			fly.Y = 900
			yourGame.Score++
		}
	}
	if rl.CheckCollisionRecs(player.Rectangle, level.Salt.Rectangle) {
		yourGame.isOver = true
	}
	for _, block := range level.Blocks {
		if !rl.CheckCollisionRecs(player.Rectangle, block.Rectangle) {
			continue
		}

		// Calculate how far the player overlaps into the block from each side
		overlapLeft := (player.X + player.Width) - block.X
		overlapRight := (block.X + block.Width) - player.X
		overlapTop := (player.Y + player.Height) - block.Y
		overlapBottom := (block.Y + block.Height) - player.Y

		// Find the smallest overlap per axis
		overlapX := min(overlapLeft, overlapRight)
		overlapY := min(overlapTop, overlapBottom)

		// Resolve the axis with the smallest penetration
		if overlapX < overlapY {
			if player.Xvel > 0 {
				player.X = block.X - player.Width
			} else {
				player.X = block.X + block.Width
			}
			player.Xvel = 0
		} else {
			if player.Yvel > 0 {
				player.Y = block.Y - player.Height
				player.Yvel = 0
				player.FallCount = 0
				player.JumpCount = 0
			} else {
				player.Y = block.Y + block.Height
				player.Yvel = 0
			}
		}
	}
}

func draw(background rl.Texture2D, tiles []rl.Vector2, level *Level, player *Actor, door *Object, game *GameState) {
	rl.BeginDrawing()
	rl.ClearBackground(rl.RayWhite)
	// Draw background tiles
	for _, tile := range tiles {
		rl.DrawTextureV(background, tile, rl.White)
	}
	if game.isOver == false {
		//draw the door
		rl.DrawTexture(door.Texture, int32(door.X), int32(door.Y), rl.White)
		rl.DrawTexture(level.Salt.Texture, int32(level.Salt.X), int32(level.Salt.Y), rl.White)

		//draw blocks
		for _, block := range level.Blocks {
			drawBlock(block)
		}

		// Draw player
		drawFrog(player)

		for _, fly := range level.Flies {
			if fly.X < Width {
				rl.DrawTexture(fly.Texture, int32(fly.X), int32(fly.Y), rl.White)
			}
		}
		//GUI
		rl.DrawText("Your score is "+strconv.Itoa(game.Score), 20, 20, TextSize, rl.DarkGray)
	} else {
		rl.DrawText("Game over", Width/3-TextSize, Height/3-TextSize, TextSize, rl.DarkGray)
		rl.DrawText("Press 'X' to play again.", Width/3-TextSize, Height/2-TextSize, TextSize, rl.DarkGray)

	}
	rl.EndDrawing()
}

func drawFrog(player *Actor) {
	// Get the source rectangle for the current frame
	src := player.Frames[player.CurrentFrame]

	dst := rl.NewRectangle(player.X, player.Y, player.Width, player.Height)
	origin := rl.NewVector2(0, 0)

	if player.FacingRight == false {
		// Flip horizontally by making source width negative
		src.Width = -src.Width
		// Shift the source rect start so it doesn't disappear
		src.X += player.Width
	}

	rl.DrawTexturePro(player.Texture, src, dst, origin, 0, rl.White)
}

func main() {
	rl.InitWindow(Width, Height, "Frog Bros")
	defer rl.CloseWindow()
	rl.SetTargetFPS(FPS)
	game := newGame()
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
	// Player starts on the left ground floor
	player := newActor(theFrogTextures["run"], BlockSize, BlockSize, 50, 70)
	//fly code
	var flyTextures [2]rl.Texture2D
	flyTextures[0] = rl.LoadTexture("images/FlyUp.png")
	defer rl.UnloadTexture(flyTextures[0])
	flyTextures[1] = rl.LoadTexture("images/FlyDown.png")
	defer rl.UnloadTexture(flyTextures[1])
	//door code
	doorTexture := rl.LoadTexture("images/door.png")
	defer rl.UnloadTexture(doorTexture)
	door := newObject(doorTexture, Width-5-doorTexture.Width, 20)
	//salt
	saltTexture := rl.LoadTexture("images/salt.png")
	defer rl.UnloadTexture(saltTexture)
	level := createLevel(blockTexture, &flyTextures, saltTexture)
	// Game loop
	for !rl.WindowShouldClose() {
		update(player, theFrogTextures, level, &flyTextures, blockTexture, saltTexture, game)
		draw(background, tiles, level, player, door, game)
	}
}
