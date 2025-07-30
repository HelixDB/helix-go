QUERY create_user(name: String, age: U32, email: String, now: I32) =>
    newUser <- AddN<User>({name: name, age: age, email: email, created_at: now, updated_at: now})
    RETURN newUser 

QUERY create_users(users: [{name: String, age: U32, email: String, now: I32}]) =>
    FOR { name, age, email, now } IN users {
        AddN<User>({name: name, age: age, email: email, created_at: now, updated_at: now})
    }
    RETURN "Success" 

QUERY update_user(id: ID, name: String, age: U32, email: String) =>
    updatedUser <- N<User>(id)::UPDATE({name: name, age: age, email: email})
    RETURN updatedUser 

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

QUERY following(id: ID) =>
    following <- N<User>(id)::Out<Follows>
    RETURN following 
