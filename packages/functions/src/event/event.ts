import "zod-openapi/extend";
import { bus } from "sst/aws/bus";
import { User } from "@terminal/core/user/index";
import { Log } from "@terminal/core/util/log";
import { Actor } from "@terminal/core/actor";

const log = Log.create({ namespace: "event" });
export const handler = bus.subscriber(
  [User.Event.Updated],
  async (event) =>
    Actor.provide(
      event.metadata.actor.type,
      event.metadata.actor.properties,
      async () => {
        log.info("received", {
          type: event.type,
          ...event.properties,
        });
      },
    ),
);
            // Optionally add to a subscribers list if needed
            // await EmailOctopus.addToSubscribersList(event.properties.subscriptionID);
            break;
          }
        }
      },
    ),
);
