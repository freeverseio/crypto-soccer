const { isTrainingGroupValid, isTrainingSpecialPlayerValid } = require('./training');

describe('training', () => {
  test('0 group sum', () => {
    const TP = 10;
    expect(() => isTrainingGroupValid(TP,0,0,0,0,0)).not.toThrow();
  });
});