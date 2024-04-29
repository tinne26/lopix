package lopix

func (self *controller) toRelativeCoords(x, y int) (float64, float64) {
	relX := float64(x)/float64(self.hiResWidth)
	relY := float64(y)/float64(self.hiResHeight)
	return max(min(relX, 1.0), 0.0), max(min(relY, 1.0), 0.0)
}

func (self *controller) toLogicalCoords(x, y int) (int, int) {
	rx, ry := self.toRelativeCoords(x, y)
	return int(rx*float64(self.logicalWidth)), int(ry*float64(self.logicalHeight))
}
