const IPFS = require('ipfs');
const OrbitDB = require('orbit-db');
const NodeRSA = require('node-rsa');
const fs = require('fs').promises;
const yargs = require('yargs');

/*--------------------------------------------------*/
/*                 Sample Usage                     */
/*--------------------------------------------------*/
// console.log('-- Usage: node request.js query -a $address);
// (e.g., -a /orbitdb/zdpuAt9Hn2fGgkhZUhbgg8jbpGsC5JxsdcA7cY1KjuhU6e4aZ/Koma)
// console.log('-- Usage: node request.js decrypt -s $sk -a $address);
// (e.g., -s sk_data.pem -a /orbitdb/zdpuAt9Hn2fGgkhZUhbgg8jbpGsC5JxsdcA7cY1KjuhU6e4aZ/Koma)
/*--------------------------------------------------*/

async function main() {
    // yargs helper
    const argv = yargs
        .command('query', 'Query data of selected orbitdb', {
            address: {
                describe: 'Address of orbitdb',
                alias: 'a',
                demandOption: true,
            }
        })
        .command('decrypt', 'Decrpyt the ciphertext of data', {
            sk: {
                describe: 'Private key',
                alias: 's',
                demandOption: true,
            },
            address: {
                describe: 'Address of orbitdb',
                alias: 'a',
                demandOption: true,
            }
        })
        .help()
        .alias('help', 'h').argv;

    // Create IPFS instance
    var db;
    const ipfsOptions = { repo: './ipfs', };
    const ipfs = await IPFS.create(ipfsOptions);

    // Create OrbitDB instance and open existed db
    const orbitdb = await OrbitDB.createInstance(ipfs);
    db = await orbitdb.open(argv.address);
    await db.load();

    if (argv._[0] == "query") {
        getAll(db, orbitdb);
        process.exit(0);

    } else if (argv._[0] == "decrypt") {
        const sk_data = await fs.readFile(argv.sk);
        const key = new NodeRSA(sk_data);

        ciphertext = await getAll(db, orbitdb);
        console.log('\n\nCiphertext: ')
        ciphertext.forEach(element => {
            console.log('\n==========================================')
            console.log('_id: ', element._id);
            console.log('desc: ', key.decrypt(element.desc, 'utf8'));
            console.log('name: ', key.decrypt(element.name, 'utf8'));
            console.log('role: ', element.role);
            console.log('price: ', key.decrypt(element.price, 'utf8'));
            console.log('artwork: ', key.decrypt(element.artwork,));
        });
        process.exit(0);
    }
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
    return allArtwork;
}

main();
