package main

import (
	"github.com/jteeuwen/glfw"
	"github.com/banthar/gl"
	"github.com/banthar/glu"
	"log"
)

var (
	// used when the user press Esc to exit
	// buffered in one unit to avoid locking on channel send
	exit chan int = make(chan int,1)
)

func main() {
	log.Printf("Starting glfw window")

	err := glfw.Init()
	if err != nil {
		log.Fatalf("Error while starting glfw library: %v", err)
	}
	defer glfw.Terminate()

	err = glfw.OpenWindow(256,256, 8,8,8, 0, 0, 0, glfw.Windowed)
	if err != nil {
		log.Fatalf("Error while opening glfw window: %v", err)
	}
	defer glfw.CloseWindow()

	glfw.SetSwapInterval(1) //vsync on

	glfw.SetWindowTitle("Colored Triangle")

	InitGL()

	glfw.SetWindowSizeCallback(func(w,h int){ InitGLWindow(w,h) })
	glfw.SetKeyCallback(func(key,state int){ HandleKey(key,state) })

	run := true
	for run && glfw.WindowParam(glfw.Opened) == 1 {
		select {
			case exitCode := <-exit:
				log.Printf("Received exit code: %d", exitCode)
				run = false
			default:
				Draw()
		}
	}
}

func HandleKey(key, state int) {
	switch key {
		case glfw.KeyEsc:
			exit <- 0
	}
}

// Initialize the OpenGL configuration with the information from the current window.
func InitGLWindow(w,h int) {
	gl.Viewport(0,0,w,h)
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	glu.Perspective(45.0, float64(w)/float64(h), 0.1, 100.0)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
}

// Initialize the OpenGL configuration that don't depend on window parameters
func InitGL() {
	gl.ShadeModel(gl.SMOOTH)
	gl.ClearColor(0,0,0,0)
	gl.ClearDepth(1)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)
	gl.Hint(gl.PERSPECTIVE_CORRECTION_HINT, gl.NICEST)
}

// Render stuff
func Draw() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.LoadIdentity()

	gl.Translatef(-1.5, 0, -6)
	gl.Begin(gl.TRIANGLES)
	gl.Color3f(1, 0, 0)
	gl.Vertex3f(0, 1, 0)
	gl.Color3f(0, 1, 0)
	gl.Vertex3f(-1, -1, 0)
	gl.Color3f(0, 0, 1)
	gl.Vertex3f(1, -1, 0)
	gl.End()

	glfw.SwapBuffers()
}
