package lopix

import "math"

import "github.com/hajimehoshi/ebiten/v2"

func (self *controller) project(hiResCanvas *ebiten.Image) {
	// TODO: could check that hiResCanvas bounds match self.hiResWidth
	//       and self.hiResHeight, but any measures in that case might
	//       do more harm than good (one frame desync is likeliest)
	switch self.scalingMode {
	case Proportional:
		self.projectProportional(hiResCanvas)
	case PixelPerfect:
		panic("unimplemented") // skip filter or not here? well, skip unless factors < 1.0
	case Stretched:
		panic("unimplemented")
	default:
		panic("invalid scaling mode")
	}
}

func floatImgSize(image *ebiten.Image) (float64, float64) {
	w, h := intImgSize(image)
	return float64(w), float64(h)
}

func intImgSize(image *ebiten.Image) (int, int) {
	if image == nil { return 0, 0 }
	bounds := image.Bounds()
	return bounds.Dx(), bounds.Dy()
}

func (self *controller) projectProportional(hiResCanvas *ebiten.Image) {
	factor := self.getProportionalScaleFactor()
	
	// pixel perfect case, no filter needed
	if factor == float64(int(factor)) {
		self.nearestDrawWithFactor(hiResCanvas, factor)
		return
	}

	switch self.scalingFilter {
	case Derivative:
		panic("unimplemented")
	case Bilinear:
		panic("unimplemented")
	case Nearest:
		self.nearestDrawWithFactor(hiResCanvas, factor)
	case Linear:
		self.linearDrawWithFactor(hiResCanvas, factor)
	default:
		panic("invalid scaling filter")
	}
}

func (self *controller) getProportionalScaleFactor() float64 {
	ofw, ofh := float64(self.hiResWidth), float64(self.hiResHeight)
	ifw, ifh := float64(self.logicalWidth), float64(self.logicalHeight)
	return min(ofw/ifw, ofh/ifh)
}

func (self *controller) getPixelPerfectScaleFactor() float64 {
	ofw, ofh := float64(self.hiResWidth), float64(self.hiResHeight)
	ifw, ifh := float64(self.logicalWidth), float64(self.logicalHeight)
	factor := min(ofw/ifw, ofh/ifh)
	perfectFactor := math.Floor(factor)
	if perfectFactor == 0 { return factor }
	return perfectFactor
}

func (self *controller) getStretchedFactors() (float64, float64) {
	ofw, ofh := float64(self.hiResWidth), float64(self.hiResHeight)
	ifw, ifh := float64(self.logicalWidth), float64(self.logicalHeight)
	return ofw/ifw, ofh/ifh
}

func (self *controller) nearestDrawWithFactor(hiResCanvas *ebiten.Image, factor float64) {
	self.opts.GeoM.Scale(factor, factor)
	self.opts.GeoM.Translate(self.xMargin, self.yMargin)
	hiResCanvas.DrawImage(self.logicalCanvas, &self.opts)
	self.opts.GeoM.Reset()
}

func (self *controller) linearDrawWithFactor(hiResCanvas *ebiten.Image, factor float64) {
	self.opts.Filter = ebiten.FilterLinear
	self.opts.GeoM.Scale(factor, factor)
	self.opts.GeoM.Translate(self.xMargin, self.yMargin)
	hiResCanvas.DrawImage(self.logicalCanvas, &self.opts)
	self.opts.GeoM.Reset()
	self.opts.Filter = ebiten.FilterNearest
}

func (self *controller) refreshMargins() {
	var xFactor, yFactor float64
	switch self.scalingMode {
	case Proportional:
		xFactor = self.getProportionalScaleFactor()
		yFactor = xFactor
	case PixelPerfect:
		xFactor = self.getPixelPerfectScaleFactor()
		yFactor = xFactor
	case Stretched:
		xFactor, yFactor = self.getStretchedFactors()
	default:
		panic("invalid scaling mode")
	}

	self.xMargin = (float64(self.hiResWidth ) - float64(self.logicalWidth )*xFactor)/2.0
	self.yMargin = (float64(self.hiResHeight) - float64(self.logicalHeight)*yFactor)/2.0
}
