package lopix

import "image"

import "github.com/hajimehoshi/ebiten/v2"

// Similar to [ebiten.Game], but without the Layout() method:
// we only need a logical resolution, which is set at [Run]().
type Game interface {
	// Updates the game logic
	Update() error

	// Draws the game contents
	Draw(*ebiten.Image)
}

// Equivalent to [ebiten.RunGame](), but expecting
// a lopix [Game] instead of an [ebiten.Game].
// 
// You must remember to [SetResolution]() before this.
func Run(game Game) error {
	w, h := pkgController.getResolution()
	if w == 0 && h == 0 {
		panic("Must set the game resolution with lopix.SetResolution(width, height) before lopix.Run()")
	} else if w < 1 || h < 1 {
		panic("invalid resolution") // should not be possible
	}
	pkgController.innerGame = game
	return ebiten.RunGame(&pkgController)
}

// Sets the windowed size to a reasonably visible but still
// pixel-perfect value.
//
// If the game is currently fullscreened, it will remain
// fullscreened, but the windowed size will still be updated.
func AutoResizeWindow() {
	pkgController.autoResizeWindow()
}

// Must be called at least once before [Run]().
// 
// Multiple resolutions are typically only relevant if you are trying
// to support different aspect ratios (e.g. ultrawide), which isn't
// even suitable for all types of games... and definitely not a common
// concern for the kind of games lopix tries to support.
//
// See also [GetResolution]().
//
// Must only be called during initialization or [Game].Update().
func SetResolution(width, height int) {
	pkgController.setResolution(width, height)
}

// Returns the last pair of values passed to [SetResolution]().
//
// With lopix, the game resolution and the size of the canvas
// passed to [Game].Draw() are always the same.
func GetResolution() (width, height int) {
	return pkgController.getResolution()
}

// Utility function, equivalent to obtaining the resolution
// and returning a rectangle of that size with (0, 0) origin.
func GetLogicalBounds() image.Rectangle {
	w, h := pkgController.getResolution()
	return image.Rect(0, 0, w, h)
}

// Queues the given callback function to be invoked after
// the current draw function and any other queued draws
// finish.
//
// Despite lopix focusing on low resolution pixel art games,
// in some cases you might still want to render vectorial UI,
// apply shaders for screen effects, draw high resolution
// backgrounds... and you can't do that on the logical canvas.
//
// This method can only be invoked during the draw stage.
// Multiple draws might be queued. See also [QueueDraw]().
//
// Must only be called from [Game].Draw() or successive
// draw callbacks.
func QueueHiResDraw(callback func(*ebiten.Image)) {
	pkgController.queueHiResDraw(callback)
}

// While drawing in high resolution, sometimes not the whole
// canvas is used due to aspect ratio mismatches. This function
// returns the active area for the high resolution canvas.
//
// Notice that [QueueHiResDraw]() callbacks receive the full
// resolution canvas in case you want to fill the black margins
// yourself.
func HiResActiveArea() image.Rectangle {
	return pkgController.hiResActiveArea()	
}

// See [QueueHiResDraw](). If you need to interleave high
// resolution and logically rendered layers, you might need
// to make use of use this function. If you aren't using
// [QueueHiResDraw]() in the first place, you should ignore
// this function.
//
// Notice that the canvas passed to the callback will be
// preemptively cleared if the previous draw was a high
// resolution draw.
//
// Must only be called from [Game].Draw() or successive
// draw callbacks.
func QueueDraw(callback func(*ebiten.Image)) {
	pkgController.queueLogicalDraw(callback)
}

// In some games and applications it's possible to spare
// GPU by using [ebiten.SetScreenClearedEveryFrame](false)
// and omitting redundant draw calls.
//
// The redraw manager allows you to synchronize this
// process with lopix itself, as there are some projections
// that would otherwise fall outside your control.
//
// By default, redraws are executed on every frame. If you
// want to manage them more efficiently, you can do the
// following:
//  - Make sure to disable ebitengine's screen clear.
//  - Opt into managed redraws with [RedrawManager.SetManaged](true).
//  - Whenever a redraw becomes necessary, issue a
//    [RedrawManager.Request]().
//  - On [Game].Draw(), if ![RedrawManager.Pending](), skip the draw.
type RedrawManager controller

// See [RedrawManager].
func Redraw() *RedrawManager {
	return (*RedrawManager)(&pkgController)
}

// Enables or disables manual redraw management. By default,
// redraw management is disabled and the screen is redrawn
// every frame.
//
// Must only be called during initialization or [Game].Update().
func (self *RedrawManager) SetManaged(managed bool) {
	pkgController.redrawSetManaged(managed)
}

// Returns whether manual redraw management is enabled or not.
func (self *RedrawManager) IsManaged() bool {
	return pkgController.redrawIsManaged()
}

// Notifies lopix that the next [Game].Draw() needs to be
// projected on the screen. Requests are typically issued
// when relevant input or events are detected during
// [Game].Update().
//
// This function can be called multiple times within a single
// update without issue, it's only setting an internal flag
// equivalent to "needs redraw".
func (self *RedrawManager) Request() {
	pkgController.redrawRequest()
}

// Returns whether a redraw is still pending. Notice that
// besides explicit requests, a redraw can also be pending
// due to a canvas resize, the modification of the scaling
// properties, etc.
func (self *RedrawManager) Pending() bool {
	return pkgController.redrawPending()
}

// Signal the redraw manager to clear both the logical screen
// and the high resolution canvas before the next [Game].Draw().
func (self *RedrawManager) ScheduleClear() {
	pkgController.redrawScheduleClear()
}

// Scaling modes can be changed through [SetScalingMode]().
//
// Letting the player change these through the game options
// is generally nice. The default mode is [Proportional].
type ScalingMode uint8
const (
	// Proportional projects the screen to be displayed as big as
	// possible while preserving the game's aspect ratio.
	Proportional ScalingMode = iota

	// Also known as "integer scaling". Depending on the screen or
	// window size, a lot of space could be left unused... but at
	// least the results tend to be as sharp as possible.
	PixelPerfect

	// Completley fill the screen no matter how ugly it might get. 
	Stretched
)

// Returns "Propotional", "Pixel-Perfect" or "Stretched".
func (self ScalingMode) String() string {
	switch self {
	case Proportional : return "Proportional"
	case PixelPerfect : return "Pixel-Perfect"
	case Stretched    : return "Stretched"
	default:
		panic("invalid ScalingMode")
	}
}

// Changes the scaling mode. The default is [Proportional].
//
// Must only be called during initialization or [Game].Update().
func SetScalingMode(mode ScalingMode) {
	pkgController.setScalingMode(mode)
}

// Returns the current scaling mode. The default is [Proportional].
func GetScalingMode() ScalingMode {
	return pkgController.getScalingMode()
}

// Scaling filters can be changed through [SetScalingFilter]().
//
// Many filters are only provided as comparison points, not
// because they necessarily offer great results. For the purposes
// of this package, I'd generally recommend the AASampling*
// or the [Hermite] filters.
//
// In some very specific cases, [Nearest] might also be useful.
type ScalingFilter uint8
const (
	Hermite ScalingFilter = iota
	AASamplingSoft
	AASamplingSharp
	Nearest
	Bicubic
	Bilinear
	SrcHermite
	SrcBicubic
	SrcBilinear
	scalingFilterEndSentinel
)

// Returns a string representation of the scaling filter.
func (self ScalingFilter) String() string {
	switch self {
	case AASamplingSoft  : return "AASamplingSoft"
	case AASamplingSharp : return "AASamplingSharp"
	case Nearest         : return "Nearest"
	case Hermite         : return "Hermite"
	case Bicubic         : return "Bicubic"
	case Bilinear        : return "Bilinear"
	case SrcHermite      : return "SrcHermite"
	case SrcBicubic      : return "SrcBicubic"
	case SrcBilinear     : return "SrcBilinear"
	default:
		panic("invalid ScalingFilter")
	}
}

// Changes the scaling filter. The default is [Hermite].
//
// Must only be called during initialization or [Game].Update().
//
// The first time you set a filter explicitly, its shader is
// also compiled. This means that this function can be effectively
// used to precompile the relevant shaders. Otherwise, the shader
// will be recompiled the first time it's actually needed in order
// to draw.
func SetScalingFilter(filter ScalingFilter) {
	pkgController.setScalingFilter(filter)
}

// Returns the current scaling filter. The default is [Hermite].
func GetScalingFilter() ScalingFilter {
	return pkgController.getScalingFilter()
}

// TODO: probably a Project(from, to *ebiten.Image, mode ScalingMode, filter ScalingFilter)
// method would be beneficial. Could also be ScalingMode.Project() or something, but
// I don't like any in particular.

// Transforms coordinates obtained from [ebiten.CursorPosition]() and
// similar functions to relative coordinates between 0 and 1.
//
// If the coordinates fall outside the active canvas they will be clamped
// to the closest point inside it, returning 0 or 1 for any clamped axis.
func ToRelativeCoords(x, y int) (float64, float64) {
	return pkgController.toRelativeCoords(x, y)
}

// Transforms coordinates obtained from [ebiten.CursorPosition]() and
// similar functions to coordinates within the game's logical resolution.
func ToLogicalCoords(x, y int) (int, int) {
	return pkgController.toLogicalCoords(x, y)
}

// func LogicalCursorPosition() (int, int) {}
// func LogicalTouchPosition(id TouchID) (int, int) {}
