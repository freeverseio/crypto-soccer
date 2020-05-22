pragma solidity >= 0.6.3;

// import "./Proxy.sol";
// import "./Assets.sol";
// import "./Market.sol";

/**
 * @title You need to have the right permissions to operate these functions
 */
 
contract Dashboard {

    address constant private NULL_ADDR = address(0x0);
    uint256 constant private FWD_GAS_LIMIT = 10000; 

    // address constant private PROXY = 0x345cA3e014Aaf5dcA488057592ee47305D9B3e10;
    // address constant private MARKET_CRYPTO = 0x345cA3e014Aaf5dcA488057592ee47305D9B3e10;
    // address constant private STAKERS = 0x345cA3e014Aaf5dcA488057592ee47305D9B3e10;
    
    address public _proxy;
    
    function setProxy(address addr) public {
        _proxy = addr;    
    }
        
    // function setSuperUser(address addr) public returns(bool) {
    //     delegate(_proxy, msg.data);
    //     return true;
    // }
   
    // function changeAllPermissions (
    //     address newCOO,
    //     address newMarket,
    //     address newMarketCrypto,
    //     address newRelay,
    // )
    //     public
    //     returns (bool) {


    //     return true;            
    // }

    /**
    * @dev execute a delegate call via fallback function
    */
    fallback () external {
        delegate(_proxy, msg.data);
    } 
    
    /**
    * @dev Delegates call. It returns the entire context execution
    * @dev NOTE: does not check if the implementation (code) address is a contract,
    *      so having an incorrect implementation could lead to unexpected results
    * @param _target Target address to perform the delegatecall
    * @param _calldata Calldata for the delegatecall
    */
    function delegate(address _target, bytes memory _calldata) internal {
        uint256 fwdGasLimit = FWD_GAS_LIMIT;
        assembly {
            let result := delegatecall(sub(gas(), fwdGasLimit), _target, add(_calldata, 0x20), mload(_calldata), 0, 0)
            let size := returndatasize()
            let ptr := mload(0x40)
            returndatacopy(ptr, 0, size)

            // revert instead of invalid() bc if the underlying call failed with invalid() it already wasted gas.
            // if the call returned error data, forward it
            switch result case 0 { revert(ptr, size) }
            default { return(ptr, size) }
        }
    }

}
