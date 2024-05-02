package lopix

import "math"

import "github.com/hajimehoshi/ebiten/v2"

var pkgShaderIndices = []uint16{0, 1, 3, 3, 1, 2}
func (self *controller) project(hiResCanvas *ebiten.Image) {
	// TODO: could check that hiResCanvas bounds match self.hiResWidth
	//       and self.hiResHeight, but any actions in that case might
	//       do more harm than good (one frame desync is likeliest)

	if self.shaders[self.internalScalingFilter] == nil {
		self.compileShader(self.internalScalingFilter)
	}
	hiResCanvas.DrawTrianglesShader(
		self.shaderVertices[:],
		pkgShaderIndices,
		self.shaders[self.internalScalingFilter],
		&self.shaderOpts,
	)
}

// --- get scale factors for each scaling mode ---

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

func (self *controller) getStretchedScaleFactors() (float64, float64) {
	ofw, ofh := float64(self.hiResWidth), float64(self.hiResHeight)
	ifw, ifh := float64(self.logicalWidth), float64(self.logicalHeight)
	return ofw/ifw, ofh/ifh
}

// --- shader properties setup ---

func (self *controller) compileShader(filter ScalingFilter) {
	var err error
	self.shaders[filter], err = ebiten.NewShader(pkgSrcKageFilters[filter])
	if err != nil {
		panic("Failed to compile shader for filter '" + filter.String() + "':\n" + err.Error())
	}
	if self.shaderOpts.Uniforms == nil {
		self.initShaderProperties()
	}
}

func (self *controller) initShaderProperties() {
	self.shaderOpts.Uniforms = make(map[string]any, 2)
	self.shaderOpts.Uniforms["OutWidth"]  = float32(self.hiResWidth ) - float32(self.xMargin*2)
	self.shaderOpts.Uniforms["OutHeight"] = float32(self.hiResHeight) - float32(self.yMargin*2)
	for i := range 4 {
		self.shaderVertices[i].ColorR = 1.0
		self.shaderVertices[i].ColorG = 1.0
		self.shaderVertices[i].ColorB = 1.0
		self.shaderVertices[i].ColorA = 1.0
	}
}
