from simulate import *



def test2EqualTeams():
    barca = createRandomTeam(roles433)
    showTeam(barca)
    total1 = 0
    total2 = 0
    nGames = 1000
    print "Two identical teams are playing lots of games..."
    for game in range(nGames):
        goals1, goals2 = playGame(barca, barca)
        # print 'Result: %s - %s' % (goals1, goals2)
        total1 += goals1
        total2 += goals2
    print "The aggregated Result should be close to equal: %s - %s, goals per game: %s - %s " % (total1, total2, total1*1.0/nGames, total2*1.0/nGames)


def test2DifferentRandomTeams():
    barca = createRandomTeam(roles433)
    madrid = createRandomTeam(roles433)
    showTeam(barca)
    showTeam(madrid)
    total1 = 0
    total2 = 0
    nGames = 1000
    print "Two different teams are playing lots of games..."
    for game in range(nGames):
        goals1, goals2 = playGame(barca, madrid)
        # print 'Result: %s - %s' % (goals1, goals2)
        total1 += goals1
        total2 += goals2
    print "The aggregated Result should be close to equal: %s - %s, goals per game: %s - %s " % (total1, total2, total1*1.0/nGames, total2*1.0/nGames)

# np.random.seed(200)
test2EqualTeams()
print "\n\n"

test2DifferentRandomTeams()
