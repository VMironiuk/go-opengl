package main

import (
	"fmt"
	"runtime"
	"strings"

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

	// Build and compile shader program
	var status int32
	var logLength int32
	// Vertex shader
	vShader := gl.CreateShader(gl.VERTEX_SHADER)
	cVertexShaderSrc, free := gl.Strs(vertexShaderSrc)
	gl.ShaderSource(vShader, 1, cVertexShaderSrc, nil)
	free()
	gl.CompileShader(vShader)
	// Check vertex shader compilation status
	gl.GetShaderiv(vShader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		gl.GetShaderiv(vShader, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(vShader, logLength, nil, gl.Str(log))
		fmt.Printf("failed to compile vertex shader: %v", log)
	}
	// Fragment shader
	fShader := gl.CreateShader(gl.FRAGMENT_SHADER)
	cFragmentShaderSrc, free := gl.Strs(fragmentShaderSrc)
	gl.ShaderSource(fShader, 1, cFragmentShaderSrc, nil)
	free()
	gl.CompileShader(fShader)
	// Check fragment shader compilation status
	gl.GetShaderiv(fShader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		gl.GetShaderiv(fShader, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(fShader, logLength, nil, gl.Str(log))
		fmt.Printf("failed to compile fragment shader: %v", log)
	}
	// Shader program
	program := gl.CreateProgram()
	gl.AttachShader(program, vShader)
	gl.AttachShader(program, fShader)
	gl.LinkProgram(program)
	// Check program link status
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))
		fmt.Printf("failed to link program: %v", log)
	}

	gl.DeleteShader(vShader)
	gl.DeleteShader(fShader)

	// Setup vetex data (adn buffer(s)) and configure vertex attributes
	var vertices = []float32{
		0.5, 0.5, 0.0, // top right
		0.5, -0.5, 0.0, // bottom right
		-0.5, -0.5, 0.0, // bottom left
		-0.5, 0.5, 0.0, // top left
	}

	var indices = []int32{
		0, 1, 3, // first triangle
		1, 2, 3, // second triangle
	}

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	var ebo uint32
	gl.GenBuffers(1, &ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

	defer gl.DeleteVertexArrays(1, &vao)
	defer gl.DeleteBuffers(1, &vbo)

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

		// Render
		gl.UseProgram(program)
		gl.BindVertexArray(vao)
		// gl.DrawArrays(gl.TRIANGLES, 0, 3)
		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))

		// Will swap the color buffer (read about "Double buffer")
		window.SwapBuffers()
		// Check if any events are triggered (like keyboard input or mouse movements)
		glfw.PollEvents()
	}
}

var vertexShaderSrc = `
#version 330 core

layout (location=0) in vec3 pos;

void main() {
	gl_Position = vec4(pos.x, pos.y, pos.z, 1.0f);
}
` + "\x00"

var fragmentShaderSrc = `
#version 330 core

out vec4 fragColor;

void main() {
	fragColor = vec4(1.0f, 0.5f, 0.2f, 1.0f);
}
` + "\x00"

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
