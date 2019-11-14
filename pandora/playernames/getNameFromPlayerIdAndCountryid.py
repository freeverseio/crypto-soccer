countryCodesFile = 'goalRevCountryCodes'
nameFile = 'goalRevNames'
surnameFile = 'goalRevSurnames'

class CountryDNA:
    def __init__(self, officialName, countryCode):
        self.officialName = officialName
        self.countryCode = countryCode


countryIdToCode = {
    54392874534: CountryDNA("Spain", 8),
    42364373724: CountryDNA("Italy", 5)
}

def getCountryIdFromPlayerId(playerId):
    if playerId % 2 == 0:
        return 54392874534
    else:
        return 42364373724

def getNameFromPlayerId(playerId):
    countryId = getCountryIdFromPlayerId(playerId)
    assert (countryId in countryIdToCode), "Player belongs to a country not yet specified by Freeverse"

    country = countryIdToCode[countryId]
    print(country.officialName)

getNameFromPlayerId(2349874534)






