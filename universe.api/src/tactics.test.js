const { checkTactics } = require('./tactics');


describe('tactics', () => {
    describe('description of tests', () => {
        const data = getDefaultData();
        const tacticPatch = getDefaultPatch();

        test('default everything OK', () => {
            nRedCards1stHalf = 0;
            expect(() => checkTactics(nRedCards1stHalf, data, tacticPatch)).not.toThrow();
        });

        test('fails one red card in 1st half', () => {
            nRedCards1stHalf = 1;
            expect(() => checkTactics(nRedCards1stHalf, data, tacticPatch)).toThrow("too many players aligned given the 1st half redcards");
        });

    });

});

function getDefaultData() {
    data = [];
    encodedSkillsAlignedPlayer = '5986310706507378352962293074805895248510699696029696'; 
    encodedSkillsNotAlignedPlayer = '15324956156947726902719058204642840311988711972191687672616'; 
    for (p = 0; p < 18; p++) {
        data.push({ 
            "encoded_skills": encodedSkillsAlignedPlayer,
            "shirt_number": p,
            "red_card": false,
            "injury_matches_left": 0,
            "timezone_idx": 4,
            "country_idx": 0,
            "league_idx": 0,
            "match_day_idx": 1,
        })
    }
    return data;
}

function getDefaultPatch() {
    return {
        tacticId: 10,
        shirt0: 0,
        shirt1: 3,
        shirt2: 4,
        shirt3: 5,
        shirt4: 6,
        shirt5: 7,
        shirt6: 8,
        shirt7: 9,
        shirt8: 10,
        shirt9: 11,
        shirt10: 12,
        substitution0Shirt: 25,
        substitution0Target: 11,
        substitution0Minute: 0,
        substitution1Shirt: 25,
        substitution1Target: 10,
        substitution1Minute: 0,
        substitution2Shirt: 25,
        substitution2Target: 10,
        substitution2Minute: 0,
        extraAttack1: false,
        extraAttack2: false,
        extraAttack3: false,
        extraAttack4: false,
        extraAttack5: false,
        extraAttack6: false,
        extraAttack7: false,
        extraAttack8: false,
        extraAttack9: false,
        extraAttack10: false 
    };
}