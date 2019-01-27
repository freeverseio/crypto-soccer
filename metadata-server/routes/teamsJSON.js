const config = require('../config.json');

module.exports = async ({playersContract, teamsContract, teamId}) => {
    try {
        var name = await teamsContract.methods.getName(teamId).call();
        var players = await teamsContract.methods.getPlayers(teamId).call();
        console.log(players);
    }
    catch (err) {
        console.error(err);
        return {};
    }

    const schema = {
        "name": name,
        "description": "It's a team... a good one for sure.",
        "image": "http://www.monkers.net/fotos/lucasartsDead.jpg",
        "external_url": "https://www.freeverse.io/",
        "attributes": [
            {
                "trait_type": "team",
                "value": teamName
            },
            {
                "display_type": "number",
                "trait_type": "generation",
                "value": 0
            }
        ]
    };

    return schema;
};
