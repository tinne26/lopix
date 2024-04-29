package lopix

import "github.com/hajimehoshi/ebiten/v2"

func (self *controller) project(hiResCanvas *ebiten.Image) {
	switch self.scalingMode {
	case Proportional:
		self.projectProportional(hiResCanvas)
	case PixelPerfect:
		panic("unimplemented") // skip filter or not here?
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
	factor := self.getProportionalScaleFactor(hiResCanvas)
	
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

func (self *controller) getProportionalScaleFactor(hiResCanvas *ebiten.Image) float64 {
	ofw, ofh := floatImgSize(hiResCanvas)
	ifw, ifh := float64(self.logicalWidth), float64(self.logicalHeight)
	return min(ofw/ifw, ofh/ifh)
}

func (self *controller) nearestDrawWithFactor(hiResCanvas *ebiten.Image, factor float64) {
	fw, fh := floatImgSize(hiResCanvas)
	self.opts.GeoM.Scale(factor, factor)
	xTrans := (fw - float64(self.logicalWidth )*factor)/2.0
	yTrans := (fh - float64(self.logicalHeight)*factor)/2.0
	self.opts.GeoM.Translate(xTrans, yTrans)
	hiResCanvas.DrawImage(self.logicalCanvas, &self.opts)
	self.opts.GeoM.Reset()
}

func (self *controller) linearDrawWithFactor(hiResCanvas *ebiten.Image, factor float64) {
	self.opts.Filter = ebiten.FilterLinear
	fw, fh := floatImgSize(hiResCanvas)
	self.opts.GeoM.Scale(factor, factor)
	xTrans := (fw - float64(self.logicalWidth )*factor)/2.0
	yTrans := (fh - float64(self.logicalHeight)*factor)/2.0
	self.opts.GeoM.Translate(xTrans, yTrans)
	hiResCanvas.DrawImage(self.logicalCanvas, &self.opts)
	self.opts.GeoM.Reset()
	self.opts.Filter = ebiten.FilterNearest
}
