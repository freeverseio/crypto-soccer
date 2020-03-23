const { makeWrapResolversPlugin } = require("graphile-utils");

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
        console.log(result)
        //  (err, res) => {
        // throw "ciao"
        //     if (err) {
        //         return;
        //     } else {
        //         if (res.rows.lengh === 0) {
        //            console.error("unexistent")
        //            return;
        //         }
        //         console.info("existent")
        //         console.log(res.rows)
        //         return;
        //     }
        // });

        return resolve();
    };
};

module.exports = makeWrapResolversPlugin({
  Mutation: {
    updateTrainingByTeamId: updateTrainingByTeamIdWrapper(),
  },
});