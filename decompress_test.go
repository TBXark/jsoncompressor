package jsoncompressor

import (
	"testing"
)

func TestDecompress(t *testing.T) {

	type MyStruct struct {
		K1 int    `json:"k1"`
		K2 string `json:"k2"`
		K3 []int  `json:"k3"`
		K5 struct {
			K6 int    `json:"k6"`
			K8 string `json:"k8"`
		} `json:"k5"`
	}

	raw := []byte(`[1,"2",[3,4],[7,"8"]]`)
	var data MyStruct
	err := Unmarshal(raw, &data)
	if err != nil {
		t.Fatalf("Decompress failed: %v", err)
	}
	t.Logf("Decompressed: %v", data)
	if data.K1 != 1 || data.K2 != "2" || len(data.K3) != 2 || data.K3[0] != 3 || data.K3[1] != 4 || data.K5.K6 != 7 || data.K5.K8 != "8" {
		t.Fatalf("Decompressed data is invalid")
	}
}
