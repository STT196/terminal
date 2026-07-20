import { text, mysqlTable, varchar, int } from "drizzle-orm/mysql-core";
import { id, timestamps } from "../drizzle/types";

export const projectTable = mysqlTable("project", {
  ...id,
  ...timestamps,
  name: varchar("name", { length: 255 }).notNull(),
  year: varchar("year", { length: 10 }).notNull(),
  technologies: text("technologies").notNull(),
  description: text("description").notNull(),
  order: int("order"),
});
