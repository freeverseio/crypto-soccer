pragma solidity ^0.5.0;

contract Stakers {

    // constants ---------------------------------------------------------------

    string constant ERR_BADSTATE  = "err-state";
    string constant ERR_BADHASH   = "err-hash";
    string constant ERR_BADHFIN   = "err-hashfin";
    string constant ERR_BADSTAKE  = "err-stake";
    string constant ERR_POSTCOND  = "err-postcon";
    string constant ERR_BADSENDER = "err-sender";

    uint public constant REQUIRED_STAKE   = 0.0001 ether;
    uint public constant MINENROLL_SECS   = 20  seconds;
    uint public constant MAXIDLE_SECS     = 60 minutes;
    uint public constant MINUNENROLL_SECS = 20  seconds;
    uint public constant MAXCHALL_SECS    = 20  seconds;



    // types ---------------------------------------------------------------

    enum State {
        UNENROLLED,
        ENROLLING,
        UNENROLLING,
        UNENROLLABLE,
        ENROLLED,
        CHALLENGE_TT,
        CHALLENGE_LI_RES,
        CHALLENGE_TT_RES,
        SLASHABLE,
        SLASHED
    }

    struct Staker {
        State   state;
        uint64  touch; // time
        bytes32 onion;  // hash to reveal
    }

    // contract state  ---------------------------------------------------------------

    mapping(address=>Staker) public stakers;
    address                  public game;

    // constructor and state automata -----------------------------------------------

    function reset(address _stacker) external {
        delete stakers[_stacker];
    }
    
    string public lastError;
    event Failed(string msg);
    function _require(bool _cond, string memory _msg) internal returns (bool) {
        require(_cond,_msg);
        if (!_cond) {
            lastError = _msg;
            emit Failed(_msg);
        }
        return _cond;
    }
    
    // constructor and state automata -----------------------------------------------

    constructor(address _game) public {
        game = _game;
    }

    function state(address _staker, uint256 _when) public view returns (State) {

        if (_when == 0) {
            _when = block.timestamp;
        }

        if (stakers[_staker].state==State.UNENROLLED) return State.UNENROLLED;
        if (stakers[_staker].state==State.SLASHED) return State.SLASHED;
        if (stakers[_staker].state==State.UNENROLLABLE) return State.UNENROLLABLE;

        uint256 touchdiff = (_when - stakers[_staker].touch);
        if (stakers[_staker].state==State.ENROLLING) {
            if (touchdiff >= MINENROLL_SECS + MAXIDLE_SECS) return State.SLASHABLE;
            if (touchdiff >= MINENROLL_SECS) return State.ENROLLED;
            return State.ENROLLING;
        }
        if (stakers[_staker].state==State.UNENROLLING) {
            if (touchdiff >= MINUNENROLL_SECS) return State.UNENROLLABLE;
            return State.UNENROLLING;
        }
        if (stakers[_staker].state==State.ENROLLED) {
            if (touchdiff >= MAXIDLE_SECS) return State.SLASHABLE;
            return State.ENROLLED;
        }
        if (stakers[_staker].state==State.CHALLENGE_TT) {
            if (touchdiff >= MAXCHALL_SECS + MAXIDLE_SECS) return State.SLASHABLE;
            if (touchdiff >= MAXCHALL_SECS) return State.CHALLENGE_TT_RES;
            return State.CHALLENGE_TT;
        }
        if (stakers[_staker].state==State.CHALLENGE_TT_RES) {
            if (touchdiff >= MAXIDLE_SECS) return State.SLASHABLE;
            return State.CHALLENGE_TT_RES;
        }
        if (stakers[_staker].state==State.CHALLENGE_LI_RES) {
            if (touchdiff >= MAXIDLE_SECS) return State.SLASHABLE;
            return State.CHALLENGE_LI_RES;
        }
    }

    // staker functions  ---------------------------------------------------

    function enroll(bytes32 _onionhash) external payable {
        State st = state(msg.sender,block.timestamp);
        
        if (!_require(st==State.UNENROLLED,ERR_BADSTATE)) return;
        if (!_require(msg.value == REQUIRED_STAKE,ERR_BADSTAKE)) return;

        stakers[msg.sender] = Staker({
            state : State.ENROLLING,
            touch : uint64(block.timestamp),
            onion : _onionhash
        });
    }

    function slash(address _staker) external {
        State st = state(_staker,block.timestamp);
        if (!_require(st==State.SLASHABLE,ERR_BADSTATE)) return;

        address(0x0).transfer(REQUIRED_STAKE);
        stakers[_staker].state = State.SLASHED;

        if (!_require(state(_staker,block.timestamp)==State.SLASHED,ERR_POSTCOND)) return;
    }

    function touch() external {
        State st = state(msg.sender,block.timestamp);
        if (!_require(st==State.ENROLLED,ERR_BADSTATE)) return;

        stakers[msg.sender].state = State.ENROLLED;
        stakers[msg.sender].touch = uint64(block.timestamp);
        if (!_require(state(msg.sender,block.timestamp)==State.ENROLLED,ERR_POSTCOND)) return;
    }

    function queryUnenroll() external {
        State st = state(msg.sender,block.timestamp);
        if (!_require(st==State.ENROLLED || st==State.ENROLLING,ERR_BADSTATE)) return;

        stakers[msg.sender].state = State.UNENROLLING;
        stakers[msg.sender].touch = uint64(block.timestamp);
        if (!_require(state(msg.sender,block.timestamp)==State.UNENROLLING,ERR_POSTCOND)) return;
    }

    function unenroll() external {
        State st = state(msg.sender,block.timestamp);
        if (!_require(st==State.UNENROLLABLE,ERR_BADSTATE)) return;

        address(msg.sender).transfer(REQUIRED_STAKE);
        delete stakers[msg.sender];
        if (!_require(state(msg.sender,block.timestamp)==State.UNENROLLED,ERR_POSTCOND)) return;
    }

    function resolveChallenge(bytes32 _hash) external {

        State st = state(msg.sender,block.timestamp);
        if (!_require(st==State.CHALLENGE_TT_RES||st==State.CHALLENGE_LI_RES,ERR_BADSTATE)) return;

        // CHALLENGE_TT_RES no need to prove anything, even if lying was mandatory.
        //   If no challengers, the staker will not be slashed.
        if (st==State.CHALLENGE_LI_RES) {
            if (!_require(uint8(_hash[31]) & 0xf == 0,ERR_BADHFIN)) return;
        }

        if (!_require(keccak256(abi.encodePacked(_hash))==stakers[msg.sender].onion,ERR_BADHASH)) return;
        stakers[msg.sender].onion = _hash;
        stakers[msg.sender].state = State.ENROLLED;
        stakers[msg.sender].touch = uint64(block.timestamp);
        
        if (!_require(state(msg.sender,block.timestamp)==State.ENROLLED,ERR_POSTCOND)) return;
    }

    // game functions  ---------------------------------------------------

    function initChallenge(address _staker) external {
        if (!_require(msg.sender == game,ERR_BADSENDER)) return;

        State st = state(_staker,block.timestamp);
        if (!_require(st==State.ENROLLED,ERR_BADSTATE)) return;

        stakers[_staker].state = State.CHALLENGE_TT;
        stakers[_staker].touch = uint64(block.timestamp);
        if (!_require(state(_staker,block.timestamp)==State.CHALLENGE_TT,ERR_POSTCOND)) return;
    }

    function lierChallenge(address _staker) external {
        if (!_require(msg.sender == game,ERR_BADSENDER)) return;

        State st = state(_staker,block.timestamp);
        require(st==State.CHALLENGE_TT,ERR_BADSTATE);

        stakers[_staker].state = State.CHALLENGE_LI_RES;
        stakers[_staker].touch = uint64(block.timestamp);
        if (!_require(state(_staker,block.timestamp)==State.CHALLENGE_LI_RES,ERR_POSTCOND)) return;
    }
}


contract SoccerSim {
    
    uint256 constant UPDATABLE = 1;
    uint256 constant UPDATED = 2;
    uint256 constant CHALLANGABLE = 3;
    uint256 constant CHALLANGED = 4;

    uint256 constant LEAGUE_COUNT = 5;

    event LeagueChallangeAvailable(uint256 leagueNo, bytes32 value);
    event LeagueChallangeSucessfull(uint256 leagueNo);
    
    mapping(uint256=>uint256) leagues;

    Stakers public stakers;

    constructor() public {
        stakers = new Stakers(address(this));
        next();
    }
    
    function next() public {
        for (uint256 i = 0;i < legueCount();i++) {
            leagues[i]=UPDATABLE;
        }
    }

    function legueCount() public pure returns(uint256) {
        return LEAGUE_COUNT;   
    }

    function canLeagueBeUpdated(uint256 _leagueNo) public view returns(bool) {
        return leagues[_leagueNo]==UPDATABLE || leagues[_leagueNo]==CHALLANGED;
    }

    function update(uint256 _leagueNo, bytes32 _value) external {
        require(canLeagueBeUpdated(_leagueNo));
        if (leagues[_leagueNo]==UPDATABLE) {
            if (_value == 0x0) {
                leagues[_leagueNo]=CHALLANGABLE;
            } else {
                leagues[_leagueNo]=UPDATED;
            }
        } else if (leagues[_leagueNo]==CHALLANGED) {
            leagues[_leagueNo]=UPDATED;
        }
        stakers.initChallenge(msg.sender);
        emit LeagueChallangeAvailable(_leagueNo,_value);
    }
    
    
    function canLeagueBeChallanged(uint256 _leagueNo) public view returns(bool) {
        return leagues[_leagueNo]==CHALLANGABLE;
    }
    
    function challange(uint256 _leagueNo, bytes32 _value) external {
        require(canLeagueBeChallanged(_leagueNo));
        leagues[_leagueNo] = CHALLANGED;

        stakers.lierChallenge(msg.sender);
        emit LeagueChallangeSucessfull(_leagueNo);
    }
    
}