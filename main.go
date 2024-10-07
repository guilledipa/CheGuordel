package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"image/color"
	"io"
	"log"
	"os"
	"strings"
	"time"
	"unicode"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	text "github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/exp/rand"
	"golang.org/x/exp/utf8string"
)

var (
	fontFaceSource *text.GoTextFaceSource
	validWords     map[string]any
)

//go:embed fonts/NotoSans-Regular.ttf
var notoSansRegularTTF []byte

// Mode defines the current game mode.
type Mode int

const (
	// modeTitle  Mode = iota
	modeGame Mode = iota
	modeGameOver
)

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
	mode       Mode
	targetWord *utf8string.String
	guesses    [][]rune
	currentRow int
	gameWon    bool
	font       text.Face
}

// Update does stuff..
func (g *Game) Update() error {
	switch g.mode {
	case modeGame:
		if g.gameWon {
			return nil // Game over, stop processing inputs
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			g.handleEnter()
		} else if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
			g.handleBackspace()
		} else {
			if inpututil.IsKeyJustPressed(ebiten.KeySemicolon) {
				g.handleLetterInput('Ñ')
				break
			}
			for _, r := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
				// Calculate the correct Ebiten key code for the letter
				key := ebiten.KeyA + ebiten.Key(unicode.ToUpper(r)-'A')
				if inpututil.IsKeyJustPressed(key) {
					g.handleLetterInput(r)
					break
				}
			}
		}
		// Check if we've already reached the maximum guesses
		if g.currentRow == maxGuesses {
			g.mode = modeGameOver
		}
	case modeGameOver:
		// Not implemented yet.
		return nil
	}
	return nil
}

// Draw draws the game screen and is called after every update.
func (g *Game) Draw(screen *ebiten.Image) {
	if g.mode == modeGameOver {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("Perdiste... La palabra era: %s", g.targetWord))
	}
	f := &text.GoTextFace{
		Source: fontFaceSource,
		Size:   normalFontSize,
	}
	// Draw grid
	for row := 0; row < maxGuesses; row++ {
		for col := 0; col < wordLength; col++ {
			x := col*tileSize + (col+1)*tileSpacing
			y := row*tileSize + (row+1)*tileSpacing
			// Draw tile background
			tileColor := color.Color(color.White)
			// Only color tiles in completed rows
			if row < g.currentRow || g.gameWon {
				tileColor = g.getTileColor(row, col)
			}
			// Draw tile
			vector.DrawFilledRect(screen, float32(x), float32(y), float32(tileSize), float32(tileSize), tileColor, false)
			// Draw letter
			if row <= g.currentRow && col < len(g.guesses[row]) {
				letter := string(g.guesses[row][col])
				// Calculate text bounds to get width and height
				textWidth, textHeight := text.Measure(letter, f, 0)
				// Calculate position to center the letter within the tile
				textX := float64(x) + (float64(tileSize)-textWidth)/2
				textY := float64(y) + (float64(tileSize)-textHeight)/2
				// Draw
				op := &text.DrawOptions{}
				op.GeoM.Translate(textX, textY)
				op.ColorScale.ScaleWithColor(color.Black)
				text.Draw(screen, letter, f, op)
				ebitenutil.DebugPrint(screen, letter)
			}
		}
	}
	// Display game over message
	if g.gameWon {
		ebitenutil.DebugPrint(screen, "Ganaste!")
	}
}

// Layout takes the outside size (pixel), and returns the logical screen size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) handleEnter() error {
	// Do nothing if current guess word length is is not equal to wordLength
	if len(g.guesses[g.currentRow]) != wordLength {
		return nil
	}

	// Check if the word is in the valid word list
	if !isValidWord(strings.ToLower(string(g.guesses[g.currentRow]))) {
		// Handle invalid word (e.g., display an error message)
		fmt.Printf("%q no es una palabra válida.\n", g.guesses[g.currentRow])
		return nil
	}
	// Check if the guess is correct
	if string(g.guesses[g.currentRow]) == g.targetWord.String() {
		g.gameWon = true
		return nil
	}
	// Go to next row if the guess is incorrect
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
		g.guesses[g.currentRow] = append(g.guesses[g.currentRow], r)
	}
	fmt.Println(g.guesses[g.currentRow])
	return nil
}

func (g *Game) getTileColor(row, col int) color.Color {
	if row >= 0 && row < len(g.guesses) && col >= 0 && col < len(g.guesses[row]) {
		letter := g.guesses[row][col]
		if g.targetWord.At(col) == letter {
			return color.RGBA{0, 255, 0, 255} // Green
		} else {
			for i := 0; i < wordLength; i++ {
				if g.targetWord.At(i) == letter {
					return color.RGBA{255, 255, 0, 255} // Yellow
				}
			}
		}
	}
	return color.RGBA{128, 128, 128, 255} // Gray
}

// Check if a given word is valid (UPdate)
func isValidWord(word string) bool {
	if _, ok := validWords[word]; ok {
		return true
	}
	return false
}

// randomWord returns a random validGuesses key.
// This approach relies on the semantics of range; range will visit every
// key/value in the map exactly once in **unspecified** order.
// I'll use this as a way of "randomly" choose a target word.
func randomWord(m map[string]any) string {
	rand.Seed(uint64(time.Now().UnixNano()))
	r := rand.Intn(len(m))
	for k := range m {
		if r == 0 {
			return k
		}
		r--
	}
	panic("randomWord: unreachable point.")
}

func init() {
	// Load font
	var err error
	fontFaceSource, err = text.NewGoTextFaceSource(bytes.NewReader(notoSansRegularTTF))
	if err != nil {
		log.Fatal(err)
	}
	// Load word list
	jsonFile, err := os.Open("constants/validGuesses.json")
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()
	bytes, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal([]byte(bytes), &validWords)
}

func main() {
	// Initialize game
	g := &Game{
		targetWord: utf8string.NewString(strings.ToUpper(randomWord(validWords))),
		guesses:    make([][]rune, maxGuesses),
		currentRow: 0,
		gameWon:    false,
	}
	fmt.Println(g.targetWord)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("CheGuordle!")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
