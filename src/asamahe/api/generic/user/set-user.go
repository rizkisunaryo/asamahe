package user

import (
	"asamahe/structs"
	"encoding/json"
	"fmt"
	"net/http"
	"sunaryo/util/security"
	"sunaryo/util/sunhttp"
	//	"sunaryo/util/template"
)

func HandleSetUser() {
	url := "/api/user/setuser/"
	http.HandleFunc(url, setUser)
}

func setUser(w http.ResponseWriter, r *http.Request) {
	var inp structs.SetUserInput
	json.NewDecoder(r.Body).Decode(&inp)

	security.CheckKey(append([]byte(``), inp.Key...), append([]byte(``), inp.Id...))

	// Update user
	// whether the user exists or not
	userUrl := ("http://localhost:9200/asamahe/user/" + inp.Id)
	userReq := userReqJson(inp.Name, inp.PicUrl)
	err := sunhttp.PostNoResponse(userUrl, userReq)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, "{\"Status\":0}")
}

func userReqJson(name string, picUrl string) []byte {
	b := []byte(`{"Name":"`)
	b = append(b, name...)
	b = append(b, []byte(`","PicUrl":"`)...)
	b = append(b, picUrl...)
	b = append(b, []byte(`"}`)...)
	return b
}
