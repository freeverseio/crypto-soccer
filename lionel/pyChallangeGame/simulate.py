import hashlib
import random
import time
from controller import Controller
from actor import Updater,Challanger

def vector_diff(a,b):
    a,b = sorted(a),sorted(b)
    ai,bi = 0,0
    only_in_first,only_in_second = [],[]

    while ai < len(a) or bi<len(b):
        cmp = ai<len(a) and bi<len(b)
        if cmp and a[ai] == b[bi]:
            ai += 1
            bi += 1
        elif (cmp and a[ai] < b[bi]) or bi==len(b):
            only_in_first.append(a[ai])
            ai +=1
        elif (cmp and a[ai] > b[bi]) or ai==len(a):
            only_in_second.append(b[bi])
            bi +=1
    
    return [only_in_first,only_in_second]
    

def simulate(simulation_games, fn_actors, fn_leagues, debug):

    def update_actors(block,actors):
        # process actors diff ----------------------------
        diff = vector_diff(fn_actors(block),actors)
        for a in diff[0]:
            assert not a in actors
            assert a[0]=='u' or a[0]=='c'
            if a[0]=='u':
                actors[a]=Updater(crtl,a,"","",0)
            else:
                actors[a]=Challanger(crtl,a,"","",0)
            if debug:
                print str(block)+":"+a+" enrolled"   

        for a in diff[1]:
            assert not a in actors
            assert a[0]=='u' or a[0]=='c'
            actors[a].query_unenroll()
            if debug:
                print str(block)+":"+a+" query_unenroll"        

        return actors

    crtl = Controller()
    bc,stakers,game = crtl.bc,crtl.stakers,crtl.game 

    stakers.MINENROLL_BLOCKS   = 6170   # 24h
    stakers.MAXIDLE_BLOCKS     = 43200  # 1week
    stakers.MINUNENROLL_BLOCKS = 12340  # 24h
    stakers.MINREVEAL_BLOCKS   = 21     # 5m
    game.VALIDATION_RESTR      = 256    # 1hour
    game.VALIDATION_OPEN       = 256    # 1hour
    game.LEAGUE_BLOCKS         = 256    # 1hour

    actors = {}
    init   = True

    block = bc.get_blockno()
    actors = update_actors(block,actors)
    bc.jump(stakers.MINENROLL_BLOCKS)

    gameno = 0
    start_time = time.time()
    while True:
        block = bc.get_blockno()
        
        # i'm alive
        if time.time() - start_time > 5:
            start_time = time.time()
            if debug:
                print "->"+str(block)

        # set up game leagues
        if init or block % game.CYCLE_BLOCKS == 0:
            if not init:
                if game.pending_leagues_to_resolve() != 0:
                    return "pending_leagues_to_resolve="+str(game.pending_leagues_to_resolve())
            else:
                init = False

            gameno = gameno + 1
            if gameno > simulation_games:
                break

            game.new_game(fn_leagues(block))
            start = game.league_starting_block
            challange_open = start + game.LEAGUE_BLOCKS
            challange_all = challange_open + game.VALIDATION_RESTR
            if debug:
                print "GAME ",gameno," starts:",block," open:",challange_open," all:",challange_all,"======"

        # update actors
        actors = update_actors(block,actors)

        # simulate actors
        allactors = actors.keys()
        random.shuffle(allactors)
        count = 0
        for a in allactors:
            action_log = actors[a].process()
            if action_log != None:
                if debug:
                    print str(block)+":"+a+" "+action_log       
                count = count + 1
                # we assume 80kgas per operation
                if count > 100:
                    break
        if count > 5:
            print str(block)+": blocksize ",count

        # next block
        bc.jump(1)
    
    # check that no actor has been slashed
    for actor in actors:
        if not stakers.exists(actors[actor].address):
            return "actor "+actor+" has been slashed"
    
    if debug:
        for actor in actors:
            print "actor",actor,"reputation:",stakers.get(actors[actor].address).reputation


    return None

def param_simulation():
    updaters = ['u1','u2','u3','u4','u5','u6','u7','u8']
    challangers = ['c1','c2','c3','c4','c5']
    for n_updaters in range(1,20):
        for n_leagues in range(10,10000):
            result = simulate(1,
                lambda blockno: updaters[:n_updaters]+challangers,
                lambda blockno: n_leagues,
                False
            )
            if result != None:
                print n_leagues," leagues with ",n_updaters, " updaters fails"
                break
    
def oneshot_simulation():
    random.seed(5)
    updaters = ["u"+str(x) for x in range(1,100)]
    challangers = ["c"+str(x) for x in range(1,4)]
    print "ok", simulate(100,
        lambda blockno: updaters+challangers,
        lambda blockno: 10,
        True
    )

oneshot_simulation()

# notes
#   - mes periode de restriccio => necesiten mes updaters
#   - menys periode de restriccio => mes updates per bloc
#   - que pasa si un joc no s'actualitza en un temps concret?
#
#   - incentivar per...
#   - x% ronda premi trusted nodes   : 50% only trusted nodes + 50% all
#   - (100-x)% ronda incentiu nous trusted : 50% round-robin  + 50% trusted nodes   
#
#   ***** -> el score de cada updater serveix per incentivar (donar jugadors cada "final de mes")
#
#   a) qualsevol updater  -> guanya qui paga mes per gas (economia a escala executant moltes transaccions a la vegada)
