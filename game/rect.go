package game


type Rect struct {
	X float64
	Y float64
	W float64
	H float64
}

func NewRect(x, y, w, h float64) Rect {
	return Rect{
		X: x, // image的x座標
		Y: y, // image的y座標
		W: w, // image的寬度
		H: h, // image的高度
	}
}

func (r Rect) MaxX() float64 {
	return r.X + r.W
}

func (r Rect) MaxY() float64 {
	return r.Y + r.H
}

// 計算兩個矩形是否有交集
// 解釋可看Readme的範例
func (r Rect) Intersects(other Rect) bool {
	return r.MaxX() >= other.X && r.X <= other.MaxX() && r.MaxY() >= other.Y && r.Y <= other.MaxY()
}