module.exports = {
  rules: {
    "max-line-length": [
      true,
      {
        limit: 180,
        "ignore-pattern": "^import [^,]+ from |^export | implements",
      },
    ],
  }
};
