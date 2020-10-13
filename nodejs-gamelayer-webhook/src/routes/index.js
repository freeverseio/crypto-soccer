const express = require('express');
const router = new express.Router();
const postOrderStateChange = require('../handlers/postOrderStateChange');

router.route('/order/status').post(postOrderStateChange);

module.exports = router;
