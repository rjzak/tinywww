package tinywww

import (
	"bytes"
	"strings"
)

type TinyRequest struct {
	Body    []byte
	Request string
	Method  string
	HTTP    string
	Headers map[string]string
	Form    map[string]string
}

func NewTinyRequestFromBuffer(buffer []byte) *TinyRequest {
	lines := bytes.Split(buffer, []byte("\n"))
	header := make(map[string]string)
	form := make(map[string]string)

	requestString := string(lines[0])
	bodyStartCounter := len(requestString) + 1 // string size plus removed newline
	requestStringParts := strings.Split(requestString, " ")
	method := requestStringParts[0]

	if strings.Contains(requestStringParts[1], "?") {
		queryString := strings.SplitN(requestStringParts[1], "?", 2)[1]
		for _, formItem := range strings.Split(queryString, "&") {
			if strings.Contains(formItem, "=") {
				parts := strings.Split(formItem, "=")
				form[parts[0]] = parts[1]
			}
		}
	}

	for _, line := range lines[1:] {
		if !bytes.Contains(line, []byte(":")) || len(line) == 0 {
			// Done with the headers
			break
		}
		bodyStartCounter += len(line) + 1 // string size plus removed newline
		parts := bytes.SplitN(line, []byte(":"), 2)
		header[strings.TrimSpace(string(parts[0]))] = strings.TrimSpace(string(parts[1]))
	}

	return &TinyRequest{
		Body:    buffer[bodyStartCounter:],
		Request: requestString,
		Method:  method,
		HTTP:    strings.TrimSpace(requestStringParts[2]),
		Headers: header,
		Form:    form,
	}
}
