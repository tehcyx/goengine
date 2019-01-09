package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/andrebq/assimp/conv"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/tehcyx/goengine/mesh"
	"github.com/tehcyx/goengine/util"
	"gopkg.in/veandco/go-sdl2.v0/sdl"
)

func newProgram(vertexShaderSource, fragmentShaderSource string) (uint32, error) {
	defer util.TimeTrack(time.Now(), "newProgram")
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	program := gl.CreateProgram()

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program, nil
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	defer util.TimeTrack(time.Now(), "compileShader")
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

var uniRoll float32
var uniYaw float32
var uniPitch float32
var uniscale float32 = 0.3
var yrot float32 = 20.0
var zrot float32
var xrot float32

func main() {

	// printBanner()

	srcFilepath := "res/models/monkey.obj"

	var window *sdl.Window
	var context sdl.GLContext
	var event sdl.Event
	var running bool
	var err error
	runtime.LockOSThread()
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

	window, err = sdl.CreateWindow(winTitle, sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		winWidth, winHeight, sdl.WINDOW_OPENGL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create window: %s\n", err)
		panic(err)
	}
	defer window.Destroy()
	context, err = window.GLCreateContext()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create context: %s\n", err)
		panic(err)
	}
	defer sdl.GLDeleteContext(context)

	// Initialize gl
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	// Configure the vertex and fragment shaders
	program, err := newProgram(vertexShader, fragmentShader)
	if err != nil {
		panic(err)
	}

	gl.UseProgram(program)

	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(winWidth)/winHeight, 0.1, 10.0)
	projectionUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	camera := mgl32.LookAtV(mgl32.Vec3{3, 3, 3}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	cameraUniform := gl.GetUniformLocation(program, gl.Str("camera\x00"))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

	model := mgl32.Ident4()
	modelUniform := gl.GetUniformLocation(program, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	textureUniform := gl.GetUniformLocation(program, gl.Str("tex\x00"))
	gl.Uniform1i(textureUniform, 0)

	gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00"))

	monkeyModel := mesh.NewMeshFromFile(srcFilepath)
	// monkeyModel := mesh.NewMesh("res/models/monkey.obj")

	scene, err := conv.LoadAsset(srcFilepath)
	if err != nil {
		panic(err)
	}
	scene.Mesh[0].Id()

	// Configure global settings
	// gl.Enable(gl.DEPTH_TEST)
	// gl.DepthFunc(gl.LESS)

	// gl.Enable(gl.CULL_FACE)

	gl.ClearColor(0.0, 1.0, 0.8, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	running = true
	for running {
		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.KeyboardEvent:
				if t.Keysym.Sym == sdl.K_ESCAPE {
					running = false
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

				xrot = float32(t.Y) / 2
				yrot = float32(t.X) / 2
				//fmt.Printf("[%dms]MouseMotion \tid:%d \tx:%d \ty:%d \txrel:%d \tyrel:%d\n", t.Timestamp, t.Which, t.X, t.Y, t.XRel, t.YRel)
			}
		}

		gl.ClearColor(0.0, 1.0, 0.8, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		gl.UseProgram(program)
		gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

		monkeyModel.Draw()

		window.GLSwap()
	}
}

const (
	winTitle  = "OpenGL Shader"
	winWidth  = 800
	winHeight = 600
)

var vertexShader = `
#version 330
uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;
in vec3 vert;
in vec2 vertTexCoord;
void main() {
    gl_Position = projection * camera * model * vec4(vert, 1);
}
` + "\x00"

var fragmentShader = `
#version 330
uniform sampler2D tex;
out vec4 outputColor;
void main() {
    outputColor = vec4(0.3, 0.5, 0.8, 1.0);
}
` + "\x00"

func printBanner() {
	fmt.Println()
	fmt.Printf(`
	Welcome to
                                _         
            ___ ___ ___ ___ ___|_|___ ___ 
           | . | . | -_|   | . | |   | -_|
           |_  |___|___|_|_|_  |_|_|_|___|
           |___|           |___|          `)
	fmt.Println()
	fmt.Println()
	fmt.Println()
}
