package podstalk

import "net/http"

var logo = MustAsset("appuioli.png")

func LogoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/png")
	w.Write(logo)
}
