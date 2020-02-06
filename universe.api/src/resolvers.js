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
            ) ON CONFLICT (verse, team_id) DO UPDATE SET
                tactic_id=${sql.value(params.tacticId)},
                shirt_0=${sql.value(params.shirt0)},
                shirt_1=${sql.value(params.shirt1)},
                shirt_2=${sql.value(params.shirt2)},
                shirt_3=${sql.value(params.shirt3)},
                shirt_4=${sql.value(params.shirt4)},
                shirt_5=${sql.value(params.shirt5)},
                shirt_6=${sql.value(params.shirt6)},
                shirt_7=${sql.value(params.shirt7)},
                shirt_8=${sql.value(params.shirt8)},
                shirt_9=${sql.value(params.shirt9)},
                shirt_10=${sql.value(params.shirt10)},
                shirt_11=${sql.value(params.shirt11)},
                shirt_12=${sql.value(params.shirt12)},
                shirt_13=${sql.value(params.shirt13)},
                extra_attack_1=${sql.value(params.extraAttack1)},
                extra_attack_2=${sql.value(params.extraAttack2)},
                extra_attack_3=${sql.value(params.extraAttack3)},
                extra_attack_4=${sql.value(params.extraAttack4)},
                extra_attack_5=${sql.value(params.extraAttack5)},
                extra_attack_6=${sql.value(params.extraAttack6)},
                extra_attack_7=${sql.value(params.extraAttack7)},
                extra_attack_8=${sql.value(params.extraAttack8)},
                extra_attack_9=${sql.value(params.extraAttack9)},
                extra_attack_10=${sql.value(params.extraAttack10)}
              `;
        const { text, values } = sql.compile(query);
        await context.pgClient.query(text, values);
        return true;// TODO return something with sense
      },
      setTraining: async (_, params, context) => {
        const query = sql.query`INSERT INTO trainings (
			verse,
			team_id,
    		special_player_shirt,
			goalkeepers_defence,
    		goalkeepers_speed,
    		goalkeepers_pass,
    		goalkeepers_shoot,
    		goalkeepers_endurance,
    		defenders_defence,
    		defenders_speed,
    		defenders_pass,
    		defenders_shoot,
    		defenders_endurance,
    		midfielders_defence,
    		midfielders_speed,
    		midfielders_pass,
    		midfielders_shoot,
    		midfielders_endurance,
    		attackers_defence,
    		attackers_speed,
    		attackers_pass,
    		attackers_shoot,
    		attackers_endurance,
    		special_player_defence,
    		special_player_speed,
    		special_player_pass,
    		special_player_shoot,
			special_player_endurance
		) VALUES (
                ${sql.value('9223372036854775807')},
                ${sql.value(params.teamId)},
                ${sql.value(params.specialPlayerShirt)},
                ${sql.value(params.goalkeepersDefence)},
                ${sql.value(params.goalkeepersSpeed)},
                ${sql.value(params.goalkeepersPass)},
                ${sql.value(params.goalkeepersShoot)},
                ${sql.value(params.goalkeepersEndurance)},
                ${sql.value(params.defendersDefence)},
                ${sql.value(params.defendersSpeed)},
                ${sql.value(params.defendersPass)},
                ${sql.value(params.defendersShoot)},
                ${sql.value(params.defendersEndurance)},
                ${sql.value(params.midfieldersDefence)},
                ${sql.value(params.midfieldersSpeed)},
                ${sql.value(params.midfieldersPass)},
                ${sql.value(params.midfieldersShoot)},
                ${sql.value(params.midfieldersEndurance)},
                ${sql.value(params.attackersDefence)},
                ${sql.value(params.attackersSpeed)},
                ${sql.value(params.attackersPass)},
                ${sql.value(params.attackersShoot)},
                ${sql.value(params.attackersEndurance)},
                ${sql.value(params.specialPlayerDefence)},
                ${sql.value(params.specialPlayerSpeed)},
                ${sql.value(params.specialPlayerPass)},
                ${sql.value(params.specialPlayerShoot)},
                ${sql.value(params.specialPlayerEndurance)}
            ) ON CONFLICT (verse, team_id) DO UPDATE SET
        special_player_shirt=${sql.value(params.specialPlayerShirt)},
			  goalkeepers_defence=${sql.value(params.goalkeepersDefence)},
    		goalkeepers_speed=${sql.value(params.goalkeepersSpeed)},
    		goalkeepers_pass=${sql.value(params.goalkeepersPass)},
    		goalkeepers_shoot=${sql.value(params.goalkeepersShoot)},
        goalkeepers_endurance=${sql.value(params.goalkeepersEndurance)},
        defenders_defence=${sql.value(params.defendersDefence)},
    		defenders_speed=${sql.value(params.defendersSpeed)},
    		defenders_pass=${sql.value(params.defendersPass)},
    		defenders_shoot=${sql.value(params.defendersShoot)},
    		defenders_endurance=${sql.value(params.defendersEndurance)},
        midfielders_defence=${sql.value(params.midfieldersDefence)},
    		midfielders_speed=${sql.value(params.midfieldersSpeed)},
    		midfielders_pass=${sql.value(params.midfieldersPass)},
    		midfielders_shoot=${sql.value(params.midfieldersShoot)},
    		midfielders_endurance=${sql.value(params.midfieldersEndurance)},
        attackers_defence=${sql.value(params.attackersDefence)},
    		attackers_speed=${sql.value(params.attackersSpeed)},
    		attackers_pass=${sql.value(params.attackersPass)},
    		attackers_shoot=${sql.value(params.attackersShoot)},
    		attackers_endurance=${sql.value(params.attackersEndurance)},
        special_player_defence=${sql.value(params.specialPlayerDefence)},
    		special_player_speed=${sql.value(params.specialPlayerSpeed)},
    		special_player_pass=${sql.value(params.specialPlayerPass)},
    		special_player_shoot=${sql.value(params.specialPlayerShoot)},
    		special_player_endurance=${sql.value(params.specialPlayerEndurance)}
        `;
        const { text, values } = sql.compile(query);
        await context.pgClient.query(text, values);
        return true;// TODO return something with sense
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
              frozen,
              red_card_matches_left,
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
                ${sql.value(false)},
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