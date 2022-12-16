package service

import (
	"context"

	"github.com/MFCaballero/pangea/appdomain"
	"github.com/MFCaballero/pangea/repository"
	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
)

type SessionData struct {
	FlashMessage string
	Form         interface{}
	User         appdomain.User
	LoggedIn     bool
}

func NewSessionManager() *scs.SessionManager {
	sessions := scs.New()
	sessions.Store = postgresstore.New(repository.DBGet().DB)

	return sessions
}

func GetSessionData(session *scs.SessionManager, ctx context.Context) SessionData {
	var data SessionData

	data.FlashMessage = session.PopString(ctx, "flash")
	data.User, data.LoggedIn = ctx.Value("user").(appdomain.User)

	data.Form = session.Pop(ctx, "form")
	if data.Form == nil {
		data.Form = map[string]string{}
	}

	return data
}
