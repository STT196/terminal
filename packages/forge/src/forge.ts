import { Forge } from "@forgeapp/sdk";
import { Resource } from "sst";

import { User } from "./routes/user";
import { UserProfile } from "./routes/user-profile";
import ProjectPage from "./routes/project";

const forge = new Forge({
  apiKey: Resource.ForgeKey.value,
  endpoint: "wss://terminal.app.forgeapp.io/websocket",
  routes: {
    user: User,
    userProfile: UserProfile,
    project: ProjectPage,
  },
});

forge.listen();
