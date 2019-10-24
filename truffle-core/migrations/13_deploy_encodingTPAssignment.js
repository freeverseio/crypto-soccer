const EncodingTPAssignment = artifacts.require('EncodingTPAssignment');

module.exports = function (deployer) {
  deployer.then(async () => {
      await deployer.deploy(EncodingTPAssignment);
    })
    .catch(console.error);
};

