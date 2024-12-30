package raybotclient_test

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	raybotclient "github.com/tuanvumaihuynh/roboflow/internal/ws/raybot_client"
)

func TestInboundPublishMsg_Operation(t *testing.T) {
	m := raybotclient.InboundPublishMsg{}
	assert.Equal(t, raybotclient.OperationPublish, m.Operation())
}

func TestInboundResponseMsg_Operation(t *testing.T) {
	m := raybotclient.InboundResponseMsg{}
	assert.Equal(t, raybotclient.OperationResponse, m.Operation())
}

func TestUnmarshalInboundPublishMsg(t *testing.T) {
	testCases := []struct {
		name      string
		data      []byte
		want      raybotclient.InboundPublishMsg
		shouldErr bool
	}{
		{
			name: "Valid publish message",
			data: []byte(`{"op":"publish","topic":"status","data":{}}`),
			want: raybotclient.InboundPublishMsg{
				Topic: raybotclient.TopicStatus,
				Data:  []byte(`{}`),
			},
			shouldErr: false,
		},
		{
			name:      "Invalid topic",
			data:      []byte(`{"op":"publish","topic":"invalid","data":{}}`),
			want:      raybotclient.InboundPublishMsg{},
			shouldErr: true,
		},
		{
			name:      "Missing topic",
			data:      []byte(`{"op":"publish","data":{}}`),
			want:      raybotclient.InboundPublishMsg{},
			shouldErr: true,
		},
		{
			name:      "Invalid JSON",
			data:      []byte(`{`),
			want:      raybotclient.InboundPublishMsg{},
			shouldErr: true,
		},
		{
			name:      "Empty message",
			data:      []byte(``),
			want:      raybotclient.InboundPublishMsg{},
			shouldErr: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			var got raybotclient.InboundPublishMsg
			err := json.Unmarshal(tt.data, &got)

			if (err != nil) != tt.shouldErr {
				t.Fatalf("Unexpected error status: got err=%v, want shouldErr=%v", err, tt.shouldErr)
			}

			if err != nil {
				return // Skip further checks if an error occurred
			}

			assert.Equal(t, tt.want.Topic, got.Topic)
			assert.JSONEq(t, string(tt.want.Data), string(got.Data)) // Compare JSON content
		})
	}
}

func TestUnmarshalInboundResponseMsg(t *testing.T) {
	testCases := []struct {
		name      string
		data      []byte
		want      raybotclient.InboundResponseMsg
		shouldErr bool
	}{
		{
			name: "Valid response message",
			data: []byte(`{"op":"response","id":"aba18adb-428a-4702-be7a-33bfc6301faf","data":{},"status":"IN_PROGRESS"}`),
			want: raybotclient.InboundResponseMsg{
				ID:     uuid.MustParse("aba18adb-428a-4702-be7a-33bfc6301faf"),
				Data:   []byte(`{}`),
				Status: raybotclient.CommandStatusInProgress,
			},
			shouldErr: false,
		},
		{
			name:      "Invalid status",
			data:      []byte(`{"op":"response","id":"aba18adb-428a-4702-be7a-33bfc6301faf","data":{},"status":"invalid"}`),
			want:      raybotclient.InboundResponseMsg{},
			shouldErr: true,
		},
		{
			name:      "Missing status",
			data:      []byte(`{"op":"response","id":"aba18adb-428a-4702-be7a-33bfc6301faf","data":{}}`),
			want:      raybotclient.InboundResponseMsg{},
			shouldErr: true,
		},
		{
			name:      "Invalid UUID",
			data:      []byte(`{"op":"response","id":"invalid","data":{},"status":"PENDING"}`),
			want:      raybotclient.InboundResponseMsg{},
			shouldErr: true,
		},
		{
			name:      "Invalid JSON",
			data:      []byte(`{`),
			want:      raybotclient.InboundResponseMsg{},
			shouldErr: true,
		},
		{
			name:      "Empty message",
			data:      []byte(``),
			want:      raybotclient.InboundResponseMsg{},
			shouldErr: true,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			var got raybotclient.InboundResponseMsg
			err := json.Unmarshal(tt.data, &got)

			if (err != nil) != tt.shouldErr {
				t.Fatalf("Unexpected error status: got err=%v, want shouldErr=%v", err, tt.shouldErr)
			}

			if err != nil {
				return // Skip further checks if an error occurred
			}

			assert.Equal(t, tt.want.ID, got.ID)
			assert.JSONEq(t, string(tt.want.Data), string(got.Data)) // Compare JSON content
			assert.Equal(t, tt.want.Status, got.Status)
		})
	}
}
