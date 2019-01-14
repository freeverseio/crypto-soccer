var express = require('express');
const schema = require('./example_schema.json');

var router = express.Router();

/* GET JSON schema for teams with id. */
router.get('/:id', function (req, res, next) {
    const id = req.params.id;
    res.send(schema);
});

module.exports = router;

