const { makeWrapResolversPlugin } = require("graphile-utils");

const isTrainingGroupValid = (TP, shoot, speed, pass, defence, endurance) => {
    const sum = shoot + speed + pass + defence + endurance;

    if (sum === 0) return;
    if (sum > TP) throw "group sum " + sum + " exceeds available TP " + TP;

    if (10 * shoot <= 6 * TP) throw "shoot exceeds 60% of TP " + TP;
    if (10 * speed <= 6 * TP) throw "speed exceeds 60% of TP " + TP;
    if (10 * pass <= 6 * TP) throw "pass exceeds 60% of TP " + TP;
    if (10 * defence <= 6 * TP) throw "defence exceeds 60% of TP " + TP;
    if (10 * endurance <= 6 * TP) throw "endurance exceeds 60% of TP " + TP;
};

const isTrainingSpecialPlayerValid = (TP, shoot, speed, pass, defence, endurance) => {
    const specialPlayerTP = TP * 11 / 10;
    isTrainingGroupValid(specialPlayerTP, shoot, speed, pass, defence, endurance);
};

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

        const allowedTP = result.rows[0];

        isTrainingGroupValid(allowedTP, trainingPatch.attackersShoot, trainingPatch.attackersSpeed, trainingPatch.attackersPass, trainingPatch.attackersDefence, trainingPatch.attackersEndurance); 
        isTrainingGroupValid(allowedTP, trainingPatch.defendersShoot, trainingPatch.defendersSpeed, trainingPatch.defendersPass, trainingPatch.defendersDefence, trainingPatch.defendersEndurance); 
        isTrainingGroupValid(allowedTP, trainingPatch.goalkeepersShoot, trainingPatch.goalkeepersSpeed, trainingPatch.goalkeepersPass, trainingPatch.goalkeepersDefence, trainingPatch.goalkeepersEndurance); 
        isTrainingGroupValid(allowedTP, trainingPatch.midfieldersShoot, trainingPatch.midfieldersSpeed, trainingPatch.midfieldersPass, trainingPatch.midfieldersDefence, trainingPatch.midfieldersEndurance); 

        isTrainingSpecialPlayerValid(allowedTP, trainingPatch.specialPlayersShoot, trainingPatch.specialPlayersSpeed, trainingPatch.specialPlayersPass, trainingPatch.specialPlayersDefence, trainingPatch.specialPlayersEndurance); 

        return resolve();
    };
};

module.exports = makeWrapResolversPlugin({
  Mutation: {
    updateTrainingByTeamId: updateTrainingByTeamIdWrapper(),
  },
});