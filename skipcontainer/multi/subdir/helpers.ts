import { type InitialData } from "@skipruntime/api";

import { ServiceInputs } from "../skip";

export type UserID = number;
export type GroupID = number;
export type User = { name: string; active?: boolean; friends: UserID[] };
export type Group = { name: string; members: UserID[] };

// Load initial data from a source-of-truth database (mocked for simplicity)
export const initialData: InitialData<ServiceInputs> = {
    users: [
        [0, [{ name: "Bob", active: true, friends: [1, 2] }]],
        [1, [{ name: "Alice", active: true, friends: [0, 2] }]],
        [2, [{ name: "Carol", active: false, friends: [0, 1] }]],
        [3, [{ name: "Eve", active: true, friends: [] }]],
    ],
    groups: [
        [1001, [{ name: "Group 1", members: [1, 2, 3] }]],
        [1002, [{ name: "Group 2", members: [0, 2] }]],
    ],
};
