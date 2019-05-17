import random
import numpy as np
from copy import deepcopy as duplicate
import datetime
from os import listdir, makedirs
from os.path import isfile, join, exists
import sha3
from pickle import dumps as serialize

from constants import *
from pylio import *
from structs import *

import __builtin__ as builtin


print serialize(2)
print serialHash(2)
print serialHash(PlayerState())
