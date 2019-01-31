from enum import Enum
import hashlib

class StackerState(Enum):
    ENROLLING   = 1
    ENROLLED    = 2
    UNENROLLING = 3

class StakingActor:
    crtl = None
    def __init__(self,crtl,address, stake, challange_hash):
        self.crtl = crtl

        # 1 slot
        self.address = address
        self.stake = stake
        # 1 slot = 1 + 32 + 1  + 1 + 32 + 4
        self.state = StackerState.ENROLLING 
        self.keepalive = self.crtl.bc.get_blockno()
        self.challange_lier = False
        self.challange_pending = False
        self.challange_revealblock = 0
        self.reputation = 0 # max 16

        # 1 slot
        self.challange_hash = challange_hash        

# Stackers manages the stakers (actors who stakes)
class Stakers:
    MINENROLL_BLOCKS   = 6170  # 24h
    MAXIDLE_BLOCKS     = 43200 # 1w
    MINUNENROLL_BLOCKS = 12340 # 48h  
    MINREVEAL_BLOCKS   = 20    # 5m

    stkrs = {}
    crtl = None
    def __init__(self,crtl):
        self.crtl = crtl

    # enroll a new stacker
    def enroll(self,staker,onionhash):
        self.stkrs[staker]=StakingActor(self.crtl,staker,1000,onionhash)

    # get a stacker info
    def get(self, staker):
        return self.stkrs[staker]

    # check if the stacker starts to be idle (so, danger to be slashed)
    def is_staker_idle(self,staker):
        return self.crtl.bc.get_blockno() -  self.stkrs[staker].keepalive > self.MAXIDLE_BLOCKS / 2

    # check if a staker can be slashed
    def can_idle_staker_be_slashed(self,staker):
        return self.crtl.bc.get_blockno() -  self.stkrs[staker].keepalive > self.MAXIDLE_BLOCKS

    # slash an staker
    def slash_idle_staker(self,staker):
        assert self.can_idle_staker_be_slashed(staker)
        # transfer founds to 0x0
        del self.stkrs[staker]

    # prove that an staker is alive
    def touch(self,staker):
        if self.stkrs[staker].state == StackerState.ENROLLING:
            self.stkrs[staker].state = StackerState.ENROLLED  

        assert self.stkrs[staker].state == StackerState.ENROLLED 
        self.stkrs[staker].keepalive=self.crtl.bc.get_blockno()
        self.crtl.game.get_simulate_league(0) # consume gas

    # check if exists
    def exists(self,staker):
        return staker in self.stkrs

    # query to unenroll an stacker
    def query_unenroll(self,staker):
        assert self.can_participate(staker)==None or self.stkrs[staker].state == StackerState.ENROLLING
        self.stkrs[staker].state = StackerState.UNENROLLING
        self.stkrs[staker].keepalive=self.crtl.bc.get_blockno()
    
    # check if can be unenrolled
    def can_unenroll(self,staker):
        if self.stkrs[staker].state != StackerState.UNENROLLING:
            return False

        if self.stkrs[staker].challange_pending:
            return False

        return self.crtl.bc.get_blockno() - self.stkrs[staker].keepalive > self.MINUNENROLL_BLOCKS

    # unenroll the staker
    def unenroll(self,staker):
        assert self.can_unenroll(staker), "!unenroll"

        # TODO return founds
        del self.stkrs[staker]

    # check if stacker can participate in game
    def can_participate(self,staker):
        if not staker in self.stkrs:
            return "stacker-not-registered"
        if self.stkrs[staker].state == StackerState.ENROLLING and self.crtl.bc.get_blockno() - self.stkrs[staker].keepalive >= self.MINENROLL_BLOCKS:
            return None
        if self.stkrs[staker].state != StackerState.ENROLLED:
            return "stacker-not-enrolled(state="+str(self.stkrs[staker].state)+")"
        if self.stkrs[staker].challange_pending:
            return "stacker-has-pending-challange"
        return None
     
    # only can be called by the game, creates a new challange
    #  for an stacker
    def init_challange(self,staker):
        assert self.can_participate(staker) == None
        if self.stkrs[staker].state == StackerState.ENROLLING:
            self.stkrs[staker].state = StackerState.ENROLLED  

        self.stkrs[staker].challange_lier=False
        self.stkrs[staker].challange_pending=True
        self.stkrs[staker].keepalive=self.crtl.bc.get_blockno()
        self.stkrs[staker].challange_revealblock=self.crtl.bc.get_blockno()+self.MINREVEAL_BLOCKS

    #  only can be called by the game, the challanger proved
    #   that updater lied
    def set_challange_lier(self,staker):
        assert self.stkrs[staker].challange_pending
        assert not self.stkrs[staker].challange_lier
        self.stkrs[staker].challange_lier=True

    #  check if staker is ready to reveal the preimage
    def can_resolve_challange(self,staker):
        return self.stkrs[staker].challange_pending and self.crtl.bc.get_blockno() > self.stkrs[staker].challange_revealblock

    #  reveal the preimage
    def resolve_challange(self,staker,preimage):
        assert self.can_resolve_challange(staker)
        expected = hashlib.sha256(preimage).hexdigest()
        assert expected == self.stkrs[staker].challange_hash, expected + "<-hash current->"+ self.stkrs[staker].challange_hash
        lier = self.stkrs[staker].challange_lier
        assert lier == preimage.endswith("0")
        self.stkrs[staker].challange_hash = preimage
        self.stkrs[staker].challange_pending = False
        self.stkrs[staker].reputation += 1
