package httputils

import(
	"fmt"
//	"mime"
	"strings"
	"net/http"
)

func RequestAcceptsJSON(r *http.Request) error {
	for _, ctype := range r.Header["Accept"] {
//
// XXX - Tad: doesn't work. results in this error: "mime: unexpected content after media subtype"
//
//		mtype, _, err := mime.ParseMediaType(ctype)
//		if err != nil {
//			return err
//		}
//
//		if mtype == "application/json" {
//			return nil
//		}

		// XXX - Tad: This is a hacky implementation, it should be more robust.
		if strings.Contains(ctype, "application/json") {
			return nil
		}
	}

	return fmt.Errorf("not json")
}
