// scripts/seed.js
// Usage: mongosh <mongo-uri> scripts/seed.js

const pdfPaths = ['/bulletins/test', '/articles/test', '/compendium_chapter/test'];
const dataPaths = ['/datasets/test'];
const normalCount = count - pdfPaths.length * pdfCount - dataPaths.length * dataCount;
const collection = 'cachetimes';
const dbName = 'cache'
const publicationCollectionID = 'test-collection'
const publicationCollectionName = 'Test Collection'
// Default release time is 1 minute from now, rounded to the nearest minute
const nowPlus1 = new Date(Date.now() + 1 * 60 * 1000);
const releaseTime = new Date(Math.round(nowPlus1.getTime() / 60000) * 60000);

const dbHandle = db.getSiblingDB(dbName);
const docs = [];
for (let i = 0; i < count; i++) {
    docs.push({
        collection_id: publicationCollectionID,
        collection_title: publicationCollectionName,
        path: path + '/' + i,
        release_time: releaseTime
    });
}
for (const path of pdfPaths) {
    for (let i = 0; i < pdfCount; i++) {
        docs.push({
            collection_id: publicationCollectionID,
            collection_title: publicationCollectionName,
            path: path + '/' + i,
            release_time: releaseTime
        });
    }
}
for (const path of dataPaths) {
    for (let i = 0; i < dataCount; i++) {
        docs.push({
            collection_id: publicationCollectionID,
            collection_title: publicationCollectionName,
            path: path + '/' + i,
            release_time: releaseTime
        });
    }
}
const totalCount = normalCount + pdfPaths.length * pdfCount + dataPaths.length * dataCount;
dbHandle[collection].insertMany(docs);
print('Inserted ' + totalCount + ' documents into cache.' + collection + ' (' + normalCount + ' normal, ' + pdfPaths.length * pdfCount + ' pdf-eligible, ' + dataPaths.length * dataCount + ' data-eligible)');
