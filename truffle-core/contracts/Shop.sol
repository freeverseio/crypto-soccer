pragma solidity >=0.4.21 <0.6.0;

import "./EncodingSkillsSetters.sol";
/**
 * @title Creation of all game assets via creation of timezones, countries and divisions
 * @dev Timezones range from 1 to 24, with timeZone = 0 being null.
 */

contract Shop is EncodingSkillsSetters{

    event ItemOffered(uint256 itemId, uint64[] leagueIds);

    uint16 constant public MAX_LEAGUES_PER_OFFER = 512;
    uint8 constant public SK_SHO = 0;
    uint8 constant public SK_SPE = 1;
    uint8 constant public SK_PAS = 2;
    uint8 constant public SK_DEF = 3;
    uint8 constant public SK_END = 4;
    
    struct ShopItem {
        uint256 skillsBoost;
        bytes32 championshipsHash;
        uint16 stock;
        uint8 matchesDuration;
        string uri;
    }        

    ShopItem[] private _shopItems;

    function getSkillsBoost(uint256 itemId) public view returns (uint256) { return _shopItems[itemId].skillsBoost; }
    function getMatchesDuration(uint256 itemId) public view returns (uint8) { return _shopItems[itemId].matchesDuration; }
    function getStock(uint256 itemId) public view returns (uint16) { return _shopItems[itemId].stock; }
    function getChampionshipsHash(uint256 itemId) public view returns (bytes32) { return _shopItems[itemId].championshipsHash; }
    function getUri(uint256 itemId) public view returns (string memory) { return _shopItems[itemId].uri; }

    // 1 round = 2 weeks, 1 year = 25 rounds, 100 years = 2500 rounds => 12bit
    // leagueId: tz (5), country (10), leagueIdxInCountry (28), round(12) = 55 bit => 64 
    function offerItem(uint256 skillsBoost, uint16 stock, uint8 matchesDuration, string memory uri, uint64[] memory leagueIds) public {
        require(leagueIds.length < MAX_LEAGUES_PER_OFFER, "too many leafs in tree");
        _shopItems.push(ShopItem(skillsBoost, merkleRootTemp(leagueIds), stock, matchesDuration, uri));
        emit ItemOffered(_shopItems.length - 1, leagueIds);
    }
    
    // TODO: this is a temp replacement for merkleRoot. 
    function merkleRootTemp(uint64[] memory leafs) public pure returns (bytes32 hash) {
        for (uint16 l = 0; l < leafs.length; l++) {
            hash = keccak256(abi.encode(hash, leafs[l]));
        } 
    }
    
    function encodeSkillsBoost(uint16[5] memory skillsVec) public pure returns (uint256 skillsBoost) {
        skillsBoost = setShoot(skillsBoost, skillsVec[SK_SHO]);
        skillsBoost = setSpeed(skillsBoost, skillsVec[SK_SPE]);
        skillsBoost = setPass(skillsBoost, skillsVec[SK_PAS]);
        skillsBoost = setDefence(skillsBoost, skillsVec[SK_DEF]);
        skillsBoost = setEndurance(skillsBoost, skillsVec[SK_END]);
    }
    
}
