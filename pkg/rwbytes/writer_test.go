package rwbytes

import (
	"bytes"
	"testing"
)

func TestWriteBytes(t *testing.T) {
	type args struct {
		in     *bytes.Buffer
		fixLen int
		datas  []byte
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := WriteBytes(tt.args.in, tt.args.fixLen, tt.args.datas)
			if (err != nil) != tt.wantErr {
				t.Errorf("WriteBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("WriteBytes() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWriteInt(t *testing.T) {
	type args struct {
		in     *bytes.Buffer
		fixLen int
		data   int
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := WriteInt(tt.args.in, tt.args.fixLen, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("WriteInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("WriteInt() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWriteIntHex(t *testing.T) {

	type args struct {
		in     *bytes.Buffer
		fixLen int
		data   int
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"1", args{new(bytes.Buffer), 2, 1}, 2, false},
		{"2", args{new(bytes.Buffer), 2, 10}, 2, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := WriteIntHex(tt.args.in, tt.args.fixLen, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("WriteIntHex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("WriteIntHex() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWriteString(t *testing.T) {
	type args struct {
		in     *bytes.Buffer
		fixLen int
		data   string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := WriteString(tt.args.in, tt.args.fixLen, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("WriteString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("WriteString() got = %v, want %v", got, tt.want)
			}
		})
	}
}
