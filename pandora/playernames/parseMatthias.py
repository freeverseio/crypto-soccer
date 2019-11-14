import csv

database_name = "matthias.csv"

IDX_NAME = 0
IDX_GENDER = 1
MIN_SCORE = 5

def isNotFemale(str):
    return strcmp(str, "M") or strcmp(str, "?M") or strcmp(str, "?")

def strcmp(str1, str2):
    return str1.casefold() == str2.casefold()

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
        names = [line[IDX_NAME] for line in allNames if belongsEnough(line[getCountryIdx(fields, countryName)], score)]
        if strcmp(countryName, "China"):
            # chinese names are written in format:  surname+name
            chineseNames = []
            for name in names:
                name = name.split("+")
                if len(name) == 2 and not strcmp(name[1], ""):
                    chineseNames.append(name[1])
            names = chineseNames
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

# for country in fields[-6:-1]:
for country in fields[2:-1]:
    namesInCountry = getNamesFromCountry(country, fields, allNames)
    str = "" if len(namesInCountry) > 100 else " - WARNING"
    print(country, len(namesInCountry), str)

