const IPFS = require('ipfs');
const OrbitDB = require('orbit-db');
const Identities = require('orbit-db-identity-provider')
const { v4: uuidv4 } = require('uuid');
const NodeRSA = require('node-rsa');
const fs = require('fs').promises;
const yargs = require('yargs');

/*--------------------------------------------------*/
/*                 Sample Usage                     */
/*--------------------------------------------------*/

// console.log('-- Usage: node setup.js create -n $dbname (e.g., -n komadb)');
// console.log('-- Usage: node setup.js upload -a $artwork -d $dbname (e.g., -a '{ "imagePath": "amber-mosaic.jpg", "name": "amber", "desc": "generated by GAN", "price": "10000" }' -d komadb)');
// console.log('-- Usage: node setup.js transfer -t $tokenID -n $newCollector -o $owner (e.g., -t XXXX-XXXX-XXXX-XXXX -nc newCollector -o owner)');

/*--------------------------------------------------*/

async function main() {
  // yargs helper
  const argv = yargs
    .command('create', 'Create a new orbitdb', {
      name: {
        describe: 'Name of orbitdb',
        alias: 'n',
        demandOption: true,
      }
    })
    .command('upload', 'Upload artworks to orbitdb', {
      artwork: {
        describe: 'Payload of artwork',
        alias: 'a',
        demandOption: true,
      },
      dbName: {
        describe: 'Name of uploaded orbitdb',
        alias: 'd',
        demandOption: true,
      }
    })
    .command('transfer', 'Trasfer artwork', {
      tokenID: {
        describe: 'UUID of artwork',
        alias: 't',
        demandOption: true,
      },
      newCollector: {
        describe: 'New Collector',
        alias: 'n',
        demandOption: true,
      },
      owner: {
        describe: 'Owner',
        alias: 'o',
        demandOption: true,
      }
    })
    .help()
    .alias('help', 'h').argv;

  var db;

  // Create IPFS instance with DS_identity
  const ipfsOptions = { repo: './ipfs', };
  const ipfs = await IPFS.create(ipfsOptions);
  const options = { id: 'DS-id' }
  const DS_identity = await Identities.createIdentity(options)

  // Create OrbitDB instance
  const orbitdb = await OrbitDB.createInstance(ipfs);

  if (argv._[0] == "create") {
    // Create a new docs db
    db = await orbitdb.docs(argv.name, {
      accessController: {
        write: ['*']
      }
    });

    
    // Write hash address to file
    try {
      const image = await fs.appendFile('./address.txt', argv.name + ' ' + db.address.toString() + ' ');

      console.log('\n== ' + argv.name + ' created ==');
	  console.log('-- Identity: ' + db.identity.publicKey);
      console.log('-- DB address: ' + db.address.toString());
      process.exit(0);
    } catch (error) {
      console.log(error);
    }

  } else if (argv._[0] == "upload") {
    try {
      // Open existed db
      db = await orbitdb.open(await queryAddress(argv.dbName));
      await db.load();

      // Parse JSON fortmat artwork data to js object
      var artData = JSON.parse(argv.artwork);

      // Read the image and save its string into db, firstly, we encrpyted ciphertext imformation using pk_data generated from "Upload_Artworks.py"
      const pk_data = await fs.readFile('pk_data.pem');
      const key = new NodeRSA(pk_data);

      const image = await fs.readFile(artData.imagePath);
      const tokenID = uuidv4();
      const hash = await db.put({ _id: tokenID, artwork: key.encrypt(image), name: key.encrypt(artData.name), desc: key.encrypt(artData.desc), price: key.encrypt(artData.price), role: 'creator' }); // '+' means current collector

      console.log('\n-- Store Successfully \n');
      console.log('tokenID:' + tokenID);
      console.log('multi-hash:' + hash);

      // Check and close db
      //getAll(db, orbitdb);
      process.exit(0);

    } catch (error) {
      console.log(error);
    }
  } else if (argv._[0] == "transfer") {
    try {
      // Now we want to transfer artwork from "owner" to "newCollector"
      // Open owner's and newCollector's db
      dbOwner = await orbitdb.open(await queryAddress(argv.owner));
      await dbOwner.load();
      dbNew = await orbitdb.open(await queryAddress(argv.newCollector));
      await dbNew.load();

      // Get trasferring artwork information from owner's db, note that if it's the one-hand transaction, no need to change creator's role
      const transferArtwork = dbOwner.get(argv.tokenID);

      // Write artwork information to newCollector's db (+)
      const hash = await dbNew.put({ _id: transferArtwork[0]._id, artwork: transferArtwork[0].artwork, name: transferArtwork[0].name, desc: transferArtwork[0].desc, price: transferArtwork[0].price, role: '+collector' });
      console.log('\n-- Store ' + hash + ' successfully \n');

      // Change artwork information from owner's db (-)
      if (transferArtwork[0].role != 'creator') {
        const hash2 = await dbOwner.put({ _id: transferArtwork[0]._id, artwork: transferArtwork[0].artwork, name: transferArtwork[0].name, desc: transferArtwork[0].desc, price: transferArtwork[0].price, role: '-collector' });
        console.log('\n-- Store ' + hash2 + ' successfully \n');
      }

      // Check and close db
      await dbOwner.close();
      getAll(dbNew, orbitdb);
	  process.exit(0);

    } catch (error) {
      console.log(error);
    }
  }
}

/* Query hash address by name */
async function queryAddress(name) {
  // Read hash address of username from file
  const allAddress = await fs.readFile('./address.txt', 'utf8');
  const addressArray = allAddress.split(' ');
  const isAddress = (element) => element == name;

  return addressArray[addressArray.findIndex(isAddress) + 1];
}

/* Get all artwork and close db */
async function getAll(db, orbitdb) {
  // Check all the artwork
  const allArtwork = db.get('');

  console.log('-- List all artwork');
  console.log(allArtwork);

  // Disconnect with db and exit successfully
  await db.close();
  await orbitdb.disconnect();
}

main();
