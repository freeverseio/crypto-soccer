from web3 import Web3
import sqlite3
from db_utils import db_connect
import os

db_name = "goalRevFull.db"
mainNames= 'dbs/main.txt'
prefixNames = 'dbs/firsts.txt'
suffixNames = 'dbs/lasts.txt'

def deleteIfExists(filename):
    if os.path.isfile(filename):
        os.remove(filename)

def purgeRepeated(names):
    purged = []
    for n in names:
        if (n not in purged) and (n != ""):
            purged.append(n)
    return purged

def readNames(filename):
    names = []
    with open(filename, 'r', newline='\n') as file:
        for line in file:
            splitted = line.rstrip('\n')
            names.append(splitted)
    return purgeRepeated(names)

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


def writeTable(filename, tableName, cur):
    sql = """
    CREATE TABLE %s (
        idx integer PRIMARY KEY,
        name text NOT NULL UNIQUE)""" % (tableName)
    cur.execute(sql)

    data = readNames(filename)
    for (n, name) in enumerate(data):
        print(n, name)
        cur.execute("INSERT INTO %s VALUES ('%i', '%s');" % (tableName, n, name))


# deleteIfExists(db_name)
con = sqlite3.connect(':memory:') # connect to the database
con = sqlite3.connect(db_name) # connect to the database
cur = con.cursor() # instantiate a cursor obj

cur.execute("PRAGMA foreign_keys = ON;")

writeTable(mainNames, "team_mainnames", cur)
writeTable(prefixNames, "team_prefixnames", cur)
writeTable(suffixNames, "team_suffixnames", cur)

con.commit()
con.close()
