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
const AssetsABI = "[{\"inputs\":[],\"constant\":true,\"name\":\"IDX_MD\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"encodedSkills\"}],\"constant\":true,\"name\":\"getSkills\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"encodedSkills\"}],\"constant\":true,\"name\":\"getSumOfSkills\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"IDX_R\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"SK_END\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"encodedSkills\"}],\"constant\":true,\"name\":\"getAggressiveness\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint8\",\"name\":\"timeZone\"},{\"type\":\"uint256\",\"name\":\"countryIdxInTZ\"},{\"type\":\"uint256\",\"name\":\"val\"}],\"constant\":true,\"name\":\"encodeTZCountryAndVal\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"DAYS_PER_ROUND\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"encodedSkills\"}],\"constant\":true,\"name\":\"getEndurance\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"MAX_PLAYER_AGE_AT_BIRTH\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"encoded\"}],\"constant\":true,\"name\":\"decodeTZCountryAndVal\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"},{\"type\":\"uint256\",\"name\":\"\"},{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"encodedSkills\"}],\"constant\":true,\"name\":\"getLeftishness\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"IDX_D\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"state\"},{\"type\":\"uint256\",\"name\":\"value\"}],\"constant\":true,\"name\":\"setPrevPlayerTeamId\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"IDX_LC\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"FREEVERSE\",\"outputs\":[{\"type\":\"address\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"IDX_LR\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"SK_SHO\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"encodedSkills\"}],\"constant\":true,\"name\":\"getPotential\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"LEAGUES_PER_DIV\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"encodedSkills\"}],\"constant\":true,\"name\":\"getSpeed\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"playerState\"}],\"constant\":true,\"name\":\"getPrevPlayerTeamId\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"state\"},{\"type\":\"uint8\",\"name\":\"currentShirtNum\"}],\"constant\":true,\"name\":\"setCurrentShirtNum\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"encodedSkills\"}],\"constant\":true,\"name\":\"getDefence\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"encodedSkills\"}],\"constant\":true,\"name\":\"getPass\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"IDX_CR\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"encodedSkills\"}],\"constant\":true,\"name\":\"getShoot\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"encodedSkills\"}],\"constant\":true,\"name\":\"getAlignedLastHalf\",\"outputs\":[{\"type\":\"bool\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"encodedSkills\"}],\"constant\":true,\"name\":\"getPlayerIdFromSkills\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"IDX_GK\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"playerState\"}],\"constant\":true,\"name\":\"getPlayerIdFromState\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"encodedSkills\"}],\"constant\":true,\"name\":\"getInjuryWeeksLeft\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"PLAYERS_PER_TEAM_INIT\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"gameDeployMonth\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"encodedSkills\"}],\"constant\":true,\"name\":\"getMonthOfBirth\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"currentRound\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"PLAYERS_PER_TEAM_MAX\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"IDX_MF\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"N_SKILLS\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint16[5]\",\"name\":\"skills\"},{\"type\":\"uint256\",\"name\":\"monthOfBirth\"},{\"type\":\"uint256\",\"name\":\"playerId\"},{\"type\":\"uint8[4]\",\"name\":\"birthTraits\"},{\"type\":\"bool\",\"name\":\"alignedLastHalf\"},{\"type\":\"bool\",\"name\":\"redCardLastGame\"},{\"type\":\"uint8\",\"name\":\"gamesNonStopping\"},{\"type\":\"uint8\",\"name\":\"injuryWeeksLeft\"}],\"constant\":true,\"name\":\"encodePlayerSkills\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"encoded\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"IDX_M\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"playerId\"},{\"type\":\"uint256\",\"name\":\"currentTeamId\"},{\"type\":\"uint8\",\"name\":\"currentShirtNum\"},{\"type\":\"uint256\",\"name\":\"prevPlayerTeamId\"},{\"type\":\"uint256\",\"name\":\"lastSaleBlock\"}],\"constant\":true,\"name\":\"encodePlayerState\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"SK_PAS\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"TEAMS_PER_LEAGUE\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"state\"},{\"type\":\"uint256\",\"name\":\"lastSaleBlock\"}],\"constant\":true,\"name\":\"setLastSaleBlock\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"NULL_ADDR\",\"outputs\":[{\"type\":\"address\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"IDX_LCR\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"constant\":true,\"name\":\"_timeZones\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"nCountriesToAdd\"},{\"type\":\"uint8\",\"name\":\"newestOrgMapIdx\"},{\"type\":\"uint8\",\"name\":\"newestSkillsIdx\"},{\"type\":\"bytes32\",\"name\":\"scoresRoot\"},{\"type\":\"uint8\",\"name\":\"updateCycleIdx\"},{\"type\":\"uint256\",\"name\":\"lastActionsSubmissionTime\"},{\"type\":\"uint256\",\"name\":\"lastUpdateTime\"},{\"type\":\"bytes32\",\"name\":\"actionsRoot\"},{\"type\":\"uint256\",\"name\":\"lastMarketClosureBlockNum\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"FREE_PLAYER_ID\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"encodedSkills\"}],\"constant\":true,\"name\":\"getForwardness\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"playerState\"},{\"type\":\"uint256\",\"name\":\"teamId\"}],\"constant\":true,\"name\":\"setCurrentTeamId\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"playerState\"}],\"constant\":true,\"name\":\"getLastSaleBlock\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"MIN_PLAYER_AGE_AT_BIRTH\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"encodedSkills\"}],\"constant\":true,\"name\":\"getSkillsVec\",\"outputs\":[{\"type\":\"uint16[5]\",\"name\":\"skills\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"encodedSkills\"}],\"constant\":true,\"name\":\"getRedCardLastGame\",\"outputs\":[{\"type\":\"bool\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"playerState\"}],\"constant\":true,\"name\":\"getCurrentTeamId\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"IDX_F\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"tactics\"}],\"constant\":true,\"name\":\"decodeTactics\",\"outputs\":[{\"type\":\"uint8[11]\",\"name\":\"lineup\"},{\"type\":\"bool[10]\",\"name\":\"extraAttack\"},{\"type\":\"uint8\",\"name\":\"tacticsId\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"encodedSkills\"}],\"constant\":true,\"name\":\"getGamesNonStopping\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"SK_DEF\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint8[11]\",\"name\":\"lineup\"},{\"type\":\"bool[10]\",\"name\":\"extraAttack\"},{\"type\":\"uint8\",\"name\":\"tacticsId\"}],\"constant\":true,\"name\":\"encodeTactics\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"playerState\"}],\"constant\":true,\"name\":\"getCurrentShirtNum\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"IDX_L\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"TEAMS_PER_DIVISION\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"IDX_C\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"SK_SPE\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"indexed\":false,\"type\":\"uint256\",\"name\":\"teamId\"},{\"indexed\":false,\"type\":\"address\",\"name\":\"to\"}],\"type\":\"event\",\"name\":\"TeamTransfer\",\"anonymous\":false},{\"inputs\":[{\"indexed\":false,\"type\":\"uint256\",\"name\":\"playerId\"},{\"indexed\":false,\"type\":\"uint256\",\"name\":\"teamIdTarget\"}],\"type\":\"event\",\"name\":\"PlayerTransfer\",\"anonymous\":false},{\"inputs\":[{\"indexed\":false,\"type\":\"uint8\",\"name\":\"timezone\"},{\"indexed\":false,\"type\":\"uint256\",\"name\":\"countryIdxInTZ\"},{\"indexed\":false,\"type\":\"uint256\",\"name\":\"divisionIdxInCountry\"}],\"type\":\"event\",\"name\":\"DivisionCreation\",\"anonymous\":false},{\"inputs\":[{\"indexed\":false,\"type\":\"uint256\",\"name\":\"playerId\"},{\"indexed\":false,\"type\":\"uint256\",\"name\":\"state\"}],\"type\":\"event\",\"name\":\"PlayerStateChange\",\"anonymous\":false},{\"inputs\":[],\"constant\":false,\"name\":\"init\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint8\",\"name\":\"tz\"}],\"constant\":false,\"name\":\"initSingleTZ\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint8\",\"name\":\"timeZone\"}],\"constant\":true,\"name\":\"getLastUpdateTime\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint8\",\"name\":\"timeZone\"}],\"constant\":true,\"name\":\"getLastActionsSubmissionTime\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint8\",\"name\":\"tz\"},{\"type\":\"bytes32\",\"name\":\"root\"}],\"constant\":false,\"name\":\"setSkillsRoot\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"nonpayable\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint8\",\"name\":\"timeZone\"},{\"type\":\"bytes32\",\"name\":\"root\"}],\"constant\":false,\"name\":\"setActionsRoot\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"nonpayable\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint8\",\"name\":\"timeZone\"}],\"constant\":true,\"name\":\"getNCountriesInTZ\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint8\",\"name\":\"timeZone\"},{\"type\":\"uint256\",\"name\":\"countryIdxInTZ\"}],\"constant\":true,\"name\":\"getNDivisionsInCountry\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint8\",\"name\":\"timeZone\"},{\"type\":\"uint256\",\"name\":\"countryIdxInTZ\"}],\"constant\":true,\"name\":\"getNLeaguesInCountry\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint8\",\"name\":\"timeZone\"},{\"type\":\"uint256\",\"name\":\"countryIdxInTZ\"}],\"constant\":true,\"name\":\"getNTeamsInCountry\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint8\",\"name\":\"timeZone\"},{\"type\":\"uint256\",\"name\":\"countryIdxInTZ\"},{\"type\":\"uint256\",\"name\":\"teamIdxInCountry\"}],\"constant\":true,\"name\":\"_teamExistsInCountry\",\"outputs\":[{\"type\":\"bool\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"teamId\"}],\"constant\":true,\"name\":\"teamExists\",\"outputs\":[{\"type\":\"bool\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint8\",\"name\":\"timeZone\"},{\"type\":\"uint256\",\"name\":\"countryIdxInTZ\"},{\"type\":\"uint256\",\"name\":\"teamIdxInCountry\"}],\"constant\":true,\"name\":\"isBotTeamInCountry\",\"outputs\":[{\"type\":\"bool\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"teamId\"}],\"constant\":true,\"name\":\"isBotTeam\",\"outputs\":[{\"type\":\"bool\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint8\",\"name\":\"timeZone\"},{\"type\":\"uint256\",\"name\":\"countryIdxInTZ\"},{\"type\":\"uint256\",\"name\":\"teamIdxInCountry\"}],\"constant\":true,\"name\":\"getOwnerTeamInCountry\",\"outputs\":[{\"type\":\"address\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"teamId\"}],\"constant\":true,\"name\":\"getOwnerTeam\",\"outputs\":[{\"type\":\"address\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"playerId\"}],\"constant\":true,\"name\":\"getOwnerPlayer\",\"outputs\":[{\"type\":\"address\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"playerId\"}],\"constant\":true,\"name\":\"playerExists\",\"outputs\":[{\"type\":\"bool\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"playerId\"}],\"constant\":true,\"name\":\"isVirtualPlayer\",\"outputs\":[{\"type\":\"bool\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint8\",\"name\":\"timeZone\"},{\"type\":\"uint256\",\"name\":\"countryIdxInTZ\"},{\"type\":\"address\",\"name\":\"addr\"}],\"constant\":false,\"name\":\"transferFirstBotToAddr\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"teamId\"},{\"type\":\"address\",\"name\":\"addr\"}],\"constant\":false,\"name\":\"transferTeam\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint8\",\"name\":\"timeZone\"},{\"type\":\"uint256\",\"name\":\"countryIdxInTZ\"},{\"type\":\"uint256\",\"name\":\"teamIdxInCountry\"},{\"type\":\"uint8\",\"name\":\"shirtNum\"}],\"constant\":true,\"name\":\"getDefaultPlayerIdForTeamInCountry\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"teamId\"}],\"constant\":true,\"name\":\"getPlayerIdsInTeam\",\"outputs\":[{\"type\":\"uint256[25]\",\"name\":\"playerIds\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"playerId\"}],\"constant\":true,\"name\":\"getPlayerSkillsAtBirth\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"playerId\"}],\"constant\":true,\"name\":\"getPlayerStateAtBirth\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"playerId\"}],\"constant\":true,\"name\":\"getPlayerState\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"dna\"},{\"type\":\"uint256\",\"name\":\"playerCreationMonth\"}],\"constant\":true,\"name\":\"computeBirthMonth\",\"outputs\":[{\"type\":\"uint16\",\"name\":\"\"},{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"dna\"},{\"type\":\"uint8\",\"name\":\"shirtNum\"}],\"constant\":true,\"name\":\"computeSkills\",\"outputs\":[{\"type\":\"uint16[5]\",\"name\":\"\"},{\"type\":\"uint8[4]\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"teamId\"},{\"type\":\"uint8\",\"name\":\"shirtNum\"}],\"constant\":true,\"name\":\"isFreeShirt\",\"outputs\":[{\"type\":\"bool\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"playerId\"}],\"constant\":true,\"name\":\"getPlayerAgeInMonths\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"teamId\"}],\"constant\":true,\"name\":\"getFreeShirt\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"playerId\"},{\"type\":\"uint256\",\"name\":\"teamIdTarget\"}],\"constant\":false,\"name\":\"transferPlayer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint8\",\"name\":\"timeZone\"}],\"constant\":true,\"name\":\"countCountries\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint8\",\"name\":\"timeZone\"},{\"type\":\"uint256\",\"name\":\"countryIdxInTZ\"}],\"constant\":true,\"name\":\"countTeams\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"}]"

// AssetsBin is the compiled bytecode used for deploying new contracts.
const AssetsBin = `0x6080604052600161014960006101000a81548160ff02191690831515021790555034801561002c57600080fd5b506155718061003c6000396000f3fe608060405234801561001057600080fd5b50600436106105995760003560e01c806387f1e880116102e5578063c37b1c251161018d578063e81e21bb116100f4578063ec71bc82116100ad578063f305a21c11610087578063f305a21c14611f97578063f8ef7b9e14611fbb578063f9d0723d14611fdf578063fa80039b1461202e57610599565b8063ec71bc8214611f0d578063ec7ecec514611f31578063f21f5a8314611f7357610599565b8063e81e21bb14611cd6578063e945e96a14611cfa578063e9e7165214611d48578063eabf6a4b14611e12578063eb78b7b714611e7c578063ec1c542314611ebe57610599565b8063cd2105e811610146578063cd2105e814611b2f578063d7b63a1114611b71578063dba6319e14611b95578063e1c7392a14611be4578063e6400ac414611bee578063e804e51914611c9457610599565b8063c37b1c251461198b578063c566b5bc146119d7578063c73f808d14611a19578063c79055d414611a5b578063cc1cc3d714611a7f578063cc7d473b14611ae957610599565b8063a3ceb7031161024c578063b3f390b311610205578063bc1a97c1116101df578063bc1a97c114611896578063c04f6d53146118dc578063c258012b1461192b578063c2bc41cd1461194957610599565b8063b3f390b314611796578063b9627097146117e0578063b96b1a301461180457610599565b8063a3ceb70314611646578063ab1b7c5e14611677578063ac5db9ee1461169b578063ad63bcbd146116bf578063af76cd0114611704578063b32aa2c11461175057610599565b8063963fcc801161029e578063963fcc80146113f7578063976daaac1461144a578063989817561461146e5780639c53e3fd146114b45780639cc62340146115b55780639f27112a146115d957610599565b806387f1e8801461129b5780638a19c8bc146112dd5780638adddc9d146112fb5780638cc9a8d51461131f5780638f3db436146113655780638f9da2141461138957610599565b806340cd05fd116104485780635adb40f5116103af57806378f4c71811610368578063800257d511610342578063800257d5146111c457806380bac709146111fc57806383c31d3b14611259578063859824311461127d57610599565b806378f4c718146110f157806379e76597146111335780637b2566a51461117557610599565b80635adb40f514610f905780635becd99914610fdf57806365b4b47614611003578063673fe242146110455780636f6c2ae01461108b5780637420a606146110cd57610599565b80634db989fd116104015780634db989fd14610d4a578063507b172314610d9957806351585b4914610de1578063547d829814610e2357806355a6f86f14610ec9578063595ef25b14610f0b57610599565b806340cd05fd14610bce578063434807ad14610bf257806348d1e9c014610c34578063492afc6914610c585780634b93f75314610cc65780634bea2a6914610d0857610599565b8063228408b0116105075780633518dd1d116104c057806337fd56af1161049a57806337fd56af14610ae157806339644f2114610b055780633c2eb36014610b4f5780633d085f9614610baa57610599565b80633518dd1d14610a2f578063369151db14610a7157806337a8630214610a9557610599565b8063228408b014610886578063258e5d90146108ec578063266576081461092e5780632a238b0a146109705780632d0e08fd146109945780633260840b146109d957610599565b80631884332c116105595780631884332c1461070a5780631a6daba21461072e5780631fc7768f1461078b5780631ffeb349146107cd57806320748ae81461080f57806321ff8ae81461086857610599565b80623e32231461059e57806292bf78146105c2578062aae8df146106045780630abcd3e51461065f5780631060c9c2146106a4578063169d2914146106e6575b600080fd5b6105a6612073565b604051808260ff1660ff16815260200191505060405180910390f35b6105ee600480360360208110156105d857600080fd5b8101908080359060200190929190505050612078565b6040518082815260200191505060405180910390f35b61063a6004803603604081101561061a57600080fd5b810190808035906020019092919080359060200190929190505050612086565b604051808361ffff1661ffff1681526020018281526020019250505060405180910390f35b61068e6004803603602081101561067557600080fd5b81019080803560ff169060200190929190505050612130565b6040518082815260200191505060405180910390f35b6106d0600480360360208110156106ba57600080fd5b810190808035906020019092919050505061215d565b6040518082815260200191505060405180910390f35b6106ee612197565b604051808260ff1660ff16815260200191505060405180910390f35b61071261219c565b604051808260ff1660ff16815260200191505060405180910390f35b6107716004803603606081101561074457600080fd5b81019080803560ff16906020019092919080359060200190929190803590602001909291905050506121a1565b604051808215151515815260200191505060405180910390f35b6107b7600480360360208110156107a157600080fd5b81019080803590602001909291905050506121b8565b6040518082815260200191505060405180910390f35b6107f9600480360360208110156107e357600080fd5b81019080803590602001909291905050506121c9565b6040518082815260200191505060405180910390f35b6108526004803603606081101561082557600080fd5b81019080803560ff16906020019092919080359060200190929190803590602001909291905050506121f5565b6040518082815260200191505060405180910390f35b610870612382565b6040518082815260200191505060405180910390f35b6108d66004803603608081101561089c57600080fd5b81019080803560ff1690602001909291908035906020019092919080359060200190929190803560ff169060200190929190505050612387565b6040518082815260200191505060405180910390f35b6109186004803603602081101561090257600080fd5b81019080803590602001909291905050506123c1565b6040518082815260200191505060405180910390f35b61095a6004803603602081101561094457600080fd5b81019080803590602001909291905050506123d3565b6040518082815260200191505060405180910390f35b6109786124b3565b604051808260ff1660ff16815260200191505060405180910390f35b6109c3600480360360208110156109aa57600080fd5b81019080803560ff1690602001909291905050506124b8565b6040518082815260200191505060405180910390f35b610a05600480360360208110156109ef57600080fd5b81019080803590602001909291905050506124e2565b604051808460ff1660ff168152602001838152602001828152602001935050505060405180910390f35b610a5b60048036036020811015610a4557600080fd5b810190808035906020019092919050505061250c565b6040518082815260200191505060405180910390f35b610a7961251d565b604051808260ff1660ff16815260200191505060405180910390f35b610acb60048036036040811015610aab57600080fd5b810190808035906020019092919080359060200190929190505050612522565b6040518082815260200191505060405180910390f35b610ae96125cc565b604051808260ff1660ff16815260200191505060405180910390f35b610b0d6125d1565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b610ba860048036036060811015610b6557600080fd5b81019080803560ff16906020019092919080359060200190929190803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506125d6565b005b610bb2612907565b604051808260ff1660ff16815260200191505060405180910390f35b610bd661290c565b604051808260ff1660ff16815260200191505060405180910390f35b610c1e60048036036020811015610c0857600080fd5b8101908080359060200190929190505050612911565b6040518082815260200191505060405180910390f35b610c3c612922565b604051808260ff1660ff16815260200191505060405180910390f35b610c8460048036036020811015610c6e57600080fd5b8101908080359060200190929190505050612927565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b610cf260048036036020811015610cdc57600080fd5b8101908080359060200190929190505050612951565b6040518082815260200191505060405180910390f35b610d3460048036036020811015610d1e57600080fd5b8101908080359060200190929190505050612963565b6040518082815260200191505060405180910390f35b610d8360048036036040811015610d6057600080fd5b8101908080359060200190929190803560ff169060200190929190505050612979565b6040518082815260200191505060405180910390f35b610dc560048036036020811015610daf57600080fd5b8101908080359060200190929190505050612a25565b604051808260ff1660ff16815260200191505060405180910390f35b610e0d60048036036020811015610df757600080fd5b8101908080359060200190929190505050612a6c565b6040518082815260200191505060405180910390f35b610e5c60048036036040811015610e3957600080fd5b8101908080359060200190929190803560ff169060200190929190505050612a7e565b6040518083600560200280838360005b83811015610e87578082015181840152602081019050610e6c565b5050505090500182600460200280838360005b83811015610eb5578082015181840152602081019050610e9a565b505050509050019250505060405180910390f35b610ef560048036036020811015610edf57600080fd5b8101908080359060200190929190505050612f30565b6040518082815260200191505060405180910390f35b610f4e60048036036060811015610f2157600080fd5b81019080803560ff1690602001909291908035906020019092919080359060200190929190505050612f42565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b610fc960048036036040811015610fa657600080fd5b81019080803560ff16906020019092919080359060200190929190505050612fc8565b6040518082815260200191505060405180910390f35b610fe7613019565b604051808260ff1660ff16815260200191505060405180910390f35b61102f6004803603602081101561101957600080fd5b810190808035906020019092919050505061301e565b6040518082815260200191505060405180910390f35b6110716004803603602081101561105b57600080fd5b8101908080359060200190929190505050613030565b604051808215151515815260200191505060405180910390f35b6110b7600480360360208110156110a157600080fd5b8101908080359060200190929190505050613043565b6040518082815260200191505060405180910390f35b6110d5613059565b604051808260ff1660ff16815260200191505060405180910390f35b61111d6004803603602081101561110757600080fd5b810190808035906020019092919050505061305e565b6040518082815260200191505060405180910390f35b61115f6004803603602081101561114957600080fd5b8101908080359060200190929190505050613074565b6040518082815260200191505060405180910390f35b6111ae6004803603604081101561118b57600080fd5b81019080803560ff16906020019092919080359060200190929190505050613085565b6040518082815260200191505060405180910390f35b6111fa600480360360408110156111da57600080fd5b8101908080359060200190929190803590602001909291905050506130dc565b005b61123f6004803603606081101561121257600080fd5b81019080803560ff1690602001909291908035906020019092919080359060200190929190505050613489565b604051808215151515815260200191505060405180910390f35b6112616134cd565b604051808260ff1660ff16815260200191505060405180910390f35b6112856134d2565b6040518082815260200191505060405180910390f35b6112c7600480360360208110156112b157600080fd5b81019080803590602001909291905050506134d9565b6040518082815260200191505060405180910390f35b6112e56134eb565b6040518082815260200191505060405180910390f35b6113036134f2565b604051808260ff1660ff16815260200191505060405180910390f35b61134b6004803603602081101561133557600080fd5b81019080803590602001909291905050506134f7565b604051808215151515815260200191505060405180910390f35b61136d613521565b604051808260ff1660ff16815260200191505060405180910390f35b6113b56004803603602081101561139f57600080fd5b8101908080359060200190929190505050613526565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6114306004803603604081101561140d57600080fd5b8101908080359060200190929190803560ff1690602001909291905050506135c9565b604051808215151515815260200191505060405180910390f35b6114526136d7565b604051808260ff1660ff16815260200191505060405180910390f35b61149a6004803603602081101561148457600080fd5b81019080803590602001909291905050506136dc565b604051808215151515815260200191505060405180910390f35b61159f60048036036101e08110156114cb57600080fd5b810190808060a001906005806020026040519081016040528092919082600560200280828437600081840152601f19601f8201169050808301925050505050509192919290803590602001909291908035906020019092919080608001906004806020026040519081016040528092919082600460200280828437600081840152601f19601f8201169050808301925050505050509192919290803515159060200190929190803515159060200190929190803560ff169060200190929190803560ff169060200190929190505050613706565b6040518082815260200191505060405180910390f35b6115bd613d15565b604051808260ff1660ff16815260200191505060405180910390f35b611630600480360360a08110156115ef57600080fd5b810190808035906020019092919080359060200190929190803560ff1690602001909291908035906020019092919080359060200190929190505050613d1a565b6040518082815260200191505060405180910390f35b6116756004803603602081101561165c57600080fd5b81019080803560ff169060200190929190505050613de8565b005b61167f613eaa565b604051808260ff1660ff16815260200191505060405180910390f35b6116a3613eaf565b604051808260ff1660ff16815260200191505060405180910390f35b6116ee600480360360208110156116d557600080fd5b81019080803560ff169060200190929190505050613eb4565b6040518082815260200191505060405180910390f35b61173a6004803603604081101561171a57600080fd5b810190808035906020019092919080359060200190929190505050613ee1565b6040518082815260200191505060405180910390f35b61177c6004803603602081101561176657600080fd5b8101908080359060200190929190505050613f85565b604051808215151515815260200191505060405180910390f35b61179e61401f565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6117e8614024565b604051808260ff1660ff16815260200191505060405180910390f35b6118306004803603602081101561181a57600080fd5b8101908080359060200190929190505050614029565b604051808a60ff1660ff1681526020018960ff1660ff1681526020018860ff1660ff1681526020018781526020018660ff1660ff168152602001858152602001848152602001838152602001828152602001995050505050505050505060405180910390f35b6118c2600480360360208110156118ac57600080fd5b81019080803590602001909291905050506140ae565b604051808215151515815260200191505060405180910390f35b611915600480360360408110156118f257600080fd5b81019080803560ff1690602001909291908035906020019092919050505061410e565b6040518082815260200191505060405180910390f35b611933614128565b6040518082815260200191505060405180910390f35b6119756004803603602081101561195f57600080fd5b810190808035906020019092919050505061412d565b6040518082815260200191505060405180910390f35b6119c1600480360360408110156119a157600080fd5b81019080803590602001909291908035906020019092919050505061413e565b6040518082815260200191505060405180910390f35b611a03600480360360208110156119ed57600080fd5b81019080803590602001909291905050506141ee565b6040518082815260200191505060405180910390f35b611a4560048036036020811015611a2f57600080fd5b8101908080359060200190929190505050614203565b6040518082815260200191505060405180910390f35b611a63614396565b604051808260ff1660ff16815260200191505060405180910390f35b611aab60048036036020811015611a9557600080fd5b810190808035906020019092919050505061439b565b6040518082600560200280838360005b83811015611ad6578082015181840152602081019050611abb565b5050505090500191505060405180910390f35b611b1560048036036020811015611aff57600080fd5b8101908080359060200190929190505050614475565b604051808215151515815260200191505060405180910390f35b611b5b60048036036020811015611b4557600080fd5b8101908080359060200190929190505050614488565b6040518082815260200191505060405180910390f35b611b7961449e565b604051808260ff1660ff16815260200191505060405180910390f35b611bce60048036036040811015611bab57600080fd5b81019080803560ff169060200190929190803590602001909291905050506144a3565b6040518082815260200191505060405180910390f35b611bec6144ec565b005b611c1a60048036036020811015611c0457600080fd5b81019080803590602001909291905050506145ce565b6040518084600b60200280838360005b83811015611c45578082015181840152602081019050611c2a565b5050505090500183600a60200280838360005b83811015611c73578082015181840152602081019050611c58565b505050509050018260ff1660ff168152602001935050505060405180910390f35b611cc060048036036020811015611caa57600080fd5b81019080803590602001909291905050506147a7565b6040518082815260200191505060405180910390f35b611cde6147b8565b604051808260ff1660ff16815260200191505060405180910390f35b611d4660048036036040811015611d1057600080fd5b8101908080359060200190929190803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506147bd565b005b611dfc60048036036102c0811015611d5f57600080fd5b81019080806101600190600b806020026040519081016040528092919082600b60200280828437600081840152601f19601f8201169050808301925050505050509192919290806101400190600a806020026040519081016040528092919082600a60200280828437600081840152601f19601f8201169050808301925050505050509192919290803560ff16906020019092919050505061484f565b6040518082815260200191505060405180910390f35b611e3e60048036036020811015611e2857600080fd5b8101908080359060200190929190505050614a08565b6040518082601960200280838360005b83811015611e69578082015181840152602081019050611e4e565b5050505090500191505060405180910390f35b611ea860048036036020811015611e9257600080fd5b8101908080359060200190929190505050614bd1565b6040518082815260200191505060405180910390f35b611ef760048036036040811015611ed457600080fd5b81019080803560ff16906020019092919080359060200190929190505050614be2565b6040518082815260200191505060405180910390f35b611f15614cac565b604051808260ff1660ff16815260200191505060405180910390f35b611f5d60048036036020811015611f4757600080fd5b8101908080359060200190929190505050614cb1565b6040518082815260200191505060405180910390f35b611f7b614cec565b604051808260ff1660ff16815260200191505060405180910390f35b611f9f614cf1565b604051808260ff1660ff16815260200191505060405180910390f35b611fc3614cf6565b604051808260ff1660ff16815260200191505060405180910390f35b61201860048036036040811015611ff557600080fd5b81019080803560ff16906020019092919080359060200190929190505050614cfb565b6040518082815260200191505060405180910390f35b61205d6004803603602081101561204457600080fd5b81019080803560ff169060200190929190505050614d15565b6040518082815260200191505060405180910390f35b600481565b600060ba82901c9050919050565b6000806101e08311612100576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601b8152602001807f696e76616c696420706c617965724372656174696f6e4d6f6e7468000000000081525060200191505060405180910390fd5b60006014858161210c57fe5b066010019050600585901c9450600c810261ffff1684038592509250509250929050565b600061213b82614d3f565b60018260ff166019811061214b57fe5b600d0201600001805490509050919050565b6000612168826123c1565b61217183612a6c565b61217a84612f30565b61218385612951565b61218c8661301e565b010101019050919050565b600181565b600481565b60006121ad848461410e565b821090509392505050565b60006007606c83901c169050919050565b60006121ee6121e76121e26121dd85614203565b6134d9565b614dca565b4203614de4565b9050919050565b600060208460ff1610612270576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260148152602001807f646566656e6365206f7574206f6620626f756e6400000000000000000000000081525060200191505060405180910390fd5b61040083106122e7576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260148152602001807f646566656e6365206f7574206f6620626f756e6400000000000000000000000081525060200191505060405180910390fd5b63100000008210612360576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260148152602001807f646566656e6365206f7574206f6620626f756e6400000000000000000000000081525060200191505060405180910390fd5b600060268560ff16901b9050601c84901b811790508281179150509392505050565b601081565b6000601260ff168260ff16106123a057600190506123b9565b6123b685858460ff16601260ff168702016121f5565b90505b949350505050565b6000613fff60ba83901c169050919050565b6000806000806123e2856124e2565b9250925092506000601260ff1682816123f757fe5b0490506124058484836121a1565b612477576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252600f8152602001807f696e76616c6964207465616d206964000000000000000000000000000000000081525060200191505060405180910390fd5b60006124848585846121f5565b90506000601260ff16848161249557fe5b0690506124a6888383600080613d1a565b9650505050505050919050565b602081565b60006124c382614d3f565b60018260ff16601981106124d357fe5b600d0201600a01549050919050565b6000806000601f602685901c166103ff601c86901c16630fffffff86169250925092509193909250565b60006007607783901c169050919050565b600181565b600065080000000000821061259f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601a8152602001807f707265764c6561677565496478206f7574206f6620626f756e6400000000000081525060200191505060405180910390fd5b741ffffffffffc0000000000000000000000000000001983169250607a82901b8317925082905092915050565b600681565b600181565b600060018460ff16601981106125e857fe5b600d020160000183815481106125fa57fe5b9060005260206000209060050201600401549050612619848483613489565b61268b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601e8152602001807f63616e6e6f74207472616e736665722061206e6f6e2d626f74207465616d000081525060200191505060405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16141561272e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252600f8152602001807f696e76616c69642061646472657373000000000000000000000000000000000081525060200191505060405180910390fd5b6127366152cf565b6000601260ff1690505b601960ff1681101561277057600182826019811061275a57fe5b6020020181815250508080600101915050612740565b5060405180604001604052808281526020018473ffffffffffffffffffffffffffffffffffffffff1681525060018660ff16601981106127ac57fe5b600d020160000185815481106127be57fe5b906000526020600020906005020160030160008481526020019081526020016000206000820151816000019060196127f79291906152f2565b5060208201518160190160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555090505060018560ff166019811061285257fe5b600d0201600001848154811061286457fe5b90600052602060002090600502016004016000815480929190600101919050555060006128928686856121f5565b90507f77b66eb1e6d2bc131b79be4213ae7f08f29351c01060e10bcc0302278067bf388185604051808381526020018273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019250505060405180910390a1505050505050565b600581565b600081565b6000600f607d83901c169050919050565b601081565b600080600080612936856124e2565b925092509250612947838383612f42565b9350505050919050565b6000613fff60e483901c169050919050565b60006507ffffffffff607a83901c169050919050565b600060208260ff16106129f4576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601c8152602001807f63757272656e7453686972744e756d206f7574206f6620626f756e640000000081525060200191505060405180910390fd5b7503e00000000000000000000000000000000000000000198316925060a58260ff16901b8317925082905092915050565b600080600160190390505b60008160ff1610612a6157612a4583826135c9565b15612a535780915050612a67565b808060019003915050612a30565b50601990505b919050565b6000613fff60c883901c169050919050565b612a86615332565b612a8e615354565b612a96615332565b612a9e615332565b6000600a8781612aaa57fe5b069050600080600060046201552f8b81612ac057fe5b0681612ac857fe5b06905060048a901c995060038960ff161015612b105760c885600060ff1660058110612af057fe5b602002019061ffff16908161ffff16815250506000925060009150612d76565b60088960ff161015612b8457602885600060ff1660058110612b2e57fe5b602002019061ffff16908161ffff168152505060a085600360ff1660058110612b5357fe5b602002019061ffff16908161ffff16815250506001925060078960ff168b0181612b7957fe5b066001019150612d75565b600a8960ff161015612bd35760a085600260ff1660058110612ba257fe5b602002019061ffff16908161ffff16815250506002925060078960ff168b0181612bc857fe5b066001019150612d74565b600c8960ff161015612c4757608285600260ff1660058110612bf157fe5b602002019061ffff16908161ffff1681525050604685600060ff1660058110612c1657fe5b602002019061ffff16908161ffff16815250506004925060078960ff168b0181612c3c57fe5b066001019150612d73565b600e8960ff161015612cbb57608285600260ff1660058110612c6557fe5b602002019061ffff16908161ffff1681525050604685600360ff1660058110612c8a57fe5b602002019061ffff16908161ffff16815250506005925060078960ff168b0181612cb057fe5b066001019150612d72565b60108960ff161015612d1e5760a085600060ff1660058110612cd957fe5b602002019061ffff16908161ffff1681525050604685600360ff1660058110612cfe57fe5b602002019061ffff16908161ffff16815250506003925060069150612d71565b60a085600060ff1660058110612d3057fe5b602002019061ffff16908161ffff1681525050604685600360ff1660058110612d5557fe5b602002019061ffff16908161ffff168152505060039250600391505b5b5b5b5b5b60338a901c9950600080600090505b600560ff168160ff161015612e63576000878260ff1660058110612da557fe5b602002015161ffff161415612de65760328c81612dbe57fe5b06888260ff1660058110612dce57fe5b602002019061ffff16908161ffff1681525050612e37565b6064878260ff1660058110612df757fe5b602002015160328e81612e0657fe5b060261ffff1681612e1357fe5b04888260ff1660058110612e2357fe5b602002019061ffff16908161ffff16815250505b60068c901c9b50878160ff1660058110612e4d57fe5b6020020151820191508080600101915050612d85565b5060fa8161ffff161015612ee4576000600560ff168260fa0361ffff1681612e8757fe5b04905060008090505b60058160ff161015612ee15781898260ff1660058110612eac57fe5b602002015101898260ff1660058110612ec157fe5b602002019061ffff16908161ffff16815250508080600101915050612e90565b50505b8660405180608001604052808760ff1660ff1681526020018660ff1660ff1681526020018560ff1660ff1681526020018460ff1660ff1681525098509850505050505050509250929050565b6000613fff60d683901c169050919050565b6000612f4d84614d3f565b612f578484614dfe565b60018460ff1660198110612f6757fe5b600d02016000018381548110612f7957fe5b9060005260206000209060050201600301600083815260200190815260200160002060190160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690509392505050565b6000612fd383614d3f565b612fdd8383614dfe565b60018360ff1660198110612fed57fe5b600d02016000018281548110612fff57fe5b906000526020600020906005020160000154905092915050565b600381565b6000613fff60f283901c169050919050565b6000600180607684901c16149050919050565b60006507ffffffffff608183901c169050919050565b600081565b60006507ffffffffff60d583901c169050919050565b60006007606f83901c169050919050565b600061309083614d3f565b61309a8383614dfe565b608060ff1660018460ff16601981106130af57fe5b600d020160000183815481106130c157fe5b90600052602060002090600502016000015402905092915050565b6130e5826140ae565b80156130f657506130f5816136dc565b5b613168576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260198152602001807f756e6578697374656e7420706c61796572206f72207465616d0000000000000081525060200191505060405180910390fd5b600061317383614cb1565b90506000819050600061318583614488565b9050838114156131fd576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260208152602001807f63616e6e6f74207472616e7366657220746f206f726967696e616c207465616d81525060200191505060405180910390fd5b613206816134f7565b1580156132195750613217846134f7565b155b61326e576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260368152602001806154c66036913960400191505060405180910390fd5b600061327984614bd1565b9050600061328686612a25565b9050601960ff168160ff1614156132e8576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602881526020018061551e6028913960400191505060405180910390fd5b6132f2848761413e565b93506132fe8482612979565b935061330a8443613ee1565b935083600080898152602001908152602001600020819055506000806000613331866124e2565b9250925092506001808460ff166019811061334857fe5b600d0201600001838154811061335a57fe5b90600052602060002090600502016003016000838152602001908152602001600020600001866019811061338a57fe5b0181905550613398896124e2565b8093508194508295505050508960018460ff16601981106133b557fe5b600d020160000183815481106133c757fe5b906000526020600020906005020160030160008381526020019081526020016000206000018560ff16601981106133fa57fe5b01819055507f54a4f48232284e6aff96e3a82633881625eb95d9b9865baed16f627a6a1cbffb8a8a604051808381526020018281526020019250505060405180910390a17f65a4d4a8a0afb474d2e9465a6a1a41bb88fd04f41152ba070421f1b3771f15d88a88604051808381526020018281526020019250505060405180910390a150505050505050505050565b60008073ffffffffffffffffffffffffffffffffffffffff166134ad858585612f42565b73ffffffffffffffffffffffffffffffffffffffff161490509392505050565b601281565b6101475481565b6000613fff60ac83901c169050919050565b6101485481565b601981565b600080600080613506856124e2565b925092509250613517838383613489565b9350505050919050565b600581565b6000613531826140ae565b6135a3576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260118152602001807f756e6578697374656e7420706c6179657200000000000000000000000000000081525060200191505060405180910390fd5b60006135b66135b184614cb1565b614488565b90506135c181612927565b915050919050565b6000806000806135d8866124e2565b9250925092506135e9838383613489565b1561363f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602a815260200180615429602a913960400191505060405180910390fd5b600060018460ff166019811061365157fe5b600d0201600001838154811061366357fe5b906000526020600020906005020160030160008381526020019081526020016000206000018660ff166019811061369657fe5b01549050600160120360ff168660ff1611156136c65760008114806136bb5750600181145b9450505050506136d1565b600181149450505050505b92915050565b600581565b6000806000806136eb856124e2565b9250925092506136fc8383836121a1565b9350505050919050565b600080600090505b600560ff168160ff1610156137b8576140008a8260ff166005811061372f57fe5b602002015161ffff16106137ab576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260128152602001807f736b696c6c206f7574206f6620626f756e64000000000000000000000000000081525060200191505060405180910390fd5b808060010191505061370e565b50600a86600060ff16600481106137cb57fe5b602002015160ff1610613846576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260168152602001807f706f74656e7469616c206f7574206f6620626f756e640000000000000000000081525060200191505060405180910390fd5b600686600160ff166004811061385857fe5b602002015160ff16106138d3576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260188152602001807f666f72776172646e657373206f7574206f6620626f756e64000000000000000081525060200191505060405180910390fd5b600886600260ff16600481106138e557fe5b602002015160ff1610613960576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260188152602001807f6c6566697473686e657373206f7574206f6620626f756e64000000000000000081525060200191505060405180910390fd5b600886600360ff166004811061397257fe5b602002015160ff16106139ed576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601b8152602001807f616767726573736976656e657373206f7574206f6620626f756e64000000000081525060200191505060405180910390fd5b600086600260ff16600481106139ff57fe5b602002015160ff161415613a7e57600086600160ff1660048110613a1f57fe5b602002015160ff1614613a7d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602b815260200180615453602b913960400191505060405180910390fd5b5b60088360ff1610613af7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601d8152602001807f67616d65734e6f6e53746f7070696e67206f7574206f6620626f756e6400000081525060200191505060405180910390fd5b6140008810613b51576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260238152602001806154066023913960400191505060405180910390fd5b600087118015613b6657506508000000000087105b613bd8576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260158152602001807f706c617965724964206f7574206f6620626f756e64000000000000000000000081525060200191505060405180910390fd5b60008090505b600560ff168160ff161015613c2a57600e600182010260ff166101000361ffff168a8260ff1660058110613c0e57fe5b602002015161ffff16901b821791508080600101915050613bde565b5060ac88901b81179050608187901b81179050607d86600060ff1660048110613c4f57fe5b602002015160ff16901b81179050607a86600160ff1660048110613c6f57fe5b602002015160ff16901b81179050607786600260ff1660048110613c8f57fe5b602002015160ff16901b81179050607685613cab576000613cae565b60015b60ff16901b81179050607584613cc5576000613cc8565b60015b60ff16901b8117905060728360ff16901b81179050606f8260ff16901b81179050606c86600360ff1660048110613cfb57fe5b602002015160ff16901b8117905098975050505050505050565b600281565b60008086118015613d3057506508000000000086105b613da2576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260158152602001807f706c617965724964206f7574206f6620626f756e64000000000000000000000081525060200191505060405180910390fd5b600060d587901b9050613db5818761413e565b9050613dc18186612979565b9050613dcd8185612522565b9050613dd98184613ee1565b90508091505095945050505050565b6001151561014960009054906101000a900460ff16151514613e72576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260178152602001807f63616e6e6f7420696e697469616c697a6520747769636500000000000000000081525060200191505060405180910390fd5b613e7b42614de4565b61014781905550613e8b81614e74565b600061014960006101000a81548160ff02191690831515021790555050565b600281565b600881565b6000613ebf82614d3f565b60018260ff1660198110613ecf57fe5b600d0201600001805490509050919050565b60006408000000008210613f5d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601a8152602001807f6c61737453616c65426c6f636b206f7574206f6620626f756e6400000000000081525060200191505060405180910390fd5b6f03ffffffff80000000000000000000001983169250605782901b8317925082905092915050565b6000613f90826140ae565b614002576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260118152602001807f756e6578697374656e7420706c6179657200000000000000000000000000000081525060200191505060405180910390fd5b600080600084815260200190815260200160002054149050919050565b600081565b600781565b6001816019811061403657fe5b600d02016000915090508060010160009054906101000a900460ff16908060060160009054906101000a900460ff16908060060160019054906101000a900460ff16908060070154908060080160009054906101000a900460ff169080600901549080600a01549080600b01549080600c0154905089565b6000808214156140c15760009050614109565b600080600084815260200190815260200160002054146140e45760019050614109565b60008060006140f2856124e2565b925092509250614103838383614fca565b93505050505b919050565b6000600860ff1661411f8484614cfb565b02905092915050565b600181565b60006007607a83901c169050919050565b60006508000000000082106141bb576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601b8152602001807f63757272656e745465616d496478206f7574206f6620626f756e64000000000081525060200191505060405180910390fd5b7a1ffffffffffc000000000000000000000000000000000000000000198316925060aa82901b8317925082905092915050565b60006407ffffffff605783901c169050919050565b600080600080614212856124e2565b9250925092506000601260ff16828161422757fe5b0490506000601260ff16838161423957fe5b0690506000608060ff16838161424b57fe5b0490506142598686856121a1565b6142cb576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252600f8152602001807f696e76616c6964207465616d206964000000000000000000000000000000000081525060200191505060405180910390fd5b600086868585604051602001808560ff1660ff1681526020018481526020018381526020018260ff1660ff1681526020019450505050506040516020818303038152906040528051906020012060001c90506000601e601060018a60ff166019811061433357fe5b600d0201600001898154811061434557fe5b906000526020600020906005020160020160008681526020019081526020016000205402601e6101475402018161437857fe5b0490506143878285838d614fe7565b98505050505050505050919050565b601081565b6143a3615332565b6143ac8261301e565b816000600581106143b957fe5b602002019061ffff16908161ffff16815250506143d582612951565b816001600581106143e257fe5b602002019061ffff16908161ffff16815250506143fe82612f30565b8160026005811061440b57fe5b602002019061ffff16908161ffff168152505061442782612a6c565b8160036005811061443457fe5b602002019061ffff16908161ffff1681525050614450826123c1565b8160046005811061445d57fe5b602002019061ffff16908161ffff1681525050919050565b6000600180607584901c16149050919050565b60006507ffffffffff60aa83901c169050919050565b600381565b60006144ae83614d3f565b8160018460ff16601981106144bf57fe5b600d0201600b01819055504260018460ff16601981106144db57fe5b600d02016009018190555092915050565b6001151561014960009054906101000a900460ff16151514614576576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260178152602001807f63616e6e6f7420696e697469616c697a6520747769636500000000000000000081525060200191505060405180910390fd5b61457f42614de4565b610147819055506000600190505b60198160ff1610156145af576145a281614e74565b808060010191505061458d565b50600061014960006101000a81548160ff021916908315150217905550565b6145d6615376565b6145de615399565b60006a02000000000000000000008410614660576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601e8152602001807f7461637469637349642073686f756c642066697420696e20363120626974000081525060200191505060405180910390fd5b603f84169050600684901c935060008090505b600a8160ff1610156146c557600180861614614690576000614693565b60015b838260ff16600a81106146a257fe5b602002019015159081151581525050600185901c94508080600101915050614673565b5060008090505b600b8160ff16101561479f57601f8516848260ff16600b81106146eb57fe5b602002019060ff16908160ff1681525050601960ff16848260ff16600b811061471057fe5b602002015160ff161061478b576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260168152602001807f696e636f7272656374206c696e65757020656e7472790000000000000000000081525060200191505060405180910390fd5b600585901c945080806001019150506146cc565b509193909250565b60006007607283901c169050919050565b600381565b60008060006147cb856124e2565b9250925092506147dd83838387615040565b7f77b66eb1e6d2bc131b79be4213ae7f08f29351c01060e10bcc0302278067bf388585604051808381526020018273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019250505060405180910390a15050505050565b600060408260ff16106148ca576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601d8152602001807f7461637469637349642073686f756c642066697420696e20362062697400000081525060200191505060405180910390fd5b60008260ff16905060008090505b600a8160ff161015614926578060010260060160ff16858260ff16600a81106148fd57fe5b602002015161490d576000614910565b60015b60ff16901b8217915080806001019150506148d8565b5060008090505b600b8160ff1610156149fc57601960ff16868260ff16600b811061494d57fe5b602002015160ff16106149c8576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260168152602001807f696e636f7272656374206c696e65757020656e7472790000000000000000000081525060200191505060405180910390fd5b8060050260100160ff16868260ff16600b81106149e157fe5b602002015160ff16901b82179150808060010191505061492d565b50809150509392505050565b614a106152cf565b6000806000614a1e856124e2565b925092509250614a2f8383836121a1565b614aa1576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252600f8152602001807f696e76616c6964207465616d206964000000000000000000000000000000000081525060200191505060405180910390fd5b614aac838383613489565b15614afd5760008090505b601960ff168160ff161015614af757614ad284848484612387565b858260ff1660198110614ae157fe5b6020020181815250508080600101915050614ab7565b50614bc9565b60008090505b601960ff168160ff161015614bc757600060018560ff1660198110614b2457fe5b600d02016000018481548110614b3657fe5b906000526020600020906005020160030160008481526020019081526020016000206000018260ff1660198110614b6957fe5b015490506000811415614b9f57614b8285858585612387565b868360ff1660198110614b9157fe5b602002018181525050614bb9565b80868360ff1660198110614baf57fe5b6020020181815250505b508080600101915050614b03565b505b505050919050565b6000601f60a583901c169050919050565b60008160018460ff1660198110614bf557fe5b600d020160040160018560ff1660198110614c0c57fe5b600d020160060160019054906101000a900460ff1660ff1660028110614c2e57fe5b018190555060018360ff1660198110614c4357fe5b600d020160060160019054906101000a900460ff1660010360018460ff1660198110614c6b57fe5b600d020160060160016101000a81548160ff021916908360ff1602179055504260018460ff1660198110614c9b57fe5b600d0201600a018190555092915050565b600481565b6000614cbc82613f85565b15614cd157614cca826123d3565b9050614ce7565b6000808381526020019081526020016000205490505b919050565b608081565b600281565b600181565b6000601060ff16614d0c8484612fc8565b02905092915050565b6000614d2082614d3f565b60018260ff1660198110614d3057fe5b600d0201600901549050919050565b60008160ff16118015614d55575060198160ff16105b614dc7576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260178152602001807f74696d655a6f6e6520646f6573206e6f7420657869737400000000000000000081525060200191505060405180910390fd5b50565b6000600c6301e13380830281614ddc57fe5b049050919050565b60006301e13380600c830281614df657fe5b049050919050565b60018260ff1660198110614e0e57fe5b600d0201600001805490508110614e70576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602781526020018061547e6027913960400191505060405180910390fd5b5050565b614e7c6153bc565b600181600001818152505060018260ff1660198110614e9757fe5b600d020160000181908060018154018082558091505090600182039060005260206000209060050201600090919290919091506000820151816000015560208201518160010160006101000a81548160ff021916908360ff160217905550604082015181600401555050506001808360ff1660198110614f1357fe5b600d0201600001600081548110614f2657fe5b90600052602060002090600502016002016000808152602001908152602001600020819055506000801b60018360ff1660198110614f6057fe5b600d0201600201600060028110614f7357fe5b01819055507fc5d195855a200aa90e2052bcc795cedbc84c2a26556b1d5113b5a30c96003a0b82600080604051808460ff1660ff168152602001838152602001828152602001935050505060405180910390a15050565b6000601260ff16614fdb858561410e565b02821090509392505050565b600080614ff48685612086565b8161ffff169150809750819250505061500b615332565b615013615354565b61501d8888612a7e565b9150915061503382848784600080600080613706565b9350505050949350505050565b61504984614d3f565b6150538484614dfe565b61505e848484613489565b156150d1576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601e8152602001807f63616e6e6f74207472616e736665722061206e6f6e2d626f74207465616d000081525060200191505060405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161415615157576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260218152602001806154a56021913960400191505060405180910390fd5b8073ffffffffffffffffffffffffffffffffffffffff1660018560ff166019811061517e57fe5b600d0201600001848154811061519057fe5b9060005260206000209060050201600301600084815260200190815260200160002060190160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff161415615243576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260228152602001806154fc6022913960400191505060405180910390fd5b8060018560ff166019811061525457fe5b600d0201600001848154811061526657fe5b9060005260206000209060050201600301600084815260200190815260200160002060190160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050505050565b604051806103200160405280601990602082028038833980820191505090505090565b8260198101928215615321579160200282015b82811115615320578251825591602001919060010190615305565b5b50905061532e91906153e0565b5090565b6040518060a00160405280600590602082028038833980820191505090505090565b6040518060800160405280600490602082028038833980820191505090505090565b604051806101600160405280600b90602082028038833980820191505090505090565b604051806101400160405280600a90602082028038833980820191505090505090565b604051806060016040528060008152602001600060ff168152602001600081525090565b61540291905b808211156153fe5760008160009055506001016153e6565b5090565b9056fe6d6f6e74684f664269727468496e556e697854696d65206f7574206f6620626f756e6463616e6e6f742071756572792061626f757420746865207368697274206f66206120426f74205465616d6c6566746973686e65732063616e206f6e6c79206265207a65726f20666f7220676f616c6b656570657273636f756e74727920646f6573206e6f7420657869737420696e20746869732074696d655a6f6e6563616e6e6f74207472616e7366657220746f2061206e756c6c206164647265737363616e6e6f74207472616e7366657220706c61796572207768656e206174206c65617374206f6e65207465616d206973206120626f74627579657220616e642073656c6c657220617265207468652073616d652061646472746172676574207465616d20666f72207472616e7366657220697320616c72656164792066756c6ca165627a7a72305820d398d5a3d2c2f43e48c6f7ebefd7b9f9aa2d877f772f4e6cbd2e0a3acfa442940029`

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

// AssetsPlayerTransferIterator is returned from FilterPlayerTransfer and is used to iterate over the raw logs and unpacked data for PlayerTransfer events raised by the Assets contract.
type AssetsPlayerTransferIterator struct {
	Event *AssetsPlayerTransfer // Event containing the contract specifics and raw log

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
func (it *AssetsPlayerTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AssetsPlayerTransfer)
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
		it.Event = new(AssetsPlayerTransfer)
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
func (it *AssetsPlayerTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AssetsPlayerTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AssetsPlayerTransfer represents a PlayerTransfer event raised by the Assets contract.
type AssetsPlayerTransfer struct {
	PlayerId     *big.Int
	TeamIdTarget *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterPlayerTransfer is a free log retrieval operation binding the contract event 0x54a4f48232284e6aff96e3a82633881625eb95d9b9865baed16f627a6a1cbffb.
//
// Solidity: event PlayerTransfer(uint256 playerId, uint256 teamIdTarget)
func (_Assets *AssetsFilterer) FilterPlayerTransfer(opts *bind.FilterOpts) (*AssetsPlayerTransferIterator, error) {

	logs, sub, err := _Assets.contract.FilterLogs(opts, "PlayerTransfer")
	if err != nil {
		return nil, err
	}
	return &AssetsPlayerTransferIterator{contract: _Assets.contract, event: "PlayerTransfer", logs: logs, sub: sub}, nil
}

// WatchPlayerTransfer is a free log subscription operation binding the contract event 0x54a4f48232284e6aff96e3a82633881625eb95d9b9865baed16f627a6a1cbffb.
//
// Solidity: event PlayerTransfer(uint256 playerId, uint256 teamIdTarget)
func (_Assets *AssetsFilterer) WatchPlayerTransfer(opts *bind.WatchOpts, sink chan<- *AssetsPlayerTransfer) (event.Subscription, error) {

	logs, sub, err := _Assets.contract.WatchLogs(opts, "PlayerTransfer")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AssetsPlayerTransfer)
				if err := _Assets.contract.UnpackLog(event, "PlayerTransfer", log); err != nil {
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
