package setting

import (
	"log"
	"testing"
)

func TestSetup(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Test Case 1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Setup()
			log.Printf("ServerSetting.HttpPort : %v", ServerSetting.HttpPort)
			log.Printf("ServerSetting.ReadTimeout : %v", ServerSetting.ReadTimeout)
			log.Printf("JWTSetting.AccessExpireTime : %v", JWTSetting.AccessExpireTime)
			log.Printf("JWTSetting.RefreshExpireTime : %v", JWTSetting.RefreshExpireTime)
			log.Printf("RedisSetting.Host : %v", RedisSetting.Host)
		})
	}
}
