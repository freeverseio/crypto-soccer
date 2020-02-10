pragma solidity >=0.5.12 <0.6.2;

import "./EncodingSkillsSetters.sol";
/**
 * @title Creation of all game assets via creation of timezones, countries and divisions
 * @dev Timezones range from 1 to 24, with timeZone = 0 being null.
 */

contract Shop is EncodingSkillsSetters{

    event ItemOffered(
        uint256 itemId,
        uint256 countriesRoot,
        uint256 championshipsRoot,
        uint256 teamsRoot,
        uint16 itemsRemaining,
        uint64 encodedBoost,
        uint8 matchesDuration,
        uint8 onlyTopInChampioniship,
        string uri
    );

    uint8 constant public SK_SHO = 0;
    uint8 constant public SK_SPE = 1;
    uint8 constant public SK_PAS = 2;
    uint8 constant public SK_DEF = 3;
    uint8 constant public SK_END = 4;
    uint8 constant public N_SKILLS = 5;

    struct ShopItem {
        // boosts from [0,..4] -> in percentage, order: shoot, speed, pass, defence, endurance
        // boosts[5] in units, for potential
        uint256 countriesRoot;
        uint256 championshipsRoot;
        uint256 teamsRoot;
        uint16 itemsRemaining;
        uint64 encodedBoost;
        uint8 matchesDuration;
        uint8 onlyTopInChampioniship;
        string uri;
    }        

    ShopItem[] private _shopItems;
    uint8 maxPercentPerSkill = 50;
    uint8 maxIncreasePotential = 1;

    function getEncodedBoost(uint256 itemId) public view returns (uint64) { return _shopItems[itemId].encodedBoost; }
    function getCountriesRoot(uint256 itemId) public view returns (uint256) { return _shopItems[itemId].countriesRoot; }
    function getChampionshipsRoot(uint256 itemId) public view returns (uint256) { return _shopItems[itemId].championshipsRoot; }
    function getTeamsRoot(uint256 itemId) public view returns (uint256) { return _shopItems[itemId].teamsRoot; }
    function getItemsRemaining(uint256 itemId) public view returns (uint16) { return _shopItems[itemId].itemsRemaining; }
    function getMatchesDuration(uint256 itemId) public view returns (uint8) { return _shopItems[itemId].matchesDuration; }
    function getOnlyTopInChampioniship(uint256 itemId) public view returns (uint8) { return _shopItems[itemId].onlyTopInChampioniship; }
    function getUri(uint256 itemId) public view returns (string memory) { return _shopItems[itemId].uri; }
    function setMaxPercentPerSkill(uint8 newVal) public { maxPercentPerSkill = newVal; } 
    function setMaxIncreasePotential(uint8 newVal) public { maxIncreasePotential = newVal; } 

    function init() public {
        _shopItems.push(ShopItem(0, 0, 0, 0, 0, 0, 0, ""));
    }

    function offerItem(
        uint8[N_SKILLS+1] memory skillsBoost,
        uint256 countriesRoot,
        uint256 championshipsRoot,
        uint256 teamsRoot,
        uint16 itemsRemaining,
        uint8 matchesDuration,
        uint8 onlyTopInChampioniship,
        string memory uri
    ) public {
        for (uint8 sk = 0; sk < N_SKILLS; sk++) {
            require(skillsBoost[sk] <= maxPercentPerSkill, "cannot offer items that boost one skill so much");
        }
        require(skillsBoost[N_SKILLS] <= maxIncreasePotential, "cannot offer items that boost potential so much");
        require(itemsRemaining > 0, "cannot offer 0 items");
        require(matchesDuration > 0, "cannot offer items that last 0 games");
        uint64 encodedBoost = encodeBoosts(skillsBoost);
        _shopItems.push(ShopItem(
            countriesRoot,
            championshipsRoot,
            teamsRoot,
            itemsRemaining,
            encodedBoost,
            matchesDuration,
            onlyTopInChampioniship,
            uri
        ));
        emit ItemOffered(
            _shopItems.length - 1,
            countriesRoot,
            championshipsRoot,
            teamsRoot,
            itemsRemaining,
            encodedBoost,
            matchesDuration,
            onlyTopInChampioniship,
            uri
        );
    }
    
    // bits: 5 skills * 6b per skill (max 64) + 2b for potential = 32b
    function encodeBoosts(uint8[N_SKILLS+1] memory skillsBoost) public pure returns(uint64 encoded) {
        require(skillsBoost[N_SKILLS] < 4, "cannot offer items that boost potential so much");
        for (uint8 sk = 0; sk <= N_SKILLS; sk++) {
            encoded |= (uint64(skillsBoost[sk]) << 6*sk);
        }
    }
    
    function decodeBoosts(uint64 encoded) public pure returns(uint8[N_SKILLS+1] memory skillsBoost) {
        for (uint8 sk = 0; sk <= N_SKILLS; sk++) {
            skillsBoost[sk] = uint8((encoded >> 6*sk) & 63);
        }
    }
}
