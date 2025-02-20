package jsoncompressor

import (
	"testing"
)

func TestCompress(t *testing.T) {
	{
		type MyStruct struct {
			K1 int    `json:"k1"`
			K2 string `json:"k2"`
			K3 []int  `json:"k3"`
			K5 struct {
				K6 int    `json:"k6"`
				K8 string `json:"k8"`
			} `json:"k5"`
		}

		data := MyStruct{
			K1: 1,
			K2: "2",
			K3: []int{3, 4},
			K5: struct {
				K6 int    `json:"k6"`
				K8 string `json:"k8"`
			}{
				K6: 7,
				K8: "8",
			},
		}
		raw, err := Marshal(data)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}
		target := `[1,"2",[3,4],[7,"8"]]`
		if string(raw) != target {
			t.Fatalf("Marshaled data is invalid, got: %s", raw)
		}
	}
	{
		raw, err := Marshal(1)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}
		target := `1`
		if string(raw) != target {
			t.Fatalf("Marshaled data is invalid, got: %s", raw)
		}
	}
	{
		raw, err := Marshal([]string{"test"})
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}
		target := `["test"]`
		if string(raw) != target {
			t.Fatalf("Marshaled data is invalid, got: %s", raw)
		}
	}

}
