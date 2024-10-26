package output

import (
	"fmt"
	"net/http"
)

func WriteFatalError(w http.ResponseWriter, err error) {
	w.WriteHeader(200)
	fmt.Fprint(w, err)
}
