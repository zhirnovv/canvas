package handling

import (
	"fmt"
	"github.com/zhirnovv/gochat/api/bodyParser"
	"net/http"
)

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	body, err := bodyParser.ParseJSON(r)

	if err != nil {
		err.WriteTo(w)
	}

	fmt.Println(body)
}
