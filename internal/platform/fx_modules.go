package platform

import (
	"net/http"

	application "notes-service/internal/application/note"
	httpadapter "notes-service/internal/adapters/http"
	mongoadapter "notes-service/internal/adapters/mongo"

	"go.uber.org/fx"
)

var mongoModule = fx.Module("mongo",
	fx.Provide(
		NewMongoClient,
		NewNoteCollection,
		fx.Annotate(mongoadapter.NewNoteRepository, fx.As(new(application.Repository))),
	),
)

var applicationModule = fx.Module("application",
	fx.Provide(
		application.NewCreateNoteUseCase,
		application.NewListNotesUseCase,
		application.NewUpdateNoteUseCase,
		application.NewDeleteNoteUseCase,
	),
)

var httpModule = fx.Module("http",
	fx.Provide(
		func(uc *application.CreateNoteUseCase) httpadapter.NoteCreator { return uc },
		func(uc *application.ListNotesUseCase) httpadapter.NoteLister { return uc },
		func(uc *application.UpdateNoteUseCase) httpadapter.NoteUpdater { return uc },
		func(uc *application.DeleteNoteUseCase) httpadapter.NoteDeleter { return uc },
		httpadapter.NewNoteHandler,
		httpadapter.NewRouter,
		NewHTTPServer,
	),
)

func Modules() []fx.Option {
	return []fx.Option{
		mongoModule,
		applicationModule,
		httpModule,
		fx.Invoke(func(*http.Server) {}),
	}
}
