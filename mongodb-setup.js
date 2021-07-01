db.createUser({
    user: "squirrel",
    pwd: "squirrel",
    roles: [
        {
            role: "readWrite", db: "squirrel"
        }
    ]
})