const isPaying = (auction) => {
    return auction.state === "PAYING";
}

export default {
    isPaying,
}