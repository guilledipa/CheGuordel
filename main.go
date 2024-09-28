package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"image/color"
	"log"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	text "github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/exp/rand"
)

var (
	robotoFaceSource *text.GoTextFaceSource
)

//go:embed fonts/Roboto-Regular.ttf
var robotoRegularTTF []byte

const (
	screenWidth    = 800
	screenHeight   = 600
	maxGuesses     = 6
	wordLength     = 5
	tileSize       = 50
	tileSpacing    = 10
	normalFontSize = 24
)

// Game implements ebiten.Game interface
type Game struct {
	// Do stuff, e.g. game state vars: current word, guesses, etc.
	targetWord string
	guesses    []string
	currentRow int
	gameWon    bool
	font       text.Face
}

// Update does stuff..
func (g *Game) Update() error {
	// Handle inputs, game logic, etc
	if g.gameWon {
		return nil // Game over, stop processing inputs
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		g.handleEnter()
	} else if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		g.handleBackspace()
	} else {
		for _, r := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
			if inpututil.IsKeyJustPressed(ebiten.Key(r)) {
				g.handleLetterInput(rune(r))
				break
			}
		}
	}
	return nil
}

// Draw draws the game screen and is called after every update.
func (g *Game) Draw(screen *ebiten.Image) {
	f := &text.GoTextFace{
		Source: robotoFaceSource,
		Size:   normalFontSize,
	}
	// Draw grid
	for row := 0; row < maxGuesses; row++ {
		for col := 0; col < wordLength; col++ {
			x := col*tileSize + (col+1)*tileSpacing
			y := row*tileSize + (row+1)*tileSpacing
			// Draw tile background
			tileColor := color.Color(color.White)
			// Only color tiles in completed rows or the current row
			if row < g.currentRow || (row == g.currentRow && col < len(g.guesses[row])) {
				tileColor = g.getTileColor(row, col)
			}
			// Draw tile
			vector.DrawFilledRect(screen, float32(x), float32(y), float32(tileSize), float32(tileSize), tileColor, false)
			// Draw letter
			if row < g.currentRow && col < len(g.guesses[row]) {
				letter := g.guesses[row][col]
				// Calculate text bounds to get width and height
				textWidth, textHeight := text.Measure(string(letter), f, 0)
				// Calculate position to center the letter within the tile
				textX := float64(x) + (float64(tileSize)-textWidth)/2
				textY := float64(y) + (float64(tileSize)-textHeight)/2 + textHeight
				// Draw
				op := &text.DrawOptions{}
				op.GeoM.Translate(textX, textY)
				op.ColorScale.ScaleWithColor(color.Black)
				text.Draw(screen, string(letter), f, op)
			}
		}
	}
	// Display game over message
	if g.gameWon {
		ebitenutil.DebugPrint(screen, "Ganaste!")
	} else if g.currentRow == maxGuesses {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("Perdiste... La palabra era: %s", g.targetWord))
	}
}

// Layout takes the outside size (pixel), and returns the logical screen size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) handleEnter() error {
	g.currentRow++
	return nil
}

func (g *Game) handleBackspace() error {
	if len(g.guesses[g.currentRow]) > 0 {
		g.guesses[g.currentRow] = g.guesses[g.currentRow][:len(g.guesses[g.currentRow])-1]
	}
	return nil
}

func (g *Game) handleLetterInput(r rune) error {
	if len(g.guesses[g.currentRow]) < wordLength {
		g.guesses[g.currentRow] = g.guesses[g.currentRow] + strings.ToUpper(string(r))
	}
	return nil
}

func (g *Game) getTileColor(row, col int) color.Color {
	if row >= 0 && row < len(g.guesses) && col >= 0 && col < len(g.guesses[row]) {
		letter := g.guesses[row][col]
		if g.targetWord[col] == letter {
			return color.RGBA{0, 255, 0, 255} // Green
		} else if strings.ContainsRune(g.targetWord, rune(letter)) {
			return color.RGBA{255, 255, 0, 255} // Yellow
		}
	}
	return color.RGBA{128, 128, 128, 255} // Gray
}

func init() {
	// Load font
	s, err := text.NewGoTextFaceSource(bytes.NewReader(robotoRegularTTF))
	if err != nil {
		log.Fatal(err)
	}
	robotoFaceSource = s
}

func main() {
	// Choose a random word (Temporal)
	wordList := []string{"MESSI", "CHETO", "CHINO", "PUCHO", "BOLSA"}
	rand.Seed(uint64(time.Now().UnixNano()))
	targetWord := wordList[rand.Intn(len(wordList))]
	fmt.Println(targetWord)
	// Initialize game
	g := &Game{
		targetWord: targetWord,
		guesses:    make([]string, maxGuesses),
		currentRow: 0,
		gameWon:    false,
	}
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("CheGuordle!")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
