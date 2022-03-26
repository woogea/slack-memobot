package mention

import (
	"testing"
)

func Test_echo(t *testing.T) {
	type args struct {
		m    *Mention
		text string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := echo(tt.args.m, tt.args.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("echo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("echo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_set(t *testing.T) {
	type args struct {
		m    *Mention
		text string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "set",
			args:    args{m: NewMention("hoge"), text: "hoge"},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := set(tt.args.m, tt.args.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("set() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_get(t *testing.T) {
	type args struct {
		m    *Mention
		text string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "get",
			args:    args{m: NewMention("hoge"), text: "hoge"},
			want:    "fuga",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := set(tt.args.m, "hoge fuga")
			if err != nil {
				t.Errorf("set() error = %v", err)
			}
			got, err := get(tt.args.m, tt.args.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_list(t *testing.T) {
	type args struct {
		m    *Mention
		text string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := list(tt.args.m, tt.args.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("list() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("list() = %v, want %v", got, tt.want)
			}
		})
	}
}
