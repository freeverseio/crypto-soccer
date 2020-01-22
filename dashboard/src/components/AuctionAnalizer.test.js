import analizer from './AuctionAnalizer';

it('isWaitingPayment', () => {
    expect(analizer.isWaitingPayment(undefined)).toBe(false);
    expect(analizer.isWaitingPayment(null)).toBe(false);
    expect(analizer.isWaitingPayment({})).toBe(false);
    expect(analizer.isWaitingPayment({"state": "STARTED"})).toBe(false);
    expect(analizer.isWaitingPayment({"state": "ASSET_FROZEN"})).toBe(false);
    expect(analizer.isWaitingPayment({"state": "PAYING"})).toBe(true);
    expect(analizer.isWaitingPayment({"state": "PAID"})).toBe(false);
    expect(analizer.isWaitingPayment({"state": "WITHDRAWAL"})).toBe(false);
    expect(analizer.isWaitingPayment({"state": "NO_BIDS"})).toBe(false);
    expect(analizer.isWaitingPayment({"state": "CANCELLED_BY_SELLER"})).toBe(false);
    expect(analizer.isWaitingPayment({"state": "FAILED"})).toBe(false);
});

it('canBePutOnSale', () => {
    expect(analizer.canBePutOnSale(undefined)).toBe(true);
    expect(analizer.canBePutOnSale(null)).toBe(true);
    expect(analizer.canBePutOnSale({})).toBe(true);
    expect(analizer.canBePutOnSale({"state": "STARTED"})).toBe(false);
    expect(analizer.canBePutOnSale({"state": "ASSET_FROZEN"})).toBe(false);
    expect(analizer.canBePutOnSale({"state": "PAYING"})).toBe(false);
    expect(analizer.canBePutOnSale({"state": "PAID"})).toBe(false);
    expect(analizer.canBePutOnSale({"state": "WITHDRAWAL"})).toBe(false);
    expect(analizer.canBePutOnSale({"state": "NO_BIDS"})).toBe(true);
    expect(analizer.canBePutOnSale({"state": "CANCELLED_BY_SELLER"})).toBe(true);
    expect(analizer.canBePutOnSale({"state": "FAILED"})).toBe(true);
});

it('isWaitingWithdrawal', () => {
    expect(analizer.isWaitingWithdrawal(undefined)).toBe(false);
    expect(analizer.isWaitingWithdrawal(null)).toBe(false);
    expect(analizer.isWaitingWithdrawal({})).toBe(false);
    expect(analizer.isWaitingWithdrawal({"state": "STARTED"})).toBe(false);
    expect(analizer.isWaitingWithdrawal({"state": "ASSET_FROZEN"})).toBe(false);
    expect(analizer.isWaitingWithdrawal({"state": "PAYING"})).toBe(false);
    expect(analizer.isWaitingWithdrawal({"state": "PAID"})).toBe(false);
    expect(analizer.isWaitingWithdrawal({"state": "WITHDRAWAL"})).toBe(true);
    expect(analizer.isWaitingWithdrawal({"state": "NO_BIDS"})).toBe(false);
    expect(analizer.isWaitingWithdrawal({"state": "CANCELLED_BY_SELLER"})).toBe(false);
    expect(analizer.isWaitingWithdrawal({"state": "FAILED"})).toBe(false);
});
