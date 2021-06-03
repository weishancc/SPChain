'use strict';

module.exports.info = 'transferArtwork callback';

const helper = require('./helper');
const contractID = 'artworks';
const version = '1.0';
const tokenID = 'afee1720-28ef-4790-8fe8-50f38681b919';
const multiHash = 'zdpuAqazmCjGhntRXpsKpcsUou9D72TzowCTgJK5jC62CoYvt';
const newOwner = 'Amber';

let bc, ctx, limitIndex;
let txIndex = 0;

/**
* Initializes the workload module before the start of the round.
* @param {BlockchainInterface} blockchain The SUT adapter instance.
* @param {object} context The SUT-specific context for the round.
* @param {object} args The user-provided arguments for the workload module.
*/
module.exports.init = async function (blockchain, context, args) {
    bc = blockchain;
    ctx = context;
    limitIndex = args.assets;
    txIndex++;

    await helper.createArtwork(bc, ctx, txIndex);

    return Promise.resolve();
}

module.exports.run = function () {
    try {
         const myArgs = {
            chaincodeFunction: 'transferArtwork',
            invokerIdentity: 'Admin@gallery.spchain.com',
            chaincodeArguments: [tokenID + ctx.clientIdx + txIndex.toString(), newOwner, multiHash]
        };
        return bc.bcObj.invokeSmartContract(ctx, contractID, version, myArgs);
    } catch (error) {
        console.log(`Client ${ctx.clientIdx}: Smart Contract threw with error: ${error}`);
    };
};

module.exports.end = async function () {
    return Promise.resolve();
};
