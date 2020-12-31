`go/synchronizer/process/calendar.go`

    func Populate:
        - Instead of retrieving teams indexs from contracts now it's from golang code.
        - Instead of updating match by match now it's on bulk

    func Reset:
        - Instead of updating match by match now it's on bulk

`go/synchronizer/process/league_processor.go`

    func resetLeague:
        - Instead of updating team by team now it's on bulk

## Results comparison

Run tests when old version:

```
./test_run_shuffle_tz_with_data.sh
INFO[0000] [DBMS] ... connected
INFO[0000] Repo root at: /home/christian/Projects/crypto-soccer
INFO[0000] Deploy by truffle: ./node_modules/.bin/truffle migrate --network local --reset
INFO[0012] PROXY=0xE95e9F1b339d45fb4eaaCB4Cce211080Ee4070c2
=== RUN   TestSuffleTzWithData
=== PAUSE TestSuffleTzWithData
=== CONT  TestSuffleTzWithData
INFO[0012] [processor|timezone 9] start process matches ...
INFO[0012] [processor|consume] shuffling timezone 10 and create calendars
INFO[0054] Total time elapsed: 41.458314506s
INFO[0054] Time updating: 3.611002729
INFO[0054] Time reseting: 29.95177106000002
INFO[0054] Time deleting: 7.894314790999996
INFO[0054] [processor|timezone 9] Retriving user actions QmceyngPuA2PGz8K8CefP28YASDDJnsaWPvow2nXo4NRJz
INFO[0054] [processor|timezone 9] loading matches from storage
INFO[0054] [processor|timezone 9] reset trainings
INFO[0054] [processor|timezone 9] processing 1st half of 64 matches
INFO[0054] [precessor|1stHalfParallelProcess] 8 workers, took 0.374338524 secs
INFO[0054] [processor|timezone 9] save user action in history
INFO[0054] [processor|timezone 9] save matches to storage
INFO[0055] [processor|timezone 9] ... end
--- PASS: TestSuffleTzWithData (42.30s)
PASS
ok      github.com/freeverseio/crypto-soccer/go/synchronizer/process    55.151s
```

Run tests with new version:

```
./test_run_shuffle_tz_with_data.sh
INFO[0000] [DBMS] ... connected
INFO[0000] Repo root at: /home/christian/Projects/crypto-soccer
INFO[0000] Deploy by truffle: ./node_modules/.bin/truffle migrate --network local --reset
INFO[0011] PROXY=0xF0d1fdcbf52CeE0a6D77F332944C0c537FE7D27D
=== RUN   TestSuffleTzWithData
=== PAUSE TestSuffleTzWithData
=== CONT  TestSuffleTzWithData
INFO[0011] [processor|timezone 9] start process matches ...
INFO[0011] [processor|consume] shuffling timezone 10 and create calendars
INFO[0033] Total time elapsed: 22.26633305s
INFO[0033] Time updating: 3.405279495
INFO[0033] Time reseting: 11.676089318000002
INFO[0033] Time deleting: 7.1836863430000015
INFO[0033] [processor|timezone 9] Retriving user actions QmceyngPuA2PGz8K8CefP28YASDDJnsaWPvow2nXo4NRJz
INFO[0033] [processor|timezone 9] loading matches from storage
INFO[0033] [processor|timezone 9] reset trainings
INFO[0033] [processor|timezone 9] processing 1st half of 64 matches
INFO[0034] [precessor|1stHalfParallelProcess] 8 workers, took 0.379486067 secs
INFO[0034] [processor|timezone 9] save user action in history
INFO[0034] [processor|timezone 9] save matches to storage
INFO[0034] [processor|timezone 9] ... end
--- PASS: TestSuffleTzWithData (23.11s)
PASS
ok      github.com/freeverseio/crypto-soccer/go/synchronizer/process    34.492s
```

---

`go/synchronizer/process/calendar.go`

    func Populate:
        - Instead of retrieving teams indexs from contracts now it's from golang code.
        - Instead of updating match by match now it's on bulk
        - Instead of having to fetch team information one by one for every team, now it has an option to pass the teams slice to the function and if the teams are not passed onto it, it will perform a single query to get the teams by tz country and leagueIdx ordered so it can perform a search.

    func Reset:
        - Instead of updating match by match now it's on bulk

    Two new functions resetByTzCountry and populatebyTzCountry which minimeze select and update querys to db

`go/synchronizer/process/league_processor.go`

    func resetLeague:
        - Instead of updating team by team now it's on bulk

    new func resetLeagues

Run tests with new version:

```
./test_run_shuffle_tz_with_data.sh
INFO[0000] [DBMS] ... connected
INFO[0000] Repo root at: /home/christian/Projects/crypto-soccer
INFO[0000] Deploy by truffle: ./node_modules/.bin/truffle migrate --network local --reset
INFO[0011] PROXY=0x7ca40Ae4c328628615787FefC9b4BCAE458dAAEA
=== RUN   TestSuffleTzWithData
=== PAUSE TestSuffleTzWithData
=== CONT  TestSuffleTzWithData
INFO[0011] [processor|timezone 9] start process matches ...
INFO[0011] [processor|consume] shuffling timezone 10 and create calendars
INFO[0017] Total time elapsed: 5.932333099s
INFO[0017] Time updating: 3.361815798
INFO[0017] Time reseting: 2.4952293340000002
INFO[0017] Time deleting: 0.073675767
INFO[0017] [processor|timezone 9] Retriving user actions QmceyngPuA2PGz8K8CefP28YASDDJnsaWPvow2nXo4NRJz
INFO[0017] [processor|timezone 9] loading matches from storage
INFO[0017] [processor|timezone 9] reset trainings
INFO[0017] [processor|timezone 9] processing 1st half of 64 matches
INFO[0017] [precessor|1stHalfParallelProcess] 8 workers, took 0.390216493 secs
INFO[0017] [processor|timezone 9] save user action in history
INFO[0017] [processor|timezone 9] save matches to storage
INFO[0018] [processor|timezone 9] ... end
--- PASS: TestSuffleTzWithData (6.77s)
PASS
ok      github.com/freeverseio/crypto-soccer/go/synchronizer/process    18.269s
```

Running addteams:

```
./addteams
INFO[0000] [DBMS] ... connected
INFO[0000] Repo root at: /home/christian/Projects/crypto-soccer
INFO[0000] Deploy by truffle: ./node_modules/.bin/truffle migrate --network local --reset
INFO[0015] PROXY=0xA5e320Bf016dfA4548F32534df3bFD3598001E8E
INFO[0016] Aqui hay 8835 teams
INFO[0016] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 16
INFO[0043] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 17
INFO[0065] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 18
INFO[0086] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 19
INFO[0107] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 20
INFO[0128] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 21
INFO[0149] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 22
INFO[0170] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 23
INFO[0190] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 24
INFO[0212] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 25
INFO[0232] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 26
INFO[0253] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 27
INFO[0273] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 28
INFO[0294] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 29
INFO[0315] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 30
INFO[0336] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 31
INFO[0357] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 32
INFO[0377] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 33
INFO[0398] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 34
INFO[0419] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 35
INFO[0440] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 36
INFO[0461] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 37
INFO[0481] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 38
INFO[0502] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 39
INFO[0522] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 40
INFO[0543] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 41
INFO[0564] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 42
INFO[0584] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 43
INFO[0605] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 44
INFO[0625] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 45
INFO[0646] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 46
INFO[0667] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 47
INFO[0688] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 48
INFO[0708] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 49
INFO[0729] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 50
INFO[0750] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 51
INFO[0771] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 52
INFO[0791] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 53
INFO[0811] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 54
INFO[0832] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 55
INFO[0853] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 56
INFO[0873] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 57
INFO[0894] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 58
INFO[0915] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 59
INFO[0936] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 60
INFO[0957] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 61
INFO[0978] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 62
INFO[0998] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 63
INFO[1019] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 64
INFO[1041] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 65
INFO[1062] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 66
INFO[1082] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 67
INFO[1103] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 68
INFO[1124] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 69
INFO[1144] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 70
INFO[1165] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 71
INFO[1185] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 72
INFO[1206] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 73
INFO[1227] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 74
INFO[1248] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 75
INFO[1268] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 76
INFO[1289] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 77
INFO[1310] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 78
INFO[1331] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 79
INFO[1352] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 80
INFO[1373] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 81
INFO[1394] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 82
INFO[1415] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 83
INFO[1436] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 84
INFO[1457] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 85
INFO[1478] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 86
INFO[1499] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 87
INFO[1519] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 88
INFO[1540] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 89
INFO[1561] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 90
INFO[1581] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 91
INFO[1602] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 92
INFO[1623] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 93
INFO[1644] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 94
INFO[1665] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 95
INFO[1686] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 96
INFO[1706] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 97
INFO[1727] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 98
INFO[1747] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 99
INFO[1768] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 100
INFO[1789] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 101
INFO[1809] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 102
INFO[1830] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 103
INFO[1850] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 104
INFO[1871] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 105
INFO[1891] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 106
INFO[1912] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 107
INFO[1933] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 108
INFO[1955] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 109
INFO[1976] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 110
INFO[1997] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 111
INFO[2017] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 112
INFO[2037] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 113
INFO[2058] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 114
INFO[2078] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 115
INFO[2099] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 116
INFO[2119] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 117
INFO[2140] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 118
INFO[2161] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 119
INFO[2181] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 120
INFO[2202] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 121
INFO[2223] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 122
INFO[2244] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 123
INFO[2265] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 124
INFO[2285] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 125
INFO[2305] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 126
INFO[2326] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 127
INFO[2347] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 128
INFO[2367] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 129
INFO[2388] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 130
INFO[2409] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 131
INFO[2429] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 132
INFO[2450] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 133
INFO[2470] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 134
INFO[2491] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 135
INFO[2512] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 136
INFO[2533] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 137
INFO[2554] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 138
INFO[2575] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 139
INFO[2595] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 140
INFO[2616] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 141
INFO[2637] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 142
INFO[2657] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 143
INFO[2678] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 144
INFO[2699] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 145
INFO[2720] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 146
INFO[2741] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 147
INFO[2761] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 148
INFO[2782] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 149
INFO[2803] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 150
INFO[2823] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 151
INFO[2845] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 152
INFO[2866] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 153
INFO[2886] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 154
INFO[2907] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 155
INFO[2927] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 156
INFO[2948] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 157
INFO[2969] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 158
INFO[2990] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 159
INFO[3011] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 160
INFO[3033] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 161
INFO[3053] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 162
INFO[3074] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 163
INFO[3095] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 164
INFO[3115] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 165
INFO[3136] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 166
INFO[3157] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 167
INFO[3178] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 168
INFO[3199] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 169
INFO[3220] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 170
INFO[3240] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 171
INFO[3261] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 172
INFO[3282] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 173
INFO[3303] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 174
INFO[3324] [processor|consume] LeaguesDivisionCreation tz 10, country 0, division 175
INFO[3344] Aqui hay 29315 teams
INFO[3344] [processor|timezone 9] start process matches ...
INFO[3344] [processor|consume] shuffling timezone 10 and create calendars
INFO[3414] [processor|timezone 9] Retriving user actions QmceyngPuA2PGz8K8CefP28YASDDJnsaWPvow2nXo4NRJz
INFO[3414] [processor|timezone 9] loading matches from storage
INFO[3414] [processor|timezone 9] reset trainings
INFO[3414] [processor|timezone 9] processing 1st half of 64 matches
INFO[3414] [precessor|1stHalfParallelProcess] 8 workers, took 0.449397503 secs
INFO[3414] [processor|timezone 9] save user action in history
INFO[3414] [processor|timezone 9] save matches to storage
INFO[3415] [processor|timezone 9] ... end
```

Shuffling took 70 secs

Run addteams with 29315 teams:

```
./addteams 
INFO[0000] [DBMS] ... connected                         
INFO[0000] Repo root at: /home/christian/Projects/crypto-soccer 
INFO[0000] Deploy by truffle: ./node_modules/.bin/truffle migrate --network local --reset 
INFO[0010] PROXY=0xEedE65eD217C735B892c8D71531ab523958fEe17 
INFO[0010] Aqui hay 29315 teams                         
INFO[0010] [processor|timezone 9] start process matches ... 
INFO[0010] [processor|consume] shuffling timezone 10 and create calendars 
INFO[0379] [processor|timezone 9] Retriving user actions QmceyngPuA2PGz8K8CefP28YASDDJnsaWPvow2nXo4NRJz 
INFO[0379] [processor|timezone 9] loading matches from storage 
INFO[0379] [processor|timezone 9] reset trainings       
INFO[0379] [processor|timezone 9] processing 1st half of 64 matches 
INFO[0379] [precessor|1stHalfParallelProcess] 8 workers, took 0.367866524 secs 
INFO[0379] [processor|timezone 9] save user action in history 
INFO[0379] [processor|timezone 9] save matches to storage
```

Shuffling took 369 seconds