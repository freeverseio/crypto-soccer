const config = require('../config.json');

module.exports = async (instance, id) => {
    const name = await instance.methods.getName(id).call();
    const image = config.players_image_base_URL + id;
    const speed = await instance.methods.getSpeed(id).call();
    const defence = await instance.methods.getDefence(id).call();
    const endurance = await instance.methods.getEndurance(id).call();
    const shoot = await instance.methods.getShoot(id).call();
    const pass = await instance.methods.getPass(id).call();

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
