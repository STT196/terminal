export const secret = {
  AirtableSecret: new sst.Secret("AirtableSecret"),
  ForgeKey: new sst.Secret("ForgeKey"),
  GithubClientID: new sst.Secret("GithubClientID"),
  GithubClientSecret: new sst.Secret("GithubClientSecret"),
  TwitchClientID: new sst.Secret("TwitchClientID"),
  TwitchClientSecret: new sst.Secret("TwitchClientSecret"),
  SlackWebhook: new sst.Secret("SlackWebhook"),
  IpinfoToken: new sst.Secret("IpinfoToken"),
  AnthropicApiKey: new sst.Secret("AnthropicApiKey"),
  SlackWebhooks: {
    Operations: new sst.Secret("SlackOperationsWebhook"),
  },
};

export const allSecrets = Object.values(secret);
