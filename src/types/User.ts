import type { Base, Timestamp } from "./Common";

export interface User extends Base, Timestamp {
  userName: string;
  altName?: string;

  email: string;
  altEmail?: string;

  password: string;
}
