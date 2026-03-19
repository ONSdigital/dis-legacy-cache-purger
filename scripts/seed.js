// scripts/seed.js
// Usage: mongosh <mongo-uri> scripts/seed.js

//TODO: make these parameters configurable via command line args
const count = 1000;
const pathPrefix = 'test/path';
const collection = 'cachetimes';
const dbName = 'cache'
const publicationCollectionName = 'test-collection'
// Default release time is 1 minute from now, rounded to the nearest minute
const nowPlus1 = new Date(Date.now() + 1 * 60 * 1000);
const releaseTime = new Date(Math.round(nowPlus1.getTime() / 60000) * 60000);

const dbHandle = db.getSiblingDB(dbName);
const docs = [];
for (let i = 0; i < count; i++) {
    docs.push({
        collection_id: publicationCollectionName,
        path: pathPrefix + '/' + i,
        release_time: releaseTime
    });
}
dbHandle[collection].insertMany(docs);
print('Inserted ' + count + ' documents into cache.' + collection);
