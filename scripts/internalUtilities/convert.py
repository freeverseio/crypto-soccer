import json
import copy

def reformat(old):
    final = copy.deepcopy(old)
    final['HomeTeam']['Training'] = {}
    final['HomeTeam']['Training']['special_player_shirt'] =  old['HomeTeam']['Training']['special_player_shirt']
    final['HomeTeam']['Training']['Goalkeepers'] = {
        'Shoot': old['HomeTeam']['Training']['goalkeepers_shoot'],
        'Speed': old['HomeTeam']['Training']['goalkeepers_speed'],
        'Pass': old['HomeTeam']['Training']['goalkeepers_pass'],
        'Defence': old['HomeTeam']['Training']['goalkeepers_defence'],
        'Endurance': old['HomeTeam']['Training']['goalkeepers_endurance'],
    } 
    final['HomeTeam']['Training']['Defenders'] = {
        'Shoot': old['HomeTeam']['Training']['defenders_shoot'],
        'Speed': old['HomeTeam']['Training']['defenders_speed'],
        'Pass': old['HomeTeam']['Training']['defenders_pass'],
        'Defence': old['HomeTeam']['Training']['defenders_defence'],
        'Endurance': old['HomeTeam']['Training']['defenders_endurance'],
    } 
    final['HomeTeam']['Training']['Midfielders'] = {
        'Shoot': old['HomeTeam']['Training']['midfielders_shoot'],
        'Speed': old['HomeTeam']['Training']['midfielders_speed'],
        'Pass': old['HomeTeam']['Training']['midfielders_pass'],
        'Defence': old['HomeTeam']['Training']['midfielders_defence'],
        'Endurance': old['HomeTeam']['Training']['midfielders_endurance'],
    } 
    final['HomeTeam']['Training']['Attackers'] = {
        'Shoot': old['HomeTeam']['Training']['attackers_shoot'],
        'Speed': old['HomeTeam']['Training']['attackers_speed'],
        'Pass': old['HomeTeam']['Training']['attackers_pass'],
        'Defence': old['HomeTeam']['Training']['attackers_defence'],
        'Endurance': old['HomeTeam']['Training']['attackers_endurance'],
    } 
    final['HomeTeam']['Training']['SpecialPlayer'] = {
        'Shoot': old['HomeTeam']['Training']['special_player_shoot'],
        'Speed': old['HomeTeam']['Training']['special_player_speed'],
        'Pass': old['HomeTeam']['Training']['special_player_pass'],
        'Defence': old['HomeTeam']['Training']['special_player_defence'],
        'Endurance': old['HomeTeam']['Training']['special_player_endurance'],
    } 
    return final



filename = '3859fc1422bc9d7e58621e77466eb42c7db8cc2305687bfe41b23bc137e14d70.1st.error.json'
with open(filename, 'r') as f:
    old = json.load(f)

with open('data.json', 'w') as f:
    json.dump(reformat(old), f, sort_keys=True, indent=4)
    