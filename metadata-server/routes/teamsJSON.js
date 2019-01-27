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
        "description": "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum",
        "image": "http://www.monkers.net/fotos/lucasartsDead.jpg",
        "external_url": "https://www.freeverse.io/",
        "attributes": [
            {
                "trait_type": "team",
                "value": name
            }
        ]
    };

    return schema;
};
