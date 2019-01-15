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

        it("/players/43452 returns 200", () => {
            chai.request(app)
                .get('/players/43452')
                .end((err, res) => {
                    res.should.have.status(200);
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

        it("/teams/43452 returns 200", () => {
            chai.request(app)
                .get('/teams/43452')
                .end((err, res) => {
                    res.should.have.status(200);
                    res.body.should.be.a('object');
                    const json = JSON.parse(res.text);
                    json.name.should.be.equal("Dave Starbelly");
                    json.description.should.be.equal("Friendly OpenSea Creature that enjoys long swims in the ocean.");
                    json.image.should.be.equal("https://storage.googleapis.com/opensea-prod.appspot.com/creature/3.png");
                });
        });
    });
});
