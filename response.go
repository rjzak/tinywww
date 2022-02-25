package tinywww

type TinyResponse struct {
	Buffer  []byte
	Headers map[string]string
}

func NewTinyResponse() *TinyResponse {
	return &TinyResponse{
		Buffer:  make([]byte, 0),
		Headers: make(map[string]string),
	}
}

func (tr *TinyResponse) SetHeaders(headers map[string]string) {
	tr.Headers = headers
}

func (tr *TinyResponse) SetHeader(key, value string) {
	tr.Headers[key] = value
}

func (tr *TinyResponse) Write(buffer []byte) (int, error) {
	tr.Buffer = buffer
	return len(tr.Buffer), nil
}

func (tr *TinyResponse) Append(buffer []byte) (int, error) {
	tr.Buffer = append(tr.Buffer, buffer...)
	return len(tr.Buffer), nil
}
