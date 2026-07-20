import { Action, Layout, Page, ctx, io } from "@forgeapp/sdk";
import { Project } from "@terminal/core/project/index";

export default new Page({
  name: "Projects",
  handler: async () => {
    return new Layout({
      title: "Projects",
      children: [
        io.display.table("", {
          getData: async (input) => {
            const projects = await Project.list();
            return {
              data: projects.slice(input.offset, input.offset + input.pageSize),
            };
          },
          columns: [
            {
              label: "Name",
              renderCell: (row) => ({
                label: row.name,
                route: "project/edit",
                params: { id: row.id },
              }),
            },
            "year",
            "technologies",
          ],
          rowMenuItems: (row) => [
            {
              label: "Edit",
              route: "project/edit",
              params: { id: row.id },
            },
            {
              label: "Delete",
              route: "project/delete",
              params: { id: row.id },
            },
          ],
          isSortable: false,
        }),
      ],
    });
  },
  routes: {
    new: new Page({
      name: "New Project",
      handler: async () => {
        const name = await io.text("Project Name");
        const year = await io.text("Year (e.g. 2025)");
        const technologies = await io.text("Technologies (comma separated)");
        const description = await io.textarea("Description (use • for bullets)");
        const orderStr = await io.text("Display Order (number, optional)");

        const order = orderStr ? parseInt(orderStr, 10) : undefined;

        await Project.create({
          name,
          year,
          technologies,
          description,
          order: isNaN(order as number) ? undefined : order,
        });

        await ctx.redirect({ route: "project" });
      },
    }),
    edit: new Page({
      name: "Edit Project",
      unlisted: true,
      handler: async () => {
        const id = String(ctx.params.id);
        const project = await Project.fromID(id);

        if (!project) {
          await io.message("Project not found");
          await ctx.redirect({ route: "project" });
          return;
        }

        const name = await io.text("Project Name", project.name);
        const year = await io.text("Year", project.year);
        const technologies = await io.text("Technologies", project.technologies);
        const description = await io.textarea("Description", project.description);
        const orderStr = await io.text("Display Order", String(project.order ?? ""));

        const order = orderStr ? parseInt(orderStr, 10) : undefined;

        await Project.update({
          id,
          name,
          year,
          technologies,
          description,
          order: isNaN(order as number) ? undefined : order,
        });

        await ctx.redirect({ route: "project" });
      },
    }),
    delete: new Action({
      name: "Delete Project",
      unlisted: true,
      handler: async () => {
        const id = String(ctx.params.id);

        const confirmed = await io.confirm(
          "Are you sure you want to delete this project?",
        );

        if (confirmed) {
          await Project.remove(id);
        }

        await ctx.redirect({ route: "project" });
      },
    }),
  },
});
