const IPFS = require('ipfs');
const OrbitDB = require('orbit-db');
const { v4: uuidv4 } = require('uuid');
const fs = require('fs').promises;

async function main() {

  // Create OrbitDB or put data to existed DB
  const args = process.argv;
  if (args.length != 4) {
    console.log('-- Usage: node setup.js create $username');
    console.log('-- Usage: node setup.js upload $data\n');

    console.log('Please give the correct argument :)\n');
    process.exit(1);
  }

  let func = args[2];
  let jsonData = args[3];
  var db;

  // Create IPFS instance
  const ipfsOptions = { repo: './ipfs', };
  const ipfs = await IPFS.create(ipfsOptions);

  // Create OrbitDB instance
  const orbitdb = await OrbitDB.createInstance(ipfs);

  if (func == "create") {
    // Create a new docs db
    db = await orbitdb.docs(jsonData, { indexBy: 'name' });
    console.log('\n-- ' + jsonData + ' created --');
    console.log('-- DB address: ' + db.address.toString());

  } else if (func == "upload") {
    // Open existed DB and put data
    db = await orbitdb.open('/orbitdb/zdpuAqUjn1CBkRJr2ycTSd98HQVG8weKgJTdc1z1f6d7agpWZ/komadb');
    await db.load();

    // Parse JSON fortmat artwork data to js object
    var artData = JSON.parse(jsonData);

    //  Use 'fs' to read the image and save its string into db
    try {
      const image = await fs.readFile(artData.imagePath);
      var hash = await db.put({ _id: uuidv4(), artwork: image, name: artData.name, desc: artData.desc, price: artData.price });
      console.log('\n-- Store ' + hash + ' successfully \n');
    } catch (error) {
      console.log(error);
    }
  }

  // Check all the artwork
  const allArtwork = db.get('');

  console.log('--List all artwork'); 
  console.log(allArtwork);

  // Disconnect with db and exit successfully
  await db.close();
  await orbitdb.disconnect();
  process.exit(0);
}

main();
