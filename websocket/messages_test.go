package websocket

import "testing"

func TestDataUpdate_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Invalid unmarshal",
			args: args{
				data: []byte(""),
			},
			wantErr: true,
		}, {
			name: "Invalid array length",
			args: args{
				data: []byte("[]"),
			},
			wantErr: true,
		}, {
			name: "Invalid channel_id",
			args: args{
				data: []byte("[\"asdas\", [], \"\", \"\"]"),
			},
			wantErr: true,
		}, {
			name: "Invalid channel_name",
			args: args{
				data: []byte("[1, [], 123, \"\"]"),
			},
			wantErr: true,
		}, {
			name: "Invalid pair",
			args: args{
				data: []byte("[1, [], \"trades\", 123]"),
			},
			wantErr: true,
		}, {
			name: "Invalid pair",
			args: args{
				data: []byte("[1, [], \"trades\", \"XBT/USD\"]"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &DataUpdate{}
			if err := u.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("DataUpdate.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
