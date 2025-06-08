package telegrambot

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/arian-nj/master-card/back/internal/server"
	"github.com/go-chi/chi/v5"
	"github.com/go-telegram/bot"
)

type ApplicationTelegram struct {
	*server.CommonGlobals
	bot      *bot.Bot
	subbedFS fs.FS
	botToken string
}

func NewTelegramApplication(globalStructs *server.CommonGlobals) *ApplicationTelegram {
	return &ApplicationTelegram{
		CommonGlobals: globalStructs,
	}
}

//go:embed static
var staticContent embed.FS

func (app *ApplicationTelegram) StartTelegramBot(ctx context.Context) error {
	app.botToken = "8052428016:AAFl8AjzSiIG3owcVIm-tSpGQ0iq_IHo78Q"
	botOpts := []bot.Option{
		bot.WithDefaultHandler(app.defaultHandler),
	}

	newBot, err := bot.New(app.botToken, botOpts...)
	if err != nil {
		panic(err)
	}
	app.bot = newBot
	app.botRoutes(newBot)

	app.Logger.Info("starting telegram bot")
	newBot.Start(ctx)
	return nil
}

func (app *ApplicationTelegram) WebsiteRoutes() *chi.Mux {
	subbedFS, errSubFS := fs.Sub(staticContent, "static")
	if errSubFS != nil {
		app.Logger.Panic(fmt.Errorf("failed to create sub fs: %w", errSubFS).Error())
	}
	app.subbedFS = subbedFS

	// http router
	mux := chi.NewRouter()

	mux.NotFound(app.NotFound)
	mux.MethodNotAllowed(app.MethodNotAllowed)

	mux.Use(app.RecoverPanic)

	mux.Use(app.CorsMiddlewareFunc)
	mux.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ext := strings.ToLower(filepath.Ext(r.URL.Path))
			switch ext {
			case ".js":
				w.Header().Set("Content-Type", "application/javascript")
			case ".css":
				w.Header().Set("Content-Type", "text/css")
			case ".html":
				w.Header().Set("Content-Type", "text/html")
			case ".wasm":
				w.Header().Set("Content-Type", "application/wasm")
			default:
				// Try to use Go's built-in MIME type detection
				if mtype := mime.TypeByExtension(ext); mtype != "" {
					w.Header().Set("Content-Type", mtype)
				}
			}
			next.ServeHTTP(w, r)
		})
	})
	mux.Handle("/*", http.FileServer(http.FS(app.subbedFS)))

	return mux
}
