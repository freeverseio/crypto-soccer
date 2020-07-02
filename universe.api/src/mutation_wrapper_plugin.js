const { makeWrapResolversPlugin } = require("graphile-utils");
const { isTrainingGroupValid, isTrainingSpecialPlayerValid } = require('./training');
const { checkTactics } = require("./tactics");

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

const updateTacticByTeamIdWrapper = propName => {
    return async (resolve, source, args, context, resolveInfo) => {
        const { teamId, tacticPatch } = args.input;
        const { pgClient } = context;

        var query = {
            text: 'SELECT encoded_skills, shirt_number, red_card, timezone_idx, country_idx, league_idx, match_day_idx, match_idx FROM players JOIN matches ON (players.team_Id = matches.home_team_id OR players.team_Id = matches.visitor_team_id)  WHERE (team_id = $1 AND state = $2);',
            values: [teamId, 'end'],
        };                    
        const resultQ1 = await pgClient.query(query);
        console.log("TONI: nResults = ", resultQ1.rowCount);
        if (resultQ1.rowCount == 0) { resolve(); } // it is 1st half, no need to check
        data = resultQ1.rows;
        // TODO: check that all entries in data share the same timezone, country, etc.
        console.log(data);
        console.log("TONI: one match found at half time ");
        console.log("quering for: ", data[0].timezone_idx, data[0].country_idx, data[0].league_idx, data[0].match_day_idx, data[0].match_idx);

        query = {
            text: 'SELECT COUNT(*) FROM match_events WHERE (team_id = $1 AND type = $2 AND timezone_idx = $3 AND country_idx = $4 AND league_idx = $5 AND match_day_idx = $6 AND match_idx = $7);',
            values: [teamId, 'red_card', data[0].timezone_idx, data[0].country_idx, data[0].league_idx, data[0].match_day_idx, data[0].match_idx],
        };                    
        const resultQ2 = await pgClient.query(query);
        if (resultQ2.rowCount === 0) {
            throw "unexistent matchevents";
        }
        const nRedCards1stHalf = resultQ2.rows[0].count;
        console.log("TONI: numRedCards at 1st half: ", nRedCards1stHalf);

        console.log(data);
    
        checkTactics(nRedCards1stHalf, data, tacticPatch);
        return resolve();
    };
};

module.exports = makeWrapResolversPlugin({
    Mutation: {
        updateTrainingByTeamId: updateTrainingByTeamIdWrapper(),
        updateTacticByTeamId: updateTacticByTeamIdWrapper(),
    },
});