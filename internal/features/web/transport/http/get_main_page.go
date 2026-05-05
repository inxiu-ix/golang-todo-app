package web_transport_http

import (
	"net/http"
	"os"
	"path"
)

func (h *WebHTTPHandler) GetMainPage(w http.ResponseWriter, r *http.Request) {
	filePath := path.Join(os.Getenv("PROJECT_ROOT"), "public/index.html")
	http.ServeFile(w, r, filePath)
}
