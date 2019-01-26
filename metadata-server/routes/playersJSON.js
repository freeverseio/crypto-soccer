const config = require('../config.json');

module.exports = async (instance, id) => {
    try {
        var name = await instance.methods.getName(id).call();
        var image = config.players_image_base_URL + id;
        var speed = await instance.methods.getSpeed(id).call();
        var defence = await instance.methods.getDefence(id).call();
        var endurance = await instance.methods.getEndurance(id).call();
        var shoot = await instance.methods.getShoot(id).call();
        var pass = await instance.methods.getPass(id).call();
    }
    catch (err) {
        console.error(err);
        return {};
    }

    const schema = {
        "name": name,
        "description": "put a description",
        "image": image,
        "attributes": [
            {
                "trait_type": "speed",
                "value": speed
            },
            {
                "trait_type": "defence",
                "value": defence
            },
            {
                "trait_type": "endurance",
                "value": endurance
            },
            {
                "trait_type": "shoot",
                "value": shoot
            },
            {
                "trait_type": "pass",
                "value": pass
            }
        ]
    };

    return schema;
};
