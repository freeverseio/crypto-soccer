// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package assets

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// AssetsABI is the input ABI used to generate the binding from.
const AssetsABI = "[{\"inputs\":[],\"constant\":true,\"name\":\"IDX_MD\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"encodedSkills\"}],\"constant\":true,\"name\":\"getSkills\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"encodedSkills\"}],\"constant\":true,\"name\":\"getSumOfSkills\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"IDX_R\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"SK_END\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"encodedSkills\"}],\"constant\":true,\"name\":\"getAggressiveness\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint8\",\"name\":\"timeZone\"},{\"type\":\"uint256\",\"name\":\"countryIdxInTZ\"},{\"type\":\"uint256\",\"name\":\"val\"}],\"constant\":true,\"name\":\"encodeTZCountryAndVal\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"DAYS_PER_ROUND\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"encodedSkills\"}],\"constant\":true,\"name\":\"getEndurance\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"MAX_PLAYER_AGE_AT_BIRTH\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"encoded\"}],\"constant\":true,\"name\":\"decodeTZCountryAndVal\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"},{\"type\":\"uint256\",\"name\":\"\"},{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"encodedSkills\"}],\"constant\":true,\"name\":\"getLeftishness\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"IDX_D\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"state\"},{\"type\":\"uint256\",\"name\":\"value\"}],\"constant\":true,\"name\":\"setPrevPlayerTeamId\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"IDX_LC\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"FREEVERSE\",\"outputs\":[{\"type\":\"address\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"IDX_LR\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"SK_SHO\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"encodedSkills\"}],\"constant\":true,\"name\":\"getPotential\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"LEAGUES_PER_DIV\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"encodedSkills\"}],\"constant\":true,\"name\":\"getSpeed\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"playerState\"}],\"constant\":true,\"name\":\"getPrevPlayerTeamId\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"state\"},{\"type\":\"uint8\",\"name\":\"currentShirtNum\"}],\"constant\":true,\"name\":\"setCurrentShirtNum\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"encodedSkills\"}],\"constant\":true,\"name\":\"getDefence\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"encodedSkills\"}],\"constant\":true,\"name\":\"getPass\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"IDX_CR\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"encodedSkills\"}],\"constant\":true,\"name\":\"getShoot\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"encodedSkills\"}],\"constant\":true,\"name\":\"getAlignedLastHalf\",\"outputs\":[{\"type\":\"bool\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"encodedSkills\"}],\"constant\":true,\"name\":\"getPlayerIdFromSkills\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"IDX_GK\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"playerState\"}],\"constant\":true,\"name\":\"getPlayerIdFromState\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"encodedSkills\"}],\"constant\":true,\"name\":\"getInjuryWeeksLeft\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"PLAYERS_PER_TEAM_INIT\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"gameDeployMonth\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"encodedSkills\"}],\"constant\":true,\"name\":\"getMonthOfBirth\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"currentRound\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"PLAYERS_PER_TEAM_MAX\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"IDX_MF\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"N_SKILLS\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint16[5]\",\"name\":\"skills\"},{\"type\":\"uint256\",\"name\":\"monthOfBirth\"},{\"type\":\"uint256\",\"name\":\"playerId\"},{\"type\":\"uint8[4]\",\"name\":\"birthTraits\"},{\"type\":\"bool\",\"name\":\"alignedLastHalf\"},{\"type\":\"bool\",\"name\":\"redCardLastGame\"},{\"type\":\"uint8\",\"name\":\"gamesNonStopping\"},{\"type\":\"uint8\",\"name\":\"injuryWeeksLeft\"}],\"constant\":true,\"name\":\"encodePlayerSkills\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"encoded\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"IDX_M\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"playerId\"},{\"type\":\"uint256\",\"name\":\"currentTeamId\"},{\"type\":\"uint8\",\"name\":\"currentShirtNum\"},{\"type\":\"uint256\",\"name\":\"prevPlayerTeamId\"},{\"type\":\"uint256\",\"name\":\"lastSaleBlock\"}],\"constant\":true,\"name\":\"encodePlayerState\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"SK_PAS\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"TEAMS_PER_LEAGUE\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"state\"},{\"type\":\"uint256\",\"name\":\"lastSaleBlock\"}],\"constant\":true,\"name\":\"setLastSaleBlock\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"NULL_ADDR\",\"outputs\":[{\"type\":\"address\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"IDX_LCR\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"constant\":true,\"name\":\"_timeZones\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"nCountriesToAdd\"},{\"type\":\"uint8\",\"name\":\"newestOrgMapIdx\"},{\"type\":\"uint8\",\"name\":\"newestSkillsIdx\"},{\"type\":\"bytes32\",\"name\":\"scoresRoot\"},{\"type\":\"uint8\",\"name\":\"updateCycleIdx\"},{\"type\":\"uint256\",\"name\":\"lastActionsSubmissionTime\"},{\"type\":\"uint256\",\"name\":\"lastUpdateTime\"},{\"type\":\"bytes32\",\"name\":\"actionsRoot\"},{\"type\":\"uint256\",\"name\":\"lastMarketClosureBlockNum\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"FREE_PLAYER_ID\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"encodedSkills\"}],\"constant\":true,\"name\":\"getForwardness\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"playerState\"},{\"type\":\"uint256\",\"name\":\"teamId\"}],\"constant\":true,\"name\":\"setCurrentTeamId\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"playerState\"}],\"constant\":true,\"name\":\"getLastSaleBlock\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"MIN_PLAYER_AGE_AT_BIRTH\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"encodedSkills\"}],\"constant\":true,\"name\":\"getSkillsVec\",\"outputs\":[{\"type\":\"uint16[5]\",\"name\":\"skills\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"encodedSkills\"}],\"constant\":true,\"name\":\"getRedCardLastGame\",\"outputs\":[{\"type\":\"bool\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"playerState\"}],\"constant\":true,\"name\":\"getCurrentTeamId\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"IDX_F\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"tactics\"}],\"constant\":true,\"name\":\"decodeTactics\",\"outputs\":[{\"type\":\"uint8[11]\",\"name\":\"lineup\"},{\"type\":\"bool[10]\",\"name\":\"extraAttack\"},{\"type\":\"uint8\",\"name\":\"tacticsId\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"encodedSkills\"}],\"constant\":true,\"name\":\"getGamesNonStopping\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"SK_DEF\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint8[11]\",\"name\":\"lineup\"},{\"type\":\"bool[10]\",\"name\":\"extraAttack\"},{\"type\":\"uint8\",\"name\":\"tacticsId\"}],\"constant\":true,\"name\":\"encodeTactics\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"playerState\"}],\"constant\":true,\"name\":\"getCurrentShirtNum\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"IDX_L\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"TEAMS_PER_DIVISION\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"IDX_C\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"SK_SPE\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"indexed\":false,\"type\":\"uint256\",\"name\":\"teamId\"},{\"indexed\":false,\"type\":\"address\",\"name\":\"to\"}],\"type\":\"event\",\"name\":\"TeamTransfer\",\"anonymous\":false},{\"inputs\":[{\"indexed\":false,\"type\":\"uint8\",\"name\":\"timezone\"},{\"indexed\":false,\"type\":\"uint256\",\"name\":\"countryIdxInTZ\"},{\"indexed\":false,\"type\":\"uint256\",\"name\":\"divisionIdxInCountry\"}],\"type\":\"event\",\"name\":\"DivisionCreation\",\"anonymous\":false},{\"inputs\":[{\"indexed\":false,\"type\":\"uint256\",\"name\":\"playerId\"},{\"indexed\":false,\"type\":\"uint256\",\"name\":\"state\"}],\"type\":\"event\",\"name\":\"PlayerStateChange\",\"anonymous\":false},{\"inputs\":[],\"constant\":false,\"name\":\"init\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint8\",\"name\":\"tz\"}],\"constant\":false,\"name\":\"initSingleTZ\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint8\",\"name\":\"timeZone\"}],\"constant\":true,\"name\":\"getLastUpdateTime\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint8\",\"name\":\"timeZone\"}],\"constant\":true,\"name\":\"getLastActionsSubmissionTime\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint8\",\"name\":\"tz\"},{\"type\":\"bytes32\",\"name\":\"root\"}],\"constant\":false,\"name\":\"setSkillsRoot\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"nonpayable\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint8\",\"name\":\"timeZone\"},{\"type\":\"bytes32\",\"name\":\"root\"}],\"constant\":false,\"name\":\"setActionsRoot\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"nonpayable\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint8\",\"name\":\"timeZone\"}],\"constant\":true,\"name\":\"getNCountriesInTZ\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint8\",\"name\":\"timeZone\"},{\"type\":\"uint256\",\"name\":\"countryIdxInTZ\"}],\"constant\":true,\"name\":\"getNDivisionsInCountry\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint8\",\"name\":\"timeZone\"},{\"type\":\"uint256\",\"name\":\"countryIdxInTZ\"}],\"constant\":true,\"name\":\"getNLeaguesInCountry\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint8\",\"name\":\"timeZone\"},{\"type\":\"uint256\",\"name\":\"countryIdxInTZ\"}],\"constant\":true,\"name\":\"getNTeamsInCountry\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint8\",\"name\":\"timeZone\"},{\"type\":\"uint256\",\"name\":\"countryIdxInTZ\"},{\"type\":\"uint256\",\"name\":\"teamIdxInCountry\"}],\"constant\":true,\"name\":\"_teamExistsInCountry\",\"outputs\":[{\"type\":\"bool\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"teamId\"}],\"constant\":true,\"name\":\"teamExists\",\"outputs\":[{\"type\":\"bool\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint8\",\"name\":\"timeZone\"},{\"type\":\"uint256\",\"name\":\"countryIdxInTZ\"},{\"type\":\"uint256\",\"name\":\"teamIdxInCountry\"}],\"constant\":true,\"name\":\"isBotTeamInCountry\",\"outputs\":[{\"type\":\"bool\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"teamId\"}],\"constant\":true,\"name\":\"isBotTeam\",\"outputs\":[{\"type\":\"bool\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint8\",\"name\":\"timeZone\"},{\"type\":\"uint256\",\"name\":\"countryIdxInTZ\"},{\"type\":\"uint256\",\"name\":\"teamIdxInCountry\"}],\"constant\":true,\"name\":\"getOwnerTeamInCountry\",\"outputs\":[{\"type\":\"address\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"teamId\"}],\"constant\":true,\"name\":\"getOwnerTeam\",\"outputs\":[{\"type\":\"address\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"playerId\"}],\"constant\":true,\"name\":\"getOwnerPlayer\",\"outputs\":[{\"type\":\"address\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"playerId\"}],\"constant\":true,\"name\":\"getCurrentTeamIdFromPlayerId\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"playerId\"}],\"constant\":true,\"name\":\"playerExists\",\"outputs\":[{\"type\":\"bool\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"playerId\"}],\"constant\":true,\"name\":\"isVirtualPlayer\",\"outputs\":[{\"type\":\"bool\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint8\",\"name\":\"timeZone\"},{\"type\":\"uint256\",\"name\":\"countryIdxInTZ\"},{\"type\":\"address\",\"name\":\"addr\"}],\"constant\":false,\"name\":\"transferFirstBotToAddr\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"teamId\"},{\"type\":\"address\",\"name\":\"addr\"}],\"constant\":false,\"name\":\"transferTeam\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint8\",\"name\":\"timeZone\"},{\"type\":\"uint256\",\"name\":\"countryIdxInTZ\"},{\"type\":\"uint256\",\"name\":\"teamIdxInCountry\"},{\"type\":\"uint8\",\"name\":\"shirtNum\"}],\"constant\":true,\"name\":\"getDefaultPlayerIdForTeamInCountry\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"teamId\"}],\"constant\":true,\"name\":\"getPlayerIdsInTeam\",\"outputs\":[{\"type\":\"uint256[25]\",\"name\":\"playerIds\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"playerId\"}],\"constant\":true,\"name\":\"getPlayerSkillsAtBirth\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"playerId\"}],\"constant\":true,\"name\":\"getPlayerStateAtBirth\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"playerId\"}],\"constant\":true,\"name\":\"getPlayerState\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"dna\"},{\"type\":\"uint256\",\"name\":\"playerCreationMonth\"}],\"constant\":true,\"name\":\"computeBirthMonth\",\"outputs\":[{\"type\":\"uint16\",\"name\":\"\"},{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"dna\"},{\"type\":\"uint8\",\"name\":\"shirtNum\"}],\"constant\":true,\"name\":\"computeSkills\",\"outputs\":[{\"type\":\"uint16[5]\",\"name\":\"\"},{\"type\":\"uint8[4]\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"teamId\"},{\"type\":\"uint8\",\"name\":\"shirtNum\"}],\"constant\":true,\"name\":\"isFreeShirt\",\"outputs\":[{\"type\":\"bool\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"playerId\"}],\"constant\":true,\"name\":\"getPlayerAgeInMonths\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"teamId\"}],\"constant\":true,\"name\":\"getFreeShirt\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"playerId\"},{\"type\":\"uint256\",\"name\":\"teamIdTarget\"}],\"constant\":false,\"name\":\"transferPlayer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint8\",\"name\":\"timeZone\"}],\"constant\":true,\"name\":\"countCountries\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint8\",\"name\":\"timeZone\"},{\"type\":\"uint256\",\"name\":\"countryIdxInTZ\"}],\"constant\":true,\"name\":\"countTeams\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"}]"

// AssetsBin is the compiled bytecode used for deploying new contracts.
const AssetsBin = `0x6080604052600161014960006101000a81548160ff02191690831515021790555034801561002c57600080fd5b5061559b8061003c6000396000f3fe608060405234801561001057600080fd5b50600436106105b45760003560e01c80638598243111610300578063c2bc41cd116101a8578063e81e21bb116100f4578063ec71bc82116100ad578063f305a21c11610087578063f305a21c14611ff4578063f8ef7b9e14612018578063f9d0723d1461203c578063fa80039b1461208b576105b4565b8063ec71bc8214611f6a578063ec7ecec514611f8e578063f21f5a8314611fd0576105b4565b8063e81e21bb14611d33578063e945e96a14611d57578063e9e7165214611da5578063eabf6a4b14611e6f578063eb78b7b714611ed9578063ec1c542314611f1b576105b4565b8063cc7d473b11610161578063dba6319e1161013b578063dba6319e14611bf2578063e1c7392a14611c41578063e6400ac414611c4b578063e804e51914611cf1576105b4565b8063cc7d473b14611b46578063cd2105e814611b8c578063d7b63a1114611bce576105b4565b8063c2bc41cd146119a6578063c37b1c25146119e8578063c566b5bc14611a34578063c73f808d14611a76578063c79055d414611ab8578063cc1cc3d714611adc576105b4565b80639f27112a11610267578063b32aa2c111610220578063b96b1a30116101fa578063b96b1a3014611861578063bc1a97c1146118f3578063c04f6d5314611939578063c258012b14611988576105b4565b8063b32aa2c1146117ad578063b3f390b3146117f3578063b96270971461183d576105b4565b80639f27112a14611636578063a3ceb703146116a3578063ab1b7c5e146116d4578063ac5db9ee146116f8578063ad63bcbd1461171c578063af76cd0114611761576105b4565b80638f9da214116102b95780638f9da214146113e6578063963fcc8014611454578063976daaac146114a757806398981756146114cb5780639c53e3fd146115115780639cc6234014611612576105b4565b806385982431146112da57806387f1e880146112f85780638a19c8bc1461133a5780638adddc9d146113585780638cc9a8d51461137c5780638f3db436146113c2576105b4565b80633d085f9611610463578063595ef25b116103ca5780637420a606116103835780637b2566a51161035d5780637b2566a5146111d2578063800257d51461122157806380bac7091461125957806383c31d3b146112b6576105b4565b80637420a6061461112a57806378f4c7181461114e57806379e7659714611190576105b4565b8063595ef25b14610f685780635adb40f514610fed5780635becd9991461103c57806365b4b47614611060578063673fe242146110a25780636f6c2ae0146110e8576105b4565b80634bea2a691161041c5780634bea2a6914610d655780634db989fd14610da7578063507b172314610df657806351585b4914610e3e578063547d829814610e8057806355a6f86f14610f26576105b4565b80633d085f9614610c0757806340cd05fd14610c2b578063434807ad14610c4f57806348d1e9c014610c91578063492afc6914610cb55780634b93f75314610d23576105b4565b8063228408b0116105225780633518dd1d116104db57806337fd56af116104b557806337fd56af14610afc57806338c96b5c14610b2057806339644f2114610b625780633c2eb36014610bac576105b4565b80633518dd1d14610a4a578063369151db14610a8c57806337a8630214610ab0576105b4565b8063228408b0146108a1578063258e5d901461090757806326657608146109495780632a238b0a1461098b5780632d0e08fd146109af5780633260840b146109f4576105b4565b80631884332c116105745780631884332c146107255780631a6daba2146107495780631fc7768f146107a65780631ffeb349146107e857806320748ae81461082a57806321ff8ae814610883576105b4565b80623e3223146105b957806292bf78146105dd578062aae8df1461061f5780630abcd3e51461067a5780631060c9c2146106bf578063169d291414610701575b600080fd5b6105c16120d0565b604051808260ff1660ff16815260200191505060405180910390f35b610609600480360360208110156105f357600080fd5b81019080803590602001909291905050506120d5565b6040518082815260200191505060405180910390f35b6106556004803603604081101561063557600080fd5b8101908080359060200190929190803590602001909291905050506120e3565b604051808361ffff1661ffff1681526020018281526020019250505060405180910390f35b6106a96004803603602081101561069057600080fd5b81019080803560ff16906020019092919050505061218d565b6040518082815260200191505060405180910390f35b6106eb600480360360208110156106d557600080fd5b81019080803590602001909291905050506121ba565b6040518082815260200191505060405180910390f35b6107096121f4565b604051808260ff1660ff16815260200191505060405180910390f35b61072d6121f9565b604051808260ff1660ff16815260200191505060405180910390f35b61078c6004803603606081101561075f57600080fd5b81019080803560ff16906020019092919080359060200190929190803590602001909291905050506121fe565b604051808215151515815260200191505060405180910390f35b6107d2600480360360208110156107bc57600080fd5b8101908080359060200190929190505050612215565b6040518082815260200191505060405180910390f35b610814600480360360208110156107fe57600080fd5b8101908080359060200190929190505050612226565b6040518082815260200191505060405180910390f35b61086d6004803603606081101561084057600080fd5b81019080803560ff1690602001909291908035906020019092919080359060200190929190505050612252565b6040518082815260200191505060405180910390f35b61088b6123df565b6040518082815260200191505060405180910390f35b6108f1600480360360808110156108b757600080fd5b81019080803560ff1690602001909291908035906020019092919080359060200190929190803560ff1690602001909291905050506123e4565b6040518082815260200191505060405180910390f35b6109336004803603602081101561091d57600080fd5b810190808035906020019092919050505061241e565b6040518082815260200191505060405180910390f35b6109756004803603602081101561095f57600080fd5b8101908080359060200190929190505050612430565b6040518082815260200191505060405180910390f35b610993612510565b604051808260ff1660ff16815260200191505060405180910390f35b6109de600480360360208110156109c557600080fd5b81019080803560ff169060200190929190505050612515565b6040518082815260200191505060405180910390f35b610a2060048036036020811015610a0a57600080fd5b810190808035906020019092919050505061253f565b604051808460ff1660ff168152602001838152602001828152602001935050505060405180910390f35b610a7660048036036020811015610a6057600080fd5b8101908080359060200190929190505050612569565b6040518082815260200191505060405180910390f35b610a9461257a565b604051808260ff1660ff16815260200191505060405180910390f35b610ae660048036036040811015610ac657600080fd5b81019080803590602001909291908035906020019092919050505061257f565b6040518082815260200191505060405180910390f35b610b04612629565b604051808260ff1660ff16815260200191505060405180910390f35b610b4c60048036036020811015610b3657600080fd5b810190808035906020019092919050505061262e565b6040518082815260200191505060405180910390f35b610b6a612648565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b610c0560048036036060811015610bc257600080fd5b81019080803560ff16906020019092919080359060200190929190803573ffffffffffffffffffffffffffffffffffffffff16906020019092919050505061264d565b005b610c0f61297e565b604051808260ff1660ff16815260200191505060405180910390f35b610c33612983565b604051808260ff1660ff16815260200191505060405180910390f35b610c7b60048036036020811015610c6557600080fd5b8101908080359060200190929190505050612988565b6040518082815260200191505060405180910390f35b610c99612999565b604051808260ff1660ff16815260200191505060405180910390f35b610ce160048036036020811015610ccb57600080fd5b810190808035906020019092919050505061299e565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b610d4f60048036036020811015610d3957600080fd5b81019080803590602001909291905050506129c8565b6040518082815260200191505060405180910390f35b610d9160048036036020811015610d7b57600080fd5b81019080803590602001909291905050506129da565b6040518082815260200191505060405180910390f35b610de060048036036040811015610dbd57600080fd5b8101908080359060200190929190803560ff1690602001909291905050506129f0565b6040518082815260200191505060405180910390f35b610e2260048036036020811015610e0c57600080fd5b8101908080359060200190929190505050612a9c565b604051808260ff1660ff16815260200191505060405180910390f35b610e6a60048036036020811015610e5457600080fd5b8101908080359060200190929190505050612ae3565b6040518082815260200191505060405180910390f35b610eb960048036036040811015610e9657600080fd5b8101908080359060200190929190803560ff169060200190929190505050612af5565b6040518083600560200280838360005b83811015610ee4578082015181840152602081019050610ec9565b5050505090500182600460200280838360005b83811015610f12578082015181840152602081019050610ef7565b505050509050019250505060405180910390f35b610f5260048036036020811015610f3c57600080fd5b8101908080359060200190929190505050612fa7565b6040518082815260200191505060405180910390f35b610fab60048036036060811015610f7e57600080fd5b81019080803560ff1690602001909291908035906020019092919080359060200190929190505050612fb9565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6110266004803603604081101561100357600080fd5b81019080803560ff1690602001909291908035906020019092919050505061303f565b6040518082815260200191505060405180910390f35b611044613090565b604051808260ff1660ff16815260200191505060405180910390f35b61108c6004803603602081101561107657600080fd5b8101908080359060200190929190505050613095565b6040518082815260200191505060405180910390f35b6110ce600480360360208110156110b857600080fd5b81019080803590602001909291905050506130a7565b604051808215151515815260200191505060405180910390f35b611114600480360360208110156110fe57600080fd5b81019080803590602001909291905050506130ba565b6040518082815260200191505060405180910390f35b6111326130d0565b604051808260ff1660ff16815260200191505060405180910390f35b61117a6004803603602081101561116457600080fd5b81019080803590602001909291905050506130d5565b6040518082815260200191505060405180910390f35b6111bc600480360360208110156111a657600080fd5b81019080803590602001909291905050506130eb565b6040518082815260200191505060405180910390f35b61120b600480360360408110156111e857600080fd5b81019080803560ff169060200190929190803590602001909291905050506130fc565b6040518082815260200191505060405180910390f35b6112576004803603604081101561123757600080fd5b810190808035906020019092919080359060200190929190505050613153565b005b61129c6004803603606081101561126f57600080fd5b81019080803560ff16906020019092919080359060200190929190803590602001909291905050506134c1565b604051808215151515815260200191505060405180910390f35b6112be613505565b604051808260ff1660ff16815260200191505060405180910390f35b6112e261350a565b6040518082815260200191505060405180910390f35b6113246004803603602081101561130e57600080fd5b8101908080359060200190929190505050613511565b6040518082815260200191505060405180910390f35b611342613523565b6040518082815260200191505060405180910390f35b61136061352a565b604051808260ff1660ff16815260200191505060405180910390f35b6113a86004803603602081101561139257600080fd5b810190808035906020019092919050505061352f565b604051808215151515815260200191505060405180910390f35b6113ca613559565b604051808260ff1660ff16815260200191505060405180910390f35b611412600480360360208110156113fc57600080fd5b810190808035906020019092919050505061355e565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b61148d6004803603604081101561146a57600080fd5b8101908080359060200190929190803560ff1690602001909291905050506135f3565b604051808215151515815260200191505060405180910390f35b6114af613701565b604051808260ff1660ff16815260200191505060405180910390f35b6114f7600480360360208110156114e157600080fd5b8101908080359060200190929190505050613706565b604051808215151515815260200191505060405180910390f35b6115fc60048036036101e081101561152857600080fd5b810190808060a001906005806020026040519081016040528092919082600560200280828437600081840152601f19601f8201169050808301925050505050509192919290803590602001909291908035906020019092919080608001906004806020026040519081016040528092919082600460200280828437600081840152601f19601f8201169050808301925050505050509192919290803515159060200190929190803515159060200190929190803560ff169060200190929190803560ff169060200190929190505050613730565b6040518082815260200191505060405180910390f35b61161a613d3f565b604051808260ff1660ff16815260200191505060405180910390f35b61168d600480360360a081101561164c57600080fd5b810190808035906020019092919080359060200190929190803560ff1690602001909291908035906020019092919080359060200190929190505050613d44565b6040518082815260200191505060405180910390f35b6116d2600480360360208110156116b957600080fd5b81019080803560ff169060200190929190505050613e12565b005b6116dc613ed4565b604051808260ff1660ff16815260200191505060405180910390f35b611700613ed9565b604051808260ff1660ff16815260200191505060405180910390f35b61174b6004803603602081101561173257600080fd5b81019080803560ff169060200190929190505050613ede565b6040518082815260200191505060405180910390f35b6117976004803603604081101561177757600080fd5b810190808035906020019092919080359060200190929190505050613f0b565b6040518082815260200191505060405180910390f35b6117d9600480360360208110156117c357600080fd5b8101908080359060200190929190505050613faf565b604051808215151515815260200191505060405180910390f35b6117fb614049565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b61184561404e565b604051808260ff1660ff16815260200191505060405180910390f35b61188d6004803603602081101561187757600080fd5b8101908080359060200190929190505050614053565b604051808a60ff1660ff1681526020018960ff1660ff1681526020018860ff1660ff1681526020018781526020018660ff1660ff168152602001858152602001848152602001838152602001828152602001995050505050505050505060405180910390f35b61191f6004803603602081101561190957600080fd5b81019080803590602001909291905050506140d8565b604051808215151515815260200191505060405180910390f35b6119726004803603604081101561194f57600080fd5b81019080803560ff16906020019092919080359060200190929190505050614138565b6040518082815260200191505060405180910390f35b611990614152565b6040518082815260200191505060405180910390f35b6119d2600480360360208110156119bc57600080fd5b8101908080359060200190929190505050614157565b6040518082815260200191505060405180910390f35b611a1e600480360360408110156119fe57600080fd5b810190808035906020019092919080359060200190929190505050614168565b6040518082815260200191505060405180910390f35b611a6060048036036020811015611a4a57600080fd5b8101908080359060200190929190505050614218565b6040518082815260200191505060405180910390f35b611aa260048036036020811015611a8c57600080fd5b810190808035906020019092919050505061422d565b6040518082815260200191505060405180910390f35b611ac06143c0565b604051808260ff1660ff16815260200191505060405180910390f35b611b0860048036036020811015611af257600080fd5b81019080803590602001909291905050506143c5565b6040518082600560200280838360005b83811015611b33578082015181840152602081019050611b18565b5050505090500191505060405180910390f35b611b7260048036036020811015611b5c57600080fd5b810190808035906020019092919050505061449f565b604051808215151515815260200191505060405180910390f35b611bb860048036036020811015611ba257600080fd5b81019080803590602001909291905050506144b2565b6040518082815260200191505060405180910390f35b611bd66144c8565b604051808260ff1660ff16815260200191505060405180910390f35b611c2b60048036036040811015611c0857600080fd5b81019080803560ff169060200190929190803590602001909291905050506144cd565b6040518082815260200191505060405180910390f35b611c49614516565b005b611c7760048036036020811015611c6157600080fd5b81019080803590602001909291905050506145f8565b6040518084600b60200280838360005b83811015611ca2578082015181840152602081019050611c87565b5050505090500183600a60200280838360005b83811015611cd0578082015181840152602081019050611cb5565b505050509050018260ff1660ff168152602001935050505060405180910390f35b611d1d60048036036020811015611d0757600080fd5b81019080803590602001909291905050506147d1565b6040518082815260200191505060405180910390f35b611d3b6147e2565b604051808260ff1660ff16815260200191505060405180910390f35b611da360048036036040811015611d6d57600080fd5b8101908080359060200190929190803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506147e7565b005b611e5960048036036102c0811015611dbc57600080fd5b81019080806101600190600b806020026040519081016040528092919082600b60200280828437600081840152601f19601f8201169050808301925050505050509192919290806101400190600a806020026040519081016040528092919082600a60200280828437600081840152601f19601f8201169050808301925050505050509192919290803560ff169060200190929190505050614879565b6040518082815260200191505060405180910390f35b611e9b60048036036020811015611e8557600080fd5b8101908080359060200190929190505050614a32565b6040518082601960200280838360005b83811015611ec6578082015181840152602081019050611eab565b5050505090500191505060405180910390f35b611f0560048036036020811015611eef57600080fd5b8101908080359060200190929190505050614bfb565b6040518082815260200191505060405180910390f35b611f5460048036036040811015611f3157600080fd5b81019080803560ff16906020019092919080359060200190929190505050614c0c565b6040518082815260200191505060405180910390f35b611f72614cd6565b604051808260ff1660ff16815260200191505060405180910390f35b611fba60048036036020811015611fa457600080fd5b8101908080359060200190929190505050614cdb565b6040518082815260200191505060405180910390f35b611fd8614d16565b604051808260ff1660ff16815260200191505060405180910390f35b611ffc614d1b565b604051808260ff1660ff16815260200191505060405180910390f35b612020614d20565b604051808260ff1660ff16815260200191505060405180910390f35b6120756004803603604081101561205257600080fd5b81019080803560ff16906020019092919080359060200190929190505050614d25565b6040518082815260200191505060405180910390f35b6120ba600480360360208110156120a157600080fd5b81019080803560ff169060200190929190505050614d3f565b6040518082815260200191505060405180910390f35b600481565b600060ba82901c9050919050565b6000806101e0831161215d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601b8152602001807f696e76616c696420706c617965724372656174696f6e4d6f6e7468000000000081525060200191505060405180910390fd5b60006014858161216957fe5b066010019050600585901c9450600c810261ffff1684038592509250509250929050565b600061219882614d69565b60018260ff16601981106121a857fe5b600d0201600001805490509050919050565b60006121c58261241e565b6121ce83612ae3565b6121d784612fa7565b6121e0856129c8565b6121e986613095565b010101019050919050565b600181565b600481565b600061220a8484614138565b821090509392505050565b60006007606c83901c169050919050565b600061224b61224461223f61223a8561422d565b613511565b614df4565b4203614e0e565b9050919050565b600060208460ff16106122cd576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260148152602001807f646566656e6365206f7574206f6620626f756e6400000000000000000000000081525060200191505060405180910390fd5b6104008310612344576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260148152602001807f646566656e6365206f7574206f6620626f756e6400000000000000000000000081525060200191505060405180910390fd5b631000000082106123bd576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260148152602001807f646566656e6365206f7574206f6620626f756e6400000000000000000000000081525060200191505060405180910390fd5b600060268560ff16901b9050601c84901b811790508281179150509392505050565b601081565b6000601260ff168260ff16106123fd5760019050612416565b61241385858460ff16601260ff16870201612252565b90505b949350505050565b6000613fff60ba83901c169050919050565b60008060008061243f8561253f565b9250925092506000601260ff16828161245457fe5b0490506124628484836121fe565b6124d4576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252600f8152602001807f696e76616c6964207465616d206964000000000000000000000000000000000081525060200191505060405180910390fd5b60006124e1858584612252565b90506000601260ff1684816124f257fe5b069050612503888383600080613d44565b9650505050505050919050565b602081565b600061252082614d69565b60018260ff166019811061253057fe5b600d0201600a01549050919050565b6000806000601f602685901c166103ff601c86901c16630fffffff86169250925092509193909250565b60006007607783901c169050919050565b600181565b60006508000000000082106125fc576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601a8152602001807f707265764c6561677565496478206f7574206f6620626f756e6400000000000081525060200191505060405180910390fd5b741ffffffffffc0000000000000000000000000000001983169250607a82901b8317925082905092915050565b600681565b600061264161263c83614cdb565b6144b2565b9050919050565b600181565b600060018460ff166019811061265f57fe5b600d0201600001838154811061267157fe5b90600052602060002090600502016004015490506126908484836134c1565b612702576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601e8152602001807f63616e6e6f74207472616e736665722061206e6f6e2d626f74207465616d000081525060200191505060405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614156127a5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252600f8152602001807f696e76616c69642061646472657373000000000000000000000000000000000081525060200191505060405180910390fd5b6127ad6152f9565b6000601260ff1690505b601960ff168110156127e75760018282601981106127d157fe5b60200201818152505080806001019150506127b7565b5060405180604001604052808281526020018473ffffffffffffffffffffffffffffffffffffffff1681525060018660ff166019811061282357fe5b600d0201600001858154811061283557fe5b9060005260206000209060050201600301600084815260200190815260200160002060008201518160000190601961286e92919061531c565b5060208201518160190160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555090505060018560ff16601981106128c957fe5b600d020160000184815481106128db57fe5b9060005260206000209060050201600401600081548092919060010191905055506000612909868685612252565b90507f77b66eb1e6d2bc131b79be4213ae7f08f29351c01060e10bcc0302278067bf388185604051808381526020018273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019250505060405180910390a1505050505050565b600581565b600081565b6000600f607d83901c169050919050565b601081565b6000806000806129ad8561253f565b9250925092506129be838383612fb9565b9350505050919050565b6000613fff60e483901c169050919050565b60006507ffffffffff607a83901c169050919050565b600060208260ff1610612a6b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601c8152602001807f63757272656e7453686972744e756d206f7574206f6620626f756e640000000081525060200191505060405180910390fd5b7503e00000000000000000000000000000000000000000198316925060a58260ff16901b8317925082905092915050565b600080600160190390505b60008160ff1610612ad857612abc83826135f3565b15612aca5780915050612ade565b808060019003915050612aa7565b50601990505b919050565b6000613fff60c883901c169050919050565b612afd61535c565b612b0561537e565b612b0d61535c565b612b1561535c565b6000600a8781612b2157fe5b069050600080600060046201552f8b81612b3757fe5b0681612b3f57fe5b06905060048a901c995060038960ff161015612b875760c885600060ff1660058110612b6757fe5b602002019061ffff16908161ffff16815250506000925060009150612ded565b60088960ff161015612bfb57602885600060ff1660058110612ba557fe5b602002019061ffff16908161ffff168152505060a085600360ff1660058110612bca57fe5b602002019061ffff16908161ffff16815250506001925060078960ff168b0181612bf057fe5b066001019150612dec565b600a8960ff161015612c4a5760a085600260ff1660058110612c1957fe5b602002019061ffff16908161ffff16815250506002925060078960ff168b0181612c3f57fe5b066001019150612deb565b600c8960ff161015612cbe57608285600260ff1660058110612c6857fe5b602002019061ffff16908161ffff1681525050604685600060ff1660058110612c8d57fe5b602002019061ffff16908161ffff16815250506004925060078960ff168b0181612cb357fe5b066001019150612dea565b600e8960ff161015612d3257608285600260ff1660058110612cdc57fe5b602002019061ffff16908161ffff1681525050604685600360ff1660058110612d0157fe5b602002019061ffff16908161ffff16815250506005925060078960ff168b0181612d2757fe5b066001019150612de9565b60108960ff161015612d955760a085600060ff1660058110612d5057fe5b602002019061ffff16908161ffff1681525050604685600360ff1660058110612d7557fe5b602002019061ffff16908161ffff16815250506003925060069150612de8565b60a085600060ff1660058110612da757fe5b602002019061ffff16908161ffff1681525050604685600360ff1660058110612dcc57fe5b602002019061ffff16908161ffff168152505060039250600391505b5b5b5b5b5b60338a901c9950600080600090505b600560ff168160ff161015612eda576000878260ff1660058110612e1c57fe5b602002015161ffff161415612e5d5760328c81612e3557fe5b06888260ff1660058110612e4557fe5b602002019061ffff16908161ffff1681525050612eae565b6064878260ff1660058110612e6e57fe5b602002015160328e81612e7d57fe5b060261ffff1681612e8a57fe5b04888260ff1660058110612e9a57fe5b602002019061ffff16908161ffff16815250505b60068c901c9b50878160ff1660058110612ec457fe5b6020020151820191508080600101915050612dfc565b5060fa8161ffff161015612f5b576000600560ff168260fa0361ffff1681612efe57fe5b04905060008090505b60058160ff161015612f585781898260ff1660058110612f2357fe5b602002015101898260ff1660058110612f3857fe5b602002019061ffff16908161ffff16815250508080600101915050612f07565b50505b8660405180608001604052808760ff1660ff1681526020018660ff1660ff1681526020018560ff1660ff1681526020018460ff1660ff1681525098509850505050505050509250929050565b6000613fff60d683901c169050919050565b6000612fc484614d69565b612fce8484614e28565b60018460ff1660198110612fde57fe5b600d02016000018381548110612ff057fe5b9060005260206000209060050201600301600083815260200190815260200160002060190160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690509392505050565b600061304a83614d69565b6130548383614e28565b60018360ff166019811061306457fe5b600d0201600001828154811061307657fe5b906000526020600020906005020160000154905092915050565b600381565b6000613fff60f283901c169050919050565b6000600180607684901c16149050919050565b60006507ffffffffff608183901c169050919050565b600081565b60006507ffffffffff60d583901c169050919050565b60006007606f83901c169050919050565b600061310783614d69565b6131118383614e28565b608060ff1660018460ff166019811061312657fe5b600d0201600001838154811061313857fe5b90600052602060002090600502016000015402905092915050565b61315c826140d8565b801561316d575061316c81613706565b5b6131df576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260198152602001807f756e6578697374656e7420706c61796572206f72207465616d0000000000000081525060200191505060405180910390fd5b60006131ea83614cdb565b9050600081905060006131fc836144b2565b905083811415613274576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260208152602001807f63616e6e6f74207472616e7366657220746f206f726967696e616c207465616d81525060200191505060405180910390fd5b61327d8161352f565b158015613290575061328e8461352f565b155b6132e5576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260368152602001806154f06036913960400191505060405180910390fd5b60006132f084614bfb565b905060006132fd86612a9c565b9050601960ff168160ff16141561335f576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260288152602001806155486028913960400191505060405180910390fd5b6133698487614168565b935061337584826129f0565b93506133818443613f0b565b9350836000808981526020019081526020016000208190555060008060006133a88661253f565b9250925092506001808460ff16601981106133bf57fe5b600d020160000183815481106133d157fe5b90600052602060002090600502016003016000838152602001908152602001600020600001866019811061340157fe5b018190555061340f8961253f565b8093508194508295505050508960018460ff166019811061342c57fe5b600d0201600001838154811061343e57fe5b906000526020600020906005020160030160008381526020019081526020016000206000018560ff166019811061347157fe5b01819055507f65a4d4a8a0afb474d2e9465a6a1a41bb88fd04f41152ba070421f1b3771f15d88a88604051808381526020018281526020019250505060405180910390a150505050505050505050565b60008073ffffffffffffffffffffffffffffffffffffffff166134e5858585612fb9565b73ffffffffffffffffffffffffffffffffffffffff161490509392505050565b601281565b6101475481565b6000613fff60ac83901c169050919050565b6101485481565b601981565b60008060008061353e8561253f565b92509250925061354f8383836134c1565b9350505050919050565b600581565b6000613569826140d8565b6135db576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260118152602001807f756e6578697374656e7420706c6179657200000000000000000000000000000081525060200191505060405180910390fd5b6135ec6135e78361262e565b61299e565b9050919050565b6000806000806136028661253f565b9250925092506136138383836134c1565b15613669576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602a815260200180615453602a913960400191505060405180910390fd5b600060018460ff166019811061367b57fe5b600d0201600001838154811061368d57fe5b906000526020600020906005020160030160008381526020019081526020016000206000018660ff16601981106136c057fe5b01549050600160120360ff168660ff1611156136f05760008114806136e55750600181145b9450505050506136fb565b600181149450505050505b92915050565b600581565b6000806000806137158561253f565b9250925092506137268383836121fe565b9350505050919050565b600080600090505b600560ff168160ff1610156137e2576140008a8260ff166005811061375957fe5b602002015161ffff16106137d5576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260128152602001807f736b696c6c206f7574206f6620626f756e64000000000000000000000000000081525060200191505060405180910390fd5b8080600101915050613738565b50600a86600060ff16600481106137f557fe5b602002015160ff1610613870576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260168152602001807f706f74656e7469616c206f7574206f6620626f756e640000000000000000000081525060200191505060405180910390fd5b600686600160ff166004811061388257fe5b602002015160ff16106138fd576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260188152602001807f666f72776172646e657373206f7574206f6620626f756e64000000000000000081525060200191505060405180910390fd5b600886600260ff166004811061390f57fe5b602002015160ff161061398a576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260188152602001807f6c6566697473686e657373206f7574206f6620626f756e64000000000000000081525060200191505060405180910390fd5b600886600360ff166004811061399c57fe5b602002015160ff1610613a17576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601b8152602001807f616767726573736976656e657373206f7574206f6620626f756e64000000000081525060200191505060405180910390fd5b600086600260ff1660048110613a2957fe5b602002015160ff161415613aa857600086600160ff1660048110613a4957fe5b602002015160ff1614613aa7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602b81526020018061547d602b913960400191505060405180910390fd5b5b60088360ff1610613b21576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601d8152602001807f67616d65734e6f6e53746f7070696e67206f7574206f6620626f756e6400000081525060200191505060405180910390fd5b6140008810613b7b576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260238152602001806154306023913960400191505060405180910390fd5b600087118015613b9057506508000000000087105b613c02576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260158152602001807f706c617965724964206f7574206f6620626f756e64000000000000000000000081525060200191505060405180910390fd5b60008090505b600560ff168160ff161015613c5457600e600182010260ff166101000361ffff168a8260ff1660058110613c3857fe5b602002015161ffff16901b821791508080600101915050613c08565b5060ac88901b81179050608187901b81179050607d86600060ff1660048110613c7957fe5b602002015160ff16901b81179050607a86600160ff1660048110613c9957fe5b602002015160ff16901b81179050607786600260ff1660048110613cb957fe5b602002015160ff16901b81179050607685613cd5576000613cd8565b60015b60ff16901b81179050607584613cef576000613cf2565b60015b60ff16901b8117905060728360ff16901b81179050606f8260ff16901b81179050606c86600360ff1660048110613d2557fe5b602002015160ff16901b8117905098975050505050505050565b600281565b60008086118015613d5a57506508000000000086105b613dcc576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260158152602001807f706c617965724964206f7574206f6620626f756e64000000000000000000000081525060200191505060405180910390fd5b600060d587901b9050613ddf8187614168565b9050613deb81866129f0565b9050613df7818561257f565b9050613e038184613f0b565b90508091505095945050505050565b6001151561014960009054906101000a900460ff16151514613e9c576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260178152602001807f63616e6e6f7420696e697469616c697a6520747769636500000000000000000081525060200191505060405180910390fd5b613ea542614e0e565b61014781905550613eb581614e9e565b600061014960006101000a81548160ff02191690831515021790555050565b600281565b600881565b6000613ee982614d69565b60018260ff1660198110613ef957fe5b600d0201600001805490509050919050565b60006408000000008210613f87576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601a8152602001807f6c61737453616c65426c6f636b206f7574206f6620626f756e6400000000000081525060200191505060405180910390fd5b6f03ffffffff80000000000000000000001983169250605782901b8317925082905092915050565b6000613fba826140d8565b61402c576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260118152602001807f756e6578697374656e7420706c6179657200000000000000000000000000000081525060200191505060405180910390fd5b600080600084815260200190815260200160002054149050919050565b600081565b600781565b6001816019811061406057fe5b600d02016000915090508060010160009054906101000a900460ff16908060060160009054906101000a900460ff16908060060160019054906101000a900460ff16908060070154908060080160009054906101000a900460ff169080600901549080600a01549080600b01549080600c0154905089565b6000808214156140eb5760009050614133565b6000806000848152602001908152602001600020541461410e5760019050614133565b600080600061411c8561253f565b92509250925061412d838383614ff4565b93505050505b919050565b6000600860ff166141498484614d25565b02905092915050565b600181565b60006007607a83901c169050919050565b60006508000000000082106141e5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601b8152602001807f63757272656e745465616d496478206f7574206f6620626f756e64000000000081525060200191505060405180910390fd5b7a1ffffffffffc000000000000000000000000000000000000000000198316925060aa82901b8317925082905092915050565b60006407ffffffff605783901c169050919050565b60008060008061423c8561253f565b9250925092506000601260ff16828161425157fe5b0490506000601260ff16838161426357fe5b0690506000608060ff16838161427557fe5b0490506142838686856121fe565b6142f5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252600f8152602001807f696e76616c6964207465616d206964000000000000000000000000000000000081525060200191505060405180910390fd5b600086868585604051602001808560ff1660ff1681526020018481526020018381526020018260ff1660ff1681526020019450505050506040516020818303038152906040528051906020012060001c90506000601e601060018a60ff166019811061435d57fe5b600d0201600001898154811061436f57fe5b906000526020600020906005020160020160008681526020019081526020016000205402601e610147540201816143a257fe5b0490506143b18285838d615011565b98505050505050505050919050565b601081565b6143cd61535c565b6143d682613095565b816000600581106143e357fe5b602002019061ffff16908161ffff16815250506143ff826129c8565b8160016005811061440c57fe5b602002019061ffff16908161ffff168152505061442882612fa7565b8160026005811061443557fe5b602002019061ffff16908161ffff168152505061445182612ae3565b8160036005811061445e57fe5b602002019061ffff16908161ffff168152505061447a8261241e565b8160046005811061448757fe5b602002019061ffff16908161ffff1681525050919050565b6000600180607584901c16149050919050565b60006507ffffffffff60aa83901c169050919050565b600381565b60006144d883614d69565b8160018460ff16601981106144e957fe5b600d0201600b01819055504260018460ff166019811061450557fe5b600d02016009018190555092915050565b6001151561014960009054906101000a900460ff161515146145a0576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260178152602001807f63616e6e6f7420696e697469616c697a6520747769636500000000000000000081525060200191505060405180910390fd5b6145a942614e0e565b610147819055506000600190505b60198160ff1610156145d9576145cc81614e9e565b80806001019150506145b7565b50600061014960006101000a81548160ff021916908315150217905550565b6146006153a0565b6146086153c3565b60006a0200000000000000000000841061468a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601e8152602001807f7461637469637349642073686f756c642066697420696e20363120626974000081525060200191505060405180910390fd5b603f84169050600684901c935060008090505b600a8160ff1610156146ef576001808616146146ba5760006146bd565b60015b838260ff16600a81106146cc57fe5b602002019015159081151581525050600185901c9450808060010191505061469d565b5060008090505b600b8160ff1610156147c957601f8516848260ff16600b811061471557fe5b602002019060ff16908160ff1681525050601960ff16848260ff16600b811061473a57fe5b602002015160ff16106147b5576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260168152602001807f696e636f7272656374206c696e65757020656e7472790000000000000000000081525060200191505060405180910390fd5b600585901c945080806001019150506146f6565b509193909250565b60006007607283901c169050919050565b600381565b60008060006147f58561253f565b9250925092506148078383838761506a565b7f77b66eb1e6d2bc131b79be4213ae7f08f29351c01060e10bcc0302278067bf388585604051808381526020018273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019250505060405180910390a15050505050565b600060408260ff16106148f4576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601d8152602001807f7461637469637349642073686f756c642066697420696e20362062697400000081525060200191505060405180910390fd5b60008260ff16905060008090505b600a8160ff161015614950578060010260060160ff16858260ff16600a811061492757fe5b602002015161493757600061493a565b60015b60ff16901b821791508080600101915050614902565b5060008090505b600b8160ff161015614a2657601960ff16868260ff16600b811061497757fe5b602002015160ff16106149f2576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260168152602001807f696e636f7272656374206c696e65757020656e7472790000000000000000000081525060200191505060405180910390fd5b8060050260100160ff16868260ff16600b8110614a0b57fe5b602002015160ff16901b821791508080600101915050614957565b50809150509392505050565b614a3a6152f9565b6000806000614a488561253f565b925092509250614a598383836121fe565b614acb576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252600f8152602001807f696e76616c6964207465616d206964000000000000000000000000000000000081525060200191505060405180910390fd5b614ad68383836134c1565b15614b275760008090505b601960ff168160ff161015614b2157614afc848484846123e4565b858260ff1660198110614b0b57fe5b6020020181815250508080600101915050614ae1565b50614bf3565b60008090505b601960ff168160ff161015614bf157600060018560ff1660198110614b4e57fe5b600d02016000018481548110614b6057fe5b906000526020600020906005020160030160008481526020019081526020016000206000018260ff1660198110614b9357fe5b015490506000811415614bc957614bac858585856123e4565b868360ff1660198110614bbb57fe5b602002018181525050614be3565b80868360ff1660198110614bd957fe5b6020020181815250505b508080600101915050614b2d565b505b505050919050565b6000601f60a583901c169050919050565b60008160018460ff1660198110614c1f57fe5b600d020160040160018560ff1660198110614c3657fe5b600d020160060160019054906101000a900460ff1660ff1660028110614c5857fe5b018190555060018360ff1660198110614c6d57fe5b600d020160060160019054906101000a900460ff1660010360018460ff1660198110614c9557fe5b600d020160060160016101000a81548160ff021916908360ff1602179055504260018460ff1660198110614cc557fe5b600d0201600a018190555092915050565b600481565b6000614ce682613faf565b15614cfb57614cf482612430565b9050614d11565b6000808381526020019081526020016000205490505b919050565b608081565b600281565b600181565b6000601060ff16614d36848461303f565b02905092915050565b6000614d4a82614d69565b60018260ff1660198110614d5a57fe5b600d0201600901549050919050565b60008160ff16118015614d7f575060198160ff16105b614df1576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260178152602001807f74696d655a6f6e6520646f6573206e6f7420657869737400000000000000000081525060200191505060405180910390fd5b50565b6000600c6301e13380830281614e0657fe5b049050919050565b60006301e13380600c830281614e2057fe5b049050919050565b60018260ff1660198110614e3857fe5b600d0201600001805490508110614e9a576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260278152602001806154a86027913960400191505060405180910390fd5b5050565b614ea66153e6565b600181600001818152505060018260ff1660198110614ec157fe5b600d020160000181908060018154018082558091505090600182039060005260206000209060050201600090919290919091506000820151816000015560208201518160010160006101000a81548160ff021916908360ff160217905550604082015181600401555050506001808360ff1660198110614f3d57fe5b600d0201600001600081548110614f5057fe5b90600052602060002090600502016002016000808152602001908152602001600020819055506000801b60018360ff1660198110614f8a57fe5b600d0201600201600060028110614f9d57fe5b01819055507fc5d195855a200aa90e2052bcc795cedbc84c2a26556b1d5113b5a30c96003a0b82600080604051808460ff1660ff168152602001838152602001828152602001935050505060405180910390a15050565b6000601260ff166150058585614138565b02821090509392505050565b60008061501e86856120e3565b8161ffff169150809750819250505061503561535c565b61503d61537e565b6150478888612af5565b9150915061505d82848784600080600080613730565b9350505050949350505050565b61507384614d69565b61507d8484614e28565b6150888484846134c1565b156150fb576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601e8152602001807f63616e6e6f74207472616e736665722061206e6f6e2d626f74207465616d000081525060200191505060405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161415615181576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260218152602001806154cf6021913960400191505060405180910390fd5b8073ffffffffffffffffffffffffffffffffffffffff1660018560ff16601981106151a857fe5b600d020160000184815481106151ba57fe5b9060005260206000209060050201600301600084815260200190815260200160002060190160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16141561526d576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260228152602001806155266022913960400191505060405180910390fd5b8060018560ff166019811061527e57fe5b600d0201600001848154811061529057fe5b9060005260206000209060050201600301600084815260200190815260200160002060190160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050505050565b604051806103200160405280601990602082028038833980820191505090505090565b826019810192821561534b579160200282015b8281111561534a57825182559160200191906001019061532f565b5b509050615358919061540a565b5090565b6040518060a00160405280600590602082028038833980820191505090505090565b6040518060800160405280600490602082028038833980820191505090505090565b604051806101600160405280600b90602082028038833980820191505090505090565b604051806101400160405280600a90602082028038833980820191505090505090565b604051806060016040528060008152602001600060ff168152602001600081525090565b61542c91905b80821115615428576000816000905550600101615410565b5090565b9056fe6d6f6e74684f664269727468496e556e697854696d65206f7574206f6620626f756e6463616e6e6f742071756572792061626f757420746865207368697274206f66206120426f74205465616d6c6566746973686e65732063616e206f6e6c79206265207a65726f20666f7220676f616c6b656570657273636f756e74727920646f6573206e6f7420657869737420696e20746869732074696d655a6f6e6563616e6e6f74207472616e7366657220746f2061206e756c6c206164647265737363616e6e6f74207472616e7366657220706c61796572207768656e206174206c65617374206f6e65207465616d206973206120626f74627579657220616e642073656c6c657220617265207468652073616d652061646472746172676574207465616d20666f72207472616e7366657220697320616c72656164792066756c6ca165627a7a72305820a61f91439b65c553a1beb28e1cc73f558fad85c5807d46f58095b7217bbd55320029`

// DeployAssets deploys a new Ethereum contract, binding an instance of Assets to it.
func DeployAssets(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Assets, error) {
	parsed, err := abi.JSON(strings.NewReader(AssetsABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(AssetsBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Assets{AssetsCaller: AssetsCaller{contract: contract}, AssetsTransactor: AssetsTransactor{contract: contract}, AssetsFilterer: AssetsFilterer{contract: contract}}, nil
}

// Assets is an auto generated Go binding around an Ethereum contract.
type Assets struct {
	AssetsCaller     // Read-only binding to the contract
	AssetsTransactor // Write-only binding to the contract
	AssetsFilterer   // Log filterer for contract events
}

// AssetsCaller is an auto generated read-only Go binding around an Ethereum contract.
type AssetsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AssetsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AssetsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AssetsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AssetsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AssetsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AssetsSession struct {
	Contract     *Assets           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AssetsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AssetsCallerSession struct {
	Contract *AssetsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// AssetsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AssetsTransactorSession struct {
	Contract     *AssetsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AssetsRaw is an auto generated low-level Go binding around an Ethereum contract.
type AssetsRaw struct {
	Contract *Assets // Generic contract binding to access the raw methods on
}

// AssetsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AssetsCallerRaw struct {
	Contract *AssetsCaller // Generic read-only contract binding to access the raw methods on
}

// AssetsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AssetsTransactorRaw struct {
	Contract *AssetsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAssets creates a new instance of Assets, bound to a specific deployed contract.
func NewAssets(address common.Address, backend bind.ContractBackend) (*Assets, error) {
	contract, err := bindAssets(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Assets{AssetsCaller: AssetsCaller{contract: contract}, AssetsTransactor: AssetsTransactor{contract: contract}, AssetsFilterer: AssetsFilterer{contract: contract}}, nil
}

// NewAssetsCaller creates a new read-only instance of Assets, bound to a specific deployed contract.
func NewAssetsCaller(address common.Address, caller bind.ContractCaller) (*AssetsCaller, error) {
	contract, err := bindAssets(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AssetsCaller{contract: contract}, nil
}

// NewAssetsTransactor creates a new write-only instance of Assets, bound to a specific deployed contract.
func NewAssetsTransactor(address common.Address, transactor bind.ContractTransactor) (*AssetsTransactor, error) {
	contract, err := bindAssets(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AssetsTransactor{contract: contract}, nil
}

// NewAssetsFilterer creates a new log filterer instance of Assets, bound to a specific deployed contract.
func NewAssetsFilterer(address common.Address, filterer bind.ContractFilterer) (*AssetsFilterer, error) {
	contract, err := bindAssets(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AssetsFilterer{contract: contract}, nil
}

// bindAssets binds a generic wrapper to an already deployed contract.
func bindAssets(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AssetsABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Assets *AssetsRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Assets.Contract.AssetsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Assets *AssetsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Assets.Contract.AssetsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Assets *AssetsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Assets.Contract.AssetsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Assets *AssetsCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Assets.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Assets *AssetsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Assets.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Assets *AssetsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Assets.Contract.contract.Transact(opts, method, params...)
}

// DAYSPERROUND is a free data retrieval call binding the contract method 0x21ff8ae8.
//
// Solidity: function DAYS_PER_ROUND() constant returns(uint256)
func (_Assets *AssetsCaller) DAYSPERROUND(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "DAYS_PER_ROUND")
	return *ret0, err
}

// DAYSPERROUND is a free data retrieval call binding the contract method 0x21ff8ae8.
//
// Solidity: function DAYS_PER_ROUND() constant returns(uint256)
func (_Assets *AssetsSession) DAYSPERROUND() (*big.Int, error) {
	return _Assets.Contract.DAYSPERROUND(&_Assets.CallOpts)
}

// DAYSPERROUND is a free data retrieval call binding the contract method 0x21ff8ae8.
//
// Solidity: function DAYS_PER_ROUND() constant returns(uint256)
func (_Assets *AssetsCallerSession) DAYSPERROUND() (*big.Int, error) {
	return _Assets.Contract.DAYSPERROUND(&_Assets.CallOpts)
}

// FREEVERSE is a free data retrieval call binding the contract method 0x39644f21.
//
// Solidity: function FREEVERSE() constant returns(address)
func (_Assets *AssetsCaller) FREEVERSE(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "FREEVERSE")
	return *ret0, err
}

// FREEVERSE is a free data retrieval call binding the contract method 0x39644f21.
//
// Solidity: function FREEVERSE() constant returns(address)
func (_Assets *AssetsSession) FREEVERSE() (common.Address, error) {
	return _Assets.Contract.FREEVERSE(&_Assets.CallOpts)
}

// FREEVERSE is a free data retrieval call binding the contract method 0x39644f21.
//
// Solidity: function FREEVERSE() constant returns(address)
func (_Assets *AssetsCallerSession) FREEVERSE() (common.Address, error) {
	return _Assets.Contract.FREEVERSE(&_Assets.CallOpts)
}

// FREEPLAYERID is a free data retrieval call binding the contract method 0xc258012b.
//
// Solidity: function FREE_PLAYER_ID() constant returns(uint256)
func (_Assets *AssetsCaller) FREEPLAYERID(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "FREE_PLAYER_ID")
	return *ret0, err
}

// FREEPLAYERID is a free data retrieval call binding the contract method 0xc258012b.
//
// Solidity: function FREE_PLAYER_ID() constant returns(uint256)
func (_Assets *AssetsSession) FREEPLAYERID() (*big.Int, error) {
	return _Assets.Contract.FREEPLAYERID(&_Assets.CallOpts)
}

// FREEPLAYERID is a free data retrieval call binding the contract method 0xc258012b.
//
// Solidity: function FREE_PLAYER_ID() constant returns(uint256)
func (_Assets *AssetsCallerSession) FREEPLAYERID() (*big.Int, error) {
	return _Assets.Contract.FREEPLAYERID(&_Assets.CallOpts)
}

// IDXC is a free data retrieval call binding the contract method 0xf305a21c.
//
// Solidity: function IDX_C() constant returns(uint8)
func (_Assets *AssetsCaller) IDXC(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "IDX_C")
	return *ret0, err
}

// IDXC is a free data retrieval call binding the contract method 0xf305a21c.
//
// Solidity: function IDX_C() constant returns(uint8)
func (_Assets *AssetsSession) IDXC() (uint8, error) {
	return _Assets.Contract.IDXC(&_Assets.CallOpts)
}

// IDXC is a free data retrieval call binding the contract method 0xf305a21c.
//
// Solidity: function IDX_C() constant returns(uint8)
func (_Assets *AssetsCallerSession) IDXC() (uint8, error) {
	return _Assets.Contract.IDXC(&_Assets.CallOpts)
}

// IDXCR is a free data retrieval call binding the contract method 0x5becd999.
//
// Solidity: function IDX_CR() constant returns(uint8)
func (_Assets *AssetsCaller) IDXCR(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "IDX_CR")
	return *ret0, err
}

// IDXCR is a free data retrieval call binding the contract method 0x5becd999.
//
// Solidity: function IDX_CR() constant returns(uint8)
func (_Assets *AssetsSession) IDXCR() (uint8, error) {
	return _Assets.Contract.IDXCR(&_Assets.CallOpts)
}

// IDXCR is a free data retrieval call binding the contract method 0x5becd999.
//
// Solidity: function IDX_CR() constant returns(uint8)
func (_Assets *AssetsCallerSession) IDXCR() (uint8, error) {
	return _Assets.Contract.IDXCR(&_Assets.CallOpts)
}

// IDXD is a free data retrieval call binding the contract method 0x369151db.
//
// Solidity: function IDX_D() constant returns(uint8)
func (_Assets *AssetsCaller) IDXD(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "IDX_D")
	return *ret0, err
}

// IDXD is a free data retrieval call binding the contract method 0x369151db.
//
// Solidity: function IDX_D() constant returns(uint8)
func (_Assets *AssetsSession) IDXD() (uint8, error) {
	return _Assets.Contract.IDXD(&_Assets.CallOpts)
}

// IDXD is a free data retrieval call binding the contract method 0x369151db.
//
// Solidity: function IDX_D() constant returns(uint8)
func (_Assets *AssetsCallerSession) IDXD() (uint8, error) {
	return _Assets.Contract.IDXD(&_Assets.CallOpts)
}

// IDXF is a free data retrieval call binding the contract method 0xd7b63a11.
//
// Solidity: function IDX_F() constant returns(uint8)
func (_Assets *AssetsCaller) IDXF(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "IDX_F")
	return *ret0, err
}

// IDXF is a free data retrieval call binding the contract method 0xd7b63a11.
//
// Solidity: function IDX_F() constant returns(uint8)
func (_Assets *AssetsSession) IDXF() (uint8, error) {
	return _Assets.Contract.IDXF(&_Assets.CallOpts)
}

// IDXF is a free data retrieval call binding the contract method 0xd7b63a11.
//
// Solidity: function IDX_F() constant returns(uint8)
func (_Assets *AssetsCallerSession) IDXF() (uint8, error) {
	return _Assets.Contract.IDXF(&_Assets.CallOpts)
}

// IDXGK is a free data retrieval call binding the contract method 0x7420a606.
//
// Solidity: function IDX_GK() constant returns(uint8)
func (_Assets *AssetsCaller) IDXGK(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "IDX_GK")
	return *ret0, err
}

// IDXGK is a free data retrieval call binding the contract method 0x7420a606.
//
// Solidity: function IDX_GK() constant returns(uint8)
func (_Assets *AssetsSession) IDXGK() (uint8, error) {
	return _Assets.Contract.IDXGK(&_Assets.CallOpts)
}

// IDXGK is a free data retrieval call binding the contract method 0x7420a606.
//
// Solidity: function IDX_GK() constant returns(uint8)
func (_Assets *AssetsCallerSession) IDXGK() (uint8, error) {
	return _Assets.Contract.IDXGK(&_Assets.CallOpts)
}

// IDXL is a free data retrieval call binding the contract method 0xec71bc82.
//
// Solidity: function IDX_L() constant returns(uint8)
func (_Assets *AssetsCaller) IDXL(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "IDX_L")
	return *ret0, err
}

// IDXL is a free data retrieval call binding the contract method 0xec71bc82.
//
// Solidity: function IDX_L() constant returns(uint8)
func (_Assets *AssetsSession) IDXL() (uint8, error) {
	return _Assets.Contract.IDXL(&_Assets.CallOpts)
}

// IDXL is a free data retrieval call binding the contract method 0xec71bc82.
//
// Solidity: function IDX_L() constant returns(uint8)
func (_Assets *AssetsCallerSession) IDXL() (uint8, error) {
	return _Assets.Contract.IDXL(&_Assets.CallOpts)
}

// IDXLC is a free data retrieval call binding the contract method 0x37fd56af.
//
// Solidity: function IDX_LC() constant returns(uint8)
func (_Assets *AssetsCaller) IDXLC(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "IDX_LC")
	return *ret0, err
}

// IDXLC is a free data retrieval call binding the contract method 0x37fd56af.
//
// Solidity: function IDX_LC() constant returns(uint8)
func (_Assets *AssetsSession) IDXLC() (uint8, error) {
	return _Assets.Contract.IDXLC(&_Assets.CallOpts)
}

// IDXLC is a free data retrieval call binding the contract method 0x37fd56af.
//
// Solidity: function IDX_LC() constant returns(uint8)
func (_Assets *AssetsCallerSession) IDXLC() (uint8, error) {
	return _Assets.Contract.IDXLC(&_Assets.CallOpts)
}

// IDXLCR is a free data retrieval call binding the contract method 0xb9627097.
//
// Solidity: function IDX_LCR() constant returns(uint8)
func (_Assets *AssetsCaller) IDXLCR(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "IDX_LCR")
	return *ret0, err
}

// IDXLCR is a free data retrieval call binding the contract method 0xb9627097.
//
// Solidity: function IDX_LCR() constant returns(uint8)
func (_Assets *AssetsSession) IDXLCR() (uint8, error) {
	return _Assets.Contract.IDXLCR(&_Assets.CallOpts)
}

// IDXLCR is a free data retrieval call binding the contract method 0xb9627097.
//
// Solidity: function IDX_LCR() constant returns(uint8)
func (_Assets *AssetsCallerSession) IDXLCR() (uint8, error) {
	return _Assets.Contract.IDXLCR(&_Assets.CallOpts)
}

// IDXLR is a free data retrieval call binding the contract method 0x3d085f96.
//
// Solidity: function IDX_LR() constant returns(uint8)
func (_Assets *AssetsCaller) IDXLR(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "IDX_LR")
	return *ret0, err
}

// IDXLR is a free data retrieval call binding the contract method 0x3d085f96.
//
// Solidity: function IDX_LR() constant returns(uint8)
func (_Assets *AssetsSession) IDXLR() (uint8, error) {
	return _Assets.Contract.IDXLR(&_Assets.CallOpts)
}

// IDXLR is a free data retrieval call binding the contract method 0x3d085f96.
//
// Solidity: function IDX_LR() constant returns(uint8)
func (_Assets *AssetsCallerSession) IDXLR() (uint8, error) {
	return _Assets.Contract.IDXLR(&_Assets.CallOpts)
}

// IDXM is a free data retrieval call binding the contract method 0x9cc62340.
//
// Solidity: function IDX_M() constant returns(uint8)
func (_Assets *AssetsCaller) IDXM(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "IDX_M")
	return *ret0, err
}

// IDXM is a free data retrieval call binding the contract method 0x9cc62340.
//
// Solidity: function IDX_M() constant returns(uint8)
func (_Assets *AssetsSession) IDXM() (uint8, error) {
	return _Assets.Contract.IDXM(&_Assets.CallOpts)
}

// IDXM is a free data retrieval call binding the contract method 0x9cc62340.
//
// Solidity: function IDX_M() constant returns(uint8)
func (_Assets *AssetsCallerSession) IDXM() (uint8, error) {
	return _Assets.Contract.IDXM(&_Assets.CallOpts)
}

// IDXMD is a free data retrieval call binding the contract method 0x003e3223.
//
// Solidity: function IDX_MD() constant returns(uint8)
func (_Assets *AssetsCaller) IDXMD(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "IDX_MD")
	return *ret0, err
}

// IDXMD is a free data retrieval call binding the contract method 0x003e3223.
//
// Solidity: function IDX_MD() constant returns(uint8)
func (_Assets *AssetsSession) IDXMD() (uint8, error) {
	return _Assets.Contract.IDXMD(&_Assets.CallOpts)
}

// IDXMD is a free data retrieval call binding the contract method 0x003e3223.
//
// Solidity: function IDX_MD() constant returns(uint8)
func (_Assets *AssetsCallerSession) IDXMD() (uint8, error) {
	return _Assets.Contract.IDXMD(&_Assets.CallOpts)
}

// IDXMF is a free data retrieval call binding the contract method 0x8f3db436.
//
// Solidity: function IDX_MF() constant returns(uint8)
func (_Assets *AssetsCaller) IDXMF(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "IDX_MF")
	return *ret0, err
}

// IDXMF is a free data retrieval call binding the contract method 0x8f3db436.
//
// Solidity: function IDX_MF() constant returns(uint8)
func (_Assets *AssetsSession) IDXMF() (uint8, error) {
	return _Assets.Contract.IDXMF(&_Assets.CallOpts)
}

// IDXMF is a free data retrieval call binding the contract method 0x8f3db436.
//
// Solidity: function IDX_MF() constant returns(uint8)
func (_Assets *AssetsCallerSession) IDXMF() (uint8, error) {
	return _Assets.Contract.IDXMF(&_Assets.CallOpts)
}

// IDXR is a free data retrieval call binding the contract method 0x169d2914.
//
// Solidity: function IDX_R() constant returns(uint8)
func (_Assets *AssetsCaller) IDXR(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "IDX_R")
	return *ret0, err
}

// IDXR is a free data retrieval call binding the contract method 0x169d2914.
//
// Solidity: function IDX_R() constant returns(uint8)
func (_Assets *AssetsSession) IDXR() (uint8, error) {
	return _Assets.Contract.IDXR(&_Assets.CallOpts)
}

// IDXR is a free data retrieval call binding the contract method 0x169d2914.
//
// Solidity: function IDX_R() constant returns(uint8)
func (_Assets *AssetsCallerSession) IDXR() (uint8, error) {
	return _Assets.Contract.IDXR(&_Assets.CallOpts)
}

// LEAGUESPERDIV is a free data retrieval call binding the contract method 0x48d1e9c0.
//
// Solidity: function LEAGUES_PER_DIV() constant returns(uint8)
func (_Assets *AssetsCaller) LEAGUESPERDIV(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "LEAGUES_PER_DIV")
	return *ret0, err
}

// LEAGUESPERDIV is a free data retrieval call binding the contract method 0x48d1e9c0.
//
// Solidity: function LEAGUES_PER_DIV() constant returns(uint8)
func (_Assets *AssetsSession) LEAGUESPERDIV() (uint8, error) {
	return _Assets.Contract.LEAGUESPERDIV(&_Assets.CallOpts)
}

// LEAGUESPERDIV is a free data retrieval call binding the contract method 0x48d1e9c0.
//
// Solidity: function LEAGUES_PER_DIV() constant returns(uint8)
func (_Assets *AssetsCallerSession) LEAGUESPERDIV() (uint8, error) {
	return _Assets.Contract.LEAGUESPERDIV(&_Assets.CallOpts)
}

// MAXPLAYERAGEATBIRTH is a free data retrieval call binding the contract method 0x2a238b0a.
//
// Solidity: function MAX_PLAYER_AGE_AT_BIRTH() constant returns(uint8)
func (_Assets *AssetsCaller) MAXPLAYERAGEATBIRTH(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "MAX_PLAYER_AGE_AT_BIRTH")
	return *ret0, err
}

// MAXPLAYERAGEATBIRTH is a free data retrieval call binding the contract method 0x2a238b0a.
//
// Solidity: function MAX_PLAYER_AGE_AT_BIRTH() constant returns(uint8)
func (_Assets *AssetsSession) MAXPLAYERAGEATBIRTH() (uint8, error) {
	return _Assets.Contract.MAXPLAYERAGEATBIRTH(&_Assets.CallOpts)
}

// MAXPLAYERAGEATBIRTH is a free data retrieval call binding the contract method 0x2a238b0a.
//
// Solidity: function MAX_PLAYER_AGE_AT_BIRTH() constant returns(uint8)
func (_Assets *AssetsCallerSession) MAXPLAYERAGEATBIRTH() (uint8, error) {
	return _Assets.Contract.MAXPLAYERAGEATBIRTH(&_Assets.CallOpts)
}

// MINPLAYERAGEATBIRTH is a free data retrieval call binding the contract method 0xc79055d4.
//
// Solidity: function MIN_PLAYER_AGE_AT_BIRTH() constant returns(uint8)
func (_Assets *AssetsCaller) MINPLAYERAGEATBIRTH(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "MIN_PLAYER_AGE_AT_BIRTH")
	return *ret0, err
}

// MINPLAYERAGEATBIRTH is a free data retrieval call binding the contract method 0xc79055d4.
//
// Solidity: function MIN_PLAYER_AGE_AT_BIRTH() constant returns(uint8)
func (_Assets *AssetsSession) MINPLAYERAGEATBIRTH() (uint8, error) {
	return _Assets.Contract.MINPLAYERAGEATBIRTH(&_Assets.CallOpts)
}

// MINPLAYERAGEATBIRTH is a free data retrieval call binding the contract method 0xc79055d4.
//
// Solidity: function MIN_PLAYER_AGE_AT_BIRTH() constant returns(uint8)
func (_Assets *AssetsCallerSession) MINPLAYERAGEATBIRTH() (uint8, error) {
	return _Assets.Contract.MINPLAYERAGEATBIRTH(&_Assets.CallOpts)
}

// NULLADDR is a free data retrieval call binding the contract method 0xb3f390b3.
//
// Solidity: function NULL_ADDR() constant returns(address)
func (_Assets *AssetsCaller) NULLADDR(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "NULL_ADDR")
	return *ret0, err
}

// NULLADDR is a free data retrieval call binding the contract method 0xb3f390b3.
//
// Solidity: function NULL_ADDR() constant returns(address)
func (_Assets *AssetsSession) NULLADDR() (common.Address, error) {
	return _Assets.Contract.NULLADDR(&_Assets.CallOpts)
}

// NULLADDR is a free data retrieval call binding the contract method 0xb3f390b3.
//
// Solidity: function NULL_ADDR() constant returns(address)
func (_Assets *AssetsCallerSession) NULLADDR() (common.Address, error) {
	return _Assets.Contract.NULLADDR(&_Assets.CallOpts)
}

// NSKILLS is a free data retrieval call binding the contract method 0x976daaac.
//
// Solidity: function N_SKILLS() constant returns(uint8)
func (_Assets *AssetsCaller) NSKILLS(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "N_SKILLS")
	return *ret0, err
}

// NSKILLS is a free data retrieval call binding the contract method 0x976daaac.
//
// Solidity: function N_SKILLS() constant returns(uint8)
func (_Assets *AssetsSession) NSKILLS() (uint8, error) {
	return _Assets.Contract.NSKILLS(&_Assets.CallOpts)
}

// NSKILLS is a free data retrieval call binding the contract method 0x976daaac.
//
// Solidity: function N_SKILLS() constant returns(uint8)
func (_Assets *AssetsCallerSession) NSKILLS() (uint8, error) {
	return _Assets.Contract.NSKILLS(&_Assets.CallOpts)
}

// PLAYERSPERTEAMINIT is a free data retrieval call binding the contract method 0x83c31d3b.
//
// Solidity: function PLAYERS_PER_TEAM_INIT() constant returns(uint8)
func (_Assets *AssetsCaller) PLAYERSPERTEAMINIT(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "PLAYERS_PER_TEAM_INIT")
	return *ret0, err
}

// PLAYERSPERTEAMINIT is a free data retrieval call binding the contract method 0x83c31d3b.
//
// Solidity: function PLAYERS_PER_TEAM_INIT() constant returns(uint8)
func (_Assets *AssetsSession) PLAYERSPERTEAMINIT() (uint8, error) {
	return _Assets.Contract.PLAYERSPERTEAMINIT(&_Assets.CallOpts)
}

// PLAYERSPERTEAMINIT is a free data retrieval call binding the contract method 0x83c31d3b.
//
// Solidity: function PLAYERS_PER_TEAM_INIT() constant returns(uint8)
func (_Assets *AssetsCallerSession) PLAYERSPERTEAMINIT() (uint8, error) {
	return _Assets.Contract.PLAYERSPERTEAMINIT(&_Assets.CallOpts)
}

// PLAYERSPERTEAMMAX is a free data retrieval call binding the contract method 0x8adddc9d.
//
// Solidity: function PLAYERS_PER_TEAM_MAX() constant returns(uint8)
func (_Assets *AssetsCaller) PLAYERSPERTEAMMAX(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "PLAYERS_PER_TEAM_MAX")
	return *ret0, err
}

// PLAYERSPERTEAMMAX is a free data retrieval call binding the contract method 0x8adddc9d.
//
// Solidity: function PLAYERS_PER_TEAM_MAX() constant returns(uint8)
func (_Assets *AssetsSession) PLAYERSPERTEAMMAX() (uint8, error) {
	return _Assets.Contract.PLAYERSPERTEAMMAX(&_Assets.CallOpts)
}

// PLAYERSPERTEAMMAX is a free data retrieval call binding the contract method 0x8adddc9d.
//
// Solidity: function PLAYERS_PER_TEAM_MAX() constant returns(uint8)
func (_Assets *AssetsCallerSession) PLAYERSPERTEAMMAX() (uint8, error) {
	return _Assets.Contract.PLAYERSPERTEAMMAX(&_Assets.CallOpts)
}

// SKDEF is a free data retrieval call binding the contract method 0xe81e21bb.
//
// Solidity: function SK_DEF() constant returns(uint8)
func (_Assets *AssetsCaller) SKDEF(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "SK_DEF")
	return *ret0, err
}

// SKDEF is a free data retrieval call binding the contract method 0xe81e21bb.
//
// Solidity: function SK_DEF() constant returns(uint8)
func (_Assets *AssetsSession) SKDEF() (uint8, error) {
	return _Assets.Contract.SKDEF(&_Assets.CallOpts)
}

// SKDEF is a free data retrieval call binding the contract method 0xe81e21bb.
//
// Solidity: function SK_DEF() constant returns(uint8)
func (_Assets *AssetsCallerSession) SKDEF() (uint8, error) {
	return _Assets.Contract.SKDEF(&_Assets.CallOpts)
}

// SKEND is a free data retrieval call binding the contract method 0x1884332c.
//
// Solidity: function SK_END() constant returns(uint8)
func (_Assets *AssetsCaller) SKEND(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "SK_END")
	return *ret0, err
}

// SKEND is a free data retrieval call binding the contract method 0x1884332c.
//
// Solidity: function SK_END() constant returns(uint8)
func (_Assets *AssetsSession) SKEND() (uint8, error) {
	return _Assets.Contract.SKEND(&_Assets.CallOpts)
}

// SKEND is a free data retrieval call binding the contract method 0x1884332c.
//
// Solidity: function SK_END() constant returns(uint8)
func (_Assets *AssetsCallerSession) SKEND() (uint8, error) {
	return _Assets.Contract.SKEND(&_Assets.CallOpts)
}

// SKPAS is a free data retrieval call binding the contract method 0xab1b7c5e.
//
// Solidity: function SK_PAS() constant returns(uint8)
func (_Assets *AssetsCaller) SKPAS(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "SK_PAS")
	return *ret0, err
}

// SKPAS is a free data retrieval call binding the contract method 0xab1b7c5e.
//
// Solidity: function SK_PAS() constant returns(uint8)
func (_Assets *AssetsSession) SKPAS() (uint8, error) {
	return _Assets.Contract.SKPAS(&_Assets.CallOpts)
}

// SKPAS is a free data retrieval call binding the contract method 0xab1b7c5e.
//
// Solidity: function SK_PAS() constant returns(uint8)
func (_Assets *AssetsCallerSession) SKPAS() (uint8, error) {
	return _Assets.Contract.SKPAS(&_Assets.CallOpts)
}

// SKSHO is a free data retrieval call binding the contract method 0x40cd05fd.
//
// Solidity: function SK_SHO() constant returns(uint8)
func (_Assets *AssetsCaller) SKSHO(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "SK_SHO")
	return *ret0, err
}

// SKSHO is a free data retrieval call binding the contract method 0x40cd05fd.
//
// Solidity: function SK_SHO() constant returns(uint8)
func (_Assets *AssetsSession) SKSHO() (uint8, error) {
	return _Assets.Contract.SKSHO(&_Assets.CallOpts)
}

// SKSHO is a free data retrieval call binding the contract method 0x40cd05fd.
//
// Solidity: function SK_SHO() constant returns(uint8)
func (_Assets *AssetsCallerSession) SKSHO() (uint8, error) {
	return _Assets.Contract.SKSHO(&_Assets.CallOpts)
}

// SKSPE is a free data retrieval call binding the contract method 0xf8ef7b9e.
//
// Solidity: function SK_SPE() constant returns(uint8)
func (_Assets *AssetsCaller) SKSPE(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "SK_SPE")
	return *ret0, err
}

// SKSPE is a free data retrieval call binding the contract method 0xf8ef7b9e.
//
// Solidity: function SK_SPE() constant returns(uint8)
func (_Assets *AssetsSession) SKSPE() (uint8, error) {
	return _Assets.Contract.SKSPE(&_Assets.CallOpts)
}

// SKSPE is a free data retrieval call binding the contract method 0xf8ef7b9e.
//
// Solidity: function SK_SPE() constant returns(uint8)
func (_Assets *AssetsCallerSession) SKSPE() (uint8, error) {
	return _Assets.Contract.SKSPE(&_Assets.CallOpts)
}

// TEAMSPERDIVISION is a free data retrieval call binding the contract method 0xf21f5a83.
//
// Solidity: function TEAMS_PER_DIVISION() constant returns(uint8)
func (_Assets *AssetsCaller) TEAMSPERDIVISION(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "TEAMS_PER_DIVISION")
	return *ret0, err
}

// TEAMSPERDIVISION is a free data retrieval call binding the contract method 0xf21f5a83.
//
// Solidity: function TEAMS_PER_DIVISION() constant returns(uint8)
func (_Assets *AssetsSession) TEAMSPERDIVISION() (uint8, error) {
	return _Assets.Contract.TEAMSPERDIVISION(&_Assets.CallOpts)
}

// TEAMSPERDIVISION is a free data retrieval call binding the contract method 0xf21f5a83.
//
// Solidity: function TEAMS_PER_DIVISION() constant returns(uint8)
func (_Assets *AssetsCallerSession) TEAMSPERDIVISION() (uint8, error) {
	return _Assets.Contract.TEAMSPERDIVISION(&_Assets.CallOpts)
}

// TEAMSPERLEAGUE is a free data retrieval call binding the contract method 0xac5db9ee.
//
// Solidity: function TEAMS_PER_LEAGUE() constant returns(uint8)
func (_Assets *AssetsCaller) TEAMSPERLEAGUE(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "TEAMS_PER_LEAGUE")
	return *ret0, err
}

// TEAMSPERLEAGUE is a free data retrieval call binding the contract method 0xac5db9ee.
//
// Solidity: function TEAMS_PER_LEAGUE() constant returns(uint8)
func (_Assets *AssetsSession) TEAMSPERLEAGUE() (uint8, error) {
	return _Assets.Contract.TEAMSPERLEAGUE(&_Assets.CallOpts)
}

// TEAMSPERLEAGUE is a free data retrieval call binding the contract method 0xac5db9ee.
//
// Solidity: function TEAMS_PER_LEAGUE() constant returns(uint8)
func (_Assets *AssetsCallerSession) TEAMSPERLEAGUE() (uint8, error) {
	return _Assets.Contract.TEAMSPERLEAGUE(&_Assets.CallOpts)
}

// TeamExistsInCountry is a free data retrieval call binding the contract method 0x1a6daba2.
//
// Solidity: function _teamExistsInCountry(uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) constant returns(bool)
func (_Assets *AssetsCaller) TeamExistsInCountry(opts *bind.CallOpts, timeZone uint8, countryIdxInTZ *big.Int, teamIdxInCountry *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "_teamExistsInCountry", timeZone, countryIdxInTZ, teamIdxInCountry)
	return *ret0, err
}

// TeamExistsInCountry is a free data retrieval call binding the contract method 0x1a6daba2.
//
// Solidity: function _teamExistsInCountry(uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) constant returns(bool)
func (_Assets *AssetsSession) TeamExistsInCountry(timeZone uint8, countryIdxInTZ *big.Int, teamIdxInCountry *big.Int) (bool, error) {
	return _Assets.Contract.TeamExistsInCountry(&_Assets.CallOpts, timeZone, countryIdxInTZ, teamIdxInCountry)
}

// TeamExistsInCountry is a free data retrieval call binding the contract method 0x1a6daba2.
//
// Solidity: function _teamExistsInCountry(uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) constant returns(bool)
func (_Assets *AssetsCallerSession) TeamExistsInCountry(timeZone uint8, countryIdxInTZ *big.Int, teamIdxInCountry *big.Int) (bool, error) {
	return _Assets.Contract.TeamExistsInCountry(&_Assets.CallOpts, timeZone, countryIdxInTZ, teamIdxInCountry)
}

// TimeZones is a free data retrieval call binding the contract method 0xb96b1a30.
//
// Solidity: function _timeZones(uint256 ) constant returns(uint8 nCountriesToAdd, uint8 newestOrgMapIdx, uint8 newestSkillsIdx, bytes32 scoresRoot, uint8 updateCycleIdx, uint256 lastActionsSubmissionTime, uint256 lastUpdateTime, bytes32 actionsRoot, uint256 lastMarketClosureBlockNum)
func (_Assets *AssetsCaller) TimeZones(opts *bind.CallOpts, arg0 *big.Int) (struct {
	NCountriesToAdd           uint8
	NewestOrgMapIdx           uint8
	NewestSkillsIdx           uint8
	ScoresRoot                [32]byte
	UpdateCycleIdx            uint8
	LastActionsSubmissionTime *big.Int
	LastUpdateTime            *big.Int
	ActionsRoot               [32]byte
	LastMarketClosureBlockNum *big.Int
}, error) {
	ret := new(struct {
		NCountriesToAdd           uint8
		NewestOrgMapIdx           uint8
		NewestSkillsIdx           uint8
		ScoresRoot                [32]byte
		UpdateCycleIdx            uint8
		LastActionsSubmissionTime *big.Int
		LastUpdateTime            *big.Int
		ActionsRoot               [32]byte
		LastMarketClosureBlockNum *big.Int
	})
	out := ret
	err := _Assets.contract.Call(opts, out, "_timeZones", arg0)
	return *ret, err
}

// TimeZones is a free data retrieval call binding the contract method 0xb96b1a30.
//
// Solidity: function _timeZones(uint256 ) constant returns(uint8 nCountriesToAdd, uint8 newestOrgMapIdx, uint8 newestSkillsIdx, bytes32 scoresRoot, uint8 updateCycleIdx, uint256 lastActionsSubmissionTime, uint256 lastUpdateTime, bytes32 actionsRoot, uint256 lastMarketClosureBlockNum)
func (_Assets *AssetsSession) TimeZones(arg0 *big.Int) (struct {
	NCountriesToAdd           uint8
	NewestOrgMapIdx           uint8
	NewestSkillsIdx           uint8
	ScoresRoot                [32]byte
	UpdateCycleIdx            uint8
	LastActionsSubmissionTime *big.Int
	LastUpdateTime            *big.Int
	ActionsRoot               [32]byte
	LastMarketClosureBlockNum *big.Int
}, error) {
	return _Assets.Contract.TimeZones(&_Assets.CallOpts, arg0)
}

// TimeZones is a free data retrieval call binding the contract method 0xb96b1a30.
//
// Solidity: function _timeZones(uint256 ) constant returns(uint8 nCountriesToAdd, uint8 newestOrgMapIdx, uint8 newestSkillsIdx, bytes32 scoresRoot, uint8 updateCycleIdx, uint256 lastActionsSubmissionTime, uint256 lastUpdateTime, bytes32 actionsRoot, uint256 lastMarketClosureBlockNum)
func (_Assets *AssetsCallerSession) TimeZones(arg0 *big.Int) (struct {
	NCountriesToAdd           uint8
	NewestOrgMapIdx           uint8
	NewestSkillsIdx           uint8
	ScoresRoot                [32]byte
	UpdateCycleIdx            uint8
	LastActionsSubmissionTime *big.Int
	LastUpdateTime            *big.Int
	ActionsRoot               [32]byte
	LastMarketClosureBlockNum *big.Int
}, error) {
	return _Assets.Contract.TimeZones(&_Assets.CallOpts, arg0)
}

// ComputeBirthMonth is a free data retrieval call binding the contract method 0x00aae8df.
//
// Solidity: function computeBirthMonth(uint256 dna, uint256 playerCreationMonth) constant returns(uint16, uint256)
func (_Assets *AssetsCaller) ComputeBirthMonth(opts *bind.CallOpts, dna *big.Int, playerCreationMonth *big.Int) (uint16, *big.Int, error) {
	var (
		ret0 = new(uint16)
		ret1 = new(*big.Int)
	)
	out := &[]interface{}{
		ret0,
		ret1,
	}
	err := _Assets.contract.Call(opts, out, "computeBirthMonth", dna, playerCreationMonth)
	return *ret0, *ret1, err
}

// ComputeBirthMonth is a free data retrieval call binding the contract method 0x00aae8df.
//
// Solidity: function computeBirthMonth(uint256 dna, uint256 playerCreationMonth) constant returns(uint16, uint256)
func (_Assets *AssetsSession) ComputeBirthMonth(dna *big.Int, playerCreationMonth *big.Int) (uint16, *big.Int, error) {
	return _Assets.Contract.ComputeBirthMonth(&_Assets.CallOpts, dna, playerCreationMonth)
}

// ComputeBirthMonth is a free data retrieval call binding the contract method 0x00aae8df.
//
// Solidity: function computeBirthMonth(uint256 dna, uint256 playerCreationMonth) constant returns(uint16, uint256)
func (_Assets *AssetsCallerSession) ComputeBirthMonth(dna *big.Int, playerCreationMonth *big.Int) (uint16, *big.Int, error) {
	return _Assets.Contract.ComputeBirthMonth(&_Assets.CallOpts, dna, playerCreationMonth)
}

// ComputeSkills is a free data retrieval call binding the contract method 0x547d8298.
//
// Solidity: function computeSkills(uint256 dna, uint8 shirtNum) constant returns(uint16[5], uint8[4])
func (_Assets *AssetsCaller) ComputeSkills(opts *bind.CallOpts, dna *big.Int, shirtNum uint8) ([5]uint16, [4]uint8, error) {
	var (
		ret0 = new([5]uint16)
		ret1 = new([4]uint8)
	)
	out := &[]interface{}{
		ret0,
		ret1,
	}
	err := _Assets.contract.Call(opts, out, "computeSkills", dna, shirtNum)
	return *ret0, *ret1, err
}

// ComputeSkills is a free data retrieval call binding the contract method 0x547d8298.
//
// Solidity: function computeSkills(uint256 dna, uint8 shirtNum) constant returns(uint16[5], uint8[4])
func (_Assets *AssetsSession) ComputeSkills(dna *big.Int, shirtNum uint8) ([5]uint16, [4]uint8, error) {
	return _Assets.Contract.ComputeSkills(&_Assets.CallOpts, dna, shirtNum)
}

// ComputeSkills is a free data retrieval call binding the contract method 0x547d8298.
//
// Solidity: function computeSkills(uint256 dna, uint8 shirtNum) constant returns(uint16[5], uint8[4])
func (_Assets *AssetsCallerSession) ComputeSkills(dna *big.Int, shirtNum uint8) ([5]uint16, [4]uint8, error) {
	return _Assets.Contract.ComputeSkills(&_Assets.CallOpts, dna, shirtNum)
}

// CountCountries is a free data retrieval call binding the contract method 0x0abcd3e5.
//
// Solidity: function countCountries(uint8 timeZone) constant returns(uint256)
func (_Assets *AssetsCaller) CountCountries(opts *bind.CallOpts, timeZone uint8) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "countCountries", timeZone)
	return *ret0, err
}

// CountCountries is a free data retrieval call binding the contract method 0x0abcd3e5.
//
// Solidity: function countCountries(uint8 timeZone) constant returns(uint256)
func (_Assets *AssetsSession) CountCountries(timeZone uint8) (*big.Int, error) {
	return _Assets.Contract.CountCountries(&_Assets.CallOpts, timeZone)
}

// CountCountries is a free data retrieval call binding the contract method 0x0abcd3e5.
//
// Solidity: function countCountries(uint8 timeZone) constant returns(uint256)
func (_Assets *AssetsCallerSession) CountCountries(timeZone uint8) (*big.Int, error) {
	return _Assets.Contract.CountCountries(&_Assets.CallOpts, timeZone)
}

// CountTeams is a free data retrieval call binding the contract method 0x7b2566a5.
//
// Solidity: function countTeams(uint8 timeZone, uint256 countryIdxInTZ) constant returns(uint256)
func (_Assets *AssetsCaller) CountTeams(opts *bind.CallOpts, timeZone uint8, countryIdxInTZ *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "countTeams", timeZone, countryIdxInTZ)
	return *ret0, err
}

// CountTeams is a free data retrieval call binding the contract method 0x7b2566a5.
//
// Solidity: function countTeams(uint8 timeZone, uint256 countryIdxInTZ) constant returns(uint256)
func (_Assets *AssetsSession) CountTeams(timeZone uint8, countryIdxInTZ *big.Int) (*big.Int, error) {
	return _Assets.Contract.CountTeams(&_Assets.CallOpts, timeZone, countryIdxInTZ)
}

// CountTeams is a free data retrieval call binding the contract method 0x7b2566a5.
//
// Solidity: function countTeams(uint8 timeZone, uint256 countryIdxInTZ) constant returns(uint256)
func (_Assets *AssetsCallerSession) CountTeams(timeZone uint8, countryIdxInTZ *big.Int) (*big.Int, error) {
	return _Assets.Contract.CountTeams(&_Assets.CallOpts, timeZone, countryIdxInTZ)
}

// CurrentRound is a free data retrieval call binding the contract method 0x8a19c8bc.
//
// Solidity: function currentRound() constant returns(uint256)
func (_Assets *AssetsCaller) CurrentRound(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "currentRound")
	return *ret0, err
}

// CurrentRound is a free data retrieval call binding the contract method 0x8a19c8bc.
//
// Solidity: function currentRound() constant returns(uint256)
func (_Assets *AssetsSession) CurrentRound() (*big.Int, error) {
	return _Assets.Contract.CurrentRound(&_Assets.CallOpts)
}

// CurrentRound is a free data retrieval call binding the contract method 0x8a19c8bc.
//
// Solidity: function currentRound() constant returns(uint256)
func (_Assets *AssetsCallerSession) CurrentRound() (*big.Int, error) {
	return _Assets.Contract.CurrentRound(&_Assets.CallOpts)
}

// DecodeTZCountryAndVal is a free data retrieval call binding the contract method 0x3260840b.
//
// Solidity: function decodeTZCountryAndVal(uint256 encoded) constant returns(uint8, uint256, uint256)
func (_Assets *AssetsCaller) DecodeTZCountryAndVal(opts *bind.CallOpts, encoded *big.Int) (uint8, *big.Int, *big.Int, error) {
	var (
		ret0 = new(uint8)
		ret1 = new(*big.Int)
		ret2 = new(*big.Int)
	)
	out := &[]interface{}{
		ret0,
		ret1,
		ret2,
	}
	err := _Assets.contract.Call(opts, out, "decodeTZCountryAndVal", encoded)
	return *ret0, *ret1, *ret2, err
}

// DecodeTZCountryAndVal is a free data retrieval call binding the contract method 0x3260840b.
//
// Solidity: function decodeTZCountryAndVal(uint256 encoded) constant returns(uint8, uint256, uint256)
func (_Assets *AssetsSession) DecodeTZCountryAndVal(encoded *big.Int) (uint8, *big.Int, *big.Int, error) {
	return _Assets.Contract.DecodeTZCountryAndVal(&_Assets.CallOpts, encoded)
}

// DecodeTZCountryAndVal is a free data retrieval call binding the contract method 0x3260840b.
//
// Solidity: function decodeTZCountryAndVal(uint256 encoded) constant returns(uint8, uint256, uint256)
func (_Assets *AssetsCallerSession) DecodeTZCountryAndVal(encoded *big.Int) (uint8, *big.Int, *big.Int, error) {
	return _Assets.Contract.DecodeTZCountryAndVal(&_Assets.CallOpts, encoded)
}

// DecodeTactics is a free data retrieval call binding the contract method 0xe6400ac4.
//
// Solidity: function decodeTactics(uint256 tactics) constant returns(uint8[11] lineup, bool[10] extraAttack, uint8 tacticsId)
func (_Assets *AssetsCaller) DecodeTactics(opts *bind.CallOpts, tactics *big.Int) (struct {
	Lineup      [11]uint8
	ExtraAttack [10]bool
	TacticsId   uint8
}, error) {
	ret := new(struct {
		Lineup      [11]uint8
		ExtraAttack [10]bool
		TacticsId   uint8
	})
	out := ret
	err := _Assets.contract.Call(opts, out, "decodeTactics", tactics)
	return *ret, err
}

// DecodeTactics is a free data retrieval call binding the contract method 0xe6400ac4.
//
// Solidity: function decodeTactics(uint256 tactics) constant returns(uint8[11] lineup, bool[10] extraAttack, uint8 tacticsId)
func (_Assets *AssetsSession) DecodeTactics(tactics *big.Int) (struct {
	Lineup      [11]uint8
	ExtraAttack [10]bool
	TacticsId   uint8
}, error) {
	return _Assets.Contract.DecodeTactics(&_Assets.CallOpts, tactics)
}

// DecodeTactics is a free data retrieval call binding the contract method 0xe6400ac4.
//
// Solidity: function decodeTactics(uint256 tactics) constant returns(uint8[11] lineup, bool[10] extraAttack, uint8 tacticsId)
func (_Assets *AssetsCallerSession) DecodeTactics(tactics *big.Int) (struct {
	Lineup      [11]uint8
	ExtraAttack [10]bool
	TacticsId   uint8
}, error) {
	return _Assets.Contract.DecodeTactics(&_Assets.CallOpts, tactics)
}

// EncodePlayerSkills is a free data retrieval call binding the contract method 0x9c53e3fd.
//
// Solidity: function encodePlayerSkills(uint16[5] skills, uint256 monthOfBirth, uint256 playerId, uint8[4] birthTraits, bool alignedLastHalf, bool redCardLastGame, uint8 gamesNonStopping, uint8 injuryWeeksLeft) constant returns(uint256 encoded)
func (_Assets *AssetsCaller) EncodePlayerSkills(opts *bind.CallOpts, skills [5]uint16, monthOfBirth *big.Int, playerId *big.Int, birthTraits [4]uint8, alignedLastHalf bool, redCardLastGame bool, gamesNonStopping uint8, injuryWeeksLeft uint8) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "encodePlayerSkills", skills, monthOfBirth, playerId, birthTraits, alignedLastHalf, redCardLastGame, gamesNonStopping, injuryWeeksLeft)
	return *ret0, err
}

// EncodePlayerSkills is a free data retrieval call binding the contract method 0x9c53e3fd.
//
// Solidity: function encodePlayerSkills(uint16[5] skills, uint256 monthOfBirth, uint256 playerId, uint8[4] birthTraits, bool alignedLastHalf, bool redCardLastGame, uint8 gamesNonStopping, uint8 injuryWeeksLeft) constant returns(uint256 encoded)
func (_Assets *AssetsSession) EncodePlayerSkills(skills [5]uint16, monthOfBirth *big.Int, playerId *big.Int, birthTraits [4]uint8, alignedLastHalf bool, redCardLastGame bool, gamesNonStopping uint8, injuryWeeksLeft uint8) (*big.Int, error) {
	return _Assets.Contract.EncodePlayerSkills(&_Assets.CallOpts, skills, monthOfBirth, playerId, birthTraits, alignedLastHalf, redCardLastGame, gamesNonStopping, injuryWeeksLeft)
}

// EncodePlayerSkills is a free data retrieval call binding the contract method 0x9c53e3fd.
//
// Solidity: function encodePlayerSkills(uint16[5] skills, uint256 monthOfBirth, uint256 playerId, uint8[4] birthTraits, bool alignedLastHalf, bool redCardLastGame, uint8 gamesNonStopping, uint8 injuryWeeksLeft) constant returns(uint256 encoded)
func (_Assets *AssetsCallerSession) EncodePlayerSkills(skills [5]uint16, monthOfBirth *big.Int, playerId *big.Int, birthTraits [4]uint8, alignedLastHalf bool, redCardLastGame bool, gamesNonStopping uint8, injuryWeeksLeft uint8) (*big.Int, error) {
	return _Assets.Contract.EncodePlayerSkills(&_Assets.CallOpts, skills, monthOfBirth, playerId, birthTraits, alignedLastHalf, redCardLastGame, gamesNonStopping, injuryWeeksLeft)
}

// EncodePlayerState is a free data retrieval call binding the contract method 0x9f27112a.
//
// Solidity: function encodePlayerState(uint256 playerId, uint256 currentTeamId, uint8 currentShirtNum, uint256 prevPlayerTeamId, uint256 lastSaleBlock) constant returns(uint256)
func (_Assets *AssetsCaller) EncodePlayerState(opts *bind.CallOpts, playerId *big.Int, currentTeamId *big.Int, currentShirtNum uint8, prevPlayerTeamId *big.Int, lastSaleBlock *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "encodePlayerState", playerId, currentTeamId, currentShirtNum, prevPlayerTeamId, lastSaleBlock)
	return *ret0, err
}

// EncodePlayerState is a free data retrieval call binding the contract method 0x9f27112a.
//
// Solidity: function encodePlayerState(uint256 playerId, uint256 currentTeamId, uint8 currentShirtNum, uint256 prevPlayerTeamId, uint256 lastSaleBlock) constant returns(uint256)
func (_Assets *AssetsSession) EncodePlayerState(playerId *big.Int, currentTeamId *big.Int, currentShirtNum uint8, prevPlayerTeamId *big.Int, lastSaleBlock *big.Int) (*big.Int, error) {
	return _Assets.Contract.EncodePlayerState(&_Assets.CallOpts, playerId, currentTeamId, currentShirtNum, prevPlayerTeamId, lastSaleBlock)
}

// EncodePlayerState is a free data retrieval call binding the contract method 0x9f27112a.
//
// Solidity: function encodePlayerState(uint256 playerId, uint256 currentTeamId, uint8 currentShirtNum, uint256 prevPlayerTeamId, uint256 lastSaleBlock) constant returns(uint256)
func (_Assets *AssetsCallerSession) EncodePlayerState(playerId *big.Int, currentTeamId *big.Int, currentShirtNum uint8, prevPlayerTeamId *big.Int, lastSaleBlock *big.Int) (*big.Int, error) {
	return _Assets.Contract.EncodePlayerState(&_Assets.CallOpts, playerId, currentTeamId, currentShirtNum, prevPlayerTeamId, lastSaleBlock)
}

// EncodeTZCountryAndVal is a free data retrieval call binding the contract method 0x20748ae8.
//
// Solidity: function encodeTZCountryAndVal(uint8 timeZone, uint256 countryIdxInTZ, uint256 val) constant returns(uint256)
func (_Assets *AssetsCaller) EncodeTZCountryAndVal(opts *bind.CallOpts, timeZone uint8, countryIdxInTZ *big.Int, val *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "encodeTZCountryAndVal", timeZone, countryIdxInTZ, val)
	return *ret0, err
}

// EncodeTZCountryAndVal is a free data retrieval call binding the contract method 0x20748ae8.
//
// Solidity: function encodeTZCountryAndVal(uint8 timeZone, uint256 countryIdxInTZ, uint256 val) constant returns(uint256)
func (_Assets *AssetsSession) EncodeTZCountryAndVal(timeZone uint8, countryIdxInTZ *big.Int, val *big.Int) (*big.Int, error) {
	return _Assets.Contract.EncodeTZCountryAndVal(&_Assets.CallOpts, timeZone, countryIdxInTZ, val)
}

// EncodeTZCountryAndVal is a free data retrieval call binding the contract method 0x20748ae8.
//
// Solidity: function encodeTZCountryAndVal(uint8 timeZone, uint256 countryIdxInTZ, uint256 val) constant returns(uint256)
func (_Assets *AssetsCallerSession) EncodeTZCountryAndVal(timeZone uint8, countryIdxInTZ *big.Int, val *big.Int) (*big.Int, error) {
	return _Assets.Contract.EncodeTZCountryAndVal(&_Assets.CallOpts, timeZone, countryIdxInTZ, val)
}

// EncodeTactics is a free data retrieval call binding the contract method 0xe9e71652.
//
// Solidity: function encodeTactics(uint8[11] lineup, bool[10] extraAttack, uint8 tacticsId) constant returns(uint256)
func (_Assets *AssetsCaller) EncodeTactics(opts *bind.CallOpts, lineup [11]uint8, extraAttack [10]bool, tacticsId uint8) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "encodeTactics", lineup, extraAttack, tacticsId)
	return *ret0, err
}

// EncodeTactics is a free data retrieval call binding the contract method 0xe9e71652.
//
// Solidity: function encodeTactics(uint8[11] lineup, bool[10] extraAttack, uint8 tacticsId) constant returns(uint256)
func (_Assets *AssetsSession) EncodeTactics(lineup [11]uint8, extraAttack [10]bool, tacticsId uint8) (*big.Int, error) {
	return _Assets.Contract.EncodeTactics(&_Assets.CallOpts, lineup, extraAttack, tacticsId)
}

// EncodeTactics is a free data retrieval call binding the contract method 0xe9e71652.
//
// Solidity: function encodeTactics(uint8[11] lineup, bool[10] extraAttack, uint8 tacticsId) constant returns(uint256)
func (_Assets *AssetsCallerSession) EncodeTactics(lineup [11]uint8, extraAttack [10]bool, tacticsId uint8) (*big.Int, error) {
	return _Assets.Contract.EncodeTactics(&_Assets.CallOpts, lineup, extraAttack, tacticsId)
}

// GameDeployMonth is a free data retrieval call binding the contract method 0x85982431.
//
// Solidity: function gameDeployMonth() constant returns(uint256)
func (_Assets *AssetsCaller) GameDeployMonth(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "gameDeployMonth")
	return *ret0, err
}

// GameDeployMonth is a free data retrieval call binding the contract method 0x85982431.
//
// Solidity: function gameDeployMonth() constant returns(uint256)
func (_Assets *AssetsSession) GameDeployMonth() (*big.Int, error) {
	return _Assets.Contract.GameDeployMonth(&_Assets.CallOpts)
}

// GameDeployMonth is a free data retrieval call binding the contract method 0x85982431.
//
// Solidity: function gameDeployMonth() constant returns(uint256)
func (_Assets *AssetsCallerSession) GameDeployMonth() (*big.Int, error) {
	return _Assets.Contract.GameDeployMonth(&_Assets.CallOpts)
}

// GetAggressiveness is a free data retrieval call binding the contract method 0x1fc7768f.
//
// Solidity: function getAggressiveness(uint256 encodedSkills) constant returns(uint256)
func (_Assets *AssetsCaller) GetAggressiveness(opts *bind.CallOpts, encodedSkills *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "getAggressiveness", encodedSkills)
	return *ret0, err
}

// GetAggressiveness is a free data retrieval call binding the contract method 0x1fc7768f.
//
// Solidity: function getAggressiveness(uint256 encodedSkills) constant returns(uint256)
func (_Assets *AssetsSession) GetAggressiveness(encodedSkills *big.Int) (*big.Int, error) {
	return _Assets.Contract.GetAggressiveness(&_Assets.CallOpts, encodedSkills)
}

// GetAggressiveness is a free data retrieval call binding the contract method 0x1fc7768f.
//
// Solidity: function getAggressiveness(uint256 encodedSkills) constant returns(uint256)
func (_Assets *AssetsCallerSession) GetAggressiveness(encodedSkills *big.Int) (*big.Int, error) {
	return _Assets.Contract.GetAggressiveness(&_Assets.CallOpts, encodedSkills)
}

// GetAlignedLastHalf is a free data retrieval call binding the contract method 0x673fe242.
//
// Solidity: function getAlignedLastHalf(uint256 encodedSkills) constant returns(bool)
func (_Assets *AssetsCaller) GetAlignedLastHalf(opts *bind.CallOpts, encodedSkills *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "getAlignedLastHalf", encodedSkills)
	return *ret0, err
}

// GetAlignedLastHalf is a free data retrieval call binding the contract method 0x673fe242.
//
// Solidity: function getAlignedLastHalf(uint256 encodedSkills) constant returns(bool)
func (_Assets *AssetsSession) GetAlignedLastHalf(encodedSkills *big.Int) (bool, error) {
	return _Assets.Contract.GetAlignedLastHalf(&_Assets.CallOpts, encodedSkills)
}

// GetAlignedLastHalf is a free data retrieval call binding the contract method 0x673fe242.
//
// Solidity: function getAlignedLastHalf(uint256 encodedSkills) constant returns(bool)
func (_Assets *AssetsCallerSession) GetAlignedLastHalf(encodedSkills *big.Int) (bool, error) {
	return _Assets.Contract.GetAlignedLastHalf(&_Assets.CallOpts, encodedSkills)
}

// GetCurrentShirtNum is a free data retrieval call binding the contract method 0xeb78b7b7.
//
// Solidity: function getCurrentShirtNum(uint256 playerState) constant returns(uint256)
func (_Assets *AssetsCaller) GetCurrentShirtNum(opts *bind.CallOpts, playerState *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "getCurrentShirtNum", playerState)
	return *ret0, err
}

// GetCurrentShirtNum is a free data retrieval call binding the contract method 0xeb78b7b7.
//
// Solidity: function getCurrentShirtNum(uint256 playerState) constant returns(uint256)
func (_Assets *AssetsSession) GetCurrentShirtNum(playerState *big.Int) (*big.Int, error) {
	return _Assets.Contract.GetCurrentShirtNum(&_Assets.CallOpts, playerState)
}

// GetCurrentShirtNum is a free data retrieval call binding the contract method 0xeb78b7b7.
//
// Solidity: function getCurrentShirtNum(uint256 playerState) constant returns(uint256)
func (_Assets *AssetsCallerSession) GetCurrentShirtNum(playerState *big.Int) (*big.Int, error) {
	return _Assets.Contract.GetCurrentShirtNum(&_Assets.CallOpts, playerState)
}

// GetCurrentTeamId is a free data retrieval call binding the contract method 0xcd2105e8.
//
// Solidity: function getCurrentTeamId(uint256 playerState) constant returns(uint256)
func (_Assets *AssetsCaller) GetCurrentTeamId(opts *bind.CallOpts, playerState *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "getCurrentTeamId", playerState)
	return *ret0, err
}

// GetCurrentTeamId is a free data retrieval call binding the contract method 0xcd2105e8.
//
// Solidity: function getCurrentTeamId(uint256 playerState) constant returns(uint256)
func (_Assets *AssetsSession) GetCurrentTeamId(playerState *big.Int) (*big.Int, error) {
	return _Assets.Contract.GetCurrentTeamId(&_Assets.CallOpts, playerState)
}

// GetCurrentTeamId is a free data retrieval call binding the contract method 0xcd2105e8.
//
// Solidity: function getCurrentTeamId(uint256 playerState) constant returns(uint256)
func (_Assets *AssetsCallerSession) GetCurrentTeamId(playerState *big.Int) (*big.Int, error) {
	return _Assets.Contract.GetCurrentTeamId(&_Assets.CallOpts, playerState)
}

// GetCurrentTeamIdFromPlayerId is a free data retrieval call binding the contract method 0x38c96b5c.
//
// Solidity: function getCurrentTeamIdFromPlayerId(uint256 playerId) constant returns(uint256)
func (_Assets *AssetsCaller) GetCurrentTeamIdFromPlayerId(opts *bind.CallOpts, playerId *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "getCurrentTeamIdFromPlayerId", playerId)
	return *ret0, err
}

// GetCurrentTeamIdFromPlayerId is a free data retrieval call binding the contract method 0x38c96b5c.
//
// Solidity: function getCurrentTeamIdFromPlayerId(uint256 playerId) constant returns(uint256)
func (_Assets *AssetsSession) GetCurrentTeamIdFromPlayerId(playerId *big.Int) (*big.Int, error) {
	return _Assets.Contract.GetCurrentTeamIdFromPlayerId(&_Assets.CallOpts, playerId)
}

// GetCurrentTeamIdFromPlayerId is a free data retrieval call binding the contract method 0x38c96b5c.
//
// Solidity: function getCurrentTeamIdFromPlayerId(uint256 playerId) constant returns(uint256)
func (_Assets *AssetsCallerSession) GetCurrentTeamIdFromPlayerId(playerId *big.Int) (*big.Int, error) {
	return _Assets.Contract.GetCurrentTeamIdFromPlayerId(&_Assets.CallOpts, playerId)
}

// GetDefaultPlayerIdForTeamInCountry is a free data retrieval call binding the contract method 0x228408b0.
//
// Solidity: function getDefaultPlayerIdForTeamInCountry(uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry, uint8 shirtNum) constant returns(uint256)
func (_Assets *AssetsCaller) GetDefaultPlayerIdForTeamInCountry(opts *bind.CallOpts, timeZone uint8, countryIdxInTZ *big.Int, teamIdxInCountry *big.Int, shirtNum uint8) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "getDefaultPlayerIdForTeamInCountry", timeZone, countryIdxInTZ, teamIdxInCountry, shirtNum)
	return *ret0, err
}

// GetDefaultPlayerIdForTeamInCountry is a free data retrieval call binding the contract method 0x228408b0.
//
// Solidity: function getDefaultPlayerIdForTeamInCountry(uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry, uint8 shirtNum) constant returns(uint256)
func (_Assets *AssetsSession) GetDefaultPlayerIdForTeamInCountry(timeZone uint8, countryIdxInTZ *big.Int, teamIdxInCountry *big.Int, shirtNum uint8) (*big.Int, error) {
	return _Assets.Contract.GetDefaultPlayerIdForTeamInCountry(&_Assets.CallOpts, timeZone, countryIdxInTZ, teamIdxInCountry, shirtNum)
}

// GetDefaultPlayerIdForTeamInCountry is a free data retrieval call binding the contract method 0x228408b0.
//
// Solidity: function getDefaultPlayerIdForTeamInCountry(uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry, uint8 shirtNum) constant returns(uint256)
func (_Assets *AssetsCallerSession) GetDefaultPlayerIdForTeamInCountry(timeZone uint8, countryIdxInTZ *big.Int, teamIdxInCountry *big.Int, shirtNum uint8) (*big.Int, error) {
	return _Assets.Contract.GetDefaultPlayerIdForTeamInCountry(&_Assets.CallOpts, timeZone, countryIdxInTZ, teamIdxInCountry, shirtNum)
}

// GetDefence is a free data retrieval call binding the contract method 0x51585b49.
//
// Solidity: function getDefence(uint256 encodedSkills) constant returns(uint256)
func (_Assets *AssetsCaller) GetDefence(opts *bind.CallOpts, encodedSkills *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "getDefence", encodedSkills)
	return *ret0, err
}

// GetDefence is a free data retrieval call binding the contract method 0x51585b49.
//
// Solidity: function getDefence(uint256 encodedSkills) constant returns(uint256)
func (_Assets *AssetsSession) GetDefence(encodedSkills *big.Int) (*big.Int, error) {
	return _Assets.Contract.GetDefence(&_Assets.CallOpts, encodedSkills)
}

// GetDefence is a free data retrieval call binding the contract method 0x51585b49.
//
// Solidity: function getDefence(uint256 encodedSkills) constant returns(uint256)
func (_Assets *AssetsCallerSession) GetDefence(encodedSkills *big.Int) (*big.Int, error) {
	return _Assets.Contract.GetDefence(&_Assets.CallOpts, encodedSkills)
}

// GetEndurance is a free data retrieval call binding the contract method 0x258e5d90.
//
// Solidity: function getEndurance(uint256 encodedSkills) constant returns(uint256)
func (_Assets *AssetsCaller) GetEndurance(opts *bind.CallOpts, encodedSkills *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "getEndurance", encodedSkills)
	return *ret0, err
}

// GetEndurance is a free data retrieval call binding the contract method 0x258e5d90.
//
// Solidity: function getEndurance(uint256 encodedSkills) constant returns(uint256)
func (_Assets *AssetsSession) GetEndurance(encodedSkills *big.Int) (*big.Int, error) {
	return _Assets.Contract.GetEndurance(&_Assets.CallOpts, encodedSkills)
}

// GetEndurance is a free data retrieval call binding the contract method 0x258e5d90.
//
// Solidity: function getEndurance(uint256 encodedSkills) constant returns(uint256)
func (_Assets *AssetsCallerSession) GetEndurance(encodedSkills *big.Int) (*big.Int, error) {
	return _Assets.Contract.GetEndurance(&_Assets.CallOpts, encodedSkills)
}

// GetForwardness is a free data retrieval call binding the contract method 0xc2bc41cd.
//
// Solidity: function getForwardness(uint256 encodedSkills) constant returns(uint256)
func (_Assets *AssetsCaller) GetForwardness(opts *bind.CallOpts, encodedSkills *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "getForwardness", encodedSkills)
	return *ret0, err
}

// GetForwardness is a free data retrieval call binding the contract method 0xc2bc41cd.
//
// Solidity: function getForwardness(uint256 encodedSkills) constant returns(uint256)
func (_Assets *AssetsSession) GetForwardness(encodedSkills *big.Int) (*big.Int, error) {
	return _Assets.Contract.GetForwardness(&_Assets.CallOpts, encodedSkills)
}

// GetForwardness is a free data retrieval call binding the contract method 0xc2bc41cd.
//
// Solidity: function getForwardness(uint256 encodedSkills) constant returns(uint256)
func (_Assets *AssetsCallerSession) GetForwardness(encodedSkills *big.Int) (*big.Int, error) {
	return _Assets.Contract.GetForwardness(&_Assets.CallOpts, encodedSkills)
}

// GetFreeShirt is a free data retrieval call binding the contract method 0x507b1723.
//
// Solidity: function getFreeShirt(uint256 teamId) constant returns(uint8)
func (_Assets *AssetsCaller) GetFreeShirt(opts *bind.CallOpts, teamId *big.Int) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "getFreeShirt", teamId)
	return *ret0, err
}

// GetFreeShirt is a free data retrieval call binding the contract method 0x507b1723.
//
// Solidity: function getFreeShirt(uint256 teamId) constant returns(uint8)
func (_Assets *AssetsSession) GetFreeShirt(teamId *big.Int) (uint8, error) {
	return _Assets.Contract.GetFreeShirt(&_Assets.CallOpts, teamId)
}

// GetFreeShirt is a free data retrieval call binding the contract method 0x507b1723.
//
// Solidity: function getFreeShirt(uint256 teamId) constant returns(uint8)
func (_Assets *AssetsCallerSession) GetFreeShirt(teamId *big.Int) (uint8, error) {
	return _Assets.Contract.GetFreeShirt(&_Assets.CallOpts, teamId)
}

// GetGamesNonStopping is a free data retrieval call binding the contract method 0xe804e519.
//
// Solidity: function getGamesNonStopping(uint256 encodedSkills) constant returns(uint256)
func (_Assets *AssetsCaller) GetGamesNonStopping(opts *bind.CallOpts, encodedSkills *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "getGamesNonStopping", encodedSkills)
	return *ret0, err
}

// GetGamesNonStopping is a free data retrieval call binding the contract method 0xe804e519.
//
// Solidity: function getGamesNonStopping(uint256 encodedSkills) constant returns(uint256)
func (_Assets *AssetsSession) GetGamesNonStopping(encodedSkills *big.Int) (*big.Int, error) {
	return _Assets.Contract.GetGamesNonStopping(&_Assets.CallOpts, encodedSkills)
}

// GetGamesNonStopping is a free data retrieval call binding the contract method 0xe804e519.
//
// Solidity: function getGamesNonStopping(uint256 encodedSkills) constant returns(uint256)
func (_Assets *AssetsCallerSession) GetGamesNonStopping(encodedSkills *big.Int) (*big.Int, error) {
	return _Assets.Contract.GetGamesNonStopping(&_Assets.CallOpts, encodedSkills)
}

// GetInjuryWeeksLeft is a free data retrieval call binding the contract method 0x79e76597.
//
// Solidity: function getInjuryWeeksLeft(uint256 encodedSkills) constant returns(uint256)
func (_Assets *AssetsCaller) GetInjuryWeeksLeft(opts *bind.CallOpts, encodedSkills *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "getInjuryWeeksLeft", encodedSkills)
	return *ret0, err
}

// GetInjuryWeeksLeft is a free data retrieval call binding the contract method 0x79e76597.
//
// Solidity: function getInjuryWeeksLeft(uint256 encodedSkills) constant returns(uint256)
func (_Assets *AssetsSession) GetInjuryWeeksLeft(encodedSkills *big.Int) (*big.Int, error) {
	return _Assets.Contract.GetInjuryWeeksLeft(&_Assets.CallOpts, encodedSkills)
}

// GetInjuryWeeksLeft is a free data retrieval call binding the contract method 0x79e76597.
//
// Solidity: function getInjuryWeeksLeft(uint256 encodedSkills) constant returns(uint256)
func (_Assets *AssetsCallerSession) GetInjuryWeeksLeft(encodedSkills *big.Int) (*big.Int, error) {
	return _Assets.Contract.GetInjuryWeeksLeft(&_Assets.CallOpts, encodedSkills)
}

// GetLastActionsSubmissionTime is a free data retrieval call binding the contract method 0xfa80039b.
//
// Solidity: function getLastActionsSubmissionTime(uint8 timeZone) constant returns(uint256)
func (_Assets *AssetsCaller) GetLastActionsSubmissionTime(opts *bind.CallOpts, timeZone uint8) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "getLastActionsSubmissionTime", timeZone)
	return *ret0, err
}

// GetLastActionsSubmissionTime is a free data retrieval call binding the contract method 0xfa80039b.
//
// Solidity: function getLastActionsSubmissionTime(uint8 timeZone) constant returns(uint256)
func (_Assets *AssetsSession) GetLastActionsSubmissionTime(timeZone uint8) (*big.Int, error) {
	return _Assets.Contract.GetLastActionsSubmissionTime(&_Assets.CallOpts, timeZone)
}

// GetLastActionsSubmissionTime is a free data retrieval call binding the contract method 0xfa80039b.
//
// Solidity: function getLastActionsSubmissionTime(uint8 timeZone) constant returns(uint256)
func (_Assets *AssetsCallerSession) GetLastActionsSubmissionTime(timeZone uint8) (*big.Int, error) {
	return _Assets.Contract.GetLastActionsSubmissionTime(&_Assets.CallOpts, timeZone)
}

// GetLastSaleBlock is a free data retrieval call binding the contract method 0xc566b5bc.
//
// Solidity: function getLastSaleBlock(uint256 playerState) constant returns(uint256)
func (_Assets *AssetsCaller) GetLastSaleBlock(opts *bind.CallOpts, playerState *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "getLastSaleBlock", playerState)
	return *ret0, err
}

// GetLastSaleBlock is a free data retrieval call binding the contract method 0xc566b5bc.
//
// Solidity: function getLastSaleBlock(uint256 playerState) constant returns(uint256)
func (_Assets *AssetsSession) GetLastSaleBlock(playerState *big.Int) (*big.Int, error) {
	return _Assets.Contract.GetLastSaleBlock(&_Assets.CallOpts, playerState)
}

// GetLastSaleBlock is a free data retrieval call binding the contract method 0xc566b5bc.
//
// Solidity: function getLastSaleBlock(uint256 playerState) constant returns(uint256)
func (_Assets *AssetsCallerSession) GetLastSaleBlock(playerState *big.Int) (*big.Int, error) {
	return _Assets.Contract.GetLastSaleBlock(&_Assets.CallOpts, playerState)
}

// GetLastUpdateTime is a free data retrieval call binding the contract method 0x2d0e08fd.
//
// Solidity: function getLastUpdateTime(uint8 timeZone) constant returns(uint256)
func (_Assets *AssetsCaller) GetLastUpdateTime(opts *bind.CallOpts, timeZone uint8) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "getLastUpdateTime", timeZone)
	return *ret0, err
}

// GetLastUpdateTime is a free data retrieval call binding the contract method 0x2d0e08fd.
//
// Solidity: function getLastUpdateTime(uint8 timeZone) constant returns(uint256)
func (_Assets *AssetsSession) GetLastUpdateTime(timeZone uint8) (*big.Int, error) {
	return _Assets.Contract.GetLastUpdateTime(&_Assets.CallOpts, timeZone)
}

// GetLastUpdateTime is a free data retrieval call binding the contract method 0x2d0e08fd.
//
// Solidity: function getLastUpdateTime(uint8 timeZone) constant returns(uint256)
func (_Assets *AssetsCallerSession) GetLastUpdateTime(timeZone uint8) (*big.Int, error) {
	return _Assets.Contract.GetLastUpdateTime(&_Assets.CallOpts, timeZone)
}

// GetLeftishness is a free data retrieval call binding the contract method 0x3518dd1d.
//
// Solidity: function getLeftishness(uint256 encodedSkills) constant returns(uint256)
func (_Assets *AssetsCaller) GetLeftishness(opts *bind.CallOpts, encodedSkills *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "getLeftishness", encodedSkills)
	return *ret0, err
}

// GetLeftishness is a free data retrieval call binding the contract method 0x3518dd1d.
//
// Solidity: function getLeftishness(uint256 encodedSkills) constant returns(uint256)
func (_Assets *AssetsSession) GetLeftishness(encodedSkills *big.Int) (*big.Int, error) {
	return _Assets.Contract.GetLeftishness(&_Assets.CallOpts, encodedSkills)
}

// GetLeftishness is a free data retrieval call binding the contract method 0x3518dd1d.
//
// Solidity: function getLeftishness(uint256 encodedSkills) constant returns(uint256)
func (_Assets *AssetsCallerSession) GetLeftishness(encodedSkills *big.Int) (*big.Int, error) {
	return _Assets.Contract.GetLeftishness(&_Assets.CallOpts, encodedSkills)
}

// GetMonthOfBirth is a free data retrieval call binding the contract method 0x87f1e880.
//
// Solidity: function getMonthOfBirth(uint256 encodedSkills) constant returns(uint256)
func (_Assets *AssetsCaller) GetMonthOfBirth(opts *bind.CallOpts, encodedSkills *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "getMonthOfBirth", encodedSkills)
	return *ret0, err
}

// GetMonthOfBirth is a free data retrieval call binding the contract method 0x87f1e880.
//
// Solidity: function getMonthOfBirth(uint256 encodedSkills) constant returns(uint256)
func (_Assets *AssetsSession) GetMonthOfBirth(encodedSkills *big.Int) (*big.Int, error) {
	return _Assets.Contract.GetMonthOfBirth(&_Assets.CallOpts, encodedSkills)
}

// GetMonthOfBirth is a free data retrieval call binding the contract method 0x87f1e880.
//
// Solidity: function getMonthOfBirth(uint256 encodedSkills) constant returns(uint256)
func (_Assets *AssetsCallerSession) GetMonthOfBirth(encodedSkills *big.Int) (*big.Int, error) {
	return _Assets.Contract.GetMonthOfBirth(&_Assets.CallOpts, encodedSkills)
}

// GetNCountriesInTZ is a free data retrieval call binding the contract method 0xad63bcbd.
//
// Solidity: function getNCountriesInTZ(uint8 timeZone) constant returns(uint256)
func (_Assets *AssetsCaller) GetNCountriesInTZ(opts *bind.CallOpts, timeZone uint8) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "getNCountriesInTZ", timeZone)
	return *ret0, err
}

// GetNCountriesInTZ is a free data retrieval call binding the contract method 0xad63bcbd.
//
// Solidity: function getNCountriesInTZ(uint8 timeZone) constant returns(uint256)
func (_Assets *AssetsSession) GetNCountriesInTZ(timeZone uint8) (*big.Int, error) {
	return _Assets.Contract.GetNCountriesInTZ(&_Assets.CallOpts, timeZone)
}

// GetNCountriesInTZ is a free data retrieval call binding the contract method 0xad63bcbd.
//
// Solidity: function getNCountriesInTZ(uint8 timeZone) constant returns(uint256)
func (_Assets *AssetsCallerSession) GetNCountriesInTZ(timeZone uint8) (*big.Int, error) {
	return _Assets.Contract.GetNCountriesInTZ(&_Assets.CallOpts, timeZone)
}

// GetNDivisionsInCountry is a free data retrieval call binding the contract method 0x5adb40f5.
//
// Solidity: function getNDivisionsInCountry(uint8 timeZone, uint256 countryIdxInTZ) constant returns(uint256)
func (_Assets *AssetsCaller) GetNDivisionsInCountry(opts *bind.CallOpts, timeZone uint8, countryIdxInTZ *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "getNDivisionsInCountry", timeZone, countryIdxInTZ)
	return *ret0, err
}

// GetNDivisionsInCountry is a free data retrieval call binding the contract method 0x5adb40f5.
//
// Solidity: function getNDivisionsInCountry(uint8 timeZone, uint256 countryIdxInTZ) constant returns(uint256)
func (_Assets *AssetsSession) GetNDivisionsInCountry(timeZone uint8, countryIdxInTZ *big.Int) (*big.Int, error) {
	return _Assets.Contract.GetNDivisionsInCountry(&_Assets.CallOpts, timeZone, countryIdxInTZ)
}

// GetNDivisionsInCountry is a free data retrieval call binding the contract method 0x5adb40f5.
//
// Solidity: function getNDivisionsInCountry(uint8 timeZone, uint256 countryIdxInTZ) constant returns(uint256)
func (_Assets *AssetsCallerSession) GetNDivisionsInCountry(timeZone uint8, countryIdxInTZ *big.Int) (*big.Int, error) {
	return _Assets.Contract.GetNDivisionsInCountry(&_Assets.CallOpts, timeZone, countryIdxInTZ)
}

// GetNLeaguesInCountry is a free data retrieval call binding the contract method 0xf9d0723d.
//
// Solidity: function getNLeaguesInCountry(uint8 timeZone, uint256 countryIdxInTZ) constant returns(uint256)
func (_Assets *AssetsCaller) GetNLeaguesInCountry(opts *bind.CallOpts, timeZone uint8, countryIdxInTZ *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "getNLeaguesInCountry", timeZone, countryIdxInTZ)
	return *ret0, err
}

// GetNLeaguesInCountry is a free data retrieval call binding the contract method 0xf9d0723d.
//
// Solidity: function getNLeaguesInCountry(uint8 timeZone, uint256 countryIdxInTZ) constant returns(uint256)
func (_Assets *AssetsSession) GetNLeaguesInCountry(timeZone uint8, countryIdxInTZ *big.Int) (*big.Int, error) {
	return _Assets.Contract.GetNLeaguesInCountry(&_Assets.CallOpts, timeZone, countryIdxInTZ)
}

// GetNLeaguesInCountry is a free data retrieval call binding the contract method 0xf9d0723d.
//
// Solidity: function getNLeaguesInCountry(uint8 timeZone, uint256 countryIdxInTZ) constant returns(uint256)
func (_Assets *AssetsCallerSession) GetNLeaguesInCountry(timeZone uint8, countryIdxInTZ *big.Int) (*big.Int, error) {
	return _Assets.Contract.GetNLeaguesInCountry(&_Assets.CallOpts, timeZone, countryIdxInTZ)
}

// GetNTeamsInCountry is a free data retrieval call binding the contract method 0xc04f6d53.
//
// Solidity: function getNTeamsInCountry(uint8 timeZone, uint256 countryIdxInTZ) constant returns(uint256)
func (_Assets *AssetsCaller) GetNTeamsInCountry(opts *bind.CallOpts, timeZone uint8, countryIdxInTZ *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "getNTeamsInCountry", timeZone, countryIdxInTZ)
	return *ret0, err
}

// GetNTeamsInCountry is a free data retrieval call binding the contract method 0xc04f6d53.
//
// Solidity: function getNTeamsInCountry(uint8 timeZone, uint256 countryIdxInTZ) constant returns(uint256)
func (_Assets *AssetsSession) GetNTeamsInCountry(timeZone uint8, countryIdxInTZ *big.Int) (*big.Int, error) {
	return _Assets.Contract.GetNTeamsInCountry(&_Assets.CallOpts, timeZone, countryIdxInTZ)
}

// GetNTeamsInCountry is a free data retrieval call binding the contract method 0xc04f6d53.
//
// Solidity: function getNTeamsInCountry(uint8 timeZone, uint256 countryIdxInTZ) constant returns(uint256)
func (_Assets *AssetsCallerSession) GetNTeamsInCountry(timeZone uint8, countryIdxInTZ *big.Int) (*big.Int, error) {
	return _Assets.Contract.GetNTeamsInCountry(&_Assets.CallOpts, timeZone, countryIdxInTZ)
}

// GetOwnerPlayer is a free data retrieval call binding the contract method 0x8f9da214.
//
// Solidity: function getOwnerPlayer(uint256 playerId) constant returns(address)
func (_Assets *AssetsCaller) GetOwnerPlayer(opts *bind.CallOpts, playerId *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "getOwnerPlayer", playerId)
	return *ret0, err
}

// GetOwnerPlayer is a free data retrieval call binding the contract method 0x8f9da214.
//
// Solidity: function getOwnerPlayer(uint256 playerId) constant returns(address)
func (_Assets *AssetsSession) GetOwnerPlayer(playerId *big.Int) (common.Address, error) {
	return _Assets.Contract.GetOwnerPlayer(&_Assets.CallOpts, playerId)
}

// GetOwnerPlayer is a free data retrieval call binding the contract method 0x8f9da214.
//
// Solidity: function getOwnerPlayer(uint256 playerId) constant returns(address)
func (_Assets *AssetsCallerSession) GetOwnerPlayer(playerId *big.Int) (common.Address, error) {
	return _Assets.Contract.GetOwnerPlayer(&_Assets.CallOpts, playerId)
}

// GetOwnerTeam is a free data retrieval call binding the contract method 0x492afc69.
//
// Solidity: function getOwnerTeam(uint256 teamId) constant returns(address)
func (_Assets *AssetsCaller) GetOwnerTeam(opts *bind.CallOpts, teamId *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "getOwnerTeam", teamId)
	return *ret0, err
}

// GetOwnerTeam is a free data retrieval call binding the contract method 0x492afc69.
//
// Solidity: function getOwnerTeam(uint256 teamId) constant returns(address)
func (_Assets *AssetsSession) GetOwnerTeam(teamId *big.Int) (common.Address, error) {
	return _Assets.Contract.GetOwnerTeam(&_Assets.CallOpts, teamId)
}

// GetOwnerTeam is a free data retrieval call binding the contract method 0x492afc69.
//
// Solidity: function getOwnerTeam(uint256 teamId) constant returns(address)
func (_Assets *AssetsCallerSession) GetOwnerTeam(teamId *big.Int) (common.Address, error) {
	return _Assets.Contract.GetOwnerTeam(&_Assets.CallOpts, teamId)
}

// GetOwnerTeamInCountry is a free data retrieval call binding the contract method 0x595ef25b.
//
// Solidity: function getOwnerTeamInCountry(uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) constant returns(address)
func (_Assets *AssetsCaller) GetOwnerTeamInCountry(opts *bind.CallOpts, timeZone uint8, countryIdxInTZ *big.Int, teamIdxInCountry *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "getOwnerTeamInCountry", timeZone, countryIdxInTZ, teamIdxInCountry)
	return *ret0, err
}

// GetOwnerTeamInCountry is a free data retrieval call binding the contract method 0x595ef25b.
//
// Solidity: function getOwnerTeamInCountry(uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) constant returns(address)
func (_Assets *AssetsSession) GetOwnerTeamInCountry(timeZone uint8, countryIdxInTZ *big.Int, teamIdxInCountry *big.Int) (common.Address, error) {
	return _Assets.Contract.GetOwnerTeamInCountry(&_Assets.CallOpts, timeZone, countryIdxInTZ, teamIdxInCountry)
}

// GetOwnerTeamInCountry is a free data retrieval call binding the contract method 0x595ef25b.
//
// Solidity: function getOwnerTeamInCountry(uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) constant returns(address)
func (_Assets *AssetsCallerSession) GetOwnerTeamInCountry(timeZone uint8, countryIdxInTZ *big.Int, teamIdxInCountry *big.Int) (common.Address, error) {
	return _Assets.Contract.GetOwnerTeamInCountry(&_Assets.CallOpts, timeZone, countryIdxInTZ, teamIdxInCountry)
}

// GetPass is a free data retrieval call binding the contract method 0x55a6f86f.
//
// Solidity: function getPass(uint256 encodedSkills) constant returns(uint256)
func (_Assets *AssetsCaller) GetPass(opts *bind.CallOpts, encodedSkills *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "getPass", encodedSkills)
	return *ret0, err
}

// GetPass is a free data retrieval call binding the contract method 0x55a6f86f.
//
// Solidity: function getPass(uint256 encodedSkills) constant returns(uint256)
func (_Assets *AssetsSession) GetPass(encodedSkills *big.Int) (*big.Int, error) {
	return _Assets.Contract.GetPass(&_Assets.CallOpts, encodedSkills)
}

// GetPass is a free data retrieval call binding the contract method 0x55a6f86f.
//
// Solidity: function getPass(uint256 encodedSkills) constant returns(uint256)
func (_Assets *AssetsCallerSession) GetPass(encodedSkills *big.Int) (*big.Int, error) {
	return _Assets.Contract.GetPass(&_Assets.CallOpts, encodedSkills)
}

// GetPlayerAgeInMonths is a free data retrieval call binding the contract method 0x1ffeb349.
//
// Solidity: function getPlayerAgeInMonths(uint256 playerId) constant returns(uint256)
func (_Assets *AssetsCaller) GetPlayerAgeInMonths(opts *bind.CallOpts, playerId *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "getPlayerAgeInMonths", playerId)
	return *ret0, err
}

// GetPlayerAgeInMonths is a free data retrieval call binding the contract method 0x1ffeb349.
//
// Solidity: function getPlayerAgeInMonths(uint256 playerId) constant returns(uint256)
func (_Assets *AssetsSession) GetPlayerAgeInMonths(playerId *big.Int) (*big.Int, error) {
	return _Assets.Contract.GetPlayerAgeInMonths(&_Assets.CallOpts, playerId)
}

// GetPlayerAgeInMonths is a free data retrieval call binding the contract method 0x1ffeb349.
//
// Solidity: function getPlayerAgeInMonths(uint256 playerId) constant returns(uint256)
func (_Assets *AssetsCallerSession) GetPlayerAgeInMonths(playerId *big.Int) (*big.Int, error) {
	return _Assets.Contract.GetPlayerAgeInMonths(&_Assets.CallOpts, playerId)
}

// GetPlayerIdFromSkills is a free data retrieval call binding the contract method 0x6f6c2ae0.
//
// Solidity: function getPlayerIdFromSkills(uint256 encodedSkills) constant returns(uint256)
func (_Assets *AssetsCaller) GetPlayerIdFromSkills(opts *bind.CallOpts, encodedSkills *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "getPlayerIdFromSkills", encodedSkills)
	return *ret0, err
}

// GetPlayerIdFromSkills is a free data retrieval call binding the contract method 0x6f6c2ae0.
//
// Solidity: function getPlayerIdFromSkills(uint256 encodedSkills) constant returns(uint256)
func (_Assets *AssetsSession) GetPlayerIdFromSkills(encodedSkills *big.Int) (*big.Int, error) {
	return _Assets.Contract.GetPlayerIdFromSkills(&_Assets.CallOpts, encodedSkills)
}

// GetPlayerIdFromSkills is a free data retrieval call binding the contract method 0x6f6c2ae0.
//
// Solidity: function getPlayerIdFromSkills(uint256 encodedSkills) constant returns(uint256)
func (_Assets *AssetsCallerSession) GetPlayerIdFromSkills(encodedSkills *big.Int) (*big.Int, error) {
	return _Assets.Contract.GetPlayerIdFromSkills(&_Assets.CallOpts, encodedSkills)
}

// GetPlayerIdFromState is a free data retrieval call binding the contract method 0x78f4c718.
//
// Solidity: function getPlayerIdFromState(uint256 playerState) constant returns(uint256)
func (_Assets *AssetsCaller) GetPlayerIdFromState(opts *bind.CallOpts, playerState *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "getPlayerIdFromState", playerState)
	return *ret0, err
}

// GetPlayerIdFromState is a free data retrieval call binding the contract method 0x78f4c718.
//
// Solidity: function getPlayerIdFromState(uint256 playerState) constant returns(uint256)
func (_Assets *AssetsSession) GetPlayerIdFromState(playerState *big.Int) (*big.Int, error) {
	return _Assets.Contract.GetPlayerIdFromState(&_Assets.CallOpts, playerState)
}

// GetPlayerIdFromState is a free data retrieval call binding the contract method 0x78f4c718.
//
// Solidity: function getPlayerIdFromState(uint256 playerState) constant returns(uint256)
func (_Assets *AssetsCallerSession) GetPlayerIdFromState(playerState *big.Int) (*big.Int, error) {
	return _Assets.Contract.GetPlayerIdFromState(&_Assets.CallOpts, playerState)
}

// GetPlayerIdsInTeam is a free data retrieval call binding the contract method 0xeabf6a4b.
//
// Solidity: function getPlayerIdsInTeam(uint256 teamId) constant returns(uint256[25] playerIds)
func (_Assets *AssetsCaller) GetPlayerIdsInTeam(opts *bind.CallOpts, teamId *big.Int) ([25]*big.Int, error) {
	var (
		ret0 = new([25]*big.Int)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "getPlayerIdsInTeam", teamId)
	return *ret0, err
}

// GetPlayerIdsInTeam is a free data retrieval call binding the contract method 0xeabf6a4b.
//
// Solidity: function getPlayerIdsInTeam(uint256 teamId) constant returns(uint256[25] playerIds)
func (_Assets *AssetsSession) GetPlayerIdsInTeam(teamId *big.Int) ([25]*big.Int, error) {
	return _Assets.Contract.GetPlayerIdsInTeam(&_Assets.CallOpts, teamId)
}

// GetPlayerIdsInTeam is a free data retrieval call binding the contract method 0xeabf6a4b.
//
// Solidity: function getPlayerIdsInTeam(uint256 teamId) constant returns(uint256[25] playerIds)
func (_Assets *AssetsCallerSession) GetPlayerIdsInTeam(teamId *big.Int) ([25]*big.Int, error) {
	return _Assets.Contract.GetPlayerIdsInTeam(&_Assets.CallOpts, teamId)
}

// GetPlayerSkillsAtBirth is a free data retrieval call binding the contract method 0xc73f808d.
//
// Solidity: function getPlayerSkillsAtBirth(uint256 playerId) constant returns(uint256)
func (_Assets *AssetsCaller) GetPlayerSkillsAtBirth(opts *bind.CallOpts, playerId *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "getPlayerSkillsAtBirth", playerId)
	return *ret0, err
}

// GetPlayerSkillsAtBirth is a free data retrieval call binding the contract method 0xc73f808d.
//
// Solidity: function getPlayerSkillsAtBirth(uint256 playerId) constant returns(uint256)
func (_Assets *AssetsSession) GetPlayerSkillsAtBirth(playerId *big.Int) (*big.Int, error) {
	return _Assets.Contract.GetPlayerSkillsAtBirth(&_Assets.CallOpts, playerId)
}

// GetPlayerSkillsAtBirth is a free data retrieval call binding the contract method 0xc73f808d.
//
// Solidity: function getPlayerSkillsAtBirth(uint256 playerId) constant returns(uint256)
func (_Assets *AssetsCallerSession) GetPlayerSkillsAtBirth(playerId *big.Int) (*big.Int, error) {
	return _Assets.Contract.GetPlayerSkillsAtBirth(&_Assets.CallOpts, playerId)
}

// GetPlayerState is a free data retrieval call binding the contract method 0xec7ecec5.
//
// Solidity: function getPlayerState(uint256 playerId) constant returns(uint256)
func (_Assets *AssetsCaller) GetPlayerState(opts *bind.CallOpts, playerId *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "getPlayerState", playerId)
	return *ret0, err
}

// GetPlayerState is a free data retrieval call binding the contract method 0xec7ecec5.
//
// Solidity: function getPlayerState(uint256 playerId) constant returns(uint256)
func (_Assets *AssetsSession) GetPlayerState(playerId *big.Int) (*big.Int, error) {
	return _Assets.Contract.GetPlayerState(&_Assets.CallOpts, playerId)
}

// GetPlayerState is a free data retrieval call binding the contract method 0xec7ecec5.
//
// Solidity: function getPlayerState(uint256 playerId) constant returns(uint256)
func (_Assets *AssetsCallerSession) GetPlayerState(playerId *big.Int) (*big.Int, error) {
	return _Assets.Contract.GetPlayerState(&_Assets.CallOpts, playerId)
}

// GetPlayerStateAtBirth is a free data retrieval call binding the contract method 0x26657608.
//
// Solidity: function getPlayerStateAtBirth(uint256 playerId) constant returns(uint256)
func (_Assets *AssetsCaller) GetPlayerStateAtBirth(opts *bind.CallOpts, playerId *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "getPlayerStateAtBirth", playerId)
	return *ret0, err
}

// GetPlayerStateAtBirth is a free data retrieval call binding the contract method 0x26657608.
//
// Solidity: function getPlayerStateAtBirth(uint256 playerId) constant returns(uint256)
func (_Assets *AssetsSession) GetPlayerStateAtBirth(playerId *big.Int) (*big.Int, error) {
	return _Assets.Contract.GetPlayerStateAtBirth(&_Assets.CallOpts, playerId)
}

// GetPlayerStateAtBirth is a free data retrieval call binding the contract method 0x26657608.
//
// Solidity: function getPlayerStateAtBirth(uint256 playerId) constant returns(uint256)
func (_Assets *AssetsCallerSession) GetPlayerStateAtBirth(playerId *big.Int) (*big.Int, error) {
	return _Assets.Contract.GetPlayerStateAtBirth(&_Assets.CallOpts, playerId)
}

// GetPotential is a free data retrieval call binding the contract method 0x434807ad.
//
// Solidity: function getPotential(uint256 encodedSkills) constant returns(uint256)
func (_Assets *AssetsCaller) GetPotential(opts *bind.CallOpts, encodedSkills *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "getPotential", encodedSkills)
	return *ret0, err
}

// GetPotential is a free data retrieval call binding the contract method 0x434807ad.
//
// Solidity: function getPotential(uint256 encodedSkills) constant returns(uint256)
func (_Assets *AssetsSession) GetPotential(encodedSkills *big.Int) (*big.Int, error) {
	return _Assets.Contract.GetPotential(&_Assets.CallOpts, encodedSkills)
}

// GetPotential is a free data retrieval call binding the contract method 0x434807ad.
//
// Solidity: function getPotential(uint256 encodedSkills) constant returns(uint256)
func (_Assets *AssetsCallerSession) GetPotential(encodedSkills *big.Int) (*big.Int, error) {
	return _Assets.Contract.GetPotential(&_Assets.CallOpts, encodedSkills)
}

// GetPrevPlayerTeamId is a free data retrieval call binding the contract method 0x4bea2a69.
//
// Solidity: function getPrevPlayerTeamId(uint256 playerState) constant returns(uint256)
func (_Assets *AssetsCaller) GetPrevPlayerTeamId(opts *bind.CallOpts, playerState *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "getPrevPlayerTeamId", playerState)
	return *ret0, err
}

// GetPrevPlayerTeamId is a free data retrieval call binding the contract method 0x4bea2a69.
//
// Solidity: function getPrevPlayerTeamId(uint256 playerState) constant returns(uint256)
func (_Assets *AssetsSession) GetPrevPlayerTeamId(playerState *big.Int) (*big.Int, error) {
	return _Assets.Contract.GetPrevPlayerTeamId(&_Assets.CallOpts, playerState)
}

// GetPrevPlayerTeamId is a free data retrieval call binding the contract method 0x4bea2a69.
//
// Solidity: function getPrevPlayerTeamId(uint256 playerState) constant returns(uint256)
func (_Assets *AssetsCallerSession) GetPrevPlayerTeamId(playerState *big.Int) (*big.Int, error) {
	return _Assets.Contract.GetPrevPlayerTeamId(&_Assets.CallOpts, playerState)
}

// GetRedCardLastGame is a free data retrieval call binding the contract method 0xcc7d473b.
//
// Solidity: function getRedCardLastGame(uint256 encodedSkills) constant returns(bool)
func (_Assets *AssetsCaller) GetRedCardLastGame(opts *bind.CallOpts, encodedSkills *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "getRedCardLastGame", encodedSkills)
	return *ret0, err
}

// GetRedCardLastGame is a free data retrieval call binding the contract method 0xcc7d473b.
//
// Solidity: function getRedCardLastGame(uint256 encodedSkills) constant returns(bool)
func (_Assets *AssetsSession) GetRedCardLastGame(encodedSkills *big.Int) (bool, error) {
	return _Assets.Contract.GetRedCardLastGame(&_Assets.CallOpts, encodedSkills)
}

// GetRedCardLastGame is a free data retrieval call binding the contract method 0xcc7d473b.
//
// Solidity: function getRedCardLastGame(uint256 encodedSkills) constant returns(bool)
func (_Assets *AssetsCallerSession) GetRedCardLastGame(encodedSkills *big.Int) (bool, error) {
	return _Assets.Contract.GetRedCardLastGame(&_Assets.CallOpts, encodedSkills)
}

// GetShoot is a free data retrieval call binding the contract method 0x65b4b476.
//
// Solidity: function getShoot(uint256 encodedSkills) constant returns(uint256)
func (_Assets *AssetsCaller) GetShoot(opts *bind.CallOpts, encodedSkills *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "getShoot", encodedSkills)
	return *ret0, err
}

// GetShoot is a free data retrieval call binding the contract method 0x65b4b476.
//
// Solidity: function getShoot(uint256 encodedSkills) constant returns(uint256)
func (_Assets *AssetsSession) GetShoot(encodedSkills *big.Int) (*big.Int, error) {
	return _Assets.Contract.GetShoot(&_Assets.CallOpts, encodedSkills)
}

// GetShoot is a free data retrieval call binding the contract method 0x65b4b476.
//
// Solidity: function getShoot(uint256 encodedSkills) constant returns(uint256)
func (_Assets *AssetsCallerSession) GetShoot(encodedSkills *big.Int) (*big.Int, error) {
	return _Assets.Contract.GetShoot(&_Assets.CallOpts, encodedSkills)
}

// GetSkills is a free data retrieval call binding the contract method 0x0092bf78.
//
// Solidity: function getSkills(uint256 encodedSkills) constant returns(uint256)
func (_Assets *AssetsCaller) GetSkills(opts *bind.CallOpts, encodedSkills *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "getSkills", encodedSkills)
	return *ret0, err
}

// GetSkills is a free data retrieval call binding the contract method 0x0092bf78.
//
// Solidity: function getSkills(uint256 encodedSkills) constant returns(uint256)
func (_Assets *AssetsSession) GetSkills(encodedSkills *big.Int) (*big.Int, error) {
	return _Assets.Contract.GetSkills(&_Assets.CallOpts, encodedSkills)
}

// GetSkills is a free data retrieval call binding the contract method 0x0092bf78.
//
// Solidity: function getSkills(uint256 encodedSkills) constant returns(uint256)
func (_Assets *AssetsCallerSession) GetSkills(encodedSkills *big.Int) (*big.Int, error) {
	return _Assets.Contract.GetSkills(&_Assets.CallOpts, encodedSkills)
}

// GetSkillsVec is a free data retrieval call binding the contract method 0xcc1cc3d7.
//
// Solidity: function getSkillsVec(uint256 encodedSkills) constant returns(uint16[5] skills)
func (_Assets *AssetsCaller) GetSkillsVec(opts *bind.CallOpts, encodedSkills *big.Int) ([5]uint16, error) {
	var (
		ret0 = new([5]uint16)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "getSkillsVec", encodedSkills)
	return *ret0, err
}

// GetSkillsVec is a free data retrieval call binding the contract method 0xcc1cc3d7.
//
// Solidity: function getSkillsVec(uint256 encodedSkills) constant returns(uint16[5] skills)
func (_Assets *AssetsSession) GetSkillsVec(encodedSkills *big.Int) ([5]uint16, error) {
	return _Assets.Contract.GetSkillsVec(&_Assets.CallOpts, encodedSkills)
}

// GetSkillsVec is a free data retrieval call binding the contract method 0xcc1cc3d7.
//
// Solidity: function getSkillsVec(uint256 encodedSkills) constant returns(uint16[5] skills)
func (_Assets *AssetsCallerSession) GetSkillsVec(encodedSkills *big.Int) ([5]uint16, error) {
	return _Assets.Contract.GetSkillsVec(&_Assets.CallOpts, encodedSkills)
}

// GetSpeed is a free data retrieval call binding the contract method 0x4b93f753.
//
// Solidity: function getSpeed(uint256 encodedSkills) constant returns(uint256)
func (_Assets *AssetsCaller) GetSpeed(opts *bind.CallOpts, encodedSkills *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "getSpeed", encodedSkills)
	return *ret0, err
}

// GetSpeed is a free data retrieval call binding the contract method 0x4b93f753.
//
// Solidity: function getSpeed(uint256 encodedSkills) constant returns(uint256)
func (_Assets *AssetsSession) GetSpeed(encodedSkills *big.Int) (*big.Int, error) {
	return _Assets.Contract.GetSpeed(&_Assets.CallOpts, encodedSkills)
}

// GetSpeed is a free data retrieval call binding the contract method 0x4b93f753.
//
// Solidity: function getSpeed(uint256 encodedSkills) constant returns(uint256)
func (_Assets *AssetsCallerSession) GetSpeed(encodedSkills *big.Int) (*big.Int, error) {
	return _Assets.Contract.GetSpeed(&_Assets.CallOpts, encodedSkills)
}

// GetSumOfSkills is a free data retrieval call binding the contract method 0x1060c9c2.
//
// Solidity: function getSumOfSkills(uint256 encodedSkills) constant returns(uint256)
func (_Assets *AssetsCaller) GetSumOfSkills(opts *bind.CallOpts, encodedSkills *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "getSumOfSkills", encodedSkills)
	return *ret0, err
}

// GetSumOfSkills is a free data retrieval call binding the contract method 0x1060c9c2.
//
// Solidity: function getSumOfSkills(uint256 encodedSkills) constant returns(uint256)
func (_Assets *AssetsSession) GetSumOfSkills(encodedSkills *big.Int) (*big.Int, error) {
	return _Assets.Contract.GetSumOfSkills(&_Assets.CallOpts, encodedSkills)
}

// GetSumOfSkills is a free data retrieval call binding the contract method 0x1060c9c2.
//
// Solidity: function getSumOfSkills(uint256 encodedSkills) constant returns(uint256)
func (_Assets *AssetsCallerSession) GetSumOfSkills(encodedSkills *big.Int) (*big.Int, error) {
	return _Assets.Contract.GetSumOfSkills(&_Assets.CallOpts, encodedSkills)
}

// IsBotTeam is a free data retrieval call binding the contract method 0x8cc9a8d5.
//
// Solidity: function isBotTeam(uint256 teamId) constant returns(bool)
func (_Assets *AssetsCaller) IsBotTeam(opts *bind.CallOpts, teamId *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "isBotTeam", teamId)
	return *ret0, err
}

// IsBotTeam is a free data retrieval call binding the contract method 0x8cc9a8d5.
//
// Solidity: function isBotTeam(uint256 teamId) constant returns(bool)
func (_Assets *AssetsSession) IsBotTeam(teamId *big.Int) (bool, error) {
	return _Assets.Contract.IsBotTeam(&_Assets.CallOpts, teamId)
}

// IsBotTeam is a free data retrieval call binding the contract method 0x8cc9a8d5.
//
// Solidity: function isBotTeam(uint256 teamId) constant returns(bool)
func (_Assets *AssetsCallerSession) IsBotTeam(teamId *big.Int) (bool, error) {
	return _Assets.Contract.IsBotTeam(&_Assets.CallOpts, teamId)
}

// IsBotTeamInCountry is a free data retrieval call binding the contract method 0x80bac709.
//
// Solidity: function isBotTeamInCountry(uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) constant returns(bool)
func (_Assets *AssetsCaller) IsBotTeamInCountry(opts *bind.CallOpts, timeZone uint8, countryIdxInTZ *big.Int, teamIdxInCountry *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "isBotTeamInCountry", timeZone, countryIdxInTZ, teamIdxInCountry)
	return *ret0, err
}

// IsBotTeamInCountry is a free data retrieval call binding the contract method 0x80bac709.
//
// Solidity: function isBotTeamInCountry(uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) constant returns(bool)
func (_Assets *AssetsSession) IsBotTeamInCountry(timeZone uint8, countryIdxInTZ *big.Int, teamIdxInCountry *big.Int) (bool, error) {
	return _Assets.Contract.IsBotTeamInCountry(&_Assets.CallOpts, timeZone, countryIdxInTZ, teamIdxInCountry)
}

// IsBotTeamInCountry is a free data retrieval call binding the contract method 0x80bac709.
//
// Solidity: function isBotTeamInCountry(uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) constant returns(bool)
func (_Assets *AssetsCallerSession) IsBotTeamInCountry(timeZone uint8, countryIdxInTZ *big.Int, teamIdxInCountry *big.Int) (bool, error) {
	return _Assets.Contract.IsBotTeamInCountry(&_Assets.CallOpts, timeZone, countryIdxInTZ, teamIdxInCountry)
}

// IsFreeShirt is a free data retrieval call binding the contract method 0x963fcc80.
//
// Solidity: function isFreeShirt(uint256 teamId, uint8 shirtNum) constant returns(bool)
func (_Assets *AssetsCaller) IsFreeShirt(opts *bind.CallOpts, teamId *big.Int, shirtNum uint8) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "isFreeShirt", teamId, shirtNum)
	return *ret0, err
}

// IsFreeShirt is a free data retrieval call binding the contract method 0x963fcc80.
//
// Solidity: function isFreeShirt(uint256 teamId, uint8 shirtNum) constant returns(bool)
func (_Assets *AssetsSession) IsFreeShirt(teamId *big.Int, shirtNum uint8) (bool, error) {
	return _Assets.Contract.IsFreeShirt(&_Assets.CallOpts, teamId, shirtNum)
}

// IsFreeShirt is a free data retrieval call binding the contract method 0x963fcc80.
//
// Solidity: function isFreeShirt(uint256 teamId, uint8 shirtNum) constant returns(bool)
func (_Assets *AssetsCallerSession) IsFreeShirt(teamId *big.Int, shirtNum uint8) (bool, error) {
	return _Assets.Contract.IsFreeShirt(&_Assets.CallOpts, teamId, shirtNum)
}

// IsVirtualPlayer is a free data retrieval call binding the contract method 0xb32aa2c1.
//
// Solidity: function isVirtualPlayer(uint256 playerId) constant returns(bool)
func (_Assets *AssetsCaller) IsVirtualPlayer(opts *bind.CallOpts, playerId *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "isVirtualPlayer", playerId)
	return *ret0, err
}

// IsVirtualPlayer is a free data retrieval call binding the contract method 0xb32aa2c1.
//
// Solidity: function isVirtualPlayer(uint256 playerId) constant returns(bool)
func (_Assets *AssetsSession) IsVirtualPlayer(playerId *big.Int) (bool, error) {
	return _Assets.Contract.IsVirtualPlayer(&_Assets.CallOpts, playerId)
}

// IsVirtualPlayer is a free data retrieval call binding the contract method 0xb32aa2c1.
//
// Solidity: function isVirtualPlayer(uint256 playerId) constant returns(bool)
func (_Assets *AssetsCallerSession) IsVirtualPlayer(playerId *big.Int) (bool, error) {
	return _Assets.Contract.IsVirtualPlayer(&_Assets.CallOpts, playerId)
}

// PlayerExists is a free data retrieval call binding the contract method 0xbc1a97c1.
//
// Solidity: function playerExists(uint256 playerId) constant returns(bool)
func (_Assets *AssetsCaller) PlayerExists(opts *bind.CallOpts, playerId *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "playerExists", playerId)
	return *ret0, err
}

// PlayerExists is a free data retrieval call binding the contract method 0xbc1a97c1.
//
// Solidity: function playerExists(uint256 playerId) constant returns(bool)
func (_Assets *AssetsSession) PlayerExists(playerId *big.Int) (bool, error) {
	return _Assets.Contract.PlayerExists(&_Assets.CallOpts, playerId)
}

// PlayerExists is a free data retrieval call binding the contract method 0xbc1a97c1.
//
// Solidity: function playerExists(uint256 playerId) constant returns(bool)
func (_Assets *AssetsCallerSession) PlayerExists(playerId *big.Int) (bool, error) {
	return _Assets.Contract.PlayerExists(&_Assets.CallOpts, playerId)
}

// SetCurrentShirtNum is a free data retrieval call binding the contract method 0x4db989fd.
//
// Solidity: function setCurrentShirtNum(uint256 state, uint8 currentShirtNum) constant returns(uint256)
func (_Assets *AssetsCaller) SetCurrentShirtNum(opts *bind.CallOpts, state *big.Int, currentShirtNum uint8) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "setCurrentShirtNum", state, currentShirtNum)
	return *ret0, err
}

// SetCurrentShirtNum is a free data retrieval call binding the contract method 0x4db989fd.
//
// Solidity: function setCurrentShirtNum(uint256 state, uint8 currentShirtNum) constant returns(uint256)
func (_Assets *AssetsSession) SetCurrentShirtNum(state *big.Int, currentShirtNum uint8) (*big.Int, error) {
	return _Assets.Contract.SetCurrentShirtNum(&_Assets.CallOpts, state, currentShirtNum)
}

// SetCurrentShirtNum is a free data retrieval call binding the contract method 0x4db989fd.
//
// Solidity: function setCurrentShirtNum(uint256 state, uint8 currentShirtNum) constant returns(uint256)
func (_Assets *AssetsCallerSession) SetCurrentShirtNum(state *big.Int, currentShirtNum uint8) (*big.Int, error) {
	return _Assets.Contract.SetCurrentShirtNum(&_Assets.CallOpts, state, currentShirtNum)
}

// SetCurrentTeamId is a free data retrieval call binding the contract method 0xc37b1c25.
//
// Solidity: function setCurrentTeamId(uint256 playerState, uint256 teamId) constant returns(uint256)
func (_Assets *AssetsCaller) SetCurrentTeamId(opts *bind.CallOpts, playerState *big.Int, teamId *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "setCurrentTeamId", playerState, teamId)
	return *ret0, err
}

// SetCurrentTeamId is a free data retrieval call binding the contract method 0xc37b1c25.
//
// Solidity: function setCurrentTeamId(uint256 playerState, uint256 teamId) constant returns(uint256)
func (_Assets *AssetsSession) SetCurrentTeamId(playerState *big.Int, teamId *big.Int) (*big.Int, error) {
	return _Assets.Contract.SetCurrentTeamId(&_Assets.CallOpts, playerState, teamId)
}

// SetCurrentTeamId is a free data retrieval call binding the contract method 0xc37b1c25.
//
// Solidity: function setCurrentTeamId(uint256 playerState, uint256 teamId) constant returns(uint256)
func (_Assets *AssetsCallerSession) SetCurrentTeamId(playerState *big.Int, teamId *big.Int) (*big.Int, error) {
	return _Assets.Contract.SetCurrentTeamId(&_Assets.CallOpts, playerState, teamId)
}

// SetLastSaleBlock is a free data retrieval call binding the contract method 0xaf76cd01.
//
// Solidity: function setLastSaleBlock(uint256 state, uint256 lastSaleBlock) constant returns(uint256)
func (_Assets *AssetsCaller) SetLastSaleBlock(opts *bind.CallOpts, state *big.Int, lastSaleBlock *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "setLastSaleBlock", state, lastSaleBlock)
	return *ret0, err
}

// SetLastSaleBlock is a free data retrieval call binding the contract method 0xaf76cd01.
//
// Solidity: function setLastSaleBlock(uint256 state, uint256 lastSaleBlock) constant returns(uint256)
func (_Assets *AssetsSession) SetLastSaleBlock(state *big.Int, lastSaleBlock *big.Int) (*big.Int, error) {
	return _Assets.Contract.SetLastSaleBlock(&_Assets.CallOpts, state, lastSaleBlock)
}

// SetLastSaleBlock is a free data retrieval call binding the contract method 0xaf76cd01.
//
// Solidity: function setLastSaleBlock(uint256 state, uint256 lastSaleBlock) constant returns(uint256)
func (_Assets *AssetsCallerSession) SetLastSaleBlock(state *big.Int, lastSaleBlock *big.Int) (*big.Int, error) {
	return _Assets.Contract.SetLastSaleBlock(&_Assets.CallOpts, state, lastSaleBlock)
}

// SetPrevPlayerTeamId is a free data retrieval call binding the contract method 0x37a86302.
//
// Solidity: function setPrevPlayerTeamId(uint256 state, uint256 value) constant returns(uint256)
func (_Assets *AssetsCaller) SetPrevPlayerTeamId(opts *bind.CallOpts, state *big.Int, value *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "setPrevPlayerTeamId", state, value)
	return *ret0, err
}

// SetPrevPlayerTeamId is a free data retrieval call binding the contract method 0x37a86302.
//
// Solidity: function setPrevPlayerTeamId(uint256 state, uint256 value) constant returns(uint256)
func (_Assets *AssetsSession) SetPrevPlayerTeamId(state *big.Int, value *big.Int) (*big.Int, error) {
	return _Assets.Contract.SetPrevPlayerTeamId(&_Assets.CallOpts, state, value)
}

// SetPrevPlayerTeamId is a free data retrieval call binding the contract method 0x37a86302.
//
// Solidity: function setPrevPlayerTeamId(uint256 state, uint256 value) constant returns(uint256)
func (_Assets *AssetsCallerSession) SetPrevPlayerTeamId(state *big.Int, value *big.Int) (*big.Int, error) {
	return _Assets.Contract.SetPrevPlayerTeamId(&_Assets.CallOpts, state, value)
}

// TeamExists is a free data retrieval call binding the contract method 0x98981756.
//
// Solidity: function teamExists(uint256 teamId) constant returns(bool)
func (_Assets *AssetsCaller) TeamExists(opts *bind.CallOpts, teamId *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Assets.contract.Call(opts, out, "teamExists", teamId)
	return *ret0, err
}

// TeamExists is a free data retrieval call binding the contract method 0x98981756.
//
// Solidity: function teamExists(uint256 teamId) constant returns(bool)
func (_Assets *AssetsSession) TeamExists(teamId *big.Int) (bool, error) {
	return _Assets.Contract.TeamExists(&_Assets.CallOpts, teamId)
}

// TeamExists is a free data retrieval call binding the contract method 0x98981756.
//
// Solidity: function teamExists(uint256 teamId) constant returns(bool)
func (_Assets *AssetsCallerSession) TeamExists(teamId *big.Int) (bool, error) {
	return _Assets.Contract.TeamExists(&_Assets.CallOpts, teamId)
}

// Init is a paid mutator transaction binding the contract method 0xe1c7392a.
//
// Solidity: function init() returns()
func (_Assets *AssetsTransactor) Init(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Assets.contract.Transact(opts, "init")
}

// Init is a paid mutator transaction binding the contract method 0xe1c7392a.
//
// Solidity: function init() returns()
func (_Assets *AssetsSession) Init() (*types.Transaction, error) {
	return _Assets.Contract.Init(&_Assets.TransactOpts)
}

// Init is a paid mutator transaction binding the contract method 0xe1c7392a.
//
// Solidity: function init() returns()
func (_Assets *AssetsTransactorSession) Init() (*types.Transaction, error) {
	return _Assets.Contract.Init(&_Assets.TransactOpts)
}

// InitSingleTZ is a paid mutator transaction binding the contract method 0xa3ceb703.
//
// Solidity: function initSingleTZ(uint8 tz) returns()
func (_Assets *AssetsTransactor) InitSingleTZ(opts *bind.TransactOpts, tz uint8) (*types.Transaction, error) {
	return _Assets.contract.Transact(opts, "initSingleTZ", tz)
}

// InitSingleTZ is a paid mutator transaction binding the contract method 0xa3ceb703.
//
// Solidity: function initSingleTZ(uint8 tz) returns()
func (_Assets *AssetsSession) InitSingleTZ(tz uint8) (*types.Transaction, error) {
	return _Assets.Contract.InitSingleTZ(&_Assets.TransactOpts, tz)
}

// InitSingleTZ is a paid mutator transaction binding the contract method 0xa3ceb703.
//
// Solidity: function initSingleTZ(uint8 tz) returns()
func (_Assets *AssetsTransactorSession) InitSingleTZ(tz uint8) (*types.Transaction, error) {
	return _Assets.Contract.InitSingleTZ(&_Assets.TransactOpts, tz)
}

// SetActionsRoot is a paid mutator transaction binding the contract method 0xdba6319e.
//
// Solidity: function setActionsRoot(uint8 timeZone, bytes32 root) returns(uint256)
func (_Assets *AssetsTransactor) SetActionsRoot(opts *bind.TransactOpts, timeZone uint8, root [32]byte) (*types.Transaction, error) {
	return _Assets.contract.Transact(opts, "setActionsRoot", timeZone, root)
}

// SetActionsRoot is a paid mutator transaction binding the contract method 0xdba6319e.
//
// Solidity: function setActionsRoot(uint8 timeZone, bytes32 root) returns(uint256)
func (_Assets *AssetsSession) SetActionsRoot(timeZone uint8, root [32]byte) (*types.Transaction, error) {
	return _Assets.Contract.SetActionsRoot(&_Assets.TransactOpts, timeZone, root)
}

// SetActionsRoot is a paid mutator transaction binding the contract method 0xdba6319e.
//
// Solidity: function setActionsRoot(uint8 timeZone, bytes32 root) returns(uint256)
func (_Assets *AssetsTransactorSession) SetActionsRoot(timeZone uint8, root [32]byte) (*types.Transaction, error) {
	return _Assets.Contract.SetActionsRoot(&_Assets.TransactOpts, timeZone, root)
}

// SetSkillsRoot is a paid mutator transaction binding the contract method 0xec1c5423.
//
// Solidity: function setSkillsRoot(uint8 tz, bytes32 root) returns(uint256)
func (_Assets *AssetsTransactor) SetSkillsRoot(opts *bind.TransactOpts, tz uint8, root [32]byte) (*types.Transaction, error) {
	return _Assets.contract.Transact(opts, "setSkillsRoot", tz, root)
}

// SetSkillsRoot is a paid mutator transaction binding the contract method 0xec1c5423.
//
// Solidity: function setSkillsRoot(uint8 tz, bytes32 root) returns(uint256)
func (_Assets *AssetsSession) SetSkillsRoot(tz uint8, root [32]byte) (*types.Transaction, error) {
	return _Assets.Contract.SetSkillsRoot(&_Assets.TransactOpts, tz, root)
}

// SetSkillsRoot is a paid mutator transaction binding the contract method 0xec1c5423.
//
// Solidity: function setSkillsRoot(uint8 tz, bytes32 root) returns(uint256)
func (_Assets *AssetsTransactorSession) SetSkillsRoot(tz uint8, root [32]byte) (*types.Transaction, error) {
	return _Assets.Contract.SetSkillsRoot(&_Assets.TransactOpts, tz, root)
}

// TransferFirstBotToAddr is a paid mutator transaction binding the contract method 0x3c2eb360.
//
// Solidity: function transferFirstBotToAddr(uint8 timeZone, uint256 countryIdxInTZ, address addr) returns()
func (_Assets *AssetsTransactor) TransferFirstBotToAddr(opts *bind.TransactOpts, timeZone uint8, countryIdxInTZ *big.Int, addr common.Address) (*types.Transaction, error) {
	return _Assets.contract.Transact(opts, "transferFirstBotToAddr", timeZone, countryIdxInTZ, addr)
}

// TransferFirstBotToAddr is a paid mutator transaction binding the contract method 0x3c2eb360.
//
// Solidity: function transferFirstBotToAddr(uint8 timeZone, uint256 countryIdxInTZ, address addr) returns()
func (_Assets *AssetsSession) TransferFirstBotToAddr(timeZone uint8, countryIdxInTZ *big.Int, addr common.Address) (*types.Transaction, error) {
	return _Assets.Contract.TransferFirstBotToAddr(&_Assets.TransactOpts, timeZone, countryIdxInTZ, addr)
}

// TransferFirstBotToAddr is a paid mutator transaction binding the contract method 0x3c2eb360.
//
// Solidity: function transferFirstBotToAddr(uint8 timeZone, uint256 countryIdxInTZ, address addr) returns()
func (_Assets *AssetsTransactorSession) TransferFirstBotToAddr(timeZone uint8, countryIdxInTZ *big.Int, addr common.Address) (*types.Transaction, error) {
	return _Assets.Contract.TransferFirstBotToAddr(&_Assets.TransactOpts, timeZone, countryIdxInTZ, addr)
}

// TransferPlayer is a paid mutator transaction binding the contract method 0x800257d5.
//
// Solidity: function transferPlayer(uint256 playerId, uint256 teamIdTarget) returns()
func (_Assets *AssetsTransactor) TransferPlayer(opts *bind.TransactOpts, playerId *big.Int, teamIdTarget *big.Int) (*types.Transaction, error) {
	return _Assets.contract.Transact(opts, "transferPlayer", playerId, teamIdTarget)
}

// TransferPlayer is a paid mutator transaction binding the contract method 0x800257d5.
//
// Solidity: function transferPlayer(uint256 playerId, uint256 teamIdTarget) returns()
func (_Assets *AssetsSession) TransferPlayer(playerId *big.Int, teamIdTarget *big.Int) (*types.Transaction, error) {
	return _Assets.Contract.TransferPlayer(&_Assets.TransactOpts, playerId, teamIdTarget)
}

// TransferPlayer is a paid mutator transaction binding the contract method 0x800257d5.
//
// Solidity: function transferPlayer(uint256 playerId, uint256 teamIdTarget) returns()
func (_Assets *AssetsTransactorSession) TransferPlayer(playerId *big.Int, teamIdTarget *big.Int) (*types.Transaction, error) {
	return _Assets.Contract.TransferPlayer(&_Assets.TransactOpts, playerId, teamIdTarget)
}

// TransferTeam is a paid mutator transaction binding the contract method 0xe945e96a.
//
// Solidity: function transferTeam(uint256 teamId, address addr) returns()
func (_Assets *AssetsTransactor) TransferTeam(opts *bind.TransactOpts, teamId *big.Int, addr common.Address) (*types.Transaction, error) {
	return _Assets.contract.Transact(opts, "transferTeam", teamId, addr)
}

// TransferTeam is a paid mutator transaction binding the contract method 0xe945e96a.
//
// Solidity: function transferTeam(uint256 teamId, address addr) returns()
func (_Assets *AssetsSession) TransferTeam(teamId *big.Int, addr common.Address) (*types.Transaction, error) {
	return _Assets.Contract.TransferTeam(&_Assets.TransactOpts, teamId, addr)
}

// TransferTeam is a paid mutator transaction binding the contract method 0xe945e96a.
//
// Solidity: function transferTeam(uint256 teamId, address addr) returns()
func (_Assets *AssetsTransactorSession) TransferTeam(teamId *big.Int, addr common.Address) (*types.Transaction, error) {
	return _Assets.Contract.TransferTeam(&_Assets.TransactOpts, teamId, addr)
}

// AssetsDivisionCreationIterator is returned from FilterDivisionCreation and is used to iterate over the raw logs and unpacked data for DivisionCreation events raised by the Assets contract.
type AssetsDivisionCreationIterator struct {
	Event *AssetsDivisionCreation // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AssetsDivisionCreationIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AssetsDivisionCreation)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AssetsDivisionCreation)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AssetsDivisionCreationIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AssetsDivisionCreationIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AssetsDivisionCreation represents a DivisionCreation event raised by the Assets contract.
type AssetsDivisionCreation struct {
	Timezone             uint8
	CountryIdxInTZ       *big.Int
	DivisionIdxInCountry *big.Int
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterDivisionCreation is a free log retrieval operation binding the contract event 0xc5d195855a200aa90e2052bcc795cedbc84c2a26556b1d5113b5a30c96003a0b.
//
// Solidity: event DivisionCreation(uint8 timezone, uint256 countryIdxInTZ, uint256 divisionIdxInCountry)
func (_Assets *AssetsFilterer) FilterDivisionCreation(opts *bind.FilterOpts) (*AssetsDivisionCreationIterator, error) {

	logs, sub, err := _Assets.contract.FilterLogs(opts, "DivisionCreation")
	if err != nil {
		return nil, err
	}
	return &AssetsDivisionCreationIterator{contract: _Assets.contract, event: "DivisionCreation", logs: logs, sub: sub}, nil
}

// WatchDivisionCreation is a free log subscription operation binding the contract event 0xc5d195855a200aa90e2052bcc795cedbc84c2a26556b1d5113b5a30c96003a0b.
//
// Solidity: event DivisionCreation(uint8 timezone, uint256 countryIdxInTZ, uint256 divisionIdxInCountry)
func (_Assets *AssetsFilterer) WatchDivisionCreation(opts *bind.WatchOpts, sink chan<- *AssetsDivisionCreation) (event.Subscription, error) {

	logs, sub, err := _Assets.contract.WatchLogs(opts, "DivisionCreation")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AssetsDivisionCreation)
				if err := _Assets.contract.UnpackLog(event, "DivisionCreation", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// AssetsPlayerStateChangeIterator is returned from FilterPlayerStateChange and is used to iterate over the raw logs and unpacked data for PlayerStateChange events raised by the Assets contract.
type AssetsPlayerStateChangeIterator struct {
	Event *AssetsPlayerStateChange // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AssetsPlayerStateChangeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AssetsPlayerStateChange)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AssetsPlayerStateChange)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AssetsPlayerStateChangeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AssetsPlayerStateChangeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AssetsPlayerStateChange represents a PlayerStateChange event raised by the Assets contract.
type AssetsPlayerStateChange struct {
	PlayerId *big.Int
	State    *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterPlayerStateChange is a free log retrieval operation binding the contract event 0x65a4d4a8a0afb474d2e9465a6a1a41bb88fd04f41152ba070421f1b3771f15d8.
//
// Solidity: event PlayerStateChange(uint256 playerId, uint256 state)
func (_Assets *AssetsFilterer) FilterPlayerStateChange(opts *bind.FilterOpts) (*AssetsPlayerStateChangeIterator, error) {

	logs, sub, err := _Assets.contract.FilterLogs(opts, "PlayerStateChange")
	if err != nil {
		return nil, err
	}
	return &AssetsPlayerStateChangeIterator{contract: _Assets.contract, event: "PlayerStateChange", logs: logs, sub: sub}, nil
}

// WatchPlayerStateChange is a free log subscription operation binding the contract event 0x65a4d4a8a0afb474d2e9465a6a1a41bb88fd04f41152ba070421f1b3771f15d8.
//
// Solidity: event PlayerStateChange(uint256 playerId, uint256 state)
func (_Assets *AssetsFilterer) WatchPlayerStateChange(opts *bind.WatchOpts, sink chan<- *AssetsPlayerStateChange) (event.Subscription, error) {

	logs, sub, err := _Assets.contract.WatchLogs(opts, "PlayerStateChange")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AssetsPlayerStateChange)
				if err := _Assets.contract.UnpackLog(event, "PlayerStateChange", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// AssetsTeamTransferIterator is returned from FilterTeamTransfer and is used to iterate over the raw logs and unpacked data for TeamTransfer events raised by the Assets contract.
type AssetsTeamTransferIterator struct {
	Event *AssetsTeamTransfer // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AssetsTeamTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AssetsTeamTransfer)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AssetsTeamTransfer)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AssetsTeamTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AssetsTeamTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AssetsTeamTransfer represents a TeamTransfer event raised by the Assets contract.
type AssetsTeamTransfer struct {
	TeamId *big.Int
	To     common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterTeamTransfer is a free log retrieval operation binding the contract event 0x77b66eb1e6d2bc131b79be4213ae7f08f29351c01060e10bcc0302278067bf38.
//
// Solidity: event TeamTransfer(uint256 teamId, address to)
func (_Assets *AssetsFilterer) FilterTeamTransfer(opts *bind.FilterOpts) (*AssetsTeamTransferIterator, error) {

	logs, sub, err := _Assets.contract.FilterLogs(opts, "TeamTransfer")
	if err != nil {
		return nil, err
	}
	return &AssetsTeamTransferIterator{contract: _Assets.contract, event: "TeamTransfer", logs: logs, sub: sub}, nil
}

// WatchTeamTransfer is a free log subscription operation binding the contract event 0x77b66eb1e6d2bc131b79be4213ae7f08f29351c01060e10bcc0302278067bf38.
//
// Solidity: event TeamTransfer(uint256 teamId, address to)
func (_Assets *AssetsFilterer) WatchTeamTransfer(opts *bind.WatchOpts, sink chan<- *AssetsTeamTransfer) (event.Subscription, error) {

	logs, sub, err := _Assets.contract.WatchLogs(opts, "TeamTransfer")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AssetsTeamTransfer)
				if err := _Assets.contract.UnpackLog(event, "TeamTransfer", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}
