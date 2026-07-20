import { Layout, Page, ctx, io } from "@forgeapp/sdk";
import * as queries from "../queries";

export const UserProfile = new Page({
  name: "User Profile",
  unlisted: true,
  handler: async () => {
    const userID = ctx.params.userID as string | undefined;
    if (!userID) {
      throw new Error("User ID is required. Navigate to this page from the User list.");
    }

    const user = await queries.getUser(userID);
    if (!user) {
      throw new Error("User not found");
    }

    return new Layout({
      title: `User Profile: ${user.name || user.email || userID}`,
      children: [
        io.display.metadata("User Information", {
          layout: "grid",
          data: [
            {
              label: "ID",
              value: user.id,
            },
            {
              label: "Name",
              value: user.name || "N/A",
            },
            {
              label: "Email",
              value: user.email || "N/A",
            },
            {
              label: "Created",
              value: user.timeCreated?.toISOString() || "N/A",
            },
            {
              label: "Updated",
              value: user.timeUpdated?.toISOString() || "N/A",
            },
          ],
        }),
      ],
    });
  },
});
              }),
            },
            "email",
            "created",
            "trackingStatus",
            {
              label: "tracking number",
              renderCell: (row) => ({
                label: row.trackingNumber || "N/A",
              }),
            },
            "fulfiller",
            {
              label: "printed",
              renderCell: (row) => ({
                label: formatters.formatBoolean(!!row.timePrinted),
              }),
            },
          ],
          isSortable: false,
        }),

        // Subscriptions Table
        io.display.heading("Subscriptions", { level: 3 }),
        io.display.table("Subscriptions", {
          getData: async (input) =>
            queries.getUserSubscriptions(userID, {
              offset: input.offset,
              pageSize: input.pageSize,
            }),
          columns: [
            "id",
            "product",
            "productVariant",
            {
              label: "price",
              renderCell: (row) => ({
                label: formatters.formatCurrency(row.price),
              }),
            },
            "quantity",
            {
              label: "schedule",
              renderCell: (row) => ({
                label: formatters.formatSchedule(row.schedule),
              }),
            },
            "next",
            {
              label: "address",
              renderCell: (row) => ({
                label: formatters.formatAddressShort(row.address),
              }),
            },
            "created",
          ],
          isSortable: false,
        }),

        // Shipping Addresses
        io.display.heading("Shipping Addresses", { level: 3 }),
        io.display.table("Addresses", {
          getData: async (input) =>
            queries.getUserAddresses(userID, {
              offset: input.offset,
              pageSize: input.pageSize,
            }),
          columns: [
            {
              label: "id",
              renderCell: (row) => ({
                label: row.id,
                route: "userProfile/editAddress",
                params: {
                  userID,
                  addressID: row.id,
                },
              }),
            },
            {
              label: "address",
              renderCell: (row) => ({
                label: formatters.formatAddressFull(row.address),
                route: "userProfile/editAddress",
                params: {
                  userID,
                  addressID: row.id,
                },
              }),
            },
            "timeCreated",
          ],
          rowMenuItems: (row) => [
            {
              label: "Edit Address",
              route: "userProfile/editAddress",
              params: {
                userID,
                addressID: row.id,
              },
            },
          ],
          isSortable: false,
        }),

        // Payment Methods (Cards)
        io.display.heading("Payment Methods", { level: 3 }),
        io.display.table("Cards", {
          getData: async (input) =>
            queries.getUserCards(userID, {
              offset: input.offset,
              pageSize: input.pageSize,
            }),
          columns: [
            "id",
            "brand",
            {
              label: "last4",
              renderCell: (row) => ({
                label: formatters.formatCardLast4(row.last4),
              }),
            },
            {
              label: "expiration",
              renderCell: (row) => ({
                label: formatters.formatCardExpiration(
                  row.expirationMonth,
                  row.expirationYear,
                ),
              }),
            },
            "timeCreated",
          ],
          isSortable: false,
        }),

        // Cart Info
        io.display.heading("Cart", { level: 3 }),
        io.display.table("Cart", {
          getData: async (input) =>
            queries.getUserCart(userID, {
              offset: input.offset,
              pageSize: input.pageSize,
            }),
          columns: [
            {
              label: "cartID",
              renderCell: (row) => ({
                label: row.cartID || "N/A",
              }),
            },
            {
              label: "items",
              renderCell: (row) => ({
                label: row.items?.toString() || "0",
              }),
            },
            {
              label: "shipping service",
              renderCell: (row) => ({
                label: row.shippingService || "N/A",
              }),
            },
          ],
          isSortable: false,
        }),
      ],
    });
  },
});
