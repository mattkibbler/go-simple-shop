package output

import (
	"fmt"
	"net/http"
)

func WriteFatalError(w http.ResponseWriter, err error) {
	fmt.Fprint(w, err)
}
