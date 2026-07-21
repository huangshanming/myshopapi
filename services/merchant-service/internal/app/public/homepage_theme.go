package public

import (
	"context"
	"mymall/pkg/appinput"
	"net/http"

	"mymall/pkg/xerr"
)

func (h *HomepageThemeHandler) PublicThemeTiles(ctx context.Context, in appinput.CallInput) (any, error) {
	list, err := h.logic.ListThemeTiles()
	if err != nil {
		return nil, xerr.New(http.StatusInternalServerError, err.Error())
	}
	return map[string]interface{}{"list": list}, nil
}
