const dayjs = require('dayjs');
const { selectMessages, selectTeamMailboxStartedAt } = require('../repositories');
const messagesView = require('../views/message');

const getMessagesResolver = async (_, { teamId, auctionId, limit, offset }) => {
  try {
    limit = parseInt(limit) ? parseInt(limit) : null;
    offset = parseInt(offset) ? parseInt(offset) : 0;
    const mailboxStartedAt = await selectTeamMailboxStartedAt({ teamId });
    const isDateValid = dayjs(mailboxStartedAt).isValid();
    const createdAt = isDateValid ? mailboxStartedAt : dayjs('2020-06-01T16:00:00.000Z').format();

    const messages = await selectMessages({
      destinatary: teamId,
      auctionId,
      createdAt,
      offset,
      limit,
    });
    return { totalCount: messages.length, nodes: messages.map(messagesView) };
  } catch (e) {
    return e;
  }
};

module.exports = getMessagesResolver;
