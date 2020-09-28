var cron = require('node-cron');
const selectLastChecked = require('./repositories/selectLastChecked');
const updateLastChecked = require('./repositories/updateLastChecked');
const HorizonService = require('./services/HorizonService');
const GamelayerService = require('./services/GamelayerService');

cron.schedule('* * * * *', async () => {
  console.log('running a task every minute');

  
  //fetch last events from offers(where inserted_at >= last_checked_at)
  // set last_checked_at = fetch from database
  // fetch last events from offers and set new current last_checked_at
  // process events
});
