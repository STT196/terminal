import { z } from "zod";
import { and, eq } from "drizzle-orm";
import { useTransaction } from "../drizzle/transaction";
import { fn } from "../util/fn";
import { createID } from "../util/id";
import { addressTable } from "./address.sql";
import { Shippo } from "../shippo";
import { cartTable } from "../cart/cart.sql";
import { subscriptionTable } from "../subscription/subscription.sql";
import { VisibleError, ErrorCodes } from "../error";
import { Actor } from "../actor";
import {
  AddressInner,
  AddressInfo,
  type AddressInner as AddressInnerType,
  type AddressInfo as AddressInfoType,
} from "./schema";

export namespace Address {
  /** @see AddressInner */
  export const Inner = AddressInner;
  export type Inner = AddressInnerType;

  /** @see AddressInfo */
  export const Info = AddressInfo;
  export type Info = AddressInfoType;

  export function list() {
    return useTransaction(async (tx) =>
      tx
        .select()
        .from(addressTable)
        .where(eq(addressTable.userID, Actor.userID()))
        .execute()
        .then((rows) => rows.map(serialize)),
    );
  }

  export const create = fn(Inner, (input) =>
    useTransaction(async (tx) => {
      const validated = await Shippo.assertValidAddress(input);
      const id = createID("userShipping");
      await tx.insert(addressTable).values({
        id,
        userID: Actor.userID(),
        address: validated,
      });
      return id;
    }),
  );

  export const remove = fn(z.string(), (input) =>
    useTransaction(async (tx) => {
      const subscriptions = await tx
        .select()
        .from(subscriptionTable)
        .where(eq(subscriptionTable.addressID, input));
      if (subscriptions.length > 0) {
        throw new VisibleError(
          "validation",
          ErrorCodes.Validation.IN_USE,
          "Address is in use by a subscription, please cancel the subscription first.",
          "id",
        );
      }

      await tx
        .update(cartTable)
        .set({ addressID: null })
        .where(eq(cartTable.addressID, input));
      const response = await tx
        .delete(addressTable)
        .where(
          and(
            eq(addressTable.id, input),
            eq(addressTable.userID, Actor.userID()),
          ),
        );
      if (response.rowsAffected === 0) {
        throw new VisibleError(
          "not_found",
          ErrorCodes.NotFound.RESOURCE_NOT_FOUND,
          "Address not found",
        );
      }
    }),
  );

  export const fromID = fn(Info.shape.id, (id) =>
    useTransaction(async (tx) => {
      const rows = await tx
        .select()
        .from(addressTable)
        .where(
          and(eq(addressTable.id, id), eq(addressTable.userID, Actor.userID())),
        )
        .limit(1);
      return rows.map(serialize).at(0);
    }),
  );

  function serialize(
    input: typeof addressTable.$inferSelect,
  ): z.infer<typeof Info> {
    return {
      id: input.id,
      name: input.address.name,
      street1: input.address.street1,
      street2: input.address.street2,
      city: input.address.city,
      province: input.address.province,
      country: input.address.country,
      zip: input.address.zip,
      phone: input.address.phone,
      created: input.timeCreated,
    };
  }
}
