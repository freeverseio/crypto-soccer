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
    return total1,total2

def test2AllPlayersEqualTeams():
    skills = np.array([50 for skill in range(5)])
    age = 20
    robotsTeam = createAllPlayersEqualTeam(skills, age, roles433)
    showTeam(robotsTeam)
    total1 = 0
    total2 = 0
    nGames = 1000
    print "Two identical teams are playing lots of games..."
    for game in range(nGames):
        goals1, goals2 = playGame(robotsTeam, robotsTeam)
        # print 'Result: %s - %s' % (goals1, goals2)
        total1 += goals1
        total2 += goals2
    print "The aggregated Result should be close to equal: %s - %s, goals per game: %s - %s " % (total1, total2, total1*1.0/nGames, total2*1.0/nGames)
    return total1,total2


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
    return total1,total2


np.random.seed(0)

# Testing equal teams games:
total1, total2 = test2EqualTeams()
success = ( total1 == 1250 and total2 == 1232 )
if not success:
    print "TEST FAILED"
else:
    print "TEST PASSED"

print "\n\n"

# Testing different teams games:
total1, total2 = test2DifferentRandomTeams()
success = (total1 == 758 and total2 == 940)
if not success:
    print "TEST FAILED"
else:
    print "TEST PASSED"

print "\n\n"


# Testing all players equal team games:
total1, total2 = test2AllPlayersEqualTeams()
success = ( total1 == 1170 and total2 == 1140 )
if not success:
    print "TEST FAILED"
else:
    print "TEST PASSED"

