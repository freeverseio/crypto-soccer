import csv

database_name = "sources/japanSurnames"
out_codes_file = "tmp/goalRevCountryCodes"
out_names_file = "tmp/goalRevNames"
out_surnames_file = "tmp/goalRevSurnames"

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
    outCodesFile = open(out_codes_file, 'a+')
    countryCode = 1053
    str = "%i,JapanSurnames" % (countryCode)
    for n in range(8):
        str += ","
    outCodesFile.write("%s\n" % str)
    outCodesFile.close()

    outNamesFile = open(out_surnames_file, 'a+')
    totalEntries = 0
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
print("Parsing Japanese surnames...")
with open(database_name, 'r', newline='\n') as file:
    for (l, line) in enumerate(file):
        thisLine = line.rstrip('\n')
        allNames.append(thisLine)

allNames = removeDuplications(allNames)

writeNames(allNames)