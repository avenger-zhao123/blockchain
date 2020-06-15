module blockchain

go 1.14

require (
	github.com/FactomProject/basen v0.0.0-20150613233007-fe3947df716e // indirect
	github.com/FactomProject/btcutilecc v0.0.0-20130527213604-d3a63a5752ec // indirect
	github.com/cmars/basen v0.0.0-20150613233007-fe3947df716e // indirect
	github.com/mr-tron/base58 v1.2.0
	github.com/stretchr/testify v1.6.1 // indirect
	github.com/syndtr/goleveldb v1.0.0
	github.com/tyler-smith/go-bip32 v0.0.0-20170922074101-2c9cfd177564
	github.com/tyler-smith/go-bip39 v1.0.2
	golang.org/x/crypto v0.0.0-20200604202706-70a84ac30bf9
)

replace ./golang.org/x/sys => ./github.com/golang/sys

replace ./golang.org/x/crypto => ./github.com/golang/crypto
