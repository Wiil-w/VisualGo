package base

type Vector2 struct {
	X, Y float64
}

func (v Vector2) Right(args ...float64) Vector2 {
	vNew := v
	vNew.X += v.validateArgs(args...)
	return vNew
}
func (v Vector2) Top(args ...float64) Vector2 {
	vNew := v
	vNew.Y += v.validateArgs(args...)
	return vNew
}
func (v Vector2) Left(args ...float64) Vector2 {
	vNew := v
	vNew.X -= v.validateArgs(args...)
	return vNew
}
func (v Vector2) Bottom(args ...float64) Vector2 {
	vNew := v
	vNew.Y -= v.validateArgs(args...)
	return vNew
}
func (v Vector2) validateArgs(args ...float64) float64 {
	if len(args) > 0 && args[0] > 0 {
		return args[0]
	} else {
		return 1.0
	}
}
