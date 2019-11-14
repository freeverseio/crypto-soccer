from web3 import Web3

countryCodesFile = 'goalRevCountryCodes'
nameFile = 'goalRevNames'
surnameFile = 'goalRevSurnames'

PURE_PURE_RATIO = 35
PURE_FOREIGN_RATIO = 20
FOREIGN_PURE_RATIO = 20
FOREIGN_FOREIGN_RATIO = 25


class CountryDNA:
    def __init__(self, officialName, countryCode):
        self.officialName = officialName
        self.countryCode = countryCode


countryIdToCode = {
    54392874534: CountryDNA("Spain", 8),
    42364373724: CountryDNA("Italy", 5)
}

# def readCountryCodes(countryCodesFile):
#     codes = []
#     with open(countryCodesFile, 'r', newline='\n') as file:
#         for line in file:
#             codes.append(line.split(","))
#     return codes


def getCountryIdFromPlayerId(playerId):
    if playerId % 2 == 0:
        return 54392874534
    else:
        return 42364373724

def getRndFromSeed(seed, maxVal):
    return Web3.toInt(Web3.keccak(seed)) % maxVal


def getRndNameFromCountry(country.countryCode, playerId, False)



def getNameFromPlayerId(playerId):
    countryId = getCountryIdFromPlayerId(playerId)
    assert (countryId in countryIdToCode), "Player belongs to a country not yet specified by Freeverse"

    country = countryIdToCode[countryId]

    maxVal = 100
    dice = getRndFromSeed(playerId, maxVal)
    if dice < (PURE_PURE_RATIO + PURE_FOREIGN_RATIO):
        name = getRndNameFromCountry(country.countryCode, playerId, True)
    else:
        name = getRndNameFromCountry(country.countryCode, playerId, False)

    print(country.officialName, country.countryCode)


getNameFromPlayerId(2349874534)






