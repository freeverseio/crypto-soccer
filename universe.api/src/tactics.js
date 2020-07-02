// INPUTS:
// - data:
//      - array where each element is a struct with fields: "shirt_number", "encoded_skills", "red_card", "injury_matches_left"
// - tacticPatch: tacticId, shirt0, ..., shirt10, substitution0Shirt, substitution0Target, substitution0Minute, ....
//      - shirtN is a value in [0, 24] for valid team players, and 25 for no-one chosen in that position
//      - substitutionShirt, as shirtN
//      - substitutionTarget is a value in [0, 10] refering to the player that will LEAVE the field

const checkTactics = (nRedCards1stHalf, data, tacticPatch) => {
    // for each shirt in shirt0...10:
    // - c
    // - get wasAligned, redCard by player Id
    nChangesAtHalfTime = 0;
    nAlignedPlayersThatCanPlay = 0;
    Object.keys(tacticPatch).forEach(function(key, index) {
        if (typeof(key) == 'string') {
            if (key.startsWith('shirt')) {
                const shirtNum = tacticPatch[key];
                const player = getPlayerByShirtNum(shirtNum, data);
                const canPlay = !player.red_card && player.injury_matches_left == 0;
                if (canPlay) {
                    nAlignedPlayersThatCanPlay++;
                    if (!wasAligned1stHalf(encodedSkills)) { nChangesAtHalfTime++; }
                }
            }
        }
    });
    if (nAlignedPlayersThatCanPlay >= 11 - nRedCards1stHalf) { throw "too many players aligned given the 1st half redcards"}
    if (nChangesAtHalfTime >= 3) { throw "too many changes at half time" }
};

function getPlayerByShirtNum(shirtNum, data) {
    for (player of data) {
        if (player.shirt_number == shirtNum) {
            return player;
        }
    }
}


module.exports = {
    checkTactics,
};

