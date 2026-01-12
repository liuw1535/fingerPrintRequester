package requester

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"

	utls "github.com/refraction-networking/utls"
)

func ForwardResponse(resp *http.Response, conn net.Conn) error {
	defer conn.Close()
	
	// Cancel read timeout for streaming
	if tcpConn, ok := conn.(*net.TCPConn); ok {
		tcpConn.SetReadDeadline(time.Time{})
	} else if utlsConn, ok := conn.(*utls.UConn); ok {
		utlsConn.SetReadDeadline(time.Time{})
	}

	// Write status line
	fmt.Fprintf(os.Stdout, "HTTP/%d.%d %s\r\n", resp.ProtoMajor, resp.ProtoMinor, resp.Status)

	// Write headers
	for k, vv := range resp.Header {
		for _, v := range vv {
			fmt.Fprintf(os.Stdout, "%s: %s\r\n", k, v)
		}
	}
	fmt.Fprintf(os.Stdout, "\r\n")
	os.Stdout.Sync()

	// Stream body chunk by chunk
	buf := make([]byte, 8192)
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			os.Stdout.Write(buf[:n])
			os.Stdout.Sync()
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			resp.Body.Close()
			return err
		}
	}
	resp.Body.Close()
	return nil
}
