package protocol

import (
	"testing"
	"time"

	"github.com/B-AJ-Amar/gokv/internal/store"
)

func TestSetxExpiry(t *testing.T) {
	mem := store.NewInMemoryStore()
	args := store.SetArgs{ExpType: store.ExpirePX, ExpVal: 50} // 50 ms expiry
	mem.Setx("expkey", []byte("val"), args)
	val, _ := mem.Get("expkey")
	if string(val) != "val" {
		t.Errorf("Expected value 'val', got %v", val)
	}
	time.Sleep(60 * time.Millisecond)
	val, _ = mem.Get("expkey")
	if val != nil {
		t.Errorf("Expected key to be expired, got %v", val)
	}
}

func TestSetxNX(t *testing.T) {
	mem := store.NewInMemoryStore()
	args := store.SetArgs{NX_XX: 1} // NX
	mem.Setx("nxkey", []byte("first"), args)
	mem.Setx("nxkey", []byte("second"), args)
	val, _ := mem.Get("nxkey")
	if string(val) != "first" {
		t.Errorf("NX failed, expected 'first', got %v", val)
	}
}

func TestSetxXX(t *testing.T) {
	mem := store.NewInMemoryStore()
	args := store.SetArgs{NX_XX: 2}                     // XX
	mem.Setx("xxkey", []byte("first"), store.SetArgs{}) // normal set
	mem.Setx("xxkey", []byte("second"), args)
	val, _ := mem.Get("xxkey")
	if string(val) != "second" {
		t.Errorf("XX failed, expected 'second', got %v", val)
	}
	// XX on non-existing key
	ret, _, _ := mem.Setx("notfound", []byte("fail"), args)
	if ret != 0 {
		t.Errorf("XX should not set non-existing key")
	}
}

func TestSetxGET(t *testing.T) {
	mem := store.NewInMemoryStore()
	mem.Setx("getkey", []byte("old"), store.SetArgs{})
	args := store.SetArgs{Get: true}
	_, oldVal, _ := mem.Setx("getkey", []byte("new"), args)
	if string(oldVal) != "old" {
		t.Errorf("GET param failed, expected old value 'old', got %v", oldVal)
	}
	val, _ := mem.Get("getkey")
	if string(val) != "new" {
		t.Errorf("GET param failed, expected new value 'new', got %v", val)
	}
}

func TestProcessSetAndGet(t *testing.T) {
	mem := store.NewInMemoryStore()
	resp := &RESP{}
	// SET key value
	setReq := &RESPReq{cmd: "set", argsLen: 3, args: []string{"set", "mykey", "value"}}
	setRes, err := resp.Process(setReq, &mem)
	if err != nil {
		t.Fatalf("Process SET failed: %v", err)
	}
	if setRes.msgType != SimpleRes || setRes.message != "OK" {
		t.Errorf("SET response incorrect: got type %d, msg %q", setRes.msgType, setRes.message)
	}

	// GET key
	getReq := &RESPReq{cmd: "get", argsLen: 2, args: []string{"get", "mykey"}}
	getRes, err := resp.Process(getReq, &mem)
	if err != nil {
		t.Fatalf("Process GET failed: %v", err)
	}
	if getRes.msgType != BulkStrRes || getRes.message != "value" {
		t.Errorf("GET response incorrect: got type %d, msg %q", getRes.msgType, getRes.message)
	}
}

func TestProcessPing(t *testing.T) {
	mem := store.NewInMemoryStore()
	resp := &RESP{}
	pingReq := &RESPReq{cmd: "ping", argsLen: 1, args: []string{"ping"}}
	pingRes, err := resp.Process(pingReq, &mem)
	if err != nil {
		t.Fatalf("Process PING failed: %v", err)
	}
	if pingRes.msgType != SimpleRes || pingRes.message != "PONG" {
		t.Errorf("PING response incorrect: got type %d, msg %q", pingRes.msgType, pingRes.message)
	}
}

func TestProcessUnknownCommand(t *testing.T) {
	mem := store.NewInMemoryStore()
	resp := &RESP{}
	req := &RESPReq{cmd: "unknown", argsLen: 1, args: []string{"unknown"}}
	res, err := resp.Process(req, &mem)
	if err != nil {
		t.Fatalf("Process unknown command failed: %v", err)
	}
	if res.msgType != ErrorRes || res.message != "ERR unknown command" {
		t.Errorf("Unknown command response incorrect: got type %d, msg %q", res.msgType, res.message)
	}
}
