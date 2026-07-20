import { z } from "zod";
import "zod-openapi/extend";
import { useTransaction } from "../drizzle/transaction";
import { projectTable } from "./project.sql";
import { eq, isNull, asc } from "drizzle-orm";
import { fn } from "../util/fn";
import { createID } from "../util/id";
import { Common } from "../common";
import { Examples } from "../examples";

export namespace Project {
  export const Info = z
    .object({
      id: z.string().openapi({
        description: Common.IdDescription,
        example: Examples.Project.id,
      }),
      name: z.string().openapi({
        description: "Name of the project.",
        example: Examples.Project.name,
      }),
      year: z.string().openapi({
        description: "Year of the project.",
        example: Examples.Project.year,
      }),
      technologies: z.string().openapi({
        description: "Technologies used in the project.",
        example: Examples.Project.technologies,
      }),
      description: z.string().openapi({
        description: "Description of the project. Supports bullet points with \\n separator.",
        example: Examples.Project.description,
      }),
      order: z.number().int().optional().openapi({
        description: "Order of the project used when displaying a sorted list.",
        example: Examples.Project.order,
      }),
    })
    .openapi({
      ref: "Project",
      description: "A project in the portfolio.",
      example: Examples.Project,
    });

  export type Info = z.infer<typeof Info>;

  export const CreateInput = z
    .object({
      name: z.string().min(1),
      year: z.string().min(1),
      technologies: z.string().min(1),
      description: z.string().min(1),
      order: z.number().int().optional(),
    })
    .openapi({
      ref: "ProjectCreate",
      description: "Input for creating a project.",
    });

  export const UpdateInput = z
    .object({
      name: z.string().optional(),
      year: z.string().optional(),
      technologies: z.string().optional(),
      description: z.string().optional(),
      order: z.number().int().optional(),
    })
    .openapi({
      ref: "ProjectUpdate",
      description: "Input for updating a project.",
    });

  export const list = () =>
    useTransaction(async (tx) => {
      const rows = await tx
        .select()
        .from(projectTable)
        .where(isNull(projectTable.timeDeleted))
        .orderBy(asc(projectTable.order));

      return rows.map((row): Info => ({
        id: row.id,
        name: row.name,
        year: row.year,
        technologies: row.technologies,
        description: row.description,
        order: row.order || undefined,
      }));
    });

  export const fromID = fn(Info.shape.id, (input) =>
    useTransaction(async (tx) => {
      const rows = await tx
        .select()
        .from(projectTable)
        .where(eq(projectTable.id, input))
        .where(isNull(projectTable.timeDeleted))
        .limit(1);

      const row = rows[0];
      if (!row) return undefined;

      return {
        id: row.id,
        name: row.name,
        year: row.year,
        technologies: row.technologies,
        description: row.description,
        order: row.order || undefined,
      } satisfies Info;
    }),
  );

  export const create = fn(CreateInput, (input) =>
    useTransaction(async (tx) => {
      const id = createID("project");
      await tx.insert(projectTable).values({
        id,
        name: input.name,
        year: input.year,
        technologies: input.technologies,
        description: input.description,
        order: input.order,
      });
      return id;
    }),
  );

  export const update = fn(
    z.object({ id: Info.shape.id, ...UpdateInput.shape }),
    (input) =>
      useTransaction(async (tx) => {
        const { id, ...data } = input;
        await tx
          .update(projectTable)
          .set(data)
          .where(eq(projectTable.id, id));
      }),
  );

  export const remove = fn(Info.shape.id, (input) =>
    useTransaction(async (tx) => {
      await tx
        .update(projectTable)
        .set({ timeDeleted: new Date() })
        .where(eq(projectTable.id, input));
    }),
  );
}
