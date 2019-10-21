// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package leagues

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

// LeaguesABI is the input ABI used to generate the binding from.
const LeaguesABI = "[{\"inputs\":[],\"constant\":true,\"name\":\"IDX_MD\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"encodedSkills\"}],\"constant\":true,\"name\":\"getSkills\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"dna\"},{\"type\":\"uint256\",\"name\":\"playerCreationMonth\"}],\"constant\":true,\"name\":\"computeBirthMonth\",\"outputs\":[{\"type\":\"uint16\",\"name\":\"\"},{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint8\",\"name\":\"timeZone\"}],\"constant\":true,\"name\":\"countCountries\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"MATCHES_PER_DAY\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"encodedSkills\"}],\"constant\":true,\"name\":\"getSumOfSkills\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"IDX_R\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"SK_END\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint8\",\"name\":\"timeZone\"},{\"type\":\"uint256\",\"name\":\"countryIdxInTZ\"},{\"type\":\"uint256\",\"name\":\"teamIdxInCountry\"}],\"constant\":true,\"name\":\"_teamExistsInCountry\",\"outputs\":[{\"type\":\"bool\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"encodedSkills\"}],\"constant\":true,\"name\":\"getAggressiveness\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"playerId\"}],\"constant\":true,\"name\":\"getPlayerAgeInMonths\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint8\",\"name\":\"timeZone\"},{\"type\":\"uint256\",\"name\":\"countryIdxInTZ\"},{\"type\":\"uint256\",\"name\":\"val\"}],\"constant\":true,\"name\":\"encodeTZCountryAndVal\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"DAYS_PER_ROUND\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint8\",\"name\":\"timeZone\"},{\"type\":\"uint256\",\"name\":\"countryIdxInTZ\"},{\"type\":\"uint256\",\"name\":\"teamIdxInCountry\"},{\"type\":\"uint8\",\"name\":\"shirtNum\"}],\"constant\":true,\"name\":\"getDefaultPlayerIdForTeamInCountry\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"encodedSkills\"}],\"constant\":true,\"name\":\"getEndurance\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"playerId\"}],\"constant\":true,\"name\":\"getPlayerStateAtBirth\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"MAX_PLAYER_AGE_AT_BIRTH\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint8\",\"name\":\"timeZone\"}],\"constant\":true,\"name\":\"getLastUpdateTime\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"encoded\"}],\"constant\":true,\"name\":\"decodeTZCountryAndVal\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"},{\"type\":\"uint256\",\"name\":\"\"},{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"encodedSkills\"}],\"constant\":true,\"name\":\"getLeftishness\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"IDX_D\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"state\"},{\"type\":\"uint256\",\"name\":\"value\"}],\"constant\":true,\"name\":\"setPrevPlayerTeamId\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"IDX_LC\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"playerId\"}],\"constant\":true,\"name\":\"getCurrentTeamIdFromPlayerId\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"FREEVERSE\",\"outputs\":[{\"type\":\"address\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint8\",\"name\":\"timeZone\"},{\"type\":\"uint256\",\"name\":\"countryIdxInTZ\"},{\"type\":\"address\",\"name\":\"addr\"}],\"constant\":false,\"name\":\"transferFirstBotToAddr\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"IDX_LR\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"SK_SHO\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"encodedSkills\"}],\"constant\":true,\"name\":\"getPotential\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"LEAGUES_PER_DIV\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"teamId\"}],\"constant\":true,\"name\":\"getOwnerTeam\",\"outputs\":[{\"type\":\"address\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"encodedSkills\"}],\"constant\":true,\"name\":\"getSpeed\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"playerState\"}],\"constant\":true,\"name\":\"getPrevPlayerTeamId\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"state\"},{\"type\":\"uint8\",\"name\":\"currentShirtNum\"}],\"constant\":true,\"name\":\"setCurrentShirtNum\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"teamId\"}],\"constant\":true,\"name\":\"getFreeShirt\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"encodedSkills\"}],\"constant\":true,\"name\":\"getDefence\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"dna\"},{\"type\":\"uint8\",\"name\":\"shirtNum\"}],\"constant\":true,\"name\":\"computeSkills\",\"outputs\":[{\"type\":\"uint16[5]\",\"name\":\"\"},{\"type\":\"uint8[4]\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"encodedSkills\"}],\"constant\":true,\"name\":\"getPass\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"MATCHDAYS\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint8\",\"name\":\"timeZone\"},{\"type\":\"uint256\",\"name\":\"countryIdxInTZ\"},{\"type\":\"uint256\",\"name\":\"teamIdxInCountry\"}],\"constant\":true,\"name\":\"getOwnerTeamInCountry\",\"outputs\":[{\"type\":\"address\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint8\",\"name\":\"timeZone\"},{\"type\":\"uint256\",\"name\":\"countryIdxInTZ\"}],\"constant\":true,\"name\":\"getNDivisionsInCountry\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"IDX_CR\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"encodedSkills\"}],\"constant\":true,\"name\":\"getShoot\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"encodedSkills\"}],\"constant\":true,\"name\":\"getAlignedLastHalf\",\"outputs\":[{\"type\":\"bool\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"encodedSkills\"}],\"constant\":true,\"name\":\"getPlayerIdFromSkills\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"IDX_GK\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"playerState\"}],\"constant\":true,\"name\":\"getPlayerIdFromState\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"encodedSkills\"}],\"constant\":true,\"name\":\"getInjuryWeeksLeft\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint8\",\"name\":\"timeZone\"},{\"type\":\"uint256\",\"name\":\"countryIdxInTZ\"}],\"constant\":true,\"name\":\"countTeams\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"playerId\"},{\"type\":\"uint256\",\"name\":\"teamIdTarget\"}],\"constant\":false,\"name\":\"transferPlayer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint8\",\"name\":\"timeZone\"},{\"type\":\"uint256\",\"name\":\"countryIdxInTZ\"},{\"type\":\"uint256\",\"name\":\"teamIdxInCountry\"}],\"constant\":true,\"name\":\"isBotTeamInCountry\",\"outputs\":[{\"type\":\"bool\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"PLAYERS_PER_TEAM_INIT\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"gameDeployMonth\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"encodedSkills\"}],\"constant\":true,\"name\":\"getMonthOfBirth\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"currentRound\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"PLAYERS_PER_TEAM_MAX\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"teamId\"}],\"constant\":true,\"name\":\"isBotTeam\",\"outputs\":[{\"type\":\"bool\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"IDX_MF\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"playerId\"}],\"constant\":true,\"name\":\"getOwnerPlayer\",\"outputs\":[{\"type\":\"address\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"teamId\"},{\"type\":\"uint8\",\"name\":\"shirtNum\"}],\"constant\":true,\"name\":\"isFreeShirt\",\"outputs\":[{\"type\":\"bool\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"N_SKILLS\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"teamId\"}],\"constant\":true,\"name\":\"teamExists\",\"outputs\":[{\"type\":\"bool\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint16[5]\",\"name\":\"skills\"},{\"type\":\"uint256\",\"name\":\"monthOfBirth\"},{\"type\":\"uint256\",\"name\":\"playerId\"},{\"type\":\"uint8[4]\",\"name\":\"birthTraits\"},{\"type\":\"bool\",\"name\":\"alignedLastHalf\"},{\"type\":\"bool\",\"name\":\"redCardLastGame\"},{\"type\":\"uint8\",\"name\":\"gamesNonStopping\"},{\"type\":\"uint8\",\"name\":\"injuryWeeksLeft\"}],\"constant\":true,\"name\":\"encodePlayerSkills\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"encoded\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"IDX_M\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"playerId\"},{\"type\":\"uint256\",\"name\":\"currentTeamId\"},{\"type\":\"uint8\",\"name\":\"currentShirtNum\"},{\"type\":\"uint256\",\"name\":\"prevPlayerTeamId\"},{\"type\":\"uint256\",\"name\":\"lastSaleBlock\"}],\"constant\":true,\"name\":\"encodePlayerState\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint8\",\"name\":\"tz\"}],\"constant\":false,\"name\":\"initSingleTZ\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"SK_PAS\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"TEAMS_PER_LEAGUE\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint8\",\"name\":\"timeZone\"}],\"constant\":true,\"name\":\"getNCountriesInTZ\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"state\"},{\"type\":\"uint256\",\"name\":\"lastSaleBlock\"}],\"constant\":true,\"name\":\"setLastSaleBlock\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"playerId\"}],\"constant\":true,\"name\":\"isVirtualPlayer\",\"outputs\":[{\"type\":\"bool\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"NULL_ADDR\",\"outputs\":[{\"type\":\"address\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"IDX_LCR\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"constant\":true,\"name\":\"_timeZones\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"nCountriesToAdd\"},{\"type\":\"uint8\",\"name\":\"newestOrgMapIdx\"},{\"type\":\"uint8\",\"name\":\"newestSkillsIdx\"},{\"type\":\"bytes32\",\"name\":\"scoresRoot\"},{\"type\":\"uint8\",\"name\":\"updateCycleIdx\"},{\"type\":\"uint256\",\"name\":\"lastActionsSubmissionTime\"},{\"type\":\"uint256\",\"name\":\"lastUpdateTime\"},{\"type\":\"bytes32\",\"name\":\"actionsRoot\"},{\"type\":\"uint256\",\"name\":\"lastMarketClosureBlockNum\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"playerId\"}],\"constant\":true,\"name\":\"playerExists\",\"outputs\":[{\"type\":\"bool\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint8\",\"name\":\"timeZone\"},{\"type\":\"uint256\",\"name\":\"countryIdxInTZ\"}],\"constant\":true,\"name\":\"getNTeamsInCountry\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"FREE_PLAYER_ID\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"encodedSkills\"}],\"constant\":true,\"name\":\"getForwardness\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"playerState\"},{\"type\":\"uint256\",\"name\":\"teamId\"}],\"constant\":true,\"name\":\"setCurrentTeamId\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"playerState\"}],\"constant\":true,\"name\":\"getLastSaleBlock\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"playerId\"}],\"constant\":true,\"name\":\"getPlayerSkillsAtBirth\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"MIN_PLAYER_AGE_AT_BIRTH\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"encodedSkills\"}],\"constant\":true,\"name\":\"getSkillsVec\",\"outputs\":[{\"type\":\"uint16[5]\",\"name\":\"skills\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"encodedSkills\"}],\"constant\":true,\"name\":\"getRedCardLastGame\",\"outputs\":[{\"type\":\"bool\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"playerState\"}],\"constant\":true,\"name\":\"getCurrentTeamId\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"IDX_F\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint8\",\"name\":\"timeZone\"},{\"type\":\"bytes32\",\"name\":\"root\"}],\"constant\":false,\"name\":\"setActionsRoot\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"nonpayable\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":false,\"name\":\"init\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"tactics\"}],\"constant\":true,\"name\":\"decodeTactics\",\"outputs\":[{\"type\":\"uint8[11]\",\"name\":\"lineup\"},{\"type\":\"bool[10]\",\"name\":\"extraAttack\"},{\"type\":\"uint8\",\"name\":\"tacticsId\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"encodedSkills\"}],\"constant\":true,\"name\":\"getGamesNonStopping\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"SK_DEF\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"teamId\"},{\"type\":\"address\",\"name\":\"addr\"}],\"constant\":false,\"name\":\"transferTeam\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint8[11]\",\"name\":\"lineup\"},{\"type\":\"bool[10]\",\"name\":\"extraAttack\"},{\"type\":\"uint8\",\"name\":\"tacticsId\"}],\"constant\":true,\"name\":\"encodeTactics\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"teamId\"}],\"constant\":true,\"name\":\"getPlayerIdsInTeam\",\"outputs\":[{\"type\":\"uint256[25]\",\"name\":\"playerIds\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"playerState\"}],\"constant\":true,\"name\":\"getCurrentShirtNum\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint8\",\"name\":\"tz\"},{\"type\":\"bytes32\",\"name\":\"root\"}],\"constant\":false,\"name\":\"setSkillsRoot\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"nonpayable\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"IDX_L\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint256\",\"name\":\"playerId\"}],\"constant\":true,\"name\":\"getPlayerState\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"TEAMS_PER_DIVISION\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"IDX_C\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"SK_SPE\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint8\",\"name\":\"timeZone\"},{\"type\":\"uint256\",\"name\":\"countryIdxInTZ\"}],\"constant\":true,\"name\":\"getNLeaguesInCountry\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint8\",\"name\":\"timeZone\"}],\"constant\":true,\"name\":\"getLastActionsSubmissionTime\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"indexed\":false,\"type\":\"uint256\",\"name\":\"teamId\"},{\"indexed\":false,\"type\":\"address\",\"name\":\"to\"}],\"type\":\"event\",\"name\":\"TeamTransfer\",\"anonymous\":false},{\"inputs\":[{\"indexed\":false,\"type\":\"uint8\",\"name\":\"timezone\"},{\"indexed\":false,\"type\":\"uint256\",\"name\":\"countryIdxInTZ\"},{\"indexed\":false,\"type\":\"uint256\",\"name\":\"divisionIdxInCountry\"}],\"type\":\"event\",\"name\":\"DivisionCreation\",\"anonymous\":false},{\"inputs\":[{\"indexed\":false,\"type\":\"uint256\",\"name\":\"playerId\"},{\"indexed\":false,\"type\":\"uint256\",\"name\":\"state\"}],\"type\":\"event\",\"name\":\"PlayerStateChange\",\"anonymous\":false},{\"inputs\":[{\"type\":\"address\",\"name\":\"addr\"}],\"constant\":false,\"name\":\"setEngineAdress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[],\"constant\":true,\"name\":\"getEngineAddress\",\"outputs\":[{\"type\":\"address\",\"name\":\"\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint8\",\"name\":\"matchday\"},{\"type\":\"uint8\",\"name\":\"matchIdxInDay\"}],\"constant\":true,\"name\":\"getTeamsInMatch\",\"outputs\":[{\"type\":\"uint8\",\"name\":\"homeIdx\"},{\"type\":\"uint8\",\"name\":\"visitorIdx\"}],\"stateMutability\":\"pure\",\"payable\":false,\"type\":\"function\"},{\"inputs\":[{\"type\":\"uint8\",\"name\":\"matchday\"},{\"type\":\"uint256[25][8]\",\"name\":\"prevLeagueState\"},{\"type\":\"uint256[8]\",\"name\":\"tacticsIds\"},{\"type\":\"uint256\",\"name\":\"currentVerseSeed\"}],\"constant\":true,\"name\":\"computeMatchday\",\"outputs\":[{\"type\":\"uint8[8]\",\"name\":\"scores\"}],\"stateMutability\":\"view\",\"payable\":false,\"type\":\"function\"}]"

// LeaguesBin is the compiled bytecode used for deploying new contracts.
const LeaguesBin = `0x6080604052600161014960006101000a81548160ff02191690831515021790555034801561002c57600080fd5b50615e93806200003d6000396000f3fe608060405234801561001057600080fd5b50600436106106565760003560e01c806380bac70911610351578063c2bc41cd116101c3578063e804e5191161010f578063ec71bc82116100ad578063f305a21c11610087578063f305a21c14612306578063f8ef7b9e1461232a578063f9d0723d1461234e578063fa80039b1461239d57610656565b8063ec71bc821461227c578063ec7ecec5146122a0578063f21f5a83146122e257610656565b8063e9e71652116100e9578063e9e71652146120b7578063eabf6a4b14612181578063eb78b7b7146121eb578063ec1c54231461222d57610656565b8063e804e51914612003578063e81e21bb14612045578063e945e96a1461206957610656565b8063cc7d473b1161017c578063d7e4e6d511610156578063d7e4e6d514611dcf578063dba6319e14611f04578063e1c7392a14611f53578063e6400ac414611f5d57610656565b8063cc7d473b14611d23578063cd2105e814611d69578063d7b63a1114611dab57610656565b8063c2bc41cd14611b83578063c37b1c2514611bc5578063c566b5bc14611c11578063c73f808d14611c53578063c79055d414611c95578063cc1cc3d714611cb957610656565b80639cc623401161029d578063b32aa2c11161023b578063b96b1a3011610215578063b96b1a3014611a3e578063bc1a97c114611ad0578063c04f6d5314611b16578063c258012b14611b6557610656565b8063b32aa2c11461198a578063b3f390b3146119d0578063b962709714611a1a57610656565b8063ab1b7c5e11610277578063ab1b7c5e146118b1578063ac5db9ee146118d5578063ad63bcbd146118f9578063af76cd011461193e57610656565b80639cc62340146117ef5780639f27112a14611813578063a3ceb7031461188057610656565b80638cc9a8d51161030a578063963fcc80116102e4578063963fcc8014611631578063976daaac1461168457806398981756146116a85780639c53e3fd146116ee57610656565b80638cc9a8d5146115595780638f3db4361461159f5780638f9da214146115c357610656565b806380bac7091461143657806383c31d3b1461149357806385982431146114b757806387f1e880146114d55780638a19c8bc146115175780638adddc9d1461153557610656565b806339644f21116104ea578063547d829811610436578063673fe242116103d457806378f4c718116103ae57806378f4c7181461132b57806379e765971461136d5780637b2566a5146113af578063800257d5146113fe57610656565b8063673fe2421461127f5780636f6c2ae0146112c55780637420a6061461130757610656565b8063595ef25b11610410578063595ef25b146111455780635adb40f5146111ca5780635becd9991461121957806365b4b4761461123d57610656565b8063547d82981461103957806355a6f86f146110df578063561b11181461112157610656565b806348d1e9c0116104a35780634bea2a691161047d5780634bea2a6914610f1e5780634db989fd14610f60578063507b172314610faf57806351585b4914610ff757610656565b806348d1e9c014610e4a578063492afc6914610e6e5780634b93f75314610edc57610656565b806339644f2114610cd15780633c2eb36014610d1b5780633d085f9614610d7657806340cd05fd14610d9a578063434807ad14610dbe5780634562a61814610e0057610656565b806320748ae8116105a95780632d0e08fd11610562578063369151db1161053c578063369151db14610bfb57806337a8630214610c1f57806337fd56af14610c6b57806338c96b5c14610c8f57610656565b80632d0e08fd14610b1e5780633260840b14610b635780633518dd1d14610bb957610656565b806320748ae81461099957806321ff8ae8146109f2578063228408b014610a10578063258e5d9014610a765780632665760814610ab85780632a238b0a14610afa57610656565b80630c85696c116106165780631884332c116105f05780631884332c146108945780631a6daba2146108b85780631fc7768f146109155780631ffeb3491461095757610656565b80630c85696c1461080a5780631060c9c21461082e578063169d29141461087057610656565b80623e32231461065b57806292bf781461067f578062aae8df146106c1578063032324c81461071c578063058672f9146107815780630abcd3e5146107c5575b600080fd5b6106636123e2565b604051808260ff1660ff16815260200191505060405180910390f35b6106ab6004803603602081101561069557600080fd5b81019080803590602001909291905050506123e7565b6040518082815260200191505060405180910390f35b6106f7600480360360408110156106d757600080fd5b8101908080359060200190929190803590602001909291905050506123f5565b604051808361ffff1661ffff1681526020018281526020019250505060405180910390f35b6107586004803603604081101561073257600080fd5b81019080803560ff169060200190929190803560ff16906020019092919050505061249f565b604051808360ff1660ff1681526020018260ff1660ff1681526020019250505060405180910390f35b6107c36004803603602081101561079757600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506125e3565b005b6107f4600480360360208110156107db57600080fd5b81019080803560ff169060200190929190505050612628565b6040518082815260200191505060405180910390f35b610812612655565b604051808260ff1660ff16815260200191505060405180910390f35b61085a6004803603602081101561084457600080fd5b810190808035906020019092919050505061265a565b6040518082815260200191505060405180910390f35b610878612694565b604051808260ff1660ff16815260200191505060405180910390f35b61089c612699565b604051808260ff1660ff16815260200191505060405180910390f35b6108fb600480360360608110156108ce57600080fd5b81019080803560ff169060200190929190803590602001909291908035906020019092919050505061269e565b604051808215151515815260200191505060405180910390f35b6109416004803603602081101561092b57600080fd5b81019080803590602001909291905050506126b5565b6040518082815260200191505060405180910390f35b6109836004803603602081101561096d57600080fd5b81019080803590602001909291905050506126c6565b6040518082815260200191505060405180910390f35b6109dc600480360360608110156109af57600080fd5b81019080803560ff16906020019092919080359060200190929190803590602001909291905050506126f2565b6040518082815260200191505060405180910390f35b6109fa61287f565b6040518082815260200191505060405180910390f35b610a6060048036036080811015610a2657600080fd5b81019080803560ff1690602001909291908035906020019092919080359060200190929190803560ff169060200190929190505050612884565b6040518082815260200191505060405180910390f35b610aa260048036036020811015610a8c57600080fd5b81019080803590602001909291905050506128be565b6040518082815260200191505060405180910390f35b610ae460048036036020811015610ace57600080fd5b81019080803590602001909291905050506128d0565b6040518082815260200191505060405180910390f35b610b026129b0565b604051808260ff1660ff16815260200191505060405180910390f35b610b4d60048036036020811015610b3457600080fd5b81019080803560ff1690602001909291905050506129b5565b6040518082815260200191505060405180910390f35b610b8f60048036036020811015610b7957600080fd5b81019080803590602001909291905050506129df565b604051808460ff1660ff168152602001838152602001828152602001935050505060405180910390f35b610be560048036036020811015610bcf57600080fd5b8101908080359060200190929190505050612a09565b6040518082815260200191505060405180910390f35b610c03612a1a565b604051808260ff1660ff16815260200191505060405180910390f35b610c5560048036036040811015610c3557600080fd5b810190808035906020019092919080359060200190929190505050612a1f565b6040518082815260200191505060405180910390f35b610c73612ac9565b604051808260ff1660ff16815260200191505060405180910390f35b610cbb60048036036020811015610ca557600080fd5b8101908080359060200190929190505050612ace565b6040518082815260200191505060405180910390f35b610cd9612ae8565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b610d7460048036036060811015610d3157600080fd5b81019080803560ff16906020019092919080359060200190929190803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050612aed565b005b610d7e612e1e565b604051808260ff1660ff16815260200191505060405180910390f35b610da2612e23565b604051808260ff1660ff16815260200191505060405180910390f35b610dea60048036036020811015610dd457600080fd5b8101908080359060200190929190505050612e28565b6040518082815260200191505060405180910390f35b610e08612e39565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b610e52612e64565b604051808260ff1660ff16815260200191505060405180910390f35b610e9a60048036036020811015610e8457600080fd5b8101908080359060200190929190505050612e69565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b610f0860048036036020811015610ef257600080fd5b8101908080359060200190929190505050612e93565b6040518082815260200191505060405180910390f35b610f4a60048036036020811015610f3457600080fd5b8101908080359060200190929190505050612ea5565b6040518082815260200191505060405180910390f35b610f9960048036036040811015610f7657600080fd5b8101908080359060200190929190803560ff169060200190929190505050612ebb565b6040518082815260200191505060405180910390f35b610fdb60048036036020811015610fc557600080fd5b8101908080359060200190929190505050612f67565b604051808260ff1660ff16815260200191505060405180910390f35b6110236004803603602081101561100d57600080fd5b8101908080359060200190929190505050612fae565b6040518082815260200191505060405180910390f35b6110726004803603604081101561104f57600080fd5b8101908080359060200190929190803560ff169060200190929190505050612fc0565b6040518083600560200280838360005b8381101561109d578082015181840152602081019050611082565b5050505090500182600460200280838360005b838110156110cb5780820151818401526020810190506110b0565b505050509050019250505060405180910390f35b61110b600480360360208110156110f557600080fd5b8101908080359060200190929190505050613472565b6040518082815260200191505060405180910390f35b611129613484565b604051808260ff1660ff16815260200191505060405180910390f35b6111886004803603606081101561115b57600080fd5b81019080803560ff1690602001909291908035906020019092919080359060200190929190505050613489565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b611203600480360360408110156111e057600080fd5b81019080803560ff1690602001909291908035906020019092919050505061350f565b6040518082815260200191505060405180910390f35b611221613560565b604051808260ff1660ff16815260200191505060405180910390f35b6112696004803603602081101561125357600080fd5b8101908080359060200190929190505050613565565b6040518082815260200191505060405180910390f35b6112ab6004803603602081101561129557600080fd5b8101908080359060200190929190505050613577565b604051808215151515815260200191505060405180910390f35b6112f1600480360360208110156112db57600080fd5b810190808035906020019092919050505061358a565b6040518082815260200191505060405180910390f35b61130f6135a0565b604051808260ff1660ff16815260200191505060405180910390f35b6113576004803603602081101561134157600080fd5b81019080803590602001909291905050506135a5565b6040518082815260200191505060405180910390f35b6113996004803603602081101561138357600080fd5b81019080803590602001909291905050506135bb565b6040518082815260200191505060405180910390f35b6113e8600480360360408110156113c557600080fd5b81019080803560ff169060200190929190803590602001909291905050506135cc565b6040518082815260200191505060405180910390f35b6114346004803603604081101561141457600080fd5b810190808035906020019092919080359060200190929190505050613623565b005b6114796004803603606081101561144c57600080fd5b81019080803560ff1690602001909291908035906020019092919080359060200190929190505050613991565b604051808215151515815260200191505060405180910390f35b61149b6139d5565b604051808260ff1660ff16815260200191505060405180910390f35b6114bf6139da565b6040518082815260200191505060405180910390f35b611501600480360360208110156114eb57600080fd5b81019080803590602001909291905050506139e1565b6040518082815260200191505060405180910390f35b61151f6139f3565b6040518082815260200191505060405180910390f35b61153d6139fa565b604051808260ff1660ff16815260200191505060405180910390f35b6115856004803603602081101561156f57600080fd5b81019080803590602001909291905050506139ff565b604051808215151515815260200191505060405180910390f35b6115a7613a29565b604051808260ff1660ff16815260200191505060405180910390f35b6115ef600480360360208110156115d957600080fd5b8101908080359060200190929190505050613a2e565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b61166a6004803603604081101561164757600080fd5b8101908080359060200190929190803560ff169060200190929190505050613ac3565b604051808215151515815260200191505060405180910390f35b61168c613bd1565b604051808260ff1660ff16815260200191505060405180910390f35b6116d4600480360360208110156116be57600080fd5b8101908080359060200190929190505050613bd6565b604051808215151515815260200191505060405180910390f35b6117d960048036036101e081101561170557600080fd5b810190808060a001906005806020026040519081016040528092919082600560200280828437600081840152601f19601f8201169050808301925050505050509192919290803590602001909291908035906020019092919080608001906004806020026040519081016040528092919082600460200280828437600081840152601f19601f8201169050808301925050505050509192919290803515159060200190929190803515159060200190929190803560ff169060200190929190803560ff169060200190929190505050613c00565b6040518082815260200191505060405180910390f35b6117f761420f565b604051808260ff1660ff16815260200191505060405180910390f35b61186a600480360360a081101561182957600080fd5b810190808035906020019092919080359060200190929190803560ff1690602001909291908035906020019092919080359060200190929190505050614214565b6040518082815260200191505060405180910390f35b6118af6004803603602081101561189657600080fd5b81019080803560ff1690602001909291905050506142e2565b005b6118b96143a4565b604051808260ff1660ff16815260200191505060405180910390f35b6118dd6143a9565b604051808260ff1660ff16815260200191505060405180910390f35b6119286004803603602081101561190f57600080fd5b81019080803560ff1690602001909291905050506143ae565b6040518082815260200191505060405180910390f35b6119746004803603604081101561195457600080fd5b8101908080359060200190929190803590602001909291905050506143db565b6040518082815260200191505060405180910390f35b6119b6600480360360208110156119a057600080fd5b810190808035906020019092919050505061447f565b604051808215151515815260200191505060405180910390f35b6119d8614519565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b611a2261451e565b604051808260ff1660ff16815260200191505060405180910390f35b611a6a60048036036020811015611a5457600080fd5b8101908080359060200190929190505050614523565b604051808a60ff1660ff1681526020018960ff1660ff1681526020018860ff1660ff1681526020018781526020018660ff1660ff168152602001858152602001848152602001838152602001828152602001995050505050505050505060405180910390f35b611afc60048036036020811015611ae657600080fd5b81019080803590602001909291905050506145a8565b604051808215151515815260200191505060405180910390f35b611b4f60048036036040811015611b2c57600080fd5b81019080803560ff16906020019092919080359060200190929190505050614608565b6040518082815260200191505060405180910390f35b611b6d614622565b6040518082815260200191505060405180910390f35b611baf60048036036020811015611b9957600080fd5b8101908080359060200190929190505050614627565b6040518082815260200191505060405180910390f35b611bfb60048036036040811015611bdb57600080fd5b810190808035906020019092919080359060200190929190505050614638565b6040518082815260200191505060405180910390f35b611c3d60048036036020811015611c2757600080fd5b81019080803590602001909291905050506146e8565b6040518082815260200191505060405180910390f35b611c7f60048036036020811015611c6957600080fd5b81019080803590602001909291905050506146fd565b6040518082815260200191505060405180910390f35b611c9d614890565b604051808260ff1660ff16815260200191505060405180910390f35b611ce560048036036020811015611ccf57600080fd5b8101908080359060200190929190505050614895565b6040518082600560200280838360005b83811015611d10578082015181840152602081019050611cf5565b5050505090500191505060405180910390f35b611d4f60048036036020811015611d3957600080fd5b810190808035906020019092919050505061496f565b604051808215151515815260200191505060405180910390f35b611d9560048036036020811015611d7f57600080fd5b8101908080359060200190929190505050614982565b6040518082815260200191505060405180910390f35b611db3614998565b604051808260ff1660ff16815260200191505060405180910390f35b611ec66004803603611a40811015611de657600080fd5b81019080803560ff169060200190929190806119000190600880602002604051908101604052809291906000905b82821015611e6957838261032002016019806020026040519081016040528092919082601960200280828437600081840152601f19601f82011690508083019250505050505081526020019060010190611e14565b5050505091929192908061010001906008806020026040519081016040528092919082600860200280828437600081840152601f19601f82011690508083019250505050505091929192908035906020019092919050505061499d565b6040518082600860200280838360005b83811015611ef1578082015181840152602081019050611ed6565b5050505090500191505060405180910390f35b611f3d60048036036040811015611f1a57600080fd5b81019080803560ff16906020019092919080359060200190929190505050614c78565b6040518082815260200191505060405180910390f35b611f5b614cc1565b005b611f8960048036036020811015611f7357600080fd5b8101908080359060200190929190505050614da3565b6040518084600b60200280838360005b83811015611fb4578082015181840152602081019050611f99565b5050505090500183600a60200280838360005b83811015611fe2578082015181840152602081019050611fc7565b505050509050018260ff1660ff168152602001935050505060405180910390f35b61202f6004803603602081101561201957600080fd5b8101908080359060200190929190505050614f7c565b6040518082815260200191505060405180910390f35b61204d614f8d565b604051808260ff1660ff16815260200191505060405180910390f35b6120b56004803603604081101561207f57600080fd5b8101908080359060200190929190803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050614f92565b005b61216b60048036036102c08110156120ce57600080fd5b81019080806101600190600b806020026040519081016040528092919082600b60200280828437600081840152601f19601f8201169050808301925050505050509192919290806101400190600a806020026040519081016040528092919082600a60200280828437600081840152601f19601f8201169050808301925050505050509192919290803560ff169060200190929190505050615024565b6040518082815260200191505060405180910390f35b6121ad6004803603602081101561219757600080fd5b81019080803590602001909291905050506151dd565b6040518082601960200280838360005b838110156121d85780820151818401526020810190506121bd565b5050505090500191505060405180910390f35b6122176004803603602081101561220157600080fd5b81019080803590602001909291905050506153a6565b6040518082815260200191505060405180910390f35b6122666004803603604081101561224357600080fd5b81019080803560ff169060200190929190803590602001909291905050506153b7565b6040518082815260200191505060405180910390f35b612284615481565b604051808260ff1660ff16815260200191505060405180910390f35b6122cc600480360360208110156122b657600080fd5b8101908080359060200190929190505050615486565b6040518082815260200191505060405180910390f35b6122ea6154c1565b604051808260ff1660ff16815260200191505060405180910390f35b61230e6154c6565b604051808260ff1660ff16815260200191505060405180910390f35b6123326154cb565b604051808260ff1660ff16815260200191505060405180910390f35b6123876004803603604081101561236457600080fd5b81019080803560ff169060200190929190803590602001909291905050506154d0565b6040518082815260200191505060405180910390f35b6123cc600480360360208110156123b357600080fd5b81019080803560ff1690602001909291905050506154ea565b6040518082815260200191505060405180910390f35b600481565b600060ba82901c9050919050565b6000806101e0831161246f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601b8152602001807f696e76616c696420706c617965724372656174696f6e4d6f6e7468000000000081525060200191505060405180910390fd5b60006014858161247b57fe5b066010019050600585901c9450600c810261ffff1684038592509250509250929050565b600080600e60ff168460ff161061251e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252600f8152602001807f77726f6e67206d6174636820646179000000000000000000000000000000000081525060200191505060405180910390fd5b600460ff168360ff161061259a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252600b8152602001807f77726f6e67206d6174636800000000000000000000000000000000000000000081525060200191505060405180910390fd5b600160080360ff168460ff1610156125c3576125b68484615514565b80925081935050506125dc565b6125d36001600803850384615514565b80935081925050505b9250929050565b8061014960016101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b600061263382615581565b60018260ff166019811061264357fe5b600d0201600001805490509050919050565b600481565b6000612665826128be565b61266e83612fae565b61267784613472565b61268085612e93565b61268986613565565b010101019050919050565b600181565b600481565b60006126aa8484614608565b821090509392505050565b60006007606c83901c169050919050565b60006126eb6126e46126df6126da856146fd565b6139e1565b61560c565b4203615626565b9050919050565b600060208460ff161061276d576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260148152602001807f646566656e6365206f7574206f6620626f756e6400000000000000000000000081525060200191505060405180910390fd5b61040083106127e4576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260148152602001807f646566656e6365206f7574206f6620626f756e6400000000000000000000000081525060200191505060405180910390fd5b6310000000821061285d576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260148152602001807f646566656e6365206f7574206f6620626f756e6400000000000000000000000081525060200191505060405180910390fd5b600060268560ff16901b9050601c84901b811790508281179150509392505050565b601081565b6000601260ff168260ff161061289d57600190506128b6565b6128b385858460ff16601260ff168702016126f2565b90505b949350505050565b6000613fff60ba83901c169050919050565b6000806000806128df856129df565b9250925092506000601260ff1682816128f457fe5b04905061290284848361269e565b612974576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252600f8152602001807f696e76616c6964207465616d206964000000000000000000000000000000000081525060200191505060405180910390fd5b60006129818585846126f2565b90506000601260ff16848161299257fe5b0690506129a3888383600080614214565b9650505050505050919050565b602081565b60006129c082615581565b60018260ff16601981106129d057fe5b600d0201600a01549050919050565b6000806000601f602685901c166103ff601c86901c16630fffffff86169250925092509193909250565b60006007607783901c169050919050565b600181565b6000650800000000008210612a9c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601a8152602001807f707265764c6561677565496478206f7574206f6620626f756e6400000000000081525060200191505060405180910390fd5b741ffffffffffc0000000000000000000000000000001983169250607a82901b8317925082905092915050565b600681565b6000612ae1612adc83615486565b614982565b9050919050565b600181565b600060018460ff1660198110612aff57fe5b600d02016000018381548110612b1157fe5b9060005260206000209060050201600401549050612b30848483613991565b612ba2576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601e8152602001807f63616e6e6f74207472616e736665722061206e6f6e2d626f74207465616d000081525060200191505060405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff161415612c45576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252600f8152602001807f696e76616c69642061646472657373000000000000000000000000000000000081525060200191505060405180910390fd5b612c4d615b39565b6000601260ff1690505b601960ff16811015612c87576001828260198110612c7157fe5b6020020181815250508080600101915050612c57565b5060405180604001604052808281526020018473ffffffffffffffffffffffffffffffffffffffff1681525060018660ff1660198110612cc357fe5b600d02016000018581548110612cd557fe5b90600052602060002090600502016003016000848152602001908152602001600020600082015181600001906019612d0e929190615b5c565b5060208201518160190160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555090505060018560ff1660198110612d6957fe5b600d02016000018481548110612d7b57fe5b9060005260206000209060050201600401600081548092919060010191905055506000612da98686856126f2565b90507f77b66eb1e6d2bc131b79be4213ae7f08f29351c01060e10bcc0302278067bf388185604051808381526020018273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019250505060405180910390a1505050505050565b600581565b600081565b6000600f607d83901c169050919050565b600061014960019054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b601081565b600080600080612e78856129df565b925092509250612e89838383613489565b9350505050919050565b6000613fff60e483901c169050919050565b60006507ffffffffff607a83901c169050919050565b600060208260ff1610612f36576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601c8152602001807f63757272656e7453686972744e756d206f7574206f6620626f756e640000000081525060200191505060405180910390fd5b7503e00000000000000000000000000000000000000000198316925060a58260ff16901b8317925082905092915050565b600080600160190390505b60008160ff1610612fa357612f878382613ac3565b15612f955780915050612fa9565b808060019003915050612f72565b50601990505b919050565b6000613fff60c883901c169050919050565b612fc8615b9c565b612fd0615bbe565b612fd8615b9c565b612fe0615b9c565b6000600a8781612fec57fe5b069050600080600060046201552f8b8161300257fe5b068161300a57fe5b06905060048a901c995060038960ff1610156130525760c885600060ff166005811061303257fe5b602002019061ffff16908161ffff168152505060009250600091506132b8565b60088960ff1610156130c657602885600060ff166005811061307057fe5b602002019061ffff16908161ffff168152505060a085600360ff166005811061309557fe5b602002019061ffff16908161ffff16815250506001925060078960ff168b01816130bb57fe5b0660010191506132b7565b600a8960ff1610156131155760a085600260ff16600581106130e457fe5b602002019061ffff16908161ffff16815250506002925060078960ff168b018161310a57fe5b0660010191506132b6565b600c8960ff16101561318957608285600260ff166005811061313357fe5b602002019061ffff16908161ffff1681525050604685600060ff166005811061315857fe5b602002019061ffff16908161ffff16815250506004925060078960ff168b018161317e57fe5b0660010191506132b5565b600e8960ff1610156131fd57608285600260ff16600581106131a757fe5b602002019061ffff16908161ffff1681525050604685600360ff16600581106131cc57fe5b602002019061ffff16908161ffff16815250506005925060078960ff168b01816131f257fe5b0660010191506132b4565b60108960ff1610156132605760a085600060ff166005811061321b57fe5b602002019061ffff16908161ffff1681525050604685600360ff166005811061324057fe5b602002019061ffff16908161ffff168152505060039250600691506132b3565b60a085600060ff166005811061327257fe5b602002019061ffff16908161ffff1681525050604685600360ff166005811061329757fe5b602002019061ffff16908161ffff168152505060039250600391505b5b5b5b5b5b60338a901c9950600080600090505b600560ff168160ff1610156133a5576000878260ff16600581106132e757fe5b602002015161ffff1614156133285760328c8161330057fe5b06888260ff166005811061331057fe5b602002019061ffff16908161ffff1681525050613379565b6064878260ff166005811061333957fe5b602002015160328e8161334857fe5b060261ffff168161335557fe5b04888260ff166005811061336557fe5b602002019061ffff16908161ffff16815250505b60068c901c9b50878160ff166005811061338f57fe5b60200201518201915080806001019150506132c7565b5060fa8161ffff161015613426576000600560ff168260fa0361ffff16816133c957fe5b04905060008090505b60058160ff1610156134235781898260ff16600581106133ee57fe5b602002015101898260ff166005811061340357fe5b602002019061ffff16908161ffff168152505080806001019150506133d2565b50505b8660405180608001604052808760ff1660ff1681526020018660ff1660ff1681526020018560ff1660ff1681526020018460ff1660ff1681525098509850505050505050509250929050565b6000613fff60d683901c169050919050565b600e81565b600061349484615581565b61349e8484615640565b60018460ff16601981106134ae57fe5b600d020160000183815481106134c057fe5b9060005260206000209060050201600301600083815260200190815260200160002060190160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690509392505050565b600061351a83615581565b6135248383615640565b60018360ff166019811061353457fe5b600d0201600001828154811061354657fe5b906000526020600020906005020160000154905092915050565b600381565b6000613fff60f283901c169050919050565b6000600180607684901c16149050919050565b60006507ffffffffff608183901c169050919050565b600081565b60006507ffffffffff60d583901c169050919050565b60006007606f83901c169050919050565b60006135d783615581565b6135e18383615640565b608060ff1660018460ff16601981106135f657fe5b600d0201600001838154811061360857fe5b90600052602060002090600502016000015402905092915050565b61362c826145a8565b801561363d575061363c81613bd6565b5b6136af576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260198152602001807f756e6578697374656e7420706c61796572206f72207465616d0000000000000081525060200191505060405180910390fd5b60006136ba83615486565b9050600081905060006136cc83614982565b905083811415613744576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260208152602001807f63616e6e6f74207472616e7366657220746f206f726967696e616c207465616d81525060200191505060405180910390fd5b61374d816139ff565b158015613760575061375e846139ff565b155b6137b5576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526036815260200180615de86036913960400191505060405180910390fd5b60006137c0846153a6565b905060006137cd86612f67565b9050601960ff168160ff16141561382f576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526028815260200180615e406028913960400191505060405180910390fd5b6138398487614638565b93506138458482612ebb565b935061385184436143db565b935083600080898152602001908152602001600020819055506000806000613878866129df565b9250925092506001808460ff166019811061388f57fe5b600d020160000183815481106138a157fe5b9060005260206000209060050201600301600083815260200190815260200160002060000186601981106138d157fe5b01819055506138df896129df565b8093508194508295505050508960018460ff16601981106138fc57fe5b600d0201600001838154811061390e57fe5b906000526020600020906005020160030160008381526020019081526020016000206000018560ff166019811061394157fe5b01819055507f65a4d4a8a0afb474d2e9465a6a1a41bb88fd04f41152ba070421f1b3771f15d88a88604051808381526020018281526020019250505060405180910390a150505050505050505050565b60008073ffffffffffffffffffffffffffffffffffffffff166139b5858585613489565b73ffffffffffffffffffffffffffffffffffffffff161490509392505050565b601281565b6101475481565b6000613fff60ac83901c169050919050565b6101485481565b601981565b600080600080613a0e856129df565b925092509250613a1f838383613991565b9350505050919050565b600581565b6000613a39826145a8565b613aab576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260118152602001807f756e6578697374656e7420706c6179657200000000000000000000000000000081525060200191505060405180910390fd5b613abc613ab783612ace565b612e69565b9050919050565b600080600080613ad2866129df565b925092509250613ae3838383613991565b15613b39576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602a815260200180615d4b602a913960400191505060405180910390fd5b600060018460ff1660198110613b4b57fe5b600d02016000018381548110613b5d57fe5b906000526020600020906005020160030160008381526020019081526020016000206000018660ff1660198110613b9057fe5b01549050600160120360ff168660ff161115613bc0576000811480613bb55750600181145b945050505050613bcb565b600181149450505050505b92915050565b600581565b600080600080613be5856129df565b925092509250613bf683838361269e565b9350505050919050565b600080600090505b600560ff168160ff161015613cb2576140008a8260ff1660058110613c2957fe5b602002015161ffff1610613ca5576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260128152602001807f736b696c6c206f7574206f6620626f756e64000000000000000000000000000081525060200191505060405180910390fd5b8080600101915050613c08565b50600a86600060ff1660048110613cc557fe5b602002015160ff1610613d40576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260168152602001807f706f74656e7469616c206f7574206f6620626f756e640000000000000000000081525060200191505060405180910390fd5b600686600160ff1660048110613d5257fe5b602002015160ff1610613dcd576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260188152602001807f666f72776172646e657373206f7574206f6620626f756e64000000000000000081525060200191505060405180910390fd5b600886600260ff1660048110613ddf57fe5b602002015160ff1610613e5a576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260188152602001807f6c6566697473686e657373206f7574206f6620626f756e64000000000000000081525060200191505060405180910390fd5b600886600360ff1660048110613e6c57fe5b602002015160ff1610613ee7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601b8152602001807f616767726573736976656e657373206f7574206f6620626f756e64000000000081525060200191505060405180910390fd5b600086600260ff1660048110613ef957fe5b602002015160ff161415613f7857600086600160ff1660048110613f1957fe5b602002015160ff1614613f77576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602b815260200180615d75602b913960400191505060405180910390fd5b5b60088360ff1610613ff1576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601d8152602001807f67616d65734e6f6e53746f7070696e67206f7574206f6620626f756e6400000081525060200191505060405180910390fd5b614000881061404b576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526023815260200180615d286023913960400191505060405180910390fd5b60008711801561406057506508000000000087105b6140d2576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260158152602001807f706c617965724964206f7574206f6620626f756e64000000000000000000000081525060200191505060405180910390fd5b60008090505b600560ff168160ff16101561412457600e600182010260ff166101000361ffff168a8260ff166005811061410857fe5b602002015161ffff16901b8217915080806001019150506140d8565b5060ac88901b81179050608187901b81179050607d86600060ff166004811061414957fe5b602002015160ff16901b81179050607a86600160ff166004811061416957fe5b602002015160ff16901b81179050607786600260ff166004811061418957fe5b602002015160ff16901b811790506076856141a55760006141a8565b60015b60ff16901b811790506075846141bf5760006141c2565b60015b60ff16901b8117905060728360ff16901b81179050606f8260ff16901b81179050606c86600360ff16600481106141f557fe5b602002015160ff16901b8117905098975050505050505050565b600281565b6000808611801561422a57506508000000000086105b61429c576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260158152602001807f706c617965724964206f7574206f6620626f756e64000000000000000000000081525060200191505060405180910390fd5b600060d587901b90506142af8187614638565b90506142bb8186612ebb565b90506142c78185612a1f565b90506142d381846143db565b90508091505095945050505050565b6001151561014960009054906101000a900460ff1615151461436c576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260178152602001807f63616e6e6f7420696e697469616c697a6520747769636500000000000000000081525060200191505060405180910390fd5b61437542615626565b61014781905550614385816156b6565b600061014960006101000a81548160ff02191690831515021790555050565b600281565b600881565b60006143b982615581565b60018260ff16601981106143c957fe5b600d0201600001805490509050919050565b60006408000000008210614457576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601a8152602001807f6c61737453616c65426c6f636b206f7574206f6620626f756e6400000000000081525060200191505060405180910390fd5b6f03ffffffff80000000000000000000001983169250605782901b8317925082905092915050565b600061448a826145a8565b6144fc576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260118152602001807f756e6578697374656e7420706c6179657200000000000000000000000000000081525060200191505060405180910390fd5b600080600084815260200190815260200160002054149050919050565b600081565b600781565b6001816019811061453057fe5b600d02016000915090508060010160009054906101000a900460ff16908060060160009054906101000a900460ff16908060060160019054906101000a900460ff16908060070154908060080160009054906101000a900460ff169080600901549080600a01549080600b01549080600c0154905089565b6000808214156145bb5760009050614603565b600080600084815260200190815260200160002054146145de5760019050614603565b60008060006145ec856129df565b9250925092506145fd83838361580c565b93505050505b919050565b6000600860ff1661461984846154d0565b02905092915050565b600181565b60006007607a83901c169050919050565b60006508000000000082106146b5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601b8152602001807f63757272656e745465616d496478206f7574206f6620626f756e64000000000081525060200191505060405180910390fd5b7a1ffffffffffc000000000000000000000000000000000000000000198316925060aa82901b8317925082905092915050565b60006407ffffffff605783901c169050919050565b60008060008061470c856129df565b9250925092506000601260ff16828161472157fe5b0490506000601260ff16838161473357fe5b0690506000608060ff16838161474557fe5b04905061475386868561269e565b6147c5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252600f8152602001807f696e76616c6964207465616d206964000000000000000000000000000000000081525060200191505060405180910390fd5b600086868585604051602001808560ff1660ff1681526020018481526020018381526020018260ff1660ff1681526020019450505050506040516020818303038152906040528051906020012060001c90506000601e601060018a60ff166019811061482d57fe5b600d0201600001898154811061483f57fe5b906000526020600020906005020160020160008681526020019081526020016000205402601e6101475402018161487257fe5b0490506148818285838d615829565b98505050505050505050919050565b601081565b61489d615b9c565b6148a682613565565b816000600581106148b357fe5b602002019061ffff16908161ffff16815250506148cf82612e93565b816001600581106148dc57fe5b602002019061ffff16908161ffff16815250506148f882613472565b8160026005811061490557fe5b602002019061ffff16908161ffff168152505061492182612fae565b8160036005811061492e57fe5b602002019061ffff16908161ffff168152505061494a826128be565b8160046005811061495757fe5b602002019061ffff16908161ffff1681525050919050565b6000600180607584901c16149050919050565b60006507ffffffffff60aa83901c169050919050565b600381565b6149a5615be0565b6149ad615c03565b60008060008090505b600460ff168160ff161015614c6c576149cf898261249f565b809350819450505060008682604051602001808381526020018260ff1660ff168152602001925050506040516020818303038152906040528051906020012060001c9050614a1b615c25565b60405180604001604052808a8760ff1660088110614a3557fe5b602002015181526020018a8660ff1660088110614a4e57fe5b60200201518152509050614a60615c47565b60405180604001604052808c8860ff1660088110614a7a57fe5b602002015181526020018c8760ff1660088110614a9357fe5b6020020151815250905061014960019054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16633d3b4ddc8483856000806040518663ffffffff1660e01b8152600401808681526020018560026000925b81841015614b485782846020020151601960200280838360005b83811015614b37578082015181840152602081019050614b1c565b505050509050019260010192614b02565b9250505084600260200280838360005b83811015614b73578082015181840152602081019050614b58565b50505050905001831515151581526020018215151515815260200195505050505050604080518083038186803b158015614bac57600080fd5b505afa158015614bc0573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052506040811015614be557600080fd5b8101908091905050965086600060028110614bfc57fe5b6020020151886002860260ff1660088110614c1357fe5b602002019060ff16908160ff168152505086600160028110614c3157fe5b6020020151886001600287020160ff1660088110614c4b57fe5b602002019060ff16908160ff168152505050505080806001019150506149b6565b50505050949350505050565b6000614c8383615581565b8160018460ff1660198110614c9457fe5b600d0201600b01819055504260018460ff1660198110614cb057fe5b600d02016009018190555092915050565b6001151561014960009054906101000a900460ff16151514614d4b576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260178152602001807f63616e6e6f7420696e697469616c697a6520747769636500000000000000000081525060200191505060405180910390fd5b614d5442615626565b610147819055506000600190505b60198160ff161015614d8457614d77816156b6565b8080600101915050614d62565b50600061014960006101000a81548160ff021916908315150217905550565b614dab615c75565b614db3615c98565b60006a02000000000000000000008410614e35576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601e8152602001807f7461637469637349642073686f756c642066697420696e20363120626974000081525060200191505060405180910390fd5b603f84169050600684901c935060008090505b600a8160ff161015614e9a57600180861614614e65576000614e68565b60015b838260ff16600a8110614e7757fe5b602002019015159081151581525050600185901c94508080600101915050614e48565b5060008090505b600b8160ff161015614f7457601f8516848260ff16600b8110614ec057fe5b602002019060ff16908160ff1681525050601960ff16848260ff16600b8110614ee557fe5b602002015160ff1610614f60576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260168152602001807f696e636f7272656374206c696e65757020656e7472790000000000000000000081525060200191505060405180910390fd5b600585901c94508080600101915050614ea1565b509193909250565b60006007607283901c169050919050565b600381565b6000806000614fa0856129df565b925092509250614fb283838387615882565b7f77b66eb1e6d2bc131b79be4213ae7f08f29351c01060e10bcc0302278067bf388585604051808381526020018273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019250505060405180910390a15050505050565b600060408260ff161061509f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601d8152602001807f7461637469637349642073686f756c642066697420696e20362062697400000081525060200191505060405180910390fd5b60008260ff16905060008090505b600a8160ff1610156150fb578060010260060160ff16858260ff16600a81106150d257fe5b60200201516150e25760006150e5565b60015b60ff16901b8217915080806001019150506150ad565b5060008090505b600b8160ff1610156151d157601960ff16868260ff16600b811061512257fe5b602002015160ff161061519d576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260168152602001807f696e636f7272656374206c696e65757020656e7472790000000000000000000081525060200191505060405180910390fd5b8060050260100160ff16868260ff16600b81106151b657fe5b602002015160ff16901b821791508080600101915050615102565b50809150509392505050565b6151e5615b39565b60008060006151f3856129df565b92509250925061520483838361269e565b615276576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252600f8152602001807f696e76616c6964207465616d206964000000000000000000000000000000000081525060200191505060405180910390fd5b615281838383613991565b156152d25760008090505b601960ff168160ff1610156152cc576152a784848484612884565b858260ff16601981106152b657fe5b602002018181525050808060010191505061528c565b5061539e565b60008090505b601960ff168160ff16101561539c57600060018560ff16601981106152f957fe5b600d0201600001848154811061530b57fe5b906000526020600020906005020160030160008481526020019081526020016000206000018260ff166019811061533e57fe5b0154905060008114156153745761535785858585612884565b868360ff166019811061536657fe5b60200201818152505061538e565b80868360ff166019811061538457fe5b6020020181815250505b5080806001019150506152d8565b505b505050919050565b6000601f60a583901c169050919050565b60008160018460ff16601981106153ca57fe5b600d020160040160018560ff16601981106153e157fe5b600d020160060160019054906101000a900460ff1660ff166002811061540357fe5b018190555060018360ff166019811061541857fe5b600d020160060160019054906101000a900460ff1660010360018460ff166019811061544057fe5b600d020160060160016101000a81548160ff021916908360ff1602179055504260018460ff166019811061547057fe5b600d0201600a018190555092915050565b600481565b60006154918261447f565b156154a65761549f826128d0565b90506154bc565b6000808381526020019081526020016000205490505b919050565b608081565b600281565b600181565b6000601060ff166154e1848461350f565b02905092915050565b60006154f582615581565b60018260ff166019811061550557fe5b600d0201600901549050919050565b600080600080905060008460ff16111561553957615536858560080301615b11565b90505b6000615549866001870101615b11565b9050600060028760ff168161555a57fe5b0660ff16141561557157818193509350505061557a565b80829350935050505b9250929050565b60008160ff16118015615597575060198160ff16105b615609576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260178152602001807f74696d655a6f6e6520646f6573206e6f7420657869737400000000000000000081525060200191505060405180910390fd5b50565b6000600c6301e1338083028161561e57fe5b049050919050565b60006301e13380600c83028161563857fe5b049050919050565b60018260ff166019811061565057fe5b600d02016000018054905081106156b2576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526027815260200180615da06027913960400191505060405180910390fd5b5050565b6156be615cbb565b600181600001818152505060018260ff16601981106156d957fe5b600d020160000181908060018154018082558091505090600182039060005260206000209060050201600090919290919091506000820151816000015560208201518160010160006101000a81548160ff021916908360ff160217905550604082015181600401555050506001808360ff166019811061575557fe5b600d020160000160008154811061576857fe5b90600052602060002090600502016002016000808152602001908152602001600020819055506000801b60018360ff16601981106157a257fe5b600d02016002016000600281106157b557fe5b01819055507fc5d195855a200aa90e2052bcc795cedbc84c2a26556b1d5113b5a30c96003a0b82600080604051808460ff1660ff168152602001838152602001828152602001935050505060405180910390a15050565b6000601260ff1661581d8585614608565b02821090509392505050565b60008061583686856123f5565b8161ffff169150809750819250505061584d615b9c565b615855615bbe565b61585f8888612fc0565b9150915061587582848784600080600080613c00565b9350505050949350505050565b61588b84615581565b6158958484615640565b6158a0848484613991565b15615913576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601e8152602001807f63616e6e6f74207472616e736665722061206e6f6e2d626f74207465616d000081525060200191505060405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161415615999576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526021815260200180615dc76021913960400191505060405180910390fd5b8073ffffffffffffffffffffffffffffffffffffffff1660018560ff16601981106159c057fe5b600d020160000184815481106159d257fe5b9060005260206000209060050201600301600084815260200190815260200160002060190160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff161415615a85576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526022815260200180615e1e6022913960400191505060405180910390fd5b8060018560ff1660198110615a9657fe5b600d02016000018481548110615aa857fe5b9060005260206000209060050201600301600084815260200190815260200160002060190160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050505050565b6000600860ff168260ff161015615b2a57819050615b34565b6001600803820390505b919050565b604051806103200160405280601990602082028038833980820191505090505090565b8260198101928215615b8b579160200282015b82811115615b8a578251825591602001919060010190615b6f565b5b509050615b989190615cdf565b5090565b6040518060a00160405280600590602082028038833980820191505090505090565b6040518060800160405280600490602082028038833980820191505090505090565b604051806101000160405280600890602082028038833980820191505090505090565b6040518060400160405280600290602082028038833980820191505090505090565b6040518060400160405280600290602082028038833980820191505090505090565b6040518061064001604052806002905b615c5f615d04565b815260200190600190039081615c575790505090565b604051806101600160405280600b90602082028038833980820191505090505090565b604051806101400160405280600a90602082028038833980820191505090505090565b604051806060016040528060008152602001600060ff168152602001600081525090565b615d0191905b80821115615cfd576000816000905550600101615ce5565b5090565b90565b60405180610320016040528060199060208202803883398082019150509050509056fe6d6f6e74684f664269727468496e556e697854696d65206f7574206f6620626f756e6463616e6e6f742071756572792061626f757420746865207368697274206f66206120426f74205465616d6c6566746973686e65732063616e206f6e6c79206265207a65726f20666f7220676f616c6b656570657273636f756e74727920646f6573206e6f7420657869737420696e20746869732074696d655a6f6e6563616e6e6f74207472616e7366657220746f2061206e756c6c206164647265737363616e6e6f74207472616e7366657220706c61796572207768656e206174206c65617374206f6e65207465616d206973206120626f74627579657220616e642073656c6c657220617265207468652073616d652061646472746172676574207465616d20666f72207472616e7366657220697320616c72656164792066756c6ca165627a7a72305820782ce876817582e8630cf3b861c3776ac136b8c8b3b9ca4e2507e20b80917cce0029`

// DeployLeagues deploys a new Ethereum contract, binding an instance of Leagues to it.
func DeployLeagues(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Leagues, error) {
	parsed, err := abi.JSON(strings.NewReader(LeaguesABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(LeaguesBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Leagues{LeaguesCaller: LeaguesCaller{contract: contract}, LeaguesTransactor: LeaguesTransactor{contract: contract}, LeaguesFilterer: LeaguesFilterer{contract: contract}}, nil
}

// Leagues is an auto generated Go binding around an Ethereum contract.
type Leagues struct {
	LeaguesCaller     // Read-only binding to the contract
	LeaguesTransactor // Write-only binding to the contract
	LeaguesFilterer   // Log filterer for contract events
}

// LeaguesCaller is an auto generated read-only Go binding around an Ethereum contract.
type LeaguesCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LeaguesTransactor is an auto generated write-only Go binding around an Ethereum contract.
type LeaguesTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LeaguesFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type LeaguesFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LeaguesSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type LeaguesSession struct {
	Contract     *Leagues          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// LeaguesCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type LeaguesCallerSession struct {
	Contract *LeaguesCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// LeaguesTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type LeaguesTransactorSession struct {
	Contract     *LeaguesTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// LeaguesRaw is an auto generated low-level Go binding around an Ethereum contract.
type LeaguesRaw struct {
	Contract *Leagues // Generic contract binding to access the raw methods on
}

// LeaguesCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type LeaguesCallerRaw struct {
	Contract *LeaguesCaller // Generic read-only contract binding to access the raw methods on
}

// LeaguesTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type LeaguesTransactorRaw struct {
	Contract *LeaguesTransactor // Generic write-only contract binding to access the raw methods on
}

// NewLeagues creates a new instance of Leagues, bound to a specific deployed contract.
func NewLeagues(address common.Address, backend bind.ContractBackend) (*Leagues, error) {
	contract, err := bindLeagues(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Leagues{LeaguesCaller: LeaguesCaller{contract: contract}, LeaguesTransactor: LeaguesTransactor{contract: contract}, LeaguesFilterer: LeaguesFilterer{contract: contract}}, nil
}

// NewLeaguesCaller creates a new read-only instance of Leagues, bound to a specific deployed contract.
func NewLeaguesCaller(address common.Address, caller bind.ContractCaller) (*LeaguesCaller, error) {
	contract, err := bindLeagues(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LeaguesCaller{contract: contract}, nil
}

// NewLeaguesTransactor creates a new write-only instance of Leagues, bound to a specific deployed contract.
func NewLeaguesTransactor(address common.Address, transactor bind.ContractTransactor) (*LeaguesTransactor, error) {
	contract, err := bindLeagues(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LeaguesTransactor{contract: contract}, nil
}

// NewLeaguesFilterer creates a new log filterer instance of Leagues, bound to a specific deployed contract.
func NewLeaguesFilterer(address common.Address, filterer bind.ContractFilterer) (*LeaguesFilterer, error) {
	contract, err := bindLeagues(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LeaguesFilterer{contract: contract}, nil
}

// bindLeagues binds a generic wrapper to an already deployed contract.
func bindLeagues(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(LeaguesABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Leagues *LeaguesRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Leagues.Contract.LeaguesCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Leagues *LeaguesRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Leagues.Contract.LeaguesTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Leagues *LeaguesRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Leagues.Contract.LeaguesTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Leagues *LeaguesCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Leagues.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Leagues *LeaguesTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Leagues.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Leagues *LeaguesTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Leagues.Contract.contract.Transact(opts, method, params...)
}

// DAYSPERROUND is a free data retrieval call binding the contract method 0x21ff8ae8.
//
// Solidity: function DAYS_PER_ROUND() constant returns(uint256)
func (_Leagues *LeaguesCaller) DAYSPERROUND(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "DAYS_PER_ROUND")
	return *ret0, err
}

// DAYSPERROUND is a free data retrieval call binding the contract method 0x21ff8ae8.
//
// Solidity: function DAYS_PER_ROUND() constant returns(uint256)
func (_Leagues *LeaguesSession) DAYSPERROUND() (*big.Int, error) {
	return _Leagues.Contract.DAYSPERROUND(&_Leagues.CallOpts)
}

// DAYSPERROUND is a free data retrieval call binding the contract method 0x21ff8ae8.
//
// Solidity: function DAYS_PER_ROUND() constant returns(uint256)
func (_Leagues *LeaguesCallerSession) DAYSPERROUND() (*big.Int, error) {
	return _Leagues.Contract.DAYSPERROUND(&_Leagues.CallOpts)
}

// FREEVERSE is a free data retrieval call binding the contract method 0x39644f21.
//
// Solidity: function FREEVERSE() constant returns(address)
func (_Leagues *LeaguesCaller) FREEVERSE(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "FREEVERSE")
	return *ret0, err
}

// FREEVERSE is a free data retrieval call binding the contract method 0x39644f21.
//
// Solidity: function FREEVERSE() constant returns(address)
func (_Leagues *LeaguesSession) FREEVERSE() (common.Address, error) {
	return _Leagues.Contract.FREEVERSE(&_Leagues.CallOpts)
}

// FREEVERSE is a free data retrieval call binding the contract method 0x39644f21.
//
// Solidity: function FREEVERSE() constant returns(address)
func (_Leagues *LeaguesCallerSession) FREEVERSE() (common.Address, error) {
	return _Leagues.Contract.FREEVERSE(&_Leagues.CallOpts)
}

// FREEPLAYERID is a free data retrieval call binding the contract method 0xc258012b.
//
// Solidity: function FREE_PLAYER_ID() constant returns(uint256)
func (_Leagues *LeaguesCaller) FREEPLAYERID(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "FREE_PLAYER_ID")
	return *ret0, err
}

// FREEPLAYERID is a free data retrieval call binding the contract method 0xc258012b.
//
// Solidity: function FREE_PLAYER_ID() constant returns(uint256)
func (_Leagues *LeaguesSession) FREEPLAYERID() (*big.Int, error) {
	return _Leagues.Contract.FREEPLAYERID(&_Leagues.CallOpts)
}

// FREEPLAYERID is a free data retrieval call binding the contract method 0xc258012b.
//
// Solidity: function FREE_PLAYER_ID() constant returns(uint256)
func (_Leagues *LeaguesCallerSession) FREEPLAYERID() (*big.Int, error) {
	return _Leagues.Contract.FREEPLAYERID(&_Leagues.CallOpts)
}

// IDXC is a free data retrieval call binding the contract method 0xf305a21c.
//
// Solidity: function IDX_C() constant returns(uint8)
func (_Leagues *LeaguesCaller) IDXC(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "IDX_C")
	return *ret0, err
}

// IDXC is a free data retrieval call binding the contract method 0xf305a21c.
//
// Solidity: function IDX_C() constant returns(uint8)
func (_Leagues *LeaguesSession) IDXC() (uint8, error) {
	return _Leagues.Contract.IDXC(&_Leagues.CallOpts)
}

// IDXC is a free data retrieval call binding the contract method 0xf305a21c.
//
// Solidity: function IDX_C() constant returns(uint8)
func (_Leagues *LeaguesCallerSession) IDXC() (uint8, error) {
	return _Leagues.Contract.IDXC(&_Leagues.CallOpts)
}

// IDXCR is a free data retrieval call binding the contract method 0x5becd999.
//
// Solidity: function IDX_CR() constant returns(uint8)
func (_Leagues *LeaguesCaller) IDXCR(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "IDX_CR")
	return *ret0, err
}

// IDXCR is a free data retrieval call binding the contract method 0x5becd999.
//
// Solidity: function IDX_CR() constant returns(uint8)
func (_Leagues *LeaguesSession) IDXCR() (uint8, error) {
	return _Leagues.Contract.IDXCR(&_Leagues.CallOpts)
}

// IDXCR is a free data retrieval call binding the contract method 0x5becd999.
//
// Solidity: function IDX_CR() constant returns(uint8)
func (_Leagues *LeaguesCallerSession) IDXCR() (uint8, error) {
	return _Leagues.Contract.IDXCR(&_Leagues.CallOpts)
}

// IDXD is a free data retrieval call binding the contract method 0x369151db.
//
// Solidity: function IDX_D() constant returns(uint8)
func (_Leagues *LeaguesCaller) IDXD(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "IDX_D")
	return *ret0, err
}

// IDXD is a free data retrieval call binding the contract method 0x369151db.
//
// Solidity: function IDX_D() constant returns(uint8)
func (_Leagues *LeaguesSession) IDXD() (uint8, error) {
	return _Leagues.Contract.IDXD(&_Leagues.CallOpts)
}

// IDXD is a free data retrieval call binding the contract method 0x369151db.
//
// Solidity: function IDX_D() constant returns(uint8)
func (_Leagues *LeaguesCallerSession) IDXD() (uint8, error) {
	return _Leagues.Contract.IDXD(&_Leagues.CallOpts)
}

// IDXF is a free data retrieval call binding the contract method 0xd7b63a11.
//
// Solidity: function IDX_F() constant returns(uint8)
func (_Leagues *LeaguesCaller) IDXF(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "IDX_F")
	return *ret0, err
}

// IDXF is a free data retrieval call binding the contract method 0xd7b63a11.
//
// Solidity: function IDX_F() constant returns(uint8)
func (_Leagues *LeaguesSession) IDXF() (uint8, error) {
	return _Leagues.Contract.IDXF(&_Leagues.CallOpts)
}

// IDXF is a free data retrieval call binding the contract method 0xd7b63a11.
//
// Solidity: function IDX_F() constant returns(uint8)
func (_Leagues *LeaguesCallerSession) IDXF() (uint8, error) {
	return _Leagues.Contract.IDXF(&_Leagues.CallOpts)
}

// IDXGK is a free data retrieval call binding the contract method 0x7420a606.
//
// Solidity: function IDX_GK() constant returns(uint8)
func (_Leagues *LeaguesCaller) IDXGK(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "IDX_GK")
	return *ret0, err
}

// IDXGK is a free data retrieval call binding the contract method 0x7420a606.
//
// Solidity: function IDX_GK() constant returns(uint8)
func (_Leagues *LeaguesSession) IDXGK() (uint8, error) {
	return _Leagues.Contract.IDXGK(&_Leagues.CallOpts)
}

// IDXGK is a free data retrieval call binding the contract method 0x7420a606.
//
// Solidity: function IDX_GK() constant returns(uint8)
func (_Leagues *LeaguesCallerSession) IDXGK() (uint8, error) {
	return _Leagues.Contract.IDXGK(&_Leagues.CallOpts)
}

// IDXL is a free data retrieval call binding the contract method 0xec71bc82.
//
// Solidity: function IDX_L() constant returns(uint8)
func (_Leagues *LeaguesCaller) IDXL(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "IDX_L")
	return *ret0, err
}

// IDXL is a free data retrieval call binding the contract method 0xec71bc82.
//
// Solidity: function IDX_L() constant returns(uint8)
func (_Leagues *LeaguesSession) IDXL() (uint8, error) {
	return _Leagues.Contract.IDXL(&_Leagues.CallOpts)
}

// IDXL is a free data retrieval call binding the contract method 0xec71bc82.
//
// Solidity: function IDX_L() constant returns(uint8)
func (_Leagues *LeaguesCallerSession) IDXL() (uint8, error) {
	return _Leagues.Contract.IDXL(&_Leagues.CallOpts)
}

// IDXLC is a free data retrieval call binding the contract method 0x37fd56af.
//
// Solidity: function IDX_LC() constant returns(uint8)
func (_Leagues *LeaguesCaller) IDXLC(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "IDX_LC")
	return *ret0, err
}

// IDXLC is a free data retrieval call binding the contract method 0x37fd56af.
//
// Solidity: function IDX_LC() constant returns(uint8)
func (_Leagues *LeaguesSession) IDXLC() (uint8, error) {
	return _Leagues.Contract.IDXLC(&_Leagues.CallOpts)
}

// IDXLC is a free data retrieval call binding the contract method 0x37fd56af.
//
// Solidity: function IDX_LC() constant returns(uint8)
func (_Leagues *LeaguesCallerSession) IDXLC() (uint8, error) {
	return _Leagues.Contract.IDXLC(&_Leagues.CallOpts)
}

// IDXLCR is a free data retrieval call binding the contract method 0xb9627097.
//
// Solidity: function IDX_LCR() constant returns(uint8)
func (_Leagues *LeaguesCaller) IDXLCR(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "IDX_LCR")
	return *ret0, err
}

// IDXLCR is a free data retrieval call binding the contract method 0xb9627097.
//
// Solidity: function IDX_LCR() constant returns(uint8)
func (_Leagues *LeaguesSession) IDXLCR() (uint8, error) {
	return _Leagues.Contract.IDXLCR(&_Leagues.CallOpts)
}

// IDXLCR is a free data retrieval call binding the contract method 0xb9627097.
//
// Solidity: function IDX_LCR() constant returns(uint8)
func (_Leagues *LeaguesCallerSession) IDXLCR() (uint8, error) {
	return _Leagues.Contract.IDXLCR(&_Leagues.CallOpts)
}

// IDXLR is a free data retrieval call binding the contract method 0x3d085f96.
//
// Solidity: function IDX_LR() constant returns(uint8)
func (_Leagues *LeaguesCaller) IDXLR(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "IDX_LR")
	return *ret0, err
}

// IDXLR is a free data retrieval call binding the contract method 0x3d085f96.
//
// Solidity: function IDX_LR() constant returns(uint8)
func (_Leagues *LeaguesSession) IDXLR() (uint8, error) {
	return _Leagues.Contract.IDXLR(&_Leagues.CallOpts)
}

// IDXLR is a free data retrieval call binding the contract method 0x3d085f96.
//
// Solidity: function IDX_LR() constant returns(uint8)
func (_Leagues *LeaguesCallerSession) IDXLR() (uint8, error) {
	return _Leagues.Contract.IDXLR(&_Leagues.CallOpts)
}

// IDXM is a free data retrieval call binding the contract method 0x9cc62340.
//
// Solidity: function IDX_M() constant returns(uint8)
func (_Leagues *LeaguesCaller) IDXM(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "IDX_M")
	return *ret0, err
}

// IDXM is a free data retrieval call binding the contract method 0x9cc62340.
//
// Solidity: function IDX_M() constant returns(uint8)
func (_Leagues *LeaguesSession) IDXM() (uint8, error) {
	return _Leagues.Contract.IDXM(&_Leagues.CallOpts)
}

// IDXM is a free data retrieval call binding the contract method 0x9cc62340.
//
// Solidity: function IDX_M() constant returns(uint8)
func (_Leagues *LeaguesCallerSession) IDXM() (uint8, error) {
	return _Leagues.Contract.IDXM(&_Leagues.CallOpts)
}

// IDXMD is a free data retrieval call binding the contract method 0x003e3223.
//
// Solidity: function IDX_MD() constant returns(uint8)
func (_Leagues *LeaguesCaller) IDXMD(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "IDX_MD")
	return *ret0, err
}

// IDXMD is a free data retrieval call binding the contract method 0x003e3223.
//
// Solidity: function IDX_MD() constant returns(uint8)
func (_Leagues *LeaguesSession) IDXMD() (uint8, error) {
	return _Leagues.Contract.IDXMD(&_Leagues.CallOpts)
}

// IDXMD is a free data retrieval call binding the contract method 0x003e3223.
//
// Solidity: function IDX_MD() constant returns(uint8)
func (_Leagues *LeaguesCallerSession) IDXMD() (uint8, error) {
	return _Leagues.Contract.IDXMD(&_Leagues.CallOpts)
}

// IDXMF is a free data retrieval call binding the contract method 0x8f3db436.
//
// Solidity: function IDX_MF() constant returns(uint8)
func (_Leagues *LeaguesCaller) IDXMF(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "IDX_MF")
	return *ret0, err
}

// IDXMF is a free data retrieval call binding the contract method 0x8f3db436.
//
// Solidity: function IDX_MF() constant returns(uint8)
func (_Leagues *LeaguesSession) IDXMF() (uint8, error) {
	return _Leagues.Contract.IDXMF(&_Leagues.CallOpts)
}

// IDXMF is a free data retrieval call binding the contract method 0x8f3db436.
//
// Solidity: function IDX_MF() constant returns(uint8)
func (_Leagues *LeaguesCallerSession) IDXMF() (uint8, error) {
	return _Leagues.Contract.IDXMF(&_Leagues.CallOpts)
}

// IDXR is a free data retrieval call binding the contract method 0x169d2914.
//
// Solidity: function IDX_R() constant returns(uint8)
func (_Leagues *LeaguesCaller) IDXR(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "IDX_R")
	return *ret0, err
}

// IDXR is a free data retrieval call binding the contract method 0x169d2914.
//
// Solidity: function IDX_R() constant returns(uint8)
func (_Leagues *LeaguesSession) IDXR() (uint8, error) {
	return _Leagues.Contract.IDXR(&_Leagues.CallOpts)
}

// IDXR is a free data retrieval call binding the contract method 0x169d2914.
//
// Solidity: function IDX_R() constant returns(uint8)
func (_Leagues *LeaguesCallerSession) IDXR() (uint8, error) {
	return _Leagues.Contract.IDXR(&_Leagues.CallOpts)
}

// LEAGUESPERDIV is a free data retrieval call binding the contract method 0x48d1e9c0.
//
// Solidity: function LEAGUES_PER_DIV() constant returns(uint8)
func (_Leagues *LeaguesCaller) LEAGUESPERDIV(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "LEAGUES_PER_DIV")
	return *ret0, err
}

// LEAGUESPERDIV is a free data retrieval call binding the contract method 0x48d1e9c0.
//
// Solidity: function LEAGUES_PER_DIV() constant returns(uint8)
func (_Leagues *LeaguesSession) LEAGUESPERDIV() (uint8, error) {
	return _Leagues.Contract.LEAGUESPERDIV(&_Leagues.CallOpts)
}

// LEAGUESPERDIV is a free data retrieval call binding the contract method 0x48d1e9c0.
//
// Solidity: function LEAGUES_PER_DIV() constant returns(uint8)
func (_Leagues *LeaguesCallerSession) LEAGUESPERDIV() (uint8, error) {
	return _Leagues.Contract.LEAGUESPERDIV(&_Leagues.CallOpts)
}

// MATCHDAYS is a free data retrieval call binding the contract method 0x561b1118.
//
// Solidity: function MATCHDAYS() constant returns(uint8)
func (_Leagues *LeaguesCaller) MATCHDAYS(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "MATCHDAYS")
	return *ret0, err
}

// MATCHDAYS is a free data retrieval call binding the contract method 0x561b1118.
//
// Solidity: function MATCHDAYS() constant returns(uint8)
func (_Leagues *LeaguesSession) MATCHDAYS() (uint8, error) {
	return _Leagues.Contract.MATCHDAYS(&_Leagues.CallOpts)
}

// MATCHDAYS is a free data retrieval call binding the contract method 0x561b1118.
//
// Solidity: function MATCHDAYS() constant returns(uint8)
func (_Leagues *LeaguesCallerSession) MATCHDAYS() (uint8, error) {
	return _Leagues.Contract.MATCHDAYS(&_Leagues.CallOpts)
}

// MATCHESPERDAY is a free data retrieval call binding the contract method 0x0c85696c.
//
// Solidity: function MATCHES_PER_DAY() constant returns(uint8)
func (_Leagues *LeaguesCaller) MATCHESPERDAY(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "MATCHES_PER_DAY")
	return *ret0, err
}

// MATCHESPERDAY is a free data retrieval call binding the contract method 0x0c85696c.
//
// Solidity: function MATCHES_PER_DAY() constant returns(uint8)
func (_Leagues *LeaguesSession) MATCHESPERDAY() (uint8, error) {
	return _Leagues.Contract.MATCHESPERDAY(&_Leagues.CallOpts)
}

// MATCHESPERDAY is a free data retrieval call binding the contract method 0x0c85696c.
//
// Solidity: function MATCHES_PER_DAY() constant returns(uint8)
func (_Leagues *LeaguesCallerSession) MATCHESPERDAY() (uint8, error) {
	return _Leagues.Contract.MATCHESPERDAY(&_Leagues.CallOpts)
}

// MAXPLAYERAGEATBIRTH is a free data retrieval call binding the contract method 0x2a238b0a.
//
// Solidity: function MAX_PLAYER_AGE_AT_BIRTH() constant returns(uint8)
func (_Leagues *LeaguesCaller) MAXPLAYERAGEATBIRTH(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "MAX_PLAYER_AGE_AT_BIRTH")
	return *ret0, err
}

// MAXPLAYERAGEATBIRTH is a free data retrieval call binding the contract method 0x2a238b0a.
//
// Solidity: function MAX_PLAYER_AGE_AT_BIRTH() constant returns(uint8)
func (_Leagues *LeaguesSession) MAXPLAYERAGEATBIRTH() (uint8, error) {
	return _Leagues.Contract.MAXPLAYERAGEATBIRTH(&_Leagues.CallOpts)
}

// MAXPLAYERAGEATBIRTH is a free data retrieval call binding the contract method 0x2a238b0a.
//
// Solidity: function MAX_PLAYER_AGE_AT_BIRTH() constant returns(uint8)
func (_Leagues *LeaguesCallerSession) MAXPLAYERAGEATBIRTH() (uint8, error) {
	return _Leagues.Contract.MAXPLAYERAGEATBIRTH(&_Leagues.CallOpts)
}

// MINPLAYERAGEATBIRTH is a free data retrieval call binding the contract method 0xc79055d4.
//
// Solidity: function MIN_PLAYER_AGE_AT_BIRTH() constant returns(uint8)
func (_Leagues *LeaguesCaller) MINPLAYERAGEATBIRTH(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "MIN_PLAYER_AGE_AT_BIRTH")
	return *ret0, err
}

// MINPLAYERAGEATBIRTH is a free data retrieval call binding the contract method 0xc79055d4.
//
// Solidity: function MIN_PLAYER_AGE_AT_BIRTH() constant returns(uint8)
func (_Leagues *LeaguesSession) MINPLAYERAGEATBIRTH() (uint8, error) {
	return _Leagues.Contract.MINPLAYERAGEATBIRTH(&_Leagues.CallOpts)
}

// MINPLAYERAGEATBIRTH is a free data retrieval call binding the contract method 0xc79055d4.
//
// Solidity: function MIN_PLAYER_AGE_AT_BIRTH() constant returns(uint8)
func (_Leagues *LeaguesCallerSession) MINPLAYERAGEATBIRTH() (uint8, error) {
	return _Leagues.Contract.MINPLAYERAGEATBIRTH(&_Leagues.CallOpts)
}

// NULLADDR is a free data retrieval call binding the contract method 0xb3f390b3.
//
// Solidity: function NULL_ADDR() constant returns(address)
func (_Leagues *LeaguesCaller) NULLADDR(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "NULL_ADDR")
	return *ret0, err
}

// NULLADDR is a free data retrieval call binding the contract method 0xb3f390b3.
//
// Solidity: function NULL_ADDR() constant returns(address)
func (_Leagues *LeaguesSession) NULLADDR() (common.Address, error) {
	return _Leagues.Contract.NULLADDR(&_Leagues.CallOpts)
}

// NULLADDR is a free data retrieval call binding the contract method 0xb3f390b3.
//
// Solidity: function NULL_ADDR() constant returns(address)
func (_Leagues *LeaguesCallerSession) NULLADDR() (common.Address, error) {
	return _Leagues.Contract.NULLADDR(&_Leagues.CallOpts)
}

// NSKILLS is a free data retrieval call binding the contract method 0x976daaac.
//
// Solidity: function N_SKILLS() constant returns(uint8)
func (_Leagues *LeaguesCaller) NSKILLS(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "N_SKILLS")
	return *ret0, err
}

// NSKILLS is a free data retrieval call binding the contract method 0x976daaac.
//
// Solidity: function N_SKILLS() constant returns(uint8)
func (_Leagues *LeaguesSession) NSKILLS() (uint8, error) {
	return _Leagues.Contract.NSKILLS(&_Leagues.CallOpts)
}

// NSKILLS is a free data retrieval call binding the contract method 0x976daaac.
//
// Solidity: function N_SKILLS() constant returns(uint8)
func (_Leagues *LeaguesCallerSession) NSKILLS() (uint8, error) {
	return _Leagues.Contract.NSKILLS(&_Leagues.CallOpts)
}

// PLAYERSPERTEAMINIT is a free data retrieval call binding the contract method 0x83c31d3b.
//
// Solidity: function PLAYERS_PER_TEAM_INIT() constant returns(uint8)
func (_Leagues *LeaguesCaller) PLAYERSPERTEAMINIT(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "PLAYERS_PER_TEAM_INIT")
	return *ret0, err
}

// PLAYERSPERTEAMINIT is a free data retrieval call binding the contract method 0x83c31d3b.
//
// Solidity: function PLAYERS_PER_TEAM_INIT() constant returns(uint8)
func (_Leagues *LeaguesSession) PLAYERSPERTEAMINIT() (uint8, error) {
	return _Leagues.Contract.PLAYERSPERTEAMINIT(&_Leagues.CallOpts)
}

// PLAYERSPERTEAMINIT is a free data retrieval call binding the contract method 0x83c31d3b.
//
// Solidity: function PLAYERS_PER_TEAM_INIT() constant returns(uint8)
func (_Leagues *LeaguesCallerSession) PLAYERSPERTEAMINIT() (uint8, error) {
	return _Leagues.Contract.PLAYERSPERTEAMINIT(&_Leagues.CallOpts)
}

// PLAYERSPERTEAMMAX is a free data retrieval call binding the contract method 0x8adddc9d.
//
// Solidity: function PLAYERS_PER_TEAM_MAX() constant returns(uint8)
func (_Leagues *LeaguesCaller) PLAYERSPERTEAMMAX(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "PLAYERS_PER_TEAM_MAX")
	return *ret0, err
}

// PLAYERSPERTEAMMAX is a free data retrieval call binding the contract method 0x8adddc9d.
//
// Solidity: function PLAYERS_PER_TEAM_MAX() constant returns(uint8)
func (_Leagues *LeaguesSession) PLAYERSPERTEAMMAX() (uint8, error) {
	return _Leagues.Contract.PLAYERSPERTEAMMAX(&_Leagues.CallOpts)
}

// PLAYERSPERTEAMMAX is a free data retrieval call binding the contract method 0x8adddc9d.
//
// Solidity: function PLAYERS_PER_TEAM_MAX() constant returns(uint8)
func (_Leagues *LeaguesCallerSession) PLAYERSPERTEAMMAX() (uint8, error) {
	return _Leagues.Contract.PLAYERSPERTEAMMAX(&_Leagues.CallOpts)
}

// SKDEF is a free data retrieval call binding the contract method 0xe81e21bb.
//
// Solidity: function SK_DEF() constant returns(uint8)
func (_Leagues *LeaguesCaller) SKDEF(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "SK_DEF")
	return *ret0, err
}

// SKDEF is a free data retrieval call binding the contract method 0xe81e21bb.
//
// Solidity: function SK_DEF() constant returns(uint8)
func (_Leagues *LeaguesSession) SKDEF() (uint8, error) {
	return _Leagues.Contract.SKDEF(&_Leagues.CallOpts)
}

// SKDEF is a free data retrieval call binding the contract method 0xe81e21bb.
//
// Solidity: function SK_DEF() constant returns(uint8)
func (_Leagues *LeaguesCallerSession) SKDEF() (uint8, error) {
	return _Leagues.Contract.SKDEF(&_Leagues.CallOpts)
}

// SKEND is a free data retrieval call binding the contract method 0x1884332c.
//
// Solidity: function SK_END() constant returns(uint8)
func (_Leagues *LeaguesCaller) SKEND(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "SK_END")
	return *ret0, err
}

// SKEND is a free data retrieval call binding the contract method 0x1884332c.
//
// Solidity: function SK_END() constant returns(uint8)
func (_Leagues *LeaguesSession) SKEND() (uint8, error) {
	return _Leagues.Contract.SKEND(&_Leagues.CallOpts)
}

// SKEND is a free data retrieval call binding the contract method 0x1884332c.
//
// Solidity: function SK_END() constant returns(uint8)
func (_Leagues *LeaguesCallerSession) SKEND() (uint8, error) {
	return _Leagues.Contract.SKEND(&_Leagues.CallOpts)
}

// SKPAS is a free data retrieval call binding the contract method 0xab1b7c5e.
//
// Solidity: function SK_PAS() constant returns(uint8)
func (_Leagues *LeaguesCaller) SKPAS(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "SK_PAS")
	return *ret0, err
}

// SKPAS is a free data retrieval call binding the contract method 0xab1b7c5e.
//
// Solidity: function SK_PAS() constant returns(uint8)
func (_Leagues *LeaguesSession) SKPAS() (uint8, error) {
	return _Leagues.Contract.SKPAS(&_Leagues.CallOpts)
}

// SKPAS is a free data retrieval call binding the contract method 0xab1b7c5e.
//
// Solidity: function SK_PAS() constant returns(uint8)
func (_Leagues *LeaguesCallerSession) SKPAS() (uint8, error) {
	return _Leagues.Contract.SKPAS(&_Leagues.CallOpts)
}

// SKSHO is a free data retrieval call binding the contract method 0x40cd05fd.
//
// Solidity: function SK_SHO() constant returns(uint8)
func (_Leagues *LeaguesCaller) SKSHO(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "SK_SHO")
	return *ret0, err
}

// SKSHO is a free data retrieval call binding the contract method 0x40cd05fd.
//
// Solidity: function SK_SHO() constant returns(uint8)
func (_Leagues *LeaguesSession) SKSHO() (uint8, error) {
	return _Leagues.Contract.SKSHO(&_Leagues.CallOpts)
}

// SKSHO is a free data retrieval call binding the contract method 0x40cd05fd.
//
// Solidity: function SK_SHO() constant returns(uint8)
func (_Leagues *LeaguesCallerSession) SKSHO() (uint8, error) {
	return _Leagues.Contract.SKSHO(&_Leagues.CallOpts)
}

// SKSPE is a free data retrieval call binding the contract method 0xf8ef7b9e.
//
// Solidity: function SK_SPE() constant returns(uint8)
func (_Leagues *LeaguesCaller) SKSPE(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "SK_SPE")
	return *ret0, err
}

// SKSPE is a free data retrieval call binding the contract method 0xf8ef7b9e.
//
// Solidity: function SK_SPE() constant returns(uint8)
func (_Leagues *LeaguesSession) SKSPE() (uint8, error) {
	return _Leagues.Contract.SKSPE(&_Leagues.CallOpts)
}

// SKSPE is a free data retrieval call binding the contract method 0xf8ef7b9e.
//
// Solidity: function SK_SPE() constant returns(uint8)
func (_Leagues *LeaguesCallerSession) SKSPE() (uint8, error) {
	return _Leagues.Contract.SKSPE(&_Leagues.CallOpts)
}

// TEAMSPERDIVISION is a free data retrieval call binding the contract method 0xf21f5a83.
//
// Solidity: function TEAMS_PER_DIVISION() constant returns(uint8)
func (_Leagues *LeaguesCaller) TEAMSPERDIVISION(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "TEAMS_PER_DIVISION")
	return *ret0, err
}

// TEAMSPERDIVISION is a free data retrieval call binding the contract method 0xf21f5a83.
//
// Solidity: function TEAMS_PER_DIVISION() constant returns(uint8)
func (_Leagues *LeaguesSession) TEAMSPERDIVISION() (uint8, error) {
	return _Leagues.Contract.TEAMSPERDIVISION(&_Leagues.CallOpts)
}

// TEAMSPERDIVISION is a free data retrieval call binding the contract method 0xf21f5a83.
//
// Solidity: function TEAMS_PER_DIVISION() constant returns(uint8)
func (_Leagues *LeaguesCallerSession) TEAMSPERDIVISION() (uint8, error) {
	return _Leagues.Contract.TEAMSPERDIVISION(&_Leagues.CallOpts)
}

// TEAMSPERLEAGUE is a free data retrieval call binding the contract method 0xac5db9ee.
//
// Solidity: function TEAMS_PER_LEAGUE() constant returns(uint8)
func (_Leagues *LeaguesCaller) TEAMSPERLEAGUE(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "TEAMS_PER_LEAGUE")
	return *ret0, err
}

// TEAMSPERLEAGUE is a free data retrieval call binding the contract method 0xac5db9ee.
//
// Solidity: function TEAMS_PER_LEAGUE() constant returns(uint8)
func (_Leagues *LeaguesSession) TEAMSPERLEAGUE() (uint8, error) {
	return _Leagues.Contract.TEAMSPERLEAGUE(&_Leagues.CallOpts)
}

// TEAMSPERLEAGUE is a free data retrieval call binding the contract method 0xac5db9ee.
//
// Solidity: function TEAMS_PER_LEAGUE() constant returns(uint8)
func (_Leagues *LeaguesCallerSession) TEAMSPERLEAGUE() (uint8, error) {
	return _Leagues.Contract.TEAMSPERLEAGUE(&_Leagues.CallOpts)
}

// TeamExistsInCountry is a free data retrieval call binding the contract method 0x1a6daba2.
//
// Solidity: function _teamExistsInCountry(uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) constant returns(bool)
func (_Leagues *LeaguesCaller) TeamExistsInCountry(opts *bind.CallOpts, timeZone uint8, countryIdxInTZ *big.Int, teamIdxInCountry *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "_teamExistsInCountry", timeZone, countryIdxInTZ, teamIdxInCountry)
	return *ret0, err
}

// TeamExistsInCountry is a free data retrieval call binding the contract method 0x1a6daba2.
//
// Solidity: function _teamExistsInCountry(uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) constant returns(bool)
func (_Leagues *LeaguesSession) TeamExistsInCountry(timeZone uint8, countryIdxInTZ *big.Int, teamIdxInCountry *big.Int) (bool, error) {
	return _Leagues.Contract.TeamExistsInCountry(&_Leagues.CallOpts, timeZone, countryIdxInTZ, teamIdxInCountry)
}

// TeamExistsInCountry is a free data retrieval call binding the contract method 0x1a6daba2.
//
// Solidity: function _teamExistsInCountry(uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) constant returns(bool)
func (_Leagues *LeaguesCallerSession) TeamExistsInCountry(timeZone uint8, countryIdxInTZ *big.Int, teamIdxInCountry *big.Int) (bool, error) {
	return _Leagues.Contract.TeamExistsInCountry(&_Leagues.CallOpts, timeZone, countryIdxInTZ, teamIdxInCountry)
}

// TimeZones is a free data retrieval call binding the contract method 0xb96b1a30.
//
// Solidity: function _timeZones(uint256 ) constant returns(uint8 nCountriesToAdd, uint8 newestOrgMapIdx, uint8 newestSkillsIdx, bytes32 scoresRoot, uint8 updateCycleIdx, uint256 lastActionsSubmissionTime, uint256 lastUpdateTime, bytes32 actionsRoot, uint256 lastMarketClosureBlockNum)
func (_Leagues *LeaguesCaller) TimeZones(opts *bind.CallOpts, arg0 *big.Int) (struct {
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
	err := _Leagues.contract.Call(opts, out, "_timeZones", arg0)
	return *ret, err
}

// TimeZones is a free data retrieval call binding the contract method 0xb96b1a30.
//
// Solidity: function _timeZones(uint256 ) constant returns(uint8 nCountriesToAdd, uint8 newestOrgMapIdx, uint8 newestSkillsIdx, bytes32 scoresRoot, uint8 updateCycleIdx, uint256 lastActionsSubmissionTime, uint256 lastUpdateTime, bytes32 actionsRoot, uint256 lastMarketClosureBlockNum)
func (_Leagues *LeaguesSession) TimeZones(arg0 *big.Int) (struct {
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
	return _Leagues.Contract.TimeZones(&_Leagues.CallOpts, arg0)
}

// TimeZones is a free data retrieval call binding the contract method 0xb96b1a30.
//
// Solidity: function _timeZones(uint256 ) constant returns(uint8 nCountriesToAdd, uint8 newestOrgMapIdx, uint8 newestSkillsIdx, bytes32 scoresRoot, uint8 updateCycleIdx, uint256 lastActionsSubmissionTime, uint256 lastUpdateTime, bytes32 actionsRoot, uint256 lastMarketClosureBlockNum)
func (_Leagues *LeaguesCallerSession) TimeZones(arg0 *big.Int) (struct {
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
	return _Leagues.Contract.TimeZones(&_Leagues.CallOpts, arg0)
}

// ComputeBirthMonth is a free data retrieval call binding the contract method 0x00aae8df.
//
// Solidity: function computeBirthMonth(uint256 dna, uint256 playerCreationMonth) constant returns(uint16, uint256)
func (_Leagues *LeaguesCaller) ComputeBirthMonth(opts *bind.CallOpts, dna *big.Int, playerCreationMonth *big.Int) (uint16, *big.Int, error) {
	var (
		ret0 = new(uint16)
		ret1 = new(*big.Int)
	)
	out := &[]interface{}{
		ret0,
		ret1,
	}
	err := _Leagues.contract.Call(opts, out, "computeBirthMonth", dna, playerCreationMonth)
	return *ret0, *ret1, err
}

// ComputeBirthMonth is a free data retrieval call binding the contract method 0x00aae8df.
//
// Solidity: function computeBirthMonth(uint256 dna, uint256 playerCreationMonth) constant returns(uint16, uint256)
func (_Leagues *LeaguesSession) ComputeBirthMonth(dna *big.Int, playerCreationMonth *big.Int) (uint16, *big.Int, error) {
	return _Leagues.Contract.ComputeBirthMonth(&_Leagues.CallOpts, dna, playerCreationMonth)
}

// ComputeBirthMonth is a free data retrieval call binding the contract method 0x00aae8df.
//
// Solidity: function computeBirthMonth(uint256 dna, uint256 playerCreationMonth) constant returns(uint16, uint256)
func (_Leagues *LeaguesCallerSession) ComputeBirthMonth(dna *big.Int, playerCreationMonth *big.Int) (uint16, *big.Int, error) {
	return _Leagues.Contract.ComputeBirthMonth(&_Leagues.CallOpts, dna, playerCreationMonth)
}

// ComputeMatchday is a free data retrieval call binding the contract method 0xd7e4e6d5.
//
// Solidity: function computeMatchday(uint8 matchday, uint256[25][8] prevLeagueState, uint256[8] tacticsIds, uint256 currentVerseSeed) constant returns(uint8[8] scores)
func (_Leagues *LeaguesCaller) ComputeMatchday(opts *bind.CallOpts, matchday uint8, prevLeagueState [8][25]*big.Int, tacticsIds [8]*big.Int, currentVerseSeed *big.Int) ([8]uint8, error) {
	var (
		ret0 = new([8]uint8)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "computeMatchday", matchday, prevLeagueState, tacticsIds, currentVerseSeed)
	return *ret0, err
}

// ComputeMatchday is a free data retrieval call binding the contract method 0xd7e4e6d5.
//
// Solidity: function computeMatchday(uint8 matchday, uint256[25][8] prevLeagueState, uint256[8] tacticsIds, uint256 currentVerseSeed) constant returns(uint8[8] scores)
func (_Leagues *LeaguesSession) ComputeMatchday(matchday uint8, prevLeagueState [8][25]*big.Int, tacticsIds [8]*big.Int, currentVerseSeed *big.Int) ([8]uint8, error) {
	return _Leagues.Contract.ComputeMatchday(&_Leagues.CallOpts, matchday, prevLeagueState, tacticsIds, currentVerseSeed)
}

// ComputeMatchday is a free data retrieval call binding the contract method 0xd7e4e6d5.
//
// Solidity: function computeMatchday(uint8 matchday, uint256[25][8] prevLeagueState, uint256[8] tacticsIds, uint256 currentVerseSeed) constant returns(uint8[8] scores)
func (_Leagues *LeaguesCallerSession) ComputeMatchday(matchday uint8, prevLeagueState [8][25]*big.Int, tacticsIds [8]*big.Int, currentVerseSeed *big.Int) ([8]uint8, error) {
	return _Leagues.Contract.ComputeMatchday(&_Leagues.CallOpts, matchday, prevLeagueState, tacticsIds, currentVerseSeed)
}

// ComputeSkills is a free data retrieval call binding the contract method 0x547d8298.
//
// Solidity: function computeSkills(uint256 dna, uint8 shirtNum) constant returns(uint16[5], uint8[4])
func (_Leagues *LeaguesCaller) ComputeSkills(opts *bind.CallOpts, dna *big.Int, shirtNum uint8) ([5]uint16, [4]uint8, error) {
	var (
		ret0 = new([5]uint16)
		ret1 = new([4]uint8)
	)
	out := &[]interface{}{
		ret0,
		ret1,
	}
	err := _Leagues.contract.Call(opts, out, "computeSkills", dna, shirtNum)
	return *ret0, *ret1, err
}

// ComputeSkills is a free data retrieval call binding the contract method 0x547d8298.
//
// Solidity: function computeSkills(uint256 dna, uint8 shirtNum) constant returns(uint16[5], uint8[4])
func (_Leagues *LeaguesSession) ComputeSkills(dna *big.Int, shirtNum uint8) ([5]uint16, [4]uint8, error) {
	return _Leagues.Contract.ComputeSkills(&_Leagues.CallOpts, dna, shirtNum)
}

// ComputeSkills is a free data retrieval call binding the contract method 0x547d8298.
//
// Solidity: function computeSkills(uint256 dna, uint8 shirtNum) constant returns(uint16[5], uint8[4])
func (_Leagues *LeaguesCallerSession) ComputeSkills(dna *big.Int, shirtNum uint8) ([5]uint16, [4]uint8, error) {
	return _Leagues.Contract.ComputeSkills(&_Leagues.CallOpts, dna, shirtNum)
}

// CountCountries is a free data retrieval call binding the contract method 0x0abcd3e5.
//
// Solidity: function countCountries(uint8 timeZone) constant returns(uint256)
func (_Leagues *LeaguesCaller) CountCountries(opts *bind.CallOpts, timeZone uint8) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "countCountries", timeZone)
	return *ret0, err
}

// CountCountries is a free data retrieval call binding the contract method 0x0abcd3e5.
//
// Solidity: function countCountries(uint8 timeZone) constant returns(uint256)
func (_Leagues *LeaguesSession) CountCountries(timeZone uint8) (*big.Int, error) {
	return _Leagues.Contract.CountCountries(&_Leagues.CallOpts, timeZone)
}

// CountCountries is a free data retrieval call binding the contract method 0x0abcd3e5.
//
// Solidity: function countCountries(uint8 timeZone) constant returns(uint256)
func (_Leagues *LeaguesCallerSession) CountCountries(timeZone uint8) (*big.Int, error) {
	return _Leagues.Contract.CountCountries(&_Leagues.CallOpts, timeZone)
}

// CountTeams is a free data retrieval call binding the contract method 0x7b2566a5.
//
// Solidity: function countTeams(uint8 timeZone, uint256 countryIdxInTZ) constant returns(uint256)
func (_Leagues *LeaguesCaller) CountTeams(opts *bind.CallOpts, timeZone uint8, countryIdxInTZ *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "countTeams", timeZone, countryIdxInTZ)
	return *ret0, err
}

// CountTeams is a free data retrieval call binding the contract method 0x7b2566a5.
//
// Solidity: function countTeams(uint8 timeZone, uint256 countryIdxInTZ) constant returns(uint256)
func (_Leagues *LeaguesSession) CountTeams(timeZone uint8, countryIdxInTZ *big.Int) (*big.Int, error) {
	return _Leagues.Contract.CountTeams(&_Leagues.CallOpts, timeZone, countryIdxInTZ)
}

// CountTeams is a free data retrieval call binding the contract method 0x7b2566a5.
//
// Solidity: function countTeams(uint8 timeZone, uint256 countryIdxInTZ) constant returns(uint256)
func (_Leagues *LeaguesCallerSession) CountTeams(timeZone uint8, countryIdxInTZ *big.Int) (*big.Int, error) {
	return _Leagues.Contract.CountTeams(&_Leagues.CallOpts, timeZone, countryIdxInTZ)
}

// CurrentRound is a free data retrieval call binding the contract method 0x8a19c8bc.
//
// Solidity: function currentRound() constant returns(uint256)
func (_Leagues *LeaguesCaller) CurrentRound(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "currentRound")
	return *ret0, err
}

// CurrentRound is a free data retrieval call binding the contract method 0x8a19c8bc.
//
// Solidity: function currentRound() constant returns(uint256)
func (_Leagues *LeaguesSession) CurrentRound() (*big.Int, error) {
	return _Leagues.Contract.CurrentRound(&_Leagues.CallOpts)
}

// CurrentRound is a free data retrieval call binding the contract method 0x8a19c8bc.
//
// Solidity: function currentRound() constant returns(uint256)
func (_Leagues *LeaguesCallerSession) CurrentRound() (*big.Int, error) {
	return _Leagues.Contract.CurrentRound(&_Leagues.CallOpts)
}

// DecodeTZCountryAndVal is a free data retrieval call binding the contract method 0x3260840b.
//
// Solidity: function decodeTZCountryAndVal(uint256 encoded) constant returns(uint8, uint256, uint256)
func (_Leagues *LeaguesCaller) DecodeTZCountryAndVal(opts *bind.CallOpts, encoded *big.Int) (uint8, *big.Int, *big.Int, error) {
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
	err := _Leagues.contract.Call(opts, out, "decodeTZCountryAndVal", encoded)
	return *ret0, *ret1, *ret2, err
}

// DecodeTZCountryAndVal is a free data retrieval call binding the contract method 0x3260840b.
//
// Solidity: function decodeTZCountryAndVal(uint256 encoded) constant returns(uint8, uint256, uint256)
func (_Leagues *LeaguesSession) DecodeTZCountryAndVal(encoded *big.Int) (uint8, *big.Int, *big.Int, error) {
	return _Leagues.Contract.DecodeTZCountryAndVal(&_Leagues.CallOpts, encoded)
}

// DecodeTZCountryAndVal is a free data retrieval call binding the contract method 0x3260840b.
//
// Solidity: function decodeTZCountryAndVal(uint256 encoded) constant returns(uint8, uint256, uint256)
func (_Leagues *LeaguesCallerSession) DecodeTZCountryAndVal(encoded *big.Int) (uint8, *big.Int, *big.Int, error) {
	return _Leagues.Contract.DecodeTZCountryAndVal(&_Leagues.CallOpts, encoded)
}

// DecodeTactics is a free data retrieval call binding the contract method 0xe6400ac4.
//
// Solidity: function decodeTactics(uint256 tactics) constant returns(uint8[11] lineup, bool[10] extraAttack, uint8 tacticsId)
func (_Leagues *LeaguesCaller) DecodeTactics(opts *bind.CallOpts, tactics *big.Int) (struct {
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
	err := _Leagues.contract.Call(opts, out, "decodeTactics", tactics)
	return *ret, err
}

// DecodeTactics is a free data retrieval call binding the contract method 0xe6400ac4.
//
// Solidity: function decodeTactics(uint256 tactics) constant returns(uint8[11] lineup, bool[10] extraAttack, uint8 tacticsId)
func (_Leagues *LeaguesSession) DecodeTactics(tactics *big.Int) (struct {
	Lineup      [11]uint8
	ExtraAttack [10]bool
	TacticsId   uint8
}, error) {
	return _Leagues.Contract.DecodeTactics(&_Leagues.CallOpts, tactics)
}

// DecodeTactics is a free data retrieval call binding the contract method 0xe6400ac4.
//
// Solidity: function decodeTactics(uint256 tactics) constant returns(uint8[11] lineup, bool[10] extraAttack, uint8 tacticsId)
func (_Leagues *LeaguesCallerSession) DecodeTactics(tactics *big.Int) (struct {
	Lineup      [11]uint8
	ExtraAttack [10]bool
	TacticsId   uint8
}, error) {
	return _Leagues.Contract.DecodeTactics(&_Leagues.CallOpts, tactics)
}

// EncodePlayerSkills is a free data retrieval call binding the contract method 0x9c53e3fd.
//
// Solidity: function encodePlayerSkills(uint16[5] skills, uint256 monthOfBirth, uint256 playerId, uint8[4] birthTraits, bool alignedLastHalf, bool redCardLastGame, uint8 gamesNonStopping, uint8 injuryWeeksLeft) constant returns(uint256 encoded)
func (_Leagues *LeaguesCaller) EncodePlayerSkills(opts *bind.CallOpts, skills [5]uint16, monthOfBirth *big.Int, playerId *big.Int, birthTraits [4]uint8, alignedLastHalf bool, redCardLastGame bool, gamesNonStopping uint8, injuryWeeksLeft uint8) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "encodePlayerSkills", skills, monthOfBirth, playerId, birthTraits, alignedLastHalf, redCardLastGame, gamesNonStopping, injuryWeeksLeft)
	return *ret0, err
}

// EncodePlayerSkills is a free data retrieval call binding the contract method 0x9c53e3fd.
//
// Solidity: function encodePlayerSkills(uint16[5] skills, uint256 monthOfBirth, uint256 playerId, uint8[4] birthTraits, bool alignedLastHalf, bool redCardLastGame, uint8 gamesNonStopping, uint8 injuryWeeksLeft) constant returns(uint256 encoded)
func (_Leagues *LeaguesSession) EncodePlayerSkills(skills [5]uint16, monthOfBirth *big.Int, playerId *big.Int, birthTraits [4]uint8, alignedLastHalf bool, redCardLastGame bool, gamesNonStopping uint8, injuryWeeksLeft uint8) (*big.Int, error) {
	return _Leagues.Contract.EncodePlayerSkills(&_Leagues.CallOpts, skills, monthOfBirth, playerId, birthTraits, alignedLastHalf, redCardLastGame, gamesNonStopping, injuryWeeksLeft)
}

// EncodePlayerSkills is a free data retrieval call binding the contract method 0x9c53e3fd.
//
// Solidity: function encodePlayerSkills(uint16[5] skills, uint256 monthOfBirth, uint256 playerId, uint8[4] birthTraits, bool alignedLastHalf, bool redCardLastGame, uint8 gamesNonStopping, uint8 injuryWeeksLeft) constant returns(uint256 encoded)
func (_Leagues *LeaguesCallerSession) EncodePlayerSkills(skills [5]uint16, monthOfBirth *big.Int, playerId *big.Int, birthTraits [4]uint8, alignedLastHalf bool, redCardLastGame bool, gamesNonStopping uint8, injuryWeeksLeft uint8) (*big.Int, error) {
	return _Leagues.Contract.EncodePlayerSkills(&_Leagues.CallOpts, skills, monthOfBirth, playerId, birthTraits, alignedLastHalf, redCardLastGame, gamesNonStopping, injuryWeeksLeft)
}

// EncodePlayerState is a free data retrieval call binding the contract method 0x9f27112a.
//
// Solidity: function encodePlayerState(uint256 playerId, uint256 currentTeamId, uint8 currentShirtNum, uint256 prevPlayerTeamId, uint256 lastSaleBlock) constant returns(uint256)
func (_Leagues *LeaguesCaller) EncodePlayerState(opts *bind.CallOpts, playerId *big.Int, currentTeamId *big.Int, currentShirtNum uint8, prevPlayerTeamId *big.Int, lastSaleBlock *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "encodePlayerState", playerId, currentTeamId, currentShirtNum, prevPlayerTeamId, lastSaleBlock)
	return *ret0, err
}

// EncodePlayerState is a free data retrieval call binding the contract method 0x9f27112a.
//
// Solidity: function encodePlayerState(uint256 playerId, uint256 currentTeamId, uint8 currentShirtNum, uint256 prevPlayerTeamId, uint256 lastSaleBlock) constant returns(uint256)
func (_Leagues *LeaguesSession) EncodePlayerState(playerId *big.Int, currentTeamId *big.Int, currentShirtNum uint8, prevPlayerTeamId *big.Int, lastSaleBlock *big.Int) (*big.Int, error) {
	return _Leagues.Contract.EncodePlayerState(&_Leagues.CallOpts, playerId, currentTeamId, currentShirtNum, prevPlayerTeamId, lastSaleBlock)
}

// EncodePlayerState is a free data retrieval call binding the contract method 0x9f27112a.
//
// Solidity: function encodePlayerState(uint256 playerId, uint256 currentTeamId, uint8 currentShirtNum, uint256 prevPlayerTeamId, uint256 lastSaleBlock) constant returns(uint256)
func (_Leagues *LeaguesCallerSession) EncodePlayerState(playerId *big.Int, currentTeamId *big.Int, currentShirtNum uint8, prevPlayerTeamId *big.Int, lastSaleBlock *big.Int) (*big.Int, error) {
	return _Leagues.Contract.EncodePlayerState(&_Leagues.CallOpts, playerId, currentTeamId, currentShirtNum, prevPlayerTeamId, lastSaleBlock)
}

// EncodeTZCountryAndVal is a free data retrieval call binding the contract method 0x20748ae8.
//
// Solidity: function encodeTZCountryAndVal(uint8 timeZone, uint256 countryIdxInTZ, uint256 val) constant returns(uint256)
func (_Leagues *LeaguesCaller) EncodeTZCountryAndVal(opts *bind.CallOpts, timeZone uint8, countryIdxInTZ *big.Int, val *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "encodeTZCountryAndVal", timeZone, countryIdxInTZ, val)
	return *ret0, err
}

// EncodeTZCountryAndVal is a free data retrieval call binding the contract method 0x20748ae8.
//
// Solidity: function encodeTZCountryAndVal(uint8 timeZone, uint256 countryIdxInTZ, uint256 val) constant returns(uint256)
func (_Leagues *LeaguesSession) EncodeTZCountryAndVal(timeZone uint8, countryIdxInTZ *big.Int, val *big.Int) (*big.Int, error) {
	return _Leagues.Contract.EncodeTZCountryAndVal(&_Leagues.CallOpts, timeZone, countryIdxInTZ, val)
}

// EncodeTZCountryAndVal is a free data retrieval call binding the contract method 0x20748ae8.
//
// Solidity: function encodeTZCountryAndVal(uint8 timeZone, uint256 countryIdxInTZ, uint256 val) constant returns(uint256)
func (_Leagues *LeaguesCallerSession) EncodeTZCountryAndVal(timeZone uint8, countryIdxInTZ *big.Int, val *big.Int) (*big.Int, error) {
	return _Leagues.Contract.EncodeTZCountryAndVal(&_Leagues.CallOpts, timeZone, countryIdxInTZ, val)
}

// EncodeTactics is a free data retrieval call binding the contract method 0xe9e71652.
//
// Solidity: function encodeTactics(uint8[11] lineup, bool[10] extraAttack, uint8 tacticsId) constant returns(uint256)
func (_Leagues *LeaguesCaller) EncodeTactics(opts *bind.CallOpts, lineup [11]uint8, extraAttack [10]bool, tacticsId uint8) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "encodeTactics", lineup, extraAttack, tacticsId)
	return *ret0, err
}

// EncodeTactics is a free data retrieval call binding the contract method 0xe9e71652.
//
// Solidity: function encodeTactics(uint8[11] lineup, bool[10] extraAttack, uint8 tacticsId) constant returns(uint256)
func (_Leagues *LeaguesSession) EncodeTactics(lineup [11]uint8, extraAttack [10]bool, tacticsId uint8) (*big.Int, error) {
	return _Leagues.Contract.EncodeTactics(&_Leagues.CallOpts, lineup, extraAttack, tacticsId)
}

// EncodeTactics is a free data retrieval call binding the contract method 0xe9e71652.
//
// Solidity: function encodeTactics(uint8[11] lineup, bool[10] extraAttack, uint8 tacticsId) constant returns(uint256)
func (_Leagues *LeaguesCallerSession) EncodeTactics(lineup [11]uint8, extraAttack [10]bool, tacticsId uint8) (*big.Int, error) {
	return _Leagues.Contract.EncodeTactics(&_Leagues.CallOpts, lineup, extraAttack, tacticsId)
}

// GameDeployMonth is a free data retrieval call binding the contract method 0x85982431.
//
// Solidity: function gameDeployMonth() constant returns(uint256)
func (_Leagues *LeaguesCaller) GameDeployMonth(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "gameDeployMonth")
	return *ret0, err
}

// GameDeployMonth is a free data retrieval call binding the contract method 0x85982431.
//
// Solidity: function gameDeployMonth() constant returns(uint256)
func (_Leagues *LeaguesSession) GameDeployMonth() (*big.Int, error) {
	return _Leagues.Contract.GameDeployMonth(&_Leagues.CallOpts)
}

// GameDeployMonth is a free data retrieval call binding the contract method 0x85982431.
//
// Solidity: function gameDeployMonth() constant returns(uint256)
func (_Leagues *LeaguesCallerSession) GameDeployMonth() (*big.Int, error) {
	return _Leagues.Contract.GameDeployMonth(&_Leagues.CallOpts)
}

// GetAggressiveness is a free data retrieval call binding the contract method 0x1fc7768f.
//
// Solidity: function getAggressiveness(uint256 encodedSkills) constant returns(uint256)
func (_Leagues *LeaguesCaller) GetAggressiveness(opts *bind.CallOpts, encodedSkills *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "getAggressiveness", encodedSkills)
	return *ret0, err
}

// GetAggressiveness is a free data retrieval call binding the contract method 0x1fc7768f.
//
// Solidity: function getAggressiveness(uint256 encodedSkills) constant returns(uint256)
func (_Leagues *LeaguesSession) GetAggressiveness(encodedSkills *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetAggressiveness(&_Leagues.CallOpts, encodedSkills)
}

// GetAggressiveness is a free data retrieval call binding the contract method 0x1fc7768f.
//
// Solidity: function getAggressiveness(uint256 encodedSkills) constant returns(uint256)
func (_Leagues *LeaguesCallerSession) GetAggressiveness(encodedSkills *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetAggressiveness(&_Leagues.CallOpts, encodedSkills)
}

// GetAlignedLastHalf is a free data retrieval call binding the contract method 0x673fe242.
//
// Solidity: function getAlignedLastHalf(uint256 encodedSkills) constant returns(bool)
func (_Leagues *LeaguesCaller) GetAlignedLastHalf(opts *bind.CallOpts, encodedSkills *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "getAlignedLastHalf", encodedSkills)
	return *ret0, err
}

// GetAlignedLastHalf is a free data retrieval call binding the contract method 0x673fe242.
//
// Solidity: function getAlignedLastHalf(uint256 encodedSkills) constant returns(bool)
func (_Leagues *LeaguesSession) GetAlignedLastHalf(encodedSkills *big.Int) (bool, error) {
	return _Leagues.Contract.GetAlignedLastHalf(&_Leagues.CallOpts, encodedSkills)
}

// GetAlignedLastHalf is a free data retrieval call binding the contract method 0x673fe242.
//
// Solidity: function getAlignedLastHalf(uint256 encodedSkills) constant returns(bool)
func (_Leagues *LeaguesCallerSession) GetAlignedLastHalf(encodedSkills *big.Int) (bool, error) {
	return _Leagues.Contract.GetAlignedLastHalf(&_Leagues.CallOpts, encodedSkills)
}

// GetCurrentShirtNum is a free data retrieval call binding the contract method 0xeb78b7b7.
//
// Solidity: function getCurrentShirtNum(uint256 playerState) constant returns(uint256)
func (_Leagues *LeaguesCaller) GetCurrentShirtNum(opts *bind.CallOpts, playerState *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "getCurrentShirtNum", playerState)
	return *ret0, err
}

// GetCurrentShirtNum is a free data retrieval call binding the contract method 0xeb78b7b7.
//
// Solidity: function getCurrentShirtNum(uint256 playerState) constant returns(uint256)
func (_Leagues *LeaguesSession) GetCurrentShirtNum(playerState *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetCurrentShirtNum(&_Leagues.CallOpts, playerState)
}

// GetCurrentShirtNum is a free data retrieval call binding the contract method 0xeb78b7b7.
//
// Solidity: function getCurrentShirtNum(uint256 playerState) constant returns(uint256)
func (_Leagues *LeaguesCallerSession) GetCurrentShirtNum(playerState *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetCurrentShirtNum(&_Leagues.CallOpts, playerState)
}

// GetCurrentTeamId is a free data retrieval call binding the contract method 0xcd2105e8.
//
// Solidity: function getCurrentTeamId(uint256 playerState) constant returns(uint256)
func (_Leagues *LeaguesCaller) GetCurrentTeamId(opts *bind.CallOpts, playerState *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "getCurrentTeamId", playerState)
	return *ret0, err
}

// GetCurrentTeamId is a free data retrieval call binding the contract method 0xcd2105e8.
//
// Solidity: function getCurrentTeamId(uint256 playerState) constant returns(uint256)
func (_Leagues *LeaguesSession) GetCurrentTeamId(playerState *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetCurrentTeamId(&_Leagues.CallOpts, playerState)
}

// GetCurrentTeamId is a free data retrieval call binding the contract method 0xcd2105e8.
//
// Solidity: function getCurrentTeamId(uint256 playerState) constant returns(uint256)
func (_Leagues *LeaguesCallerSession) GetCurrentTeamId(playerState *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetCurrentTeamId(&_Leagues.CallOpts, playerState)
}

// GetCurrentTeamIdFromPlayerId is a free data retrieval call binding the contract method 0x38c96b5c.
//
// Solidity: function getCurrentTeamIdFromPlayerId(uint256 playerId) constant returns(uint256)
func (_Leagues *LeaguesCaller) GetCurrentTeamIdFromPlayerId(opts *bind.CallOpts, playerId *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "getCurrentTeamIdFromPlayerId", playerId)
	return *ret0, err
}

// GetCurrentTeamIdFromPlayerId is a free data retrieval call binding the contract method 0x38c96b5c.
//
// Solidity: function getCurrentTeamIdFromPlayerId(uint256 playerId) constant returns(uint256)
func (_Leagues *LeaguesSession) GetCurrentTeamIdFromPlayerId(playerId *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetCurrentTeamIdFromPlayerId(&_Leagues.CallOpts, playerId)
}

// GetCurrentTeamIdFromPlayerId is a free data retrieval call binding the contract method 0x38c96b5c.
//
// Solidity: function getCurrentTeamIdFromPlayerId(uint256 playerId) constant returns(uint256)
func (_Leagues *LeaguesCallerSession) GetCurrentTeamIdFromPlayerId(playerId *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetCurrentTeamIdFromPlayerId(&_Leagues.CallOpts, playerId)
}

// GetDefaultPlayerIdForTeamInCountry is a free data retrieval call binding the contract method 0x228408b0.
//
// Solidity: function getDefaultPlayerIdForTeamInCountry(uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry, uint8 shirtNum) constant returns(uint256)
func (_Leagues *LeaguesCaller) GetDefaultPlayerIdForTeamInCountry(opts *bind.CallOpts, timeZone uint8, countryIdxInTZ *big.Int, teamIdxInCountry *big.Int, shirtNum uint8) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "getDefaultPlayerIdForTeamInCountry", timeZone, countryIdxInTZ, teamIdxInCountry, shirtNum)
	return *ret0, err
}

// GetDefaultPlayerIdForTeamInCountry is a free data retrieval call binding the contract method 0x228408b0.
//
// Solidity: function getDefaultPlayerIdForTeamInCountry(uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry, uint8 shirtNum) constant returns(uint256)
func (_Leagues *LeaguesSession) GetDefaultPlayerIdForTeamInCountry(timeZone uint8, countryIdxInTZ *big.Int, teamIdxInCountry *big.Int, shirtNum uint8) (*big.Int, error) {
	return _Leagues.Contract.GetDefaultPlayerIdForTeamInCountry(&_Leagues.CallOpts, timeZone, countryIdxInTZ, teamIdxInCountry, shirtNum)
}

// GetDefaultPlayerIdForTeamInCountry is a free data retrieval call binding the contract method 0x228408b0.
//
// Solidity: function getDefaultPlayerIdForTeamInCountry(uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry, uint8 shirtNum) constant returns(uint256)
func (_Leagues *LeaguesCallerSession) GetDefaultPlayerIdForTeamInCountry(timeZone uint8, countryIdxInTZ *big.Int, teamIdxInCountry *big.Int, shirtNum uint8) (*big.Int, error) {
	return _Leagues.Contract.GetDefaultPlayerIdForTeamInCountry(&_Leagues.CallOpts, timeZone, countryIdxInTZ, teamIdxInCountry, shirtNum)
}

// GetDefence is a free data retrieval call binding the contract method 0x51585b49.
//
// Solidity: function getDefence(uint256 encodedSkills) constant returns(uint256)
func (_Leagues *LeaguesCaller) GetDefence(opts *bind.CallOpts, encodedSkills *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "getDefence", encodedSkills)
	return *ret0, err
}

// GetDefence is a free data retrieval call binding the contract method 0x51585b49.
//
// Solidity: function getDefence(uint256 encodedSkills) constant returns(uint256)
func (_Leagues *LeaguesSession) GetDefence(encodedSkills *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetDefence(&_Leagues.CallOpts, encodedSkills)
}

// GetDefence is a free data retrieval call binding the contract method 0x51585b49.
//
// Solidity: function getDefence(uint256 encodedSkills) constant returns(uint256)
func (_Leagues *LeaguesCallerSession) GetDefence(encodedSkills *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetDefence(&_Leagues.CallOpts, encodedSkills)
}

// GetEndurance is a free data retrieval call binding the contract method 0x258e5d90.
//
// Solidity: function getEndurance(uint256 encodedSkills) constant returns(uint256)
func (_Leagues *LeaguesCaller) GetEndurance(opts *bind.CallOpts, encodedSkills *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "getEndurance", encodedSkills)
	return *ret0, err
}

// GetEndurance is a free data retrieval call binding the contract method 0x258e5d90.
//
// Solidity: function getEndurance(uint256 encodedSkills) constant returns(uint256)
func (_Leagues *LeaguesSession) GetEndurance(encodedSkills *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetEndurance(&_Leagues.CallOpts, encodedSkills)
}

// GetEndurance is a free data retrieval call binding the contract method 0x258e5d90.
//
// Solidity: function getEndurance(uint256 encodedSkills) constant returns(uint256)
func (_Leagues *LeaguesCallerSession) GetEndurance(encodedSkills *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetEndurance(&_Leagues.CallOpts, encodedSkills)
}

// GetEngineAddress is a free data retrieval call binding the contract method 0x4562a618.
//
// Solidity: function getEngineAddress() constant returns(address)
func (_Leagues *LeaguesCaller) GetEngineAddress(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "getEngineAddress")
	return *ret0, err
}

// GetEngineAddress is a free data retrieval call binding the contract method 0x4562a618.
//
// Solidity: function getEngineAddress() constant returns(address)
func (_Leagues *LeaguesSession) GetEngineAddress() (common.Address, error) {
	return _Leagues.Contract.GetEngineAddress(&_Leagues.CallOpts)
}

// GetEngineAddress is a free data retrieval call binding the contract method 0x4562a618.
//
// Solidity: function getEngineAddress() constant returns(address)
func (_Leagues *LeaguesCallerSession) GetEngineAddress() (common.Address, error) {
	return _Leagues.Contract.GetEngineAddress(&_Leagues.CallOpts)
}

// GetForwardness is a free data retrieval call binding the contract method 0xc2bc41cd.
//
// Solidity: function getForwardness(uint256 encodedSkills) constant returns(uint256)
func (_Leagues *LeaguesCaller) GetForwardness(opts *bind.CallOpts, encodedSkills *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "getForwardness", encodedSkills)
	return *ret0, err
}

// GetForwardness is a free data retrieval call binding the contract method 0xc2bc41cd.
//
// Solidity: function getForwardness(uint256 encodedSkills) constant returns(uint256)
func (_Leagues *LeaguesSession) GetForwardness(encodedSkills *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetForwardness(&_Leagues.CallOpts, encodedSkills)
}

// GetForwardness is a free data retrieval call binding the contract method 0xc2bc41cd.
//
// Solidity: function getForwardness(uint256 encodedSkills) constant returns(uint256)
func (_Leagues *LeaguesCallerSession) GetForwardness(encodedSkills *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetForwardness(&_Leagues.CallOpts, encodedSkills)
}

// GetFreeShirt is a free data retrieval call binding the contract method 0x507b1723.
//
// Solidity: function getFreeShirt(uint256 teamId) constant returns(uint8)
func (_Leagues *LeaguesCaller) GetFreeShirt(opts *bind.CallOpts, teamId *big.Int) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "getFreeShirt", teamId)
	return *ret0, err
}

// GetFreeShirt is a free data retrieval call binding the contract method 0x507b1723.
//
// Solidity: function getFreeShirt(uint256 teamId) constant returns(uint8)
func (_Leagues *LeaguesSession) GetFreeShirt(teamId *big.Int) (uint8, error) {
	return _Leagues.Contract.GetFreeShirt(&_Leagues.CallOpts, teamId)
}

// GetFreeShirt is a free data retrieval call binding the contract method 0x507b1723.
//
// Solidity: function getFreeShirt(uint256 teamId) constant returns(uint8)
func (_Leagues *LeaguesCallerSession) GetFreeShirt(teamId *big.Int) (uint8, error) {
	return _Leagues.Contract.GetFreeShirt(&_Leagues.CallOpts, teamId)
}

// GetGamesNonStopping is a free data retrieval call binding the contract method 0xe804e519.
//
// Solidity: function getGamesNonStopping(uint256 encodedSkills) constant returns(uint256)
func (_Leagues *LeaguesCaller) GetGamesNonStopping(opts *bind.CallOpts, encodedSkills *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "getGamesNonStopping", encodedSkills)
	return *ret0, err
}

// GetGamesNonStopping is a free data retrieval call binding the contract method 0xe804e519.
//
// Solidity: function getGamesNonStopping(uint256 encodedSkills) constant returns(uint256)
func (_Leagues *LeaguesSession) GetGamesNonStopping(encodedSkills *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetGamesNonStopping(&_Leagues.CallOpts, encodedSkills)
}

// GetGamesNonStopping is a free data retrieval call binding the contract method 0xe804e519.
//
// Solidity: function getGamesNonStopping(uint256 encodedSkills) constant returns(uint256)
func (_Leagues *LeaguesCallerSession) GetGamesNonStopping(encodedSkills *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetGamesNonStopping(&_Leagues.CallOpts, encodedSkills)
}

// GetInjuryWeeksLeft is a free data retrieval call binding the contract method 0x79e76597.
//
// Solidity: function getInjuryWeeksLeft(uint256 encodedSkills) constant returns(uint256)
func (_Leagues *LeaguesCaller) GetInjuryWeeksLeft(opts *bind.CallOpts, encodedSkills *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "getInjuryWeeksLeft", encodedSkills)
	return *ret0, err
}

// GetInjuryWeeksLeft is a free data retrieval call binding the contract method 0x79e76597.
//
// Solidity: function getInjuryWeeksLeft(uint256 encodedSkills) constant returns(uint256)
func (_Leagues *LeaguesSession) GetInjuryWeeksLeft(encodedSkills *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetInjuryWeeksLeft(&_Leagues.CallOpts, encodedSkills)
}

// GetInjuryWeeksLeft is a free data retrieval call binding the contract method 0x79e76597.
//
// Solidity: function getInjuryWeeksLeft(uint256 encodedSkills) constant returns(uint256)
func (_Leagues *LeaguesCallerSession) GetInjuryWeeksLeft(encodedSkills *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetInjuryWeeksLeft(&_Leagues.CallOpts, encodedSkills)
}

// GetLastActionsSubmissionTime is a free data retrieval call binding the contract method 0xfa80039b.
//
// Solidity: function getLastActionsSubmissionTime(uint8 timeZone) constant returns(uint256)
func (_Leagues *LeaguesCaller) GetLastActionsSubmissionTime(opts *bind.CallOpts, timeZone uint8) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "getLastActionsSubmissionTime", timeZone)
	return *ret0, err
}

// GetLastActionsSubmissionTime is a free data retrieval call binding the contract method 0xfa80039b.
//
// Solidity: function getLastActionsSubmissionTime(uint8 timeZone) constant returns(uint256)
func (_Leagues *LeaguesSession) GetLastActionsSubmissionTime(timeZone uint8) (*big.Int, error) {
	return _Leagues.Contract.GetLastActionsSubmissionTime(&_Leagues.CallOpts, timeZone)
}

// GetLastActionsSubmissionTime is a free data retrieval call binding the contract method 0xfa80039b.
//
// Solidity: function getLastActionsSubmissionTime(uint8 timeZone) constant returns(uint256)
func (_Leagues *LeaguesCallerSession) GetLastActionsSubmissionTime(timeZone uint8) (*big.Int, error) {
	return _Leagues.Contract.GetLastActionsSubmissionTime(&_Leagues.CallOpts, timeZone)
}

// GetLastSaleBlock is a free data retrieval call binding the contract method 0xc566b5bc.
//
// Solidity: function getLastSaleBlock(uint256 playerState) constant returns(uint256)
func (_Leagues *LeaguesCaller) GetLastSaleBlock(opts *bind.CallOpts, playerState *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "getLastSaleBlock", playerState)
	return *ret0, err
}

// GetLastSaleBlock is a free data retrieval call binding the contract method 0xc566b5bc.
//
// Solidity: function getLastSaleBlock(uint256 playerState) constant returns(uint256)
func (_Leagues *LeaguesSession) GetLastSaleBlock(playerState *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetLastSaleBlock(&_Leagues.CallOpts, playerState)
}

// GetLastSaleBlock is a free data retrieval call binding the contract method 0xc566b5bc.
//
// Solidity: function getLastSaleBlock(uint256 playerState) constant returns(uint256)
func (_Leagues *LeaguesCallerSession) GetLastSaleBlock(playerState *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetLastSaleBlock(&_Leagues.CallOpts, playerState)
}

// GetLastUpdateTime is a free data retrieval call binding the contract method 0x2d0e08fd.
//
// Solidity: function getLastUpdateTime(uint8 timeZone) constant returns(uint256)
func (_Leagues *LeaguesCaller) GetLastUpdateTime(opts *bind.CallOpts, timeZone uint8) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "getLastUpdateTime", timeZone)
	return *ret0, err
}

// GetLastUpdateTime is a free data retrieval call binding the contract method 0x2d0e08fd.
//
// Solidity: function getLastUpdateTime(uint8 timeZone) constant returns(uint256)
func (_Leagues *LeaguesSession) GetLastUpdateTime(timeZone uint8) (*big.Int, error) {
	return _Leagues.Contract.GetLastUpdateTime(&_Leagues.CallOpts, timeZone)
}

// GetLastUpdateTime is a free data retrieval call binding the contract method 0x2d0e08fd.
//
// Solidity: function getLastUpdateTime(uint8 timeZone) constant returns(uint256)
func (_Leagues *LeaguesCallerSession) GetLastUpdateTime(timeZone uint8) (*big.Int, error) {
	return _Leagues.Contract.GetLastUpdateTime(&_Leagues.CallOpts, timeZone)
}

// GetLeftishness is a free data retrieval call binding the contract method 0x3518dd1d.
//
// Solidity: function getLeftishness(uint256 encodedSkills) constant returns(uint256)
func (_Leagues *LeaguesCaller) GetLeftishness(opts *bind.CallOpts, encodedSkills *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "getLeftishness", encodedSkills)
	return *ret0, err
}

// GetLeftishness is a free data retrieval call binding the contract method 0x3518dd1d.
//
// Solidity: function getLeftishness(uint256 encodedSkills) constant returns(uint256)
func (_Leagues *LeaguesSession) GetLeftishness(encodedSkills *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetLeftishness(&_Leagues.CallOpts, encodedSkills)
}

// GetLeftishness is a free data retrieval call binding the contract method 0x3518dd1d.
//
// Solidity: function getLeftishness(uint256 encodedSkills) constant returns(uint256)
func (_Leagues *LeaguesCallerSession) GetLeftishness(encodedSkills *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetLeftishness(&_Leagues.CallOpts, encodedSkills)
}

// GetMonthOfBirth is a free data retrieval call binding the contract method 0x87f1e880.
//
// Solidity: function getMonthOfBirth(uint256 encodedSkills) constant returns(uint256)
func (_Leagues *LeaguesCaller) GetMonthOfBirth(opts *bind.CallOpts, encodedSkills *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "getMonthOfBirth", encodedSkills)
	return *ret0, err
}

// GetMonthOfBirth is a free data retrieval call binding the contract method 0x87f1e880.
//
// Solidity: function getMonthOfBirth(uint256 encodedSkills) constant returns(uint256)
func (_Leagues *LeaguesSession) GetMonthOfBirth(encodedSkills *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetMonthOfBirth(&_Leagues.CallOpts, encodedSkills)
}

// GetMonthOfBirth is a free data retrieval call binding the contract method 0x87f1e880.
//
// Solidity: function getMonthOfBirth(uint256 encodedSkills) constant returns(uint256)
func (_Leagues *LeaguesCallerSession) GetMonthOfBirth(encodedSkills *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetMonthOfBirth(&_Leagues.CallOpts, encodedSkills)
}

// GetNCountriesInTZ is a free data retrieval call binding the contract method 0xad63bcbd.
//
// Solidity: function getNCountriesInTZ(uint8 timeZone) constant returns(uint256)
func (_Leagues *LeaguesCaller) GetNCountriesInTZ(opts *bind.CallOpts, timeZone uint8) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "getNCountriesInTZ", timeZone)
	return *ret0, err
}

// GetNCountriesInTZ is a free data retrieval call binding the contract method 0xad63bcbd.
//
// Solidity: function getNCountriesInTZ(uint8 timeZone) constant returns(uint256)
func (_Leagues *LeaguesSession) GetNCountriesInTZ(timeZone uint8) (*big.Int, error) {
	return _Leagues.Contract.GetNCountriesInTZ(&_Leagues.CallOpts, timeZone)
}

// GetNCountriesInTZ is a free data retrieval call binding the contract method 0xad63bcbd.
//
// Solidity: function getNCountriesInTZ(uint8 timeZone) constant returns(uint256)
func (_Leagues *LeaguesCallerSession) GetNCountriesInTZ(timeZone uint8) (*big.Int, error) {
	return _Leagues.Contract.GetNCountriesInTZ(&_Leagues.CallOpts, timeZone)
}

// GetNDivisionsInCountry is a free data retrieval call binding the contract method 0x5adb40f5.
//
// Solidity: function getNDivisionsInCountry(uint8 timeZone, uint256 countryIdxInTZ) constant returns(uint256)
func (_Leagues *LeaguesCaller) GetNDivisionsInCountry(opts *bind.CallOpts, timeZone uint8, countryIdxInTZ *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "getNDivisionsInCountry", timeZone, countryIdxInTZ)
	return *ret0, err
}

// GetNDivisionsInCountry is a free data retrieval call binding the contract method 0x5adb40f5.
//
// Solidity: function getNDivisionsInCountry(uint8 timeZone, uint256 countryIdxInTZ) constant returns(uint256)
func (_Leagues *LeaguesSession) GetNDivisionsInCountry(timeZone uint8, countryIdxInTZ *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetNDivisionsInCountry(&_Leagues.CallOpts, timeZone, countryIdxInTZ)
}

// GetNDivisionsInCountry is a free data retrieval call binding the contract method 0x5adb40f5.
//
// Solidity: function getNDivisionsInCountry(uint8 timeZone, uint256 countryIdxInTZ) constant returns(uint256)
func (_Leagues *LeaguesCallerSession) GetNDivisionsInCountry(timeZone uint8, countryIdxInTZ *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetNDivisionsInCountry(&_Leagues.CallOpts, timeZone, countryIdxInTZ)
}

// GetNLeaguesInCountry is a free data retrieval call binding the contract method 0xf9d0723d.
//
// Solidity: function getNLeaguesInCountry(uint8 timeZone, uint256 countryIdxInTZ) constant returns(uint256)
func (_Leagues *LeaguesCaller) GetNLeaguesInCountry(opts *bind.CallOpts, timeZone uint8, countryIdxInTZ *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "getNLeaguesInCountry", timeZone, countryIdxInTZ)
	return *ret0, err
}

// GetNLeaguesInCountry is a free data retrieval call binding the contract method 0xf9d0723d.
//
// Solidity: function getNLeaguesInCountry(uint8 timeZone, uint256 countryIdxInTZ) constant returns(uint256)
func (_Leagues *LeaguesSession) GetNLeaguesInCountry(timeZone uint8, countryIdxInTZ *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetNLeaguesInCountry(&_Leagues.CallOpts, timeZone, countryIdxInTZ)
}

// GetNLeaguesInCountry is a free data retrieval call binding the contract method 0xf9d0723d.
//
// Solidity: function getNLeaguesInCountry(uint8 timeZone, uint256 countryIdxInTZ) constant returns(uint256)
func (_Leagues *LeaguesCallerSession) GetNLeaguesInCountry(timeZone uint8, countryIdxInTZ *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetNLeaguesInCountry(&_Leagues.CallOpts, timeZone, countryIdxInTZ)
}

// GetNTeamsInCountry is a free data retrieval call binding the contract method 0xc04f6d53.
//
// Solidity: function getNTeamsInCountry(uint8 timeZone, uint256 countryIdxInTZ) constant returns(uint256)
func (_Leagues *LeaguesCaller) GetNTeamsInCountry(opts *bind.CallOpts, timeZone uint8, countryIdxInTZ *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "getNTeamsInCountry", timeZone, countryIdxInTZ)
	return *ret0, err
}

// GetNTeamsInCountry is a free data retrieval call binding the contract method 0xc04f6d53.
//
// Solidity: function getNTeamsInCountry(uint8 timeZone, uint256 countryIdxInTZ) constant returns(uint256)
func (_Leagues *LeaguesSession) GetNTeamsInCountry(timeZone uint8, countryIdxInTZ *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetNTeamsInCountry(&_Leagues.CallOpts, timeZone, countryIdxInTZ)
}

// GetNTeamsInCountry is a free data retrieval call binding the contract method 0xc04f6d53.
//
// Solidity: function getNTeamsInCountry(uint8 timeZone, uint256 countryIdxInTZ) constant returns(uint256)
func (_Leagues *LeaguesCallerSession) GetNTeamsInCountry(timeZone uint8, countryIdxInTZ *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetNTeamsInCountry(&_Leagues.CallOpts, timeZone, countryIdxInTZ)
}

// GetOwnerPlayer is a free data retrieval call binding the contract method 0x8f9da214.
//
// Solidity: function getOwnerPlayer(uint256 playerId) constant returns(address)
func (_Leagues *LeaguesCaller) GetOwnerPlayer(opts *bind.CallOpts, playerId *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "getOwnerPlayer", playerId)
	return *ret0, err
}

// GetOwnerPlayer is a free data retrieval call binding the contract method 0x8f9da214.
//
// Solidity: function getOwnerPlayer(uint256 playerId) constant returns(address)
func (_Leagues *LeaguesSession) GetOwnerPlayer(playerId *big.Int) (common.Address, error) {
	return _Leagues.Contract.GetOwnerPlayer(&_Leagues.CallOpts, playerId)
}

// GetOwnerPlayer is a free data retrieval call binding the contract method 0x8f9da214.
//
// Solidity: function getOwnerPlayer(uint256 playerId) constant returns(address)
func (_Leagues *LeaguesCallerSession) GetOwnerPlayer(playerId *big.Int) (common.Address, error) {
	return _Leagues.Contract.GetOwnerPlayer(&_Leagues.CallOpts, playerId)
}

// GetOwnerTeam is a free data retrieval call binding the contract method 0x492afc69.
//
// Solidity: function getOwnerTeam(uint256 teamId) constant returns(address)
func (_Leagues *LeaguesCaller) GetOwnerTeam(opts *bind.CallOpts, teamId *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "getOwnerTeam", teamId)
	return *ret0, err
}

// GetOwnerTeam is a free data retrieval call binding the contract method 0x492afc69.
//
// Solidity: function getOwnerTeam(uint256 teamId) constant returns(address)
func (_Leagues *LeaguesSession) GetOwnerTeam(teamId *big.Int) (common.Address, error) {
	return _Leagues.Contract.GetOwnerTeam(&_Leagues.CallOpts, teamId)
}

// GetOwnerTeam is a free data retrieval call binding the contract method 0x492afc69.
//
// Solidity: function getOwnerTeam(uint256 teamId) constant returns(address)
func (_Leagues *LeaguesCallerSession) GetOwnerTeam(teamId *big.Int) (common.Address, error) {
	return _Leagues.Contract.GetOwnerTeam(&_Leagues.CallOpts, teamId)
}

// GetOwnerTeamInCountry is a free data retrieval call binding the contract method 0x595ef25b.
//
// Solidity: function getOwnerTeamInCountry(uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) constant returns(address)
func (_Leagues *LeaguesCaller) GetOwnerTeamInCountry(opts *bind.CallOpts, timeZone uint8, countryIdxInTZ *big.Int, teamIdxInCountry *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "getOwnerTeamInCountry", timeZone, countryIdxInTZ, teamIdxInCountry)
	return *ret0, err
}

// GetOwnerTeamInCountry is a free data retrieval call binding the contract method 0x595ef25b.
//
// Solidity: function getOwnerTeamInCountry(uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) constant returns(address)
func (_Leagues *LeaguesSession) GetOwnerTeamInCountry(timeZone uint8, countryIdxInTZ *big.Int, teamIdxInCountry *big.Int) (common.Address, error) {
	return _Leagues.Contract.GetOwnerTeamInCountry(&_Leagues.CallOpts, timeZone, countryIdxInTZ, teamIdxInCountry)
}

// GetOwnerTeamInCountry is a free data retrieval call binding the contract method 0x595ef25b.
//
// Solidity: function getOwnerTeamInCountry(uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) constant returns(address)
func (_Leagues *LeaguesCallerSession) GetOwnerTeamInCountry(timeZone uint8, countryIdxInTZ *big.Int, teamIdxInCountry *big.Int) (common.Address, error) {
	return _Leagues.Contract.GetOwnerTeamInCountry(&_Leagues.CallOpts, timeZone, countryIdxInTZ, teamIdxInCountry)
}

// GetPass is a free data retrieval call binding the contract method 0x55a6f86f.
//
// Solidity: function getPass(uint256 encodedSkills) constant returns(uint256)
func (_Leagues *LeaguesCaller) GetPass(opts *bind.CallOpts, encodedSkills *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "getPass", encodedSkills)
	return *ret0, err
}

// GetPass is a free data retrieval call binding the contract method 0x55a6f86f.
//
// Solidity: function getPass(uint256 encodedSkills) constant returns(uint256)
func (_Leagues *LeaguesSession) GetPass(encodedSkills *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetPass(&_Leagues.CallOpts, encodedSkills)
}

// GetPass is a free data retrieval call binding the contract method 0x55a6f86f.
//
// Solidity: function getPass(uint256 encodedSkills) constant returns(uint256)
func (_Leagues *LeaguesCallerSession) GetPass(encodedSkills *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetPass(&_Leagues.CallOpts, encodedSkills)
}

// GetPlayerAgeInMonths is a free data retrieval call binding the contract method 0x1ffeb349.
//
// Solidity: function getPlayerAgeInMonths(uint256 playerId) constant returns(uint256)
func (_Leagues *LeaguesCaller) GetPlayerAgeInMonths(opts *bind.CallOpts, playerId *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "getPlayerAgeInMonths", playerId)
	return *ret0, err
}

// GetPlayerAgeInMonths is a free data retrieval call binding the contract method 0x1ffeb349.
//
// Solidity: function getPlayerAgeInMonths(uint256 playerId) constant returns(uint256)
func (_Leagues *LeaguesSession) GetPlayerAgeInMonths(playerId *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetPlayerAgeInMonths(&_Leagues.CallOpts, playerId)
}

// GetPlayerAgeInMonths is a free data retrieval call binding the contract method 0x1ffeb349.
//
// Solidity: function getPlayerAgeInMonths(uint256 playerId) constant returns(uint256)
func (_Leagues *LeaguesCallerSession) GetPlayerAgeInMonths(playerId *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetPlayerAgeInMonths(&_Leagues.CallOpts, playerId)
}

// GetPlayerIdFromSkills is a free data retrieval call binding the contract method 0x6f6c2ae0.
//
// Solidity: function getPlayerIdFromSkills(uint256 encodedSkills) constant returns(uint256)
func (_Leagues *LeaguesCaller) GetPlayerIdFromSkills(opts *bind.CallOpts, encodedSkills *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "getPlayerIdFromSkills", encodedSkills)
	return *ret0, err
}

// GetPlayerIdFromSkills is a free data retrieval call binding the contract method 0x6f6c2ae0.
//
// Solidity: function getPlayerIdFromSkills(uint256 encodedSkills) constant returns(uint256)
func (_Leagues *LeaguesSession) GetPlayerIdFromSkills(encodedSkills *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetPlayerIdFromSkills(&_Leagues.CallOpts, encodedSkills)
}

// GetPlayerIdFromSkills is a free data retrieval call binding the contract method 0x6f6c2ae0.
//
// Solidity: function getPlayerIdFromSkills(uint256 encodedSkills) constant returns(uint256)
func (_Leagues *LeaguesCallerSession) GetPlayerIdFromSkills(encodedSkills *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetPlayerIdFromSkills(&_Leagues.CallOpts, encodedSkills)
}

// GetPlayerIdFromState is a free data retrieval call binding the contract method 0x78f4c718.
//
// Solidity: function getPlayerIdFromState(uint256 playerState) constant returns(uint256)
func (_Leagues *LeaguesCaller) GetPlayerIdFromState(opts *bind.CallOpts, playerState *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "getPlayerIdFromState", playerState)
	return *ret0, err
}

// GetPlayerIdFromState is a free data retrieval call binding the contract method 0x78f4c718.
//
// Solidity: function getPlayerIdFromState(uint256 playerState) constant returns(uint256)
func (_Leagues *LeaguesSession) GetPlayerIdFromState(playerState *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetPlayerIdFromState(&_Leagues.CallOpts, playerState)
}

// GetPlayerIdFromState is a free data retrieval call binding the contract method 0x78f4c718.
//
// Solidity: function getPlayerIdFromState(uint256 playerState) constant returns(uint256)
func (_Leagues *LeaguesCallerSession) GetPlayerIdFromState(playerState *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetPlayerIdFromState(&_Leagues.CallOpts, playerState)
}

// GetPlayerIdsInTeam is a free data retrieval call binding the contract method 0xeabf6a4b.
//
// Solidity: function getPlayerIdsInTeam(uint256 teamId) constant returns(uint256[25] playerIds)
func (_Leagues *LeaguesCaller) GetPlayerIdsInTeam(opts *bind.CallOpts, teamId *big.Int) ([25]*big.Int, error) {
	var (
		ret0 = new([25]*big.Int)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "getPlayerIdsInTeam", teamId)
	return *ret0, err
}

// GetPlayerIdsInTeam is a free data retrieval call binding the contract method 0xeabf6a4b.
//
// Solidity: function getPlayerIdsInTeam(uint256 teamId) constant returns(uint256[25] playerIds)
func (_Leagues *LeaguesSession) GetPlayerIdsInTeam(teamId *big.Int) ([25]*big.Int, error) {
	return _Leagues.Contract.GetPlayerIdsInTeam(&_Leagues.CallOpts, teamId)
}

// GetPlayerIdsInTeam is a free data retrieval call binding the contract method 0xeabf6a4b.
//
// Solidity: function getPlayerIdsInTeam(uint256 teamId) constant returns(uint256[25] playerIds)
func (_Leagues *LeaguesCallerSession) GetPlayerIdsInTeam(teamId *big.Int) ([25]*big.Int, error) {
	return _Leagues.Contract.GetPlayerIdsInTeam(&_Leagues.CallOpts, teamId)
}

// GetPlayerSkillsAtBirth is a free data retrieval call binding the contract method 0xc73f808d.
//
// Solidity: function getPlayerSkillsAtBirth(uint256 playerId) constant returns(uint256)
func (_Leagues *LeaguesCaller) GetPlayerSkillsAtBirth(opts *bind.CallOpts, playerId *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "getPlayerSkillsAtBirth", playerId)
	return *ret0, err
}

// GetPlayerSkillsAtBirth is a free data retrieval call binding the contract method 0xc73f808d.
//
// Solidity: function getPlayerSkillsAtBirth(uint256 playerId) constant returns(uint256)
func (_Leagues *LeaguesSession) GetPlayerSkillsAtBirth(playerId *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetPlayerSkillsAtBirth(&_Leagues.CallOpts, playerId)
}

// GetPlayerSkillsAtBirth is a free data retrieval call binding the contract method 0xc73f808d.
//
// Solidity: function getPlayerSkillsAtBirth(uint256 playerId) constant returns(uint256)
func (_Leagues *LeaguesCallerSession) GetPlayerSkillsAtBirth(playerId *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetPlayerSkillsAtBirth(&_Leagues.CallOpts, playerId)
}

// GetPlayerState is a free data retrieval call binding the contract method 0xec7ecec5.
//
// Solidity: function getPlayerState(uint256 playerId) constant returns(uint256)
func (_Leagues *LeaguesCaller) GetPlayerState(opts *bind.CallOpts, playerId *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "getPlayerState", playerId)
	return *ret0, err
}

// GetPlayerState is a free data retrieval call binding the contract method 0xec7ecec5.
//
// Solidity: function getPlayerState(uint256 playerId) constant returns(uint256)
func (_Leagues *LeaguesSession) GetPlayerState(playerId *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetPlayerState(&_Leagues.CallOpts, playerId)
}

// GetPlayerState is a free data retrieval call binding the contract method 0xec7ecec5.
//
// Solidity: function getPlayerState(uint256 playerId) constant returns(uint256)
func (_Leagues *LeaguesCallerSession) GetPlayerState(playerId *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetPlayerState(&_Leagues.CallOpts, playerId)
}

// GetPlayerStateAtBirth is a free data retrieval call binding the contract method 0x26657608.
//
// Solidity: function getPlayerStateAtBirth(uint256 playerId) constant returns(uint256)
func (_Leagues *LeaguesCaller) GetPlayerStateAtBirth(opts *bind.CallOpts, playerId *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "getPlayerStateAtBirth", playerId)
	return *ret0, err
}

// GetPlayerStateAtBirth is a free data retrieval call binding the contract method 0x26657608.
//
// Solidity: function getPlayerStateAtBirth(uint256 playerId) constant returns(uint256)
func (_Leagues *LeaguesSession) GetPlayerStateAtBirth(playerId *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetPlayerStateAtBirth(&_Leagues.CallOpts, playerId)
}

// GetPlayerStateAtBirth is a free data retrieval call binding the contract method 0x26657608.
//
// Solidity: function getPlayerStateAtBirth(uint256 playerId) constant returns(uint256)
func (_Leagues *LeaguesCallerSession) GetPlayerStateAtBirth(playerId *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetPlayerStateAtBirth(&_Leagues.CallOpts, playerId)
}

// GetPotential is a free data retrieval call binding the contract method 0x434807ad.
//
// Solidity: function getPotential(uint256 encodedSkills) constant returns(uint256)
func (_Leagues *LeaguesCaller) GetPotential(opts *bind.CallOpts, encodedSkills *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "getPotential", encodedSkills)
	return *ret0, err
}

// GetPotential is a free data retrieval call binding the contract method 0x434807ad.
//
// Solidity: function getPotential(uint256 encodedSkills) constant returns(uint256)
func (_Leagues *LeaguesSession) GetPotential(encodedSkills *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetPotential(&_Leagues.CallOpts, encodedSkills)
}

// GetPotential is a free data retrieval call binding the contract method 0x434807ad.
//
// Solidity: function getPotential(uint256 encodedSkills) constant returns(uint256)
func (_Leagues *LeaguesCallerSession) GetPotential(encodedSkills *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetPotential(&_Leagues.CallOpts, encodedSkills)
}

// GetPrevPlayerTeamId is a free data retrieval call binding the contract method 0x4bea2a69.
//
// Solidity: function getPrevPlayerTeamId(uint256 playerState) constant returns(uint256)
func (_Leagues *LeaguesCaller) GetPrevPlayerTeamId(opts *bind.CallOpts, playerState *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "getPrevPlayerTeamId", playerState)
	return *ret0, err
}

// GetPrevPlayerTeamId is a free data retrieval call binding the contract method 0x4bea2a69.
//
// Solidity: function getPrevPlayerTeamId(uint256 playerState) constant returns(uint256)
func (_Leagues *LeaguesSession) GetPrevPlayerTeamId(playerState *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetPrevPlayerTeamId(&_Leagues.CallOpts, playerState)
}

// GetPrevPlayerTeamId is a free data retrieval call binding the contract method 0x4bea2a69.
//
// Solidity: function getPrevPlayerTeamId(uint256 playerState) constant returns(uint256)
func (_Leagues *LeaguesCallerSession) GetPrevPlayerTeamId(playerState *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetPrevPlayerTeamId(&_Leagues.CallOpts, playerState)
}

// GetRedCardLastGame is a free data retrieval call binding the contract method 0xcc7d473b.
//
// Solidity: function getRedCardLastGame(uint256 encodedSkills) constant returns(bool)
func (_Leagues *LeaguesCaller) GetRedCardLastGame(opts *bind.CallOpts, encodedSkills *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "getRedCardLastGame", encodedSkills)
	return *ret0, err
}

// GetRedCardLastGame is a free data retrieval call binding the contract method 0xcc7d473b.
//
// Solidity: function getRedCardLastGame(uint256 encodedSkills) constant returns(bool)
func (_Leagues *LeaguesSession) GetRedCardLastGame(encodedSkills *big.Int) (bool, error) {
	return _Leagues.Contract.GetRedCardLastGame(&_Leagues.CallOpts, encodedSkills)
}

// GetRedCardLastGame is a free data retrieval call binding the contract method 0xcc7d473b.
//
// Solidity: function getRedCardLastGame(uint256 encodedSkills) constant returns(bool)
func (_Leagues *LeaguesCallerSession) GetRedCardLastGame(encodedSkills *big.Int) (bool, error) {
	return _Leagues.Contract.GetRedCardLastGame(&_Leagues.CallOpts, encodedSkills)
}

// GetShoot is a free data retrieval call binding the contract method 0x65b4b476.
//
// Solidity: function getShoot(uint256 encodedSkills) constant returns(uint256)
func (_Leagues *LeaguesCaller) GetShoot(opts *bind.CallOpts, encodedSkills *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "getShoot", encodedSkills)
	return *ret0, err
}

// GetShoot is a free data retrieval call binding the contract method 0x65b4b476.
//
// Solidity: function getShoot(uint256 encodedSkills) constant returns(uint256)
func (_Leagues *LeaguesSession) GetShoot(encodedSkills *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetShoot(&_Leagues.CallOpts, encodedSkills)
}

// GetShoot is a free data retrieval call binding the contract method 0x65b4b476.
//
// Solidity: function getShoot(uint256 encodedSkills) constant returns(uint256)
func (_Leagues *LeaguesCallerSession) GetShoot(encodedSkills *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetShoot(&_Leagues.CallOpts, encodedSkills)
}

// GetSkills is a free data retrieval call binding the contract method 0x0092bf78.
//
// Solidity: function getSkills(uint256 encodedSkills) constant returns(uint256)
func (_Leagues *LeaguesCaller) GetSkills(opts *bind.CallOpts, encodedSkills *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "getSkills", encodedSkills)
	return *ret0, err
}

// GetSkills is a free data retrieval call binding the contract method 0x0092bf78.
//
// Solidity: function getSkills(uint256 encodedSkills) constant returns(uint256)
func (_Leagues *LeaguesSession) GetSkills(encodedSkills *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetSkills(&_Leagues.CallOpts, encodedSkills)
}

// GetSkills is a free data retrieval call binding the contract method 0x0092bf78.
//
// Solidity: function getSkills(uint256 encodedSkills) constant returns(uint256)
func (_Leagues *LeaguesCallerSession) GetSkills(encodedSkills *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetSkills(&_Leagues.CallOpts, encodedSkills)
}

// GetSkillsVec is a free data retrieval call binding the contract method 0xcc1cc3d7.
//
// Solidity: function getSkillsVec(uint256 encodedSkills) constant returns(uint16[5] skills)
func (_Leagues *LeaguesCaller) GetSkillsVec(opts *bind.CallOpts, encodedSkills *big.Int) ([5]uint16, error) {
	var (
		ret0 = new([5]uint16)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "getSkillsVec", encodedSkills)
	return *ret0, err
}

// GetSkillsVec is a free data retrieval call binding the contract method 0xcc1cc3d7.
//
// Solidity: function getSkillsVec(uint256 encodedSkills) constant returns(uint16[5] skills)
func (_Leagues *LeaguesSession) GetSkillsVec(encodedSkills *big.Int) ([5]uint16, error) {
	return _Leagues.Contract.GetSkillsVec(&_Leagues.CallOpts, encodedSkills)
}

// GetSkillsVec is a free data retrieval call binding the contract method 0xcc1cc3d7.
//
// Solidity: function getSkillsVec(uint256 encodedSkills) constant returns(uint16[5] skills)
func (_Leagues *LeaguesCallerSession) GetSkillsVec(encodedSkills *big.Int) ([5]uint16, error) {
	return _Leagues.Contract.GetSkillsVec(&_Leagues.CallOpts, encodedSkills)
}

// GetSpeed is a free data retrieval call binding the contract method 0x4b93f753.
//
// Solidity: function getSpeed(uint256 encodedSkills) constant returns(uint256)
func (_Leagues *LeaguesCaller) GetSpeed(opts *bind.CallOpts, encodedSkills *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "getSpeed", encodedSkills)
	return *ret0, err
}

// GetSpeed is a free data retrieval call binding the contract method 0x4b93f753.
//
// Solidity: function getSpeed(uint256 encodedSkills) constant returns(uint256)
func (_Leagues *LeaguesSession) GetSpeed(encodedSkills *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetSpeed(&_Leagues.CallOpts, encodedSkills)
}

// GetSpeed is a free data retrieval call binding the contract method 0x4b93f753.
//
// Solidity: function getSpeed(uint256 encodedSkills) constant returns(uint256)
func (_Leagues *LeaguesCallerSession) GetSpeed(encodedSkills *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetSpeed(&_Leagues.CallOpts, encodedSkills)
}

// GetSumOfSkills is a free data retrieval call binding the contract method 0x1060c9c2.
//
// Solidity: function getSumOfSkills(uint256 encodedSkills) constant returns(uint256)
func (_Leagues *LeaguesCaller) GetSumOfSkills(opts *bind.CallOpts, encodedSkills *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "getSumOfSkills", encodedSkills)
	return *ret0, err
}

// GetSumOfSkills is a free data retrieval call binding the contract method 0x1060c9c2.
//
// Solidity: function getSumOfSkills(uint256 encodedSkills) constant returns(uint256)
func (_Leagues *LeaguesSession) GetSumOfSkills(encodedSkills *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetSumOfSkills(&_Leagues.CallOpts, encodedSkills)
}

// GetSumOfSkills is a free data retrieval call binding the contract method 0x1060c9c2.
//
// Solidity: function getSumOfSkills(uint256 encodedSkills) constant returns(uint256)
func (_Leagues *LeaguesCallerSession) GetSumOfSkills(encodedSkills *big.Int) (*big.Int, error) {
	return _Leagues.Contract.GetSumOfSkills(&_Leagues.CallOpts, encodedSkills)
}

// GetTeamsInMatch is a free data retrieval call binding the contract method 0x032324c8.
//
// Solidity: function getTeamsInMatch(uint8 matchday, uint8 matchIdxInDay) constant returns(uint8 homeIdx, uint8 visitorIdx)
func (_Leagues *LeaguesCaller) GetTeamsInMatch(opts *bind.CallOpts, matchday uint8, matchIdxInDay uint8) (struct {
	HomeIdx    uint8
	VisitorIdx uint8
}, error) {
	ret := new(struct {
		HomeIdx    uint8
		VisitorIdx uint8
	})
	out := ret
	err := _Leagues.contract.Call(opts, out, "getTeamsInMatch", matchday, matchIdxInDay)
	return *ret, err
}

// GetTeamsInMatch is a free data retrieval call binding the contract method 0x032324c8.
//
// Solidity: function getTeamsInMatch(uint8 matchday, uint8 matchIdxInDay) constant returns(uint8 homeIdx, uint8 visitorIdx)
func (_Leagues *LeaguesSession) GetTeamsInMatch(matchday uint8, matchIdxInDay uint8) (struct {
	HomeIdx    uint8
	VisitorIdx uint8
}, error) {
	return _Leagues.Contract.GetTeamsInMatch(&_Leagues.CallOpts, matchday, matchIdxInDay)
}

// GetTeamsInMatch is a free data retrieval call binding the contract method 0x032324c8.
//
// Solidity: function getTeamsInMatch(uint8 matchday, uint8 matchIdxInDay) constant returns(uint8 homeIdx, uint8 visitorIdx)
func (_Leagues *LeaguesCallerSession) GetTeamsInMatch(matchday uint8, matchIdxInDay uint8) (struct {
	HomeIdx    uint8
	VisitorIdx uint8
}, error) {
	return _Leagues.Contract.GetTeamsInMatch(&_Leagues.CallOpts, matchday, matchIdxInDay)
}

// IsBotTeam is a free data retrieval call binding the contract method 0x8cc9a8d5.
//
// Solidity: function isBotTeam(uint256 teamId) constant returns(bool)
func (_Leagues *LeaguesCaller) IsBotTeam(opts *bind.CallOpts, teamId *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "isBotTeam", teamId)
	return *ret0, err
}

// IsBotTeam is a free data retrieval call binding the contract method 0x8cc9a8d5.
//
// Solidity: function isBotTeam(uint256 teamId) constant returns(bool)
func (_Leagues *LeaguesSession) IsBotTeam(teamId *big.Int) (bool, error) {
	return _Leagues.Contract.IsBotTeam(&_Leagues.CallOpts, teamId)
}

// IsBotTeam is a free data retrieval call binding the contract method 0x8cc9a8d5.
//
// Solidity: function isBotTeam(uint256 teamId) constant returns(bool)
func (_Leagues *LeaguesCallerSession) IsBotTeam(teamId *big.Int) (bool, error) {
	return _Leagues.Contract.IsBotTeam(&_Leagues.CallOpts, teamId)
}

// IsBotTeamInCountry is a free data retrieval call binding the contract method 0x80bac709.
//
// Solidity: function isBotTeamInCountry(uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) constant returns(bool)
func (_Leagues *LeaguesCaller) IsBotTeamInCountry(opts *bind.CallOpts, timeZone uint8, countryIdxInTZ *big.Int, teamIdxInCountry *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "isBotTeamInCountry", timeZone, countryIdxInTZ, teamIdxInCountry)
	return *ret0, err
}

// IsBotTeamInCountry is a free data retrieval call binding the contract method 0x80bac709.
//
// Solidity: function isBotTeamInCountry(uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) constant returns(bool)
func (_Leagues *LeaguesSession) IsBotTeamInCountry(timeZone uint8, countryIdxInTZ *big.Int, teamIdxInCountry *big.Int) (bool, error) {
	return _Leagues.Contract.IsBotTeamInCountry(&_Leagues.CallOpts, timeZone, countryIdxInTZ, teamIdxInCountry)
}

// IsBotTeamInCountry is a free data retrieval call binding the contract method 0x80bac709.
//
// Solidity: function isBotTeamInCountry(uint8 timeZone, uint256 countryIdxInTZ, uint256 teamIdxInCountry) constant returns(bool)
func (_Leagues *LeaguesCallerSession) IsBotTeamInCountry(timeZone uint8, countryIdxInTZ *big.Int, teamIdxInCountry *big.Int) (bool, error) {
	return _Leagues.Contract.IsBotTeamInCountry(&_Leagues.CallOpts, timeZone, countryIdxInTZ, teamIdxInCountry)
}

// IsFreeShirt is a free data retrieval call binding the contract method 0x963fcc80.
//
// Solidity: function isFreeShirt(uint256 teamId, uint8 shirtNum) constant returns(bool)
func (_Leagues *LeaguesCaller) IsFreeShirt(opts *bind.CallOpts, teamId *big.Int, shirtNum uint8) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "isFreeShirt", teamId, shirtNum)
	return *ret0, err
}

// IsFreeShirt is a free data retrieval call binding the contract method 0x963fcc80.
//
// Solidity: function isFreeShirt(uint256 teamId, uint8 shirtNum) constant returns(bool)
func (_Leagues *LeaguesSession) IsFreeShirt(teamId *big.Int, shirtNum uint8) (bool, error) {
	return _Leagues.Contract.IsFreeShirt(&_Leagues.CallOpts, teamId, shirtNum)
}

// IsFreeShirt is a free data retrieval call binding the contract method 0x963fcc80.
//
// Solidity: function isFreeShirt(uint256 teamId, uint8 shirtNum) constant returns(bool)
func (_Leagues *LeaguesCallerSession) IsFreeShirt(teamId *big.Int, shirtNum uint8) (bool, error) {
	return _Leagues.Contract.IsFreeShirt(&_Leagues.CallOpts, teamId, shirtNum)
}

// IsVirtualPlayer is a free data retrieval call binding the contract method 0xb32aa2c1.
//
// Solidity: function isVirtualPlayer(uint256 playerId) constant returns(bool)
func (_Leagues *LeaguesCaller) IsVirtualPlayer(opts *bind.CallOpts, playerId *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "isVirtualPlayer", playerId)
	return *ret0, err
}

// IsVirtualPlayer is a free data retrieval call binding the contract method 0xb32aa2c1.
//
// Solidity: function isVirtualPlayer(uint256 playerId) constant returns(bool)
func (_Leagues *LeaguesSession) IsVirtualPlayer(playerId *big.Int) (bool, error) {
	return _Leagues.Contract.IsVirtualPlayer(&_Leagues.CallOpts, playerId)
}

// IsVirtualPlayer is a free data retrieval call binding the contract method 0xb32aa2c1.
//
// Solidity: function isVirtualPlayer(uint256 playerId) constant returns(bool)
func (_Leagues *LeaguesCallerSession) IsVirtualPlayer(playerId *big.Int) (bool, error) {
	return _Leagues.Contract.IsVirtualPlayer(&_Leagues.CallOpts, playerId)
}

// PlayerExists is a free data retrieval call binding the contract method 0xbc1a97c1.
//
// Solidity: function playerExists(uint256 playerId) constant returns(bool)
func (_Leagues *LeaguesCaller) PlayerExists(opts *bind.CallOpts, playerId *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "playerExists", playerId)
	return *ret0, err
}

// PlayerExists is a free data retrieval call binding the contract method 0xbc1a97c1.
//
// Solidity: function playerExists(uint256 playerId) constant returns(bool)
func (_Leagues *LeaguesSession) PlayerExists(playerId *big.Int) (bool, error) {
	return _Leagues.Contract.PlayerExists(&_Leagues.CallOpts, playerId)
}

// PlayerExists is a free data retrieval call binding the contract method 0xbc1a97c1.
//
// Solidity: function playerExists(uint256 playerId) constant returns(bool)
func (_Leagues *LeaguesCallerSession) PlayerExists(playerId *big.Int) (bool, error) {
	return _Leagues.Contract.PlayerExists(&_Leagues.CallOpts, playerId)
}

// SetCurrentShirtNum is a free data retrieval call binding the contract method 0x4db989fd.
//
// Solidity: function setCurrentShirtNum(uint256 state, uint8 currentShirtNum) constant returns(uint256)
func (_Leagues *LeaguesCaller) SetCurrentShirtNum(opts *bind.CallOpts, state *big.Int, currentShirtNum uint8) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "setCurrentShirtNum", state, currentShirtNum)
	return *ret0, err
}

// SetCurrentShirtNum is a free data retrieval call binding the contract method 0x4db989fd.
//
// Solidity: function setCurrentShirtNum(uint256 state, uint8 currentShirtNum) constant returns(uint256)
func (_Leagues *LeaguesSession) SetCurrentShirtNum(state *big.Int, currentShirtNum uint8) (*big.Int, error) {
	return _Leagues.Contract.SetCurrentShirtNum(&_Leagues.CallOpts, state, currentShirtNum)
}

// SetCurrentShirtNum is a free data retrieval call binding the contract method 0x4db989fd.
//
// Solidity: function setCurrentShirtNum(uint256 state, uint8 currentShirtNum) constant returns(uint256)
func (_Leagues *LeaguesCallerSession) SetCurrentShirtNum(state *big.Int, currentShirtNum uint8) (*big.Int, error) {
	return _Leagues.Contract.SetCurrentShirtNum(&_Leagues.CallOpts, state, currentShirtNum)
}

// SetCurrentTeamId is a free data retrieval call binding the contract method 0xc37b1c25.
//
// Solidity: function setCurrentTeamId(uint256 playerState, uint256 teamId) constant returns(uint256)
func (_Leagues *LeaguesCaller) SetCurrentTeamId(opts *bind.CallOpts, playerState *big.Int, teamId *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "setCurrentTeamId", playerState, teamId)
	return *ret0, err
}

// SetCurrentTeamId is a free data retrieval call binding the contract method 0xc37b1c25.
//
// Solidity: function setCurrentTeamId(uint256 playerState, uint256 teamId) constant returns(uint256)
func (_Leagues *LeaguesSession) SetCurrentTeamId(playerState *big.Int, teamId *big.Int) (*big.Int, error) {
	return _Leagues.Contract.SetCurrentTeamId(&_Leagues.CallOpts, playerState, teamId)
}

// SetCurrentTeamId is a free data retrieval call binding the contract method 0xc37b1c25.
//
// Solidity: function setCurrentTeamId(uint256 playerState, uint256 teamId) constant returns(uint256)
func (_Leagues *LeaguesCallerSession) SetCurrentTeamId(playerState *big.Int, teamId *big.Int) (*big.Int, error) {
	return _Leagues.Contract.SetCurrentTeamId(&_Leagues.CallOpts, playerState, teamId)
}

// SetLastSaleBlock is a free data retrieval call binding the contract method 0xaf76cd01.
//
// Solidity: function setLastSaleBlock(uint256 state, uint256 lastSaleBlock) constant returns(uint256)
func (_Leagues *LeaguesCaller) SetLastSaleBlock(opts *bind.CallOpts, state *big.Int, lastSaleBlock *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "setLastSaleBlock", state, lastSaleBlock)
	return *ret0, err
}

// SetLastSaleBlock is a free data retrieval call binding the contract method 0xaf76cd01.
//
// Solidity: function setLastSaleBlock(uint256 state, uint256 lastSaleBlock) constant returns(uint256)
func (_Leagues *LeaguesSession) SetLastSaleBlock(state *big.Int, lastSaleBlock *big.Int) (*big.Int, error) {
	return _Leagues.Contract.SetLastSaleBlock(&_Leagues.CallOpts, state, lastSaleBlock)
}

// SetLastSaleBlock is a free data retrieval call binding the contract method 0xaf76cd01.
//
// Solidity: function setLastSaleBlock(uint256 state, uint256 lastSaleBlock) constant returns(uint256)
func (_Leagues *LeaguesCallerSession) SetLastSaleBlock(state *big.Int, lastSaleBlock *big.Int) (*big.Int, error) {
	return _Leagues.Contract.SetLastSaleBlock(&_Leagues.CallOpts, state, lastSaleBlock)
}

// SetPrevPlayerTeamId is a free data retrieval call binding the contract method 0x37a86302.
//
// Solidity: function setPrevPlayerTeamId(uint256 state, uint256 value) constant returns(uint256)
func (_Leagues *LeaguesCaller) SetPrevPlayerTeamId(opts *bind.CallOpts, state *big.Int, value *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "setPrevPlayerTeamId", state, value)
	return *ret0, err
}

// SetPrevPlayerTeamId is a free data retrieval call binding the contract method 0x37a86302.
//
// Solidity: function setPrevPlayerTeamId(uint256 state, uint256 value) constant returns(uint256)
func (_Leagues *LeaguesSession) SetPrevPlayerTeamId(state *big.Int, value *big.Int) (*big.Int, error) {
	return _Leagues.Contract.SetPrevPlayerTeamId(&_Leagues.CallOpts, state, value)
}

// SetPrevPlayerTeamId is a free data retrieval call binding the contract method 0x37a86302.
//
// Solidity: function setPrevPlayerTeamId(uint256 state, uint256 value) constant returns(uint256)
func (_Leagues *LeaguesCallerSession) SetPrevPlayerTeamId(state *big.Int, value *big.Int) (*big.Int, error) {
	return _Leagues.Contract.SetPrevPlayerTeamId(&_Leagues.CallOpts, state, value)
}

// TeamExists is a free data retrieval call binding the contract method 0x98981756.
//
// Solidity: function teamExists(uint256 teamId) constant returns(bool)
func (_Leagues *LeaguesCaller) TeamExists(opts *bind.CallOpts, teamId *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Leagues.contract.Call(opts, out, "teamExists", teamId)
	return *ret0, err
}

// TeamExists is a free data retrieval call binding the contract method 0x98981756.
//
// Solidity: function teamExists(uint256 teamId) constant returns(bool)
func (_Leagues *LeaguesSession) TeamExists(teamId *big.Int) (bool, error) {
	return _Leagues.Contract.TeamExists(&_Leagues.CallOpts, teamId)
}

// TeamExists is a free data retrieval call binding the contract method 0x98981756.
//
// Solidity: function teamExists(uint256 teamId) constant returns(bool)
func (_Leagues *LeaguesCallerSession) TeamExists(teamId *big.Int) (bool, error) {
	return _Leagues.Contract.TeamExists(&_Leagues.CallOpts, teamId)
}

// Init is a paid mutator transaction binding the contract method 0xe1c7392a.
//
// Solidity: function init() returns()
func (_Leagues *LeaguesTransactor) Init(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Leagues.contract.Transact(opts, "init")
}

// Init is a paid mutator transaction binding the contract method 0xe1c7392a.
//
// Solidity: function init() returns()
func (_Leagues *LeaguesSession) Init() (*types.Transaction, error) {
	return _Leagues.Contract.Init(&_Leagues.TransactOpts)
}

// Init is a paid mutator transaction binding the contract method 0xe1c7392a.
//
// Solidity: function init() returns()
func (_Leagues *LeaguesTransactorSession) Init() (*types.Transaction, error) {
	return _Leagues.Contract.Init(&_Leagues.TransactOpts)
}

// InitSingleTZ is a paid mutator transaction binding the contract method 0xa3ceb703.
//
// Solidity: function initSingleTZ(uint8 tz) returns()
func (_Leagues *LeaguesTransactor) InitSingleTZ(opts *bind.TransactOpts, tz uint8) (*types.Transaction, error) {
	return _Leagues.contract.Transact(opts, "initSingleTZ", tz)
}

// InitSingleTZ is a paid mutator transaction binding the contract method 0xa3ceb703.
//
// Solidity: function initSingleTZ(uint8 tz) returns()
func (_Leagues *LeaguesSession) InitSingleTZ(tz uint8) (*types.Transaction, error) {
	return _Leagues.Contract.InitSingleTZ(&_Leagues.TransactOpts, tz)
}

// InitSingleTZ is a paid mutator transaction binding the contract method 0xa3ceb703.
//
// Solidity: function initSingleTZ(uint8 tz) returns()
func (_Leagues *LeaguesTransactorSession) InitSingleTZ(tz uint8) (*types.Transaction, error) {
	return _Leagues.Contract.InitSingleTZ(&_Leagues.TransactOpts, tz)
}

// SetActionsRoot is a paid mutator transaction binding the contract method 0xdba6319e.
//
// Solidity: function setActionsRoot(uint8 timeZone, bytes32 root) returns(uint256)
func (_Leagues *LeaguesTransactor) SetActionsRoot(opts *bind.TransactOpts, timeZone uint8, root [32]byte) (*types.Transaction, error) {
	return _Leagues.contract.Transact(opts, "setActionsRoot", timeZone, root)
}

// SetActionsRoot is a paid mutator transaction binding the contract method 0xdba6319e.
//
// Solidity: function setActionsRoot(uint8 timeZone, bytes32 root) returns(uint256)
func (_Leagues *LeaguesSession) SetActionsRoot(timeZone uint8, root [32]byte) (*types.Transaction, error) {
	return _Leagues.Contract.SetActionsRoot(&_Leagues.TransactOpts, timeZone, root)
}

// SetActionsRoot is a paid mutator transaction binding the contract method 0xdba6319e.
//
// Solidity: function setActionsRoot(uint8 timeZone, bytes32 root) returns(uint256)
func (_Leagues *LeaguesTransactorSession) SetActionsRoot(timeZone uint8, root [32]byte) (*types.Transaction, error) {
	return _Leagues.Contract.SetActionsRoot(&_Leagues.TransactOpts, timeZone, root)
}

// SetEngineAdress is a paid mutator transaction binding the contract method 0x058672f9.
//
// Solidity: function setEngineAdress(address addr) returns()
func (_Leagues *LeaguesTransactor) SetEngineAdress(opts *bind.TransactOpts, addr common.Address) (*types.Transaction, error) {
	return _Leagues.contract.Transact(opts, "setEngineAdress", addr)
}

// SetEngineAdress is a paid mutator transaction binding the contract method 0x058672f9.
//
// Solidity: function setEngineAdress(address addr) returns()
func (_Leagues *LeaguesSession) SetEngineAdress(addr common.Address) (*types.Transaction, error) {
	return _Leagues.Contract.SetEngineAdress(&_Leagues.TransactOpts, addr)
}

// SetEngineAdress is a paid mutator transaction binding the contract method 0x058672f9.
//
// Solidity: function setEngineAdress(address addr) returns()
func (_Leagues *LeaguesTransactorSession) SetEngineAdress(addr common.Address) (*types.Transaction, error) {
	return _Leagues.Contract.SetEngineAdress(&_Leagues.TransactOpts, addr)
}

// SetSkillsRoot is a paid mutator transaction binding the contract method 0xec1c5423.
//
// Solidity: function setSkillsRoot(uint8 tz, bytes32 root) returns(uint256)
func (_Leagues *LeaguesTransactor) SetSkillsRoot(opts *bind.TransactOpts, tz uint8, root [32]byte) (*types.Transaction, error) {
	return _Leagues.contract.Transact(opts, "setSkillsRoot", tz, root)
}

// SetSkillsRoot is a paid mutator transaction binding the contract method 0xec1c5423.
//
// Solidity: function setSkillsRoot(uint8 tz, bytes32 root) returns(uint256)
func (_Leagues *LeaguesSession) SetSkillsRoot(tz uint8, root [32]byte) (*types.Transaction, error) {
	return _Leagues.Contract.SetSkillsRoot(&_Leagues.TransactOpts, tz, root)
}

// SetSkillsRoot is a paid mutator transaction binding the contract method 0xec1c5423.
//
// Solidity: function setSkillsRoot(uint8 tz, bytes32 root) returns(uint256)
func (_Leagues *LeaguesTransactorSession) SetSkillsRoot(tz uint8, root [32]byte) (*types.Transaction, error) {
	return _Leagues.Contract.SetSkillsRoot(&_Leagues.TransactOpts, tz, root)
}

// TransferFirstBotToAddr is a paid mutator transaction binding the contract method 0x3c2eb360.
//
// Solidity: function transferFirstBotToAddr(uint8 timeZone, uint256 countryIdxInTZ, address addr) returns()
func (_Leagues *LeaguesTransactor) TransferFirstBotToAddr(opts *bind.TransactOpts, timeZone uint8, countryIdxInTZ *big.Int, addr common.Address) (*types.Transaction, error) {
	return _Leagues.contract.Transact(opts, "transferFirstBotToAddr", timeZone, countryIdxInTZ, addr)
}

// TransferFirstBotToAddr is a paid mutator transaction binding the contract method 0x3c2eb360.
//
// Solidity: function transferFirstBotToAddr(uint8 timeZone, uint256 countryIdxInTZ, address addr) returns()
func (_Leagues *LeaguesSession) TransferFirstBotToAddr(timeZone uint8, countryIdxInTZ *big.Int, addr common.Address) (*types.Transaction, error) {
	return _Leagues.Contract.TransferFirstBotToAddr(&_Leagues.TransactOpts, timeZone, countryIdxInTZ, addr)
}

// TransferFirstBotToAddr is a paid mutator transaction binding the contract method 0x3c2eb360.
//
// Solidity: function transferFirstBotToAddr(uint8 timeZone, uint256 countryIdxInTZ, address addr) returns()
func (_Leagues *LeaguesTransactorSession) TransferFirstBotToAddr(timeZone uint8, countryIdxInTZ *big.Int, addr common.Address) (*types.Transaction, error) {
	return _Leagues.Contract.TransferFirstBotToAddr(&_Leagues.TransactOpts, timeZone, countryIdxInTZ, addr)
}

// TransferPlayer is a paid mutator transaction binding the contract method 0x800257d5.
//
// Solidity: function transferPlayer(uint256 playerId, uint256 teamIdTarget) returns()
func (_Leagues *LeaguesTransactor) TransferPlayer(opts *bind.TransactOpts, playerId *big.Int, teamIdTarget *big.Int) (*types.Transaction, error) {
	return _Leagues.contract.Transact(opts, "transferPlayer", playerId, teamIdTarget)
}

// TransferPlayer is a paid mutator transaction binding the contract method 0x800257d5.
//
// Solidity: function transferPlayer(uint256 playerId, uint256 teamIdTarget) returns()
func (_Leagues *LeaguesSession) TransferPlayer(playerId *big.Int, teamIdTarget *big.Int) (*types.Transaction, error) {
	return _Leagues.Contract.TransferPlayer(&_Leagues.TransactOpts, playerId, teamIdTarget)
}

// TransferPlayer is a paid mutator transaction binding the contract method 0x800257d5.
//
// Solidity: function transferPlayer(uint256 playerId, uint256 teamIdTarget) returns()
func (_Leagues *LeaguesTransactorSession) TransferPlayer(playerId *big.Int, teamIdTarget *big.Int) (*types.Transaction, error) {
	return _Leagues.Contract.TransferPlayer(&_Leagues.TransactOpts, playerId, teamIdTarget)
}

// TransferTeam is a paid mutator transaction binding the contract method 0xe945e96a.
//
// Solidity: function transferTeam(uint256 teamId, address addr) returns()
func (_Leagues *LeaguesTransactor) TransferTeam(opts *bind.TransactOpts, teamId *big.Int, addr common.Address) (*types.Transaction, error) {
	return _Leagues.contract.Transact(opts, "transferTeam", teamId, addr)
}

// TransferTeam is a paid mutator transaction binding the contract method 0xe945e96a.
//
// Solidity: function transferTeam(uint256 teamId, address addr) returns()
func (_Leagues *LeaguesSession) TransferTeam(teamId *big.Int, addr common.Address) (*types.Transaction, error) {
	return _Leagues.Contract.TransferTeam(&_Leagues.TransactOpts, teamId, addr)
}

// TransferTeam is a paid mutator transaction binding the contract method 0xe945e96a.
//
// Solidity: function transferTeam(uint256 teamId, address addr) returns()
func (_Leagues *LeaguesTransactorSession) TransferTeam(teamId *big.Int, addr common.Address) (*types.Transaction, error) {
	return _Leagues.Contract.TransferTeam(&_Leagues.TransactOpts, teamId, addr)
}

// LeaguesDivisionCreationIterator is returned from FilterDivisionCreation and is used to iterate over the raw logs and unpacked data for DivisionCreation events raised by the Leagues contract.
type LeaguesDivisionCreationIterator struct {
	Event *LeaguesDivisionCreation // Event containing the contract specifics and raw log

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
func (it *LeaguesDivisionCreationIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LeaguesDivisionCreation)
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
		it.Event = new(LeaguesDivisionCreation)
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
func (it *LeaguesDivisionCreationIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LeaguesDivisionCreationIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LeaguesDivisionCreation represents a DivisionCreation event raised by the Leagues contract.
type LeaguesDivisionCreation struct {
	Timezone             uint8
	CountryIdxInTZ       *big.Int
	DivisionIdxInCountry *big.Int
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterDivisionCreation is a free log retrieval operation binding the contract event 0xc5d195855a200aa90e2052bcc795cedbc84c2a26556b1d5113b5a30c96003a0b.
//
// Solidity: event DivisionCreation(uint8 timezone, uint256 countryIdxInTZ, uint256 divisionIdxInCountry)
func (_Leagues *LeaguesFilterer) FilterDivisionCreation(opts *bind.FilterOpts) (*LeaguesDivisionCreationIterator, error) {

	logs, sub, err := _Leagues.contract.FilterLogs(opts, "DivisionCreation")
	if err != nil {
		return nil, err
	}
	return &LeaguesDivisionCreationIterator{contract: _Leagues.contract, event: "DivisionCreation", logs: logs, sub: sub}, nil
}

// WatchDivisionCreation is a free log subscription operation binding the contract event 0xc5d195855a200aa90e2052bcc795cedbc84c2a26556b1d5113b5a30c96003a0b.
//
// Solidity: event DivisionCreation(uint8 timezone, uint256 countryIdxInTZ, uint256 divisionIdxInCountry)
func (_Leagues *LeaguesFilterer) WatchDivisionCreation(opts *bind.WatchOpts, sink chan<- *LeaguesDivisionCreation) (event.Subscription, error) {

	logs, sub, err := _Leagues.contract.WatchLogs(opts, "DivisionCreation")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LeaguesDivisionCreation)
				if err := _Leagues.contract.UnpackLog(event, "DivisionCreation", log); err != nil {
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

// LeaguesPlayerStateChangeIterator is returned from FilterPlayerStateChange and is used to iterate over the raw logs and unpacked data for PlayerStateChange events raised by the Leagues contract.
type LeaguesPlayerStateChangeIterator struct {
	Event *LeaguesPlayerStateChange // Event containing the contract specifics and raw log

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
func (it *LeaguesPlayerStateChangeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LeaguesPlayerStateChange)
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
		it.Event = new(LeaguesPlayerStateChange)
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
func (it *LeaguesPlayerStateChangeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LeaguesPlayerStateChangeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LeaguesPlayerStateChange represents a PlayerStateChange event raised by the Leagues contract.
type LeaguesPlayerStateChange struct {
	PlayerId *big.Int
	State    *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterPlayerStateChange is a free log retrieval operation binding the contract event 0x65a4d4a8a0afb474d2e9465a6a1a41bb88fd04f41152ba070421f1b3771f15d8.
//
// Solidity: event PlayerStateChange(uint256 playerId, uint256 state)
func (_Leagues *LeaguesFilterer) FilterPlayerStateChange(opts *bind.FilterOpts) (*LeaguesPlayerStateChangeIterator, error) {

	logs, sub, err := _Leagues.contract.FilterLogs(opts, "PlayerStateChange")
	if err != nil {
		return nil, err
	}
	return &LeaguesPlayerStateChangeIterator{contract: _Leagues.contract, event: "PlayerStateChange", logs: logs, sub: sub}, nil
}

// WatchPlayerStateChange is a free log subscription operation binding the contract event 0x65a4d4a8a0afb474d2e9465a6a1a41bb88fd04f41152ba070421f1b3771f15d8.
//
// Solidity: event PlayerStateChange(uint256 playerId, uint256 state)
func (_Leagues *LeaguesFilterer) WatchPlayerStateChange(opts *bind.WatchOpts, sink chan<- *LeaguesPlayerStateChange) (event.Subscription, error) {

	logs, sub, err := _Leagues.contract.WatchLogs(opts, "PlayerStateChange")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LeaguesPlayerStateChange)
				if err := _Leagues.contract.UnpackLog(event, "PlayerStateChange", log); err != nil {
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

// LeaguesTeamTransferIterator is returned from FilterTeamTransfer and is used to iterate over the raw logs and unpacked data for TeamTransfer events raised by the Leagues contract.
type LeaguesTeamTransferIterator struct {
	Event *LeaguesTeamTransfer // Event containing the contract specifics and raw log

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
func (it *LeaguesTeamTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LeaguesTeamTransfer)
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
		it.Event = new(LeaguesTeamTransfer)
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
func (it *LeaguesTeamTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LeaguesTeamTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LeaguesTeamTransfer represents a TeamTransfer event raised by the Leagues contract.
type LeaguesTeamTransfer struct {
	TeamId *big.Int
	To     common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterTeamTransfer is a free log retrieval operation binding the contract event 0x77b66eb1e6d2bc131b79be4213ae7f08f29351c01060e10bcc0302278067bf38.
//
// Solidity: event TeamTransfer(uint256 teamId, address to)
func (_Leagues *LeaguesFilterer) FilterTeamTransfer(opts *bind.FilterOpts) (*LeaguesTeamTransferIterator, error) {

	logs, sub, err := _Leagues.contract.FilterLogs(opts, "TeamTransfer")
	if err != nil {
		return nil, err
	}
	return &LeaguesTeamTransferIterator{contract: _Leagues.contract, event: "TeamTransfer", logs: logs, sub: sub}, nil
}

// WatchTeamTransfer is a free log subscription operation binding the contract event 0x77b66eb1e6d2bc131b79be4213ae7f08f29351c01060e10bcc0302278067bf38.
//
// Solidity: event TeamTransfer(uint256 teamId, address to)
func (_Leagues *LeaguesFilterer) WatchTeamTransfer(opts *bind.WatchOpts, sink chan<- *LeaguesTeamTransfer) (event.Subscription, error) {

	logs, sub, err := _Leagues.contract.WatchLogs(opts, "TeamTransfer")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LeaguesTeamTransfer)
				if err := _Leagues.contract.UnpackLog(event, "TeamTransfer", log); err != nil {
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
