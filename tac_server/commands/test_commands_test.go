package commands

import (
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/AnmolKothari22/tac/tac_server/input_output"
	"github.com/AnmolKothari22/tac/tac_server/shards"
)

func TestCommands(t *testing.T) {
	tests := []struct {
		name     string
		command  func(*sync.RWMutex, []string, *map[string]shards.Stored_values) (string, inputoutput.WriteType)
		parts    []string
		initial  map[string]string
		wantMsg  string
		wantType inputoutput.WriteType
		wantMap  map[string]string
	}{
		{
			name: "set creates key", command: Handle_set, parts: []string{"set", "count", "4"},
			initial: map[string]string{}, wantMsg: "OK", wantType: inputoutput.Normal,
			wantMap: map[string]string{"count": "4"},
		},
		{
			name: "set replaces value", command: Handle_set, parts: []string{"set", "count", "5"},
			initial: map[string]string{"count": "4"}, wantMsg: "OK", wantType: inputoutput.Normal,
			wantMap: map[string]string{"count": "5"},
		},
		{
			name: "get existing key", command: Handle_get, parts: []string{"get", "count"},
			initial: map[string]string{"count": "4"}, wantMsg: "4", wantType: inputoutput.Normal,
			wantMap: map[string]string{"count": "4"},
		},
		{
			name: "get empty value", command: Handle_get, parts: []string{"get", "empty"},
			initial: map[string]string{"empty": ""}, wantMsg: "", wantType: inputoutput.Normal,
			wantMap: map[string]string{"empty": ""},
		},
		{
			name: "get missing key", command: Handle_get, parts: []string{"get", "missing"},
			initial: map[string]string{}, wantMsg: "", wantType: inputoutput.KeyMiss,
			wantMap: map[string]string{},
		},
		{
			name: "inc by default", command: Handle_inc, parts: []string{"inc", "count"},
			initial: map[string]string{"count": "4"}, wantMsg: "key: count increased by 1", wantType: inputoutput.Normal,
			wantMap: map[string]string{"count": "5"},
		},
		{
			name: "inc by explicit amount", command: Handle_inc, parts: []string{"inc", "count", "-2"},
			initial: map[string]string{"count": "4"}, wantMsg: "key: count increased by -2", wantType: inputoutput.Normal,
			wantMap: map[string]string{"count": "2"},
		},
		{
			name: "inc missing key", command: Handle_inc, parts: []string{"inc", "missing"},
			initial: map[string]string{}, wantMsg: "key not found", wantType: inputoutput.KeyMiss,
			wantMap: map[string]string{},
		},
		{
			name: "inc nonnumeric value", command: Handle_inc, parts: []string{"inc", "count"},
			initial: map[string]string{"count": "nope"}, wantMsg: "value is NAN", wantType: inputoutput.Error,
			wantMap: map[string]string{"count": "nope"},
		},
		{
			name: "inc nonnumeric amount", command: Handle_inc, parts: []string{"inc", "count", "nope"},
			initial: map[string]string{"count": "4"}, wantMsg: "NAN not allowed/out of bound number", wantType: inputoutput.Error,
			wantMap: map[string]string{"count": "4"},
		},
		{
			name: "inc overflow", command: Handle_inc, parts: []string{"inc", "count", "1"},
			initial: map[string]string{"count": "9223372036854775807"}, wantMsg: "value going out of bound ", wantType: inputoutput.Error,
			wantMap: map[string]string{"count": "9223372036854775807"},
		},
		{
			name: "dec by default", command: Handle_dec, parts: []string{"dec", "count"},
			initial: map[string]string{"count": "4"}, wantMsg: "key: count decreased by 1", wantType: inputoutput.Normal,
			wantMap: map[string]string{"count": "3"},
		},
		{
			name: "dec by explicit amount", command: Handle_dec, parts: []string{"dec", "count", "-2"},
			initial: map[string]string{"count": "4"}, wantMsg: "key: count decreased by -2", wantType: inputoutput.Normal,
			wantMap: map[string]string{"count": "6"},
		},
		{
			name: "dec missing key", command: Handle_dec, parts: []string{"dec", "missing"},
			initial: map[string]string{}, wantMsg: "key not found", wantType: inputoutput.KeyMiss,
			wantMap: map[string]string{},
		},
		{
			name: "dec nonnumeric value", command: Handle_dec, parts: []string{"dec", "count"},
			initial: map[string]string{"count": "nope"}, wantMsg: "value is NAN", wantType: inputoutput.Error,
			wantMap: map[string]string{"count": "nope"},
		},
		{
			name: "dec nonnumeric amount", command: Handle_dec, parts: []string{"dec", "count", "nope"},
			initial: map[string]string{"count": "4"}, wantMsg: "decreasing by a NAN not allowed", wantType: inputoutput.Error,
			wantMap: map[string]string{"count": "4"},
		},
		{
			name: "dec overflow", command: Handle_dec, parts: []string{"dec", "count", "1"},
			initial: map[string]string{"count": "-9223372036854775808"}, wantMsg: "value going out of bound ", wantType: inputoutput.Error,
			wantMap: map[string]string{"count": "-9223372036854775808"},
		},
		{
			name: "del existing key", command: Handle_del, parts: []string{"del", "count"},
			initial: map[string]string{"count": "4"}, wantMsg: "deleted key: count with value: 4", wantType: inputoutput.Normal,
			wantMap: map[string]string{},
		},
		{
			name: "del missing key", command: Handle_del, parts: []string{"del", "missing"},
			initial: map[string]string{}, wantMsg: "key not found", wantType: inputoutput.KeyMiss,
			wantMap: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lock := &sync.RWMutex{}
			store := make(map[string]shards.Stored_values, len(tt.initial))
			for key, value := range tt.initial {
				store[key] = shards.Stored_values{Value:value ,Created_at: time.Now(),Expire_after: 0}
			}

			msg, typ := tt.command(lock, tt.parts, &store)
			if msg != tt.wantMsg {
				t.Errorf("message = %q, want %q", msg, tt.wantMsg)
			}
			if typ != tt.wantType {
				t.Errorf("type = %v, want %v", typ, tt.wantType)
			}
			if !reflect.DeepEqual(store, tt.wantMap) {
				t.Errorf("store = %#v, want %#v", store, tt.wantMap)
			}
		})
	}
}
