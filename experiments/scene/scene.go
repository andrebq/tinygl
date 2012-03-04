// see README.mkd for mor information
package main

import (
	"github.com/jteeuwen/glfw"
	"github.com/banthar/gl"
	"github.com/banthar/glu"
	"log"
	"fmt"
	"image/color"
)

// Represent a vector with color information
type Vector struct {
	X, Y, Z float32
	Color color.Color
}

// Return the string representation of the vector coords
func (v *Vector) Coords() string {
	return fmt.Sprintf("(%.2f,%.2f,%.2f)", v.X, v.Y, v.Z)
}

// Return the string representation of the vector color.
// If the color is a pre-defined return the name, otherwise return the color components
func (v *Vector) ColorName() string {
	if v.Color == color.Black {
		return "Black"
	} else if v.Color == color.White {
		return "White"
	} else if v.Color == color.Transparent {
		return "Transparent"
	} else if v.Color == color.Opaque {
		return "Opaque"
	} else {
		r,g,b,a := v.Color.RGBA()
		return fmt.Sprintf("[%d,%d,%d,%d]",r,g,b,a)
	}
	return ""
}

// Represents a collection of vectors with color information
type Mesh struct {
	V []*Vector
}

// Render the mesh
// The V array must be a multiple of 3 since it uses triangles to render the mesh.
// If two or more triangles share the same edge the user can pass the same pointer twice
func (m *Mesh) Render() {
	count := 0
	gl.Begin(gl.TRIANGLES)
	for _, v := range(m.V) {
		if count % 3 == 0 && count > 0 {
			gl.End()
			gl.Begin(gl.TRIANGLES)
			count = 0
		}
		gl.Color3f(1,0,0)
		gl.Vertex3f(v.X, v.Y, v.Z)
		count++
	}
	if count > 0 {
		gl.End()
	}
}

type Render interface {
	Render()
}

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

	// create a mesh with 3 vertices (a triangle)
	m := &Mesh{make([]*Vector,3)}
	m.V[0] = &Vector{0,1,0, color.NRGBA{1, 0, 0, 0}}
	m.V[1] = &Vector{-1,-1,0, color.NRGBA{0, 1, 0, 0}}
	m.V[2] = &Vector{1, -1, 0, color.NRGBA{0, 0, 1, 0}}

	run := true
	for run && glfw.WindowParam(glfw.Opened) == 1 {
		select {
			case exitCode := <-exit:
				log.Printf("Received exit code: %d", exitCode)
				run = false
			default:
				Draw(m)
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
func Draw(r Render) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.LoadIdentity()

	gl.Translatef(-1.5, 0, -6)
/*
	gl.Begin(gl.TRIANGLES)
	gl.Color3f(1, 0, 0)
	gl.Vertex3f(0, 1, 0)

	gl.Color3f(1, 0, 0)
	gl.Vertex3f(-1, -1, 0)

	gl.Color3f(1, 0, 0)
	gl.Vertex3f(1, -1, 0)

	gl.End()
*/
	r.Render()
	glfw.SwapBuffers()
}
