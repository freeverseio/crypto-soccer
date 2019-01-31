import hashlib

class League:
    updater = None
    result  = None

# TODO change name to `SimultaneousLeagues`
class Game:
    VALIDATION_RESTR  = 66
    VALIDATION_OPEN   = 34
    LEAGUE_BLOCKS     = 100

    crtl = None
    def __init__(self,crtl):
        self.crtl = crtl

    # leagues ------------------------------------------------

    def get_league_count(self):
        return len(self.leagues)

    # game ---------------------------------------------------

    def new_game(self,leagues):
        self.VALIDATION_BLOCKS = self.VALIDATION_RESTR + self.VALIDATION_OPEN
        self.CYCLE_BLOCKS = self.LEAGUE_BLOCKS + self.VALIDATION_BLOCKS

        self.rand = self.crtl.bc.get_last_blockhash()
        self.league_starting_block = self.crtl.bc.get_blockno()
        self.leagues=[]
        for _ in range(0,leagues):   
            self.leagues.append(League())         

    def pending_leagues_to_resolve(self):
        count = 0
        for l in self.leagues:   
            if l.result == None:
                count = count + 1
        return count

    # updater functions --------------------------------------

    def bitdiff(self,x,y):
        bx = bin(int(x,16))[2:].zfill(256)
        by = bin(int(y,16))[2:].zfill(256)
        diff  = 16
        for i in range(16):
            if bx[i]==by[i]:
                diff = diff - 1
            else:
                return diff
        return diff
    # TODO change `sender` for `stacker`
    def is_accepting_updater_update(self,staker,league_id):
        if self.crtl.stakers.can_participate(staker) != None:
            return "stacker-cannot-participate:"+self.crtl.stakers.can_participate(staker)

        # false if aleady updated
        if self.leagues[league_id].result != None:
            return "league-already-updated"

        start_block = self.league_starting_block + self.LEAGUE_BLOCKS
        end_block = start_block + self.VALIDATION_BLOCKS

        # check if it's in the current update window
        if self.crtl.bc.get_blockno() < start_block or self.crtl.bc.get_blockno() > end_block:
            return "not-in-current-update-window("+str(start_block)+"..."+str(end_block)+")"

        everybody_blockno = start_block + self.crtl.game.VALIDATION_RESTR

        if self.crtl.bc.get_blockno() < start_block and self.crtl.bc.get_blockno()>everybody_blockno:
            return None

        # 0 at start_block, 16 at start_block+VALIDATION_RESTR
        expecteddiff = (16 * ( self.crtl.bc.get_blockno()  - start_block )) / self.crtl.game.VALIDATION_RESTR
        stakerhash = hashlib.sha256(staker).hexdigest()
        currentdiff = self.bitdiff(stakerhash,self.rand)
        if currentdiff > expecteddiff:
            return "not-in-incremental-window(current="+str(currentdiff)+",requiered="+str(expecteddiff)+")"

        return None

    def get_simulate_league(self,league_id):
        return hashlib.sha256("league"+str(league_id)).hexdigest()

    def update(self,sender,league_id,result):
        assert self.is_accepting_updater_update(sender,league_id)== None

        self.leagues[league_id].updater = sender
        self.leagues[league_id].result = result
        # reveal time is the half of the validation open time
        self.crtl.stakers.init_challange(sender)

    # challanger functions --------------------------------------

    # TODO add `staker`

    def is_accepting_challanger_challange(self,staker,league_id):
        assert self.crtl.stakers.can_participate(staker) == None
        return self.leagues[league_id].result != None

    def get_updater_result(self,staker,league_id):
        assert self.crtl.stakers.can_participate(staker) == None
        return self.leagues[league_id].result
 
    def challange(self,league_id,staker,result):
        assert self.crtl.stakers.can_participate(staker) == None
        simulation_ok = self.get_simulate_league(league_id) == result
        if simulation_ok and self.leagues[league_id].result != result:
            self.crtl.stakers.set_challange_lier(self.leagues[league_id].updater)
            self.leagues[league_id].updater = None
            self.leagues[league_id].result = None