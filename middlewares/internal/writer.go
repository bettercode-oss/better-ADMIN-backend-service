// MIT License

// Copyright (c) 2021 vearne

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
package internal

import (
	"better-admin-backend-service/middlewares/internal/buffpool"
	"bytes"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

// Copy and paste from https://github.com/vearne/gin-timeout

type ResponseWriter struct {
	gin.ResponseWriter
	h    http.Header
	body *bytes.Buffer

	code        int
	mu          sync.Mutex
	timedOut    bool
	wroteHeader bool
}

func NewResponseWriter(c *gin.Context) *ResponseWriter {
	buffer := buffpool.GetBuff()
	writer := &ResponseWriter{
		body:           buffer,
		ResponseWriter: c.Writer,
		h:              make(http.Header),
	}
	return writer
}

func (w *ResponseWriter) Write(b []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.timedOut {
		return 0, nil
	}

	return w.body.Write(b)
}

func (w *ResponseWriter) WriteHeader(code int) {
	//checkWriteHeaderCode(code)
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.timedOut {
		return
	}
	w.writeHeader(code)
}

func (w *ResponseWriter) writeHeader(code int) {
	w.wroteHeader = true
	w.code = code
}

func (w *ResponseWriter) WriteHeaderNow() {}

func (w *ResponseWriter) Header() http.Header {
	return w.h
}

func (w *ResponseWriter) Body() *bytes.Buffer {
	return w.body
}

//func checkWriteHeaderCode(code int) {
//	if code < 100 || code > 999 {
//		panic(fmt.Sprintf("invalid WriteHeader code %v", code))
//	}
//}

func (w *ResponseWriter) Done(c *gin.Context) {
	dst := w.ResponseWriter.Header()
	for k, vv := range w.Header() {
		dst[k] = vv
	}

	if !w.wroteHeader {
		w.code = http.StatusOK
	}

	w.ResponseWriter.WriteHeader(w.code)
	if !(w.code == http.StatusNoContent || w.code == http.StatusNotModified) {
		_, err := w.ResponseWriter.Write(w.body.Bytes())
		if err != nil {
			panic(err)
		}
		buffpool.PutBuff(w.body)
	}
}
