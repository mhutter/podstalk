package podstalk

import (
	"net/http"
	"os"
	"time"
)

var (
	logo = MustAsset("appuioli.png")
	info os.FileInfo
)

func init() {
	info, _ = AssetInfo("appuioli.png")
}

func LogoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Last-Modified", info.ModTime().Format(time.RFC1123))
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.Write(logo)
}
