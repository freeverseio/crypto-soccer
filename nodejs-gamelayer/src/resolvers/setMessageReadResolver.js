const { updateMessageRead } = require('../repositories');

const setMessageReadResolver = async (_, { id }) => {
  try {
    await updateMessageRead({ id });

    return true;
  } catch (e) {
    return e;
  }
};

module.exports = setMessageReadResolver;
