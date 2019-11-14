from web3 import Web3

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
# FOREIGN_FOREIGN_RATIO = 0

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
    42364373721: CountryDNA("China", 51, 51,  DEFAULT_MIX_RATIOS)
}

def readNames(filename):
    names = []
    with open(filename, 'r', newline='\n') as file:
        for line in file:
            splitted = line.rstrip('\n').split(",")
            names.append([int(splitted[0]), splitted[1]])
    return names

def getCountryIdFromPlayerId(playerId):
    if playerId % 3 == 0:
        return 54392874534
    elif playerId % 3 == 1:
        return 42364373724
    else:
        return 42364373721

def getRndFromSeed(seed, maxVal, iterations):
    for i in range(iterations):
        seed = Web3.toInt(Web3.keccak(seed))
    return seed % maxVal


def getRndNameFromCountry(countryCodeForNames, playerId, onlyThisCountry):
    allNames = readNames(namesFile)
    if onlyThisCountry:
        names = [name[1] for name in allNames if name[0] == countryCodeForNames]
    else:
        names = [name[1] for name in allNames if name[0] != countryCodeForNames]
    return names[getRndFromSeed(playerId, len(names), 2)]

def getRndSurnameFromCountry(countryCodeForSurnames, playerId, onlyThisCountry):
    allNames = readNames(surnamesFile)
    if onlyThisCountry:
        names = [name[1] for name in allNames if name[0] == countryCodeForSurnames]
    else:
        names = [name[1] for name in allNames if name[0] != countryCodeForSurnames]
    return names[getRndFromSeed(playerId, len(names), 4)]


def getNameFromPlayerId(playerId):
    countryId = getCountryIdFromPlayerId(playerId)
    assert (countryId in countryIdToCountrySpec), "Player belongs to a country not yet specified by Freeverse"
    country = countryIdToCountrySpec[countryId]

    maxVal = 100
    dice = getRndFromSeed(playerId, maxVal, 1)
    if dice < (country.mixRatios[0] + country.mixRatios[1]):  #PURE_PURE_RATIO + PURE_FOREIGN_RATIO:
        name = getRndNameFromCountry(country.countryCodeForNames, playerId, True)
    else:
        name = getRndNameFromCountry(country.countryCodeForNames, playerId, False)
    return name


def getSurnameFromPlayerId(playerId):
    countryId = getCountryIdFromPlayerId(playerId)
    assert (countryId in countryIdToCountrySpec), "Player belongs to a country not yet specified by Freeverse"
    country = countryIdToCountrySpec[countryId]

    maxVal = 100
    dice = getRndFromSeed(playerId+2, maxVal, 3)
    if dice < (country.mixRatios[0] + country.mixRatios[2]): # PURE_PURE_RATIO + FOREIGN_PURE_RATIO:
        name = getRndSurnameFromCountry(country.countryCodeForSurnames, playerId, True)
    else:
        name = getRndSurnameFromCountry(country.countryCodeForSurnames, playerId, False)
    return name

def getFullNameFromPlayerId(playerId):
    return getNameFromPlayerId(playerId) + " " + getSurnameFromPlayerId(playerId)

# TEST 1
str = ''
for i in range(10):
    str += getNameFromPlayerId(i)
if str == 'PrimianoLidianoSongEfraínDorandoChuOtilioVladimiroLeonBenigno':
    print("NAMES TEST PASSED")
else:
    print("NAMES TEST FAILED: ", str)

print("------------")

# TEST 2
str = ''
for i in range(10):
    str += getSurnameFromPlayerId(i)
if str == 'IrahetaDiuguidSummerBaratzVisonLuMolieriLoeraPrimaPozuelos':
    print("SURNAMES TEST PASSED")
else:
    print("SURNAMES TEST FAILED: ", str)

print("------------")

# Print examples
if True:
    countries = ["Spain", "Italy", "China"]
    for i in range(20):
        print(countries[i % 3]+":", getFullNameFromPlayerId(i))




