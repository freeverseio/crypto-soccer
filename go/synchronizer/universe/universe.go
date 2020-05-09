package universe

import (
	"crypto/sha256"
	"errors"
	"sort"

	"github.com/freeverseio/crypto-soccer/go/storage"
)

type Universe struct {
	players []storage.Player
}

func (b *Universe) Append(players ...storage.Player) error {
	for _, player := range players {
		if player.EncodedSkills == nil {
			return errors.New("encodedSkills is nil")
		}
		if player.EncodedState == nil {
			return errors.New("encodedState is nil")
		}
		if i := sort.Search(len(b.players), func(i int) bool {
			return b.players[i].PlayerId.Cmp(player.PlayerId) == 0
		}); i != len(b.players) {
			return errors.New("player already in universe")
		}

		b.players = append(b.players, player)
	}
	return nil
}

func (b *Universe) sort() {
	sort.Slice(b.players[:], func(i, j int) bool {
		return b.players[i].PlayerId.Cmp(b.players[j].PlayerId) == -1
	})
}

func (b *Universe) Hash() ([32]byte, error) {
	b.sort()

	var result [32]byte
	h := sha256.New()
	for _, player := range b.players {
		h.Write(player.EncodedSkills.Bytes())
		h.Write(player.EncodedState.Bytes())
	}
	hash := h.Sum(nil)
	if len(hash) != 32 {
		return result, errors.New("Hash is not 32 byte")
	}
	copy(result[:], hash)
	return result, nil
}

func (b Universe) Size() int {
	return len(b.players)
}
