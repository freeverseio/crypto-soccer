const GetSocialIdValidation = require('./GetSocialIdValidation.js');
const Web3 = require('web3');

test('Team Hash', async () => {
  const web3 = new Web3(Web3.givenProvider || 'ws://localhost:8545');
  getSocialIdValidation = new GetSocialIdValidation({
    web3,
    getSocialId: '',
    teamId: '3',
    signature: '',
  });
  const hash = getSocialIdValidation.hash();

  expect(hash).toBe('0x074b4277787bca36334cf57f0507141ef743a08d7690dba02af123626e6955d0');
});

test('Team PrefixedHash', async () => {
  const web3 = new Web3(Web3.givenProvider || 'ws://localhost:8545');
  const signature =
    '3feac668bb718f492638b9b58d1f294379cdc8bde40074f5e49c3f80f28190e121f0fd08227c64a643dd032748ef772b0d1cf1500f649345521c133290c941a91b';
  getSocialIdValidation = new GetSocialIdValidation({
    web3,
    getSocialId: 'ciao',
    teamId: '4',
    signature,
  });
  const hash = getSocialIdValidation.prefixedHash();

  expect(hash).toBe('0x34e71acf56fc4ad1a6e219dcc96bf5111d8092a6ab64308281a7c77525d2a404');
});

test('Team signerAddress', async () => {
  const web3 = new Web3(Web3.givenProvider || 'ws://localhost:8545');
  const signature =
    '3feac668bb718f492638b9b58d1f294379cdc8bde40074f5e49c3f80f28190e121f0fd08227c64a643dd032748ef772b0d1cf1500f649345521c133290c941a91b';
  getSocialIdValidation = new GetSocialIdValidation({
    web3,
    getSocialId: 'ciao',
    teamId: '4',
    signature,
  });
  const signerAddress = await getSocialIdValidation.signerAddress();

  expect(signerAddress).toBe('0x291081e5a1bF0b9dF6633e4868C88e1FA48900e7');
});
