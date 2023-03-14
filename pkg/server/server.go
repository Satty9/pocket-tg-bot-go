package server

import (
	"log"
	"net/http"
	"strconv"

	"github.com/satty9/pocket-tg-bot-go/pkg/repository"
	"github.com/zhashkevych/go-pocket-sdk"
)

type AuthorizationServer struct {
	server          *http.Server
	pocketClient    *pocket.Client
	tokenRepository repository.TokenRepositorier
	redirectURL     string
}

func NewAuthorizationServer(newPocketClient *pocket.Client, newTokenRepository repository.TokenRepositorier, newRedirectURL string) *AuthorizationServer {
	return &AuthorizationServer{
		// server don`t init
		pocketClient:    newPocketClient,
		tokenRepository: newTokenRepository,
		redirectURL:     newRedirectURL,
	}
}

func (server *AuthorizationServer) Start() error {
	// create server
	server.server = &http.Server{
		Addr:    ":80", // server open on port 80
		Handler: server,
	}

	return server.server.ListenAndServe()
}

func (server *AuthorizationServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	// server wait only for GET method
	if request.Method != http.MethodGet {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// request must have "chat_id"
	chatIDParam := request.URL.Query().Get("chat_id")
	if chatIDParam == "" {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	chatID, err := strconv.ParseInt(chatIDParam, 10, 64)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	requestToken, err := server.tokenRepository.Get(chatID, repository.RequestTokens)
	if err != nil {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	authResponse, err := server.pocketClient.Authorize(request.Context(), requestToken)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	// write value to repository
	err = server.tokenRepository.Save(chatID, authResponse.AccessToken, repository.AccessTokens)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("New user finish authorization\nchat_id: %d\n request_token: %s\naccess_token%s\n", chatID, requestToken, authResponse.AccessToken)

	// make redirect in user browser
	writer.Header().Add("Location", server.redirectURL)
	writer.WriteHeader(http.StatusMovedPermanently)
}
