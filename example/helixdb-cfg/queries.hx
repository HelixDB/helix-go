QUERY create_user(name: String, age: U32, email: String, now: I32) =>
    user <- AddN<User>({name: name, age: age, email: email, created_at: now, updated_at: now})
    RETURN user 

QUERY create_users(users: [{name: String, age: U32, email: String, now: I32}]) =>
    FOR { name, age, email, now } IN users {
        AddN<User>({name: name, age: age, email: email, created_at: now, updated_at: now})
    }
    RETURN "Success" 

QUERY update_user(id: ID, name: String, age: U32, email: String) =>
    updated_user <- N<User>(id)::UPDATE({name: name, age: age, email: email})
    RETURN updated_user 

QUERY get_users() =>
    users <- N<User>
    RETURN users 

QUERY get_user_by_id(id: ID) =>
    user <- N<User>(id)
    RETURN user 

QUERY follow(followerId: ID, followedId: ID) =>
    follower <- N<User>(followerId)
    followed <- N<User>(followedId)
    AddE<Follows>::From(follower)::To(followed)
    RETURN "Success" 

QUERY followers(id: ID) =>
    followers <- N<User>(id)::In<Follows>
    RETURN followers 

QUERY follower_count(id: ID) =>
    count <- N<User>(id)::In<Follows>::COUNT
    RETURN count 

QUERY following(id: ID) =>
    following <- N<User>(id)::Out<Follows>
    RETURN following 

QUERY following_count(id: ID) =>
    count <- N<User>(id)::Out<Follows>::COUNT
    RETURN count 
