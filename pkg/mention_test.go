package mention

import (
	"testing"
)

func Test_mention(t *testing.T) {
	type args struct {
		cmd  string
		text string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "echo",
			args: args{cmd: "echo", text: "hoge"},
			want: "hoge",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mention(tt.args.cmd, tt.args.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("mention() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("mention() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_set(t *testing.T) {
	type args struct {
		text    string
		getText string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr error
	}{
		{
			name: "Set and Get",
			args: args{text: "hoge fuga", getText: "hoge"},
			want: "fuga",
		},
		{
			name:    "Set and Get notfound",
			args:    args{text: "hoge fuga", getText: "fuga"},
			want:    "",
			wantErr: &NotfoundError{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _ = set(tt.args.text)
			got, err := get(tt.args.getText)
			if err != nil {
				if err != tt.wantErr {
					t.Errorf("get() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
			if got != tt.want {
				t.Errorf("set() = %v, want %v", got, tt.want)
			}
		})
	}
}
