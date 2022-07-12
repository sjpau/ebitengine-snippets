/*
	w, h := screen.Size()
	xscale, yscale := float64(w)/ScreenWidth, float64(h)/ScreenHeight
	scaling := int(math.Min(xscale, yscale))
	Width, Height := ScreenWidth*scaling, ScreenHeight*scaling
	ysf, xsf := float64(scaling), float64(scaling)
	marginHor := (w - Width) / 2
	marginVer := (h - Height) / 2

	o := &ebiten.DrawImageOptions{}
	if float64(scaling) < 1 {
		// Resolution is too low
		marginHor, marginVer = 0, 0
		xsf, ysf = float64(w)/ScreenWidth, float64(h)/ScreenHeight
		if w >= ScreenWidth {
			xsf = 1.0
			marginHor = (w - ScreenWidth) / 2
		}
		if h >= ScreenHeight {
			ysf = 1.0
			marginVer = (h - ScreenHeight) / 2
		}
	}
	o.GeoM.Scale(xsf, ysf)
	if marginHor != 0 || marginVer != 0 {
		screen.Fill(color.RGBA{24, 24, 24, 255})
		o.GeoM.Translate(float64(marginHor), float64(marginVer))
	}
	if scaling >= 1 {
	}
*/
