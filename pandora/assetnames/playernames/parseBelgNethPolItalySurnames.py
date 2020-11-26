import csv

database_name = "sources/Belgium-Netherlands_Italy_Poland/BelgianSurnames.csv"
out_codes_file = "tmp/goalRevCountryCodes"
out_names_file = "tmp/goalRevNames"
out_surnames_file = "tmp/goalRevSurnames"

def findNonValid(names):
    for name in names:
        if strcmp(name, "") or contains(name, "?") or contains(name, "/") or contains(name, "+") or contains(name, "-"):
            print("WARNING: apparently invalid name: ", name)

def contains(str, substr):
    return substr.casefold() in str.casefold()

def strcmp(str1, str2):
    return str1.casefold() == str2.casefold()

def removeDuplications(names):
    return list(dict.fromkeys(names))

def writeNames(allNames, countryCode, countryName):
    outCodesFile = open(out_codes_file, 'a+')
    str = "%i,%s" % (countryCode, countryName)
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


def writeAllNamesForCountry(database_name, countryCode, countryName):
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
                allNames.append(thisLine[0])
    writeNames(allNames, countryCode, countryName)


database_name = "sources/Belgium-Netherlands_Italy_Poland/BelgianSurnames.csv"
countryCode = 1141
countryName = "BelgiumSurnames"
writeAllNamesForCountry(database_name, countryCode, countryName)

database_name = "sources/Belgium-Netherlands_Italy_Poland/PolishSurnames.csv"
countryCode = 1125
countryName = "PolandSurnames"
writeAllNamesForCountry(database_name, countryCode, countryName)

database_name = "sources/Belgium-Netherlands_Italy_Poland/ItalianSurnames.csv"
countryCode = 1105
countryName = "ItalySurnames"
writeAllNamesForCountry(database_name, countryCode, countryName)

database_name = "sources/Belgium-Netherlands_Italy_Poland/NetherlandsSurnames.csv"
countryCode = 1112
countryName = "NetherlandsSurnames"
writeAllNamesForCountry(database_name, countryCode, countryName)
