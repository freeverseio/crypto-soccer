package process

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go-synchronizer/storage"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/testutils"
)

func TestSyncTeamWithNoTeam(t *testing.T) {
	storage, err := storage.NewSqlite3("../sql/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	blockchain := testutils.DefaultSimulatedBlockchain()

	p := NewEventProcessor(nil, storage, blockchain.Assets)
	p.Process()

	count, err := storage.TeamCount()
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("Expected 0 received %v", count)
	}
}

func TestSyncTeams(t *testing.T) {
	storage, err := storage.NewSqlite3("../sql/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	ganache := testutils.NewGanache()

	owner := ganache.CreateAccountWithBalance("1000000000000000000") // 1 eth
	ganache.DeployContracts(owner)

	alice := ganache.CreateAccountWithBalance("50000000000000000000") // 50 eth
	bob := ganache.CreateAccountWithBalance("50000000000000000000")   // 50 eth
	carol := ganache.CreateAccountWithBalance("50000000000000000000") // 50 eth

	ganache.CreateTeam("A", alice)
	ganache.CreateTeam("B", bob)
	ganache.CreateTeam("C", carol)

	p := NewEventProcessor(ganache.Client, storage, ganache.Assets)

	if err := p.Process(); err != nil {
		t.Fatal(err)
	} else {
		if count, err := storage.TeamCount(); err != nil {
			t.Fatal(err)
		} else if count != 3 {
			t.Fatalf("Expected 3 actual %v", count)
		}
	}

	if team, err := storage.GetTeam(1); err != nil {
		t.Fatal(err)
	} else if team.Id != 1 {
		t.Fatalf("Expected 1 result %v", team.Id)
	} else if team.Name != "A" {
		t.Fatalf("Expected A result %v", team.Name)
	}
	if team, err := storage.GetTeam(2); err != nil {
		t.Fatal(err)
	} else if team.Id != 2 {
		t.Fatalf("Expected 2 result %v", team.Id)
	} else if team.Name != "B" {
		t.Fatalf("Expected B result %v", team.Name)
	}
	if team, err := storage.GetTeam(3); err != nil {
		t.Fatal(err)
	} else if team.Id != 3 {
		t.Fatalf("Expected 3 result %v", team.Id)
	} else if team.Name != "C" {
		t.Fatalf("Expected C result %v", team.Name)
	}

	if count, err := storage.PlayerCount(); err != nil {
		t.Fatal(err)
	} else if count != 33 {
		t.Fatalf("Expected 33 players actual %v", count)
	} else {
		expected := [...]string{
			"233252715009006782633031027761500368726392237034442587091272386971883798528",
			"318054903599544818898214922780748669125244523168847882546880797338668892160",
			"402870895946060375023378743329593182812670531289755665070483278784440041472",
			"233246676059966483839595065386808619099667612920538174708039849198961557504",
			"473539170526046530280254058813683500460702397139160344923488436430435254272",
			"353386237041947793367532468205453061269991757550625035325248317181940203520",
			"353403491144493866664180862555547597561259459502628276141554174390970613760",
			"381664848594113858985454672333068162950546717863300999977021095927569448960",
			"339255342677080671879386431565238714729508219307736333831860595908080041984",
			"537134018345346822328115811322903124789533758363876423880129798732146278400",
			"268579735070596519652329684924137907766567486688240860110918941462765764608",
			"416991437533458158185997477881027911304220802259315242624379580083720945664",
			"494738747306414473062544110739785544475727072322825803303183958551302242304",
			"254450134414545871673985320095152824883805746524231490396204908308238172160",
			"402859680318657323977431378886472957997857225827425868984632256705733853184",
			"466467899639685404603524672484854437359502225768479130564542021787549958144",
			"353397021204888665144933922110153134848076583622701597924415614305903312896",
			"318051884111865453612258873350753550866316842529560175747914685824590938112",
			"374597029054879643750246933904075235084110991542703174622764266855627489280",
			"318047139529848269325145452534659898585261817658510741867660537689463062528",
			"296842386597626004951478481656212150643240820863734176334143511589872992256",
			"360451899864503556458597564152707025212461613237699331451786637633889763328",
			"466472644774545519732570037537572124505955162208660125056994662042621181952",
			"332190974063597261619836509839300883578614413447835399513192663733598420992",
			"544207876781071242157545429376734050411565992255470857669831199671883661312",
			"424061845569893605341842180816454807244534398466905886045825050311665909760",
			"416991869471798331918068953343162830685957515668744295988500680935187939328",
			"247384471486711560946595821962064970542712520200194786485455597980671279104",
			"353393139578112664333854928854103711488362586963287637607325663028039385088",
			"325126605344846551604001339464582372317097307176268433699287022698345005056",
			"409936558110407890278455288943505078640571276833081245331861774255800713216",
			"282726589460611469340748942310917476626701578513914952488447029373646667776",
			"226171092345702538835181272959923482695082141611923362697371995870641258496",
		}
		for i, player := range expected {
			if result, err := storage.GetPlayer(uint64(i + 1)); err != nil {
				t.Fatal(err)
			} else if result.State != player {
				//t.Fatalf("Expecting player state to be %v actual %v", player, result.State.String())
			}
		}
	}

	ganache.CreateTeam("D", alice)
	if err := p.Process(); err != nil {
		t.Fatal(err)
	} else {
		if count, err := storage.TeamCount(); err != nil {
			t.Fatal(err)
		} else if count != 4 {
			t.Fatalf("Expected 4 actual %v", count)
		}
	}
	if team, err := storage.GetTeam(4); err != nil {
		t.Fatal(err)
	} else if team.Id != 4 {
		t.Fatalf("Expected 4 result %v", team.Id)
	} else if team.Name != "D" {
		t.Fatalf("Expected D result %v", team.Name)
	}

	if count, err := storage.PlayerCount(); err != nil {
		t.Fatal(err)
	} else if count != 44 {
		t.Fatalf("Expected 44 players actual %v", count)
	}
}
