import csv

database_name = "matthias.csv"

vals = []
nFields = 0
IDX_GENDER = 1
with open(database_name, 'r', newline='\n') as file:
    for (l, line) in enumerate(file):
        if l == 0:
            fields = line.split(";")
            nFields = len(fields)
        else:
            thisLine = line.split(";")
            if thisLine[IDX_GENDER] == "M":
                print(l, thisLine)
                assert (len(thisLine) == nFields, "wrong num of fields")
                vals.append(thisLine)


a = 2
    # for row in reader:
    #     print(row[0], row[1])
