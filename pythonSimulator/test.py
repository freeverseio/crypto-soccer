from simulate import *


barca = createRandomTeam(0,roles433)

showTeam(barca)

total1 = 0
total2 = 0
for game in range(1000):
    goals1, goals2 = playGame(barca, barca, game)
    print 'Result: %s - %s' % (goals1, goals2)
    total1 += goals1
    total2 += goals2
print 'Aggregated Result: %s - %s' % (total1, total2)
