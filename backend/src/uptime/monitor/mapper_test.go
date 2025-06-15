package monitors

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMapMonitorType(t *testing.T) {
	type args struct {
		typeTag MonitorConfigType
	}
	tests := []struct {
		name string
		args args
		want IMonitor
	}{
		{
			name: "Valid Http Monitor Type",
			args: args{typeTag: Http},
			want: &HttpMonitor{},
		},
		{
			name: "Valid Ping Monitor Type",
			args: args{typeTag: Ping},
			want: &PingMonitor{},
		},
		{
			name: "Invalid Monitor Type",
			args: args{typeTag: "InvalidType"},
			want: nil,
		},
		{
			name: "Empty Monitor Type",
			args: args{typeTag: ""},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, MapMonitorType(tt.args.typeTag), "MapMonitorType(%v)", tt.args.typeTag)
		})
	}
}
