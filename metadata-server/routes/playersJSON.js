
module.exports = async (instance, id) => {
    const name = await instance.methods.getName(id).call();

    const schema = {
        "name": name
    };

    return schema;
};
