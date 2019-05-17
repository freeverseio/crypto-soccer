import hashlib

class GameData:
    updater = None
    result  = None

class GamesController:
    crtl = None
    def __init__(self,crtl):
        self.crtl = crtl

    # leagues ------------------------------------------------

    def get_games_count(self):
        return len(self.games)

    def new_game(self,games):
        self.VALIDATION_BLOCKS = self.VALIDATION_RESTR + self.VALIDATION_OPEN
        self.CYCLE_BLOCKS = self.PLAY_BLOCKS + self.VALIDATION_BLOCKS

        self.rand = self.crtl.bc.get_last_blockhash()
        self.league_starting_block = self.crtl.bc.get_blockno()
        self.games=[]
        for _ in range(0,games):   
            self.games.append(GameData())         

    def pending_games_to_resolve(self):
        count = 0
        for l in self.games:   
            if l.updater == None:
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

    def is_accepting_updater_update(self,staker,game_id):
        if self.crtl.stakers.can_participate(staker) != None:
            return "stacker-cannot-participate:"+self.crtl.stakers.can_participate(staker)

        # false if aleady updated
        if self.games[game_id].updater != None:
            return "league-already-updated"

        start_block = self.league_starting_block + self.PLAY_BLOCKS
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

    def updated(self,sender,game_id):
        assert self.is_accepting_updater_update(sender,game_id)== None
        self.games[game_id].updater = sender
        # reveal time is the half of the validation open time
        self.crtl.stakers.init_challange(sender)

    # challanger functions --------------------------------------
    def is_accepting_challanger_challange(self,staker,game_id):
        assert self.crtl.stakers.can_participate(staker) == None
        return self.games[game_id].updater != None

 
    def challanged(self,game_id,staker):
        assert self.crtl.stakers.can_participate(staker) == None
        self.crtl.stakers.set_challange_lier(self.games[game_id].updater)
        self.games[game_id].updater = None


class SimultaneousLeagues(GamesController):
    
    VALIDATION_RESTR  = 66
    VALIDATION_OPEN   = 34
    PLAY_BLOCKS       = 100

    def __init__(self,crtl):
        GamesController.__init__(self,crtl)

    def get_simulate_league(self,game_id):
        # returns compressed_game_state = initStatesHash, finalStatesHash, scores
        return hashlib.sha256("league"+str(game_id)).hexdigest()

    def update(self,sender,game_id,result):
        self.games[game_id].result = result
        self.updated(sender,game_id)

    def get_updater_result(self,game_id):
        return self.games[game_id].result

    def challange(self,game_id,staker,result):
        simulation_ok = self.get_simulate_league(game_id) == result
        if simulation_ok and self.games[game_id].result != result:
            self.challanged(game_id,staker)
            self.games[game_id].result = None

