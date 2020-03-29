package gql

type TransferFirstBotToAddrInput struct {
	Timezone             int32
	CountryIdxInTimezone string
	Address              string
}

func (b *Resolver) TransferFirstBotToAddr(input TransferFirstBotToAddrInput) bool {
	if b.c != nil {
		b.c <- input
	}

	return true
}
