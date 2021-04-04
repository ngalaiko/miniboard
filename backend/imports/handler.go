package imports

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/ngalaiko/miniboard/backend/authorizations"
	"github.com/ngalaiko/miniboard/backend/httpx"
	"github.com/ngalaiko/miniboard/backend/imports/parser"
	"github.com/ngalaiko/miniboard/backend/operations"
	"github.com/ngalaiko/miniboard/backend/subscriptions"
	"github.com/ngalaiko/miniboard/backend/tags"
)

type logger interface {
	Error(string, ...interface{})
}

type tagsService interface {
	Create(ctx context.Context, userID string, title string) (*tags.Tag, error)
	GetByTitle(ctx context.Context, userID string, title string) (*tags.Tag, error)
}

type subscriptionsService interface {
	Create(ctx context.Context, userID string, url *url.URL, tagIDs []string) (*subscriptions.UserSubscription, error)
}

type operationService interface {
	Create(context.Context, string, operations.Task) (*operations.Operation, error)
}

// Handler handles import endpoints.
type Handler struct {
	logger               logger
	tagsService          tagsService
	subscriptionsService subscriptionsService
	operationService     operationService
}

// NewHandler returns a new imports handler.
func NewHandler(logger logger, tagsService tagsService, subscriptionsService subscriptionsService, operationService operationService) *Handler {
	return &Handler{
		logger:               logger,
		tagsService:          tagsService,
		subscriptionsService: subscriptionsService,
		operationService:     operationService,
	}
}

// Create returns http handler to handle create import.
func (h *Handler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, auth := authorizations.FromContext(r.Context())
		if !auth {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			h.logger.Error("failed to read request body: %s", err)
			httpx.InternalError(w, h.logger)
			return
		}

		if err := r.Body.Close(); err != nil {
			h.logger.Error("failed to close request body: %s", err)
			httpx.InternalError(w, h.logger)
			return
		}

		opml, err := parser.ParseOPML(body)
		if err != nil {
			httpx.Error(w, h.logger, fmt.Errorf("failed to parse opml: %w", err), http.StatusBadRequest)
			return
		}

		operation, err := h.operationService.Create(r.Context(), token.UserID, h.create(token.UserID, opml))
		switch {
		case err == nil:
			httpx.JSON(w, h.logger, operation, http.StatusOK)
		default:
			h.logger.Error("failed to import: %s", err)
			httpx.InternalError(w, h.logger)
		}
	}
}

func (h *Handler) create(userID string, opml *parser.OPML) operations.Task {
	return func(ctx context.Context, operation *operations.Operation, status chan<- *operations.Operation) error {
		response := struct {
			Tags          []*tags.Tag                       `json:"tags,omitempty"`
			Subscriptions []*subscriptions.UserSubscription `json:"subscriptions,omitempty"`
			Errors        map[string]string                 `json:"errors,omitempty"`
		}{
			Errors: map[string]string{},
		}
		for _, opmlTag := range opml.Tags {
			tag, err := h.tagsService.GetByTitle(ctx, userID, opmlTag.Title)
			switch {
			case err == nil:
			case errors.Is(err, tags.ErrNotFound):
				tag, err = h.tagsService.Create(ctx, userID, opmlTag.Title)
				if err != nil {
					h.logger.Error("failed to create tag: %s", err)
					return fmt.Errorf("internal error")
				}
				response.Tags = append(response.Tags, tag)
			default:
				h.logger.Error("failed to lookup tag by title: %s", err)
				return fmt.Errorf("internal error")
			}

			for _, opmlFeed := range opmlTag.Feeds {
				parsedURL, err := url.Parse(opmlFeed.URL)
				if err != nil {
					return fmt.Errorf("failed to parse feed url: %w", err)
				}

				subscription, err := h.subscriptionsService.Create(ctx, userID, parsedURL, []string{tag.ID})
				switch {
				case err == nil:
				case errors.Is(err, subscriptions.ErrAlreadyExists):
					continue
				case errors.Is(err, subscriptions.ErrFailedToDownloadSubscription),
					errors.Is(err, subscriptions.ErrFailedToParseSubscription):
					response.Errors[opmlFeed.URL] = err.Error()
					continue
				default:
					h.logger.Error("failed to create subscription: %s", err)
					return fmt.Errorf("internal error")
				}

				response.Subscriptions = append(response.Subscriptions, subscription)
			}
		}

		operation.Success(response)
		status <- operation
		return nil
	}
}
