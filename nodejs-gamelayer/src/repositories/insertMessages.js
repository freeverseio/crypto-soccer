const PostgresSQLService = require('../services/PostgresSQLService');
const copyFrom = require('pg-copy-streams').from;
const Readable = require('stream').Readable;

const insertMessages = async (messages) => {
  const pool = await PostgresSQLService.getPool();
  try {
    const client = await pool.connect();
    const done = () => {
      client.release();
    };
    const stream = client.query(
      copyFrom(
        'COPY inbox (destinatary, category, auction_id, title, text_message, custom_image_url, metadata) FROM STDIN'
      )
    );
    const rs = new Readable();
    let currentIndex = 0;
    rs._read = () => {
      if (currentIndex === messages.length) {
        rs.push(null);
      } else {
        const message = messages[currentIndex];
        rs.push(
          message.destinatary +
            '\t' +
            message.category +
            '\t' +
            message.auctionId +
            '\t' +
            message.title +
            '\t' +
            message.text +
            '\t' +
            message.customImageUrl +
            '\t' +
            message.metadata +
            '\n'
        );
        currentIndex = currentIndex + 1;
      }
    };
    const onError = (strErr) => {
      console.error('Something went wrong on broadcast message:', strErr);
      done();
    };
    rs.on('error', onError);
    stream.on('error', onError);
    stream.on('end', done);
    rs.pipe(stream);
  } catch (e) {
    throw e;
  }
};

module.exports = insertMessages;
