const config = require('../config.json');

module.exports = async ({playersContract, teamsContract, playerId}) => {
    try {
        var name = await playersContract.methods.getName(playerId).call();
        var image = config.players_image_base_URL + playerId;
        var speed = await playersContract.methods.getSpeed(playerId).call();
        var defence = await playersContract.methods.getDefence(playerId).call();
        var endurance = await playersContract.methods.getEndurance(playerId).call();
        var shoot = await playersContract.methods.getShoot(playerId).call();
        var pass = await playersContract.methods.getPass(playerId).call();
        const teamId = await playersContract.methods.getTeam(playerId).call();
        var teamName = await teamsContract.methods.getName(teamId).call();
    }
    catch (err) {
        console.error(err);
        return {};
    }

    const schema = {
        "name": name,
        "description": "Edson Arantes do Nascimento (Brazilian Portuguese: [ˈɛtsõ (w)ɐˈɾɐ̃tʃiz du nɐsiˈmẽtu]; born 23 October 1940), known as Pelé ([peˈlɛ]), is a Brazilian retired professional footballer who played as a forward. He is regarded by many in the sport, including football writers, players, and fans, as the greatest player of all time. In 1999, he was voted World Player of the Century by the International Federation of Football History & Statistics (IFFHS), and was one of the two joint winners of the FIFA Player of the Century award. That same year, Pelé was elected Athlete of the Century by the International Olympic Committee. According to the IFFHS, Pelé is the most successful domestic league goal-scorer in football history scoring 650 goals in 694 League matches, and in total 1281 goals in 1363 games, which included unofficial friendlies and is a Guinness World Record.[1][2][3][4][5] During his playing days, Pelé was for a period the best-paplayerId athlete in the world.",
        "image": "https://srv.latostadora.com/designall.dll/guybrush_threepwood--i:1413852880551413850;w:520;m:1;b:FFFFFF.jpg",
        "external_url": "https://www.freeverse.io/",
        "attributes": [
            {
                "trait_type": "speed",
                "value": Number(speed),
                "max_value": 100
            },
            {
                "trait_type": "defence",
                "value": Number(defence),
                "max_value": 100
            },
            {
                "trait_type": "endurance",
                "value": Number(endurance),
                "max_value": 100
            },
            {
                "trait_type": "shoot",
                "value": Number(shoot),
                "max_value": 100
            },
            {
                "trait_type": "pass",
                "value": Number(pass),
                "max_value": 100
            },
            {
                "trait_type": "team",
                "value": teamName
            },
            {
                "display_type": "boost_number",
                "trait_type": "shoot_power",
                "value": 10
            },
            {
                "display_type": "boost_percentage",
                "trait_type": "pass_increase",
                "value": 5
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
