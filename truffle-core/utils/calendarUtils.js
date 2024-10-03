    // Inputs:
    // - verse: the verse to be played
    // - TZForRound1: the timezone that played the very first games
    // Outputs: for the verse to be played:
    // - timezone: which timezone plays
    // - matchDay: what matchDay of the league it corresponds to: a number in [0, 13]
    // - half: whether it corresponds to the first half (half = 0), or to the second half (half = 1)
    function nextTimeZoneToPlay(verse, TZForRound1) {
      const NULL_TIMEZONE = 0; 
      const VERSES_PER_DAY = 96; 
      const MATCHDAYS_PER_ROUND = 14;

      // if _currentVerse = 0, we should be updating TZForRound1
      // recall that timeZones range from 1...24 (not from 0...24)

      let half = verse % 4;
      const delta = 9 * 4 + half;
      let tz, dia;

      // if the half is greater or equal to 2 and verse is less than delta, return NULL_TIMEZONE
      if (half >= 2 && verse < delta) return { timezone: NULL_TIMEZONE, matchDay: 0, turnInDay: 0 };

      if (half < 2) {
          tz = TZForRound1 + Math.floor((verse - half) % VERSES_PER_DAY / 4);
          dia = 2 * Math.floor((verse - 4 * (tz - TZForRound1) - half) / VERSES_PER_DAY);
      } else {
          tz = TZForRound1 + Math.floor((verse - delta) % VERSES_PER_DAY / 4);
          dia = 1 + 2 * Math.floor((verse - 4 * (tz - TZForRound1) - delta) / VERSES_PER_DAY);
          half -= 2;
      }

      const timezone = normalizeTZ(tz);
      const matchDay = dia % MATCHDAYS_PER_ROUND;

      return { timezone, matchDay, half };
  }

  function normalizeTZ(tz) {
      return 1 + ((24 + tz - 1) % 24);
  }

  // Returns the round at which a league is (the first league played is round 0, the next league is round 1, etc.)
  function getCurrentRound(tz, TZForRound1, verse) {
      const VERSES_PER_ROUND = 672; /// 96 * 7days
      if (verse < VERSES_PER_ROUND) return 0;
      const round = Math.floor(verse / VERSES_PER_ROUND);
      const deltaN = (tz >= TZForRound1) ? (tz - TZForRound1) : ((tz + 24) - TZForRound1);
      if (verse < 4 * deltaN + round * VERSES_PER_ROUND) {
          return round - 1;
      } else {
          return round;
      }
  }


  // Returns the Unix timestamp in UTC (seconds) corresponding to the start of a match's first half 
  // Inputs:
  // - tz: the timezone where the match belongs
  // - round: the round of a league (the first league played is round 0, the next league is round 1, etc.)
  // - matchDay: what matchDay of the league it corresponds to: a number in [0, 13]
  // - TZForRound1: the timezone that played the very first games
  // - firstVerseTimeStamp: the timestamp the very first games where played at 
  // Outputs:
  // - timeUTC: the Unix timestamp in UTC (seconds) corresponding to the start of a match's first half 
  function getMatch1stHalfUTC(tz, round, matchDay, TZForRound1, firstVerseTimeStamp) {
      const DAYS_PER_ROUND = 7;
      if (tz <= 0 || tz >= 25) {
          throw new Error("timezone out of range");
      }

      const deltaN = (tz >= TZForRound1) ? 
          (tz - TZForRound1) : 
          ((tz + 24) - TZForRound1);

      let timeUTC;
      if (matchDay % 2 === 0) {
          timeUTC = firstVerseTimeStamp + (deltaN + 12 * matchDay + 24 * DAYS_PER_ROUND * round) * 3600;
      } else {
          timeUTC = firstVerseTimeStamp + (19 + 2 * deltaN + 24 * (matchDay - 1) + 48 * DAYS_PER_ROUND * round) * 1800;
      }
      return timeUTC;
  }

  function getMatchHalfUTC(tz, round, matchDay, half, TZForRound1, firstVerseTimeStamp) {
      const SECS_BETWEEN_VERSES = 900; /// 15 mins
      const extraSeconds = half == 0 ? 0 : SECS_BETWEEN_VERSES;
      return getMatch1stHalfUTC(tz, round, matchDay, TZForRound1, firstVerseTimeStamp) + extraSeconds;
  }

  function calendarInfo(verse, TZForRound1, firstVerseTimeStamp) {
      const NULL_TIMEZONE = 0;
      const timeZoneData = nextTimeZoneToPlay(verse, TZForRound1);
      if (timeZoneData.timezone == NULL_TIMEZONE) {
          return {timezone: NULL_TIMEZONE, matchDay: null, half: null, leagueRound: null, timestamp: null};
      }
      const leagueRound = getCurrentRound(timeZoneData.timezone, TZForRound1, verse);
      const timestamp = getMatchHalfUTC(timeZoneData.timezone, leagueRound, timeZoneData.matchDay, timeZoneData.half, TZForRound1, firstVerseTimeStamp);
      return {
          timezone: timeZoneData.timezone,
          matchDay: timeZoneData.matchDay,
          half: timeZoneData.half,
          leagueRound,
          timestamp: timestamp,
      };
  }

  module.exports = {
    nextTimeZoneToPlay,
    getCurrentRound,
    getMatch1stHalfUTC,
    getMatchHalfUTC,
    calendarInfo,
  }