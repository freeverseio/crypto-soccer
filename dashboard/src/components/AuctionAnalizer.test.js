import analizer from './AuctionAnalizer';

it('isPaying with PAYING auction', () => {
    const auction = {
        "state": "PAYING",
    };

    const result = analizer.isPaying(auction)
    expect(result).toBe(true);
});

it('isPaying with PAID auction', () => {
    const auction = {
        "state": "PAID",
    };

    const result = analizer.isPaying(auction)
    expect(result).toBe(false);
});