const postOrderStateChange = (req, res) => {
  try {
    const [body] = req.body;
    const {
      escrow: {
        name,
        status,
        trustee_shortlink: { hash: hashTrusteeShortLink },
        shortlink: { hash },
      },
    } = body;

    console.log(
      `Received: ${name}\n--------\nStatus: ${status}\n--------\nTrustee Shortlink hash: ${hashTrusteeShortLink}\n--------\nShortlink Hash: ${hash}\n--------\n`
    );

    // user horizon service to call notary new mutation

    res.sendStatus(200);
  } catch (e) {
    res.sendStatus(500);
  }
};

module.exports = postOrderStateChange;
