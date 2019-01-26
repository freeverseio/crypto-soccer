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
        "image": "https://srv.latostadora.com/designall.dll/guybrush_threepwood--i:1413852880551413850;w:520;m:1;b:FFFFFF.jpg",
        "external_url": "https://www.freeverse.io/",
        "attributes": [
            {
                "trait_type": "speed",
                "value": Number(speed)
            },
            {
                "trait_type": "defence",
                "value": Number(defence)
            },
            {
                "trait_type": "endurance",
                "value": Number(endurance)
            },
            {
                "trait_type": "shoot",
                "value": Number(shoot)
            },
            {
                "trait_type": "pass",
                "value": Number(pass)
            }
        ]
    };

    return schema;
};
