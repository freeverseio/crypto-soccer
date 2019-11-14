import csv

database_name = "englishSurnames.csv"

# pctwhite: 	Percent Non-Hispanic White Only
# pctblack: 	Percent Non-Hispanic Black Only
# pctapi: 	    Percent Non-Hispanic Asian and Pacific Islander Only
# pctaian: 	    Percent Non-Hispanic American Indian and Alaskan Native Only
# pct2prace: 	Percent Non-Hispanic of Two or More Races
# pcthispanic:  Percent Hispanic Origin

IDX_NAME = 0
MIN_SCORE = 45

def isNotFemale(str):
    return strcmp(str, "M") or strcmp(str, "?M") or strcmp(str, "?")

def isNotEmpty(str):
    return not strcmp(str, "")

def findNonValid(names):
    for name in names:
        if strcmp(name, "") or contains(name, "?") or contains(name, "/") or contains(name, "+") or contains(name, "-"):
            print("WARNING: apparently invalid name: ", name)

def contains(str, substr):
    return substr.casefold() in str.casefold()

def strcmp(str1, str2):
    return str1.casefold() == str2.casefold()

def belongsEnough(str, minScore = MIN_SCORE):
    if strcmp(str, "") or strcmp(str, "(S)"):
        return False
    return abs(float(str)) > minScore


def getCountryIdx(fields, countryName):
    for (f, field) in enumerate(fields):
        if strcmp(field, countryName):
            return f
    assert False, "country not found"


def getNamesFromCountry(countryName, fields, allNames):
    names = []
    score = MIN_SCORE
    while len(names) < 100 and score > 1:
        if (score < MIN_SCORE):
            print("reducing min_score for: ", countryName)
        names = [line[IDX_NAME] for line in allNames if belongsEnough(line[getCountryIdx(fields, countryName)], score)]
        score -= 10
    return names



allNames = []
nFields = 0
with open(database_name, 'r', newline='\n') as file:
    for (l, line) in enumerate(file):
        if l == 0:
            fields = line.rstrip('\n').split(",")
            nFields = len(fields)
        else:
            thisLine = line.rstrip('\n').split(",")
            assert len(thisLine) == nFields, "wrong num of fields"
            allNames.append(thisLine)


# Append Country Codes:
def writeCountryCodes():
    outCountryCodesFile = open("goalRevCountryCodes", 'a+')
    readableCountryNames = ["nonHispWhite", "nonHispBlack", "nonHispAsianPacificIslander", "nonHispAmerIndian", "twoOrMore", "hispanic"]
    # we will use country codes starting from 1000 not to overlap with previous dataset
    for (c, country) in enumerate(fields[5:]):
        countryCode = getCountryIdx(fields, country)
        str = "%s" % readableCountryNames[c]
        for n in range(9):
            str += ","
        outCountryCodesFile.write("%i,%s\n" % (countryCode + 1000, str))
    outCountryCodesFile.close()

def writeSurnames():
    outNamesFile = open("goalRevSurnames", 'a+')
    totalEntries = 0
    for country in fields[5:]:
        countryCode = 1000 + getCountryIdx(fields, country)
        namesInCountry = getNamesFromCountry(country, fields, allNames)
        for name in namesInCountry:
            name = name.title()
            str = "%i,%s\n" % (countryCode, name)
            outNamesFile.write(str)

        findNonValid(namesInCountry)
        str = "" if len(namesInCountry) > 100 else " - WARNING"
        print(country, len(namesInCountry), str)
        totalEntries += len(namesInCountry)

    outNamesFile.close()
    print("Total Entries: ", totalEntries)


writeCountryCodes()
writeSurnames()