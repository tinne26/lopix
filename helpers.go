package lopix

import "image"
import "image/color"

import "github.com/hajimehoshi/ebiten/v2"

var pkgReusableVertices []ebiten.Vertex
var pkgMiniMask *ebiten.Image
func init() {
	pkgReusableVertices = make([]ebiten.Vertex, 4)
	img := ebiten.NewImage(3, 3)
	img.Fill(color.White)
	pkgMiniMask = img.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
}

// Helper function to draw a single pixel.
func DrawPixel(target *ebiten.Image, x, y int, rgba color.RGBA) {
	DrawRect(target, image.Rect(x, y, x + 1, y + 1), rgba)
}

// Helper function to draw filled rectangles. Unlike ebiten.Image.Fill(),
// it doesn't replace the contents of the rect area, it draws over it.
func DrawRect(target *ebiten.Image, rect image.Rectangle, rgba color.RGBA) {
	bounds := rect.Bounds()
	r, g, b, a := rgba.RGBA()
	fr, fg, fb, fa := float32(r)/65535, float32(g)/65535, float32(b)/65535, float32(a)/65535

	pkgReusableVertices[0].DstX = float32(bounds.Min.X) // top left
	pkgReusableVertices[0].DstY = float32(bounds.Min.Y)
	pkgReusableVertices[1].DstX = float32(bounds.Min.X) // bottom left
	pkgReusableVertices[1].DstY = float32(bounds.Max.Y)
	pkgReusableVertices[2].DstX = float32(bounds.Max.X) // top right
	pkgReusableVertices[2].DstY = float32(bounds.Min.Y)
	pkgReusableVertices[3].DstX = float32(bounds.Max.X) // bottom right
	pkgReusableVertices[3].DstY = float32(bounds.Max.Y)
	for i := 0; i < 4; i++ {
		pkgReusableVertices[i].SrcX = 1.0
		pkgReusableVertices[i].SrcY = 1.0
		pkgReusableVertices[i].ColorR = fr
		pkgReusableVertices[i].ColorG = fg
		pkgReusableVertices[i].ColorB = fb
		pkgReusableVertices[i].ColorA = fa
	}	
	target.DrawTriangles(pkgReusableVertices, []uint16{0, 2, 1, 1, 2, 3}, pkgMiniMask, nil)
}
