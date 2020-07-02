const { makeWrapResolversPlugin } = require("graphile-utils");
const { isTrainingGroupValid, isTrainingSpecialPlayerValid } = require('./training');

const updateTrainingByTeamIdWrapper = propName => {
    return async (resolve, source, args, context, resolveInfo) => {
        const { teamId, trainingPatch } = args.input;
        const { pgClient } = context;

        const query = {
            text: 'SELECT training_points FROM teams WHERE team_id = $1',
            values: [teamId],
        };

        const result = await pgClient.query(query);
        if (result.rowCount === 0) {
            throw "unexistent team";
        }

        const allowedTP = result.rows[0].training_points;

        isTrainingGroupValid(allowedTP, trainingPatch.attackersShoot, trainingPatch.attackersSpeed, trainingPatch.attackersPass, trainingPatch.attackersDefence, trainingPatch.attackersEndurance);
        isTrainingGroupValid(allowedTP, trainingPatch.defendersShoot, trainingPatch.defendersSpeed, trainingPatch.defendersPass, trainingPatch.defendersDefence, trainingPatch.defendersEndurance);
        isTrainingGroupValid(allowedTP, trainingPatch.goalkeepersShoot, trainingPatch.goalkeepersSpeed, trainingPatch.goalkeepersPass, trainingPatch.goalkeepersDefence, trainingPatch.goalkeepersEndurance);
        isTrainingGroupValid(allowedTP, trainingPatch.midfieldersShoot, trainingPatch.midfieldersSpeed, trainingPatch.midfieldersPass, trainingPatch.midfieldersDefence, trainingPatch.midfieldersEndurance);

        isTrainingSpecialPlayerValid(allowedTP, trainingPatch.specialPlayerShoot, trainingPatch.specialPlayerSpeed, trainingPatch.specialPlayerPass, trainingPatch.specialPlayerDefence, trainingPatch.specialPlayerEndurance);

        return resolve();
    };
};


// tacticPatch: tacticId, shirt0, ..., shirt10, substitution0Shirt, substitution0Target, substitution0Minute, ....
// - shirtN is a value in [0, 24] for valid team players, and 25 for no-one chosen in that position
// - substitutionShirt, as shirtN
// - substitutionTarget is a value in [0, 10] refering to the player that will LEAVE the field
const updateTacticByTeamIdWrapper = propName => {
    return async (resolve, source, args, context, resolveInfo) => {
        const { teamId, tacticPatch } = args.input;
        const { pgClient } = context;

        var query = {
            text: 'SELECT encoded_skills, shirt_number, red_card, timezone_idx, country_idx, league_idx, match_day_idx, match_idx FROM players JOIN matches ON (players.team_Id = matches.home_team_id OR players.team_Id = matches.visitor_team_id)  WHERE (team_id = $1 AND state = $2);',
            values: [teamId, 'end'],
        };                    
        const resultQ1 = await pgClient.query(query);
        // if (resultQ1.rowCount > 1) { throw "more than one match at half time"; }
        if (resultQ1.rowCount == 0) { resolve(); }
        r = resultQ1.rows[0];
        console.log("TONI: one match found at half time ");
        console.log("quering for: ", r.timezone_idx, r.country_idx, r.league_idx, r.match_day_idx, r.match_idx);
        query = {
            text: 'SELECT COUNT(*) FROM match_events WHERE (team_id = $1 AND type = $2 AND timezone_idx = $3 AND country_idx = $4 AND league_idx = $5 AND match_day_idx = $6 AND match_idx = $7);',
            // text: 'SELECT primary_player_id FROM match_events WHERE (team_id = $1 AND type = $2 AND timezone_idx = $3 AND country_idx = $4 AND league_idx = $5 AND match_day_idx = $6 AND match_idx = $7);',
            values: [teamId, 'red_card', r.timezone_idx, r.country_idx, r.league_idx, r.match_day_idx, r.match_idx],
        };                    
        const resultQ2 = await pgClient.query(query);
        if (resultQ2.rowCount === 0) {
            throw "unexistent matchevents";
        }
        const nRedCards1stHalf = resultQ2.rows[0].count;
        console.log("TONI: numRedCards at 1st half: ", nRedCards1stHalf);


        // // for each shirt in shirt0...10:
        // // - get playerId with (teamId, shirtN) 
        // // - get wasAligned, redCard by player Id
    
        // Object.keys(tacticPatch).forEach(function(key, index) {
        //     if (typeof(key) == 'string') {
        //         if (key.startsWith('shirt')) {
        //             const shirt = tacticPatch[key];
        //             const query = {
        //                 text: 'SELECT red_card, shirt_number FROM players WHERE team_id = $1',
        //                 values: [teamId, shirt],
        //             };                    
        //             const result = await pgClient.query(query);
        //             // the result can be null, for example, if the player does not exist, has been sold, etc.
        //             if (result.rowCount > 0) {
        //                 // result.rows[0].red_card;
        //                 // result.rows[0].yellow_card_1st_half;
        //             }                    
        //         }
        //     }
        // });
        // // I will need:
        // // - is2ndHalf
        // //  - alignedEndOf1stHalf => JS code duplicated from Solidity, from encoded_skills
        // //  - get last MATCH given teamId
        // //      SELECT red_card, shirt_number FROM matches WHERE team_id = $1
        //         player_id, red_cards, shirt_number, timezone_idx, country_idx, league_idx, match_day_idx, match_idx) 
        //         JOIN: matches, players
        //         WHERE: state == half, homeTeamID == teamId OR visitorTeamID == teamId, team_id == teamId, 

        //         // this gives me all player_ids
        //         primeray red 
        //         math_events
        //         WHERE: timezone_idx, country_idx, league_idx, match_day_idx, match_idx) 

        // //  -   is2ndHalf = (state == half)
        // //  -   was there a red card in 1st half
        // checkTactics(resultOfQueries);
        return resolve();
    };
};

module.exports = makeWrapResolversPlugin({
    Mutation: {
        updateTrainingByTeamId: updateTrainingByTeamIdWrapper(),
        updateTacticByTeamId: updateTacticByTeamIdWrapper(),
    },
});