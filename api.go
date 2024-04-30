package lopix

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

// Sets the windowed size to a reasonable, pixel-perfect value
// that will usually take most of the window, but not all of it.
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
func SetResolution(width, height int) {
	pkgController.setResolution(width, height)
}

// Returns the last pair of values passed to [SetResolution]().
//
// With lopix, the game resolution and the size of the canvas
// passed to Game.Draw() are always the same.
func GetResolution() (width, height int) {
	return pkgController.getResolution()
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
func QueueHiResDraw(callback func(*ebiten.Image)) {
	pkgController.queueHiResDraw(callback)
}

// See [QueueHiResDraw](). If you need to interleave high
// resolution and logically rendered layers, you might need
// to make use of use this function. If you aren't using
// [QueueHiResDraw]() in the first place, you should ignore
// this function.
//
// Notice that the canvas passed to the callback will be
// preemptive cleared if the previous draw was a high
// resolution draw.
func QueueDraw(callback func(*ebiten.Image)) {
	pkgController.queueLogicalDraw(callback)
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

// Utility method to get the next scaling mode constant.
func (self ScalingMode) Next() ScalingMode {
	switch self {
	case Proportional : return PixelPerfect
	case PixelPerfect : return Stretched
	case Stretched    : return Proportional
	default:
		panic("invalid ScalingMode")
	}
}

// Changes the scaling mode. The default is [Proportional].
func SetScalingMode(mode ScalingMode) {
	pkgController.setScalingMode(mode)
}

// Returns the current scaling mode. The default is [Proportional].
func GetScalingMode() ScalingMode {
	return pkgController.getScalingMode()
}

// Scaling filters can be changed through [SetScalingFilter]().
//
// The default [Derivative] scaling filter is pretty much the
// only reason to use this package. Other filters can serve
// as comparison points for the dev, but they should not be
// exposed as configurable settings for the player.
//
// In some very specific cases, [Nearest] or [Bilinear] might
// be preferred over the default algorithm, but this would
// be more of an aesthetic choice than anything else.
type ScalingFilter uint8
const (
	Derivative ScalingFilter = iota 
	Nearest 
	Linear
	Bilinear
)

// Changes the scaling filter. The default is [Derivative].
func SetScalingFilter(filter ScalingFilter) {
	pkgController.setScalingFilter(filter)
}

// Returns the current scaling filter. The default is [Derivative].
func GetScalingFilter() ScalingFilter {
	return pkgController.getScalingFilter()
}

// TODO: probably a Project(from, to *ebiten.Image, mode ScalingMode, filter ScalingFilter)
// method would be beneficial. Could also be ScalingMode.Project() or something, but
// I don't like any in particular.

// Transforms coordinates obtained from [ebiten.CursorPosition]() and
// similar functions to relative coordinates between 0 and 1.
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
