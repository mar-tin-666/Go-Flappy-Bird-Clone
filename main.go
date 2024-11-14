package main

import (
	"image/color"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Screen and game settings
const (
	screenWidth      = 500
	screenHeight     = 500
	birdSize         = 20
	gravity          = 0.5
	jumpStrength     = -6
	pipeWidth        = 50
	pipeGap          = 150
	pipeSpeed        = 2
	pipeSpacing      = 250
	startMessageX    = screenWidth/2 - 110
	startMessageY    = screenHeight/2 - 30
	gameOverMessageX = screenWidth/2 - 100
	gameOverMessageY = screenHeight / 2
)

// Text messages
var (
	textWindowTitle  = "Go Flappy Bird Clone"
	textStartMessage = "Press SPACE to start and for jump"
	textGameOver     = "Game Over! Press R to restart"
	textScorePrefix  = "Score: "
)

// Colors
var (
	colorBackground = color.RGBA{120, 200, 240, 255} // Blue (sky)
	colorBird       = color.RGBA{255, 255, 0, 255}   // Yellow
	colorPipe       = color.RGBA{40, 140, 40, 255}   // Green
)

// Pipe structure
type Pipe struct {
	x      float64
	height float64
	passed bool
}

// Game structure
type Game struct {
	birdY         float64
	birdVelocity  float64
	pipes         []Pipe
	score         int
	gameOver      bool
	started       bool
	nextPipeSpawn float64
}

// Resets the game state
func (g *Game) reset() {
	g.birdY = screenHeight / 2
	g.birdVelocity = 0
	g.pipes = nil
	g.score = 0
	g.gameOver = false
	g.started = false
	g.nextPipeSpawn = screenWidth
}

// Initialize a random number generator
var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// Spawns a new pipe with a random gap position
func (g *Game) spawnPipe() {
	topHeight := float64(rng.Intn(screenHeight-pipeGap-birdSize*2) + birdSize)
	g.pipes = append(g.pipes, Pipe{x: screenWidth, height: topHeight, passed: false})
}

// Game update logic
func (g *Game) Update() error {
	// Game starts with the first space press
	if !g.started {
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			g.started = true
		}
		return nil
	}

	// Handle game over reset
	if g.gameOver {
		if ebiten.IsKeyPressed(ebiten.KeyR) {
			g.reset()
		}
		return nil
	}

	// Handle bird jump
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.birdVelocity = jumpStrength
	}

	// Update bird position
	g.birdVelocity += gravity
	g.birdY += g.birdVelocity

	// Check if the bird hits the ground
	if g.birdY > screenHeight-birdSize {
		g.birdY = screenHeight - birdSize
		g.gameOver = true
	}

	// Check if the bird hits the top
	if g.birdY < 0 {
		g.birdY = 0
		g.gameOver = true
	}

	// Spawn new pipes at intervals
	if g.nextPipeSpawn <= 0 {
		g.spawnPipe()
		g.nextPipeSpawn = pipeSpacing
	}
	g.nextPipeSpawn -= pipeSpeed

	// Update pipes and check for collisions
	for i := range g.pipes {
		g.pipes[i].x -= pipeSpeed

		// Check if the bird passes through the pipes for scoring
		if !g.pipes[i].passed && g.pipes[i].x+pipeWidth < screenWidth/2 {
			g.score++
			g.pipes[i].passed = true
		}

		// Collision detection with pipes
		if (screenWidth/2 >= g.pipes[i].x && screenWidth/2 <= g.pipes[i].x+pipeWidth) &&
			(g.birdY <= g.pipes[i].height || g.birdY+birdSize >= g.pipes[i].height+pipeGap) {
			g.gameOver = true
		}
	}

	// Remove pipes that have moved off screen
	if len(g.pipes) > 0 && g.pipes[0].x+pipeWidth < 0 {
		g.pipes = g.pipes[1:]
	}

	return nil
}

// Draw a filled rectangle with the specified color
func drawFilledRect(screen *ebiten.Image, x, y, width, height float32, clr color.Color) {
	vector.DrawFilledRect(screen, x, y, width, height, clr, false) // false means no anti-aliasing
}

// Drawing logic for the game
func (g *Game) Draw(screen *ebiten.Image) {
	// Clear the screen with background color
	screen.Fill(colorBackground)

	// Draw the bird
	drawFilledRect(screen, float32(screenWidth/2-birdSize/2), float32(g.birdY), birdSize, birdSize, colorBird)

	// Draw pipes
	for _, pipe := range g.pipes {
		// Top pipe
		drawFilledRect(screen, float32(pipe.x), 0, pipeWidth, float32(pipe.height), colorPipe)
		// Bottom pipe
		drawFilledRect(screen, float32(pipe.x), float32(pipe.height+pipeGap), pipeWidth, float32(screenHeight-pipe.height-pipeGap), colorPipe)
	}

	// Display score in the top left corner
	ebitenutil.DebugPrintAt(screen, textScorePrefix+strconv.Itoa(g.score), 10, 10)

	// Display start message if game hasn't started
	if !g.started {
		ebitenutil.DebugPrintAt(screen, textStartMessage, startMessageX, startMessageY)
	}

	// Display Game Over message if game is over
	if g.gameOver {
		ebitenutil.DebugPrintAt(screen, textGameOver, gameOverMessageX, gameOverMessageY)
	}
}

// Defines the layout of the game window
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	game := &Game{}
	game.reset()

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle(textWindowTitle)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
