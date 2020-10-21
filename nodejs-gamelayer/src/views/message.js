const dayjs = require('dayjs');

const message = ({
  id,
  destinatary,
  category,
  auction_id,
  title,
  text,
  custom_image_url,
  metadata,
  is_read,
  created_at,
}) => {
  return {
    id,
    destinatary,
    category,
    auctionId: auction_id,
    title,
    text,
    customImageUrl: custom_image_url,
    metadata,
    isRead: is_read,
    createdAt: dayjs(created_at).format(),
  };
};

module.exports = message;
