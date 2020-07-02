// INPUTS:
// - data:
//      - array where each element is a struct with fields: "shirt_number", "encoded_skills", red_card
// - tacticPatch: tacticId, shirt0, ..., shirt10, substitution0Shirt, substitution0Target, substitution0Minute, ....
//      - shirtN is a value in [0, 24] for valid team players, and 25 for no-one chosen in that position
//      - substitutionShirt, as shirtN
//      - substitutionTarget is a value in [0, 10] refering to the player that will LEAVE the field

const checkTactics = (nRedCards1stHalf, data, tacticPatch) => {
    return
    // for each shirt in shirt0...10:
    // - get playerId with (teamId, shirtN) 
    // - get wasAligned, redCard by player Id
    Object.keys(tacticPatch).forEach(function(key, index) {
        if (typeof(key) == 'string') {
            if (key.startsWith('shirt')) {
                const shirt = tacticPatch[key];
            }
        }
    });



};


module.exports = {
    checkTactics,
};

