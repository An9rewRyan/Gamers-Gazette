package views

import (
	"d/go/utils/database"
	"d/go/utils/html"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func ApiArticleDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	article, err := database.Select_article(vars["id"])
	if err != nil {
		data := map[string]string{
			"Articles": "404! No article with such id is found!",
		}
		html.Render(w, r, data, "api.html")
		return
	}
	b, err := json.Marshal(&article)
	if err != nil {
		fmt.Println(err)
	}
	data := map[string]string{
		"Articles": "Deleted article: " + string(b),
	}
	_, err = database.Delete_article(vars["id"])
	if err != nil {
		fmt.Println(err)
	}
	html.Render(w, r, data, "api.html")
}
