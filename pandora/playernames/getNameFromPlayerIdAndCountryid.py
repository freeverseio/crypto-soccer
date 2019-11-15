from web3 import Web3
import sqlite3
from db_utils import db_connect


countryCodesFile = 'goalRevCountryCodes'
namesFile = 'goalRevNames'
surnamesFile = 'goalRevSurnames'

PURE_PURE_RATIO = 35
PURE_FOREIGN_RATIO = 20
FOREIGN_PURE_RATIO = 20
FOREIGN_FOREIGN_RATIO = 25

# PURE_PURE_RATIO = 100
# PURE_FOREIGN_RATIO = 0
# FOREIGN_PURE_RATIO = 0
# # FOREIGN_FOREIGN_RATIO = 0

DEFAULT_MIX_RATIOS = [PURE_PURE_RATIO, PURE_FOREIGN_RATIO, FOREIGN_PURE_RATIO, FOREIGN_FOREIGN_RATIO]
class CountryDNA:
    def __init__(self, officialName, countryCodeForNames, countryCodeForSurnames, mixRatios):
        self.officialName = officialName
        self.countryCodeForNames = countryCodeForNames
        self.countryCodeForSurnames = countryCodeForSurnames
        self.mixRatios = mixRatios


countryIdToCountrySpec = {
    54392874534: CountryDNA("Spain", 8, 1010, DEFAULT_MIX_RATIOS),
    42364373724: CountryDNA("Italy", 5, 1010, DEFAULT_MIX_RATIOS),
    42364373721: CountryDNA("China", 51, 51,  DEFAULT_MIX_RATIOS),
    42364373722: CountryDNA("UK", 2, 1005, DEFAULT_MIX_RATIOS)
}

def db_connect():
    """ create a database connection to a database that resides
        in the memory
    """
    conn = None
    try:
        conn = sqlite3.connect(':memory:')
        print(sqlite3.version)
    except Error as e:
        print(e)
    finally:
        if conn:
            conn.close()

def readCodes(filename):
    names = []
    with open(filename, 'r', newline='\n') as file:
        for line in file:
            splitted = line.rstrip('\n').split(",")
            names.append([int(splitted[0]), splitted[1], splitted[2], splitted[3]])
    return names

def readNames(filename):
    names = []
    with open(filename, 'r', newline='\n') as file:
        for line in file:
            splitted = line.rstrip('\n').split(",")
            names.append([int(splitted[0]), splitted[1]])
    return names

def getCountryIdFromPlayerId(playerId):
    if playerId % 4 == 0:
        return 54392874534
    elif playerId % 4 == 1:
        return 42364373724
    elif playerId % 4 == 2:
        return 42364373721
    else:
        return 42364373722

def getRndFromSeed(seed, maxVal, iterations):
    for i in range(iterations):
        seed = Web3.toInt(Web3.keccak(seed))
    return seed % maxVal


def getRndNameFromCountry(allNames, countryCodeForNames, playerId, onlyThisCountry):
    if onlyThisCountry:
        names = [name[1] for name in allNames if name[0] == countryCodeForNames]
    else:
        names = [name[1] for name in allNames if name[0] != countryCodeForNames]
    return names[getRndFromSeed(playerId, len(names), 2)]

def getRndSurnameFromCountry(allNames, countryCodeForSurnames, playerId, onlyThisCountry):
    if onlyThisCountry:
        names = [name[1] for name in allNames if name[0] == countryCodeForSurnames]
    else:
        names = [name[1] for name in allNames if name[0] != countryCodeForSurnames]
    return names[getRndFromSeed(playerId, len(names), 4)]


def getCountryNameFromPlayerId(playerId):
    countryId = getCountryIdFromPlayerId(playerId)
    assert (countryId in countryIdToCountrySpec), "Player belongs to a country not yet specified by Freeverse"
    country = countryIdToCountrySpec[countryId]
    return country.officialName



def getNameFromPlayerId(playerId, allNames):
    countryId = getCountryIdFromPlayerId(playerId)
    assert (countryId in countryIdToCountrySpec), "Player belongs to a country not yet specified by Freeverse"
    country = countryIdToCountrySpec[countryId]

    maxVal = 100
    dice = getRndFromSeed(playerId, maxVal, 1)
    if dice < (country.mixRatios[0] + country.mixRatios[1]):  #PURE_PURE_RATIO + PURE_FOREIGN_RATIO:
        name = getRndNameFromCountry(allNames, country.countryCodeForNames, playerId, True)
    else:
        name = getRndNameFromCountry(allNames, country.countryCodeForNames, playerId, False)
    return name


def getSurnameFromPlayerId(playerId, allSurnames):
    countryId = getCountryIdFromPlayerId(playerId)
    assert (countryId in countryIdToCountrySpec), "Player belongs to a country not yet specified by Freeverse"
    country = countryIdToCountrySpec[countryId]

    maxVal = 100
    dice = getRndFromSeed(playerId+2, maxVal, 3)
    if dice < (country.mixRatios[0] + country.mixRatios[2]): # PURE_PURE_RATIO + FOREIGN_PURE_RATIO:
        name = getRndSurnameFromCountry(allSurnames, country.countryCodeForSurnames, playerId, True)
    else:
        name = getRndSurnameFromCountry(allSurnames, country.countryCodeForSurnames, playerId, False)
    return name

def getFullNameFromPlayerId(playerId):
    return getNameFromPlayerId(playerId, allNames)# + " " + getSurnameFromPlayerId(playerId)


allNames = readNames(namesFile)
allSurnames = readNames(surnamesFile)
allCodes = readCodes(countryCodesFile)

# # TEST 1
# str = ''
# for i in range(10):
#     str += getNameFromPlayerId(i, allNames)
# if str == 'PrimianoLidianoSongConallLaurentinoLiGangNikWaalkeDuilio':
#     print("NAMES TEST PASSED")
# else:
#     print("NAMES TEST FAILED: ", str)
#
# print("------------")
#
# # TEST 2
# str = ''
# for i in range(10):
#     str += getSurnameFromPlayerId(i, allSurnames)
# if str == 'IrahetaDiuguidSummerHsuVisonTineoJiaJorgensenChaimowitzPozuelos':
#     print("SURNAMES TEST PASSED")
# else:
#     print("SURNAMES TEST FAILED: ", str)
#
# print("------------")

# Print examples
# if False:
#     countries = ["Spain", "Italy", "China"]
#     for i in range(50):
#         a = 4*i+0
#         print(getCountryNameFromPlayerId(a)+":", getFullNameFromPlayerId(a))




con = sqlite3.connect(':memory:') # connect to the database
cur = con.cursor() # instantiate a cursor obj
countries_sql = """
CREATE TABLE countries (
    country_id integer PRIMARY KEY,
    country_name_0 text NOT NULL,
    country_name_1 text,
    country_name_2 text)"""
cur.execute(countries_sql)
names_sql = """
CREATE TABLE names (
    name text PRIMARY KEY,
    country_id integer REFERENCES countries(country_id))"""
cur.execute(names_sql)
surnames_sql = """
CREATE TABLE surnames (
    surname text PRIMARY KEY,
    country_id integer REFERENCES countries(country_id))"""
cur.execute(surnames_sql)

for code in allCodes:
    cur.execute("INSERT INTO countries VALUES ('%i', '%s', '%s', '%s');" %(code[0],code[1], code[2], code[3]))

for name in allNames:
    cur.execute("INSERT INTO names VALUES ('%s', '%i');" %(name[1],name[0]))

for surname in allSurnames:
    cur.execute("INSERT INTO surnames VALUES ('%s', '%i');" %(surname[1],surname[0]))


