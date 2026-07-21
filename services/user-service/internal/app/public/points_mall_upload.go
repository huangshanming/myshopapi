package public

import (
	"net/http"

	"mymall/pkg/httpserver"
	"mymall/services/user-service/internal/uploadpath"
)

func ServePointsMallUpload(w http.ResponseWriter, r *http.Request) {
	p := uploadpath.Abs("points-mall", httpserver.PathParam(r, "file"))
	http.ServeFile(w, r, p)
}
