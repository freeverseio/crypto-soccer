const { Pool } = require('pg');
const postgreSQLConfig = require('../config');

class PostgresSQLService {
  constructor() {
    this.connectionString = postgreSQLConfig.connectionString;
    this.pool = null;
  }

  getPool() {
    if (!this.pool) {
      this.pool = new Pool({ connectionString: this.connectionString });
    }

    return this.pool;
  }
}

module.exports = new PostgresSQLService();
