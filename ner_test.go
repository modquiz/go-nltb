package nltb

import (
	"reflect"
	"testing"
)

func TestNER_Init(t *testing.T) {
	tests := []struct {
		name string
		n    *NER
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.n.Init()
		})
	}
}

func TestNER_Do(t *testing.T) {
	type args struct {
		byteString []byte
	}
	tests := []struct {
		name string
		p    *NER
		args args
		want []TaggedWord
	}{
		{
			name: "Test 1",
			args: "This is a test of the TestNER",
			want: []TaggedWord{},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.Do(tt.args.byteString); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NER.Do() = %v, want %v", got, tt.want)
			}
		})
	}
}
