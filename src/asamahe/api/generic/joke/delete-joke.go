package joke

import (
	"asamahe/structs"
	"encoding/json"
	"fmt"
	"net/http"
	"sunaryo/util/security"
	"sunaryo/util/sunhttp"
	//	"sunaryo/util/template"
)

func HandleDeleteJoke() {
	url := "/api/joke/deletejoke/"
	http.HandleFunc(url, deleteJoke)
}

func deleteJoke(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var inp structs.DeleteJokeInput
	json.NewDecoder(r.Body).Decode(&inp)

	security.CheckKey(append([]byte(``), inp.Key...), append([]byte(``), inp.UserId...))

	// Get the joke by JokeId
	urlGetJoke := ("http://localhost:9200/asamahe/joke/" + inp.JokeId)
	var jokeRespStruct structs.JokeResponse
	err := sunhttp.GetResponse(urlGetJoke, &jokeRespStruct)
	if err != nil {
		panic(err)
	}

	// Check whether UserId (the deleter) is the Joker itself
	if inp.UserId == jokeRespStruct.Source.Joker {
		// If yes, he can delete his joke
		url := ("http://localhost:9200/asamahe/joke/" + inp.JokeId)
		err := sunhttp.Delete(url)
		if err != nil {
			panic(err)
		}
	} else {
		return
	}

	fmt.Fprintf(w, "{\"Status\":0}")
}
