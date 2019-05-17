from blockchain import Blockchain
from stakers import Stakers

class Controller:
    bc = None
    game = None
    stakers = None

    def __init__(self, gameController):
        self.bc = Blockchain()
        self.stakers = Stakers(self)
        self.game = gameController(self)

