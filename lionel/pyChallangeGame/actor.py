import hashlib
import random

from stakers import StackerState

class Actor:
    crtl = None
    def __init__(self, crtl, name, address , seed, layers):
        self.crtl = crtl

        if address != "" :
            self.address = address
        else:
            self.address = hashlib.sha256(name).hexdigest()[:40]

        if seed == "" :
            seed = str(random.randrange(1000000000))

        if layers == 0 :
            layers = 500

        self.name = name
        self.hashonion = []
        r = seed
        for _ in range(0, layers):
            self.hashonion.append(r)
            r = hashlib.sha256(r).hexdigest()

        self.crtl.stakers.enroll(self.address,self.hashonion.pop())

    def query_unenroll(self):
        self.crtl.stakers.query_unenroll(self.address)

    def process(self):
        state = self.crtl.stakers.get(self.address).state

        if state == StackerState.UNENROLLING:
            if self.crtl.stakers.can_unenroll(self.address):
                self.crtl.stakers.unenroll(self.address)
                return "unenrolled"
            return "na-waiting-for-unenrolling"
                
        # { inv : self.state == StackerState.ENROLLED }

        # proof that I am alive and that I can burn some gas
        if self.crtl.stakers.is_staker_idle(self.address):
            self.crtl.stakers.touch(self.address)
            return "touched"

        return None


class Updater(Actor):

    def __init__(self, crtl, name, address,seed,layers):
        Actor.__init__(self,crtl, name, address,seed,layers)

    def process(self):
        action = Actor.process(self)
        if action != None:
            return action

        # reveal hashonion
        if self.crtl.stakers.can_resolve_challange(self.address):
            self.crtl.stakers.resolve_challange(self.address,self.hashonion.pop())
            return "onionhash-revealed"

        # update
        leagueIds = range(0,self.crtl.game.get_league_count())
        random.shuffle(leagueIds)
        for leagueId in leagueIds:
            ok =  self.crtl.game.is_accepting_updater_update(self.address,leagueId)
            if ok == None:
                lie = self.hashonion[-1].endswith("0")
                if not lie:
                    result = self.crtl.game.get_simulate_league(leagueId)
                else:
                    result = hashlib.sha256("").hexdigest()

                self.crtl.game.update(self.address,leagueId,result)
                return "updated(league="+str(leagueId)+",lying="+str(lie)+")"

        return None

class Challanger(Actor):
    def __init__(self, crtl, name, address,seed,layers):
        Actor.__init__(self,crtl, name, address,seed,layers)

    def process(self):
        action = Actor.process(self)
        if action != None:
            return action

        leagueIds = range(0,self.crtl.game.get_league_count())
        random.shuffle(leagueIds)
        for leagueId in leagueIds:
            if self.crtl.game.is_accepting_challanger_challange(self.address,leagueId):
                expected = self.crtl.game.get_simulate_league(leagueId)
                current = self.crtl.game.get_updater_result(self.address,leagueId)
                if expected != current:
                    self.crtl.game.challange(leagueId,self.address,expected)
                    return "challanged(league="+str(leagueId)+")"

        return None