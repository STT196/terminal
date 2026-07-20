import { useTransaction } from "@terminal/core/drizzle/transaction";
import { desc, isNotNull, and, SQL, like, or } from "@terminal/core/drizzle/index";
import { userTable } from "@terminal/core/user/user.sql";

type PaginationInput = {
  offset: number;
  pageSize: number;
};

/**
 * Get user by ID
 */
export async function getUser(userID: string) {
  return useTransaction(async (tx) =>
    tx
      .select()
      .from(userTable)
      .where(eq(userTable.id, userID))
      .then((rows) => rows[0]),
  );
}

/**
 * Get all users with optional search filter by name or email
 * @param requireNameAndEmail - If true, filters out users who don't have both a name and email (default: true)
 */
export async function getAllUsers(
  pagination: PaginationInput,
  queryTerm?: string,
  requireNameAndEmail: boolean = true,
) {
  const whereConditions: SQL[] = [];
  
  if (requireNameAndEmail) {
    whereConditions.push(isNotNull(userTable.name));
    whereConditions.push(isNotNull(userTable.email));
  }
  
  if (queryTerm) {
    whereConditions.push(
      or(
        like(userTable.id, "%" + queryTerm + "%"),
        like(userTable.name, "%" + queryTerm + "%"),
        like(userTable.email, "%" + queryTerm + "%"),
      )!,
    );
  }

  return useTransaction(async (tx) => ({
    data: await tx
      .select()
      .from(userTable)
      .where(whereConditions.length > 0 ? and(...whereConditions) : sql`true`)
      .orderBy(desc(userTable.id))
      .offset(pagination.offset)
      .limit(pagination.pageSize),
  }));
}

