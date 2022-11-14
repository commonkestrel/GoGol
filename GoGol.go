package main

import (
	"image"
	"os"
	"path"
	"runtime"
	"time"

	_ "image/png"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const (
    SCREENX, SCREENY = 960, 540
    BoardX, BoardY = 48, 27 // resolution/20
)

var (
    win *pixelgl.Window
    imd *imdraw.IMDraw
    running bool
)

type Board struct {
    cells [BoardX][BoardY]bool
}

func NewBoard() *Board {
    return &Board{}
}

func (b *Board) Set(x, y int, state bool) {
    b.cells[x][y] = state
}

func (b *Board) Get(x, y int) bool {
    return b.cells[x][y]
}

func (b *Board) Neighbours(x, y int) int {
    // loop min and max if they are outside the bounds            
    min_y := y-1
    max_y := y+1
    if min_y < 0 {
        min_y = BoardY-1
    } else if max_y >= BoardY {
        max_y = 0
    }

    min_x := x-1
    max_x := x+1
    if min_x < 0 {
        min_x = BoardX-1
    } else if max_x >= BoardX {
        max_x = 0
    }

    // check each neighbor, incrementing count if alive
    var neighbours int

    if b.Get(min_x, min_y) {neighbours++}
    if b.Get(x, min_y) {neighbours++}
    if b.Get(max_x, min_y) {neighbours++}

    if b.Get(min_x, y) {neighbours++}
    if b.Get(max_x, y) {neighbours++}

    if b.Get(min_x, max_y) {neighbours++}
    if b.Get(x, max_y) {neighbours++}
    if b.Get(max_x, y) {neighbours++}

    return neighbours
}

func (b *Board) Update() {
    for x := 0; x < BoardX; x++ {
        for y := 0; y < BoardY; y++ {
            neighbours := b.Neighbours(x, y)
            if b.Get(x, y) { // check if the cell is alive
                if neighbours != 2 && neighbours != 3 {
                    b.Set(x, y, false)
                }
            } else if neighbours == 3 {
                b.Set(x, y, true)
            }
        }
    }
}

func (b *Board) Draw() {
    imd.Color = colornames.White
    for x := 0; x < BoardX; x++ {
        for y := 0; y < BoardY; y++ {
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
        //Maximized: true,
        Icon:      []pixel.Picture{icon},
        VSync: true,
    }
    win, err = pixelgl.NewWindow(cfg)
    if err != nil {
        panic(err)
    }

    imd = imdraw.New(nil)
    
    ticker := time.NewTicker(time.Second/2)
    board := NewBoard()

    for !win.Closed() {
        imd.Clear()

        if win.JustPressed(pixelgl.KeyEscape) {
            win.SetClosed(true)
        }

        if win.JustPressed(pixelgl.KeySpace) {
            running = !running
        }

        if running {
            select {
            case <-ticker.C:
                board.Update()
            default:
            }
        } else if win.JustPressed(pixelgl.MouseButtonLeft) {
            pos := win.MousePosition()
            cell_x, cell_y := int(pos.X/20), int(pos.Y/20)
            board.Set(cell_x, cell_y, !board.Get(cell_x, cell_y))
        }
            
        board.Draw()

        win.Clear(colornames.Black)
        imd.Draw(win)
        win.Update()
    }
}

func main() {
    pixelgl.Run(run)
}