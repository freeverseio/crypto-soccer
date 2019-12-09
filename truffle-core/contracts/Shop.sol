pragma solidity >=0.4.21 <0.6.0;

// import "./EncodingSkills.sol";
// import "./EncodingIDs.sol";
// import "./SortIdxsAnySize.sol";
/**
 * @title Creation of all game assets via creation of timezones, countries and divisions
 * @dev Timezones range from 1 to 24, with timeZone = 0 being null.
 */

contract Shop {

    struct ShopItem {
        uint256 skillsBoost;
        uint256 duration;
        uint16 stock;
        bytes32 championshipsHash;
        string uri;
    }        

    mapping(uint256 => ShopItem) private _itemIdToItem;

    function getSkillsBoost(uint256 itemId) private view returns (uint256) { return _itemIdToItem[itemId].skillsBoost; }
    function getDuration(uint256 itemId) private view returns (uint256) { return _itemIdToItem[itemId].duration; }
    function getStock(uint256 itemId) private view returns (uint16) { return _itemIdToItem[itemId].stock; }
    function getChampionshipsHash(uint256 itemId) private view returns (bytes32) { return _itemIdToItem[itemId].championshipsHash; }
    function getUri(uint256 itemId) private view returns (string memory) { return _itemIdToItem[itemId].uri; }
    
}
