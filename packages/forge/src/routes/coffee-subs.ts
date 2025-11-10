import { notLike } from "@terminal/core/drizzle/index";
import { productTable } from "@terminal/core/product/product.sql";
import { createSubscriptionsPage } from "../pages/subscriptions-page";

export const Subs = createSubscriptionsPage({
  name: "Subs: Coffee",
  productFilter: notLike(productTable.name, "cron"),
  routeName: "subsCoffee",
});
