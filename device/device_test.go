package device

import "testing"

var (
	mTop = make(map[int64]float32)
	mBot = make(map[int64]float32)
)

func TestConstructReqAction(t *testing.T) {
	req := constructReq(mTop, mBot)
	if req.Action == "" {
		t.Error("Req is absent", req.Action)
	}
}
func TestConstructReqMetaMAC(t *testing.T) {
	req := constructReq(mTop, mBot)
	if req.Meta.MAC == "" {
		t.Error("MAC is absent", req.Meta.MAC)
	}
}
func TestConstructReqMetaType(t *testing.T) {
	req := constructReq(mTop, mBot)
	if req.Meta.Type == "" {
		t.Error("Type is absent", req.Meta.Type)
	}
}

func TestConstructReqMetaName(t *testing.T) {
	req := constructReq(mTop, mBot)
	if req.Meta.Name == "" {
		t.Error("Name is absent", req.Meta.Name)
	}
}

func TestConstructReqData(t *testing.T) {
	req := constructReq(mTop, mBot)
	if req.Data == nil {
		t.Error("Nil data", req.Data)
	}
}
