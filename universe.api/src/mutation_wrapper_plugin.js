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

module.exports = makeWrapResolversPlugin({
    Mutation: {
        updateTrainingByTeamId: updateTrainingByTeamIdWrapper(),
    },
});