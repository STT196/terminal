import { Layout, Page, io } from "@forgeapp/sdk";
import * as queries from "../queries";

export const User = new Page({
  name: "User",
  handler: async (c) => {
    return new Layout({
      title: "User",
      children: [
        io.display.table("", {
          getData: async (input) => {
            const queryTerm = input.queryTerm?.trim() || undefined;
            return queries.getAllUsers(
              {
                offset: input.offset,
                pageSize: input.pageSize,
              },
              queryTerm,
            );
          },
          columns: [
            {
              label: "id",
              renderCell: (row) => ({
                label: row.id,
                route: "userProfile",
                params: {
                  userID: row.id,
                },
              }),
            },
            "name",
            "email",
            "timeCreated",
          ],
          isSortable: false,
        }),
      ],
    });
  },
});
          throw new Error("Selected address not found");
        }

        // Get products and variants
        const products = await Product.list();
        const results = await io.group(
          products.flatMap((product) => {
            return product.variants.map((variant) =>
              io.input.number(`${product.name} - ${variant.name}`, {
                defaultValue: 0,
              }),
            );
          }),
        );

        // Build items object
        const items = {} as Record<string, number>;
        let variantIndex = 0;
        for (const product of products) {
          for (const variant of product.variants) {
            const quantity = results[variantIndex] as number;
            if (quantity > 0) {
              items[variant.id] = quantity;
            }
            variantIndex++;
          }
        }

        if (Object.keys(items).length === 0) {
          throw new Error("No items selected");
        }

        // Create the order
        await OrderM.createInternal({
          email: user.email,
          items,
          address: addressData,
        });

        await ctx.redirect({
          route: "user",
        });
      },
    }),
  },
});
