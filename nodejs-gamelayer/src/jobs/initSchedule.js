var cron = require('node-cron');
var LeaderboardService = require('../services/LeaderboardService');

const initSchedule = () => {
  cron.schedule(
    '15 13,22 * * *',
    () => {
      console.log('Updating gamelayer leaderboards');
      LeaderboardService.saveLeaderboards();
    },
    {
      scheduled: true,
      timezone: 'Europe/Madrid',
    }
  );
};

module.exports = initSchedule;
