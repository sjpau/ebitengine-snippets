package misc

import (
	"embed"

	"github.com/tinne26/etxt"
	"github.com/tinne26/etxt/ecache"
)

type Context struct {
	Files     *embed.FS
	Fonts     *etxt.FontLibrary
	FontCache *ecache.DefaultCache
}

func GenerateContext(fs *embed.FS) (*Context, error) {
	fonts := etxt.NewFontLibrary()
	loaded, _, err := fonts.ParseEmbedDirFonts("assets/fnt", fs)
	if err != nil {
		return nil, err
	}
	if loaded < 1 {
		panic("Failed to load fonts!")
	}
	/*32M Caching*/
	fontCache, err := ecache.NewDefaultCache(32 * 1024 * 1024)
	if err != nil {
		return nil, err
	}
	return &Context{
		Files:     fs,
		Fonts:     fonts,
		FontCache: fontCache,
	}, nil
}
