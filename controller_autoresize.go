package lopix

import "github.com/hajimehoshi/ebiten/v2"

func (self *controller) autoResizeWindow() {
	if self.logicalWidth < 1 || self.logicalHeight < 1 {
		panic("can't auto-resize before setting the game resolution")
	}
	setMaxMultRawWindowSize(self.logicalWidth, self.logicalHeight, 0.37)
}

// TODO: margin proportion is not ideal. it's an okish desirable reference,
//       but we should still evaluate the overshooting case and take it
//       if the error for not taking it is significantly higher. we would
//       need to be way smarter...
func setMaxMultRawWindowSize(width, height int, marginProportion float64) {
	// there's an issue on ebitengine that affects this whole process,
	// see hajimehoshi/ebiten/issues/2978
	scaledWidth, scaledHeight := findMaxMultRawWindowSize(width, height, marginProportion)
	ebiten.SetWindowSize(scaledWidth, scaledHeight)
}

func findMaxMultRawWindowSize(width, height int, marginProportion float64) (int, int) {
	monitor := ebiten.Monitor()
	scale := monitor.DeviceScaleFactor()
	fsWidth, fsHeight := monitor.Size()
	if fsWidth <= 0 || fsHeight <= 0 { // fallback for mobile
		fsWidth, fsHeight = 480, 480 // maybe even 640 should be fine
	}
	fsWidth, fsHeight = int(float64(fsWidth)*scale), int(float64(fsHeight)*scale)
	margin := int(float64(min(fsWidth, fsHeight))*(marginProportion/scale))
	maxWidthMult  := (fsWidth  - margin)/width
	maxHeightMult := (fsHeight - margin)/height
	if maxWidthMult < maxHeightMult { maxHeightMult = maxWidthMult }
	if maxHeightMult < maxWidthMult { maxWidthMult = maxHeightMult }
	if maxWidthMult <= 0 || maxHeightMult <= 0 {
		maxWidthMult  = 1
		maxHeightMult = 1
	}

	width, height = width*maxWidthMult, height*maxHeightMult
	scaledWidth  := int(float64(width )/scale)
	scaledHeight := int(float64(height)/scale)
	return scaledWidth, scaledHeight
}
