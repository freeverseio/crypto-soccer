pragma solidity ^0.5.0;

contract Stakers {

    // constants ---------------------------------------------------------------

    string constant ERR_BADSTATE  = "err-state";
    string constant ERR_BADHASH   = "err-hash";
    string constant ERR_BADHFIN   = "err-hashfin";
    string constant ERR_BADSTAKE  = "err-stake";
    string constant ERR_POSTCOND  = "err-postcon";
    string constant ERR_BADSENDER = "err-sender";

    uint public constant REQUIRED_STAKE   = 4 ether;
    uint public constant MINENROLL_SECS   = 1 days;
    uint public constant MAXIDLE_SECS     = 1 weeks;
    uint public constant MINUNENROLL_SECS = 2 days;
    uint public constant MAXCHALL_SECS    = 5 minutes;

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
        require(st==State.UNENROLLED,ERR_BADSTATE);

        require(msg.value == REQUIRED_STAKE,ERR_BADSTAKE);
        stakers[msg.sender] = Staker({
            state : State.ENROLLING,
            touch : uint64(block.timestamp),
            onion : _onionhash
        });
    }

    function slash(address _staker) external {
        State st = state(_staker,block.timestamp);
        require(st==State.SLASHABLE,ERR_BADSTATE);

        address(0x0).transfer(REQUIRED_STAKE);
        stakers[_staker].state = State.SLASHED;

        require(state(_staker,block.timestamp)==State.SLASHED,ERR_POSTCOND);
    }

    function touch() external {
        State st = state(msg.sender,block.timestamp);
        require(st==State.ENROLLED||st==State.SLASHABLE,ERR_BADSTATE);

        stakers[msg.sender].state = State.ENROLLED;
        stakers[msg.sender].touch = uint64(block.timestamp);
        require(state(msg.sender,block.timestamp)==State.ENROLLED,ERR_POSTCOND);
    }

    function queryUnenroll() external {
        State st = state(msg.sender,block.timestamp);
        require(st==State.ENROLLED || st==State.ENROLLING,ERR_BADSTATE);

        stakers[msg.sender].state = State.UNENROLLING;
        stakers[msg.sender].touch = uint64(block.timestamp);
        require(state(msg.sender,block.timestamp)==State.UNENROLLING,ERR_POSTCOND);
    }

    function unenroll() external {
        State st = state(msg.sender,block.timestamp);
        require(st==State.UNENROLLABLE,ERR_BADSTATE);

        address(msg.sender).transfer(REQUIRED_STAKE);
        delete stakers[msg.sender];
        require(state(msg.sender,block.timestamp)==State.UNENROLLED,ERR_POSTCOND);
    }

    function resolveChallenge(bytes32 _hash) external {

        State st = state(msg.sender,block.timestamp);
        require(st==State.CHALLENGE_TT_RES||st==State.CHALLENGE_LI_RES,ERR_BADSTATE);

        // CHALLENGE_TT_RES no need to prove anything, even if lying was mandatory.
        //   If no challengers, the staker will not be slashed.
        if (st==State.CHALLENGE_LI_RES) {
            require(uint8(_hash[31]) & 0xf == 0,ERR_BADHFIN);
        }

        require(keccak256(abi.encodePacked(_hash))==stakers[msg.sender].onion,ERR_BADHASH);
        stakers[msg.sender].onion = _hash;
        stakers[msg.sender].state = State.ENROLLED;
        stakers[msg.sender].touch = uint64(block.timestamp);
        require(state(msg.sender,block.timestamp)==State.ENROLLED,ERR_POSTCOND);
    }

    // game functions  ---------------------------------------------------

    function initChallenge(address _staker) external {
        require(msg.sender == game,ERR_BADSENDER);

        State st = state(_staker,block.timestamp);
        require(st==State.ENROLLED,ERR_BADSTATE);

        stakers[_staker].state = State.CHALLENGE_TT;
        stakers[_staker].touch = uint64(block.timestamp);
        require(state(_staker,block.timestamp)==State.CHALLENGE_TT,ERR_POSTCOND);
    }

    function lierChallenge(address _staker) external {
        require(msg.sender == game,ERR_BADSENDER);

        State st = state(_staker,block.timestamp);
        require(st==State.CHALLENGE_TT,ERR_BADSTATE);

        stakers[_staker].state = State.CHALLENGE_LI_RES;
        stakers[_staker].touch = uint64(block.timestamp);
        require(state(_staker,block.timestamp)==State.CHALLENGE_LI_RES,ERR_POSTCOND);
    }
}

/*
State graphviz diagram (some typos)

digraph G {
    UNENROLLED
    ENROLLING
    ENROLLED
    UNENROLLING
    SLASHABLE
    CHALLENGE_TT [label="CHALLENGE\nTRUTHTELLER"]
    CHALLENGE_TT_RES [label="CHALLENGE\nTRUTHTELLER\nRESOLVE"]
    CHALLENGE_LI_RES [label="CHALLENGE\nLIER\nRESOLVE"]
    ENROLLED -> ENROLLED [label="touch()"]
    ENROLLED -> UNENROLLING [label="query\nunenroll()"]
    UNENROLLED -> ENROLLING [label="enroll()"]
    ENROLLING -> ENROLLED [label="MINENROLL\nBLOCKS"]
    ENROLLED -> SLASHABLE [label="MAXIDLE\nBLOCKS"]
    ENROLLED -> CHALLENGE_TT [label="updated()"]
    CHALLENGE_TT_RES -> SLASHABLE [label="MAXIDLE\nBLOCKS"]
    CHALLENGE_LI_RES -> SLASHABLE [label="MAXIDLE\nBLOCKS"]
    CHALLENGE_TT -> CHALLENGE_LI_RES [label="lierChallenge()"]
    CHALLENGE_TT -> CHALLENGE_TT_RES [label="MAXCHALL\nBLOCKS"]
    CHALLENGE_TT_RES -> ENROLLED [label="resolve\nchallenge()"]
    CHALLENGE_LI_RES -> ENROLLED [label="resolve\nchallenge()"]
    ENROLLING -> UNENROLLING [label="query\nunenroll()"]
    UNENROLLING -> UNENROLLABLE [label="MINUNENROLL\nBLOCKS"]
    UNENROLLABLE -> UNENROLLED [label="unenroll()"]
    SLASHABLE -> SLASHED [label="slash()"]
    SLASHABLE -> ENROLLED [label="touch()"]

}
*/
