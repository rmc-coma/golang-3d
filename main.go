package main

import (
	"fmt"
	"runtime"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"golang.org/x/image/math/f32"
)

func renderTriangle() {
	const (
		windowWidth  = 960
		windowHeight = 540
	)

	runtime.LockOSThread()
	if err := glfw.Init(); err != nil {
		panic(fmt.Errorf("could not initialize glfw: %v", err))
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 5)
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.DoubleBuffer, glfw.True)
	glfw.WindowHint(glfw.DepthBits, 24)

	win, winErr := glfw.CreateWindow(windowWidth, windowHeight, "Hello world", nil, nil)
	if winErr != nil {
		panic(fmt.Errorf("could not create opengl renderer: %v", winErr))
	}

	win.MakeContextCurrent()
	if glErr := gl.Init(); glErr != nil {
		panic(glErr)
	}

	gl.ClearColor(0.0, 0.0, 0.0, 1.0)

	var vertices = [9]float32{-0.5, -0.5, 0.0, 0.0, 0.5, 0.0, 0.5, -0.5, 0.0}
	var vboId uint32
	gl.GenBuffers(1, &vboId)

	gl.BindBuffer(gl.ARRAY_BUFFER, vboId)
	gl.BufferData(gl.ARRAY_BUFFER, 36, gl.Ptr(&vertices[0]), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	var vaoId uint32
	gl.GenVertexArrays(1, &vaoId)
	gl.BindVertexArray(vaoId)
	gl.BindBuffer(gl.ARRAY_BUFFER, vboId)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, gl.Ptr(nil))
	gl.EnableVertexAttribArray(0)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)

	vertexShaderSource :=
		`#version 450 core
		layout (location = 0) in vec3 vertex;

		void main() {
			gl_Position = vec4(vertex, 1.0f);
		}`

	vertexShader := gl.CreateShader(gl.VERTEX_SHADER)
	csource, freeCsource := gl.Strs(vertexShaderSource)
	gl.ShaderSource(vertexShader, 1, csource, nil)
	freeCsource()
	gl.CompileShader(vertexShader)
	var status int32
	gl.GetShaderiv(vertexShader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		panic("Failed to compile vertex shader")
	}

	fragmentShaderSource :=
		`#version 450 core
		out vec4 fragmentColor;

		void main() {
			fragmentColor = vec4(0.5f, 0.5f, 0.5f, 1.0f);
		}`

	fragmentShader := gl.CreateShader(gl.FRAGMENT_SHADER)
	csource2, freeCsource2 := gl.Strs(fragmentShaderSource)
	gl.ShaderSource(fragmentShader, 1, csource2, nil)
	freeCsource2()
	gl.CompileShader(fragmentShader)
	gl.GetShaderiv(fragmentShader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		panic("Failed to compile fragment shader")
	}

	shaderProgram := gl.CreateProgram()
	gl.AttachShader(shaderProgram, vertexShader)
	gl.AttachShader(shaderProgram, fragmentShader)
	gl.LinkProgram(shaderProgram)
	gl.GetProgramiv(shaderProgram, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		panic("Failed to link program")
	}

	for !win.ShouldClose() && win.GetKey(glfw.KeyEscape) == 0 {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		gl.UseProgram(shaderProgram)
		gl.BindVertexArray(vaoId)
		gl.DrawArrays(gl.TRIANGLES, 0, 3)
		gl.BindVertexArray(0)
		gl.UseProgram(0)

		glfw.PollEvents()
		win.SwapBuffers()
	}

	win.Destroy()
}

func main() {
	var vec f32.Vec3 = PolarCoordinatesToDirectionVector(3.14159265359 / 4.0, 3.14159265359 / 4.0)
	fmt.Println(vec[0], vec[1], vec[2])
}
