package testing

import (
	"encoding/json"
	"pixel_clash/api/ws/types"
	"pixel_clash/ctypes"
	"testing"


	"github.com/stretchr/testify/assert"
)
func CheckTwoPlayersJoinedMessage(t *testing.T, cli *Cli) {
	for range 3 {
		msg, err := cli.Receive()
		if err != nil {
			t.Fatal(err)
		}
		if msg.Type == "response" {
			byteResponse, err := json.Marshal(msg.Data)
			if err != nil {
				t.Fatalf("marshalling server response: %s", err)
			}
			var resp types.ServerResponse
			if err := json.Unmarshal(byteResponse, &resp); err != nil {
				t.Fatalf("unmarshalling server response: %s", err)
			}
			assert.Equal(t, *(cli.LastReq), resp.ID)
		} else if msg.Type == "event" {
			byteEvent, err := json.Marshal(msg.Data)
			if err != nil {
				t.Fatalf("marshalling server response: %s", err)
			}
			var event ctypes.ServerEvent
			if err := json.Unmarshal(byteEvent, &event); err != nil {
				t.Fatalf("unmarshalling server response: %s", err)
			}
			if event.Type == "find_game_response" {
				t.Logf("ok\n")
			} else if event.Type == "waiting_change" {
				 t.Logf("ok\n")
			} else {
				t.Fatalf("wrong event type: %s", event.Type)
			}
		} else {
			t.Fatalf("unknown message type: %s", msg)
		}
	}
}
	

func TestGameFlow(t *testing.T) {
	URL := "ws://localhost:8080/user/join"
	cli1, err := NewCli(
		"kirill",
		ctypes.Game{
			Size:            2,
			FieldSize:       10,
			Time:           30,
			ThresholdSqare: 4,
		},
		URL,
	)
	if err != nil {
		t.Fatalf("creating client: %s", err)
	}
	defer cli1.Conn.Close()

	cli2, err := NewCli(
		"anya",
		ctypes.Game{
			Size:            2,
			FieldSize:       10,
			Time:           30,
			ThresholdSqare: 4,
		},
		URL,
	)
	if err != nil {
		t.Fatalf("creating client: %s", err)
	}
	defer cli2.Conn.Close()

	//find_game
	
	if cli1.LastReq, err = cli1.Send(types.FindGameRequest{Nickname: cli1.Nickname, GameType: cli1.GamePref}, "find_game"); err != nil {
		t.Fatalf("finding game: %s", err)
	}
	if cli2.LastReq, err = cli2.Send(types.FindGameRequest{Nickname: cli2.Nickname, GameType: cli2.GamePref}, "find_game"); err != nil {
		t.Fatalf("finding game: %s", err)
	}
	CheckTwoPlayersJoinedMessage(t, cli1)
	CheckTwoPlayersJoinedMessage(t, cli2)
}