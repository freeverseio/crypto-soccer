import csv

database_name = "sources/spanish-names/hombres.csv"
out_codes_file = "tmp/goalRevCountryCodes"
out_names_file = "tmp/goalRevNames"
out_surnames_file = "tmp/goalRevSurnames"

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
    return removeDuplications(names)

def removeDuplications(names):
    return list(dict.fromkeys(names))

def writeNames(allNames):
    outNamesFile = open(out_names_file, 'a+')
    totalEntries = 0
    countryCode = 100
    findNonValid(allNames)
    allNames = removeDuplications(allNames)
    for name in allNames:
        name = name.title()
        str = "%i,%s\n" % (countryCode, name)
        outNamesFile.write(str)

    str = "" if len(allNames) > 100 else " - WARNING"
    print(len(allNames), str)
    totalEntries += len(allNames)

    outNamesFile.close()
    print("Total Entries: ", totalEntries)


allNames = []
nFields = 0
with open(database_name, 'r', newline='\n') as file:
    for (l, line) in enumerate(file):
        if l == 0:
            fields = line.rstrip('\n').split(",")
            nFields = len(fields)
        else:
            if l > 800:
                continue
            thisLine = line.rstrip('\n').split(",")
            assert len(thisLine) == nFields, "wrong num of fields"
            if float(thisLine[2]) > 40:
                continue
            name0 = thisLine[0].split(" ")
            allNames.append(name0[0])


writeNames(allNames)