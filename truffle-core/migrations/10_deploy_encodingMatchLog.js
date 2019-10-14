const CardsAndInjuries = artifacts.require('CardsAndInjuries');

module.exports = function (deployer) {
  deployer.then(async () => {
      await deployer.deploy(CardsAndInjuries);
    })
    .catch(console.error);
};

