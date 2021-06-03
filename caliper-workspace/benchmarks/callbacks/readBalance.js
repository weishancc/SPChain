'use strict';

module.exports.info = 'readBalance callback';

const helper = require('./helper');
const contractID = 'wallets';
const version = '1.0';

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

    // We have already added creator Koma's wallet in addWallet.js, thus just adding collector here
    await helper.createWallet(bc, ctx, txIndex, 'Amber');

    return Promise.resolve();
}

module.exports.run = function () {

    try {
        if (txIndex == limitIndex) {
            txIndex = 0;
        }

        const myArgs = {
            chaincodeFunction: 'readBalance',
            invokerIdentity: 'Admin@gallery.spchain.com',
            chaincodeArguments: ['Amber' + ctx.clientIdx + txIndex.toString()]
        };
        return bc.bcObj.querySmartContract(ctx, contractID, version, myArgs);
    } catch (error) {
        console.log(`Client ${ctx.clientIdx}: Smart Contract threw with error: ${error}`);
    };
};

module.exports.end = async function () {
    return Promise.resolve();
};
