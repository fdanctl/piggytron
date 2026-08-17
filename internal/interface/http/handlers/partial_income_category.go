package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/fdanctl/piggytron/internal/application/appincomecategory"
	"github.com/fdanctl/piggytron/internal/errs"
	"github.com/fdanctl/piggytron/internal/interface/http/httperror"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/fdanctl/piggytron/web/templates/components"
	"github.com/fdanctl/piggytron/web/templates/layouts"
	"github.com/fdanctl/piggytron/web/templates/partials"
	"github.com/fdanctl/piggytron/web/views"
)

// IncomeCategoriesHandler renders the "new income category" dialog (GET)
// and creates the category (POST).
type IncomeCategoriesHandler struct {
	service *appincomecategory.Service
}

func NewIncomeCategoriesHandler(s *appincomecategory.Service) *IncomeCategoriesHandler {
	return &IncomeCategoriesHandler{
		service: s,
	}
}

func (h *IncomeCategoriesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	action := r.PathValue("action")
	if action != "" {
		switch action {
		case "archive":
			if r.Method != http.MethodPost {
				http.NotFound(w, r)
			}
			h.PostArchive(w, r)
		}
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)

	case http.MethodPost:
		h.Post(w, r)

	case http.MethodDelete:
		h.Delete(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Get renders the new income category form in a dialog.
func (h *IncomeCategoriesHandler) Get(w http.ResponseWriter, r *http.Request) {
	sessionInfo, err := middleware.SessionInfoFromCtx(r.Context())
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}
	id := r.PathValue("id")
	view := views.NewIncomeCategoryForm()
	title := "New Income Category"
	if id != "" {
		cat, err := h.service.FindCategory(r.Context(), id, sessionInfo.UserID)
		if err != nil {
			httperror.SendError(w, r, err)
			return
		}
		view.Name = cat.Name()
		title = "Edit Income Category"
	}
	form := partials.IncomeCategoryForm(id, *view)
	components.DialogWrapper(
		"",
		components.DialogHeader("", title, nil),
		form,
		nil,
		nil,
	).Render(r.Context(), w)
}

// Post validates and creates the income category, appending the new item
// to the income categories list out-of-band.
func (h *IncomeCategoriesHandler) Post(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context())
	sessionInfo, err := middleware.SessionInfoFromCtx(r.Context())
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}

	id := r.PathValue("id")
	name := r.FormValue("name")
	view := views.IncomeCategoryForm{
		Name: name,
	}

	msgs := view.Validate()
	if len(msgs) > 0 {
		logger.Info("invalid form", "error", msgs)
		w.WriteHeader(http.StatusUnprocessableEntity)
		partials.IncomeCategoryForm(id, view).Render(r.Context(), w)
		return
	}

	if id == "" {
		category, err := h.service.CreateCategory(r.Context(), sessionInfo.UserID, name)
		if err != nil {
			view.SetError(err)
			form := partials.IncomeCategoryForm(id, view)
			httperror.SendFormError(w, r, err, form)
			return
		}

		icView := views.IncomeCategory{
			ID:   category.ID(),
			Name: category.Name(),
		}

		w.Header().Set("HX-Trigger", "incomeCategoryAdded")
		templ.Join(
			partials.IncomeCategoryForm(id, view),
			layouts.OOBWraper(
				"",
				"beforeend:#income-cat ul",
				nil,
				partials.CategoryItem(icView, templ.Attributes{"style": "animation-delay: 0s;"}),
			),
			components.SendToast(
				components.Success,
				"Category added",
			),
		).Render(r.Context(), w)
	} else {
		err := h.service.UpdateCategory(r.Context(), id, sessionInfo.UserID, name)
		if err != nil {
			view.SetError(err)
			form := partials.IncomeCategoryForm(id, view)
			httperror.SendFormError(w, r, err, form)
			return
		}
		w.Header().Set(
			"HX-Trigger",
			fmt.Sprintf(`{
			"closeModal": true,
			"contentPush": {
				"url": "/categories/%s"
			}
			}`, id),
		)
		templ.Join(
			partials.IncomeCategoryForm(id, view),
			components.SendToast(
				components.Success,
				"Category updated",
			),
		).Render(r.Context(), w)
	}
}

func (h *IncomeCategoriesHandler) Delete(w http.ResponseWriter, r *http.Request) {
	sessionInfo, err := middleware.SessionInfoFromCtx(r.Context())
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}
	id := r.PathValue("id")
	if id == "" {
		err := errs.NewAppError(
			errs.KindBadRequest,
			"Account ID is required",
			errors.New("no id passed"),
			"IncomeCategoriesHandler.Delete",
		)
		httperror.SendError(w, r, err)
		return
	}

	err = h.service.Delete(r.Context(), id, sessionInfo.UserID)
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}
	w.Header().Set(
		"HX-Trigger",
		`{"contentPush": { "url": "/categories","transition": "true" }}`,
	)
}

func (h *IncomeCategoriesHandler) PostArchive(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context())
	id := r.PathValue("id")
	logger.Debug("PostArchive", "id", id)

	err := h.service.Archive(r.Context(), id)
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}
	w.Header().Set(
		"HX-Trigger",
		`{"contentPush": { "url": "/categories" }}`,
	)
}
