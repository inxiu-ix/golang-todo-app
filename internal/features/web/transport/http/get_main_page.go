package web_transport_http

import (
	"net/http"
	"os"
	"path"

	core_logger "github.com/inxiu-ix/golang-todo-app/internal/core/logger"
	core_http_response "github.com/inxiu-ix/golang-todo-app/internal/core/transport/http/response"
)

func (h *WebHTTPHandler) GetMainPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	filePath := path.Join(os.Getenv("PROJECT_ROOT"), "public/index.html")

	f, err := os.Open(filePath)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to open main page")
		return
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to stat main page")
		return
	}

	http.ServeContent(w, r, "index.html", fi.ModTime(), f)
}
