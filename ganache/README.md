# deploy contract to ganache
```
./deploy.sh
./run_ganache.sh
go run main.go
```

# enroll all stakers and run service
```
cd ..
go run main.go --config ganache/config_ganache.yaml staker enroll 0x68dB32D26d9529B2a142927c6f1af248fc6Ba7e9
go run main.go --config ganache/config_ganache.yaml staker enroll 0x35bb6eF95c72bf4804334BB9d6A3c77Bef18d81B
go run main.go --config ganache/config_ganache.yaml staker enroll 0x26d8094A90AD1A4440600fa18073730Ef17F5eCE
go run main.go --config ganache/config_ganache.yaml staker enroll 0xb53cC19aD713e00cB71542d064215784c908D387
go run main.go --config ganache/config_ganache.yaml staker enroll 0xcb98ce2619f90f54052524Fb79b03E0261b01BEE
go run main.go --config ganache/config_ganache.yaml staker enroll 0x4C4Cb54BdBad6c96805b92b8063F3c75B24F65eB
go run main.go --config ganache/config_ganache.yaml staker enroll 0x67B27Df78bb0f199EDBd568e655BD9b2B202866d
go run main.go --config ganache/config_ganache.yaml staker enroll 0xE93E2b43fC45CCEcc056A6Ea400972298A304b4B
go run main.go --config ganache/config_ganache.yaml staker enroll 0x288Ab710f8DEc0b13753Fec71161E50Ee0cDA7e6
go run main.go service start
```
