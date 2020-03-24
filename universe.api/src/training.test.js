const { isTrainingGroupValid, isTrainingSpecialPlayerValid } = require('./training');

describe('training', () => {
    describe('group', () => {
        const TP = 10;

        test('group with sum 0', () => {
            expect(() => isTrainingGroupValid(TP, 0, 0, 0, 0, 0)).not.toThrow();
        });

        test('group with sum exceeding the TP', () => {
            expect(() => isTrainingGroupValid(TP, 5, 5, 5, 0, 0)).toThrow("group sum 15 exceeds available TP 10");
        });

        test('each element is 60% of TP', () => {
            expect(() => isTrainingGroupValid(TP, 6, 0, 0, 0, 0)).not.toThrow();
            expect(() => isTrainingGroupValid(TP, 0, 6, 0, 0, 0)).not.toThrow();
            expect(() => isTrainingGroupValid(TP, 0, 0, 6, 0, 0)).not.toThrow();
            expect(() => isTrainingGroupValid(TP, 0, 0, 0, 6, 0)).not.toThrow();
            expect(() => isTrainingGroupValid(TP, 0, 0, 0, 0, 6)).not.toThrow();
        });

        test('each element exceeding 60% of TP', () => {
            expect(() => isTrainingGroupValid(TP, 7, 0, 0, 0, 0)).toThrow("shoot exceeds 60% of TP 10");
            expect(() => isTrainingGroupValid(TP, 0, 7, 0, 0, 0)).toThrow("speed exceeds 60% of TP 10");
            expect(() => isTrainingGroupValid(TP, 0, 0, 7, 0, 0)).toThrow("pass exceeds 60% of TP 10");
            expect(() => isTrainingGroupValid(TP, 0, 0, 0, 7, 0)).toThrow("defence exceeds 60% of TP 10");
            expect(() => isTrainingGroupValid(TP, 0, 0, 0, 0, 7)).toThrow("endurance exceeds 60% of TP 10");
        });

        test('group with sum < TP', () => {
            expect(() => isTrainingGroupValid(TP, 2, 2, 2, 1, 2)).not.toThrow();
        });
    });

    describe('special player', () => {
        const TP = 10;

        test('group with sum 0', () => {
            expect(() => isTrainingSpecialPlayerValid(TP, 0, 0, 0, 0, 0)).not.toThrow();
        });

        test('group with sum exceeding the TP', () => {
            expect(() => isTrainingSpecialPlayerValid(TP, 5, 5, 5, 0, 0)).toThrow("group sum 15 exceeds available TP 11");
        });

        test('each element is 60% of TP', () => {
            expect(() => isTrainingSpecialPlayerValid(TP, 6, 0, 0, 0, 0)).not.toThrow();
            expect(() => isTrainingSpecialPlayerValid(TP, 0, 6, 0, 0, 0)).not.toThrow();
            expect(() => isTrainingSpecialPlayerValid(TP, 0, 0, 6, 0, 0)).not.toThrow();
            expect(() => isTrainingSpecialPlayerValid(TP, 0, 0, 0, 6, 0)).not.toThrow();
            expect(() => isTrainingSpecialPlayerValid(TP, 0, 0, 0, 0, 6)).not.toThrow();
        });

        test('each element exceeding 60% of TP', () => {
            expect(() => isTrainingSpecialPlayerValid(TP, 7, 0, 0, 0, 0)).toThrow("shoot exceeds 60% of TP 11");
            expect(() => isTrainingSpecialPlayerValid(TP, 0, 7, 0, 0, 0)).toThrow("speed exceeds 60% of TP 11");
            expect(() => isTrainingSpecialPlayerValid(TP, 0, 0, 7, 0, 0)).toThrow("pass exceeds 60% of TP 11");
            expect(() => isTrainingSpecialPlayerValid(TP, 0, 0, 0, 7, 0)).toThrow("defence exceeds 60% of TP 11");
            expect(() => isTrainingSpecialPlayerValid(TP, 0, 0, 0, 0, 7)).toThrow("endurance exceeds 60% of TP 11");
        });

        test('group with sum < TP', () => {
            expect(() => isTrainingSpecialPlayerValid(TP, 2, 2, 2, 1, 2)).not.toThrow();
        });
    });
});