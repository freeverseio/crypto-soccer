package useractions_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/relay/storage"
	"github.com/freeverseio/crypto-soccer/go/useractions"
)

func TestMarshal(t *testing.T) {
	var ua useractions.UserActions
	data, err := ua.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `{"tactics":[],"trainings":[]}` {
		t.Fatalf("Wrong %v", string(data))
	}
	ua.Tactics = append(ua.Tactics, storage.Tactic{Verse: 3, TeamID: "ciao"})
	ua.Trainings = append(ua.Trainings, storage.Training{Verse: 5, TeamID: "pippo"})
	data, err = ua.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `{"tactics":[{"verse":3,"team_id":"ciao","tactic_id":0,"shirt_0":0,"shirt_1":0,"shirt_2":0,"shirt_3":0,"shirt_4":0,"shirt_5":0,"shirt_6":0,"shirt_7":0,"shirt_8":0,"shirt_9":0,"shirt_10":0,"shirt_11":0,"shirt_12":0,"shirt_13":0,"extra_attack_1":false,"extra_attack_2":false,"extra_attack_3":false,"extra_attack_4":false,"extra_attack_5":false,"extra_attack_6":false,"extra_attack_7":false,"extra_attack_8":false,"extra_attack_9":false,"extra_attack_10":false}],"trainings":[{"verse":5,"team_id":"pippo","special_player_shirt":0,"goalkeepers_defence":0,"goalkeepers_speed":0,"goalkeepers_pass":0,"goalkeepers_shoot":0,"goalkeepers_endurance":0,"defenders_defence":0,"defenders_speed":0,"defenders_pass":0,"defenders_shoot":0,"defenders_endurance":0,"midfielders_defence":0,"midfielders_speed":0,"midfielders_pass":0,"midfielders_shoot":0,"midfielders_endurance":0,"attackers_defence":0,"attackers_speed":0,"attackers_pass":0,"attackers_shoot":0,"attackers_endurance":0,"special_player_defence":0,"special_player_speed":0,"special_player_pass":0,"special_player_shoot":0,"special_player_endurance":0}]}` {
		t.Fatalf("Wrong %v", string(data))
	}
}

func TestUnmarshal(t *testing.T) {
	var ua useractions.UserActions
	err := ua.Unmarshal([]byte(`{"tactics":[],"trainings":[]}`))
	if err != nil {
		t.Fatal(err)
	}
	if len(ua.Tactics) != 0 {
		t.Fatal("Tactics not empty")
	}
	if len(ua.Trainings) != 0 {
		t.Fatal("Training not empty")
	}
}

func TestIpfsPushAndPull(t *testing.T) {
	var ua useractions.UserActions
	tactic := storage.Tactic{}
	tactic.TeamID = "ciao"
	ua.Tactics = append(ua.Tactics, tactic)
	cif, err := ua.PushToIpfs("localhost:5001")
	if err != nil {
		t.Fatal(err)
	}
	if cif != "QmS8S6a3uesR2N4sYMc18yz5yP6Wcge84L2xWxQW2EVaMJ" {
		t.Fatalf("Wrong cif %v", cif)
	}
	training := storage.Training{}
	training.TeamID = "pippo"
	ua.Trainings = append(ua.Trainings, training)
	cif, err = ua.PushToIpfs("localhost:5001")
	if err != nil {
		t.Fatal(err)
	}
	if cif != "QmQgu5v92T8vD9xbxVP4tfyBEexKNDiTwgLPZYLZNhH75K" {
		t.Fatalf("Wrong cif %v", cif)
	}
	var ua2 useractions.UserActions
	err = ua2.PullFromIpfs("localhost:5001", cif)
	if err != nil {
		t.Fatal(err)
	}
	if !ua2.Equal(&ua) {
		t.Fatalf("Expected %v but %v", ua, ua2)
	}
}
