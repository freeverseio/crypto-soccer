const processOffersHistories = async () => {
  const offerLastChecked = await selectLastChecked({ entity: 'offer' });

  const lastOffersHistories = await HorizonService.getLastOfferHistories({
    lastChecked: offerLastChecked,
  });

  const newLastChecked =
    lastOffersHistories[lastOffersHistories.length].insertedAt;

  await updateLastChecked({ entity: 'offer', lastChecked: newLastChecked });

  for (const offerHistory of lastOffersHistories) {
    switch (offerHistory.state) {
      case 'started':
        //notify seller of new offer received which is the same as new higher offer
        await GamelayerService.setMessage({
          destinatary: offerHistory.seller,
          category: 'offer',
          auctionId: offerHistory.auctionId,
          text: 'Blablbalba',
          customImageUrl: '',
          metadata: '',
        });
        //notify all the buyers(all offers where playerId=playerId and state is started) that new higher offer has been made
        const offerers = await HorizonService.getOfferersByPlayerId({
          playerId,
        });
        for (const offerer of offerers) {
          await GamelayerService.setMessage({
            destinatary: offerer,
            category: 'offer',
            auctionId: offerHistory.auctionId,
            text: 'Blablbalba',
            customImageUrl: '',
            metadata: '',
          });
        }
        break;

      case 'accepted':
        //notify buyer of accepted offer
        await GamelayerService.setMessage({
          destinatary: offerHistory.buyer,
          category: 'offer',
          auctionId: offerHistory.auctionId,
          text: 'Blablbalba',
          customImageUrl: '',
          metadata: '',
        });
        break;
      case 'rejected':
        //notify buyer of rejected offer
        await GamelayerService.setMessage({
          destinatary: offerHistory.buyer,
          category: 'offer',
          auctionId: offerHistory.auctionId,
          text: 'Blablbalba',
          customImageUrl: '',
          metadata: '',
        });
        break;
    }
  }
};

module.exports = processOffersHistories;
