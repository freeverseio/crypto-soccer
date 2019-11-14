from web3 import Web3

countryCodesFile = 'goalRevCountryCodes'
namesFile = 'goalRevNames'
surnamesFile = 'goalRevSurnames'

PURE_PURE_RATIO = 35
PURE_FOREIGN_RATIO = 20
FOREIGN_PURE_RATIO = 20
FOREIGN_FOREIGN_RATIO = 25


class CountryDNA:
    def __init__(self, officialName, countryCodeForNames, countryCodeForSurnames):
        self.officialName = officialName
        self.countryCodeForNames = countryCodeForNames
        self.countryCodeForSurnames = countryCodeForSurnames


countryIdToCode = {
    54392874534: CountryDNA("Spain", 8, 1010),
    42364373724: CountryDNA("Italy", 5, 1010)
}

# def readCountryCodes(countryCodesFile):
#     codes = []
#     with open(countryCodesFile, 'r', newline='\n') as file:
#         for line in file:
#             codes.append(line.split(","))
#     return codes

def readNames(filename):
    names = []
    with open(filename, 'r', newline='\n') as file:
        for line in file:
            splitted = line.rstrip('\n').split(",")
            names.append([int(splitted[0]), splitted[1]])
    return names

def getCountryIdFromPlayerId(playerId):
    if playerId % 2 == 0:
        return 54392874534
    else:
        return 42364373724

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
    assert (countryId in countryIdToCode), "Player belongs to a country not yet specified by Freeverse"
    country = countryIdToCode[countryId]

    maxVal = 100
    dice = getRndFromSeed(playerId, maxVal, 1)
    if dice < (PURE_PURE_RATIO + PURE_FOREIGN_RATIO):
        name = getRndNameFromCountry(country.countryCodeForNames, playerId, True)
    else:
        name = getRndNameFromCountry(country.countryCodeForNames, playerId, False)

    print(name, " from ", country.officialName, country.countryCodeForNames)
    return name


def getSurnameFromPlayerId(playerId):
    countryId = getCountryIdFromPlayerId(playerId)
    assert (countryId in countryIdToCode), "Player belongs to a country not yet specified by Freeverse"
    country = countryIdToCode[countryId]

    maxVal = 100
    dice = getRndFromSeed(playerId+2, maxVal, 3)
    if dice < (PURE_PURE_RATIO + FOREIGN_PURE_RATIO):
        name = getRndSurnameFromCountry(country.countryCodeForSurnames, playerId, True)
    else:
        name = getRndSurnameFromCountry(country.countryCodeForSurnames, playerId, False)

    print(name, " from ", country.officialName, country.countryCodeForSurnames)
    return name


# TEST 1
str = ''
for i in range(10):
    str += getNameFromPlayerId(i)
if str == 'PrimianoLidianoMargaritoMansuetoLaurentinoLiOtilioVladimiroWaalkeDuilio':
    print("TEST PASSED")
else:
    print("TEST FAILED: ", str)

print("------------")

# TEST 2
str = ''
for i in range(10):
    str += getSurnameFromPlayerId(i)
if str == 'IrahetaDiuguidHiremathBaratzVisonTineoMolieriLoeraChaimowitzPozuelos':
    print("TEST PASSED")
else:
    print("TEST FAILED: ", str)




