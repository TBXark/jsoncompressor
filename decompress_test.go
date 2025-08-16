package jsoncompressor

import (
	"testing"
)

func TestDecompress(t *testing.T) {

	{
		type MyStruct struct {
			K1 int    `json:"k1"`
			K2 string `json:"k2"`
			K3 []int  `json:"k3"`
			K5 *struct {
				K6 int    `json:"k6"`
				K8 string `json:"k8"`
			} `json:"k5"`
		}

		raw := []byte(`[1,"2",[3,4],[7,"8"]]`)
		var data MyStruct
		err := Unmarshal(raw, &data)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}
		t.Logf("Unmarshaled: %+v", data)
		if data.K1 != 1 || data.K2 != "2" || len(data.K3) != 2 || data.K3[0] != 3 || data.K3[1] != 4 || data.K5.K6 != 7 || data.K5.K8 != "8" {
			t.Fatalf("Unmarshaled data is invalid")
		}

		dataV2 := &MyStruct{}
		err = Unmarshal(raw, &dataV2)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}
		t.Logf("Unmarshaled: %+v", dataV2)
	}
	{
		raw := []byte(`1`)
		var data int
		err := Unmarshal(raw, &data)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}
		t.Logf("Unmarshaled: %+v", data)
	}
	{
		raw := []byte(`["test"]`)
		var data []string
		err := Unmarshal(raw, &data)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}
		t.Logf("Unmarshaled: %v", data)
	}
}

func TestDecompressTagWithoutNameAndLengthMismatch(t *testing.T) {
	type S struct {
		A int `json:",omitempty"`
		B int `json:"b"`
	}
	{
		raw := []byte(`[10,20]`)
		var s S
		if err := Unmarshal(raw, &s); err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}
		if s.A != 10 || s.B != 20 {
			t.Fatalf("Unexpected values: %+v", s)
		}
	}
	{
		raw := []byte(`[1,2,3]`)
		var s S
		if err := Unmarshal(raw, &s); err == nil {
			t.Fatalf("expected error on length mismatch (too many)")
		}
	}
	{
		raw := []byte(`[1]`)
		var s S
		if err := Unmarshal(raw, &s); err == nil {
			t.Fatalf("expected error on length mismatch (too few)")
		}
	}
}
