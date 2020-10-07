const { insertMessage } = require('../repositories');

const setMessageResolver = async (
  _,
  { input: { id, destinatary, category, auctionId, title, text, customImageUrl, metadata } }
) => {
  try {
    if (id) {
      return Error("Can't accept id in set message");
    }
    const idFromDb = await insertMessage({
      destinatary,
      category,
      auctionId,
      title,
      text,
      customImageUrl,
      metadata,
    });

    return idFromDb;
  } catch (e) {
    return e;
  }
};

module.exports = setMessageResolver;
