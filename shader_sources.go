package lopix

import _ "embed"

// TODO: consider using quasilyte's minifier and paste code directly

//go:embed filters/nearest.kage
var _nearest []byte

//go:embed filters/derivative.kage
var _derivative []byte

//go:embed filters/hermite.kage
var _hermite []byte

//go:embed filters/bicubic.kage
var _bicubic []byte

//go:embed filters/bilinear.kage
var _bilinear []byte

//go:embed filters/src_derivative.kage
var _srcDerivative []byte

//go:embed filters/src_hermite.kage
var _srcHermite []byte

//go:embed filters/src_bicubic.kage
var _srcBicubic []byte

//go:embed filters/src_bilinear.kage
var _srcBilinear []byte

var pkgSrcKageFilters [scalingFilterEndSentinel][]byte
func init() {
	pkgSrcKageFilters[Nearest] = _nearest
	pkgSrcKageFilters[Derivative] = _derivative
	pkgSrcKageFilters[Hermite] = _hermite
	pkgSrcKageFilters[Bicubic] = _bicubic
	pkgSrcKageFilters[Bilinear] = _bilinear
	
	pkgSrcKageFilters[SrcDerivative] = _srcDerivative
	pkgSrcKageFilters[SrcHermite] = _srcHermite
	pkgSrcKageFilters[SrcBicubic] = _srcBicubic
	pkgSrcKageFilters[SrcBilinear] = _srcBilinear
}
