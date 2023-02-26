// Copyright 2023 Xinhe Wang
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package uwsgi

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
)

func init() {
	caddy.RegisterModule(Transport{})
}

type Transport struct {
}

// CaddyModule returns the Caddy module information.
func (Transport) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.reverse_proxy.transport.uwsgi",
		New: func() caddy.Module { return new(Transport) },
	}
}

var headerNameReplacer = strings.NewReplacer("-", "_")

func writeBlockVar(buffer *bytes.Buffer, s string) {
	b := []byte(s)
	binary.Write(buffer, binary.LittleEndian, uint16(len(b)))
	buffer.Write(b)
}

// generateBlockVars returns the packet body of WSGI block vars generated from http.Request.
func generateBlockVars(req *http.Request) *bytes.Buffer {
	method := req.Method
	if method == "" {
		method = "GET"
	}
	vars := map[string]string{
		"REQUEST_METHOD":  method,
		"SERVER_PROTOCOL": req.Proto,
		"REQUEST_URI":     req.URL.RequestURI(),
		"QUERY_STRING":    req.URL.RawQuery,
		"HTTP_HOST":       req.Host,
		"REMOTE_ADDR":     req.RemoteAddr,
	}
	if req.TLS != nil {
		vars["HTTPS"] = "on"
	}
	for name, value := range req.Header {
		vars["HTTP_"+headerNameReplacer.Replace(strings.ToUpper(name))] = strings.Join(value, ", ")
	}

	var packetBody bytes.Buffer
	for key, val := range vars {
		writeBlockVar(&packetBody, key)
		writeBlockVar(&packetBody, val)
	}
	return &packetBody
}

func (t Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	conn, err := net.Dial("tcp", req.URL.Host)
	if err != nil {
		return nil, err
	}

	blockVars := generateBlockVars(req)
	conn.Write([]byte{0})                                            // modifier1
	binary.Write(conn, binary.LittleEndian, uint16(blockVars.Len())) // datasize
	conn.Write([]byte{0})                                            // modifier2
	io.Copy(conn, blockVars)                                         // packet body

	if req.Body != nil {
		io.Copy(conn, req.Body)
		req.Body.Close()
	}

	return http.ReadResponse(bufio.NewReader(conn), req)
}

// UnmarshalCaddyfile implements caddyfile.Unmarshaler.
func (t *Transport) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		if d.NextArg() {
			// too many args
			return d.ArgErr()
		}
	}
	return nil
}

var (
	_ http.RoundTripper     = (*Transport)(nil)
	_ caddyfile.Unmarshaler = (*Transport)(nil)
)
