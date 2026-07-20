import { z } from "zod";
import "zod-openapi/extend";
import { Result, ErrorResponses, validator, authRequired } from "./common";
import { Project } from "@terminal/core/project/index";
import { Hono } from "hono";
import { describeRoute } from "hono-openapi";
import { Examples } from "@terminal/core/examples";

export namespace ProjectApi {
  export const route = new Hono()
    .get(
      "/",
      describeRoute({
        tags: ["Project"],
        summary: "List projects",
        description: "List all projects in the portfolio.",
        security: [],
        responses: {
          200: {
            content: {
              "application/json": {
                schema: Result(
                  Project.Info.array().openapi({
                    description: "A list of projects.",
                    example: [Examples.Project],
                  }),
                ),
                example: { data: [Examples.Project] },
              },
            },
            description: "A list of projects.",
          },
          429: ErrorResponses[429],
          500: ErrorResponses[500],
        },
      }),
      async (c) => {
        return c.json(
          {
            data: await Project.list(),
          },
          200,
        );
      },
    )
    .get(
      "/:id",
      describeRoute({
        tags: ["Project"],
        summary: "Get project",
        description: "Get a project by ID from the portfolio.",
        security: [],
        responses: {
          200: {
            content: {
              "application/json": {
                schema: Result(
                  Project.Info.openapi({
                    description: "The requested project.",
                    example: Examples.Project,
                  }),
                ),
                example: { data: Examples.Project },
              },
            },
            description: "The requested project.",
          },
          404: ErrorResponses[404],
          429: ErrorResponses[429],
          500: ErrorResponses[500],
        },
      }),
      validator(
        "param",
        z.object({
          id: Project.Info.shape.id.openapi({
            description: "ID of the project to get.",
            example: Examples.Project.id,
          }),
        }),
      ),
      async (c) => {
        const project = await Project.fromID(c.req.valid("param").id);
        if (!project) {
          return c.json(
            {
              type: "not_found",
              code: "resource_not_found",
              message: "Project not found",
            },
            404,
          );
        }
        return c.json({ data: project }, 200);
      },
    )
    .post(
      "/",
      describeRoute({
        tags: ["Project"],
        summary: "Create project",
        description: "Create a new project in the portfolio.",
        responses: {
          201: {
            content: {
              "application/json": {
                schema: Result(z.string().openapi({
                  description: "The ID of the created project.",
                  example: Examples.Project.id,
                })),
              },
            },
            description: "Project created successfully.",
          },
          400: ErrorResponses[400],
          429: ErrorResponses[429],
          500: ErrorResponses[500],
        },
      }),
      authRequired,
      validator("json", Project.CreateInput),
      async (c) => {
        const id = await Project.create(c.req.valid("json"));
        return c.json({ data: id }, 201);
      },
    )
    .patch(
      "/:id",
      describeRoute({
        tags: ["Project"],
        summary: "Update project",
        description: "Update an existing project in the portfolio.",
        responses: {
          200: {
            content: {
              "application/json": {
                schema: Result(z.literal("ok").openapi({
                  example: "ok",
                })),
              },
            },
            description: "Project updated successfully.",
          },
          400: ErrorResponses[400],
          404: ErrorResponses[404],
          429: ErrorResponses[429],
          500: ErrorResponses[500],
        },
      }),
      authRequired,
      validator(
        "param",
        z.object({
          id: Project.Info.shape.id.openapi({
            description: "ID of the project to update.",
            example: Examples.Project.id,
          }),
        }),
      ),
      validator("json", Project.UpdateInput),
      async (c) => {
        await Project.update({
          id: c.req.valid("param").id,
          ...c.req.valid("json"),
        });
        return c.json({ data: "ok" as const }, 200);
      },
    )
    .delete(
      "/:id",
      describeRoute({
        tags: ["Project"],
        summary: "Delete project",
        description: "Delete a project from the portfolio (soft delete).",
        responses: {
          200: {
            content: {
              "application/json": {
                schema: Result(z.literal("ok").openapi({
                  example: "ok",
                })),
              },
            },
            description: "Project deleted successfully.",
          },
          404: ErrorResponses[404],
          429: ErrorResponses[429],
          500: ErrorResponses[500],
        },
      }),
      authRequired,
      validator(
        "param",
        z.object({
          id: Project.Info.shape.id.openapi({
            description: "ID of the project to delete.",
            example: Examples.Project.id,
          }),
        }),
      ),
      async (c) => {
        await Project.remove(c.req.valid("param").id);
        return c.json({ data: "ok" as const }, 200);
      },
    );
}
