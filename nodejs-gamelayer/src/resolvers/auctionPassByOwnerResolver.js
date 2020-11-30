const auctionPassByOwnerResolver = async (parent, args, context, info, schema) => {
  return await info.mergeInfo.delegateToSchema({
    schema,
    operation: 'query',
    fieldName: 'hasAuctionPass',
    args: {
      owner: parent.owner,
    },
    context,
    info,
  });
};

module.exports = auctionPassByOwnerResolver;
