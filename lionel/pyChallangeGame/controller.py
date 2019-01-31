from blockchain import Blockchain
from stakers import Stakers
from game import Game

class Controller:
    bc = None
    game = None
    stakers = None

    def __init__(self):
        self.bc = Blockchain()
        self.stakers = Stakers(self)
        self.game = Game(self)

