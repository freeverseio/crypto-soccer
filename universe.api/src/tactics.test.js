const { checkTactics } = require('./tactics');


describe('tactics', () => {
    describe('description of tests', () => {
        const data = getDefaultData();
        const tacticPatch = getDefaultPatch();

        test('default everything OK', () => {
            nRedCards1stHalf = 0;
            expect(() => checkTactics(nRedCards1stHalf, data, tacticPatch)).not.toThrow();
        });

        test('fails when tacticId is too large', () => {
            var tacticPatchNew = {};
            Object.assign(tacticPatchNew, tacticPatch);
            tacticPatchNew.tacticId = 2;
            expect(() => checkTactics(nRedCards1stHalf, data, tacticPatchNew)).not.toThrow();
            tacticPatchNew.tacticId = 8;
            expect(() => checkTactics(nRedCards1stHalf, data, tacticPatchNew)).not.toThrow();
            tacticPatchNew.tacticId = 9;
            expect(() => checkTactics(nRedCards1stHalf, data, tacticPatchNew)).toThrow("tacticId supported only up to 8, received 9");
        });

        test('fails when shirtNum is too large', () => {
            var tacticPatchNew = {};
            Object.assign(tacticPatchNew, tacticPatch);
            tacticPatchNew.shirt1 = 25;
            expect(() => checkTactics(nRedCards1stHalf, data, tacticPatchNew)).not.toThrow();
            tacticPatchNew.shirt1 = 26;
            expect(() => checkTactics(nRedCards1stHalf, data, tacticPatchNew)).toThrow("shirtNum too large: 26");
        });

        test('fails one red card in 1st half', () => {
            expect(() => checkTactics(nRedCards1stHalf = 1, data, tacticPatch)).toThrow("too many players aligned given the 1st half redcards: 11");
            expect(() => checkTactics(nRedCards1stHalf = 5, data, tacticPatch)).toThrow("too many players aligned given the 1st half redcards: 11");
            expect(() => checkTactics(nRedCards1stHalf = 13, data, tacticPatch)).toThrow("too many players aligned given the 1st half redcards: 11");
        });

        test('changes at half time: 0, 1, 2,3 work, but > 3 will fail', () => {
            var nRedCards1stHalf = 0;
            var tacticPatchNew = {};
            Object.assign(tacticPatchNew, tacticPatch);
            expect(() => checkTactics(nRedCards1stHalf, data, tacticPatchNew)).not.toThrow();
            tacticPatchNew.shirt0 = 14;
            expect(() => checkTactics(nRedCards1stHalf, data, tacticPatchNew)).not.toThrow();
            tacticPatchNew.shirt3 = 15;
            expect(() => checkTactics(nRedCards1stHalf, data, tacticPatchNew)).not.toThrow();
            tacticPatchNew.shirt5 = 16;
            expect(() => checkTactics(nRedCards1stHalf, data, tacticPatchNew)).not.toThrow();
            tacticPatchNew.shirt8 = 17;
            expect(() => checkTactics(nRedCards1stHalf, data, tacticPatchNew)).toThrow("too many changes at half time");
        });
    });
});

function getDefaultData() {
    data = [];
    encodedSkillsAlignedPlayer = '5986310706507378352962293074805895248510699696029696'; 
    encodedSkillsNotAlignedPlayer = '15324956156947726902719058204642840311988711972191687672616'; 
    for (p = 0; p < 11; p++) {
        data.push({ 
            "encoded_skills": encodedSkillsAlignedPlayer,
            "shirt_number": p,
            "red_card": false,
            "injury_matches_left": 0,
            "timezone_idx": 4,
            "country_idx": 0,
            "league_idx": 0,
            "match_day_idx": 1,
        });
    }
    for (p = 11; p < 18; p++) {
        data.push({ 
            "encoded_skills": encodedSkillsNotAlignedPlayer,
            "shirt_number": p,
            "red_card": false,
            "injury_matches_left": 0,
            "timezone_idx": 4,
            "country_idx": 0,
            "league_idx": 0,
            "match_day_idx": 1,
        });
    }
    return data;
}

function getDefaultPatch() {
    return {
        tacticId: 8,
        shirt0: 0,
        shirt1: 1,
        shirt2: 2,
        shirt3: 3,
        shirt4: 4,
        shirt5: 5,
        shirt6: 6,
        shirt7: 7,
        shirt8: 8,
        shirt9: 9,
        shirt10: 10,
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