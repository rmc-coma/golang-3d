package main

import (
	"math"

	"golang.org/x/image/math/f32"
)

func PolarCoordinatesToDirectionVector(phi float32, theta float32) f32.Vec3 {
	var cos_theta = math.Cos(float64(theta));
	return f32.Vec3{float32(cos_theta * math.Sin(float64(phi))), float32(math.Sin(float64(theta))), float32(cos_theta * math.Cos(float64(phi)))}
}

type Coords [3]int64

type Ray struct {
	position f32.Vec3
	direction f32.Vec3
}

func RaycastVoxel(ray Ray, voxels [][][]uint32, voxelSize int64) Coords {
	var currentVoxel = Coords{int64(ray.position[0]), int64(ray.position[1]), int64(ray.position[2])}
	var step = 1
	var invDir = f32.Vec3{1.0 / ray.direction[0], 1.0 / ray.direction[1], 1.0 / ray.direction[2]}
	var T = f32.Vec3{0.0, 0.0, 0.0}
	var DeltaT = f32.Vec3{float32(voxelSize) * invDir[0], float32(voxelSize) * invDir[1], float32(voxelSize) * invDir[2]}
	var increment = [3]bool{false, false, false}
	const traversalLimit = 100000
	for currentVoxel[0] < traversalLimit && currentVoxel[1] < traversalLimit && currentVoxel[2] < traversalLimit {
		increment[0] = T[0] <= T[1] && T[0] <= T[2]
		increment[1] = T[1] <= T[0] && T[1] <= T[2]
		increment[2] = T[2] <= T[0] && T[2] <= T[1]
		if increment[0] { T[0] += DeltaT[0] }
		if increment[1] { T[1] += DeltaT[1] }
		if increment[2] { T[2] += DeltaT[2] }
		if increment[0] { currentVoxel[0] += int64(step) }
		if increment[1] { currentVoxel[1] += int64(step) }
		if increment[2] { currentVoxel[2] += int64(step) }
		if voxels[currentVoxel[2]][currentVoxel[1]][currentVoxel[0]] != 0 {
			return currentVoxel
		}
	}
	return Coords{traversalLimit, traversalLimit, traversalLimit}
}

func CreateRaysGrid(lookAt f32.Vec2, screenDimensions [2]uint32, fovAngle float32) [][]f32.Vec3 {
	var fov = f32.Vec2{fovAngle, fovAngle / (float32(screenDimensions[0]) / float32(screenDimensions[1]))}
	var fovPerPixel = f32.Vec2{fov[0] / float32(screenDimensions[0]), fov[1] / float32(screenDimensions[1])}
	var halfScreenFov = f32.Vec2{(float32(screenDimensions[0]) / 2.0 - 0.5) * fovPerPixel[0], (float32(screenDimensions[1]) / 2.0 - 0.5) * fovPerPixel[1]}
	var margin = f32.Vec2{halfScreenFov[0] + lookAt[0], halfScreenFov[1] + lookAt[1]}
	var rays = make([][]f32.Vec3, screenDimensions[1])
	for j := uint32(0); j < screenDimensions[1]; j++ {
		rays[j] = make([]f32.Vec3,screenDimensions[0])
		for i := uint32(0); i < screenDimensions[0]; i++ {
			rays[j][i] = PolarCoordinatesToDirectionVector(fovPerPixel[0] * float32(i) - margin[0], fovPerPixel[1] * float32(j) - margin[1])
		}
	}
	return rays
}
