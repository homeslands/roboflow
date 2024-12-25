package raybotclient_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/tuanvumaihuynh/roboflow/internal/model"
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

func TestUnmarshalInboundMsg(t *testing.T) {
	testCases := []struct {
		name      string
		data      []byte
		want      raybotclient.InboundMsg
		shouldErr bool
	}{
		{
			name: "Valid publish message",
			data: []byte(`{"op":"publish","topic":"STOP","data":{}}`),
			want: &raybotclient.InboundPublishMsg{
				Topic: model.RaybotCommandTypeStop,
				Data:  []byte(`{}`),
			},
			shouldErr: false,
		},
		{
			name: "Valid response message",
			data: []byte(`{"op":"response","id":"aba18adb-428a-4702-be7a-33bfc6301faf","data":{}}`),
			want: &raybotclient.InboundResponseMsg{
				ID:   uuid.MustParse("aba18adb-428a-4702-be7a-33bfc6301faf"),
				Data: []byte(`{}`),
			},
			shouldErr: false,
		},
		{
			name:      "Invalid operation",
			data:      []byte(`{"op":"invalid","id":"123","data":{}}`),
			want:      nil,
			shouldErr: true,
		},
		{
			name:      "Missing operation",
			data:      []byte(`{"id":"123","data":{}}`),
			want:      nil,
			shouldErr: true,
		},
		{
			name:      "Invalid JSON",
			data:      []byte(`{`),
			want:      nil,
			shouldErr: true,
		},
		{
			name:      "Empty message",
			data:      []byte(``),
			want:      nil,
			shouldErr: true,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := raybotclient.UnmarshalInboundMsg(tt.data)

			if (err != nil) != tt.shouldErr {
				t.Fatalf("Unexpected error status: got err=%v, want shouldErr=%v", err, tt.shouldErr)
			}

			if err != nil {
				return // Skip further checks if an error occurred
			}

			// Check types match
			assert.IsType(t, tt.want, got)

			// Compare field by field based on the concrete type
			switch expected := tt.want.(type) {
			case *raybotclient.InboundPublishMsg:
				actual, ok := got.(*raybotclient.InboundPublishMsg)
				if !ok {
					t.Fatalf("Expected type *InboundPublishMsg, but got %T", got)
				}
				assert.Equal(t, expected.Topic, actual.Topic)
				assert.JSONEq(t, string(expected.Data), string(actual.Data)) // Compare JSON content

			case *raybotclient.InboundResponseMsg:
				actual, ok := got.(*raybotclient.InboundResponseMsg)
				if !ok {
					t.Fatalf("Expected type *InboundResponseMsg, but got %T", got)
				}
				assert.Equal(t, expected.ID, actual.ID)
				assert.JSONEq(t, string(expected.Data), string(actual.Data)) // Compare JSON content

			default:
				t.Fatalf("Unhandled type %T in test case", tt.want)
			}
		})
	}
}
