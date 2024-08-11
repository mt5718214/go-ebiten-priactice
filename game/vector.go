package game

import "math"


type Vector struct {
	X float64
	Y float64
}

func (v *Vector) Normalize() Vector {
	// 計算向量長度
	length := math.Sqrt(v.X*v.X + v.Y*v.Y)
	if length == 0 {
		return Vector{0, 0} // Avoid division by zero
	}
	return Vector{
		X: v.X / length,
		Y: v.Y / length,
	}
}