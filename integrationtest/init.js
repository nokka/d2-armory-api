db = db.getSiblingDB('armory');
db.createUser({
    user: "armory",
    pwd: "not_secure_at_all",
    roles: [{
        role: "readWrite",
        db: "armory"
    }]
});

db.createCollection("character");
db.createCollection("statistics");

// Index characters for name in ascending order.
db.character.createIndex({ id: 1 });

// Index statistics for character name in ascending order.
db.statistics.createIndex({ character: 1 });
