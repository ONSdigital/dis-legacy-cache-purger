// scripts/seed.js
// Usage: mongosh <mongo-uri> scripts/seed.js
// Parameters injected via --eval by the Makefile seed target

const totalPdfCount = bulletinsCount + articlesCount + compendiumCount;
const dataPaths = ['/datasets/test'];
const normalCount = count - totalPdfCount - dataPaths.length * dataCount - previousCount;
const collection = 'cachetimes';
const dbName = 'cache';
const publicationCollectionID = collectionID;
const publicationCollectionName = collectionID;
// Default release time is 1 minute from now, rounded to the nearest minute
const nowPlus1 = new Date(Date.now() + 1 * 60 * 1000);
const releaseTime = new Date(Math.round(nowPlus1.getTime() / 60000) * 60000);

const dbHandle = db.getSiblingDB(dbName);
const docs = [];
for (let i = 0; i < normalCount; i++) {
    docs.push({
        collection_id: publicationCollectionID,
        collection_title: publicationCollectionName,
        path: path + '/' + i,
        release_time: releaseTime
    });
}
for (const { path, n } of [
    { path: '/bulletins/test',          n: bulletinsCount },
    { path: '/articles/test',           n: articlesCount },
    { path: '/compendium/test', n: compendiumCount },
]) {
    for (let i = 0; i < n; i++) {
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
for (let i = 0; i < previousCount; i++) {
    docs.push({
        collection_id: publicationCollectionID,
        collection_title: publicationCollectionName,
        path: '/datasets/test/' + i + '/current/previous/v' + i,
        release_time: releaseTime
    });
}
const totalCount = normalCount + totalPdfCount + dataPaths.length * dataCount + previousCount;
dbHandle[collection].insertMany(docs);
print('Inserted ' + totalCount + ' documents into cache.' + collection + ' (' + normalCount + ' normal, ' + totalPdfCount + ' pdf-eligible [bulletins=' + bulletinsCount + ' articles=' + articlesCount + ' compendium=' + compendiumCount + '], ' + dataPaths.length * dataCount + ' data-eligible, ' + previousCount + ' previous/excluded)');
