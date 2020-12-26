//nolint:funlen,scopelint,testpackage
package modelzone

import (
	"reflect"
	"testing"
)

func TestZone_Diff(t *testing.T) {
	type args struct {
		secondaryzones Zonestatemap
	}

	tests := []struct {
		name             string
		primaryzones     Zonestatemap
		args             args
		wantAddedzones   Zonestatemap
		wantDeletedzones Zonestatemap
		wantChangedzones Zonestatemap
	}{
		{
			name:             "OK-identical",
			primaryzones:     map[string]int32{"example.org.": 1234, "example.eu.": 4321},
			args:             args{map[string]int32{"example.org.": 1234, "example.eu.": 4321}},
			wantAddedzones:   map[string]int32{},
			wantDeletedzones: map[string]int32{},
			wantChangedzones: map[string]int32{},
		},
		{
			name:             "OK-identical-different-order",
			primaryzones:     map[string]int32{"example.org.": 1234, "example.eu.": 4321},
			args:             args{map[string]int32{"example.eu.": 4321, "example.org.": 1234}},
			wantAddedzones:   map[string]int32{},
			wantDeletedzones: map[string]int32{},
			wantChangedzones: map[string]int32{},
		},
		{
			name:             "OK-added-zone",
			primaryzones:     map[string]int32{"example.org.": 1234, "example.eu.": 4321, "test.org.": 1324},
			args:             args{map[string]int32{"example.eu.": 4321, "example.org.": 1234}},
			wantAddedzones:   map[string]int32{"test.org.": 1324},
			wantDeletedzones: map[string]int32{},
			wantChangedzones: map[string]int32{},
		},
		{
			name:             "OK-deleted-zone",
			primaryzones:     map[string]int32{"example.org.": 1234, "example.eu.": 4321},
			args:             args{map[string]int32{"example.eu.": 4321, "example.org.": 1234, "test.org.": 1324}},
			wantAddedzones:   map[string]int32{},
			wantDeletedzones: map[string]int32{"test.org.": 1324},
			wantChangedzones: map[string]int32{},
		},
		{
			name:             "OK-changed-zone",
			primaryzones:     map[string]int32{"example.org.": 1234, "example.eu.": 4321, "test.org.": 1325},
			args:             args{map[string]int32{"example.eu.": 4321, "example.org.": 1234, "test.org.": 1324}},
			wantAddedzones:   map[string]int32{},
			wantDeletedzones: map[string]int32{},
			wantChangedzones: map[string]int32{"test.org.": 1325},
		},
		{
			name:             "OK-all",
			primaryzones:     map[string]int32{"example.org.": 1234, "example.eu.": 4321, "test.org.": 1325},
			args:             args{map[string]int32{"example.eu.": 4321, "test.org.": 1324, "test.eu.": 4231}},
			wantAddedzones:   map[string]int32{"example.org.": 1234},
			wantDeletedzones: map[string]int32{"test.eu.": 4231},
			wantChangedzones: map[string]int32{"test.org.": 1325},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAddedzones, gotDeletedzones, gotChangedzones := tt.primaryzones.Diff(tt.args.secondaryzones)
			if !reflect.DeepEqual(gotAddedzones, tt.wantAddedzones) {
				t.Errorf("Diff() gotAddedzones = %v, want %v", gotAddedzones, tt.wantAddedzones)
			}
			if !reflect.DeepEqual(gotDeletedzones, tt.wantDeletedzones) {
				t.Errorf("Diff() gotDeletedzones = %v, want %v", gotDeletedzones, tt.wantDeletedzones)
			}
			if !reflect.DeepEqual(gotChangedzones, tt.wantChangedzones) {
				t.Errorf("Diff() gotChangedzones = %v, want %v", gotChangedzones, tt.wantChangedzones)
			}
		})
	}
}
