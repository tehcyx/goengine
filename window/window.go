package window

import (
	"fmt"
	"os"

	"github.com/go-gl/gl/v2.1/gl"
	"gopkg.in/veandco/go-sdl2.v0/sdl"
)

type Window struct {
	window   *sdl.Window
	context  sdl.GLContext
	event    sdl.Event
	isClosed bool
}

// NewWindow Creates a new window and returns a struct with the necessary accessors to handle the window
func NewWindow(winHeight, winWidth int32, title string) *Window {
	w := new(Window)
	w.create(winHeight, winWidth, title)
	return w
}

func (w *Window) create(winHeight, winWidth int32, title string) {
	var err error
	w.isClosed = false
	if err = sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_CORE)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 3)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 3)

	sdl.GLSetAttribute(sdl.GL_RED_SIZE, 8)
	sdl.GLSetAttribute(sdl.GL_GREEN_SIZE, 8)
	sdl.GLSetAttribute(sdl.GL_BLUE_SIZE, 8)
	sdl.GLSetAttribute(sdl.GL_ALPHA_SIZE, 8)
	sdl.GLSetAttribute(sdl.GL_BUFFER_SIZE, 32)
	sdl.GLSetAttribute(sdl.GL_DOUBLEBUFFER, 1)

	w.window, err = sdl.CreateWindow(title, sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, winWidth, winHeight, sdl.WINDOW_OPENGL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create window: %s\n", err)
		panic(err)
	}
	defer w.window.Destroy()

	w.context, err = w.window.GLCreateContext()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create context: %s\n", err)
		panic(err)
	}
	defer sdl.GLDeleteContext(w.context)

	// Initialize gl
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)
}

// Clear window
func (w *Window) Clear(red, green, blue, alpha float32) {
	gl.ClearColor(red, green, blue, alpha)
}

// Update window content
func (w *Window) Update() {
	// sdl.GL_SwapWindow(w.window)
	w.window.GLSwap() // is that right???

	for w.event = sdl.PollEvent(); w.event != nil; w.event = sdl.PollEvent() {
		switch t := w.event.(type) {
		case *sdl.QuitEvent:
			w.isClosed = false
		case *sdl.KeyboardEvent:
			if t.Keysym.Sym == sdl.K_ESCAPE {
				w.isClosed = false
			}
		case *sdl.MouseButtonEvent:
			if t.Type == sdl.MOUSEBUTTONUP {
				if t.Button == sdl.BUTTON_LEFT {
					fmt.Printf("Left Mouse %d\n", 1)
				} else if t.Button == sdl.BUTTON_RIGHT {
					fmt.Printf("Right Mouse %d\n", 1)
				}
			}
		case *sdl.MouseMotionEvent:
			// var zrot float32
			// var xrot float32
			// var yrot float32

			// xrot = float32(t.Y) / 2
			// yrot = float32(t.X) / 2
			//fmt.Printf("[%dms]MouseMotion \tid:%d \tx:%d \ty:%d \txrel:%d \tyrel:%d\n", t.Timestamp, t.Which, t.X, t.Y, t.XRel, t.YRel)
		}
	}
}

// IsClosed check if window is closed
func (w *Window) IsClosed() bool {
	return w.isClosed
}
