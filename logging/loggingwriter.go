package logging

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
)

type LoggingWriter struct {
	writer http.ResponseWriter
	bytes  int64
	code   int
}

func NewLoggingWriter(writer http.ResponseWriter) *LoggingWriter {
	return &LoggingWriter{writer: writer}
}

func (lw *LoggingWriter) Write(data []byte) (count int, err error) {
	fmt.Printf("SH26 %+v\n", lw.writer)
	count, err = lw.writer.Write(data)
	fmt.Printf("SH27 %+v\n", lw.writer)
	lw.bytes += int64(count)
	return
}

func (lw *LoggingWriter) GetInternal() http.ResponseWriter {
	return lw.writer
}

func (lw *LoggingWriter) WriteHeader(code int) {
	fmt.Printf("SH28 %+v\n", lw.writer)
	lw.writer.WriteHeader(code)
	fmt.Printf("SH29 %+v\n", lw.writer)
	if code == 0 {
		code = 200
	}
	lw.code = code
}

func (lw *LoggingWriter) Header() http.Header {
	return lw.writer.Header()
}

func (lw *LoggingWriter) Flush() {
	fl := lw.writer.(http.Flusher)
	fmt.Printf("SH30 %+v\n", lw.writer)
	fl.Flush()
}

func (lw *LoggingWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	fmt.Printf("SH31 %+v\n", lw.writer)
	hij, ok := lw.writer.(http.Hijacker)
	fmt.Printf("SH32 %+v\n", lw.writer)
	if ok {
		fmt.Printf("SH33 %+v\n", lw.writer)
		return hij.Hijack()
	}

	fmt.Printf("SH34 %+v\n", lw.writer)
	return nil, nil, fmt.Errorf("could not hijack connection")
}

func (lw *LoggingWriter) GetBytes() int64 {
	return lw.bytes
}

func (lw *LoggingWriter) GetCode() int {
	return lw.code
}
