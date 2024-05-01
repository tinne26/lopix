package lopix

import "cmp"

func clamp[T cmp.Ordered](x, a, b T) T {
	if x < a { return a }
	if x > b { return b }
	return x
}

func (self *controller) toRelativeCoords(x, y int) (float64, float64) {
	relX := (float64(x) - self.xMargin)/(float64(self.hiResWidth ) - self.xMargin*2)
	relY := (float64(y) - self.yMargin)/(float64(self.hiResHeight) - self.yMargin*2)
	return clamp(relX, 0.0, 1.0), clamp(relY, 0.0, 1.0)
}

func (self *controller) toLogicalCoords(x, y int) (int, int) {
	rx, ry := self.toRelativeCoords(x, y)
	return int(rx*float64(self.logicalWidth)), int(ry*float64(self.logicalHeight))
}
