import hashlib

class Blockchain:

    lastblock_hash = None

    def __init__(self):
        self.block = 1

    def jump(self,blocks):
        self.block += blocks
    
    def set_last_block_hash(self,hash):
        self.lastblock_hash = hash

    def get_last_blockhash(self):
        if self.lastblock_hash != None:
            return self.lastblock_hash
        return hashlib.sha256("bh"+str(self.block)).hexdigest()

    def get_blockno(self):
        return self.block
