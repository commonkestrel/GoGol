package main

import (
	"image"
	"os"
	"path"
	"runtime"

	"image/color"
	_ "image/png"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const (
    SCREENX, SCREENY = 960, 540
    BOARDX, BOARDY = 48, 27 // resolution/20
)

var (
    win *pixelgl.Window
    imd *imdraw.IMDraw
)

type board struct {
    cells [BOARDX][BOARDY]bool
    neighbours [BOARDX][BOARDY]uint8
}

func NewBoard() *board {
    return &board{}
}

func (b *board) Set(x, y int, state bool) {
    b.cells[x][y] = state
}

func (b *board) Get(x, y int) bool {
    return b.cells[x][y]
}

func (b *board) Neighbours(x, y int) int {
    return int(b.neighbours[x][y])
}

func (b *board) UpdateNeighbors() {
    for x := 0; x < BOARDX; x++ {
        for y := 0; y < BOARDY; y++ {
            var neighbours uint8
        }
    }
}

func (b *board) Update() {
    for x := 0; x < BOARDX; x++ {
        for y := 0; y < BOARDY; y++ {
            neighbours := b.Neighbours(x, y)
            if b.Get(x, y) { // check if the cell is alive
                if neighbours > 3 || neighbours < 2 {
                    b.Set(x, y, false)
                }
            } else if neighbours == 3 {
                b.Set(x, y, true)
            }
        }
    }
}

func (b *board) Draw() {
    imd.Color = colornames.White
    for x := 0; x < BOARDX; x++ {
        for y := 0; y < BOARDY; y++ {
            if b.Get(x, y) {
                imd.Push(pixel.V(float64(x*20), float64(y*20)))
                imd.Push(pixel.V(float64(x*20)+20, float64(y*20)+20))
                imd.Rectangle(0)
            }
        }
    }
}

// used for loading icons and sprites
func LoadPicture(path string) (pixel.Picture, error) {
    // loads and decodes PNG
    file, err := os.Open(path)
    if err != nil {
        panic(err)
    }
    defer file.Close()
    img, _, err := image.Decode(file)
    if err != nil {
        panic(err)
    }
    // converts to Pixel picture
    return pixel.PictureDataFromImage(img), nil
}

// returns the absolute path of a path relative to the file's parent directory
func relative(relative string) string {
    _, filepath, _, _ := runtime.Caller(0)
    dir := path.Dir(filepath)
    return path.Join(dir, relative)
}

func run() {
    iconpath := relative("icon.png")
    icon, err := LoadPicture(iconpath)
    if err != nil {
        panic(err)
    }

    cfg := pixelgl.WindowConfig{
        Title:     "Go Pixel",
        Bounds:    pixel.R(0, 0, SCREENX, SCREENY),
        Maximized: true,
        Icon:      []pixel.Picture{icon},
        VSync: true,
    }
    win, err = pixelgl.NewWindow(cfg)
    if err != nil {
        panic(err)
    }

    imd = imdraw.New(nil)

    for !win.Closed() {
        imd.Clear()

        // run game loop here

        win.Clear(colornames.Black)
        imd.Draw(win)
        win.Update()
    }
}

func main() {
    pixelgl.Run(run)
}