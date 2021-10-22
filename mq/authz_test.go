package mq

import (
	"github.com/igridnet/users/models"
	"testing"
)

func TestAuthorize(t *testing.T) {
	type args struct {
		node      models.Node
		topic     string
		operation int
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:    "",
			args:    args{
				node:      models.Node{
					UUID:    "temperature",
					Addr:    "",
					Key:     "password",
					Name:    "",
					Type:    int(models.ActuatorNode),
					Region:  "region",
					Latd:    "",
					Long:    "",
					Created: 0,
					Master:  "",
				},
				topic:     "region/temperaturet",
				operation: SubscribeOperation,
			},
			want:    false,
			wantErr:false,
		},
	}
		for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Authorize(tt.args.node, tt.args.topic, tt.args.operation)
			if (err != nil) != tt.wantErr {
				t.Errorf("Authorize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Authorize() got = %v, want %v", got, tt.want)
			}
		})
	}
}