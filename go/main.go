package main

import "C"
import (
	"github.com/gen2brain/raylib-go/raylib"
)

var ()

func main() {
	// Shader Code
	//--------------------------------------------------------------------------------------
	vertCode := `
		#version 330

		// Input vertex attributes
		in vec3 vertexPosition;
		in vec4 vertexColor;

		// Input uniform values
		uniform mat4 mvp;

		// Output vertex attributes (to geometry shader)
		out vec4 geomColor;

		void main() {
		    // Send vertex attributes to fragment shader
		    geomColor = vertexColor;

		    // Calculate final vertex position
		    gl_Position = mvp*vec4(vertexPosition, 1.0);
		}`
	geomCode := `
		#version 330

		layout (triangles) in;
		layout (triangle_strip) out;
		layout (max_vertices = 256) out;

		in vec4 geomColor[];

		out vec4 fragColor;

		// Input uniform values
		uniform mat4 mvp;
		uniform mat4 matModel;
		uniform mat4 matView;
		uniform mat4 matProjection;

		void main() {
		    int i;
	    	int j;

		    for (i = 0; i+2 < gl_in.length(); i+=3) {
				for (j = 0; j < 64; j++) {
		            vec3 randomOffset = vec3(
		                fract(sin(dot(vec3(j, i, 0), vec3(12.9898, 78.233, 54.53))) * 43758.5453),
		                0.0,
		                fract(sin(dot(vec3(j + 1, i + 1, 0), vec3(12.9898, 78.233, 54.53))) * 43758.5453)
		            );
	    			vec4 grassPos = mix(mix(gl_in[i].gl_Position, gl_in[i+1].gl_Position, randomOffset.x), gl_in[i+2].gl_Position, randomOffset.z);

	    			gl_Position = grassPos + matProjection * (matView * matModel * vec4(0,0,0,0) + vec4(0.02, 0, 0, 0));
   					fragColor = vec4(0,0.4,0,1);
	    			EmitVertex();
	    			gl_Position = grassPos + mvp * vec4(0, 0.2, 0, 0) +  + matProjection * (matView * matModel * vec4(0,0,0,0) + vec4(0, 0.1, 0, 0));
   					fragColor = vec4(0,0.5,0,1);
	    			EmitVertex();
	    			gl_Position = grassPos + matProjection * (matView * matModel * vec4(0,0,0,0) + vec4(-0.02, 0, 0, 0));
   					fragColor = vec4(0,0.4,0,1);
	    			EmitVertex();
  						EndPrimitive();
				}
			}
		}`
	fragCode := `
	    #version 330

	    in vec4 fragColor;

	    out vec4 finalColor;

	    void main() {
	        finalColor = fragColor;
	    }`

	// Initialization
	//--------------------------------------------------------------------------------------
	screenWidth := int32(800)
	screenHeight := int32(450)

	rl.InitWindow(screenWidth, screenHeight, "Geom raylib test")

	// Define the camera to look into our 3d world
	camera := rl.Camera{}
	camera.Position = rl.NewVector3(0.0, 1.0, 3.0)
	camera.Target = rl.NewVector3(0.0, 0.0, 0.0)
	camera.Up = rl.NewVector3(0.0, 1.0, 0.0)
	camera.Fovy = 45.0
	camera.Projection = rl.CameraPerspective

	obj := rl.LoadModel("square.obj") // Load OBJ model
	shader := rl.LoadShaderFromMemory(vertCode, "", "", geomCode, fragCode)
	obj.Materials.Shader = shader

	position := rl.NewVector3(0.0, 0.0, 0.0) // Set model position
	//----------------------------------------------------------------------------------
	// Set render cycle to 60 fps
	rl.SetTargetFPS(60)
	for !rl.WindowShouldClose() {
		// Draw
		//----------------------------------------------------------------------------------
		rl.UpdateCamera(&camera, rl.CameraFirstPerson)
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		rl.BeginMode3D(camera)
		rl.DrawModel(obj, position, 1.0, rl.White) // Draw square
		rl.DrawGrid(20, 10.0)                      // Draw a grid
		rl.EndMode3D()

		rl.EndDrawing()
		//----------------------------------------------------------------------------------
	}

	// De-Initialization
	//--------------------------------------------------------------------------------------
	rl.UnloadModel(obj)     // Unload model
	rl.UnloadShader(shader) // Unload shader

	rl.CloseWindow()
	//--------------------------------------------------------------------------------------
}
