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
