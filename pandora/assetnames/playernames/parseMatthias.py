import csv

database_name = "sources/matthias.csv"
out_codes_file = "tmp/goalRevCountryCodes"
out_names_file = "tmp/goalRevNames"
out_surnames_file = "tmp/goalRevSurnames"

IDX_NAME = 0
IDX_GENDER = 1
MIN_SCORE = 5

def isNotFemale(str):
    return strcmp(str, "M") or strcmp(str, "?M") or strcmp(str, "?")

def isNotEmpty(str):
    return not strcmp(str, "")

def findNonValid(names):
    for name in names:
        if strcmp(name, "") or contains(name, "?") or contains(name, "/") or contains(name, "+"):# or contains(name, "-"):
            print("WARNING: apparently invalid name: ", name)

def strcmp(str1, str2):
    return str1.casefold() == str2.casefold()

def contains(str, substr):
    return substr.casefold() in str.casefold()

def belongsEnough(str, minScore = MIN_SCORE):
    if strcmp(str, ""):
        return False
    return abs(int(str)) > minScore


def getCountryIdx(fields, countryName):
    for (f, field) in enumerate(fields):
        if strcmp(field, countryName):
            return f
    assert False, "country not found"


def getNamesFromCountry(countryName, fields, allNames, getSurname = False):
    names = []
    score = MIN_SCORE
    while len(names) < 100 and score > 1:
        if (score < MIN_SCORE):
            print("reducing min_score for: ", countryName)
        names = [line[IDX_NAME] for line in allNames \
                 if belongsEnough(line[getCountryIdx(fields, countryName)], score) \
                 and isNotEmpty(line[IDX_NAME])]
        if strcmp(countryName, "China"):
            # chinese names are written in format:  surname+name
            chineseNames = []
            idx = 0 if getSurname == True else 1
            for name in names:
                name = name.split("+")
                if len(name) == 2 and not strcmp(name[idx], ""):
                    chineseNames.append(name[idx])
            names = chineseNames
        if strcmp(countryName, "Korea"):
            # korean names are all compound, and here written in format:  name1+name2
            names = [name.replace("+", "-") for name in names]
        if strcmp(countryName, "Arabia/Persia"):
            # some persian names are written in format name+surname
            names = [name.split("+")[0] for name in names]
        score -= 1
    return removeDuplications(names)

def removeDuplications(names):
    return list(dict.fromkeys(names))

allNames = []
nFields = 0
with open(database_name, 'r', newline='\n') as file:
    for (l, line) in enumerate(file):
        if l == 0:
            fields = line.rstrip('\n').split(";")
            nFields = len(fields)
        else:
            thisLine = line.rstrip('\n').replace("'","").split(";")
            if isNotFemale(thisLine[IDX_GENDER]):
                assert len(thisLine) == nFields, "wrong num of fields"
                allNames.append(thisLine)


def writeCountryCodes():
    outCountryCodesFile = open(out_codes_file, 'w')
    for country in fields[2:-1]:
        countryCode = getCountryIdx(fields, country)
        altNames = country.split("/")
        str = "%i" % (countryCode)
        for name in altNames:
            str += ",%s" % name
        for n in range(10-len(altNames)):
            str += ","
        outCountryCodesFile.write("%s\n" % str)
    # add china code = 1051
    outCountryCodesFile.write("1051,ChinaSurnames,,,,,,,,,\n")
    outCountryCodesFile.close()

def writeNames():
    totalEntries = 0
    outNamesFile = open(out_names_file, 'w')
    # for country in fields[-6:-1]:
    for country in fields[2:-1]:
        countryCode = getCountryIdx(fields, country)
        namesInCountry = getNamesFromCountry(country, fields, allNames)
        for name in namesInCountry:
            str = "%i,%s\n" % (countryCode, name)
            outNamesFile.write(str)

        str = "" if len(namesInCountry) > 100 else " - WARNING"
        print(country, len(namesInCountry), str)
        totalEntries += len(namesInCountry)
        findNonValid(namesInCountry)
    outNamesFile.close()

def writeSurnames():
    totalEntries = 0
    outNamesFile = open(out_surnames_file, 'w')
    country = "China"
    countryCode = getCountryIdx(fields, country)
    namesInCountry = getNamesFromCountry(country, fields, allNames, getSurname = True)
    different = []
    chinaCodeAsSurnames = 1000 + countryCode
    for name in namesInCountry:
        if name not in different:
            str = "%i,%s\n" % (chinaCodeAsSurnames, name)
            outNamesFile.write(str)
            different.append(name)
            totalEntries += 1

    print(country, totalEntries)
    findNonValid(namesInCountry)
    outNamesFile.close()

writeCountryCodes()
writeNames()
writeSurnames()
