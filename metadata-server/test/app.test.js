// Import the dependencies for testing
const chai = require('chai');
const chaiHttp = require('chai-http');
const app = require('../app');

// Configure chai
chai.use(chaiHttp);
chai.should();

describe("routing", () => {
    it("/ returns 404", () => {
        chai.request(app)
            .get('/')
            .end((err, res) => {
                res.should.have.status(404);
                res.body.should.be.a('object');
            });
    });

    describe("players", () => {
        it("/players returns 404", () => {
            chai.request(app)
                .get('/players')
                .end((err, res) => {
                    res.should.have.status(404);
                    res.body.should.be.a('object');
                });
        });
    });

    describe("teams", () => {
        it("/teams returns 404", () => {
            chai.request(app)
                .get('/teams')
                .end((err, res) => {
                    res.should.have.status(404);
                    res.body.should.be.a('object');
                });
        });
    });

    describe("images", () => {
        it("/images/README.md" , () => {
            chai.request(app)
                .get('/images/README.md')
                .end((err, res) => {
                    res.should.have.status(200);
                    res.body.should.be.a('object');
                });
        });
    });
});
