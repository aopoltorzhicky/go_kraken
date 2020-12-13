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
			name: "Length 4, Invalid channel_id",
			args: args{
				data: []byte("[\"asdas\", [], \"\", \"\"]"),
			},
			wantErr: true,
		}, {
			name: "Length 4, Invalid channel_name",
			args: args{
				data: []byte("[1, [], 123, \"\"]"),
			},
			wantErr: true,
		}, {
			name: "Length 4, Invalid pair",
			args: args{
				data: []byte("[1, [], \"trades\", 123]"),
			},
			wantErr: true,
		}, {
			name: "Length 4, Valid data",
			args: args{
				data: []byte("[1, [], \"trades\", \"XBT/USD\"]"),
			},
			wantErr: false,
		}, {
			name: "Length 3, invalid channel name",
			args: args{
				data: []byte(`[1, 123, {}]`),
			},
			wantErr: true,
		},
		{
			name: "Length 3, invalid sequence number kind",
			args: args{
				data: []byte(`[1, "", ""]`),
			},
			wantErr: true,
		},
		{
			name: "Length 3, sequence missing from object",
			args: args{
				data: []byte(`[1, "", {}]`),
			},
			wantErr: true,
		},
		{
			name: "Length 3, sequence not a number",
			args: args{
				data: []byte(`[1, "", {"sequence": "foo"}]`),
			},
			wantErr: true,
		},
		{
			name: "Length 3, valid data",
			args: args{
				data: []byte(`[1, "", {"sequence": 456}]`),
			},
			wantErr: false,
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
