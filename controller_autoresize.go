package lopix

import "github.com/hajimehoshi/ebiten/v2"

func (self *controller) autoResizeWindow() {
	if self.logicalWidth < 1 || self.logicalHeight < 1 {
		panic("can't auto-resize before setting the game resolution")
	}
	setMaxMultRawWindowSize(self.logicalWidth, self.logicalHeight, 128)
}

func setMaxMultRawWindowSize(width, height int, logicalMargin int) {
	scaledWidth, scaledHeight := findMaxMultRawWindowSize(width, height, logicalMargin)
	ebiten.SetWindowSize(scaledWidth, scaledHeight)
}

func findMaxMultRawWindowSize(width, height int, logicalMargin int) (int, int) {
	scale := ebiten.DeviceScaleFactor()
	fsWidth, fsHeight := ebiten.ScreenSizeInFullscreen()
	if fsWidth <= 0 || fsHeight <= 0 { // fallback
		fsWidth, fsHeight = 480, 480 // maybe even 640 should be fine
	}
	maxWidthMult  := (fsWidth  - logicalMargin)/width
	maxHeightMult := (fsHeight - logicalMargin)/height
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
