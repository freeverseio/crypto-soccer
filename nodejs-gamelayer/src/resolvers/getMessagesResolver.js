const dayjs = require('dayjs');
const {
  selectMessages,
  selectTeamMailboxStartedAt,
} = require('../repositories');

const getMessagesResolver = async (_, { teamId, limit, after }) => {
  try {
    limit = parseInt(limit) ? parseInt(limit) : null;
    after = parseInt(after) ? parseInt(after) : 0;
    const mailboxStartedAt = await selectTeamMailboxStartedAt({ teamId });
    const isDateValid = dayjs(mailboxStartedAt).isValid();
    const createdAt = isDateValid
      ? mailboxStartedAt
      : dayjs('2020-06-01T16:00:00.000Z').format();

    const messages = await selectMessages({
      destinatary: teamId,
      createdAt,
      after,
      limit,
    });
    return messages;
  } catch (e) {
    return e;
  }
};

module.exports = getMessagesResolver;
