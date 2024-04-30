package lopix

import "image"
import "math"

import "github.com/hajimehoshi/ebiten/v2"

var pkgController controller

type controller struct {
	innerGame

	scalingMode ScalingMode
	scalingFilter ScalingFilter
	logicalWidth int
	logicalHeight int
	hiResWidth  int
	hiResHeight int

	rawLogicalCanvas *ebiten.Image
	logicalCanvas *ebiten.Image
	queuedDraws []queuedDraw

	opts ebiten.DrawImageOptions
}

type innerGame = Game
type queuedDraw struct {
	callback func(*ebiten.Image)
	isHighResolution bool
}

func (self *controller) Draw(hiResCanvas *ebiten.Image) {
	self.innerGame.Draw(self.logicalCanvas)
	
	var drawIndex int = 0
	var prevDrawWasHiRes bool = false
	for drawIndex < len(self.queuedDraws) {
		queued := self.queuedDraws[drawIndex]
		if queued.isHighResolution {
			prevDrawWasHiRes = true
			queued.callback(hiResCanvas)
		} else {
			if prevDrawWasHiRes {
				self.project(hiResCanvas)
				self.logicalCanvas.Clear()
				prevDrawWasHiRes = false
			}
			queued.callback(self.logicalCanvas)
		}
		
		drawIndex += 1
	}
	self.queuedDraws = self.queuedDraws[ : 0]

	// final projection
	if !prevDrawWasHiRes {
		self.project(hiResCanvas)
	}

	// respect ebiten.IsScreenClearedEveryFrame()
	if ebiten.IsScreenClearedEveryFrame() {
		self.logicalCanvas.Clear()
	}
}

func (self *controller) Layout(logicWinWidth, logicWinHeight int) (int, int) {
	monitor := ebiten.Monitor()
	scale := monitor.DeviceScaleFactor()
	self.hiResWidth  = int(float64(logicWinWidth)*scale)
	self.hiResHeight = int(float64(logicWinHeight)*scale)
	return self.hiResWidth, self.hiResHeight
}

func (self *controller) LayoutF(logicWinWidth, logicWinHeight float64) (float64, float64) {
	monitor := ebiten.Monitor()
	scale := monitor.DeviceScaleFactor()
	outWidth  := math.Ceil(logicWinWidth*scale)
	outHeight := math.Ceil(logicWinHeight*scale)
	self.hiResWidth, self.hiResHeight = int(outWidth), int(outHeight)
	return outWidth, outHeight
}

func (self *controller) setResolution(width, height int) {
	if width < 1 || height < 1 { panic("Game resolution must be at least (1, 1)") }
	if width != self.logicalWidth || height != self.logicalHeight {
		rawWidth, rawHeight := intImgSize(self.rawLogicalCanvas)
		if width <= rawWidth && height <= rawHeight {
			rect := image.Rect(0, 0, width, height)
			self.logicalCanvas = self.rawLogicalCanvas.SubImage(rect).(*ebiten.Image)
			self.logicalCanvas.Clear()
		} else {
			self.rawLogicalCanvas = ebiten.NewImage(width, height)
			self.logicalCanvas = self.rawLogicalCanvas
		}
		self.logicalWidth, self.logicalHeight = width, height
	}
}

func (self *controller) getResolution() (int, int) {
	return self.logicalWidth, self.logicalHeight
}

func (self *controller) setScalingMode(mode ScalingMode) {
	self.scalingMode = mode
}

func (self *controller) getScalingMode() ScalingMode {
	return self.scalingMode
}

func (self *controller) setScalingFilter(filter ScalingFilter) {
	self.scalingFilter = filter
}

func (self *controller) getScalingFilter() ScalingFilter {
	return self.scalingFilter
}

func (self *controller) queueLogicalDraw(callback func(*ebiten.Image)) {

}

func (self *controller) queueHiResDraw(callback func(*ebiten.Image)) {

}

