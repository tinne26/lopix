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
	xMargin float64
	yMargin float64
	
	needsRedraw bool
	redrawManaged bool
	inDraw bool

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
	self.inDraw = true
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
	if !prevDrawWasHiRes && (!self.redrawManaged || self.needsRedraw) {
		self.project(hiResCanvas)
		self.needsRedraw = false
	}

	// respect ebiten.IsScreenClearedEveryFrame()
	if ebiten.IsScreenClearedEveryFrame() {
		self.logicalCanvas.Clear()
	}
	self.inDraw = false
}

func (self *controller) Layout(logicWinWidth, logicWinHeight int) (int, int) {
	monitor := ebiten.Monitor()
	scale := monitor.DeviceScaleFactor()
	hiResWidth  := int(float64(logicWinWidth)*scale)
	hiResHeight := int(float64(logicWinHeight)*scale)
	if hiResWidth != self.hiResWidth || hiResHeight != self.hiResHeight {
		self.hiResWidth  = hiResWidth
		self.hiResHeight = hiResHeight
		self.redrawRequest()
		self.refreshMargins()
	}
	return self.hiResWidth, self.hiResHeight
}

func (self *controller) LayoutF(logicWinWidth, logicWinHeight float64) (float64, float64) {
	monitor := ebiten.Monitor()
	scale := monitor.DeviceScaleFactor()
	outWidth  := math.Ceil(logicWinWidth*scale)
	outHeight := math.Ceil(logicWinHeight*scale)
	if int(outWidth) != self.hiResWidth || int(outHeight) != self.hiResHeight {
		self.hiResWidth, self.hiResHeight = int(outWidth), int(outHeight)
		self.redrawRequest()
	}
	return outWidth, outHeight
}

func (self *controller) setResolution(width, height int) {
	if self.inDraw { panic("can't change resolution during draw stage") }
	if width < 1 || height < 1 { panic("Game resolution must be at least (1, 1)") }
	if width != self.logicalWidth || height != self.logicalHeight {
		rawWidth, rawHeight := intImgSize(self.rawLogicalCanvas)
		self.redrawRequest()
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
	if self.inDraw { panic("can't change scaling mode during draw stage") }
	if mode != self.scalingMode {
		self.scalingMode = mode
		self.redrawRequest()
		self.refreshMargins()
	}
}

func (self *controller) getScalingMode() ScalingMode {
	return self.scalingMode
}

func (self *controller) setScalingFilter(filter ScalingFilter) {
	if self.inDraw { panic("can't change scaling filter during draw stage") }
	if filter != self.scalingFilter {
		self.scalingFilter = filter
		self.redrawRequest()
	}
}

func (self *controller) getScalingFilter() ScalingFilter {
	return self.scalingFilter
}

func (self *controller) queueLogicalDraw(callback func(*ebiten.Image)) {
	if !self.inDraw { panic("can't queue draw outside draw stage") }
	self.queuedDraws = append(self.queuedDraws, queuedDraw{ callback, false })
}

func (self *controller) queueHiResDraw(callback func(*ebiten.Image)) {
	if !self.inDraw { panic("can't queue draw outside draw stage") }
	self.queuedDraws = append(self.queuedDraws, queuedDraw{ callback, true })
}

func (self *controller) redrawSetManaged(managed bool) {
	if self.inDraw { panic("can't modify RedrawManager during draw stage") }
	self.redrawManaged = managed
}

func (self *controller) redrawIsManaged() bool {
	return self.redrawManaged
}

func (self *controller) redrawRequest() {
	if self.inDraw { panic("can't modify RedrawManager during draw stage") }
	self.needsRedraw = true
}

func (self *controller) redrawPending() bool {
	return self.needsRedraw
}
