package font

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	charsInRow = 128
	charRows   = 32
)

type Font interface {
	Size() (charW, charH int)
	PrintChar(r rune, x, y int, color, bgColor sdl.Color)
	Close()
}

type FontImpl struct {
	fname        string
	charW, charH int

	renderer *sdl.Renderer
	atlas    *sdl.Texture
	ttfFont  *ttf.Font
}

func Open(renderer *sdl.Renderer, fname string, size, charW, charH int) (Font, error) {
	ttfFont, err := ttf.OpenFont(fname, size)
	if err != nil {
		return nil, fmt.Errorf("could not open font \"%s\" (size %d): %v", fname, size, err)
	}

	if charW == -1 {
		charW = size
	}
	if charH == -1 {
		charH = size
	}

	font := &FontImpl{
		fname:    fname,
		charW:    charW,
		charH:    charH,
		renderer: renderer,
		ttfFont:  ttfFont,
	}

	err = font.RenderAtlas()
	if err != nil {
		return nil, err
	}

	return font, nil
}

func (f *FontImpl) printCharToAtlas(atlas *sdl.Surface, r rune, color sdl.Color) error {
	surface, err := f.ttfFont.RenderGlyphBlended(r, color)
	if err != nil {
		if err.Error() == "Text has zero width" {
			return nil
		}
		return fmt.Errorf("could not render char %d: %w", r, err)
	}
	defer surface.Free()

	x := r % charsInRow * int32(f.charW)
	y := r / charsInRow * int32(f.charH)

	dst := sdl.Rect{X: x, Y: y, W: surface.W, H: surface.H}
	err = surface.Blit(nil, atlas, &dst)
	if err != nil {
		return fmt.Errorf("could not blit char %d", r)
	}

	return nil
}

func (f *FontImpl) RenderAtlas() error {
	atlasW := f.charW * charsInRow
	atlasH := f.charH * charRows
	fmt.Println("ATLAS SIZE ", atlasW, atlasH)
	color := sdl.Color{R: 255, G: 255, B: 255, A: 255}

	atlasSurface, err := sdl.CreateRGBSurface(0, int32(atlasW), int32(atlasH), 32, 0xFF0000, 0xFF00, 0xFF, 0xFF000000)
	defer atlasSurface.Free()

	for r := rune(33); r < 1280; r++ {
		err := f.printCharToAtlas(atlasSurface, r, color)
		if err != nil {
			return err
		}

		// Skip range
		if r == 512 {
			r = 1023
		}
	}

	f.atlas, err = f.renderer.CreateTextureFromSurface(atlasSurface)
	if err != nil {
		return fmt.Errorf("could not create atlas texture: %v", err)
	}

	return nil
}

func (f *FontImpl) Close() {
	if f.atlas != nil {
		f.atlas.Destroy()
	}
	if f.ttfFont != nil {
		f.ttfFont.Close()
	}
}

func (f *FontImpl) Size() (charW, charH int) {
	charW, charH = f.charW, f.charH
	return
}

func (f *FontImpl) PrintChar(r rune, x, y int, color, bgColor sdl.Color) {
	w, h := int32(f.charW), int32(f.charH)
	srcX := r % int32(charsInRow) * int32(f.charW)
	srcY := r / int32(charsInRow) * int32(f.charH)
	src := sdl.Rect{X: srcX, Y: srcY, W: w, H: h}
	dst := sdl.Rect{X: int32(x), Y: int32(y), W: w, H: h}

	if bgColor.A != 0 {
		f.renderer.SetDrawColor(bgColor.R, bgColor.G, bgColor.B, bgColor.A)
		f.renderer.FillRect(&dst)
	}

	f.atlas.SetColorMod(color.R, color.G, color.B)
	f.renderer.Copy(f.atlas, &src, &dst)
}
