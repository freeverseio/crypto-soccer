const canBePutOnSale = (auction) => {
    if (!auction) return true;
    if (!auction.state) return true;

    return auction.state === "NO_BIDS" ||
        auction.state === "CANCELLED_BY_SELLER" ||
        auction.state === "FAILED";
}

const isWaitingPayment = (auction) => {
    if (!auction) return false;
    if (!auction.state) return false;

    return auction.state === "PAYING";
}

const isWaitingWithdrawal = (auction) => {
    if (!auction) return false;
    if (!auction.state) return false;

    return auction.state === "WITHDRAWAL";
}

export default {
    canBePutOnSale,
    isWaitingPayment,
    isWaitingWithdrawal,
}