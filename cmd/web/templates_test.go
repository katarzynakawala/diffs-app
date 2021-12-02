package main

import (
	"testing"
	"time"
)

func TestHumanDate(t *testing.T) {
	tests := []struct {
		name string
		tm time.Time
		want string
	}{
		{
			name: "UTC",
			tm: time.Date(2021, 12, 12, 11, 0, 0, 0, time.UTC),
			want: "12 Dec 2021 at 11:00",
		},
		{
			name: "Empty",
			tm: time.Time{},
			want: "",
		},
		{
			name: "CET",
			tm: time.Date(2021, 12, 12, 11, 0, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "12 Dec 2021 at 10:00",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.tm)

			if hd != tt.want {
				t.Errorf("want %q; got %q", tt.want, hd)
			}
		})
	}
}