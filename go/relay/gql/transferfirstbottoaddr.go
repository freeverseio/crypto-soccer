package gql

func (b *Resolver) TransferFirstBotToAddr(input TransferFirstBotToAddrInput) bool {
	if b.c != nil {
		b.c <- input
	}

	return true
}
