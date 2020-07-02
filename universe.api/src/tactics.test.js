const { checkTactics } = require('./tactics');

describe('tactics', () => {
    describe('description of tests', () => {
        const TP = 10;

        test('name of test', () => {
            expect(() => checkTactics()).not.toThrow();
        });
    });


});