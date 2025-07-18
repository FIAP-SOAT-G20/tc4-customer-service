// MongoDB initialization script
// This script runs when the MongoDB container starts for the first time

// Switch to the application database
db = db.getSiblingDB("fastfood_10soat_g22_tc4");

// Create application user with read/write permissions
db.createUser({
    user: "app_user",
    pwd: "app_password",
    roles: [
        {
            role: "readWrite",
            db: "fastfood_10soat_g22_tc4",
        },
    ],
});

// Create customers collection with schema validation
db.createCollection("customers", {
    validator: {
        $jsonSchema: {
            bsonType: "object",
            required: ["name", "email", "cpf", "created_at", "updated_at"],
            properties: {
                _id: {
                    bsonType: "long",
                    description: "must be a long and is required",
                },
                name: {
                    bsonType: "string",
                    description: "must be a string and is required",
                },
                email: {
                    bsonType: "string",
                    pattern: "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$",
                    description: "must be a valid email address and is required",
                },
                cpf: {
                    bsonType: "string",
                    pattern: "^[0-9]{11}$|^[0-9]{3}\\.[0-9]{3}\\.[0-9]{3}-[0-9]{2}$",
                    description: "must be a valid CPF format and is required",
                },
                created_at: {
                    bsonType: "date",
                    description: "must be a date and is required",
                },
                updated_at: {
                    bsonType: "date",
                    description: "must be a date and is required",
                },
            },
        },
    },
});

// Create indexes for better performance
db.customers.createIndex({cpf: 1}, {unique: true});
db.customers.createIndex({email: 1}, {unique: true});

print("MongoDB initialization completed successfully");
