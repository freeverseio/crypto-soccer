import csv

database_name = "matthias.csv"

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
            print(name)
            a=2

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


def getNamesFromCountry(countryName, fields, allNames):
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
            for name in names:
                name = name.split("+")
                if len(name) == 2 and not strcmp(name[1], ""):
                    chineseNames.append(name[1])
            names = chineseNames
        if strcmp(countryName, "Korea"):
            # korean names are all compound, and here written in format:  name1+name2
            names = [name.replace("+", "-") for name in names]
        score -= 1
    return names



allNames = []
nFields = 0
with open(database_name, 'r', newline='\n') as file:
    for (l, line) in enumerate(file):
        if l == 0:
            fields = line.rstrip('\n').split(";")
            nFields = len(fields)
        else:
            thisLine = line.rstrip('\n').split(";")
            if isNotFemale(thisLine[IDX_GENDER]):
                assert len(thisLine) == nFields, "wrong num of fields"
                allNames.append(thisLine)


# Writing Country Codes:
outCountryCodesFile = open("goalRevCountryCodes", 'w')
for country in fields[2:-1]:
    countryCode = getCountryIdx(fields, country)
    altNames = country.split("/")
    for n in range(10-len(altNames)):
        altNames.append(",")
    str = "%i" % (countryCode)
    for name in altNames:
        str += ",%s" % name
    outCountryCodesFile.write("%s\n" % str)
outCountryCodesFile.close()

# Writing actual names
totalEntries = 0
outNamesFile = open("goalRevNames", 'w')
for country in fields[-6:-1]:
# for country in fields[2:-1]:
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

print("Total Entries: ", totalEntries)

