import sqlite3
from db_utils import db_connect
import os

international_db = "../country_list/country-list.sql"
db_name = "goalRev.db"
countryCodesFile = 'tmp/goalRevCountryCodes'
namesFile = 'tmp/goalRevNames'
surnamesFile = 'tmp/goalRevSurnames'

class CountryDNA:
    def __init__(self, officialName, countryCodeForNames, countryCodeForSurnames, mixRatios):
        self.officialName = officialName
        self.countryCodeForNames = countryCodeForNames
        self.countryCodeForSurnames = countryCodeForSurnames
        self.mixRatios = mixRatios


def getCountryId(tz, countryIdxInTz):
    return tz * 1000000 + countryIdxInTz


def deleteIfExists(filename):
    if os.path.isfile(filename):
        os.remove(filename)

def readCodes(filename):
    names = []
    with open(filename, 'r', newline='\n') as file:
        for line in file:
            splitted = line.rstrip('\n').split(",")
            names.append([int(splitted[0]), splitted[1], splitted[2], splitted[3]])
    return names

def readNames(filename):
    names = []
    with open(filename, 'r', newline='\n') as file:
        for line in file:
            l = line.replace('\x9a','').replace('\x8a','')
            splitted = l.rstrip('\n').rstrip('\x9a').split(",")
            names.append([int(splitted[0]), splitted[1].title()])
    return names

def db_connect():
    """ create a database connection to a database that resides
        in the memory
    """
    conn = None
    try:
        conn = sqlite3.connect(':memory:')
        print(sqlite3.version)
    except Error as e:
        print(e)
    finally:
        if conn:
            conn.close()

def getNumInCountry(country_code, allNames, allSurnames):
    return len([name for name in allNames if name[0] == country_code] \
               + [name for name in allSurnames if name[0] == country_code] )

def executeScriptsFromFile(c, filename):
    # Open and read the file as a single buffer
    fd = open(filename, 'r')
    sqlFile = fd.read()
    fd.close()

    # all SQL commands (split on ';')
    # sqlCommands = sqlFile.split(';')
    sqlCommands = filter(None, sqlFile.split(';'))
    # Execute every command from the input file
    for command in sqlCommands:
        # This will skip and report errors
        # For example, if the tables do not yet exist, this will skip over
        # the DROP TABLE commands
        try:
            c.execute(command)
        except Exception as inst:
            print("Command skipped: ", inst)


allNames = readNames(namesFile)
allSurnames = readNames(surnamesFile)
allCodes = readCodes(countryCodesFile)

deleteIfExists(db_name)
# con = sqlite3.connect(':memory:') # connect to the database
con = sqlite3.connect(db_name) # connect to the database
cur = con.cursor() # instantiate a cursor obj
executeScriptsFromFile(cur, international_db)

# ---------- Country codes ----------
written_per_country = {}
countries_sql = """
CREATE TABLE surname_regions (
    name txt NOT NULL PRIMARY KEY)"""
cur.execute(countries_sql)
properNaming = {
    1051: "Chinese",
    1005: "NonHispWhite",
    1006: "NonHispBlack",
    1007: "PacificIslander",
    1008: "AmericanIndian",
    1009: "TwoOrMoreRegions",
    1010: "Hispanic",
    1100: "Spanish",
    1053: "Japanese"
}

for code in allCodes:
    if code[0] >= 1000:
        cur.execute("INSERT INTO surname_regions VALUES ('%s');" %(properNaming[code[0]]))

import sys
sys.exit()


# ---------- Names ----------
names_sql = """
CREATE TABLE names (
    name text NOT NULL,
    country_code integer REFERENCES countries(country_code),
    idx_in_country integer NON NULL,
    PRIMARY KEY (name, country_code))"""
cur.execute(names_sql)

for name in allNames:
    exists = cur.execute("SELECT COUNT(*) FROM names WHERE (name='%s' AND country_code='%i');" % (name[1], name[0]))
    row = cur.fetchone()
    if row[0] == 0: # avoid repeating
        cur.execute("INSERT INTO names VALUES ('%s', '%i', '%i');" %(name[1], name[0], written_per_country[name[0]]))
        print(name[1], name[0], written_per_country[name[0]])
        written_per_country[name[0]] += 1

# ---------- Surnames ----------
surnames_sql = """
CREATE TABLE surnames (
    surname text NOT NULL,
    country_code integer REFERENCES countries(country_code),
    idx_in_country integer NON NULL,
    PRIMARY KEY (surname, country_code))"""
cur.execute(surnames_sql)
for surname in allSurnames:
    cur.execute("INSERT INTO surnames VALUES ('%s', '%i', '%i');" %(surname[1], surname[0], written_per_country[surname[0]]))
    # print(surname[1], surname[0], written_per_country[surname[0]])
    written_per_country[surname[0]] += 1



con.commit()
con.close()

