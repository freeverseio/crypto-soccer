const resolvers = (sql, assets, from) => {
  return {
    Mutation: {
      transferFirstBotToAddr: async (_, { timezone, countryIdxInTimezone, address }) => {
        const gas = await assets.methods
          .transferFirstBotToAddr(timezone, countryIdxInTimezone, address)
          .estimateGas();
        await assets.methods
          .transferFirstBotToAddr(timezone, countryIdxInTimezone, address)
          .send({ from, gas });
        return true;
      },
      setTactic: async (_, params, context) => {
        const query = sql.query`INSERT INTO tactics (
						verse,
						team_id,
						tactic_id,
                        shirt_0,
                        shirt_1,
                        shirt_2,
                        shirt_3,
                        shirt_4,
                        shirt_5,
                        shirt_6,
                        shirt_7,
                        shirt_8,
                        shirt_9,
                        shirt_10,
                        shirt_11,
                        shirt_12,
                        shirt_13,
                        extra_attack_1,
                        extra_attack_2,
                        extra_attack_3,
                        extra_attack_4,
                        extra_attack_5,
                        extra_attack_6,
                        extra_attack_7,
                        extra_attack_8,
                        extra_attack_9,
                        extra_attack_10
		) VALUES (
                ${sql.value('9223372036854775807')},
                ${sql.value(params.teamId)},
                ${sql.value(params.tacticId)},
                ${sql.value(params.shirt0)}, 
                ${sql.value(params.shirt1)}, 
                ${sql.value(params.shirt2)}, 
                ${sql.value(params.shirt3)}, 
                ${sql.value(params.shirt4)}, 
                ${sql.value(params.shirt5)}, 
                ${sql.value(params.shirt6)}, 
                ${sql.value(params.shirt7)}, 
                ${sql.value(params.shirt8)}, 
                ${sql.value(params.shirt9)}, 
                ${sql.value(params.shirt10)}, 
                ${sql.value(params.shirt11)}, 
                ${sql.value(params.shirt12)}, 
                ${sql.value(params.shirt13)}, 
                ${sql.value(params.extraAttack1)}, 
                ${sql.value(params.extraAttack2)}, 
                ${sql.value(params.extraAttack3)}, 
                ${sql.value(params.extraAttack4)}, 
                ${sql.value(params.extraAttack5)}, 
                ${sql.value(params.extraAttack6)}, 
                ${sql.value(params.extraAttack7)}, 
                ${sql.value(params.extraAttack8)}, 
                ${sql.value(params.extraAttack9)}, 
                ${sql.value(params.extraAttack10)} 
            )`;
        const { text, values } = sql.compile(query);
        await context.pgClient.query(text, values);
        return true;// TODO return something with sense
      },
    }
  };
};

module.exports = resolvers;