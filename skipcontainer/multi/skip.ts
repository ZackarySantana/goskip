import {
    type EagerCollection,
    type Json,
    type Resource,
    OneToManyMapper,
} from "@skipruntime/api";

import {
    initialData,
    type Group,
    type GroupID,
    type User,
    type UserID,
} from "./subdir/helpers";

import { runService } from "@skipruntime/server";

// Type alias for inputs to our service
export type ServiceInputs = {
    users: EagerCollection<UserID, User>;
    groups: EagerCollection<GroupID, Group>;
};

// Type alias for inputs to the active friends resource
type ResourceInputs = {
    users: EagerCollection<UserID, User>;
    actives: EagerCollection<GroupID, UserID>;
};

// Mapper function to compute the active users of each group
class ActiveUsers extends OneToManyMapper<GroupID, Group, UserID> {
    constructor(private users: EagerCollection<UserID, User>) {
        super();
    }

    mapValue(group: Group): UserID[] {
        return group.members.filter((uid) => this.users.getUnique(uid).active);
    }
}

// Mapper function to filter out those active users who are also friends with `user`
class FilterFriends extends OneToManyMapper<GroupID, UserID, UserID> {
    constructor(private readonly user: User) {
        super();
    }

    mapValue(uid: UserID): UserID[] {
        return this.user.friends.includes(uid) ? [uid] : [];
    }
}

class ActiveFriends implements Resource<ResourceInputs> {
    private readonly uid: UserID;

    constructor(params: Json) {
        if (typeof params != "number")
            throw new Error("Missing required number parameter 'uid'");
        this.uid = params;
    }

    instantiate(inputs: ResourceInputs): EagerCollection<GroupID, UserID> {
        const user = inputs.users.getUnique(this.uid);
        return inputs.actives.map(FilterFriends, user);
    }
}

// Specify and run the reactive service
await runService(
    {
        initialData,
        resources: { active_friends: ActiveFriends },
        createGraph(input: ServiceInputs): ResourceInputs {
            const actives = input.groups.map(ActiveUsers, input.users);
            return { users: input.users, actives };
        },
    },
    { streaming_port: 8080, control_port: 8081 }
);
