const { Pool } = require('pg');

class PostgresSQLService {
  constructor() {
    this.configPG = {
      database: 'game',
      host: 'gamedb',
      max: 10,
      password: 'freeverse',
      port: 5432,
      user: 'freeverse',
    };
    this.pool = null;
  }

  getPool() {
    if (!this.pool) {
      this.pool = new Pool(this.configPG);
    }

    return this.pool;
  }
}

module.exports = new PostgresSQLService();
