letters = ["A", "C", "F", "Z"]

for one in letters:
    for two in letters:
        if one != two:
            print("%s. %s." %(one, two))
