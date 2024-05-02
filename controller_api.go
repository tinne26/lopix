package lopix

import "image"
import "math"

import "github.com/hajimehoshi/ebiten/v2"

var pkgController controller

type controller struct {
	innerGame

	scalingMode ScalingMode
	scalingFilter ScalingFilter
	internalScalingFilter ScalingFilter
	logicalWidth int
	logicalHeight int
	hiResWidth  int
	hiResHeight int
	
	xMargin float64
	yMargin float64
	xFactor float64
	yFactor float64
	
	needsRedraw bool
	needsClear bool
	redrawManaged bool
	inDraw bool

	rawLogicalCanvas *ebiten.Image
	logicalCanvas *ebiten.Image
	queuedDraws []queuedDraw
	
	shaderOpts ebiten.DrawTrianglesShaderOptions
	shaderVertices [4]ebiten.Vertex
	shaders [scalingFilterEndSentinel]*ebiten.Shader
}

type innerGame = Game
type queuedDraw struct {
	callback func(*ebiten.Image)
	isHighResolution bool
}

func (self *controller) Draw(hiResCanvas *ebiten.Image) {
	self.inDraw = true
	self.considerDrawClear(hiResCanvas)
	self.innerGame.Draw(self.logicalCanvas)
	
	var drawIndex int = 0
	var prevDrawWasHiRes bool = false
	for drawIndex < len(self.queuedDraws) {
		queued := self.queuedDraws[drawIndex]
		if queued.isHighResolution {
			if !prevDrawWasHiRes {
				self.project(hiResCanvas)
				self.logicalCanvas.Clear()
			}
			prevDrawWasHiRes = true
			queued.callback(hiResCanvas)
		} else {
			prevDrawWasHiRes = false
			queued.callback(self.logicalCanvas)
		}
		
		drawIndex += 1
	}
	self.queuedDraws = self.queuedDraws[ : 0]

	// final projection
	if !prevDrawWasHiRes && (!self.redrawManaged || self.needsRedraw) {
		self.project(hiResCanvas)
	}
	self.needsRedraw = false

	// respect ebiten.IsScreenClearedEveryFrame()
	if ebiten.IsScreenClearedEveryFrame() {
		self.logicalCanvas.Clear()
	}
	self.inDraw = false
}

func (self *controller) considerDrawClear(hiResCanvas *ebiten.Image) {
	if !self.needsClear { return }
	self.needsClear = false
	hiResCanvas.Clear()
	self.logicalCanvas.Clear()
}

func (self *controller) Layout(logicWinWidth, logicWinHeight int) (int, int) {
	monitor := ebiten.Monitor()
	scale := monitor.DeviceScaleFactor()
	hiResWidth  := int(float64(logicWinWidth)*scale)
	hiResHeight := int(float64(logicWinHeight)*scale)
	if hiResWidth != self.hiResWidth || hiResHeight != self.hiResHeight {
		self.hiResWidth, self.hiResHeight = hiResWidth, hiResHeight
		self.redrawRequest()
		self.notifyCanvasChange()
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
		self.notifyCanvasChange()
	}
	return outWidth, outHeight
}

// func (self *controller) DrawFinalScreen(screen ebiten.FinalScreen, canvas *ebiten.Image, _ ebiten.GeoM) {
// 	screen.DrawImage(canvas, nil)
// }

func (self *controller) setResolution(width, height int) {
	if self.inDraw { panic("can't change resolution during draw stage") }
	if width < 1 || height < 1 { panic("Game resolution must be at least (1, 1)") }
	if width != self.logicalWidth || height != self.logicalHeight {
		var rawWidth, rawHeight int
		if self.rawLogicalCanvas != nil {
			bounds := self.rawLogicalCanvas.Bounds()
			rawWidth, rawHeight = bounds.Dx(), bounds.Dy()
		}
		self.redrawRequest()
		self.notifyCanvasChange()
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
		self.notifyCanvasChange()
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
		self.notifyCanvasChange()
	}
	
	if self.shaders[filter] == nil {
		self.compileShader(filter)
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

func (self *controller) hiResActiveArea() image.Rectangle {
	xm, ym := int(self.xMargin), int(self.yMargin)
	return image.Rectangle{
		Min: image.Pt(xm, ym),
		Max: image.Pt(self.hiResWidth - xm, self.hiResHeight - ym),
	}
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

func (self *controller) redrawScheduleClear() {
	self.needsClear = true
}

// --- misc helpers ---

func (self *controller) notifyCanvasChange() {
	// compute scaling factors
	switch self.scalingMode {
	case Proportional:
		self.xFactor = self.getProportionalScaleFactor()
		self.yFactor = self.xFactor
	case PixelPerfect:
		self.xFactor = self.getPixelPerfectScaleFactor()
		self.yFactor = self.xFactor
	case Stretched:
		self.xFactor, self.yFactor = self.getStretchedScaleFactors()
	default:
		panic("invalid scaling mode")
	}

	// update the internal scaling filter,
	// which simplifies to Nearest in some cases
	if self.xFactor == self.yFactor && self.xFactor == float64(int(self.xFactor)) {
		//self.internalScalingFilter = Nearest // TODO: restore optimizaiton when everything is stable
		self.internalScalingFilter = self.scalingFilter
	} else {
		self.internalScalingFilter = self.scalingFilter
	}

	self.xMargin = (float64(self.hiResWidth ) - float64(self.logicalWidth )*self.xFactor)/2.0
	self.yMargin = (float64(self.hiResHeight) - float64(self.logicalHeight)*self.yFactor)/2.0

	self.shaderOpts.Images[0] = self.logicalCanvas
	if self.shaderOpts.Uniforms != nil {
		self.shaderOpts.Uniforms["OutWidth"]  = float32(self.hiResWidth ) - float32(self.xMargin*2)
		self.shaderOpts.Uniforms["OutHeight"] = float32(self.hiResHeight) - float32(self.yMargin*2)
	}

	// set shader vertex positions, clockwise order, starting at top left
	fxMargin, fyMargin := float32(self.xMargin), float32(self.yMargin)
	// self.shaderVertices[0].SrcX = float32(0)
	// self.shaderVertices[0].SrcY = float32(0)
	self.shaderVertices[0].DstX = fxMargin
	self.shaderVertices[0].DstY = fyMargin

	self.shaderVertices[1].SrcX = float32(self.logicalWidth)
	//self.shaderVertices[1].SrcY = float32(0)
	self.shaderVertices[1].DstX = float32(self.hiResWidth) - fxMargin
	self.shaderVertices[1].DstY = fyMargin

	self.shaderVertices[2].SrcX = self.shaderVertices[1].SrcX
	self.shaderVertices[2].SrcY = float32(self.logicalHeight)
	self.shaderVertices[2].DstX = self.shaderVertices[1].DstX
	self.shaderVertices[2].DstY = float32(self.hiResHeight) - fyMargin

	//self.shaderVertices[3].SrcX = float32(0)
	self.shaderVertices[3].SrcY = self.shaderVertices[2].SrcY
	self.shaderVertices[3].DstX = fxMargin
	self.shaderVertices[3].DstY = self.shaderVertices[2].DstY
}
