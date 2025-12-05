import { z } from "zod";
import { Common } from "../common";
import { Examples } from "../examples";

/**
 * Address schema definitions.
 * Separated from index.ts to avoid circular dependencies with Order/Shippo.
 */

export const AddressInner = z
  .object({
    name: z.string().openapi({
      description: "The recipient's name.",
      example: Examples.Shipping.name,
    }),
    street1: z.string().openapi({
      description: "Street of the address.",
      example: Examples.Shipping.street1,
    }),
    street2: z.string().optional().openapi({
      description: "Apartment, suite, etc. of the address.",
      example: Examples.Shipping.street2,
    }),
    city: z.string().openapi({
      description: "City of the address.",
      example: Examples.Shipping.city,
    }),
    province: z.string().optional().openapi({
      description: "Province or state of the address.",
      example: Examples.Shipping.province,
    }),
    country: z
      .string()
      .length(2, "Country must be a 2 character country code (ISO 3166-1)")
      .openapi({
        description: "ISO 3166-1 alpha-2 country code of the address.",
        example: Examples.Shipping.country,
      }),
    zip: z.string().openapi({
      description: "Zip code of the address.",
      example: Examples.Shipping.zip,
    }),
    phone: z.string().optional().openapi({
      description: "Phone number of the recipient.",
      example: Examples.Shipping.phone,
    }),
  })
  .openapi({
    description: "Address information.",
    example: Examples.Address,
  });

export type AddressInner = z.infer<typeof AddressInner>;

export const AddressInfo = z
  .object({
    id: z.string().openapi({
      description: Common.IdDescription,
      example: Examples.Shipping.id,
    }),
    ...AddressInner.shape,
    created: z.coerce.date().openapi({
      description: "Date the address was created.",
      example: Examples.Shipping.created,
    }),
  })
  .openapi({
    ref: "Address",
    description: "Physical address associated with a Terminal shop user.",
    example: Examples.Shipping,
  });

export type AddressInfo = z.infer<typeof AddressInfo>;
