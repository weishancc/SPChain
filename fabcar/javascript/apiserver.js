var bodyParser = require('body-parser');
var express = require('express');
var app = express();
app.use(bodyParser.json());

// Setting for Hyperledger Fabric
const { FileSystemWallet, Gateway } = require('fabric-network');
const path = require('path');
const ccpPath = path.resolve(__dirname, '..', '..', 'first-network', 'connection-gallery.json')

app.post('/spchain/grantConsent/', async function (req, res) {
    try {
		console.log('Received: ');
		console.log('\npk_DS: ' + req.body.pk_DS);
		console.log('\npk_DC: ' + req.body.pk_DC);
		console.log('\nPolicy: ' + req.body.policy);
		console.log('\nEnhash: ' + req.body.enhash);
		console.log('\npk_enc: ' + req.body.pk_enc);

        // Create a new file system based wallet for managing identities.
        const walletPath = path.join(process.cwd(), 'wallet');
        const wallet = new FileSystemWallet(walletPath);
        console.log(`\nWallet path: ${walletPath}`);

        // Check to see if we've already enrolled the user.
        const userExists = await wallet.exists('Koma');
        if (!userExists) {
            console.log('An identity for the user "Koma" does not exist in the wallet');
            console.log('Run the registerGalleryUser.js application before retrying');
			require('child_process').fork('registerGalleryUser.js');
            return;
        }

        // Create a new gateway for connecting to our peer node.
        const gateway = new Gateway();
        await gateway.connect(ccpPath, { wallet, identity: 'Koma', discovery: { enabled: true, asLocalhost: true } });

        // Get the network (channel) our contract is deployed to.
        const network = await gateway.getNetwork('mychannel');

        // Get the contract from the network.
        const contract = network.getContract('3A');

        // Submit the specified transaction.
        await contract.submitTransaction('grantConsent', req.body.pk_DS, req.body.pk_DC, req.body.policy, req.body.enhash, req.body.pk_enc);
        console.log('grantConsent has been submitted');
        res.send('grantConsent has been submitted');

        // Disconnect from the gateway.
        await gateway.disconnect();

    } catch (error) {
        console.error(`Failed to submit transaction grantConsent: ${error}`);
        process.exit(1);
    }
})

// app.get('/spchain/queryAll', async function (req, res) {
//     try {

//         // Create a new file system based wallet for managing identities.
//         const walletPath = path.join(process.cwd(), 'wallet');
//         const wallet = new FileSystemWallet(walletPath);
//         console.log(`Wallet path: ${walletPath}`);

//         // Check to see if we've already enrolled the user.
//         const userExists = await wallet.exists('Koma');
//         if (!userExists) {
//             console.log('An identity for the user "Koma" does not exist in the wallet');
//             console.log('Run the registerUser.js application before retrying');
// 			require('child_process').fork('registerUser.js');
//             return;
//         }

//         // Create a new gateway for connecting to our peer node.
//         const gateway = new Gateway();
//         await gateway.connect(ccpPath, { wallet, identity: 'Koma', discovery: { enabled: true, asLocalhost: true } });

//         // Get the network (channel) our contract is deployed to.
//         const network = await gateway.getNetwork('mychannel');

//         // Get the contract from the network.
//         const contract = network.getContract('artwork');

//         // Evaluate the specified transaction.
//         const result = await contract.evaluateTransaction('queryAll');
//         console.log(`Transaction has been evaluated, result is: ${result.toString()}`);
//         res.status(200).json({response: result.toString()});

//     } catch (error) {
//         console.error(`Failed to evaluate transaction: ${error}`);
//         res.status(500).json({error: error});
//         process.exit(1);
//     }
// });


app.get('/spchain/readArtwork/', async function (req, res) {
    try {

        // Create a new file system based wallet for managing identities.
        const walletPath = path.join(process.cwd(), 'wallet');
        const wallet = new FileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        // Check to see if we've already enrolled the user.
        const userExists = await wallet.exists('Koma');
        if (!userExists) {
            console.log('An identity for the user "Koma" does not exist in the wallet');
            console.log('Run the registerGalleryUser.js application before retrying');
			require('child_process').fork('registerGalleryUser.js');
            return;
        }

        // Create a new gateway for connecting to our peer node.
        const gateway = new Gateway();
        await gateway.connect(ccpPath, { wallet, identity: 'Koma', discovery: { enabled: true, asLocalhost: true } });

        // Get the network (channel) our contract is deployed to.
        const network = await gateway.getNetwork('mychannel');

        // Get the contract from the network.
        const contract = network.getContract('artworks');

        // Evaluate the specified transaction.
        const result = await contract.evaluateTransaction('readArtwork', req.body.tokenID);
        console.log(`readArtwork has been evaluated, result is: ${result.toString()}`);
        res.status(200).json({response: result.toString()});

    } catch (error) {
        console.error(`Failed to evaluate readArtwork: ${error}`);
        res.status(500).json({error: error});
        process.exit(1);
    }
});

app.post('/spchain/uploadArtwork/', async function (req, res) {
    try {

        // Create a new file system based wallet for managing identities.
        const walletPath = path.join(process.cwd(), 'wallet');
        const wallet = new FileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        // Check to see if we've already enrolled the user.
        const userExists = await wallet.exists('Koma');
        if (!userExists) {
            console.log('An identity for the user "Koma" does not exist in the wallet');
            console.log('Run the registerGalleryUser.js application before retrying');
			require('child_process').fork('registerGalleryUser.js');
            return;
        }

        // Create a new gateway for connecting to our peer node.
        const gateway = new Gateway();
        await gateway.connect(ccpPath, { wallet, identity: 'Koma', discovery: { enabled: true, asLocalhost: true } });

        // Get the network (channel) our contract is deployed to.
        const network = await gateway.getNetwork('mychannel');

        // Get the contract from the network.
        const contract = network.getContract('artworks');

        // Submit the specified transaction.
        await contract.submitTransaction('uploadArtwork', req.body.tokenID, req.body.multihash, req.body.owner, req.body.creator);
        console.log('uploadArtwork has been submitted');
        res.send('uploadArtwork has been submitted');

        // Disconnect from the gateway.
        await gateway.disconnect();

    } catch (error) {
        console.error(`Failed to submit uploadArtwork: ${error}`);
        process.exit(1);
    }
})

app.post('/spchain/transferArtwork/', async function (req, res) {
    try {

        // Create a new file system based wallet for managing identities.
        const walletPath = path.join(process.cwd(), 'wallet');
        const wallet = new FileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        // Check to see if we've already enrolled the user.
        const userExists = await wallet.exists('Koma');
        if (!userExists) {
            console.log('An identity for the user "Koma" does not exist in the wallet');
            console.log('Run the registerGalleryUser.js application before retrying');
			require('child_process').fork('registerGalleryUser.js');
            return;
        }

        // Create a new gateway for connecting to our peer node.
        const gateway = new Gateway();
        await gateway.connect(ccpPath, { wallet, identity: 'Koma', discovery: { enabled: true, asLocalhost: true } });

        // Get the network (channel) our contract is deployed to.
        const network = await gateway.getNetwork('mychannel');

        // Get the contract from the network.
        const contract = network.getContract('artworks');

        // Submit the specified transaction.
        await contract.submitTransaction('transferArtwork', req.body.tokenID, req.body.newCollector, req.body.multihash);
        console.log('transferArtwork has been submitted');
        res.send('transferArtwork has been submitted');

        // Disconnect from the gateway.
        await gateway.disconnect();

    } catch (error) {
        console.error(`Failed to submit transferArtwork: ${error}`);
        process.exit(1);
    }
})

app.post('/spchain/addLog/', async function (req, res) {
    try {

        // Create a new file system based wallet for managing identities.
        const walletPath = path.join(process.cwd(), 'wallet');
        const wallet = new FileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        // Check to see if we've already enrolled the user.
        const userExists = await wallet.exists('Koma');
        if (!userExists) {
            console.log('An identity for the user "Koma" does not exist in the wallet');
            console.log('Run the registerGalleryUser.js application before retrying');
			require('child_process').fork('registerGalleryUser.js');
            return;
        }

        // Create a new gateway for connecting to our peer node.
        const gateway = new Gateway();
        await gateway.connect(ccpPath, { wallet, identity: 'Koma', discovery: { enabled: true, asLocalhost: true } });

        // Get the network (channel) our contract is deployed to.
        const network = await gateway.getNetwork('mychannel');

        // Get the contract from the network.
        const contract = network.getContract('logs');

        // Submit the specified transaction.
        await contract.submitTransaction('addLog', req.body.pk_DS, req.body.pk_DC, req.body.pk_DP, req.body.sk_data, req.body.status, req.body.operation);
        console.log('addLog has been submitted');
        res.send('addLog has been submitted');

        // Disconnect from the gateway.
        await gateway.disconnect();

    } catch (error) {
        console.error(`Failed to submit addLog: ${error}`);
        process.exit(1);
    }
})

app.post('/spchain/transferBalance/', async function (req, res) {
    try {

        // Create a new file system based wallet for managing identities.
        const walletPath = path.join(process.cwd(), 'wallet');
        const wallet = new FileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        // Check to see if we've already enrolled the user.
        const userExists = await wallet.exists('Koma');
        if (!userExists) {
            console.log('An identity for the user "Koma" does not exist in the wallet');
            console.log('Run the registerGalleryUser.js application before retrying');
			require('child_process').fork('registerGalleryUser.js');
            return;
        }

        // Create a new gateway for connecting to our peer node.
        const gateway = new Gateway();
        await gateway.connect(ccpPath, { wallet, identity: 'Koma', discovery: { enabled: true, asLocalhost: true } });

        // Get the network (channel) our contract is deployed to.
        const network = await gateway.getNetwork('mychannel');

        // Get the contract from the network.
        const contract = network.getContract('wallets');

        // Submit the specified transaction.
        await contract.submitTransaction('transferBalance', req.body.newCollector, req.body.collector, req.body.creator, req.body.price, req.body.r_y);
        console.log('transferBalance has been submitted');
        res.send('transferBalance has been submitted');

        // Disconnect from the gateway.
        await gateway.disconnect();

    } catch (error) {
        console.error(`Failed to submit transferBalance: ${error}`);
        process.exit(1);
    }
})

// app.put('/spchain/transferArtwork/:tokenid', async function (req, res) {
//     try {

//         // Create a new file system based wallet for managing identities.
//         const walletPath = path.join(process.cwd(), 'wallet');
//         const wallet = new FileSystemWallet(walletPath);
//         console.log(`Wallet path: ${walletPath}`);

//         // Check to see if we've already enrolled the user.
//         const userExists = await wallet.exists('Koma');
//         if (!userExists) {
//             console.log('An identity for the user "Koma" does not exist in the wallet');
//             console.log('Run the registerUser.js application before retrying');
// 			require('child_process').fork('registerUser.js');
//             return;
//         }

//         // Create a new gateway for connecting to our peer node.
//         const gateway = new Gateway();
//         await gateway.connect(ccpPath, { wallet, identity: 'Koma', discovery: { enabled: true, asLocalhost: true } });

//         // Get the network (channel) our contract is deployed to.
//         const network = await gateway.getNetwork('mychannel');

//         // Get the contract from the network.
//         const contract = network.getContract('artwork');

//         // Submit the specified transaction.
//         await contract.submitTransaction('transferArtwork', req.params.tokenid, req.body.newowner);
//         console.log('Transaction has been submitted');
//         res.send('Transaction has been submitted');

//         // Disconnect from the gateway.
//         await gateway.disconnect();

//     } catch (error) {
//         console.error(`Failed to submit transaction: ${error}`);
//         process.exit(1);
//     }	
// })

// app.get('/spchain/deleteArtwork/:tokenid', async function (req, res) {
//     try {

//         // Create a new file system based wallet for managing identities.
//         const walletPath = path.join(process.cwd(), 'wallet');
//         const wallet = new FileSystemWallet(walletPath);
//         console.log(`Wallet path: ${walletPath}`);

//         // Check to see if we've already enrolled the user.
//         const userExists = await wallet.exists('Koma');
//         if (!userExists) {
//             console.log('An identity for the user "Koma" does not exist in the wallet');
//             console.log('Run the registerUser.js application before retrying');
// 			require('child_process').fork('registerUser.js');
//             return;
//         }

//         // Create a new gateway for connecting to our peer node.
//         const gateway = new Gateway();
//         await gateway.connect(ccpPath, { wallet, identity: 'Koma', discovery: { enabled: true, asLocalhost: true } });

//         // Get the network (channel) our contract is deployed to.
//         const network = await gateway.getNetwork('mychannel');

//         // Get the contract from the network.
//         const contract = network.getContract('artwork');

//         // Evaluate the specified transaction.
//         const result = await contract.evaluateTransaction('deleteArtwork', req.params.tokenid);
//         console.log(`Transaction has been evaluated, result is: ${result.toString()}`);
//         res.status(200).json({response: result.toString()});

//     } catch (error) {
//         console.error(`Failed to evaluate transaction: ${error}`);
//         res.status(500).json({error: error});
//         process.exit(1);
//     }
// });

// app.get('/spchain/getHistoryForArtwork/:tokenid', async function (req, res) {
//     try {

//         // Create a new file system based wallet for managing identities.
//         const walletPath = path.join(process.cwd(), 'wallet');
//         const wallet = new FileSystemWallet(walletPath);
//         console.log(`Wallet path: ${walletPath}`);

//         // Check to see if we've already enrolled the user.
//         const userExists = await wallet.exists('Koma');
//         if (!userExists) {
//             console.log('An identity for the user "Koma" does not exist in the wallet');
//             console.log('Run the registerUser.js application before retrying');
// 			require('child_process').fork('registerUser.js');
//             return;
//         }

//         // Create a new gateway for connecting to our peer node.
//         const gateway = new Gateway();
//         await gateway.connect(ccpPath, { wallet, identity: 'Koma', discovery: { enabled: true, asLocalhost: true } });

//         // Get the network (channel) our contract is deployed to.
//         const network = await gateway.getNetwork('mychannel');

//         // Get the contract from the network.
//         const contract = network.getContract('artwork');

//         // Evaluate the specified transaction.
//         const result = await contract.evaluateTransaction('getHistoryForArtwork', req.params.tokenid);
//         console.log(`Transaction has been evaluated, result is: ${result.toString()}`);
//         res.status(200).json({response: result.toString()});

//     } catch (error) {
//         console.error(`Failed to evaluate transaction: ${error}`);
//         res.status(500).json({error: error});
//         process.exit(1);
//     }
// });

app.listen(8080);
