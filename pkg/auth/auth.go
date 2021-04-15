package auth

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
	pg "github.com/vgarvardt/go-oauth2-pg/v4"
	"github.com/vgarvardt/go-pg-adapter/pgx4adapter"
	"os"
	"time"
)

var (
	ClientStore  *pg.ClientStore
	TokenManager *manage.Manager
	TokenServer  *server.Server
	tokenStore   *pg.TokenStore
)

func init() {
	pgConn, err := pgx.Connect(context.TODO(), os.Getenv("DB_URI"))
	if err != nil {
		log.Error(err)
		return
	}
	log.Info("Auth DB Connect OK")
	TokenManager = manage.NewDefaultManager()
	adapter := pgx4adapter.NewConn(pgConn)
	tokenStore, err = pg.NewTokenStore(adapter, pg.WithTokenStoreGCInterval(time.Minute))
	if err != nil {
		log.Error(err)
		return
	}
	log.Info("Token Store Okay")
	ClientStore, err = pg.NewClientStore(adapter)
	if err != nil {
		log.Error(err)
		return
	}
	log.Info("Client Store Okay")
	TokenManager.MapTokenStorage(tokenStore)
	TokenManager.MapClientStorage(ClientStore)
	TokenManager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)
	TokenManager.MapAccessGenerate(generates.NewJWTAccessGenerate("", []byte("00000000"), jwt.SigningMethodHS512))
	TokenServer = server.NewServer(server.NewConfig(), TokenManager)
}

func CleanUp() {
	tokenStore.Close()
}
