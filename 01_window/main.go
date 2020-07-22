package main

import (
	"runtime"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	windowWidth  = 800
	windowHeight = 600
	windowTitle  = "Learn OpenGL :: Window"
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func main() {
	// Initialize GLFW
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	defer glfw.Terminate()
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	// Create a window
	window, err := glfw.CreateWindow(windowWidth, windowHeight, windowTitle, nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	// Register a callback to adjust the viewport's size each time the window is resized
	window.SetFramebufferSizeCallback(framebufferSizeCallback)

	// Important! Call gl.Init() only under the presence of an active OpenGL context,
	// i.e., after MakeContextCurrent().
	if err := gl.Init(); err != nil {
		panic(err)
	}
	// Set viewport
	gl.Viewport(0, 0, windowWidth, windowHeight)

	// Render loop
	for !window.ShouldClose() {
		// Check whether the user pressed the ESC (read more about this function in its declaration).
		processInput(window)

		// Rendering commands here...

		// Just to test if things actually work we want to clear the screen with a color of our choice.
		// At the start of each render iteration we always want to clear the screen otherwise we would
		// still see the results from the previous iteration (this could be the effect you're looking
		// for, but usually you don't). We can clear the screen's color buffer using the gl.Clear()
		// function where we pass in buffer bits to specify which buffer we would like to clear. The
		// possible bits we can set are gl.COLOR_BUFFER_BIT, gl.DEPTH_BUFFER_BIT and
		// gl.STENCIL_BUFFER_BIT. Right now we only care about the color values so we only clear the
		// color buffer.
		// Note that we also set a color via gl.ClearColor() to clear the screen with. Whenever we call
		// gl.Clear() and clear the color buffer, the entire color buffer will be filled with the color as
		// configured by gl.ClearColor(). This will result in a dark green-blueish color.
		// As you might recall from the OpenGL tutorial, the gl.ClearColor() function is a state-setting
		// function and gl.Clear() is a state-using function in that it uses the current state to retrieve
		// the clearing color from.
		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		// Will swap the color buffer (read about "Double buffer")
		window.SwapBuffers()
		// Check if any events are triggered (like keyboard input or mouse movements)
		glfw.PollEvents()
	}
}

// This function is used as callback to adjust the viewport's size each time the window is resized.
func framebufferSizeCallback(w *glfw.Window, width int, height int) {
	gl.Viewport(0, 0, int32(width), int32(height))
}

// Check whether the user has pressed the escape key (if it's not pressed, w.GetKey() returns
// glfw.Release). If the user did press the escape key, we close GLFW by setting w.ShouldClose
// property to true. The next condition check of the main while (render) loop will then fail and
// the application closes.
func processInput(w *glfw.Window) {
	if w.GetKey(glfw.KeyEscape) == glfw.Press {
		w.SetShouldClose(true)
	}
}
