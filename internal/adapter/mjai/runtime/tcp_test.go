package runtime_test

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strings"
	"testing"

	mjairuntime "github.com/Apricot-S/mjai-manue-go/internal/adapter/mjai/runtime"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/ai"
)

func TestRunTCP_RespondsToEachServerMessage(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Listen() failed: %v", err)
	}
	defer ln.Close()

	serverErr := make(chan error, 1)
	go func() {
		conn, err := ln.Accept()
		if err != nil {
			serverErr <- err
			return
		}
		defer conn.Close()

		r := bufio.NewReader(conn)
		w := bufio.NewWriter(conn)
		exchanges := []struct {
			in   string
			want string
		}{
			{
				in:   `{"type":"hello","protocol":"mjsonp","protocol_version":3}`,
				want: `{"type":"join","name":"tsumogiri","room":"room"}`,
			},
			{
				in:   `{"type":"start_game","gametype":"tonpu","id":0,"names":["A","B","C","D"]}`,
				want: `{"type":"none"}`,
			},
		}
		for _, tt := range exchanges {
			if _, err := fmt.Fprintln(w, tt.in); err != nil {
				serverErr <- err
				return
			}
			if err := w.Flush(); err != nil {
				serverErr <- err
				return
			}
			got, err := r.ReadString('\n')
			if err != nil {
				serverErr <- err
				return
			}
			if strings.TrimSuffix(got, "\n") != tt.want {
				serverErr <- fmt.Errorf("response = %q, want %q", got, tt.want+"\n")
				return
			}
		}

		if _, err := fmt.Fprintln(w, `{"type":"end_game","scores":[25000,25000,25000,25000]}`); err != nil {
			serverErr <- err
			return
		}
		if err := w.Flush(); err != nil {
			serverErr <- err
			return
		}
		got, err := r.ReadString('\n')
		if err == nil {
			serverErr <- fmt.Errorf("end_game response = %q, want connection close without response", got)
			return
		}
		serverErr <- nil
	}()

	err = mjairuntime.RunTCP(mjairuntime.TCPConfig{
		Name:  "tsumogiri",
		URL:   "mjsonp://" + ln.Addr().String() + "/room",
		Agent: ai.NewTsumogiriAgent(),
	})
	if err != nil {
		t.Fatalf("RunTCP() failed: %v", err)
	}
	if err := <-serverErr; err != nil {
		t.Fatalf("server failed: %v", err)
	}
}

func TestRunTCP_InvalidURLIsUsageError(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr string
	}{
		{
			name:    "invalid URL",
			url:     "://",
			wantErr: "missing protocol scheme",
		},
		{
			name:    "unsupported scheme",
			url:     "stdio://127.0.0.1:11600/room",
			wantErr: `unsupported URL scheme "stdio"`,
		},
		{
			name:    "missing host",
			url:     "mjsonp:///room",
			wantErr: "mjsonp URL requires host:port",
		},
		{
			name:    "missing port",
			url:     "mjsonp://127.0.0.1/room",
			wantErr: "mjsonp URL requires port",
		},
		{
			name:    "missing room",
			url:     "mjsonp://127.0.0.1:11600",
			wantErr: "mjsonp URL requires room path",
		},
		{
			name:    "nested room path",
			url:     "mjsonp://127.0.0.1:11600/room/extra",
			wantErr: "mjsonp URL requires room path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mjairuntime.RunTCP(mjairuntime.TCPConfig{
				Name:  "tsumogiri",
				URL:   tt.url,
				Agent: ai.NewTsumogiriAgent(),
			})
			if err == nil {
				t.Fatal("RunTCP() succeeded unexpectedly")
			}
			if _, ok := errors.AsType[*mjairuntime.UsageError](err); !ok {
				t.Errorf("errors.AsType[*UsageError](%v) = false, want true", err)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("error = %q, want to contain %q", err.Error(), tt.wantErr)
			}
		})
	}
}
