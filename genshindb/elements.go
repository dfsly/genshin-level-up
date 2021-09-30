package genshindb

import (
	"bytes"
	_ "embed"
	"encoding/base64"
	"image/png"

	"github.com/nfnt/resize"
)

var (
	//go:embed icons/anemo.png
	anemo []byte
	//go:embed icons/cryo.png
	cryo []byte
	//go:embed icons/dendro.png
	dendro []byte
	//go:embed icons/electro.png
	electro []byte
	//go:embed icons/geo.png
	geo []byte
	//go:embed icons/hydro.png
	hydro []byte
	//go:embed icons/pyro.png
	pyro []byte
)

var ElementIcons = map[string]string{}

func init() {
	ElementIcons["anemo"] = resizeAndToDataURI(anemo)
	ElementIcons["cryo"] = resizeAndToDataURI(cryo)
	ElementIcons["dendro"] = resizeAndToDataURI(dendro)
	ElementIcons["electro"] = resizeAndToDataURI(electro)
	ElementIcons["geo"] = resizeAndToDataURI(geo)
	ElementIcons["hydro"] = resizeAndToDataURI(hydro)
	ElementIcons["pyro"] = resizeAndToDataURI(pyro)
}

func resizeAndToDataURI(d []byte) string {
	i, _ := png.Decode(bytes.NewBuffer(d))
	next := resize.Resize(32, 0, i, resize.Lanczos3)

	b := bytes.NewBuffer(nil)
	_ = png.Encode(b, next)
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(b.Bytes())
}
