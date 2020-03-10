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
      createSpecialPlayer: async (_, params, context) => {
        const { playerId, name, defence, speed, pass, shoot, endurance, preferredPosition, potential, dayOfBirth } = params;
        const query = sql.query`INSERT INTO players (
              name,
              player_id,
              team_id, 
              defence, 
              speed, 
              pass, 
              shoot, 
              endurance, 
              shirt_number, 
              preferred_position, 
              potential,
              day_of_birth,
              encoded_skills,
              encoded_state,
              red_card,
              injury_matches_left) VALUES (
                ${sql.value(name)},
                ${sql.value(playerId)},
                ${sql.value('1')}, 
                ${sql.value(defence)}, 
                ${sql.value(speed)},
                ${sql.value(pass)},
                ${sql.value(shoot)},
                ${sql.value(endurance)},
                ${sql.value(0)},
                ${sql.value(preferredPosition)},
                ${sql.value(potential)},
                ${sql.value(dayOfBirth)},
                ${sql.value('')},
                ${sql.value('')},
                ${sql.value(0)},
                ${sql.value(0)}
            )`;
        const { text, values } = sql.compile(query);
        await context.pgClient.query(text, values);
        return true;// TODO return something with sense
      },
      deleteSpecialPlayer: async (_, { playerId }, context) => {
        const query = sql.query`DELETE FROM players WHERE team_id=${sql.value('1')} AND player_id=${sql.value(playerId)};`;
        const { text, values } = sql.compile(query);
        await context.pgClient.query(text, values);
        return true;// TODO return something with sense
      },
    }
  };
};

module.exports = resolvers;