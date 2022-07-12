package writer

import (
	"image"
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tinne26/etxt"
	"golang.org/x/image/math/fixed"
)

// - \x01: use current x as next line start x
// - \x02: reset next line start x to original value
// - \x03: pause *6
// - \x04: pause *10
// - \x05: pause *13
// - \x07: set text hue to main color
// - \x08: set text hue to disabled color
// - \x09: set text hue to name color
// - \x0A: line break
// - \x0B: 50% line height shift, to be used after normal line breaks
// - \x0C: shake off
// - \x0D: shake on
// - \x0E: shake on soft
// - \x10: no-delay period .
// - \x11: intermittent off
// - \x12: intermittent on

type Writer struct {
	renderer     *etxt.Renderer
	text         string
	index        int
	delay        int
	intermittent int
}

const (
	CycleTicks = 120
	delaystd   = 3
)

func New(t string, ctx *Context) *Writer {
	r := etxt.NewStdRenderer()
	r.SetCacheHandler(ctx.FontCache.NewHandler())
	r.SetFont(ctx.Fonts.GetFont("Lunchtime Doubly So Regular"))
	r.SetAlign(etxt.YCenter, etxt.XCenter)
	r.SetQuantizationMode(etxt.QuantizeNone)

	return &Writer{
		renderer: r,
		text:     t,
		delay:    delaystd,
	}
}

func (self *Writer) ReachedEnd() bool {
	return self.index == len(self.text)
}

func (self *Writer) SkipToEnd() {
	self.delay = 0
	self.index = len(self.text)
}

func (self *Writer) Tick() {
	self.intermittent = (self.intermittent + 1) % CycleTicks
	if self.ReachedEnd() {
		return
	}

	self.delay -= 1
	if self.delay <= 0 {
		self.delay = delaystd

		char := self.text[self.index]
		if char >= 32 {
			if char == 44 {
				self.delay *= 5
			} // ,
			if char == 46 || char == 33 || char == 63 {
				self.delay *= 8
			} // .!?
			if char == 58 {
				self.delay *= 8
			} // :
		} else {
			switch char {
			case 3:
				self.delay *= 6
			case 4:
				self.delay *= 10
			case 5:
				self.delay *= 13
			default:
				self.delay = 0
			}
		}
		self.index += 1

		// support utf8
		for self.index < len(self.text) && self.text[self.index] >= 128 {
			self.index += 1
		}
	}
}

func (self *Writer) Draw(screen *ebiten.Image, fontSize int, bounds image.Rectangle, opacity uint8) {
	self.renderer.SetTarget(screen)
	self.renderer.SetSizePx(fontSize)
	mainColor := color.RGBA{240, 240, 240, opacity}
	disabledColor := color.RGBA{140, 140, 140, opacity}
	nameColor := color.RGBA{255, 99, 112, opacity}
	self.renderer.SetColor(mainColor)

	shakeIntensity := 0
	intermittentOn := false
	comingFromLineBreak := true
	lineAdvance := self.renderer.GetLineAdvance()
	feed := self.renderer.NewFeed(fixed.P(bounds.Min.X, bounds.Min.Y))
	originalBreakX := feed.LineBreakX

	index := 0
	for index < self.index {
		ascii := self.text[index]
		if ascii < 32 {
			// special characters and control codes
			switch ascii {
			case 1: // use current x as next line start x
				feed.LineBreakX = feed.Position.X
			case 2: // reset next line start x to original value
				feed.LineBreakX = originalBreakX
			case 3, 4, 5: // 6 ticks pause
				// explicit pauses, already processed in Update()
			case 7:
				self.renderer.SetColor(mainColor)
			case 8:
				self.renderer.SetColor(disabledColor)
			case 9:
				self.renderer.SetColor(nameColor)
			case 10: // line break
				feed.LineBreak()
				comingFromLineBreak = true
			case 11: // 50% line break
				feed.Position.Y += lineAdvance / 2
			case 12:
				shakeIntensity = 0
			case 13:
				shakeIntensity = 64
			case 14:
				shakeIntensity = 32
			case 16: // no-delay period.
				feed.Draw('.')
			case 17: // intermittent off
				intermittentOn = false
			case 18: // intermittent on
				intermittentOn = true
			default:
				panic(ascii)
			}
			index += 1
		} else if ascii == 32 {
			if !comingFromLineBreak {
				feed.Advance(' ')
			}
			index += 1
		} else { // draw next word
			comingFromLineBreak = false
			word := self.NextWord(index)

			// measure word to see if it fits
			width := self.renderer.SelectionRect(word).Width
			if (feed.Position.X + width).Ceil() > bounds.Max.X {
				feed.LineBreak()
			}

			// abort if we are going beyond the vertical working area
			if feed.Position.Y.Floor() >= bounds.Max.Y {
				return
			}

			// consider intermittency
			intrm := (intermittentOn && self.intermittent < CycleTicks/2)
			var clr color.RGBA
			if intrm {
				clr = self.renderer.GetColor().(color.RGBA)
				self.renderer.SetColor(color.RGBA{clr.R, clr.G, clr.B, clr.A / 2})
			}

			// draw the word character by character to apply shaking
			for _, codePoint := range word {
				if shakeIntensity > 0 {
					preY := feed.Position.Y
					vibr := fixed.Int26_6(rand.Intn(shakeIntensity))
					if rand.Intn(2) == 0 {
						vibr = -vibr
					}
					feed.Position.Y += vibr
					feed.Draw(codePoint)
					feed.Position.Y = preY
				} else {
					feed.Draw(codePoint)
				}
			}
			index += len(word)

			// restore color if intermittency active
			if intrm {
				self.renderer.SetColor(clr)
			}
		}

		// jump to next line if necessary
		if !comingFromLineBreak && feed.Position.X.Ceil() > bounds.Max.X {
			feed.LineBreak()
			comingFromLineBreak = true
		}
	}
}

func (self *Writer) NextWord(index int) string {
	start := index
	for index < self.index {
		if self.text[index] <= 32 {
			return self.text[start:index]
		}
		index += 1
	}
	return self.text[start:self.index]
}
