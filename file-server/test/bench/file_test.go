package bench

import "testing"

func BenchmarkFileAddFromFile(b *testing.B) {
	for n := 0; n < b.N; n++ {
		err := fileAddFromFile()
		if err != nil {
			b.Log(err)
		}
	}
}

func BenchmarkFileAddFromPipe(b *testing.B) {
	for n := 0; n < b.N; n++ {
		err := fileAddFromPipe()
		if err != nil {
			b.Log(err)
		}
	}
}
