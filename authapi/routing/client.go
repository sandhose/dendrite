package routing

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"

	"github.com/matrix-org/util"

	"github.com/matrix-org/dendrite/authapi/api"
	"github.com/matrix-org/dendrite/authapi/storage"
)

var charset = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

const idLength = 8

func makeID() string {
	b := make([]rune, idLength)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func RegisterClient(req *http.Request, database storage.Database) util.JSONResponse {
	ctx := req.Context()

	var err error
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return util.JSONResponse{
			Code: http.StatusInternalServerError,
			JSON: fmt.Sprintf("%v", err),
		}
	}

	r := make(map[string]interface{})
	err = json.Unmarshal(body, &r)
	if err != nil {
		return util.JSONResponse{
			Code: http.StatusInternalServerError,
			JSON: fmt.Sprintf("%v", err),
		}
	}

	clientID := makeID()
	client := api.NewClient(clientID)
	err = client.Fill(r)
	if err != nil {
		return util.JSONResponse{
			Code: http.StatusInternalServerError,
			JSON: fmt.Sprintf("%v", err),
		}
	}

	err = database.CreateClient(ctx, client)
	if err != nil {
		return util.JSONResponse{
			Code: http.StatusInternalServerError,
			JSON: fmt.Sprintf("%v", err),
		}
	}

	resp := client.Serialize()

	return util.JSONResponse{
		Code: http.StatusCreated,
		JSON: resp,
	}
}
