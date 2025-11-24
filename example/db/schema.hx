N::User {
    name: String,
    age: U32,
    email: String,
    created_at: I32,
}

E::Follows {
    From: User,
    To: User,
    Properties: {
        since: I32,
    }
}

V::Preference {
    preference: String,
}

E::UserPreference {
    From: User,
    To: Preference,
}
