const dayjs = require('dayjs');
const { selectNumUnreadMessages, selectTeamMailboxStartedAt } = require('../repositories');

const getNumUnreadMessagesResolver = async (teamId) => {
  const mailboxStartedAt = await selectTeamMailboxStartedAt({ teamId });
  const isDateValid = dayjs(mailboxStartedAt).isValid();
  const createdAt = isDateValid ? mailboxStartedAt : dayjs('2020-06-01T16:00:00.000Z').format();

  const { num: numUnreadMessages } = await selectNumUnreadMessages({ destinatary: teamId, createdAt });

  return parseInt(numUnreadMessages);
};

module.exports = getNumUnreadMessagesResolver;
