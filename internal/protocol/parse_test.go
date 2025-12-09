package protocol

import (
	"testing"

	"bufio"
	"strings"
)

func TestParseSetxNX(t *testing.T) {
	input := "*4\r\n$3\r\nSET\r\n$5\r\nmykey\r\n$5\r\nvalue\r\n$2\r\nNX\r\n"
	reader := bufio.NewReader(strings.NewReader(input))
	resp := &RESP{}
	req, err := resp.Parse(reader)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if req.cmd != "set" {
		t.Errorf("Expected cmd 'set', got '%s'", req.cmd)
	}
	if req.argsLen != 4 {
		t.Errorf("Expected argsLen 4, got %d", req.argsLen)
	}
	if strings.ToUpper(req.args[3]) != "NX" {
		t.Errorf("Expected NX flag, got %v", req.args[3])
	}
}

func TestParseSetxXX(t *testing.T) {
	input := "*4\r\n$3\r\nSET\r\n$5\r\nmykey\r\n$5\r\nvalue\r\n$2\r\nXX\r\n"
	reader := bufio.NewReader(strings.NewReader(input))
	resp := &RESP{}
	req, err := resp.Parse(reader)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if req.cmd != "set" {
		t.Errorf("Expected cmd 'set', got '%s'", req.cmd)
	}
	if req.argsLen != 4 {
		t.Errorf("Expected argsLen 4, got %d", req.argsLen)
	}
	if strings.ToUpper(req.args[3]) != "XX" {
		t.Errorf("Expected XX flag, got %v", req.args[3])
	}
}

func TestParseSetxEX(t *testing.T) {
	input := "*5\r\n$3\r\nSET\r\n$5\r\nmykey\r\n$5\r\nvalue\r\n$2\r\nEX\r\n$2\r\n10\r\n"
	reader := bufio.NewReader(strings.NewReader(input))
	resp := &RESP{}
	req, err := resp.Parse(reader)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if req.cmd != "set" {
		t.Errorf("Expected cmd 'set', got '%s'", req.cmd)
	}
	if req.argsLen != 5 {
		t.Errorf("Expected argsLen 5, got %d", req.argsLen)
	}
	if strings.ToUpper(req.args[3]) != "EX" || req.args[4] != "10" {
		t.Errorf("Expected EX 10, got %v %v", req.args[3], req.args[4])
	}
}

func TestParseSetxPX(t *testing.T) {
	input := "*5\r\n$3\r\nSET\r\n$5\r\nmykey\r\n$5\r\nvalue\r\n$2\r\nPX\r\n$3\r\n500\r\n"
	reader := bufio.NewReader(strings.NewReader(input))
	resp := &RESP{}
	req, err := resp.Parse(reader)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if req.cmd != "set" {
		t.Errorf("Expected cmd 'set', got '%s'", req.cmd)
	}
	if req.argsLen != 5 {
		t.Errorf("Expected argsLen 5, got %d", req.argsLen)
	}
	if strings.ToUpper(req.args[3]) != "PX" || req.args[4] != "500" {
		t.Errorf("Expected PX 500, got %v %v", req.args[3], req.args[4])
	}
}

func TestParseSetxKEEPTTL(t *testing.T) {
	input := "*4\r\n$3\r\nSET\r\n$5\r\nmykey\r\n$5\r\nvalue\r\n$7\r\nKEEPTTL\r\n"
	reader := bufio.NewReader(strings.NewReader(input))
	resp := &RESP{}
	req, err := resp.Parse(reader)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if req.cmd != "set" {
		t.Errorf("Expected cmd 'set', got '%s'", req.cmd)
	}
	if req.argsLen != 4 {
		t.Errorf("Expected argsLen 4, got %d", req.argsLen)
	}
	if strings.ToUpper(req.args[3]) != "KEEPTTL" {
		t.Errorf("Expected KEEPTTL flag, got %v", req.args[3])
	}
}

func TestParseSetxGET(t *testing.T) {
	input := "*4\r\n$3\r\nSET\r\n$5\r\nmykey\r\n$5\r\nvalue\r\n$3\r\nGET\r\n"
	reader := bufio.NewReader(strings.NewReader(input))
	resp := &RESP{}
	req, err := resp.Parse(reader)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if req.cmd != "set" {
		t.Errorf("Expected cmd 'set', got '%s'", req.cmd)
	}
	if req.argsLen != 4 {
		t.Errorf("Expected argsLen 4, got %d", req.argsLen)
	}
	if strings.ToUpper(req.args[3]) != "GET" {
		t.Errorf("Expected GET flag, got %v", req.args[3])
	}
}

func TestParseRESPCommands(t *testing.T) {
	tests := []struct {
		input       string
		wantCmd     string
		wantArgs    []string
		wantArgsLen int
		wantErr     bool
	}{
		{
			"*3\r\n$3\r\nSET\r\n$5\r\nmykey\r\n$5\r\nvalue\r\n",
			"set", []string{"set", "mykey", "value"}, 3, false,
		},
		{
			"*1\r\n$4\r\nPING\r\n",
			"ping", []string{"ping"}, 1, false,
		},
		{
			"*2\r\n$3\r\nGET\r\n$5\r\nmykey\r\n",
			"get", []string{"get", "mykey"}, 2, false,
		},
		{
			"*2\r\n$3\r\nDEL\r\n$5\r\nmykey\r\n",
			"del", []string{"del", "mykey"}, 2, false,
		},
		{
			"*2\r\n$4\r\nINCR\r\n$3\r\nctr\r\n",
			"incr", []string{"incr", "ctr"}, 2, false,
		},
		{
			"*2\r\n$6\r\nEXISTS\r\n$5\r\nmykey\r\n",
			"exists", []string{"exists", "mykey"}, 2, false,
		},
		// {
		// 	"*2\r\n$4\r\nTYPE\r\n$5\r\nmykey\r\n",
		// 	"type", []string{"type", "mykey"}, 2, false,
		// },
	}

	for _, tt := range tests {
		reader := bufio.NewReader(strings.NewReader(tt.input))
		resp := &RESP{}
		req, err := resp.Parse(reader)
		if tt.wantErr {
			if err == nil {
				t.Errorf("Expected error for input %q, got nil", tt.input)
			}
			continue
		}
		if err != nil {
			t.Errorf("Unexpected error for input %q: %v", tt.input, err)
			continue
		}
		if req.cmd != strings.ToLower(tt.wantCmd) {
			t.Errorf("Expected cmd %q, got %q", tt.wantCmd, req.cmd)
		}
		if req.argsLen != tt.wantArgsLen {
			t.Errorf("Expected argsLen %d, got %d", tt.wantArgsLen, req.argsLen)
		}
		if len(req.args) != len(tt.wantArgs) {
			t.Errorf("Expected args %v, got %v", tt.wantArgs, req.args)
		} else {
			for i := range tt.wantArgs {
				got := req.args[i]
				if i == 0 {
					got = strings.ToLower(got)
				}
				if got != tt.wantArgs[i] {
					t.Errorf("Expected arg[%d] %q, got %q", i, tt.wantArgs[i], req.args[i])
				}
			}
		}
	}
}
