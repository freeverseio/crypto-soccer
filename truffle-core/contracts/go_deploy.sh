#!/bin/sh

MY_DIR=`dirname "$0"`
SOLC=solcjs
ABIGEN=~/go/bin/abigen
BUILD_DIR=${MY_DIR}/build
ASSETS=${MY_DIR}/assets/Assets.sol
CORE="${MY_DIR}/core/Engine.sol ${MY_DIR}/core/LeagueChallengeable.sol ${MY_DIR}/core/LeagueUpdatable.sol ${MY_DIR}/core/LeagueUsersAlongData.sol ${MY_DIR}/core/Leagues.sol ${MY_DIR}/core/LeaguesBase.sol ${MY_DIR}/core/LeaguesComputer.sol ${MY_DIR}/core/LeaguesScheduler.sol ${MY_DIR}/core/LeaguesScore.sol"

CORE_DEPS="${MY_DIR}/game_controller/GameControllerInterface.sol ${MY_DIR}/state/LeagueState.sol ${MY_DIR}/game_controller/GameControllerInterface.sol ${MY_DIR}/state/LeagueState.sol"
STATES="${MY_DIR}/state/LeagueState.sol ${MY_DIR}/state/PlayerState.sol ${MY_DIR}/state/PlayerState3D.sol ${MY_DIR}/state/TeamState.sol"

# building assets
(cd ${MY_DIR} && ${SOLC} --abi ${ASSETS} ${STATES} -o ${BUILD_DIR})
(cd ${MY_DIR} && ${SOLC} --bin ${ASSETS} ${STATES} -o ${BUILD_DIR})

(cd ${MY_DIR} && ${SOLC} --abi ${CORE} ${CORE_DEPS} ${STATES} -o ${BUILD_DIR})
(cd ${MY_DIR} && ${SOLC} --bin ${CORE} ${CORE_DEPS} ${STATES} -o ${BUILD_DIR})

# building league

#mkdir -p ${MY_DIR}/assets
#${ABIGEN} --bin=${BUILD_DIR}/___data_test_sol_SoccerSim.bin --abi=${BUILD_DIR}/___data_test_sol_SoccerSim.abi --pkg=lionel --out=${MY_DIR}/lionel/lionel.go
#${ABIGEN} --bin=${BUILD_DIR}/___data_test_sol_SoccerSim.bin --abi=${BUILD_DIR}/___data_test_sol_SoccerSim.abi --pkg=lionel --out=${MY_DIR}/lionel/lionel.go

GO_DESTDIR=${BUILD_DIR}/go
mkdir -p ${BUILD_DIR}/go

${ABIGEN} --bin=${BUILD_DIR}/__assets_Assets_sol_Assets.bin --abi=${BUILD_DIR}/__assets_Assets_sol_Assets.abi --pkg=assets --out=${GO_DESTDIR}/Assets.go

${ABIGEN} --bin=${BUILD_DIR}/__core_Engine_sol_Engine.bin --abi=${BUILD_DIR}/__core_Engine_sol_Engine.abi --pkg=core --out=${GO_DESTDIR}/Engine.go
${ABIGEN} --bin=${BUILD_DIR}/__core_LeaguesScore_sol_LeaguesScore.bin --abi=${BUILD_DIR}/__core_LeaguesScore_sol_LeaguesScore.abi --pkg=core --out=${GO_DESTDIR}/LeaguesScore.go
${ABIGEN} --bin=${BUILD_DIR}/__core_Leagues_sol_Leagues.bin --abi=${BUILD_DIR}/__core_Leagues_sol_Leagues.abi --pkg=core --out=${GO_DESTDIR}/Leagues.go
${ABIGEN} --bin=${BUILD_DIR}/__core_LeagueChallengeable_sol_LeagueChallengeable.bin --abi=${BUILD_DIR}/__core_LeagueChallengeable_sol_LeagueChallengeable.abi --pkg=core --out=${GO_DESTDIR}/LeagueChallengeable.go
${ABIGEN} --bin=${BUILD_DIR}/__core_LeagueUpdatable_sol_LeagueUpdatable.bin --abi=${BUILD_DIR}/__core_LeagueUpdatable_sol_LeagueUpdatable.abi --pkg=core --out=${GO_DESTDIR}/LeagueUpdatable.go
${ABIGEN} --bin=${BUILD_DIR}/__core_LeagueUsersAlongData_sol_LeagueUsersAlongData.bin --abi=${BUILD_DIR}/__core_LeagueUsersAlongData_sol_LeagueUsersAlongData.abi --pkg=core --out=${GO_DESTDIR}/LeagueUsersAlongData.go
${ABIGEN} --bin=${BUILD_DIR}/__core_LeaguesBase_sol_LeaguesBase.bin --abi=${BUILD_DIR}/__core_LeaguesBase_sol_LeaguesBase.abi --pkg=core --out=${GO_DESTDIR}/LeaguesBase.go
${ABIGEN} --bin=${BUILD_DIR}/__core_LeaguesComputer_sol_LeaguesComputer.bin --abi=${BUILD_DIR}/__core_LeaguesComputer_sol_LeaguesComputer.abi --pkg=core --out=${GO_DESTDIR}/LeaguesComputer.go
${ABIGEN} --bin=${BUILD_DIR}/__core_LeaguesScheduler_sol_LeaguesScheduler.bin --abi=${BUILD_DIR}/__core_LeaguesScheduler_sol_LeaguesScheduler.abi --pkg=core --out=${GO_DESTDIR}/LeaguesScheduler.go

${ABIGEN} --bin=${BUILD_DIR}/__state_LeagueState_sol_LeagueState.bin --abi=${BUILD_DIR}/__state_LeagueState_sol_LeagueState.abi --pkg=state --out=${GO_DESTDIR}/LeagueState.go
${ABIGEN} --bin=${BUILD_DIR}/__state_PlayerState3D_sol_PlayerState3D.bin --abi=${BUILD_DIR}/__state_PlayerState3D_sol_PlayerState3D.abi --pkg=state --out=${GO_DESTDIR}/PlayerState3D.go
${ABIGEN} --bin=${BUILD_DIR}/__state_PlayerState_sol_PlayerState.bin --abi=${BUILD_DIR}/__state_PlayerState_sol_PlayerState.abi --pkg=state --out=${GO_DESTDIR}/PlayerState.go
${ABIGEN} --bin=${BUILD_DIR}/__state_TeamState_sol_TeamState.bin --abi=${BUILD_DIR}/__state_TeamState_sol_TeamState.abi --pkg=state --out=${GO_DESTDIR}/TeamState.go

${ABIGEN} --bin=${BUILD_DIR}/__game_controller_GameControllerInterface_sol_GameControllerInterface.bin --abi=${BUILD_DIR}/__game_controller_GameControllerInterface_sol_GameControllerInterface.abi --pkg=game_controller --out=${GO_DESTDIR}/GameControllerInterface.go
