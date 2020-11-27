package routing

import (
	"net/http"

	"github.com/ory/fosite"

	"github.com/matrix-org/dendrite/authapi/api"
	"github.com/matrix-org/dendrite/authapi/storage"
	userapi "github.com/matrix-org/dendrite/userapi/api"
	"github.com/matrix-org/dendrite/userapi/storage/accounts"
	"github.com/matrix-org/util"
)

func Authorize(rw http.ResponseWriter, req *http.Request, oauth2Provider fosite.OAuth2Provider, db storage.Database, accountDB accounts.Database, userAPI userapi.UserInternalAPI) {
	ctx := req.Context()
	log := util.GetLogger(ctx)

	ar, err := oauth2Provider.NewAuthorizeRequest(ctx, req)
	if err != nil {
		log.WithError(err).Error("Invalid authorization request")
		oauth2Provider.WriteAuthorizeError(rw, ar, err)
		return
	}

	req.ParseForm()
	if req.PostForm.Get("localpart") == "" {
		rw.Header().Set("Content-Type", "text/html; charset=utf-8")
		rw.Write([]byte(`
			<!DOCTYPE html>
			<html>
				<head>
					<meta charset="UTF-8" />
					<title>Login</title>
				</head>
				<body>
					<form method="POST">
						<h1>Login</h1>

						<div>
							<label for="localpart">Localpart</label>
							<input type="text" name="localpart" id="localpart" />
						</div>

						<div>
							<label for="password">Password</label>
							<input type="text" name="password" id="password" />
						</div>

						<button type="submit">Login</button>
					</form>
				</body>
			</html>
		`))
		return
	}

	localpart := req.PostForm.Get("localpart")
	password := req.PostForm.Get("password")
	account, err := accountDB.GetAccountByPassword(ctx, localpart, password)
	if err != nil {
		log.WithError(err).Error("Invalid credentials")
		return
	}

	for _, scope := range ar.GetRequestedScopes() {
		ar.GrantScope(scope)
	}

	var performRes userapi.PerformDeviceCreationResponse
	if err := userAPI.PerformDeviceCreation(ctx, &userapi.PerformDeviceCreationRequest{
		AccessToken: util.RandomString(20),
		IPAddr:      req.RemoteAddr,
		UserAgent:   req.UserAgent(),
		Localpart:   localpart,
	}, &performRes); err != nil {
		log.WithError(err).Error("Failed to perform device creation")
		return
	}

	// TODO: create the device
	session := api.NewSession(account.Localpart, performRes.Device.ID)

	response, err := oauth2Provider.NewAuthorizeResponse(ctx, ar, session)
	if err != nil {
		log.WithError(err).Error("Could not fullfil authorization request")
		oauth2Provider.WriteAuthorizeError(rw, ar, err)
		return
	}

	if err := db.CreateSession(ctx, ar.GetID(), session); err != nil {
		// TODO: write error
		return
	}

	oauth2Provider.WriteAuthorizeResponse(rw, ar, response)
}
