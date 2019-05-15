from controller import Controller
from actor import Updater,Challanger
from game import SimultaneousLeagues

import hashlib
import random

def test_staker_idle_and_slash():
    crtl = Controller(SimultaneousLeagues)
    bc,stakers = crtl.bc,crtl.stakers

    crtl.stakers.enroll("stk1","onionh")
    bc.jump(stakers.MINENROLL_BLOCKS)
    
    # set to idle state
    assert not crtl.stakers.is_staker_idle("stk1")
    bc.jump(stakers.MAXIDLE_BLOCKS/2 + 1 )
    assert stakers.is_staker_idle("stk1")
    assert not stakers.can_idle_staker_be_slashed("stk1")
    
    stakers.touch("stk1")
    assert not stakers.is_staker_idle("stk1")
    assert stakers.can_participate("stk1") == None

    # set to idle state again and slash
    bc.jump(stakers.MAXIDLE_BLOCKS + 1)
    assert stakers.is_staker_idle("stk1")
    assert stakers.can_idle_staker_be_slashed("stk1")

    stakers.slash_idle_staker("stk1")
    assert not stakers.can_participate("stk1") == None

def test_staker_enroll_unenroll():
    crtl = Controller(SimultaneousLeagues)
    bc,stakers = crtl.bc,crtl.stakers 

    stakers.enroll("stk1","onionh")
    stakers.query_unenroll("stk1")
    assert not stakers.can_participate("stk1") == None
    assert not stakers.can_unenroll("stk1")
    bc.jump(stakers.MINUNENROLL_BLOCKS + 1)
    assert stakers.can_unenroll("stk1")
    assert not stakers.can_participate("stk1") == None
    stakers.unenroll("stk1")
    assert not stakers.can_participate("stk1") == None

def test_staker_challange():
    crtl = Controller(SimultaneousLeagues)
    bc,stakers = crtl.bc,crtl.stakers 

    # generate a sequence of ok,ok,lier,ok
    """
    a0c8aafae9e169965dbefd4eb9180d13a60b2a4d62301bb660dc2130549085e6
    7d2e7b5e1e339c6289568accbdabb2f990480664b3afc777ffa4c51085b70c30
    96a90635f5143c2bbe715e4acfee430c25ba24701a0ec91a6001427de474a53c
    f999afc3e06ceb47e2fbeb82ef370d420c701fad413595040df3352fb8e0d85c
    """
    r0_ok = str("a0c8aafae9e169965dbefd4eb9180d13a60b2a4d62301bb660dc2130549085e6") 
    r1_lier = hashlib.sha256(r0_ok).hexdigest() # finishes with zero
    r2_ok = hashlib.sha256(r1_lier).hexdigest()
    r3_ok = hashlib.sha256(r2_ok).hexdigest()

    stakers.enroll("stk1",r3_ok)
    bc.jump(stakers.MINENROLL_BLOCKS+1)

    # round 1, reveal r2_ok
    stakers.init_challange("stk1")
    assert not stakers.can_participate("stk1") == None
    assert not stakers.can_resolve_challange("stk1")
    bc.jump(stakers.MINREVEAL_BLOCKS+1)

    assert stakers.can_resolve_challange("stk1")
    stakers.resolve_challange("stk1",r2_ok) 
    assert stakers.can_participate("stk1") == None
    assert stakers.get("stk1").reputation == 1

    # round 2, reveal r1_lier
    stakers.init_challange("stk1")
    stakers.set_challange_lier("stk1")
    bc.jump(stakers.MINREVEAL_BLOCKS+1)
    
    stakers.resolve_challange("stk1",r1_lier)
    assert stakers.can_participate("stk1") == None
    assert stakers.get("stk1").reputation == 2

    # round 3, reveal r0_ok
    stakers.init_challange("stk1")
    bc.jump(stakers.MINREVEAL_BLOCKS+1)
 
    stakers.resolve_challange("stk1",r0_ok)
    assert stakers.can_participate("stk1") == None

def test_updater_window():
    crtl = Controller(SimultaneousLeagues)
    bc,stakers,game = crtl.bc,crtl.stakers,crtl.game 

    u1 = Updater(
        crtl,
        "u1",
        "0000000000000000000000000000000000000000000000000000000000000000",
        "a0c8aafae9e169965dbefd4eb9180d13a60b2a4d62301bb660dc2130549085e6",
        4
    )
    bc.jump(stakers.MINENROLL_BLOCKS)

    # start game
    game.new_game(1)
    bc.jump(game.PLAY_BLOCKS + 1)

    # check incremental window
    assert game.is_accepting_updater_update(u1.address,0) == "not-in-incremental-window(current=16,requiered=0)"
    bc.jump(game.VALIDATION_RESTR / 4) 
    assert game.is_accepting_updater_update(u1.address,0) == "not-in-incremental-window(current=16,requiered=4)"
    bc.jump(game.VALIDATION_RESTR / 4) 
    assert game.is_accepting_updater_update(u1.address,0) == "not-in-incremental-window(current=16,requiered=8)"

    # jump to 2/3, should update
    bc.jump(game.VALIDATION_RESTR / 2)
    assert game.is_accepting_updater_update(u1.address,0) == None
    assert u1.process() == "updated(league=0,lying=False)"

    # now not accepting more updates
    assert game.is_accepting_updater_update(u1.address,0) == "stacker-cannot-participate:stacker-has-pending-challange"
    assert u1.process() == None
    
def test_challanging_game():
    crtl = Controller(SimultaneousLeagues)
    bc,stakers,game = crtl.bc,crtl.stakers,crtl.game 

    u1 = Updater(
        crtl,
        "u1",
        "0000000000000000000000000000000000000000000000000000000000000000",
        "a0c8aafae9e169965dbefd4eb9180d13a60b2a4d62301bb660dc2130549085e6",
        4
    )
    c1 = Challanger(crtl,"c1","00","",0)
    bc.jump(stakers.MINENROLL_BLOCKS)

    game.new_game(1)

    # start game #1, jump all_updaters window, update
    bc.jump(game.PLAY_BLOCKS + 1)
    bc.jump(game.VALIDATION_RESTR)
    assert game.pending_games_to_resolve() == 1
    assert u1.process() == "updated(league=0,lying=False)"
    assert game.pending_games_to_resolve() == 0
    bc.jump(game.CYCLE_BLOCKS - game.VALIDATION_RESTR - game.PLAY_BLOCKS)
    assert u1.process() == "onionhash-revealed"

    # start game #2
    game.new_game(1)
    bc.jump(game.PLAY_BLOCKS + 1)
    bc.jump(game.VALIDATION_RESTR)
    assert u1.process() == "updated(league=0,lying=True)"

    # challanger challanges league 0
    assert c1.process() == "challanged(league=0)"

    # u1 reveals lier hash 
    bc.jump(game.CYCLE_BLOCKS - game.VALIDATION_RESTR - game.PLAY_BLOCKS)
    assert u1.process() == "onionhash-revealed"

test_staker_enroll_unenroll()
test_staker_idle_and_slash()
test_staker_challange()
test_updater_window()
test_challanging_game()
print "done."