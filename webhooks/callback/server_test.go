package callback

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
)

func ExampleHandler() {
	// output all received messages to Stdout
	receiver := &JSONReceiver{json.NewEncoder(os.Stdout)}

	req := httptest.NewRequest(
		http.MethodPost,
		"/",
		bytes.NewBufferString(`[
        	{"text": "Hey Rasa!"},
        	{"image": "http://example.com/image.jpg"}
		]`),
	)
	res := httptest.NewRecorder()
	(&Handler{receiver}).ServeHTTP(res, req)

	fmt.Println(res.Result().Status)

	// Output:
	// {"text":"Hey Rasa!"}
	// {"image":"http://example.com/image.jpg"}
	// 200 OK
}
