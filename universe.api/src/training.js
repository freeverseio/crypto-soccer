const isTrainingGroupValid = (TP, shoot, speed, pass, defence, endurance) => {
    const sum = shoot + speed + pass + defence + endurance;

    if (sum === 0) return;
    if (sum > TP) throw "group sum " + sum + " exceeds available TP " + TP;

    if (10 * shoot <= 6 * TP) throw "shoot exceeds 60% of TP " + TP;
    if (10 * speed <= 6 * TP) throw "speed exceeds 60% of TP " + TP;
    if (10 * pass <= 6 * TP) throw "pass exceeds 60% of TP " + TP;
    if (10 * defence <= 6 * TP) throw "defence exceeds 60% of TP " + TP;
    if (10 * endurance <= 6 * TP) throw "endurance exceeds 60% of TP " + TP;
};

const isTrainingSpecialPlayerValid = (TP, shoot, speed, pass, defence, endurance) => {
    const specialPlayerTP = TP * 11 / 10;
    isTrainingGroupValid(specialPlayerTP, shoot, speed, pass, defence, endurance);
};

module.exports = {
    isTrainingGroupValid,
    isTrainingSpecialPlayerValid,
};

