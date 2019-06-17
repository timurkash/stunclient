package main

import "testing"

func TestDoDial(t *testing.T){
	_, software, err := doDial("localhost:3479", "timur:realm:pass")
	if err != nil {
		t.Error("bad dial", err)
	}
	if software != "gortc/stund" {
		t.Error("bad software", software)
	}
	_, software, err = doDial("localhost:3479", "timur:realm:pass_")
	if err == nil {
		t.Error("should error")
	}
}

func BenchmarkMain(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _, err := doDial("localhost:3479", "timur:realm:pass")
		if err != nil {
			b.Fatal(err)
		}
	}
}