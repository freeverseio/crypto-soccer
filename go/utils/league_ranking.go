package utils

import "math/big"

func right(x *big.Int, n uint) *big.Int {
	return new(big.Int).Rsh(x, n)
}

func left(x *big.Int, n uint) *big.Int {
	return new(big.Int).Lsh(x, n)
}

func and(x *big.Int, n int64) *big.Int {
	return new(big.Int).And(x, big.NewInt(n))
}

func andBN(x *big.Int, y *big.Int) *big.Int {
	return new(big.Int).And(x, y)
}

func or(x *big.Int, n int64) *big.Int {
	return new(big.Int).Or(x, big.NewInt(n))
}

func orBN(x *big.Int, y *big.Int) *big.Int {
	return new(big.Int).Or(x, y)
}

func lessThan(x *big.Int, n int64) bool {
	return x.Cmp(big.NewInt(n)) == -1
}

func largerThan(x *big.Int, n int64) bool {
	return x.Cmp(big.NewInt(n)) == 1
}

func equals(x *big.Int, n int64) bool {
	return x.Cmp(big.NewInt(n)) == 0
}

func twoToPow(n uint64) int64 {
	return 2 << (n - 1)
}

func decodeTZCountryAndValGo(encoded *big.Int) (uint8, *big.Int, *big.Int) {
	return uint8(and(right(encoded, 38), 31).Int64()), and(right(encoded, 28), 1023), and(encoded, 268435455)
}
