const EncodingSkills = artifacts.require('EncodingSkills');
const EncodingState = artifacts.require('EncodingState');

module.exports = function (deployer) {
  deployer.then(async () => {
      await deployer.deploy(EncodingSkills);
      await deployer.deploy(EncodingState);
    })
    .catch(console.error);
};

