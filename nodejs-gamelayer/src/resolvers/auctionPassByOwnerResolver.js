const auctionPassByOwnerResolver = async (parent, args, context, info, schema) => {
  const result = await info.mergeInfo.delegateToSchema({
    schema,
    operation: 'query',
    fieldName: 'hasAuctionPass',
    args: {
      input: { owner: parent.owner },
    },
    context,
    info,
  });
  if (!result) {
    return false;
  }

  return result;
};

module.exports = auctionPassByOwnerResolver;
