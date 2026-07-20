import { database } from "./database";
import { allSecrets } from "./secret";

export const bus = new sst.aws.Bus("Bus");

bus.subscribe("Event", {
  handler: "./packages/functions/src/event/event.handler",
  link: [database, ...allSecrets],
  timeout: "5 minutes",
  permissions: [
    {
      actions: ["ses:SendEmail"],
      resources: ["*"],
    },
  ],
});
